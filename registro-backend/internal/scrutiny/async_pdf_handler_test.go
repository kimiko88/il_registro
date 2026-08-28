package scrutiny_test

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

func TestEnqueueAsyncScrutinyPdf_NoWorkerClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := scrutiny.NewService(nil, nil, nil, nil, nil)
	h := scrutiny.NewHandler(svc, nil)

	r.POST("/api/v1/scrutiny/class/:classId/async-pdf", h.EnqueueAsyncScrutinyPdf)

	req := httptest.NewRequest("POST", "/api/v1/scrutiny/class/class-123/async-pdf", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestGetPdfJobStatus_NoWorkerClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := scrutiny.NewService(nil, nil, nil, nil, nil)
	h := scrutiny.NewHandler(svc, nil)

	r.GET("/api/v1/scrutiny/pdf-jobs/:job_id", h.GetPdfJobStatus)

	req := httptest.NewRequest("GET", "/api/v1/scrutiny/pdf-jobs/test-job-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestTaskPayloadStructure(t *testing.T) {
	payload := pdfworker.ScrutinyPdfPayload{
		JobID:       "job-1",
		ClassID:     "class-1",
		SchoolYear:  "2025-2026",
		Period:      "1",
		RequestedBy: "teacher-1",
	}

	data, err := json.Marshal(payload)
	assert.NoError(t, err)

	var unmarshaled pdfworker.ScrutinyPdfPayload
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, "job-1", unmarshaled.JobID)
	assert.Equal(t, "class-1", unmarshaled.ClassID)
	assert.Equal(t, "2025-2026", unmarshaled.SchoolYear)
}
