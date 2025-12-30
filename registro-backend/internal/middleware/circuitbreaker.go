package middleware

import (
	"fmt"
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
	settings.Timeout = 5 * time.Second

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
