package notes

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, note *StudentNote) error
	Update(ctx context.Context, note *StudentNote) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*StudentNote, error)
	List(ctx context.Context, filter NoteFilter) ([]StudentNote, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, n *StudentNote) error {
	n.ID = uuid.New().String()
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	query := `
		INSERT INTO student_notes (id, school_id, student_id, teacher_id, class_id, subject_id, type, note, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		n.ID, n.SchoolID, n.StudentID, n.TeacherID, n.ClassID, n.SubjectID,
		n.Type, n.Note, n.Date, n.CreatedAt, n.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) Update(ctx context.Context, n *StudentNote) error {
	n.UpdatedAt = time.Now()
	query := `
		UPDATE student_notes 
		SET type=$1, note=$2, date=$3, subject_id=$4, updated_at=$5
		WHERE id=$6
	`
	_, err := r.db.ExecContext(ctx, query, n.Type, n.Note, n.Date, n.SubjectID, n.UpdatedAt, n.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM student_notes WHERE id=$1", id)
	return err
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*StudentNote, error) {
	query := `
		SELECT id, school_id, student_id, teacher_id, class_id, subject_id, type, note, to_char(date, 'YYYY-MM-DD'), created_at, updated_at
		FROM student_notes WHERE id=$1
	`
	var n StudentNote
	var subjectID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&n.ID, &n.SchoolID, &n.StudentID, &n.TeacherID, &n.ClassID, &subjectID,
		&n.Type, &n.Note, &n.Date, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if subjectID.Valid {
		s := subjectID.String
		n.SubjectID = &s
	}
	return &n, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter NoteFilter) ([]StudentNote, error) {
	query := `
		SELECT n.id, n.school_id, n.student_id, n.teacher_id, n.class_id, n.subject_id, n.type, n.note, to_char(n.date, 'YYYY-MM-DD'), n.created_at, n.updated_at,
		       u.first_name || ' ' || u.last_name as teacher_name,
		       s.name as subject_name
		FROM student_notes n
		JOIN users u ON n.teacher_id = u.id
		LEFT JOIN subjects s ON n.subject_id = s.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.StudentID != "" {
		query += fmt.Sprintf(" AND n.student_id = $%d", argIdx)
		args = append(args, filter.StudentID)
		argIdx++
	}
	if filter.ClassID != "" {
		query += fmt.Sprintf(" AND n.class_id = $%d", argIdx)
		args = append(args, filter.ClassID)
		argIdx++
	}
	if filter.TeacherID != "" {
		query += fmt.Sprintf(" AND n.teacher_id = $%d", argIdx)
		args = append(args, filter.TeacherID)
		argIdx++
	}
	if filter.Type != "" {
		query += fmt.Sprintf(" AND n.type = $%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.DateFrom != "" {
		query += fmt.Sprintf(" AND n.date >= $%d", argIdx)
		args = append(args, filter.DateFrom)
	}

	query += " ORDER BY n.date DESC, n.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []StudentNote
	for rows.Next() {
		var n StudentNote
		var subjectID sql.NullString
		var subjectName sql.NullString
		if err := rows.Scan(
			&n.ID, &n.SchoolID, &n.StudentID, &n.TeacherID, &n.ClassID, &subjectID,
			&n.Type, &n.Note, &n.Date, &n.CreatedAt, &n.UpdatedAt,
			&n.TeacherName, &subjectName,
		); err != nil {
			return nil, err
		}
		if subjectID.Valid {
			s := subjectID.String
			n.SubjectID = &s
		}
		if subjectName.Valid {
			n.SubjectName = subjectName.String
		}
		notes = append(notes, n)
	}
	return notes, nil
}
