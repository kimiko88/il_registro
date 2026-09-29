package middleware

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

// ETagMiddleware implements conditional HTTP caching (RFC 7232).
// For GET requests, it hashes the response body and sends an ETag header.
// If the client sends a matching If-None-Match header, the server returns
// 304 Not Modified with zero body payload, saving bandwidth and rendering time.
func ETagMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Next()
			return
		}

		blw := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		if c.Writer.Status() != http.StatusOK {
			_, _ = blw.ResponseWriter.Write(blw.body.Bytes())
			return
		}

		bodyBytes := blw.body.Bytes()
		h := sha1.New()
		h.Write(bodyBytes)
		etag := `"` + hex.EncodeToString(h.Sum(nil)) + `"`

		c.Header("ETag", etag)
		c.Header("Cache-Control", "private, no-cache")

		clientETag := c.GetHeader("If-None-Match")
		if clientETag != "" && (clientETag == etag || clientETag == "W/"+etag || clientETag == "*") {
			c.Status(http.StatusNotModified)
			c.Writer.WriteHeaderNow()
			return
		}

		_, _ = blw.ResponseWriter.Write(bodyBytes)
	}
}
