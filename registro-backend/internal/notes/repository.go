package notes

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"registro-backend/pkg/crypto"
)

type Repository interface {
	Create(ctx context.Context, note *StudentNote) error
	Update(ctx context.Context, note *StudentNote) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*StudentNote, error)
	List(ctx context.Context, filter NoteFilter) ([]StudentNote, error)
	ApproveNote(ctx context.Context, id string, approverID string) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, n *StudentNote) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	encryptedNote, err := crypto.EncryptString(n.Note)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO student_notes (id, school_id, student_id, teacher_id, class_id, subject_id, type, note, date, is_approved, approved_by, approved_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	var appBy sql.NullString
	if n.ApprovedBy != "" {
		appBy = sql.NullString{String: n.ApprovedBy, Valid: true}
	}

	_, err = r.db.ExecContext(ctx, query,
		n.ID, n.SchoolID, n.StudentID, n.TeacherID, n.ClassID, n.SubjectID,
		n.Type, encryptedNote, n.Date, n.IsApproved, appBy, n.ApprovedAt, n.CreatedAt, n.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) Update(ctx context.Context, n *StudentNote) error {
	n.UpdatedAt = time.Now()
	encryptedNote, err := crypto.EncryptString(n.Note)
	if err != nil {
		return err
	}
	query := `
		UPDATE student_notes 
		SET type=$1, note=$2, date=$3, subject_id=$4, updated_at=$5
		WHERE id=$6
	`
	_, err = r.db.ExecContext(ctx, query, n.Type, encryptedNote, n.Date, n.SubjectID, n.UpdatedAt, n.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM student_notes WHERE id=$1", id)
	return err
}

func (r *PostgresRepository) ApproveNote(ctx context.Context, id string, approverID string) error {
	query := `
		UPDATE student_notes
		SET is_approved = true, approved_by = $1, approved_at = NOW(), updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, approverID, id)
	return err
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*StudentNote, error) {
	query := `
		SELECT id, school_id, student_id, teacher_id, class_id, subject_id, type, note, to_char(date, 'YYYY-MM-DD'),
		       COALESCE(is_approved, true), COALESCE(approved_by, ''), approved_at, created_at, updated_at
		FROM student_notes WHERE id=$1
	`
	var n StudentNote
	var subjectID sql.NullString
	var appBy sql.NullString
	var appAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&n.ID, &n.SchoolID, &n.StudentID, &n.TeacherID, &n.ClassID, &subjectID,
		&n.Type, &n.Note, &n.Date, &n.IsApproved, &appBy, &appAt, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if subjectID.Valid {
		s := subjectID.String
		n.SubjectID = &s
	}
	if appBy.Valid {
		n.ApprovedBy = appBy.String
	}
	if appAt.Valid {
		n.ApprovedAt = &appAt.Time
	}

	decrypted, err := crypto.DecryptString(n.Note)
	if err == nil {
		n.Note = decrypted
	}

	return &n, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter NoteFilter) ([]StudentNote, error) {
	query := `
		SELECT n.id, n.school_id, n.student_id, n.teacher_id, n.class_id, n.subject_id, n.type, n.note, to_char(n.date, 'YYYY-MM-DD'),
		       COALESCE(n.is_approved, true), COALESCE(n.approved_by, ''), n.approved_at, n.created_at, n.updated_at,
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
		argIdx++
	}

	// Hide unapproved notes from students and parents
	if filter.ActorRole == "student" || filter.ActorRole == "parent" {
		query += " AND COALESCE(n.is_approved, true) = true"
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
		var appBy sql.NullString
		var appAt sql.NullTime

		if err := rows.Scan(
			&n.ID, &n.SchoolID, &n.StudentID, &n.TeacherID, &n.ClassID, &subjectID,
			&n.Type, &n.Note, &n.Date, &n.IsApproved, &appBy, &appAt, &n.CreatedAt, &n.UpdatedAt,
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
		if appBy.Valid {
			n.ApprovedBy = appBy.String
		}
		if appAt.Valid {
			n.ApprovedAt = &appAt.Time
		}

		decrypted, err := crypto.DecryptString(n.Note)
		if err == nil {
			n.Note = decrypted
		}

		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}
