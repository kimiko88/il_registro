package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/pdfworker"
	"registro-backend/internal/scrutiny"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAsyncPdfWorkerIntegration_HandlerRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := scrutiny.NewService(nil, nil, nil, nil, nil)
	// Without live Redis connection, handler safely returns 503 ServiceUnavailable
	h := scrutiny.NewHandler(svc, nil)

	v1 := r.Group("/api/v1")
	v1.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-user-1")
		c.Set("role", "teacher")
		c.Next()
	})
	h.RegisterRoutes(v1)


	t.Run("POST /api/v1/scrutiny/class/:classId/async-pdf safely handles unconfigured queue", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/scrutiny/class/class-10A/async-pdf?semester=2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp["error"], "not configured")
	})

	t.Run("GET /api/v1/scrutiny/pdf-jobs/:job_id safely handles unconfigured queue", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/scrutiny/pdf-jobs/non-existent-job-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("Verify Task Payload Serialization", func(t *testing.T) {
		payload := pdfworker.ScrutinyPdfPayload{
			JobID:       "test-uuid-9999",
			ClassID:     "class-10A",
			SchoolYear:  "2025-2026",
			Period:      "2",
			RequestedBy: "teacher-user-1",
		}

		task, err := pdfworker.NewScrutinyPdfTask(payload)
		assert.NoError(t, err)
		assert.Equal(t, "pdf:generate_scrutiny", task.Type())
		assert.NotEmpty(t, task.Payload())
	})
}
