package notifications

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	SaveToken(ctx context.Context, token *PushToken) error
	GetUserTokens(ctx context.Context, userID string) ([]PushToken, error)
	DeleteToken(ctx context.Context, userID, deviceToken string) error

	CreateDBNotification(ctx context.Context, n *DBNotification) error
	ListDBNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]DBNotification, error)
	MarkAsRead(ctx context.Context, userID, notificationID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveToken(ctx context.Context, token *PushToken) error {
	token.ID = uuid.New().String()
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	if token.Platform == "" {
		token.Platform = "android"
	}

	query := `
		INSERT INTO user_push_tokens (id, user_id, device_token, platform, created_at, updated_at)
		VALUES ($1, $2::uuid, $3, $4, $5, $6)
		ON CONFLICT (user_id, device_token) DO UPDATE
		SET platform = EXCLUDED.platform, updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.DeviceToken, token.Platform, token.CreatedAt, token.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetUserTokens(ctx context.Context, userID string) ([]PushToken, error) {
	query := `
		SELECT id, user_id, device_token, platform, created_at, updated_at
		FROM user_push_tokens
		WHERE user_id = $1::uuid
		ORDER BY updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []PushToken
	for rows.Next() {
		var t PushToken
		if err := rows.Scan(&t.ID, &t.UserID, &t.DeviceToken, &t.Platform, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, rows.Err()
}

func (r *PostgresRepository) DeleteToken(ctx context.Context, userID, deviceToken string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM user_push_tokens WHERE user_id = $1::uuid AND device_token = $2", userID, deviceToken)
	return err
}

func (r *PostgresRepository) CreateDBNotification(ctx context.Context, n *DBNotification) error {
	n.ID = uuid.New().String()
	n.CreatedAt = time.Now()
	if n.Type == "" {
		n.Type = "info"
	}

	query := `
		INSERT INTO db_notifications (id, user_id, title, body, type, payload, created_at)
		VALUES ($1, $2::uuid, $3, $4, $5, $6, $7)
	`
	payloadJSON, _ := json.Marshal(n.Payload)
	_, err := r.db.ExecContext(ctx, query, n.ID, n.UserID, n.Title, n.Body, n.Type, payloadJSON, n.CreatedAt)
	return err
}

func (r *PostgresRepository) ListDBNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]DBNotification, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, user_id, title, body, type, COALESCE(payload, '{}'::jsonb), read_at, created_at
		FROM db_notifications
		WHERE user_id = $1::uuid AND ($2 = false OR read_at IS NULL)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(ctx, query, userID, unreadOnly, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []DBNotification
	for rows.Next() {
		var n DBNotification
		var payloadRaw []byte
		var readAt sql.NullTime
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.Type, &payloadRaw, &readAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		if readAt.Valid {
			n.ReadAt = &readAt.Time
		}
		if len(payloadRaw) > 0 {
			_ = json.Unmarshal(payloadRaw, &n.Payload)
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) MarkAsRead(ctx context.Context, userID, notificationID string) error {
	query := `UPDATE db_notifications SET read_at = NOW() WHERE id = $1::uuid AND user_id = $2::uuid`
	_, err := r.db.ExecContext(ctx, query, notificationID, userID)
	return err
}

func (r *PostgresRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	query := `UPDATE db_notifications SET read_at = NOW() WHERE user_id = $1::uuid AND read_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
