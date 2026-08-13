package unit

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"registro-backend/internal/middleware"
)

func TestSecurity_CORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Allows authorized origin from whitelist and sets credentials header", func(t *testing.T) {
		os.Setenv("ALLOWED_ORIGINS", "http://localhost:5173,https://app.scuola.it")
		defer os.Unsetenv("ALLOWED_ORIGINS")

		r := gin.New()
		r.Use(middleware.CORSMiddleware())
		r.GET("/api/v1/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Contains(t, w.Header().Get("Vary"), "Origin")
	})

	t.Run("Does not set Access-Control-Allow-Origin header for unauthorized origin", func(t *testing.T) {
		os.Setenv("ALLOWED_ORIGINS", "http://localhost:5173")
		defer os.Unsetenv("ALLOWED_ORIGINS")

		r := gin.New()
		r.Use(middleware.CORSMiddleware())
		r.GET("/api/v1/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
		req.Header.Set("Origin", "http://malicious-site.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Responds with 204 No Content for OPTIONS preflight request", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.CORSMiddleware())
		r.POST("/api/v1/data", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodOptions, "/api/v1/data", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
