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
	r.POST("/catalog", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"created": true})
	})
	r.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db failure"})
	})

	// 1. Initial request -> 200 OK + ETag header
	req1, _ := http.NewRequest("GET", "/catalog", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	etag := w1.Header().Get("ETag")
	assert.NotEmpty(t, etag)
	assert.Contains(t, w1.Body.String(), "1A")

	// 2. Request with matching exact If-None-Match -> 304 Not Modified and empty body
	req2, _ := http.NewRequest("GET", "/catalog", nil)
	req2.Header.Set("If-None-Match", etag)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotModified, w2.Code)
	assert.Empty(t, w2.Body.String())

	// 3. Request with weak validator prefix (W/...) -> 304 Not Modified
	reqWeak, _ := http.NewRequest("GET", "/catalog", nil)
	reqWeak.Header.Set("If-None-Match", "W/"+etag)
	wWeak := httptest.NewRecorder()
	r.ServeHTTP(wWeak, reqWeak)
	assert.Equal(t, http.StatusNotModified, wWeak.Code)

	// 4. Request with wildcard (*) validator -> 304 Not Modified
	reqWildcard, _ := http.NewRequest("GET", "/catalog", nil)
	reqWildcard.Header.Set("If-None-Match", "*")
	wWildcard := httptest.NewRecorder()
	r.ServeHTTP(wWildcard, reqWildcard)
	assert.Equal(t, http.StatusNotModified, wWildcard.Code)

	// 5. Request with mismatched If-None-Match -> 200 OK with body
	req3, _ := http.NewRequest("GET", "/catalog", nil)
	req3.Header.Set("If-None-Match", `"different-etag"`)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Contains(t, w3.Body.String(), "1A")

	// 6. Non-GET request (POST) -> bypasses ETag computation
	reqPost, _ := http.NewRequest("POST", "/catalog", nil)
	wPost := httptest.NewRecorder()
	r.ServeHTTP(wPost, reqPost)
	assert.Equal(t, http.StatusCreated, wPost.Code)
	assert.Empty(t, wPost.Header().Get("ETag"))

	// 7. Non-200 GET request (e.g. 500 Internal Server Error) -> no ETag, returns error body
	reqErr, _ := http.NewRequest("GET", "/error", nil)
	wErr := httptest.NewRecorder()
	r.ServeHTTP(wErr, reqErr)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)
	assert.Empty(t, wErr.Header().Get("ETag"))
	assert.Contains(t, wErr.Body.String(), "db failure")
}
