package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestBodyLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestBodyLimitMiddleware(1024)) // 1 KB limit for testing
	r.POST("/test", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "payload too large"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"bytes": len(body)})
	})

	// 1. Small payload under limit -> 200 OK
	smallPayload := bytes.Repeat([]byte("a"), 500)
	req1, _ := http.NewRequest("POST", "/test", bytes.NewReader(smallPayload))
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 2. Large payload over limit -> 413 Request Entity Too Large
	largePayload := bytes.Repeat([]byte("a"), 2048)
	req2, _ := http.NewRequest("POST", "/test", bytes.NewReader(largePayload))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusRequestEntityTooLarge, w2.Code)

	// 3. Multipart request over limit -> skipped by middleware, allowed to pass through
	req3, _ := http.NewRequest("POST", "/test", bytes.NewReader(largePayload))
	req3.Header.Set("Content-Type", "multipart/form-data; boundary=something")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
