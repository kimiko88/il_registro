package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestBodyLimitMiddleware limits the HTTP request body size to maxBytes (default 2 MB)
// to protect against Out-Of-Memory Denial of Service attacks.
// Multipart file uploads are exempted from this limit, as they are independently
// guarded by upload.MaxUploadSize.
func RequestBodyLimitMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if strings.HasPrefix(contentType, "multipart/form-data") {
			c.Next()
			return
		}

		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}
