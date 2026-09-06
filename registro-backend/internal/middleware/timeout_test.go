package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimeoutMiddleware_FastRequestPasses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(TimeoutMiddleware(500 * time.Millisecond))
	r.GET("/fast", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/fast", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"ok":true`)
}

func TestTimeoutMiddleware_WebSocketExcluded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(TimeoutMiddleware(50 * time.Millisecond))
	r.GET("/ws-mock", func(c *gin.Context) {
		deadline, ok := c.Request.Context().Deadline()
		assert.False(t, ok, "websocket should not have a timeout deadline")
		_ = deadline
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ws-mock", nil)
	req.Header.Set("Upgrade", "websocket")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTimeoutMiddleware_ContextCancelledOnSlowQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(TimeoutMiddleware(50 * time.Millisecond))
	r.GET("/slow", func(c *gin.Context) {
		select {
		case <-time.After(150 * time.Millisecond):
			c.JSON(http.StatusOK, gin.H{"ok": true})
		case <-c.Request.Context().Done():
			// Handler detects context cancellation (as database queries do)
			return
		}
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/slow", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
	assert.Contains(t, w.Body.String(), "REQUEST_TIMEOUT")
}
