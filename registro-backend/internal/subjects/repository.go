package subjects

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, subject *Subject) error
	List(ctx context.Context, schoolID string) ([]Subject, error)
	Get(ctx context.Context, id string) (*Subject, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, subject *Subject) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, s *Subject) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.CreatedAt = time.Now()

	query := `
		INSERT INTO subjects (id, school_id, name, code, description, is_mandatory, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.SchoolID, s.Name, s.Code, s.Description, s.IsMandatory, s.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, schoolID string) ([]Subject, error) {
	query := `
		SELECT id, school_id, name, code, description, is_mandatory, created_at
		FROM subjects
		WHERE school_id = $1
		ORDER BY name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var s Subject
		var code, desc sql.NullString
		err := rows.Scan(
			&s.ID, &s.SchoolID, &s.Name, &code, &desc, &s.IsMandatory, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		s.Code = code.String
		s.Description = desc.String
		subjects = append(subjects, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subjects, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*Subject, error) {
	query := `SELECT id, school_id, name, code, description, is_mandatory, created_at FROM subjects WHERE id = $1`
	var s Subject
	var code, desc sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.SchoolID, &s.Name, &code, &desc, &s.IsMandatory, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	s.Code = code.String
	s.Description = desc.String
	return &s, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM subjects WHERE id = $1", id)
	return err
}

func (r *PostgresRepository) Update(ctx context.Context, s *Subject) error {
	query := `UPDATE subjects SET name=$1, code=$2, description=$3, is_mandatory=$4 WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query, s.Name, s.Code, s.Description, s.IsMandatory, s.ID)
	return err
}
