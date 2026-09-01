package pdfworker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

type ScrutinyPdfHandlerFunc func(ctx context.Context, payload ScrutinyPdfPayload) ([]byte, error)
type ReportCardPdfHandlerFunc func(ctx context.Context, payload ReportCardPdfPayload) ([]byte, error)
type RegisterPdfHandlerFunc func(ctx context.Context, payload RegisterPdfPayload) ([]byte, error)
type VerbalePdfHandlerFunc func(ctx context.Context, payload VerbalePdfPayload) ([]byte, error)

type Server struct {
	asynqServer           *asynq.Server
	client                *Client
	scrutinyHandlerFunc   ScrutinyPdfHandlerFunc
	reportCardHandlerFunc ReportCardPdfHandlerFunc
	registerHandlerFunc   RegisterPdfHandlerFunc
	verbaleHandlerFunc    VerbalePdfHandlerFunc
}

func NewServer(redisAddr string, concurrency int, client *Client, scrutinyHandler ScrutinyPdfHandlerFunc) *Server {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &Server{
		asynqServer:         srv,
		client:              client,
		scrutinyHandlerFunc: scrutinyHandler,
	}
}

func (s *Server) SetReportCardHandler(h ReportCardPdfHandlerFunc) {
	s.reportCardHandlerFunc = h
}

func (s *Server) SetRegisterHandler(h RegisterPdfHandlerFunc) {
	s.registerHandlerFunc = h
}

func (s *Server) SetVerbaleHandler(h VerbalePdfHandlerFunc) {
	s.verbaleHandlerFunc = h
}

func (s *Server) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeScrutinyPdf, s.handleScrutinyPdfTask)
	mux.HandleFunc(TypeReportCardPdf, s.handleReportCardPdfTask)
	mux.HandleFunc(TypeRegisterPdf, s.handleRegisterPdfTask)
	mux.HandleFunc(TypeVerbalePdf, s.handleVerbalePdfTask)

	log.Println("[PdfWorker] Server listening for background PDF jobs...")
	return s.asynqServer.Start(mux)
}

func (s *Server) Shutdown() {
	if s.asynqServer != nil {
		s.asynqServer.Shutdown()
	}
}

func (s *Server) handleScrutinyPdfTask(ctx context.Context, t *asynq.Task) error {
	var p ScrutinyPdfPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", err)
	}

	log.Printf("[PdfWorker] Processing Scrutiny PDF Job ID: %s for Class: %s", p.JobID, p.ClassID)

	_ = s.client.UpdateJobStatus(ctx, &JobStatus{
		JobID:    p.JobID,
		Status:   "processing",
		TaskType: TypeScrutinyPdf,
	})

	if s.scrutinyHandlerFunc == nil {
		err := fmt.Errorf("no handler function registered for scrutiny pdf")
		_ = s.client.UpdateJobStatus(ctx, &JobStatus{
			JobID:    p.JobID,
			Status:   "failed",
			TaskType: TypeScrutinyPdf,
			Error:    err.Error(),
		})
		return err
	}

	pdfBytes, err := s.scrutinyHandlerFunc(ctx, p)
	if err != nil {
		log.Printf("[PdfWorker] Failed PDF generation for Job ID %s: %v", p.JobID, err)
		_ = s.client.UpdateJobStatus(ctx, &JobStatus{
			JobID:    p.JobID,
			Status:   "failed",
			TaskType: TypeScrutinyPdf,
			Error:    err.Error(),
		})
		return err
	}

	_ = pdfBytes
	log.Printf("[PdfWorker] Successfully completed Scrutiny PDF Job ID: %s", p.JobID)

	return s.client.UpdateJobStatus(ctx, &JobStatus{
		JobID:       p.JobID,
		Status:      "completed",
		TaskType:    TypeScrutinyPdf,
		DownloadURL: fmt.Sprintf("/api/v1/scrutiny/classes/%s/pdf", p.ClassID),
	})
}

