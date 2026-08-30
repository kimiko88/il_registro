package pdfworker

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// Task types
const (
	TypeScrutinyPdf   = "pdf:generate_scrutiny"
	TypeReportCardPdf = "pdf:generate_report_card"
)

// ScrutinyPdfPayload contains parameters for generating a scrutiny PDF
type ScrutinyPdfPayload struct {
	JobID       string `json:"job_id"`
	ClassID     string `json:"class_id"`
	SchoolYear  string `json:"school_year"`
	Period      string `json:"period"`
	RequestedBy string `json:"requested_by"`
}

// ReportCardPdfPayload contains parameters for generating an individual student report card PDF
type ReportCardPdfPayload struct {
	JobID       string `json:"job_id"`
	StudentID   string `json:"student_id"`
	SchoolYear  string `json:"school_year"`
	Period      string `json:"period"`
	RequestedBy string `json:"requested_by"`
}

// NewScrutinyPdfTask creates an asynq.Task for scrutiny PDF generation
func NewScrutinyPdfTask(payload ScrutinyPdfPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal scrutiny payload: %w", err)
	}
	return asynq.NewTask(TypeScrutinyPdf, data, asynq.MaxRetry(3)), nil
}

// NewReportCardPdfTask creates an asynq.Task for report card PDF generation
func NewReportCardPdfTask(payload ReportCardPdfPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report card payload: %w", err)
	}
	return asynq.NewTask(TypeReportCardPdf, data, asynq.MaxRetry(3)), nil
}
