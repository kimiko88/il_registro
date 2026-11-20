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
				if strings.EqualFold(key, "Content-Type") || strings.EqualFold(key, "Content-Length") {
					continue
				}
				c.Header(key, value)
			}
			c.Header("X-Content-Type-Options", "nosniff")
			contentType := safeIdempotencyContentType(copyResponse.headers["Content-Type"])
			c.Writer.Header().Set("Content-Type", contentType)
			c.Data(copyResponse.status, contentType, copyResponse.body)
			c.Abort()
			return
		}

		if !owner {
			return
		}
		origWriter := c.Writer
		writer := &responseWriter{ResponseWriter: origWriter, body: &bytes.Buffer{}}
		c.Writer = writer
		completed := false
		defer func() {
			c.Writer = origWriter
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
		c.Writer = origWriter

		status := writer.status
		if status == 0 {
			status = http.StatusOK
		}
		var response *cachedResponse
		var bodyBytes []byte
		if !writer.passthrough && !writer.exceeded {
			bodyBytes = append([]byte(nil), writer.body.Bytes()...)
		}
		if status >= 200 && status < 500 && !writer.exceeded && len(bodyBytes) <= maxIdempotencyResponse {
			headers := make(map[string]string)
			for key, values := range origWriter.Header() {
				if len(values) > 0 {
					headers[key] = values[0]
				}
			}
			response = &cachedResponse{
				status:     status,
				body:       append([]byte(nil), bodyBytes...),
				headers:    headers,
				storedAt:   time.Now(),
				replayable: true,
			}
		} else if status >= 200 && status < 500 {
			headers := make(map[string]string)
			for key, values := range origWriter.Header() {
				if len(values) > 0 {
					headers[key] = values[0]
				}
			}
			response = &cachedResponse{
				status:     status,
				headers:    headers,
				storedAt:   time.Now(),
				replayable: false,
			}
		}

		if !writer.passthrough {
			c.Header("X-Content-Type-Options", "nosniff")
			safeCT := safeIdempotencyContentType(origWriter.Header().Get("Content-Type"))
			origWriter.Header().Set("Content-Type", safeCT)
			if len(bodyBytes) == 0 && (status == http.StatusNoContent || status == http.StatusNotModified) {
				origWriter.WriteHeader(status)
			} else if len(bodyBytes) > 0 {
				c.Data(status, safeCT, bodyBytes)
			} else if writer.wroteHeader {
				origWriter.WriteHeader(status)
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

func safeIdempotencyContentType(raw string) string {
	mediaType := strings.ToLower(strings.TrimSpace(strings.SplitN(raw, ";", 2)[0]))
	switch mediaType {
	case "", "text/html", "application/xhtml+xml", "image/svg+xml", "text/xml", "application/xml":
		return "application/json; charset=utf-8"
	default:
		return raw
	}
}

// responseWriter buffers cacheable responses up to maxIdempotencyResponse and falls back
// to pass-through streaming for larger payloads, enforcing safe content types.
type responseWriter struct {
	gin.ResponseWriter
	body        *bytes.Buffer
	status      int
	wroteHeader bool
	passthrough bool
	exceeded    bool
}

func (rw *responseWriter) Status() int {
	if rw.status != 0 {
		return rw.status
	}
	return http.StatusOK
}

func (rw *responseWriter) Size() int {
	if rw.passthrough {
		return rw.ResponseWriter.Size()
	}
	return rw.body.Len()
}

func (rw *responseWriter) Written() bool {
	if rw.passthrough {
		return rw.ResponseWriter.Written()
	}
	return rw.wroteHeader || rw.body.Len() > 0
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.wroteHeader = true
	if rw.passthrough {
		rw.ResponseWriter.WriteHeader(status)
	}
}

func (rw *responseWriter) WriteHeaderNow() {
	if !rw.passthrough {
		rw.switchToPassthrough()
	}
	rw.ResponseWriter.WriteHeaderNow()
}

func (rw *responseWriter) Flush() {
	if !rw.passthrough {
		rw.switchToPassthrough()
	}
	rw.ResponseWriter.Flush()
}

func (rw *responseWriter) switchToPassthrough() {
	if rw.passthrough {
		return
	}
	rw.passthrough = true
	hdr := rw.ResponseWriter.Header()
	hdr.Set("X-Content-Type-Options", "nosniff")
	hdr.Set("Content-Type", safeIdempotencyContentType(hdr.Get("Content-Type")))
	if rw.wroteHeader && rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}
	if rw.body.Len() > 0 {
		buffered := append([]byte(nil), rw.body.Bytes()...)
		rw.body.Reset()
		_, _ = rw.ResponseWriter.Write(buffered)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
		rw.wroteHeader = true
	}
	if rw.passthrough {
		return rw.ResponseWriter.Write(b)
	}
	if rw.body.Len()+len(b) > maxIdempotencyResponse {
		rw.exceeded = true
		rw.switchToPassthrough()
		return rw.ResponseWriter.Write(b)
	}
	return rw.body.Write(b)
}

func (rw *responseWriter) WriteString(s string) (int, error) {
	return rw.Write([]byte(s))
}