func (s *Server) handleReportCardPdfTask(ctx context.Context, t *asynq.Task) error {
	var p ReportCardPdfPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal report card task payload: %w", err)
	}

	log.Printf("[PdfWorker] Processing ReportCard PDF Job ID: %s for Student: %s", p.JobID, p.StudentID)

	_ = s.client.UpdateJobStatus(ctx, &JobStatus{
		JobID:    p.JobID,
		Status:   "processing",
		TaskType: TypeReportCardPdf,
	})

	if s.reportCardHandlerFunc == nil {
		err := fmt.Errorf("no handler function registered for report card pdf")
		_ = s.client.UpdateJobStatus(ctx, &JobStatus{
			JobID:    p.JobID,
			Status:   "failed",
			TaskType: TypeReportCardPdf,
			Error:    err.Error(),
		})
		return err
	}

	pdfBytes, err := s.reportCardHandlerFunc(ctx, p)
	if err != nil {
		log.Printf("[PdfWorker] Failed ReportCard PDF generation for Job ID %s: %v", p.JobID, err)
		_ = s.client.UpdateJobStatus(ctx, &JobStatus{
			JobID:    p.JobID,
			Status:   "failed",
			TaskType: TypeReportCardPdf,
			Error:    err.Error(),
		})
		return err
	}

	_ = pdfBytes
	log.Printf("[PdfWorker] Successfully completed ReportCard PDF Job ID: %s", p.JobID)

	return s.client.UpdateJobStatus(ctx, &JobStatus{
		JobID:       p.JobID,
		Status:      "completed",
		TaskType:    TypeReportCardPdf,
		DownloadURL: fmt.Sprintf("/api/v1/scrutiny/export/%s/pdf", p.StudentID),
	})
}

func (s *Server) handleRegisterPdfTask(ctx context.Context, t *asynq.Task) error {
	var p RegisterPdfPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal register task payload: %w", err)
	}

	log.Printf("[PdfWorker] Processing Register PDF Job ID: %s for Class: %s", p.JobID, p.ClassID)

	_ = s.client.UpdateJobStatus(ctx, &JobStatus{
		JobID:    p.JobID,
		Status:   "processing",
		TaskType: TypeRegisterPdf,
	})

	if s.registerHandlerFunc == nil {
		err := fmt.Errorf("no handler function registered for register pdf")
		_ = s.client.UpdateJobStatus(ctx, &JobStatus{
			JobID:    p.JobID,
			Status:   "failed",
			TaskType: TypeRegisterPdf,
			Error:    err.Error(),
		})
		return err
	}

	pdfBytes, err := s.registerHandlerFunc(ctx, p)
	if err != nil {
		log.Printf("[PdfWorker] Failed Register PDF generation for Job ID %s: %v", p.JobID, err)
		_ = s.client.UpdateJobStatus(ctx, &JobStatus{
			JobID:    p.JobID,
			Status:   "failed",
			TaskType: TypeRegisterPdf,
			Error:    err.Error(),
		})
		return err
	}

	_ = pdfBytes
	log.Printf("[PdfWorker] Successfully completed Register PDF Job ID: %s", p.JobID)

	return s.client.UpdateJobStatus(ctx, &JobStatus{
		JobID:       p.JobID,
		Status:      "completed",
		TaskType:    TypeRegisterPdf,
		DownloadURL: fmt.Sprintf("/api/v1/grades/export?format=pdf&class_id=%s", p.ClassID),
	})
}

func (s *Server) handleVerbalePdfTask(ctx context.Context, t *asynq.Task) error {
	var p VerbalePdfPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal verbale task payload: %w", err)
	}

	log.Printf("[PdfWorker] Processing Verbale PDF Job ID: %s for Verbale ID: %s", p.JobID, p.VerbaleID)

	_ = s.client.UpdateJobStatus(ctx, &JobStatus{
		JobID:    p.JobID,
		Status:   "processing",
		TaskType: TypeVerbalePdf,
	})

	if s.verbaleHandlerFunc == nil {
		err := fmt.Errorf("no handler function registered for verbale pdf")
		_ = s.client.UpdateJobStatus(ctx, &JobStatus{
			JobID:    p.JobID,
			Status:   "failed",
			TaskType: TypeVerbalePdf,
			Error:    err.Error(),
		})
		return err
	}

	pdfBytes, err := s.verbaleHandlerFunc(ctx, p)
	if err != nil {
		log.Printf("[PdfWorker] Failed Verbale PDF generation for Job ID %s: %v", p.JobID, err)
		_ = s.client.UpdateJobStatus(ctx, &JobStatus{
			JobID:    p.JobID,
			Status:   "failed",
			TaskType: TypeVerbalePdf,
			Error:    err.Error(),
		})
		return err
	}

	_ = pdfBytes
	log.Printf("[PdfWorker] Successfully completed Verbale PDF Job ID: %s", p.JobID)

	return s.client.UpdateJobStatus(ctx, &JobStatus{
		JobID:       p.JobID,
		Status:      "completed",
		TaskType:    TypeVerbalePdf,
		DownloadURL: fmt.Sprintf("/api/v1/verbali/%s/pdf", p.VerbaleID),
	})
}
