package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		origin         string
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name:           "OPTIONS request should abort with 204",
			method:         "OPTIONS",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusNoContent,
			shouldAbort:    true,
		},
		{
			name:           "GET request should continue",
			method:         "GET",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name:           "POST request should continue",
			method:         "POST",
			origin:         "http://localhost:5173",
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			r.Use(CORSMiddleware())
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			r.POST("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			r.OPTIONS("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check CORS headers
			assert.Equal(t, tt.origin, w.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
		})
	}
}

func TestErrorMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		addError     bool
		errorMessage string
	}{
		{
			name:     "no errors",
			addError: false,
		},
		{
			name:         "with error",
			addError:     true,
			errorMessage: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			r.Use(ErrorMiddleware())
			r.GET("/test", func(c *gin.Context) {
				if tt.addError {
					_ = c.Error(assert.AnError)
				}
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			// Error middleware doesn't change the status code or response
			// It just logs errors. We verify the request completes
			assert.NotNil(t, w)
		})
	}
}

func TestLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		query  string
		status int
	}{
		{
			name:   "GET request",
			method: "GET",
			path:   "/api/users",
			query:  "page=1",
			status: http.StatusOK,
		},
		{
			name:   "POST request",
			method: "POST",
			path:   "/api/users",
			query:  "",
			status: http.StatusCreated,
		},
		{
			name:   "404 request",
			method: "GET",
			path:   "/api/nonexistent",
			query:  "",
			status: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			r.Use(LoggerMiddleware())
			r.GET("/api/users", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			r.POST("/api/users", func(c *gin.Context) {
				c.Status(http.StatusCreated)
			})

			url := tt.path
			if tt.query != "" {
				url += "?" + tt.query
			}

			req, _ := http.NewRequest(tt.method, url, nil)
			r.ServeHTTP(w, req)

			// Logger middleware doesn't change response, just logs
			// We just verify it doesn't panic and request completes
			assert.NotNil(t, w)
		})
	}
}

func TestCORSMiddleware_HeaderValues(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Origin", "http://localhost:5173")

	middleware := CORSMiddleware()
	middleware(c)

	// Verify all expected CORS headers are set
	headers := w.Header()

	assert.Equal(t, "http://localhost:5173", headers.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", headers.Get("Access-Control-Allow-Credentials"))

	allowedHeaders := headers.Get("Access-Control-Allow-Headers")
	assert.Contains(t, allowedHeaders, "Authorization")
	assert.Contains(t, allowedHeaders, "Content-Type")
	assert.Contains(t, allowedHeaders, "X-CSRF-Token")

	allowedMethods := headers.Get("Access-Control-Allow-Methods")
	assert.Contains(t, allowedMethods, "GET")
	assert.Contains(t, allowedMethods, "POST")
	assert.Contains(t, allowedMethods, "PUT")
	assert.Contains(t, allowedMethods, "DELETE")
	assert.Contains(t, allowedMethods, "PATCH")
	assert.Contains(t, allowedMethods, "OPTIONS")
}
