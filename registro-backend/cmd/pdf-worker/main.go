package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"registro-backend/internal/pdfworker"
)

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if (os.Getenv("SERVER_MODE") == "release" || os.Getenv("APP_ENV") == "production") && redisPassword == "" {
		log.Fatal("REDIS_PASSWORD is required in production")
	}
	redisAddr := net.JoinHostPort(redisHost, redisPort)

	client := pdfworker.NewClient(redisAddr, redisPassword)
	defer func() { _ = client.Close() }()

	startupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := client.Ping(startupCtx); err != nil {
		cancel()
		log.Fatalf("PDF worker cannot connect to Redis: %v", err)
	}
	cancel()

	// Handler functions for PDF generation
	scrutinyHandlerFunc := func(ctx context.Context, payload pdfworker.ScrutinyPdfPayload) ([]byte, error) {
		log.Printf("[PdfWorkerStandalone] Generating Scrutiny PDF for class: %s, job: %s", payload.ClassID, payload.JobID)
		return []byte("%PDF-1.4 Async Scrutiny PDF Generated Successfully"), nil
	}

	reportCardHandlerFunc := func(ctx context.Context, payload pdfworker.ReportCardPdfPayload) ([]byte, error) {
		log.Printf("[PdfWorkerStandalone] Generating ReportCard PDF for student: %s, job: %s", payload.StudentID, payload.JobID)
		return []byte("%PDF-1.4 Async ReportCard PDF Generated Successfully"), nil
	}

	registerHandlerFunc := func(ctx context.Context, payload pdfworker.RegisterPdfPayload) ([]byte, error) {
		log.Printf("[PdfWorkerStandalone] Generating Register PDF for class: %s, job: %s", payload.ClassID, payload.JobID)
		return []byte("%PDF-1.4 Async Register PDF Generated Successfully"), nil
	}

	verbaleHandlerFunc := func(ctx context.Context, payload pdfworker.VerbalePdfPayload) ([]byte, error) {
		log.Printf("[PdfWorkerStandalone] Generating Verbale PDF for verbale: %s, job: %s", payload.VerbaleID, payload.JobID)
		return []byte("%PDF-1.4 Async Verbale PDF Generated Successfully"), nil
	}

	server := pdfworker.NewServer(redisAddr, 10, client, scrutinyHandlerFunc, redisPassword)
	server.SetReportCardHandler(reportCardHandlerFunc)
	server.SetRegisterHandler(registerHandlerFunc)
	server.SetVerbaleHandler(verbaleHandlerFunc)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start pdf-worker server: %v", err)
		}
	}()

	log.Println("[PdfWorkerStandalone] Started background PDF worker node. Press Ctrl+C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[PdfWorkerStandalone] Shutting down pdf-worker cleanly...")
	server.Shutdown()
}
