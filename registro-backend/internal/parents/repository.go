package parents

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetChildrenByParentUserID(ctx context.Context, parentUserID string) ([]string, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetChildrenByParentUserID(ctx context.Context, parentUserID string) ([]string, error) {
	query := `
		SELECT u.id::text
		FROM student_parents sp
		JOIN parents p ON sp.parent_id = p.id
		JOIN students s ON sp.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE p.user_id = $1::uuid
	`
	rows, err := r.db.QueryContext(ctx, query, parentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, rows.Err()
}
