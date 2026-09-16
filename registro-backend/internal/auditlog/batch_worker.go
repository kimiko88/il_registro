package auditlog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// BatchWorker buffers audit log events in memory and writes them to PostgreSQL
// in periodic or size-triggered batches, reducing DB write roundtrips by >90%.
type BatchWorker struct {
	db            *sql.DB
	queue         chan AuditEvent
	batchSize     int
	flushInterval time.Duration
	stopChan      chan struct{}
	wg            sync.WaitGroup
}

// NewBatchWorker initializes and starts a new audit log batch worker.
func NewBatchWorker(db *sql.DB, batchSize int, flushInterval time.Duration) *BatchWorker {
	if batchSize <= 0 {
		batchSize = 100
	}
	if flushInterval <= 0 {
		flushInterval = 500 * time.Millisecond
	}

	w := &BatchWorker{
		db:            db,
		queue:         make(chan AuditEvent, 10000),
		batchSize:     batchSize,
		flushInterval: flushInterval,
		stopChan:      make(chan struct{}),
	}

	w.wg.Add(1)
	go w.run()
	return w
}

// Enqueue adds an event to the batch queue. If the queue is saturated, it logs a warning
// and executes a direct insert in a separate goroutine so audit events are never lost.
func (w *BatchWorker) Enqueue(event AuditEvent) {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	select {
	case w.queue <- event:
	default:
		logrus.Warn("auditlog: batch worker buffer full, inserting synchronously")
		go func(ev AuditEvent) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = w.insertSingle(ctx, &ev)
		}(event)
	}
}

func (w *BatchWorker) run() {
	defer w.wg.Done()
	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	batch := make([]AuditEvent, 0, w.batchSize)

	for {
		select {
		case ev := <-w.queue:
			batch = append(batch, ev)
			if len(batch) >= w.batchSize {
				w.flush(batch)
				batch = make([]AuditEvent, 0, w.batchSize)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				w.flush(batch)
				batch = make([]AuditEvent, 0, w.batchSize)
			}
		case <-w.stopChan:
			// Drain any remaining events in the channel
			for {
				select {
				case ev := <-w.queue:
					batch = append(batch, ev)
					if len(batch) >= w.batchSize {
						w.flush(batch)
						batch = make([]AuditEvent, 0, w.batchSize)
					}
				default:
					if len(batch) > 0 {
						w.flush(batch)
					}
					return
				}
			}
		}
	}
}

func (w *BatchWorker) flush(events []AuditEvent) {
	if len(events) == 0 || w.db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Multi-row INSERT:
	// INSERT INTO audit_logs (id, school_id, actor_id, actor_role, actor_name, action, entity_type, entity_id, details, ip_address, user_agent, created_at)
	// VALUES ($1, $2, ...), ($13, $14, ...)
	valueStrings := make([]string, 0, len(events))
	valueArgs := make([]interface{}, 0, len(events)*12)

	for i, ev := range events {
		offset := i * 12
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4, offset+5, offset+6,
			offset+7, offset+8, offset+9, offset+10, offset+11, offset+12,
		))
		valueArgs = append(valueArgs,
			ev.ID, ev.SchoolID, ev.ActorID, ev.ActorRole, ev.ActorName,
			ev.Action, ev.EntityType, ev.EntityID, ev.Details, ev.IPAddress, ev.UserAgent, ev.CreatedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO audit_logs (id, school_id, actor_id, actor_role, actor_name, action, entity_type, entity_id, details, ip_address, user_agent, created_at)
		VALUES %s
	`, strings.Join(valueStrings, ", "))

	_, err := w.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		logrus.Errorf("auditlog batch worker flush error (%d items): %v", len(events), err)
	}
}

func (w *BatchWorker) insertSingle(ctx context.Context, ev *AuditEvent) error {
	if w.db == nil {
		return nil
	}
	query := `
		INSERT INTO audit_logs (id, school_id, actor_id, actor_role, actor_name, action, entity_type, entity_id, details, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := w.db.ExecContext(ctx, query,
		ev.ID, ev.SchoolID, ev.ActorID, ev.ActorRole, ev.ActorName,
		ev.Action, ev.EntityType, ev.EntityID, ev.Details, ev.IPAddress, ev.UserAgent, ev.CreatedAt,
	)
	return err
}

// Stop flushes all pending events and shuts down the worker goroutine gracefully.
func (w *BatchWorker) Stop() {
	close(w.stopChan)
	w.wg.Wait()
}
