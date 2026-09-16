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
