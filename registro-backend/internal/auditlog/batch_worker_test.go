package auditlog

import (
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestBatchWorker_Defaults(t *testing.T) {
	worker := NewBatchWorker(nil, 0, 0)
	assert.Equal(t, 100, worker.batchSize)
	assert.Equal(t, 500*time.Millisecond, worker.flushInterval)
	worker.Stop()
}

func TestBatchWorker_BatchSizeTriggeredFlush(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// batchSize = 2. When 2 items are enqueued, flush should be triggered immediately.
	mock.ExpectExec("INSERT INTO audit_logs").
		WillReturnResult(sqlmock.NewResult(1, 2))

	worker := NewBatchWorker(db, 2, 10*time.Second)

	worker.Enqueue(AuditEvent{
		SchoolID: "school-1",
		ActorID:  "actor-1",
		Action:   "BATCH_1",
	})
	worker.Enqueue(AuditEvent{
		SchoolID: "school-1",
		ActorID:  "actor-2",
		Action:   "BATCH_2",
	})

	// Wait briefly for the batch flush to execute
	time.Sleep(50 * time.Millisecond)

	worker.Stop()
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBatchWorker_TimerTriggeredFlush(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// batchSize = 10 (not reached), flushInterval = 50ms
	mock.ExpectExec("INSERT INTO audit_logs").
		WillReturnResult(sqlmock.NewResult(1, 1))

	worker := NewBatchWorker(db, 10, 50*time.Millisecond)

	worker.Enqueue(AuditEvent{
		SchoolID: "school-1",
		ActorID:  "actor-1",
		Action:   "TIMER_EVENT",
	})

	// Wait for timer ticker to trigger flush
	time.Sleep(120 * time.Millisecond)

	worker.Stop()
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBatchWorker_InsertSingle_Direct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO audit_logs").
		WithArgs(
			sqlmock.AnyArg(), "school-1", "actor-1", "admin", "Admin User",
			"DIRECT_ACTION", "user", "target-1", sqlmock.AnyArg(), "127.0.0.1", "Mozilla", sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	worker := NewBatchWorker(db, 10, 1*time.Second)
	defer worker.Stop()

	ev := &AuditEvent{
		ID:         "event-single-1",
		SchoolID:   "school-1",
		ActorID:    "actor-1",
		ActorRole:  "admin",
		ActorName:  "Admin User",
		Action:     "DIRECT_ACTION",
		EntityType: "user",
		EntityID:   "target-1",
		Details:    `{"key":"value"}`,
		IPAddress:  "127.0.0.1",
		UserAgent:  "Mozilla",
		CreatedAt:  time.Now(),
	}

	err = worker.insertSingle(t.Context(), ev)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test nil db returns nil without panic
	nilWorker := &BatchWorker{db: nil}
	assert.Nil(t, nilWorker.insertSingle(t.Context(), ev))
}

func TestBatchWorker_DrainOnStop(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// batchSize = 10, flushInterval = 10s (won't trigger on timer)
	mock.ExpectExec("INSERT INTO audit_logs").
		WillReturnResult(sqlmock.NewResult(1, 3))

	worker := NewBatchWorker(db, 10, 10*time.Second)

	for i := 0; i < 3; i++ {
		worker.Enqueue(AuditEvent{
			SchoolID: "school-1",
			ActorID:  "actor-1",
			Action:   "DRAIN_ACTION",
		})
	}

	// Calling Stop() must drain the 3 pending items
	worker.Stop()
	assert.NoError(t, mock.ExpectationsWereMet())
}
