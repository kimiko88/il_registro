package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const idempotencyTTL = 10 * time.Minute
const idempotencyHeader = "Idempotency-Key"

type cachedResponse struct {
	status   int
	body     []byte
	headers  map[string]string
	storedAt time.Time
}

var (
	idempotencyCache sync.Map
)

// IdempotencyMiddleware prevents duplicate mutations when a client retries a
// request whose response was lost in transit (e.g. during an offline-sync replay).
//
// On first request with a given Idempotency-Key the response is cached for
// idempotencyTTL (10 min). Subsequent requests with the same key receive the
// cached response immediately without re-executing the handler.
//
// Only applies to mutating methods: POST, PUT, PATCH.
// Keys without the header pass through unchanged.
func IdempotencyMiddleware() gin.HandlerFunc {
	// Background goroutine to evict expired entries every 5 minutes.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			idempotencyCache.Range(func(key, value any) bool {
				if cached, ok := value.(cachedResponse); ok {
					if now.Sub(cached.storedAt) > idempotencyTTL {
						idempotencyCache.Delete(key)
					}
				}
				return true
			})
		}
	}()

	return func(c *gin.Context) {
		method := c.Request.Method
		if method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch {
			c.Next()
			return
		}

		key := c.GetHeader(idempotencyHeader)
		if key == "" {
			c.Next()
			return
		}

		// Check if we have a cached response for this key.
		if cached, ok := idempotencyCache.Load(key); ok {
			entry := cached.(cachedResponse)
			// Re-play original response headers
			for k, v := range entry.headers {
				c.Header(k, v)
			}
			c.Header("X-Idempotency-Replayed", "true")
			c.Data(entry.status, "application/json", entry.body)
			c.Abort()
			return
		}

		// Intercept the response writer so we can cache it.
		rw := &responseWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = rw

		c.Next()

		// Cache the response.
		headers := make(map[string]string)
		for k, vals := range rw.Header() {
			if len(vals) > 0 {
				headers[k] = vals[0]
			}
		}
		idempotencyCache.Store(key, cachedResponse{
			status:   rw.status,
			body:     rw.body.Bytes(),
			headers:  headers,
			storedAt: time.Now(),
		})
	}
}

// responseWriter wraps gin.ResponseWriter to capture body and status code.
type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteString(s string) (int, error) {
	rw.body.WriteString(s)
	return rw.ResponseWriter.WriteString(s)
}
