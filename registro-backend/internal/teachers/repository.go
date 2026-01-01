package teachers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context, schoolID string) ([]Teacher, error)
	Get(ctx context.Context, id string) (*Teacher, error)
	GetByUserID(ctx context.Context, userID string) (*Teacher, error)
	GetBySubject(ctx context.Context, subjectID string) ([]Teacher, error)

	GetSubjects(ctx context.Context, teacherID string) ([]TeacherSubject, error)
	AssignSubject(ctx context.Context, teacherID, subjectID string) error
	RemoveSubject(ctx context.Context, teacherID, subjectID string) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context, schoolID string) ([]Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.school_id, t.hiring_date, t.qualification, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		WHERE t.school_id = $1
		ORDER BY u.last_name, u.first_name
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []Teacher
	for rows.Next() {
		var t Teacher
		var hDate sql.NullTime
		var qual sql.NullString
		err := rows.Scan(
			&t.ID, &t.UserID, &t.SchoolID, &hDate, &qual, &t.CreatedAt, &t.UpdatedAt,
			&t.FirstName, &t.LastName, &t.Email,
		)
		if err != nil {
			return nil, err
		}
		if hDate.Valid {
			d := hDate.Time.Format("2006-01-02")
			t.HiringDate = &d
		}
		t.Qualification = qual.String
		teachers = append(teachers, t)
	}
	return teachers, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.school_id, t.hiring_date, t.qualification, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = $1
	`
	var t Teacher
	var hDate sql.NullTime
	var qual sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.UserID, &t.SchoolID, &hDate, &qual, &t.CreatedAt, &t.UpdatedAt,
		&t.FirstName, &t.LastName, &t.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("teacher not found")
		}
		return nil, err
	}
	if hDate.Valid {
		d := hDate.Time.Format("2006-01-02")
		t.HiringDate = &d
	}
	t.Qualification = qual.String
	return &t, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*Teacher, error) {
	query := `SELECT id FROM teachers WHERE user_id = $1`
	var id string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *PostgresRepository) GetSubjects(ctx context.Context, teacherID string) ([]TeacherSubject, error) {
	query := `
		SELECT ts.id, ts.teacher_id, ts.subject_id, s.name, ts.created_at
		FROM teacher_subjects ts
		JOIN subjects s ON ts.subject_id = s.id
		WHERE ts.teacher_id = $1
		ORDER BY s.name
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TeacherSubject
	for rows.Next() {
		var ts TeacherSubject
		err := rows.Scan(&ts.ID, &ts.TeacherID, &ts.SubjectID, &ts.SubjectName, &ts.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, ts)
	}
	return results, nil
}

func (r *PostgresRepository) GetBySubject(ctx context.Context, subjectID string) ([]Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.school_id, t.hiring_date, t.qualification, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		JOIN teacher_subjects ts ON t.id = ts.teacher_id
		WHERE ts.subject_id = $1
		ORDER BY u.last_name, u.first_name
	`
	rows, err := r.db.QueryContext(ctx, query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []Teacher
	for rows.Next() {
		var t Teacher
		var hDate sql.NullTime
		var qual sql.NullString
		err := rows.Scan(
			&t.ID, &t.UserID, &t.SchoolID, &hDate, &qual, &t.CreatedAt, &t.UpdatedAt,
			&t.FirstName, &t.LastName, &t.Email,
		)
		if err != nil {
			return nil, err
		}
		if hDate.Valid {
			d := hDate.Time.Format("2006-01-02")
			t.HiringDate = &d
		}
		t.Qualification = qual.String
		teachers = append(teachers, t)
	}
	return teachers, nil
}

func (r *PostgresRepository) AssignSubject(ctx context.Context, teacherID, subjectID string) error {
	id := uuid.New().String()
	query := `INSERT INTO teacher_subjects (id, teacher_id, subject_id, created_at)
	          VALUES ($1, $2, $3, NOW())
			  ON CONFLICT (teacher_id, subject_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, id, teacherID, subjectID)
	return err
}

func (r *PostgresRepository) RemoveSubject(ctx context.Context, teacherID, subjectID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM teacher_subjects WHERE teacher_id = $1 AND subject_id = $2", teacherID, subjectID)
	return err
}
