package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestBodyLimitMiddleware limits the HTTP request body size to maxBytes (default 2 MB)
// to protect against Out-Of-Memory Denial of Service attacks.
// Multipart uploads receive a separate bounded 22 MiB request ceiling to allow
// the 20 MiB file limit plus multipart framing without permitting unbounded bodies.
func RequestBodyLimitMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		limit := maxBytes
		if strings.HasPrefix(strings.ToLower(contentType), "multipart/form-data") {
			// Allow the documented 20 MiB upload ceiling plus bounded multipart headers.
			// Never exempt multipart requests: parsers may otherwise spool unbounded data to disk.
			limit = 22 << 20
		}
		if c.Request.ContentLength > limit {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		}
		c.Next()
	}
}
