package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"registro-backend/internal/middleware"
)

func TestSecurity_SecurityHeadersMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sets mandatory security headers on HTTP responses", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.SecurityHeadersMiddleware())
		r.GET("/api/v1/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
		assert.Contains(t, w.Header().Get("Permissions-Policy"), "camera=()")
		assert.Contains(t, w.Header().Get("Content-Security-Policy"), "script-src 'self' 'nonce-")
	})

	t.Run("Sets HSTS header on HTTPS requests", func(t *testing.T) {
		t.Setenv("TRUST_PROXY_HEADERS", "true")
		r := gin.New()
		r.Use(middleware.SecurityHeadersMiddleware())
		r.GET("/api/v1/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
		req.Header.Set("X-Forwarded-Proto", "https")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Strict-Transport-Security"), "max-age=63072000")
	})
}
