package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestETagMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(ETagMiddleware())
	r.GET("/catalog", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"classes": []string{"1A", "2B", "3C"}})
	})

	// 1. Initial request -> 200 OK + ETag header
	req1, _ := http.NewRequest("GET", "/catalog", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	etag := w1.Header().Get("ETag")
	assert.NotEmpty(t, etag)
	assert.Contains(t, w1.Body.String(), "1A")

	// 2. Request with matching If-None-Match -> 304 Not Modified and empty body
	req2, _ := http.NewRequest("GET", "/catalog", nil)
	req2.Header.Set("If-None-Match", etag)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotModified, w2.Code)
	assert.Empty(t, w2.Body.String())

	// 3. Request with mismatched If-None-Match -> 200 OK with body
	req3, _ := http.NewRequest("GET", "/catalog", nil)
	req3.Header.Set("If-None-Match", `"different-etag"`)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Contains(t, w3.Body.String(), "1A")
}
