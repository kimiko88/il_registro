package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	idempotencyTTL         = 10 * time.Minute
	idempotencyHeader      = "Idempotency-Key"
	maxIdempotencyBody     = 2 << 20
	maxIdempotencyResponse = 64 << 10
	maxIdempotencyEntries  = 1000
)

type cachedResponse struct {
	status     int
	body       []byte
	headers    map[string]string
	storedAt   time.Time
	replayable bool
}

type idempotencyEntry struct {
	requestHash string
	done        chan struct{}
	response    *cachedResponse
}

var (
	idempotencyMu    sync.Mutex
	idempotencyCache = make(map[string]*idempotencyEntry)
)

func idempotencyScopeKey(userID, schoolID, method, path, clientKey string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{userID, schoolID, method, path, clientKey}, "\x00")))
	return hex.EncodeToString(sum[:])
}

// IdempotencyMiddleware prevents duplicate mutations and prevents response
// replay across users, schools, routes or payloads. Entries are process-local:
// multi-replica deployments still need a shared Redis/database idempotency store.
func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch {
			c.Next()
			return
		}

		clientKey := strings.TrimSpace(c.GetHeader(idempotencyHeader))
		if clientKey == "" {
			c.Next()
			return
		}
		if len(clientKey) > 128 || strings.IndexFunc(clientKey, func(r rune) bool { return r < 0x21 || r == 0x7f }) >= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key non valido"})
			return
		}

		userID := c.GetString("user_id")
		schoolID := c.GetString("school_id")
		if userID == "" || schoolID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "identità non disponibile per l'idempotenza"})
			return
		}

		// Hash the complete bounded request body so the same key cannot be
		// reused with a different mutation. Large uploads bypass this middleware
		// cache rather than buffering an unbounded payload.
		body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxIdempotencyBody+1))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "impossibile leggere il corpo della richiesta"})
			return
		}
		c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body), c.Request.Body))
		if len(body) > maxIdempotencyBody {
			c.Next()
			return
		}
		bodyHash := sha256.Sum256(body)
		requestHash := hex.EncodeToString(bodyHash[:])
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
			path += "?" + rawQuery
		}
		cacheKey := idempotencyScopeKey(userID, schoolID, method, path, clientKey)

		var entry *idempotencyEntry
		owner := false
		for {
			now := time.Now()
			idempotencyMu.Lock()
			for key, candidate := range idempotencyCache {
				if candidate.response != nil && now.Sub(candidate.response.storedAt) > idempotencyTTL {
					delete(idempotencyCache, key)
				}
			}

			existing := idempotencyCache[cacheKey]
			if existing == nil {
				if len(idempotencyCache) >= maxIdempotencyEntries {
					idempotencyMu.Unlock()
					c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "servizio temporaneamente occupato: cache di idempotenza piena"})
					return
				}
				entry = &idempotencyEntry{requestHash: requestHash, done: make(chan struct{})}
				idempotencyCache[cacheKey] = entry
				owner = true
				idempotencyMu.Unlock()
				break
			}
			if existing.requestHash != requestHash {
				idempotencyMu.Unlock()
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Idempotency-Key già utilizzata con un payload differente"})
				return
			}
			done := existing.done
			idempotencyMu.Unlock()

			select {
			case <-done:
			case <-c.Request.Context().Done():
				c.Abort()
				return
			}

			idempotencyMu.Lock()
			response := existing.response
			if response == nil {
				if idempotencyCache[cacheKey] == existing {
					delete(idempotencyCache, cacheKey)
				}
				idempotencyMu.Unlock()
				continue
			}
			if time.Since(response.storedAt) > idempotencyTTL {
				if idempotencyCache[cacheKey] == existing {
					delete(idempotencyCache, cacheKey)
				}
				idempotencyMu.Unlock()
				continue
			}
			copyResponse := cachedResponse{
				status:     response.status,
				body:       append([]byte(nil), response.body...),
				headers:    make(map[string]string, len(response.headers)),
				storedAt:   response.storedAt,
				replayable: response.replayable,
			}
			for key, value := range response.headers {
				copyResponse.headers[key] = value
			}
			idempotencyMu.Unlock()

			c.Header("X-Idempotency-Replayed", "true")
			if !copyResponse.replayable {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"error": "operazione già completata; la risposta supera il limite di replay, ricaricare la risorsa prima di riprovare",
				})
				return
			}
			for key, value := range copyResponse.headers {
				c.Header(key, value)
			}
			contentType := copyResponse.headers["Content-Type"]
			if contentType == "" {
				contentType = "application/json"
			}
			c.Data(copyResponse.status, contentType, copyResponse.body)
			c.Abort()
			return
		}

		if !owner {
			return
		}
		writer := &responseWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = writer
		completed := false
		defer func() {
			if completed {
				return
			}
			idempotencyMu.Lock()
			if idempotencyCache[cacheKey] == entry {
				delete(idempotencyCache, cacheKey)
			}
			close(entry.done)
			idempotencyMu.Unlock()
		}()

		c.Next()
		completed = true

		status := writer.status
		if status == 0 {
			status = http.StatusOK
		}
		var response *cachedResponse
		if status >= 200 && status < 500 {
			headers := make(map[string]string)
			for key, values := range writer.Header() {
				if len(values) > 0 {
					headers[key] = values[0]
				}
			}
			body := writer.body.Bytes()
			replayable := len(body) <= maxIdempotencyResponse
			if !replayable {
				body = nil
			}
			response = &cachedResponse{
				status:     status,
				body:       append([]byte(nil), body...),
				headers:    headers,
				storedAt:   time.Now(),
				replayable: replayable,
			}
		}

		idempotencyMu.Lock()
		entry.response = response
		if response == nil && idempotencyCache[cacheKey] == entry {
			delete(idempotencyCache, cacheKey)
		}
		close(entry.done)
		idempotencyMu.Unlock()
	}
}

// responseWriter captures the response body without changing normal Gin writes.
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
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	_, _ = rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteString(s string) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	_, _ = rw.body.WriteString(s)
	return rw.ResponseWriter.WriteString(s)
}
