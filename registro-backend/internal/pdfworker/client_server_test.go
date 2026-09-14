package pdfworker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func fastTestClient() *Client {
	return &Client{
		rdb: redis.NewClient(&redis.Options{
			Addr:        "127.0.0.1:63799",
			MaxRetries:  -1,
			DialTimeout: 10 * time.Millisecond,
		}),
	}
}

func TestTaskCreation_AllTypes(t *testing.T) {
	// Register PDF Task
	regPayload := RegisterPdfPayload{
		JobID:       "job-reg-1",
		ClassID:     "class-1",
		SubjectID:   "math",
		SchoolYear:  "2025/2026",
		Period:      "Q1",
		RequestedBy: "u1",
	}
	regTask, err := NewRegisterPdfTask(regPayload)
	assert.NoError(t, err)
	assert.NotNil(t, regTask)
	assert.Equal(t, TypeRegisterPdf, regTask.Type())

	// Verbale PDF Task
	verbPayload := VerbalePdfPayload{
		JobID:       "job-verb-1",
		VerbaleID:   "verb-1",
		RequestedBy: "u2",
	}
	verbTask, err := NewVerbalePdfTask(verbPayload)
	assert.NoError(t, err)
	assert.NotNil(t, verbTask)
	assert.Equal(t, TypeVerbalePdf, verbTask.Type())
}

func TestClient_LifecycleAndErrors(t *testing.T) {
	client := fastTestClient()
	defer client.Close()

	// Canceled context for instant returns without waiting for network timeouts
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// UpdateJobStatus error
	status := &JobStatus{
		JobID:    "job-test",
		Status:   "pending",
		TaskType: TypeScrutinyPdf,
	}
	err := client.UpdateJobStatus(ctx, status)
	assert.Error(t, err)

	// GetJobStatus error
	_, err = client.GetJobStatus(ctx, "job-test")
	assert.Error(t, err)

	// NewClient constructor
	liveClient := NewClient("127.0.0.1:63799")
	assert.NotNil(t, liveClient)
	_ = liveClient.Close()

	// Empty client close
	emptyClient := &Client{}
	assert.NoError(t, emptyClient.Close())
}

func TestServer_Handlers(t *testing.T) {
	client := fastTestClient()
	defer client.Close()

	srv := &Server{
		client: client,
	}

	srv.SetReportCardHandler(nil)
	srv.SetRegisterHandler(nil)
	srv.SetVerbaleHandler(nil)
	srv.Shutdown()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// 1. Scrutiny task handlers
	t.Run("ScrutinyPdfTask", func(t *testing.T) {
		// Invalid JSON
		badTask := asynq.NewTask(TypeScrutinyPdf, []byte("{invalid"))
		err := srv.handleScrutinyPdfTask(ctx, badTask)
		assert.Error(t, err)

		// Nil handler
		validPayload := ScrutinyPdfPayload{JobID: "job-s1", ClassID: "c1"}
		validTask, _ := NewScrutinyPdfTask(validPayload)
		err = srv.handleScrutinyPdfTask(ctx, validTask)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no handler function")

		// Handler returns error
		srv.scrutinyHandlerFunc = func(ctx context.Context, payload ScrutinyPdfPayload) ([]byte, error) {
			return nil, errors.New("pdf generation failed")
		}
		err = srv.handleScrutinyPdfTask(ctx, validTask)
		assert.Error(t, err)
		assert.Equal(t, "pdf generation failed", err.Error())

		// Handler succeeds
		srv.scrutinyHandlerFunc = func(ctx context.Context, payload ScrutinyPdfPayload) ([]byte, error) {
			return []byte("dummy-pdf"), nil
		}
		_ = srv.handleScrutinyPdfTask(ctx, validTask)
	})

	// 2. Report card task handlers
	t.Run("ReportCardPdfTask", func(t *testing.T) {
		// Invalid JSON
		badTask := asynq.NewTask(TypeReportCardPdf, []byte("{invalid"))
		err := srv.handleReportCardPdfTask(ctx, badTask)
		assert.Error(t, err)

		// Nil handler
		validPayload := ReportCardPdfPayload{JobID: "job-rc1", StudentID: "st1"}
		validTask, _ := NewReportCardPdfTask(validPayload)
		err = srv.handleReportCardPdfTask(ctx, validTask)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no handler function")

		// Handler returns error
		srv.SetReportCardHandler(func(ctx context.Context, payload ReportCardPdfPayload) ([]byte, error) {
			return nil, errors.New("report card failed")
		})
		err = srv.handleReportCardPdfTask(ctx, validTask)
		assert.Error(t, err)
		assert.Equal(t, "report card failed", err.Error())

		// Handler succeeds
		srv.SetReportCardHandler(func(ctx context.Context, payload ReportCardPdfPayload) ([]byte, error) {
			return []byte("rc-pdf"), nil
		})
		_ = srv.handleReportCardPdfTask(ctx, validTask)
	})

	// 3. Register task handlers
	t.Run("RegisterPdfTask", func(t *testing.T) {
		// Invalid JSON
		badTask := asynq.NewTask(TypeRegisterPdf, []byte("{invalid"))
		err := srv.handleRegisterPdfTask(ctx, badTask)
		assert.Error(t, err)

		// Nil handler
		validPayload := RegisterPdfPayload{JobID: "job-reg1", ClassID: "c1"}
		validTask, _ := NewRegisterPdfTask(validPayload)
		err = srv.handleRegisterPdfTask(ctx, validTask)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no handler function")

		// Handler returns error
		srv.SetRegisterHandler(func(ctx context.Context, payload RegisterPdfPayload) ([]byte, error) {
			return nil, errors.New("register pdf failed")
		})
		err = srv.handleRegisterPdfTask(ctx, validTask)
		assert.Error(t, err)
		assert.Equal(t, "register pdf failed", err.Error())

		// Handler succeeds
		srv.SetRegisterHandler(func(ctx context.Context, payload RegisterPdfPayload) ([]byte, error) {
			return []byte("reg-pdf"), nil
		})
		_ = srv.handleRegisterPdfTask(ctx, validTask)
	})

	// 4. Verbale task handlers
	t.Run("VerbalePdfTask", func(t *testing.T) {
		// Invalid JSON
		badTask := asynq.NewTask(TypeVerbalePdf, []byte("{invalid"))
		err := srv.handleVerbalePdfTask(ctx, badTask)
		assert.Error(t, err)

		// Nil handler
		validPayload := VerbalePdfPayload{JobID: "job-v1", VerbaleID: "verb1"}
		validTask, _ := NewVerbalePdfTask(validPayload)
		err = srv.handleVerbalePdfTask(ctx, validTask)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no handler function")

		// Handler returns error
		srv.SetVerbaleHandler(func(ctx context.Context, payload VerbalePdfPayload) ([]byte, error) {
			return nil, errors.New("verbale pdf failed")
		})
		err = srv.handleVerbalePdfTask(ctx, validTask)
		assert.Error(t, err)
		assert.Equal(t, "verbale pdf failed", err.Error())

		// Handler succeeds
		srv.SetVerbaleHandler(func(ctx context.Context, payload VerbalePdfPayload) ([]byte, error) {
			return []byte("verb-pdf"), nil
		})
		_ = srv.handleVerbalePdfTask(ctx, validTask)
	})
}
