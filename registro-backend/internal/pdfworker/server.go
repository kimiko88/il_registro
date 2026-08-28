package pdfworker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

type ScrutinyPdfHandlerFunc func(ctx context.Context, payload ScrutinyPdfPayload) ([]byte, error)

type Server struct {
	asynqServer         *asynq.Server
	client              *Client
	scrutinyHandlerFunc ScrutinyPdfHandlerFunc
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

func (s *Server) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeScrutinyPdf, s.handleScrutinyPdfTask)

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

	// Update status to processing
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

	// Mark as completed
	return s.client.UpdateJobStatus(ctx, &JobStatus{
		JobID:       p.JobID,
		Status:      "completed",
		TaskType:    TypeScrutinyPdf,
		DownloadURL: fmt.Sprintf("/api/v1/scrutiny/classes/%s/pdf", p.ClassID),
	})
}
