package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"registro-backend/internal/middleware"
)

func TestSecurity_AuthRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Blocks excessive login attempts with 429 Too Many Requests", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.AuthRateLimitMiddleware())
		r.POST("/api/v1/auth/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// 1st request allowed
		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		req1.RemoteAddr = "192.168.1.50:12345"
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Immediate 2nd request from same IP should be blocked with 429
		reqBlocked := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		reqBlocked.RemoteAddr = "192.168.1.50:12345"
		wBlocked := httptest.NewRecorder()
		r.ServeHTTP(wBlocked, reqBlocked)

		assert.Equal(t, http.StatusTooManyRequests, wBlocked.Code)
		assert.Contains(t, wBlocked.Body.String(), "AUTH_RATE_LIMIT_EXCEEDED")
	})

	t.Run("Does not block traffic from a different IP address", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.AuthRateLimitMiddleware())
		r.POST("/api/v1/auth/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// 1st IP exhausts burst
		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		req1.RemoteAddr = "10.0.0.1:12345"
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// 2nd IP should still be allowed
		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		req2.RemoteAddr = "10.0.0.2:12345"
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)
	})
}
