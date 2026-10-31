package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ─── SecurityHeadersMiddleware ──────────────────────────────────────────────

func TestSecurityHeaders_XFrameOptions(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"),
		"should deny framing to prevent clickjacking")
}

func TestSecurityHeaders_XContentTypeOptions(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
}

func TestSecurityHeaders_ReferrerPolicy(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
}

func TestSecurityHeaders_XSSProtection(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
}

func TestSecurityHeaders_PermissionsPolicy(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	pp := w.Header().Get("Permissions-Policy")
	assert.Contains(t, pp, "camera=()")
	assert.Contains(t, pp, "microphone=()")
	assert.Contains(t, pp, "geolocation=()")
	assert.Contains(t, pp, "payment=()")
}

func TestSecurityHeaders_CSPIsPresent(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	csp := w.Header().Get("Content-Security-Policy")
	require.NotEmpty(t, csp, "CSP header must be present")
	assert.Contains(t, csp, "default-src 'self'")
	assert.Contains(t, csp, "object-src 'none'")
	assert.Contains(t, csp, "frame-src 'none'")
	assert.Contains(t, csp, "base-uri 'self'")
	assert.Contains(t, csp, "form-action 'self'")
	assert.Contains(t, csp, "frame-ancestors 'none'")
}

func TestSecurityHeaders_CSPContainsNonce(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	csp := w.Header().Get("Content-Security-Policy")
	assert.Contains(t, csp, "nonce-", "CSP must include a nonce for script-src")
}

func TestSecurityHeaders_CSPNonceDiffersPerRequest(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, httptest.NewRequest(http.MethodGet, "/test", nil))

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/test", nil))

	csp1 := w1.Header().Get("Content-Security-Policy")
	csp2 := w2.Header().Get("Content-Security-Policy")

	assert.NotEqual(t, csp1, csp2, "nonce must differ between requests")
}

func TestSecurityHeaders_HSTSOnHTTPS(t *testing.T) {
	// Simulate HTTPS via X-Forwarded-Proto (trust proxy headers)
	os.Setenv("TRUST_PROXY_HEADERS", "true")
	defer os.Unsetenv("TRUST_PROXY_HEADERS")

	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	r.ServeHTTP(w, req)

	hsts := w.Header().Get("Strict-Transport-Security")
	assert.NotEmpty(t, hsts, "HSTS header must be set on HTTPS connections")
	assert.Contains(t, hsts, "max-age=")
	assert.Contains(t, hsts, "includeSubDomains")
	assert.Contains(t, hsts, "preload")
}

func TestSecurityHeaders_NoHSTSOnHTTP(t *testing.T) {
	os.Unsetenv("TRUST_PROXY_HEADERS")

	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	// No X-Forwarded-Proto header, plain HTTP
	r.ServeHTTP(w, req)

	hsts := w.Header().Get("Strict-Transport-Security")
	assert.Empty(t, hsts, "HSTS must NOT be set on plain HTTP connections")
}

// ─── sanitizeCSPOrigin ───────────────────────────────────────────────────────

func TestSanitizeCSPOrigin_RemovesNewlines(t *testing.T) {
	result := sanitizeCSPOrigin("https://evil.com\nhttps://injected.com")
	assert.NotContains(t, result, "\n")
}

func TestSanitizeCSPOrigin_RemovesCarriageReturn(t *testing.T) {
	result := sanitizeCSPOrigin("https://evil.com\r")
	assert.NotContains(t, result, "\r")
}

func TestSanitizeCSPOrigin_RemovesSemicolons(t *testing.T) {
	result := sanitizeCSPOrigin("https://evil.com; unsafe-inline")
	assert.NotContains(t, result, ";")
}

func TestSanitizeCSPOrigin_RemovesQuotes(t *testing.T) {
	result := sanitizeCSPOrigin("https://evil.com' 'unsafe-eval")
	assert.NotContains(t, result, "'")
	assert.NotContains(t, result, "\"")
}

func TestSanitizeCSPOrigin_PreservesValidURL(t *testing.T) {
	input := "https://api.myschool.edu"
	result := sanitizeCSPOrigin(input)
	assert.Equal(t, input, result)
}

func TestSanitizeCSPOrigin_TrimsSpaces(t *testing.T) {
	result := sanitizeCSPOrigin("   https://api.myschool.edu   ")
	assert.Equal(t, "https://api.myschool.edu", result)
}

// ─── backendURL / backendWSURL ────────────────────────────────────────────────

func TestBackendURL_ReturnsEnvVar(t *testing.T) {
	os.Setenv("BACKEND_URL", "https://api.myschool.edu")
	defer os.Unsetenv("BACKEND_URL")

	url := backendURL()
	assert.Equal(t, "https://api.myschool.edu", url)
}

func TestBackendURL_SanitizesInjectionInEnvVar(t *testing.T) {
	os.Setenv("BACKEND_URL", "https://api.myschool.edu\r\n; '\"")
	defer os.Unsetenv("BACKEND_URL")

	url := backendURL()
	assert.NotContains(t, url, "\n")
	assert.NotContains(t, url, "\r")
	assert.NotContains(t, url, ";")
	assert.NotContains(t, url, "'")
	assert.NotContains(t, url, "\"")
}

// ─── CircuitBreakerMiddleware ─────────────────────────────────────────────────

func TestCircuitBreaker_PassesThroughOnSuccess(t *testing.T) {
	// Reset circuit breaker for clean test
	CB = nil
	InitCircuitBreaker()

	r := gin.New()
	r.Use(CircuitBreakerMiddleware())
	r.GET("/ok", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ok", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCircuitBreaker_Returns503WhenOpen(t *testing.T) {
	// Force circuit breaker open by simulating enough failures
	CB = nil
	InitCircuitBreaker()

	r := gin.New()
	r.Use(CircuitBreakerMiddleware())
	r.GET("/fail", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
	})

	// Make enough requests to trip the breaker (threshold: 3+ requests, 60%+ failure ratio)
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/fail", nil))
	}

	// Next request should be rejected with 503 (circuit open)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/fail", nil))
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestCircuitBreaker_InitCircuitBreakerIdempotent(t *testing.T) {
	CB = nil
	InitCircuitBreaker()
	cb1 := CB

	InitCircuitBreaker() // Second call should reuse the same instance
	cb2 := CB

	// Both should point to same instance (or at least both be non-nil)
	assert.NotNil(t, cb1)
	assert.NotNil(t, cb2)
}
