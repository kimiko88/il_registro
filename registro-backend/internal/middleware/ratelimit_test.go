package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimit_TieredLimitsAndIETFHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	mb := newMemoryBackend()

	// 1. Test Export Limit (burst = 2)
	ipExport := "10.10.10.1"
	assert.True(t, mb.allowExport(ctx, ipExport), "1st export allowed")
	assert.True(t, mb.allowExport(ctx, ipExport), "2nd export allowed (burst=2)")
	assert.False(t, mb.allowExport(ctx, ipExport), "3rd export rate limited")

	// 2. Test Upload Limit (burst = 5)
	ipUpload := "10.10.10.2"
	for i := 0; i < 5; i++ {
		assert.True(t, mb.allowUpload(ctx, ipUpload), "upload burst request %d allowed", i+1)
	}
	assert.False(t, mb.allowUpload(ctx, ipUpload), "6th upload rate limited")
}

func TestExportRateLimitMiddleware_Headers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ExportRateLimitMiddleware())
	r.GET("/test-export", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 1st request should succeed and have RateLimit headers
	req1, _ := http.NewRequest(http.MethodGet, "/test-export", nil)
	req1.RemoteAddr = "172.16.0.1:1234"
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "5", w1.Header().Get("RateLimit-Limit"))
	assert.NotEmpty(t, w1.Header().Get("RateLimit-Remaining"))
}

func TestUploadRateLimitMiddleware_Headers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(UploadRateLimitMiddleware())
	r.POST("/test-upload", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req1, _ := http.NewRequest(http.MethodPost, "/test-upload", nil)
	req1.RemoteAddr = "172.16.0.2:1234"
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "15", w1.Header().Get("RateLimit-Limit"))
	assert.NotEmpty(t, w1.Header().Get("RateLimit-Remaining"))
}
