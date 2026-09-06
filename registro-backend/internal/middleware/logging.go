package middleware

import (
	"time"

	"registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if logger.Log != nil {
			// Read auth context set by the JWT middleware (may be nil for public routes).
			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")
			schoolID, _ := c.Get("school_id")

			fields := logrus.Fields{
				"method":     c.Request.Method,
				"path":       path,
				"status":     c.Writer.Status(),
				"latency":    time.Since(start),
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}
			// Only append identity fields when present to keep public-route logs clean.
			if userID != nil && userID != "" {
				fields["user_id"] = userID
				fields["role"] = role
				fields["school_id"] = schoolID
			}
			if raw != "" {
				fields["query"] = raw
			}

			logger.Log.WithFields(fields).Info("request")
		}
	}
}
