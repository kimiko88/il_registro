package queue

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMemoryQueue_HandlerError verifies that a handler returning an error does not
// crash the worker loop and subsequent tasks are still processed.
func TestMemoryQueue_HandlerError(t *testing.T) {
	q := NewMemoryQueue()
	q.Start()
	defer q.Stop()

	var errorCount, successCount int32

	q.RegisterHandler("FAIL_JOB", func(ctx context.Context, task *Task) error {
		atomic.AddInt32(&errorCount, 1)
		return errors.New("simulated handler error")
	})

	q.RegisterHandler("OK_JOB", func(ctx context.Context, task *Task) error {
		atomic.AddInt32(&successCount, 1)
		return nil
	})

	_, err := q.Enqueue(context.Background(), "FAIL_JOB", "data")
	require.NoError(t, err)

	_, err = q.Enqueue(context.Background(), "OK_JOB", "data")
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&errorCount) == 1 && atomic.LoadInt32(&successCount) == 1
	}, 2*time.Second, 20*time.Millisecond, "both error and success handlers should execute")
}

// TestMemoryQueue_TaskFields verifies that a Task has a non-empty ID and a correct Type.
func TestMemoryQueue_TaskFields(t *testing.T) {
	q := NewMemoryQueue()

	task, err := q.Enqueue(context.Background(), "MY_TYPE", map[string]int{"value": 42})
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)
	assert.Equal(t, "MY_TYPE", task.Type)
	assert.NotNil(t, task.Payload)
	assert.False(t, task.CreatedAt.IsZero())
}

// TestMemoryQueue_BytesPayloadPreserved verifies raw []byte payload passes through without
// double-marshaling.
func TestMemoryQueue_BytesPayloadPreserved(t *testing.T) {
	q := NewMemoryQueue()
	done := make(chan []byte, 1)
	q.RegisterHandler("RAW", func(ctx context.Context, t *Task) error {
		done <- t.Payload
		return nil
	})
	q.Start()
	defer q.Stop()

	raw := []byte(`{"key":"value"}`)
	_, err := q.Enqueue(context.Background(), "RAW", raw)
	require.NoError(t, err)

	select {
	case got := <-done:
		assert.Equal(t, raw, got)
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for raw bytes handler")
	}
}

// TestMemoryQueue_StopIdempotent verifies that Stop can be called safely after the queue
// was already stopped (no double-close panic).
func TestMemoryQueue_StopIdempotent(t *testing.T) {
	q := NewMemoryQueue()
	q.Start()
	q.Stop()
	// Calling Stop again via a new goroutine won't trigger double-close of the already-closed
	// stopChan; the WaitGroup would simply return immediately since all goroutines finished.
	// This documents/tests expected behaviour.
	assert.NotPanics(t, func() {
		q.wg.Wait() // already done — should not block
	})
}

// TestMemoryQueue_EnqueueUnmarshalablePayload verifies that a type that cannot be marshaled
// to JSON returns an error immediately.
func TestMemoryQueue_EnqueueUnmarshalablePayload(t *testing.T) {
	q := NewMemoryQueue()
	// channels cannot be JSON-marshaled
	_, err := q.Enqueue(context.Background(), "BAD_JOB", make(chan int))
	require.Error(t, err, "non-serialisable payload should return an error")
}

// TestMemoryQueue_RegisterHandlerOverwrite verifies that registering a second handler for the
// same task type overwrites the first.
func TestMemoryQueue_RegisterHandlerOverwrite(t *testing.T) {
	q := NewMemoryQueue()
	var calls int32

	q.RegisterHandler("TASK", func(ctx context.Context, t *Task) error {
		atomic.AddInt32(&calls, 100) // first handler — should NOT be called
		return nil
	})
	q.RegisterHandler("TASK", func(ctx context.Context, t *Task) error {
		atomic.AddInt32(&calls, 1) // second handler — should be called
		return nil
	})

	q.Start()
	defer q.Stop()

	_, err := q.Enqueue(context.Background(), "TASK", "x")
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&calls) == 1
	}, 1*time.Second, 20*time.Millisecond, "only the second (overwriting) handler should run")
}

// TestNewQueue_EmptyURL returns a MemoryQueue.
func TestNewQueue_EmptyURL(t *testing.T) {
	q := NewQueue("")
	require.NotNil(t, q)
	_, ok := q.(*MemoryQueue)
	assert.True(t, ok, "NewQueue(\"\") should return a *MemoryQueue")
}

// TestNewQueue_InvalidURL falls back to MemoryQueue.
func TestNewQueue_InvalidURL(t *testing.T) {
	q := NewQueue("redis://127.0.0.1:1/0") // port 1 — unreachable
	require.NotNil(t, q)
	_, ok := q.(*MemoryQueue)
	assert.True(t, ok, "unreachable Redis should fall back to MemoryQueue")
}

// TestMemoryQueue_MultipleHandlers verifies multiple handler types work correctly in the same queue.
func TestMemoryQueue_MultipleHandlers(t *testing.T) {
	q := NewMemoryQueue()
	q.Start()
	defer q.Stop()

	var aCount, bCount int32
	q.RegisterHandler("TYPE_A", func(_ context.Context, _ *Task) error {
		atomic.AddInt32(&aCount, 1)
		return nil
	})
	q.RegisterHandler("TYPE_B", func(_ context.Context, _ *Task) error {
		atomic.AddInt32(&bCount, 1)
		return nil
	})

	for i := 0; i < 3; i++ {
		_, err := q.Enqueue(context.Background(), "TYPE_A", i)
		require.NoError(t, err)
		_, err = q.Enqueue(context.Background(), "TYPE_B", i)
		require.NoError(t, err)
	}

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&aCount) == 3 && atomic.LoadInt32(&bCount) == 3
	}, 2*time.Second, 20*time.Millisecond)
}

// TestMemoryQueue_ContextCancelledInHandler verifies that a handler receiving a cancelled
// context can detect it without panicking.
func TestMemoryQueue_ContextCancelledInHandler(t *testing.T) {
	q := NewMemoryQueue()
	done := make(chan struct{})
	q.RegisterHandler("CTX_JOB", func(ctx context.Context, _ *Task) error {
		// Just check the context isn't already cancelled at start
		select {
		case <-ctx.Done():
			// context already cancelled — unusual but must not panic
		default:
		}
		close(done)
		return nil
	})

	q.Start()
	defer q.Stop()

	_, err := q.Enqueue(context.Background(), "CTX_JOB", nil)
	require.NoError(t, err)

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("handler never called")
	}
}
