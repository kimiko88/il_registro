package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── IdempotencyMiddleware ────────────────────────────────────────────────────

func TestIdempotencyMiddleware_ReplayOnSecondRequest(t *testing.T) {
	r := gin.New()
	r.Use(IdempotencyMiddleware())

	var callCount int
	r.POST("/idempotent", func(c *gin.Context) {
		callCount++
		c.JSON(http.StatusCreated, gin.H{"count": callCount})
	})

	// Use a unique key per test to avoid cache collisions with other tests.
	key := "idempotency-replay-test-" + t.Name()

	makeReq := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/idempotent", strings.NewReader("{}"))
		req.Header.Set("Idempotency-Key", key)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		return w
	}

	w1 := makeReq()
	assert.Equal(t, http.StatusCreated, w1.Code)

	w2 := makeReq()
	assert.Equal(t, http.StatusCreated, w2.Code)
	assert.Equal(t, "true", w2.Header().Get("X-Idempotency-Replayed"),
		"second request should be marked as replayed")

	assert.Equal(t, 1, callCount, "handler must only be called once")
}

func TestIdempotencyMiddleware_PassthroughOnGET(t *testing.T) {
	r := gin.New()
	r.Use(IdempotencyMiddleware())
	r.GET("/resource", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/resource", nil)
	req.Header.Set("Idempotency-Key", "some-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("X-Idempotency-Replayed"),
		"GET requests must not be cached")
}

func TestIdempotencyMiddleware_PassthroughWithNoKey(t *testing.T) {
	r := gin.New()
	r.Use(IdempotencyMiddleware())
	r.POST("/data", func(c *gin.Context) {
		c.String(http.StatusOK, "handled")
	})

	req, _ := http.NewRequest(http.MethodPost, "/data", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("X-Idempotency-Replayed"))
}

func TestIdempotencyMiddleware_PassthroughOnDELETE(t *testing.T) {
	r := gin.New()
	r.Use(IdempotencyMiddleware())
	r.DELETE("/item", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/item", nil)
	req.Header.Set("Idempotency-Key", "delete-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Header().Get("X-Idempotency-Replayed"),
		"DELETE is not a mutating method covered by idempotency middleware")
}

// ─── IPRateLimiter ────────────────────────────────────────────────────────────

func TestIPRateLimiter_GetLimiterSameIP(t *testing.T) {
	rl := NewIPRateLimiter(5, 10)
	defer rl.Close()

	l1 := rl.GetLimiter("1.2.3.4")
	l2 := rl.GetLimiter("1.2.3.4")
	assert.Equal(t, l1, l2, "same IP must return the same limiter instance")
}

func TestIPRateLimiter_GetLimiterDifferentIPs(t *testing.T) {
	rl := NewIPRateLimiter(5, 10)
	defer rl.Close()

	l1 := rl.GetLimiter("1.1.1.1")
	l2 := rl.GetLimiter("2.2.2.2")
	assert.NotSame(t, l1, l2, "different IPs must get independent limiters")
}

func TestIPRateLimiter_Close(t *testing.T) {
	rl := NewIPRateLimiter(5, 10)
	assert.NotPanics(t, func() {
		rl.Close()
		rl.Close() // double-close must be safe
	})
}

// ─── SecurityHeadersMiddleware ────────────────────────────────────────────────

func TestSecurityHeadersMiddleware_BaseHeaders(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	h := w.Header()
	assert.Equal(t, "DENY", h.Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", h.Get("X-Content-Type-Options"))
	assert.Equal(t, "strict-origin-when-cross-origin", h.Get("Referrer-Policy"))
	assert.Equal(t, "1; mode=block", h.Get("X-XSS-Protection"))
	assert.NotEmpty(t, h.Get("Permissions-Policy"))
	assert.NotEmpty(t, h.Get("Content-Security-Policy"))
}

func TestSecurityHeadersMiddleware_CSPNonce(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())

	var capturedNonce interface{}
	r.GET("/nonce", func(c *gin.Context) {
		capturedNonce, _ = c.Get("csp_nonce")
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/nonce", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	nonce, _ := capturedNonce.(string)
	assert.NotEmpty(t, nonce, "CSP nonce must be set in Gin context")
	csp := w.Header().Get("Content-Security-Policy")
	assert.Contains(t, csp, nonce, "CSP header must contain the nonce")
}

func TestSecurityHeadersMiddleware_HSTSOnHTTPS(t *testing.T) {
	t.Setenv("TRUST_PROXY_HEADERS", "true")
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/secure", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	hsts := w.Header().Get("Strict-Transport-Security")
	assert.Contains(t, hsts, "max-age=63072000")
	assert.Contains(t, hsts, "includeSubDomains")
}

func TestSecurityHeadersMiddleware_HSTSNotSetOverHTTP(t *testing.T) {
	t.Setenv("TRUST_PROXY_HEADERS", "false")
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/plain", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/plain", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Empty(t, w.Header().Get("Strict-Transport-Security"),
		"HSTS must not be set for plain HTTP requests")
}

// ─── sanitizeCSPOrigin ────────────────────────────────────────────────────────

func TestSanitizeCSPOrigin(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal", "https://example.com", "https://example.com"},
		{"crlf injection", "https://bad\r\nexample.com", "https://badexample.com"},
		{"semicolon injection", "https://example.com;evil", "https://example.comevil"},
		{"single quote", "https://example.com'script'", "https://example.comscript"},
		{"double quote", `https://example.com"test"`, "https://example.comtest"},
		{"whitespace trim", "  https://trimmed.com  ", "https://trimmed.com"},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := sanitizeCSPOrigin(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

// ─── RequestBodyLimitMiddleware ───────────────────────────────────────────────

func TestRequestBodyLimitMiddleware_AllowsSmallBody(t *testing.T) {
	r := gin.New()
	r.Use(RequestBodyLimitMiddleware(1024)) // 1 KB
	r.POST("/data", func(c *gin.Context) { c.Status(http.StatusOK) })

	body := strings.NewReader(`{"key":"small"}`)
	req, _ := http.NewRequest(http.MethodPost, "/data", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestBodyLimitMiddleware_SkipsMultipart(t *testing.T) {
	r := gin.New()
	r.Use(RequestBodyLimitMiddleware(1)) // 1 byte limit
	r.POST("/upload", func(c *gin.Context) { c.Status(http.StatusOK) })

	// multipart/form-data should bypass the limit
	req, _ := http.NewRequest(http.MethodPost, "/upload",
		strings.NewReader("huge multipart data..."))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=xxx")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestBodyLimitMiddleware_NilBody(t *testing.T) {
	r := gin.New()
	r.Use(RequestBodyLimitMiddleware(1024))
	r.GET("/get", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/get", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── ETagMiddleware ───────────────────────────────────────────────────────────

func TestETagMiddleware_SetsETagOnGET(t *testing.T) {
	r := gin.New()
	r.Use(ETagMiddleware())
	r.GET("/resource", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	req, _ := http.NewRequest(http.MethodGet, "/resource", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	etag := w.Header().Get("ETag")
	require.NotEmpty(t, etag, "ETag must be set on GET response")
	assert.Equal(t, "hello world", w.Body.String())
}

func TestETagMiddleware_304OnMatchingIfNoneMatch(t *testing.T) {
	r := gin.New()
	r.Use(ETagMiddleware())
	r.GET("/resource2", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	// First request to get the ETag
	req1, _ := http.NewRequest(http.MethodGet, "/resource2", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	etag := w1.Header().Get("ETag")
	require.NotEmpty(t, etag)

	// Second request with matching If-None-Match → 304
	req2, _ := http.NewRequest(http.MethodGet, "/resource2", nil)
	req2.Header.Set("If-None-Match", etag)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotModified, w2.Code)
	assert.Empty(t, w2.Body.String(), "304 response body must be empty")
}

func TestETagMiddleware_PassthroughOnPOST(t *testing.T) {
	r := gin.New()
	r.Use(ETagMiddleware())
	r.POST("/data", func(c *gin.Context) {
		c.String(http.StatusCreated, "created")
	})

	req, _ := http.NewRequest(http.MethodPost, "/data", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Empty(t, w.Header().Get("ETag"), "ETag must not be set for POST")
}

func TestETagMiddleware_NonOKStatusNoETag(t *testing.T) {
	r := gin.New()
	r.Use(ETagMiddleware())
	r.GET("/notfound", func(c *gin.Context) {
		c.String(http.StatusNotFound, "not found")
	})

	req, _ := http.NewRequest(http.MethodGet, "/notfound", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Empty(t, w.Header().Get("ETag"), "ETag must not be set for non-200 responses")
}

func TestETagMiddleware_WildcardIfNoneMatch(t *testing.T) {
	r := gin.New()
	r.Use(ETagMiddleware())
	r.GET("/wildcard", func(c *gin.Context) {
		c.String(http.StatusOK, "content")
	})

	req, _ := http.NewRequest(http.MethodGet, "/wildcard", nil)
	req.Header.Set("If-None-Match", "*")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotModified, w.Code)
}

// ─── RateLimitMiddleware (global) ─────────────────────────────────────────────

func TestRateLimitMiddleware_AllowsFirstRequest(t *testing.T) {
	r := gin.New()
	r.Use(RateLimitMiddleware())
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.200.200:1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "120", w.Header().Get("RateLimit-Limit"))
}

func TestAuthRateLimitMiddleware_AllowsFirstRequest(t *testing.T) {
	r := gin.New()
	r.Use(AuthRateLimitMiddleware())
	r.POST("/login", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "10.5.5.5:9999"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "10", w.Header().Get("RateLimit-Limit"))
}

// ─── TimeoutMiddleware ────────────────────────────────────────────────────────

func TestTimeoutMiddleware_CompletesFast(t *testing.T) {
	r := gin.New()
	r.Use(TimeoutMiddleware(500 * time.Millisecond))
	r.GET("/fast", func(c *gin.Context) {
		c.String(http.StatusOK, "done")
	})

	req, _ := http.NewRequest(http.MethodGet, "/fast", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "done", w.Body.String())
}

func TestTimeoutMiddleware_WebSocketSkipsTimeout(t *testing.T) {
	r := gin.New()
	r.Use(TimeoutMiddleware(10 * time.Millisecond))
	r.GET("/ws", func(c *gin.Context) {
		// Simulate slow WS handler
		time.Sleep(50 * time.Millisecond)
		c.Status(http.StatusSwitchingProtocols)
	})

	req, _ := http.NewRequest(http.MethodGet, "/ws", nil)
	req.Header.Set("Upgrade", "websocket")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// WebSocket requests bypass the timeout, so they get the handler's response
	assert.Equal(t, http.StatusSwitchingProtocols, w.Code)
}

// ─── backendURL / backendWSURL helpers ────────────────────────────────────────

func TestBackendURL_DefaultDev(t *testing.T) {
	t.Setenv("BACKEND_URL", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("GIN_MODE", "")
	url := backendURL()
	assert.Equal(t, "http://localhost:8080", url)
}

func TestBackendURL_FromEnv(t *testing.T) {
	t.Setenv("BACKEND_URL", "https://api.school.it")
	url := backendURL()
	assert.Equal(t, "https://api.school.it", url)
}

func TestBackendWSURL_DerivedFromHTTPS(t *testing.T) {
	t.Setenv("BACKEND_URL", "https://api.school.it")
	t.Setenv("BACKEND_WS_URL", "")
	ws := backendWSURL()
	assert.Equal(t, "wss://api.school.it", ws)
}

func TestBackendWSURL_DerivedFromHTTP(t *testing.T) {
	t.Setenv("BACKEND_URL", "http://localhost:8080")
	t.Setenv("BACKEND_WS_URL", "")
	ws := backendWSURL()
	assert.Equal(t, "ws://localhost:8080", ws)
}

func TestBackendWSURL_FromEnv(t *testing.T) {
	t.Setenv("BACKEND_WS_URL", "wss://ws.school.it")
	ws := backendWSURL()
	assert.Equal(t, "wss://ws.school.it", ws)
}
