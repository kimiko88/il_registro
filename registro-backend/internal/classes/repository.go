package classes

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, class *Class) error
	List(ctx context.Context, schoolID string) ([]Class, error)
	Get(ctx context.Context, id string) (*Class, error)
	Update(ctx context.Context, class *Class) error
	Delete(ctx context.Context, id string) error
	ListByTeacher(ctx context.Context, teacherUserID string) ([]Class, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// ... existing methods ...

func (r *PostgresRepository) ListByTeacher(ctx context.Context, teacherUserID string) ([]Class, error) {
	// Combine classes where user is coordinator OR assigned as teacher (via class_subjects)
	query := `
		SELECT DISTINCT c.id, c.school_id, c.name, COALESCE(c.section, ''), c.academic_year, c.coordinator_id, c.created_at, c.updated_at
		FROM classes c
		LEFT JOIN class_subjects cs ON c.id = cs.class_id
		LEFT JOIN teachers t ON cs.teacher_id = t.id
		WHERE t.user_id = $1 OR c.coordinator_id = $1
		ORDER BY c.name
	`
	rows, err := r.db.QueryContext(ctx, query, teacherUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		var coord sql.NullString
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.Name, &c.Section, &c.AcademicYear, &coord, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.CoordinatorID = coord.String
		classes = append(classes, c)
	}
	return classes, nil
}

func (r *PostgresRepository) Create(ctx context.Context, c *Class) error {
	c.ID = uuid.New().String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	query := `INSERT INTO classes (id, school_id, name, section, academic_year, coordinator_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.SchoolID, c.Name, c.Section, c.AcademicYear,
		sql.NullString{String: c.CoordinatorID, Valid: c.CoordinatorID != ""},
		c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, schoolID string) ([]Class, error) {
	query := `SELECT id, school_id, name, COALESCE(section, ''), academic_year, coordinator_id, created_at, updated_at 
	          FROM classes WHERE school_id = $1 ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		var coord sql.NullString
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.Name, &c.Section, &c.AcademicYear, &coord, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.CoordinatorID = coord.String
		classes = append(classes, c)
	}
	return classes, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*Class, error) {
	query := `SELECT id, school_id, name, COALESCE(section, ''), academic_year, coordinator_id, created_at, updated_at 
	          FROM classes WHERE id = $1`
	var c Class
	var coord sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.SchoolID, &c.Name, &c.Section, &c.AcademicYear, &coord, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.CoordinatorID = coord.String
	return &c, nil
}

func (r *PostgresRepository) Update(ctx context.Context, c *Class) error {
	c.UpdatedAt = time.Now()
	query := `UPDATE classes SET name=$1, section=$2, academic_year=$3, coordinator_id=$4, updated_at=$5 WHERE id=$6`
	res, err := r.db.ExecContext(ctx, query,
		c.Name, c.Section, c.AcademicYear,
		sql.NullString{String: c.CoordinatorID, Valid: c.CoordinatorID != ""},
		c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("class not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM classes WHERE id = $1", id)
	return err
}
