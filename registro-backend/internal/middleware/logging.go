package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"registro-backend/internal/metrics"
	"registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("req-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := strings.TrimSpace(c.Request.Header.Get("X-Request-ID"))
		if reqID == "" {
			reqID = generateRequestID()
		}
		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		duration := time.Since(start)
		metrics.DefaultRegistry.RecordRequest(c.Request.Method, path, c.Writer.Status(), duration)

		if logger.Log != nil {
			// Read auth context set by the JWT middleware (may be nil for public routes).
			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")
			schoolID, _ := c.Get("school_id")

			fields := logrus.Fields{
				"request_id": reqID,
				"method":     c.Request.Method,
				"path":       path,
				"status":     c.Writer.Status(),
				"latency":    duration,
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
