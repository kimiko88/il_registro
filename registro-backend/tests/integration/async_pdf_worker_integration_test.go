package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/grades"
	"registro-backend/internal/pdfworker"
	"registro-backend/internal/scrutiny"
	"registro-backend/internal/verbali"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAsyncPdfWorkerIntegration_HandlerRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := scrutiny.NewService(nil, nil, nil, nil, nil)
	hScrutiny := scrutiny.NewHandler(svc, nil)

	gradesH := grades.NewHandler(nil, nil)
	verbaliH := verbali.NewHandler(nil)

	v1 := r.Group("/api/v1")
	v1.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-user-1")
		c.Set("role", "teacher")
		c.Next()
	})
	hScrutiny.RegisterRoutes(v1)
	gradesH.RegisterRoutes(v1)
	verbaliH.RegisterRoutes(v1)

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

	t.Run("POST /api/v1/scrutiny/export/:studentId/async-pdf safely handles unconfigured queue", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/scrutiny/export/student-123/async-pdf?semester=2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("POST /api/v1/grades/export/async-pdf safely handles unconfigured queue", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/grades/export/async-pdf?class_id=class-10A", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("POST /api/v1/verbali/:id/async-pdf safely handles unconfigured queue", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/verbali/verb-123/async-pdf", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("GET /api/v1/scrutiny/pdf-jobs/:job_id safely handles unconfigured queue", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/scrutiny/pdf-jobs/non-existent-job-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("Verify Task Payload Serialization", func(t *testing.T) {
		scPayload := pdfworker.ScrutinyPdfPayload{
			JobID:       "test-uuid-9999",
			ClassID:     "class-10A",
			SchoolYear:  "2025-2026",
			Period:      "2",
			RequestedBy: "teacher-user-1",
		}
		scTask, err := pdfworker.NewScrutinyPdfTask(scPayload)
		assert.NoError(t, err)
		assert.Equal(t, "pdf:generate_scrutiny", scTask.Type())

		rcPayload := pdfworker.ReportCardPdfPayload{
			JobID:       "test-uuid-8888",
			StudentID:   "student-1",
			ClassID:     "class-10A",
			SchoolYear:  "2025-2026",
			Period:      "2",
			RequestedBy: "teacher-user-1",
		}
		rcTask, err := pdfworker.NewReportCardPdfTask(rcPayload)
		assert.NoError(t, err)
		assert.Equal(t, "pdf:generate_report_card", rcTask.Type())

		regPayload := pdfworker.RegisterPdfPayload{
			JobID:       "test-uuid-7777",
			ClassID:     "class-10A",
			SubjectID:   "subj-math",
			SchoolYear:  "2025-2026",
			Period:      "2",
			RequestedBy: "teacher-user-1",
		}
		regTask, err := pdfworker.NewRegisterPdfTask(regPayload)
		assert.NoError(t, err)
		assert.Equal(t, "pdf:generate_register", regTask.Type())

		verbPayload := pdfworker.VerbalePdfPayload{
			JobID:       "test-uuid-6666",
			VerbaleID:   "verb-123",
			RequestedBy: "teacher-user-1",
		}
		verbTask, err := pdfworker.NewVerbalePdfTask(verbPayload)
		assert.NoError(t, err)
		assert.Equal(t, "pdf:generate_verbale", verbTask.Type())
	})
}
