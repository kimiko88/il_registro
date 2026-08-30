package pdfworker_test

import (
	"testing"
	"time"

	"registro-backend/internal/pdfworker"

	"github.com/stretchr/testify/assert"
)

func TestNewScrutinyPdfTask(t *testing.T) {
	payload := pdfworker.ScrutinyPdfPayload{
		JobID:       "test-job-123",
		ClassID:     "class-456",
		SchoolYear:  "2025-2026",
		Period:      "1",
		RequestedBy: "teacher-789",
	}

	task, err := pdfworker.NewScrutinyPdfTask(payload)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, pdfworker.TypeScrutinyPdf, task.Type())
}

func TestNewReportCardPdfTask(t *testing.T) {
	payload := pdfworker.ReportCardPdfPayload{
		JobID:       "test-job-rc-1",
		StudentID:   "student-uuid-42",
		SchoolYear:  "2025-2026",
		Period:      "2",
		RequestedBy: "coordinator-1",
	}

	task, err := pdfworker.NewReportCardPdfTask(payload)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, pdfworker.TypeReportCardPdf, task.Type())
}

func TestJobStatusStruct(t *testing.T) {
	now := time.Now()
	status := pdfworker.JobStatus{
		JobID:     "job-999",
		Status:    "pending",
		TaskType:  pdfworker.TypeScrutinyPdf,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "job-999", status.JobID)
	assert.Equal(t, "pending", status.Status)
	assert.Equal(t, pdfworker.TypeScrutinyPdf, status.TaskType)
	assert.Equal(t, now, status.CreatedAt)
	assert.Equal(t, now, status.UpdatedAt)
}

