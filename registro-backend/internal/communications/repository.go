package communications

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, msg *Message) error
	List(ctx context.Context, userID string) ([]*Message, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, msg *Message) error {
	msg.ID = uuid.New().String()
	msg.CreatedAt = time.Now()

	query := `
		INSERT INTO communications (id, sender_id, receiver_ids, subject, body, type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	// receiver_ids stored as text array for simplicity in this MVP, or easier via pq.Array
	_, err := r.db.ExecContext(ctx, query,
		msg.ID, msg.SenderID, pq.Array(msg.ReceiverIDs), msg.Subject, msg.Body, msg.Type, msg.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, userID string) ([]*Message, error) {
	// Very basic check: Sender or Included in Receivers
	query := `
		SELECT id, sender_id, receiver_ids, subject, body, type, created_at
		FROM communications
		WHERE sender_id = $1 OR $1 = ANY(receiver_ids)
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*Message
	for rows.Next() {
		m := &Message{}
		var receivers []string
		if err := rows.Scan(
			&m.ID, &m.SenderID, pq.Array(&receivers), &m.Subject, &m.Body, &m.Type, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		m.ReceiverIDs = receivers
		msgs = append(msgs, m)
	}
	return msgs, nil
}
