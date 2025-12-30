package middleware

import (
	"time"

	"registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if logger.Log != nil {
			logger.Log.Info("request",
				"method", c.Request.Method,
				"path", path,
				"query", raw,
				"status", c.Writer.Status(),
				"latency", time.Since(start),
				"client_ip", c.ClientIP(),
			)
		}
	}
}
