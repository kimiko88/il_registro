package accessibility

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, f *AccessibilityFeedback) error
	List(ctx context.Context, schoolID, status string, limit, offset int) ([]AccessibilityFeedback, int, error)
	GetByID(ctx context.Context, id string) (*AccessibilityFeedback, error)
	UpdateStatus(ctx context.Context, id, status, responseNotes string) error
	GetUserPreferences(ctx context.Context, userID string) (string, error)
	UpsertUserPreferences(ctx context.Context, userID string, settingsJSON string) error
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, f *AccessibilityFeedback) error {
	query := `
		INSERT INTO accessibility_feedbacks (
			protocol_number, name, email, barrier_type, description,
			user_id, school_id, user_agent, ip_address, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		f.ProtocolNumber, f.Name, f.Email, f.BarrierType, f.Description,
		f.UserID, f.SchoolID, f.UserAgent, f.IPAddress, f.Status,
	).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
}

func (r *postgresRepository) List(ctx context.Context, schoolID, status string, limit, offset int) ([]AccessibilityFeedback, int, error) {
	whereClauses := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1

	if schoolID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("school_id = $%d", argIndex))
		args = append(args, schoolID)
		argIndex++
	}

	if status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM accessibility_feedbacks WHERE %s", fmt.Sprintf("%s", whereClauses[0]))
	if len(whereClauses) > 1 {
		countQuery = fmt.Sprintf("SELECT COUNT(*) FROM accessibility_feedbacks WHERE %s", joinClauses(whereClauses))
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT id, protocol_number, name, email, barrier_type, description,
		       user_id, school_id, user_agent, ip_address, status, response_notes,
		       resolved_at, created_at, updated_at
		FROM accessibility_feedbacks
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, joinClauses(whereClauses), argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feedbacks []AccessibilityFeedback
	for rows.Next() {
		var f AccessibilityFeedback
		if err := rows.Scan(
			&f.ID, &f.ProtocolNumber, &f.Name, &f.Email, &f.BarrierType, &f.Description,
			&f.UserID, &f.SchoolID, &f.UserAgent, &f.IPAddress, &f.Status, &f.ResponseNotes,
			&f.ResolvedAt, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		feedbacks = append(feedbacks, f)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}


func (r *postgresRepository) GetByID(ctx context.Context, id string) (*AccessibilityFeedback, error) {
	query := `
		SELECT id, protocol_number, name, email, barrier_type, description,
		       user_id, school_id, user_agent, ip_address, status, response_notes,
		       resolved_at, created_at, updated_at
		FROM accessibility_feedbacks
		WHERE id = $1
	`
	var f AccessibilityFeedback
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID, &f.ProtocolNumber, &f.Name, &f.Email, &f.BarrierType, &f.Description,
		&f.UserID, &f.SchoolID, &f.UserAgent, &f.IPAddress, &f.Status, &f.ResponseNotes,
		&f.ResolvedAt, &f.CreatedAt, &f.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, id, status, responseNotes string) error {
	query := `
		UPDATE accessibility_feedbacks
		SET status = $2, response_notes = $3, updated_at = NOW(),
		    resolved_at = CASE WHEN $2 = 'resolved' THEN NOW() ELSE resolved_at END
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, status, responseNotes)
	return err
}

func (r *postgresRepository) GetUserPreferences(ctx context.Context, userID string) (string, error) {
	query := `SELECT settings::text FROM user_accessibility_preferences WHERE user_id = $1`
	var settings string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&settings)
	if err == sql.ErrNoRows {
		return "{}", nil
	}
	if err != nil {
		return "", err
	}
	return settings, nil
}

func (r *postgresRepository) UpsertUserPreferences(ctx context.Context, userID string, settingsJSON string) error {
	query := `
		INSERT INTO user_accessibility_preferences (user_id, settings, updated_at)
		VALUES ($1, $2::jsonb, NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET settings = EXCLUDED.settings, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, userID, settingsJSON)
	return err
}

func joinClauses(clauses []string) string {
	res := ""
	for i, c := range clauses {
		if i > 0 {
			res += " AND "
		}
		res += c
	}
	return res
}
