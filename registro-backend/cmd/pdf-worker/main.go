package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"registro-backend/internal/config"
	"registro-backend/internal/pdfworker"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: error reading config file: %v", err)
	}


	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	if cfg.Redis.Host == "" {
		redisAddr = "localhost:6379"
	}

	client := pdfworker.NewClient(redisAddr)
	defer client.Close()

	// Handler function for PDF generation
	scrutinyHandlerFunc := func(ctx context.Context, payload pdfworker.ScrutinyPdfPayload) ([]byte, error) {
		log.Printf("[PdfWorkerStandalone] Generating Scrutiny PDF for class: %s, job: %s", payload.ClassID, payload.JobID)
		// Dummy PDF bytes generated in background queue runner
		return []byte("%PDF-1.4 Async Scrutiny PDF Generated Successfully"), nil
	}

	server := pdfworker.NewServer(redisAddr, 10, client, scrutinyHandlerFunc)

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
