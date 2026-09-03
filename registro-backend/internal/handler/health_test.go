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
}
