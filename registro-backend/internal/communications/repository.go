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
	Delete(ctx context.Context, id string) error
	Sign(ctx context.Context, communicationID string, userID string) error
	GetSignatures(ctx context.Context, communicationID string) ([]string, error)
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
	_, err := r.db.ExecContext(ctx, query,
		msg.ID, msg.SenderID, pq.Array(msg.ReceiverIDs), msg.Subject, msg.Body, msg.Type, msg.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, userID string) ([]*Message, error) {
	query := `
		SELECT c.id, c.sender_id, c.receiver_ids, c.subject, c.body, c.type, c.created_at,
		       EXISTS(SELECT 1 FROM communication_signatures cs WHERE cs.communication_id = c.id AND cs.user_id = $1) AS is_signed
		FROM communications c
		WHERE c.sender_id::text = $1 OR $1 = ANY(c.receiver_ids)
		ORDER BY c.created_at DESC
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
			&m.ID, &m.SenderID, pq.Array(&receivers), &m.Subject, &m.Body, &m.Type, &m.CreatedAt, &m.IsSigned,
		); err != nil {
			return nil, err
		}
		m.ReceiverIDs = receivers
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM communications WHERE id = $1", id)
	return err
}

func (r *PostgresRepository) Sign(ctx context.Context, communicationID string, userID string) error {
	query := `
		INSERT INTO communication_signatures (communication_id, user_id, signed_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (communication_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, communicationID, userID)
	return err
}

func (r *PostgresRepository) GetSignatures(ctx context.Context, communicationID string) ([]string, error) {
	query := `
		SELECT COALESCE(u.first_name || ' ' || u.last_name, '') AS name
		FROM communication_signatures cs
		JOIN users u ON cs.user_id = u.id
		WHERE cs.communication_id = $1
		ORDER BY cs.signed_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, communicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}
