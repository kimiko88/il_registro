package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("postgres", "postgres://localhost:5432/test?sslmode=disable")
	assert.NoError(t, err)

	h := NewHealthHandler(db)

	// 1. Health endpoint
	{
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		h.Health(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"UP"`)
	}

	// 2. Metrics endpoint
	{
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		h.Metrics(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "db_open_connections")
	}

	// 3. Ready endpoint - DB closed (returns 503)
	{
		_ = db.Close()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		h.Ready(c)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"DOWN"`)
	}

	// 4. Ping endpoint - GET and HEAD
	{
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
		h.Ping(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"ok":true`)
	}
	{
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodHead, "/api/v1/ping", nil)
		h.Ping(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Body.String())
	}

	// 5. Ping route matching via Gin Engine
	{
		r := gin.New()
		r.GET("/api/v1/ping", h.Ping)
		r.HEAD("/api/v1/ping", h.Ping)

		reqHead := httptest.NewRequest(http.MethodHead, "/api/v1/ping", nil)
		wHead := httptest.NewRecorder()
		r.ServeHTTP(wHead, reqHead)
		assert.Equal(t, http.StatusOK, wHead.Code)

		reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
		wGet := httptest.NewRecorder()
		r.ServeHTTP(wGet, reqGet)
		assert.Equal(t, http.StatusOK, wGet.Code)
		assert.Contains(t, wGet.Body.String(), `"ok":true`)
	}

	// 6. Live probe (Kubernetes Liveness)
	{
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		h.Live(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"alive"`)
		assert.Contains(t, w.Body.String(), `"timestamp"`)
	}
}
