package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware sets a deadline on the request context.
// All downstream database queries and downstream HTTP calls that respect
// c.Request.Context() will be automatically cancelled when the timeout expires,
// preventing hanging database connections and resource starvation.
//
// WebSocket upgrade requests are automatically excluded from the timeout.
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip timeout for WebSocket upgrades (long-lived persistent connections)
		if c.GetHeader("Upgrade") == "websocket" {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		// If the context timed out during execution and no response was written yet,
		// return a 504 Gateway Timeout.
		if errors.Is(ctx.Err(), context.DeadlineExceeded) && !c.Writer.Written() {
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"code":  "REQUEST_TIMEOUT",
				"error": "Il server ha impiegato troppo tempo a rispondere. Riprova più tardi.",
			})
		}
	}
}
