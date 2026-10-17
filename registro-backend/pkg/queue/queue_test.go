package queue

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryQueue_EnqueueAndProcess(t *testing.T) {
	q := NewMemoryQueue()
	var processedCount int32

	q.RegisterHandler("TEST_JOB", func(ctx context.Context, task *Task) error {
		atomic.AddInt32(&processedCount, 1)
		return nil
	})

	q.Start()
	defer q.Stop()

	task, err := q.Enqueue(context.Background(), "TEST_JOB", map[string]string{"foo": "bar"})
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&processedCount) == 1
	}, 1*time.Second, 20*time.Millisecond)
}

func TestMemoryQueue_EnqueueBytesPayload(t *testing.T) {
	q := NewMemoryQueue()
	var receivedPayload []byte
	processed := make(chan struct{})

	q.RegisterHandler("RAW_BYTES_JOB", func(ctx context.Context, task *Task) error {
		receivedPayload = task.Payload
		close(processed)
		return nil
	})

	q.Start()
	defer q.Stop()

	payload := []byte(`{"message":"direct-bytes"}`)
	task, err := q.Enqueue(context.Background(), "RAW_BYTES_JOB", payload)
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)

	select {
	case <-processed:
		assert.Equal(t, payload, receivedPayload)
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for raw bytes job processing")
	}
}

func TestMemoryQueue_UnhandledTaskDoesNotBlockOrPanic(t *testing.T) {
	q := NewMemoryQueue()
	q.Start()
	defer q.Stop()

	// Enqueue a task type without any handler
	task, err := q.Enqueue(context.Background(), "UNKNOWN_JOB", map[string]string{"k": "v"})
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)

	// Ensure queue continues functioning by then enqueuing a handled task
	var handled bool
	handledChan := make(chan struct{})
	q.RegisterHandler("VALID_JOB", func(ctx context.Context, task *Task) error {
		handled = true
		close(handledChan)
		return nil
	})

	_, err = q.Enqueue(context.Background(), "VALID_JOB", 123)
	require.NoError(t, err)

	select {
	case <-handledChan:
		assert.True(t, handled)
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for valid job processing")
	}
}

func TestMemoryQueue_QueueFull(t *testing.T) {
	q := NewMemoryQueue()
	// Do NOT call q.Start(), so workers don't drain the channel

	// Fill the 1000 items buffer
	for i := 0; i < 1000; i++ {
		_, err := q.Enqueue(context.Background(), "FILL_JOB", i)
		require.NoError(t, err)
	}

	// 1001st task should return queue full error
	_, err := q.Enqueue(context.Background(), "OVERFLOW_JOB", "overflow")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "queue is full")
}

func TestNewQueue_FallbackToMemory(t *testing.T) {
	// Empty URL -> MemoryQueue
	q1 := NewQueue("")
	assert.IsType(t, &MemoryQueue{}, q1)

	// Invalid / unreachable Redis URL -> falls back gracefully to MemoryQueue
	q2 := NewQueue("redis://localhost:99999/0")
	assert.IsType(t, &MemoryQueue{}, q2)
}
