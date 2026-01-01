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
	AssignSubject(ctx context.Context, classID string, subjectID string, teacherUserID *string, hours float64) error
	UnassignSubject(ctx context.Context, assignmentID string) error
	GetClassSubjects(ctx context.Context, classID string) ([]ClassSubject, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// ... existing ListByTeacher ...

func (r *PostgresRepository) AssignSubject(ctx context.Context, classID string, subjectID string, teacherUserID *string, hours float64) error {
	var teacherID *string
	if teacherUserID != nil {
		// Resolve teacher_id from teachers table
		// We assume the teacher belongs to same school as class, or just lookup by user_id globally (user_id is unique enough usually, but strictly school-bound).
		// Better: look up teacher by user_id.
		query := `SELECT id FROM teachers WHERE user_id = $1`
		var tid string
		err := r.db.QueryRowContext(ctx, query, *teacherUserID).Scan(&tid)
		if err != nil {
			return errors.New("teacher profile not found for user")
		}
		teacherID = &tid
	}

	// Insert
	id := uuid.New().String()
	// Constraint is (class_id, subject_id, teacher_id).
	// If we want to replace the teacher for a subject, we should probably delete old ones for this subject?
	// But let's support additive for now as per schema.
	query := `INSERT INTO class_subjects (id, class_id, subject_id, teacher_id, hours_per_week)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, id, classID, subjectID, teacherID, hours)
	return err
}

func (r *PostgresRepository) UnassignSubject(ctx context.Context, assignmentID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM class_subjects WHERE id = $1", assignmentID)
	return err
}

func (r *PostgresRepository) GetClassSubjects(ctx context.Context, classID string) ([]ClassSubject, error) {
	query := `
		SELECT cs.id, cs.class_id, cs.subject_id, s.name, cs.teacher_id, u.last_name, u.first_name, cs.hours_per_week
		FROM class_subjects cs
		JOIN subjects s ON cs.subject_id = s.id
		LEFT JOIN teachers t ON cs.teacher_id = t.id
		LEFT JOIN users u ON t.user_id = u.id
		WHERE cs.class_id = $1
		ORDER BY s.name
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ClassSubject
	for rows.Next() {
		var cs ClassSubject
		var tLen, tFirst sql.NullString
		err := rows.Scan(
			&cs.ID, &cs.ClassID, &cs.SubjectID, &cs.SubjectName, &cs.TeacherID, &tLen, &tFirst, &cs.HoursPerWeek,
		)
		if err != nil {
			return nil, err
		}
		if tLen.Valid {
			cs.TeacherName = tLen.String + " " + tFirst.String
		}
		results = append(results, cs)
	}
	return results, nil
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
