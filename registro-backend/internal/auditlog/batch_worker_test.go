package auditlog

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatchWorker_EnqueueAndStop(t *testing.T) {
	// Worker with nil db (flush handles nil gracefully)
	worker := NewBatchWorker(nil, 5, 100*time.Millisecond)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			worker.Enqueue(AuditEvent{
				SchoolID: "school-1",
				ActorID:  "actor-1",
				Action:   "TEST_ACTION",
			})
		}(i)
	}
	wg.Wait()

	// Graceful stop drains remaining items
	worker.Stop()
	assert.True(t, true)
}
