package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

var CB *gobreaker.CircuitBreaker

func InitCircuitBreaker() {
	var settings gobreaker.Settings
	settings.Name = "Supabase"
	settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}
	// 30s before attempting half-open (was 5s — too aggressive, caused flapping).
	settings.Timeout = 30 * time.Second
	// Allow at most 2 trial requests in half-open state before declaring fully open/closed.
	settings.MaxRequests = 2
	settings.OnStateChange = func(name string, from, to gobreaker.State) {
		log.Printf("[CircuitBreaker] %s: %s → %s", name, from.String(), to.String())
	}

	CB = gobreaker.NewCircuitBreaker(settings)
}

// CircuitBreakerMiddleware wraps handlers with a circuit breaker logic
// Note: This is simpler; for granular per-service CB, apply at Service level
func CircuitBreakerMiddleware() gin.HandlerFunc {
	if CB == nil {
		InitCircuitBreaker()
	}

	return func(c *gin.Context) {
		_, err := CB.Execute(func() (interface{}, error) {
			c.Next()
			if len(c.Errors) > 0 {
				return nil, fmt.Errorf("request failed")
			}
			if c.Writer.Status() >= 500 {
				return nil, fmt.Errorf("server error")
			}
			return nil, nil
		})

		if err != nil && err.Error() == "circuit breaker is open" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service temporarily unavailable"})
			c.Abort()
		}
	}
}
