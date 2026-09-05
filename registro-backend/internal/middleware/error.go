package middleware

import (
	"registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware catches unhandled Gin errors and logs them with structured logger.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				if logger.Log != nil {
					logger.Log.Errorf("[MiddlewareError] %s: %v", c.Request.URL.Path, e.Err)
				}
			}
		}
	}
}

