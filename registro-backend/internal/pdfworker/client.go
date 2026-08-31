package pdfworker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// JobStatus represents the state of an async PDF job
type JobStatus struct {
	JobID       string    `json:"job_id"`
	Status      string    `json:"status"` // pending, processing, completed, failed
	TaskType    string    `json:"task_type"`
	DownloadURL string    `json:"download_url,omitempty"`
	Error       string    `json:"error,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Client struct {
	asynqClient *asynq.Client
	rdb         *redis.Client
}

func NewClient(redisAddr string) *Client {
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	return &Client{
		asynqClient: asynqClient,
		rdb:         rdb,
	}
}

func (c *Client) Close() error {
	if c.asynqClient != nil {
		c.asynqClient.Close()
	}
	if c.rdb != nil {
		return c.rdb.Close()
	}
	return nil
}

// EnqueueScrutinyPdf enqueues a scrutiny PDF job and stores initial pending status in Redis
func (c *Client) EnqueueScrutinyPdf(ctx context.Context, payload ScrutinyPdfPayload) (*JobStatus, error) {
	task, err := NewScrutinyPdfTask(payload)
	if err != nil {
		return nil, err
	}

	info, err := c.asynqClient.EnqueueContext(ctx, task, asynq.Queue("critical"), asynq.Retention(1*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue task: %w", err)
	}

	jobStatus := &JobStatus{
		JobID:     payload.JobID,
		Status:    "pending",
		TaskType:  TypeScrutinyPdf,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.UpdateJobStatus(ctx, jobStatus); err != nil {
		return nil, err
	}

	_ = info
	return jobStatus, nil
}

// EnqueueReportCardPdf enqueues a report card PDF job and stores initial pending status in Redis
func (c *Client) EnqueueReportCardPdf(ctx context.Context, payload ReportCardPdfPayload) (*JobStatus, error) {
	task, err := NewReportCardPdfTask(payload)
	if err != nil {
		return nil, err
	}

	info, err := c.asynqClient.EnqueueContext(ctx, task, asynq.Queue("default"), asynq.Retention(1*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue report card task: %w", err)
	}

	jobStatus := &JobStatus{
		JobID:     payload.JobID,
		Status:    "pending",
		TaskType:  TypeReportCardPdf,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.UpdateJobStatus(ctx, jobStatus); err != nil {
		return nil, err
	}

	_ = info
	return jobStatus, nil
}

// EnqueueRegisterPdf enqueues a register PDF job and stores initial pending status in Redis
func (c *Client) EnqueueRegisterPdf(ctx context.Context, payload RegisterPdfPayload) (*JobStatus, error) {
	task, err := NewRegisterPdfTask(payload)
	if err != nil {
		return nil, err
	}

	info, err := c.asynqClient.EnqueueContext(ctx, task, asynq.Queue("default"), asynq.Retention(1*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue register task: %w", err)
	}

	jobStatus := &JobStatus{
		JobID:     payload.JobID,
		Status:    "pending",
		TaskType:  TypeRegisterPdf,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.UpdateJobStatus(ctx, jobStatus); err != nil {
		return nil, err
	}

	_ = info
	return jobStatus, nil
}

// EnqueueVerbalePdf enqueues a verbale PDF job and stores initial pending status in Redis
func (c *Client) EnqueueVerbalePdf(ctx context.Context, payload VerbalePdfPayload) (*JobStatus, error) {
	task, err := NewVerbalePdfTask(payload)
	if err != nil {
		return nil, err
	}

	info, err := c.asynqClient.EnqueueContext(ctx, task, asynq.Queue("default"), asynq.Retention(1*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue verbale task: %w", err)
	}

	jobStatus := &JobStatus{
		JobID:     payload.JobID,
		Status:    "pending",
		TaskType:  TypeVerbalePdf,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.UpdateJobStatus(ctx, jobStatus); err != nil {
		return nil, err
	}

	_ = info
	return jobStatus, nil
}

// UpdateJobStatus updates job status JSON in Redis with 24h expiration
func (c *Client) UpdateJobStatus(ctx context.Context, status *JobStatus) error {
	status.UpdatedAt = time.Now()
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("pdf_job:%s", status.JobID)
	return c.rdb.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetJobStatus retrieves job status from Redis
func (c *Client) GetJobStatus(ctx context.Context, jobID string) (*JobStatus, error) {
	key := fmt.Sprintf("pdf_job:%s", jobID)
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("job not found")
		}
		return nil, err
	}

	var status JobStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, err
	}
	return &status, nil
}
