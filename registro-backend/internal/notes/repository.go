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
	DeleteWithReason(ctx context.Context, id string, reason string) error
	Get(ctx context.Context, id string) (*StudentNote, error)
	List(ctx context.Context, filter NoteFilter) ([]StudentNote, error)
	ApproveNote(ctx context.Context, id string, approverID string) error
	MarkAsViewedByParent(ctx context.Context, id string) error
	MarkManyAsViewedByParent(ctx context.Context, ids []string) error
	IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error)
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
	if n.TargetRole == "" {
		n.TargetRole = "all"
	}

	encryptedNote, err := crypto.EncryptString(n.Note)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO student_notes (id, school_id, student_id, teacher_id, class_id, subject_id, type, note, date, is_reserved, target_role, is_approved, approved_by, approved_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	var appBy sql.NullString
	if n.ApprovedBy != "" {
		appBy = sql.NullString{String: n.ApprovedBy, Valid: true}
	}

	_, err = r.db.ExecContext(ctx, query,
		n.ID, n.SchoolID, n.StudentID, n.TeacherID, n.ClassID, n.SubjectID,
		n.Type, encryptedNote, n.Date, n.IsReserved, n.TargetRole, n.IsApproved, appBy, n.ApprovedAt, n.CreatedAt, n.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) Update(ctx context.Context, n *StudentNote) error {
	n.UpdatedAt = time.Now()
	if n.TargetRole == "" {
		n.TargetRole = "all"
	}
	encryptedNote, err := crypto.EncryptString(n.Note)
	if err != nil {
		return err
	}
	query := `
		UPDATE student_notes 
		SET type=$1, note=$2, date=$3, subject_id=$4, is_reserved=$5, target_role=$6, updated_at=$7
		WHERE id=$8
	`
	_, err = r.db.ExecContext(ctx, query, n.Type, encryptedNote, n.Date, n.SubjectID, n.IsReserved, n.TargetRole, n.UpdatedAt, n.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM student_notes WHERE id=$1", id)
	return err
}

func (r *PostgresRepository) DeleteWithReason(ctx context.Context, id string, reason string) error {
	if reason != "" {
		_, _ = r.db.ExecContext(ctx, "UPDATE student_notes SET deletion_reason = $1 WHERE id = $2", reason, id)
	}
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

func (r *PostgresRepository) MarkAsViewedByParent(ctx context.Context, id string) error {
	query := `
		UPDATE student_notes
		SET is_viewed_by_parent = true, parent_viewed_at = NOW()
		WHERE id = $1 AND COALESCE(is_viewed_by_parent, false) = false
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// MarkManyAsViewedByParent esegue un singolo UPDATE per tutte le note non ancora viste dal genitore.
// Questo evita le N query individuali che venivano eseguite all'interno di una GET in ListNotes.
func (r *PostgresRepository) MarkManyAsViewedByParent(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	query := `
		UPDATE student_notes
		SET is_viewed_by_parent = true, parent_viewed_at = NOW()
		WHERE id = ANY($1) AND COALESCE(is_viewed_by_parent, false) = false
	`
	_, err := r.db.ExecContext(ctx, query, ids)
	return err
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*StudentNote, error) {
	query := `
		SELECT id, school_id, student_id, teacher_id, class_id, subject_id, type, note, to_char(date, 'YYYY-MM-DD'),
		       COALESCE(is_reserved, false), COALESCE(target_role, 'all'),
		       COALESCE(is_approved, true), COALESCE(approved_by::text, ''), approved_at, created_at, updated_at
		FROM student_notes WHERE id=$1
	`
	var n StudentNote
	var subjectID sql.NullString
	var appBy sql.NullString
	var appAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&n.ID, &n.SchoolID, &n.StudentID, &n.TeacherID, &n.ClassID, &subjectID,
		&n.Type, &n.Note, &n.Date, &n.IsReserved, &n.TargetRole,
		&n.IsApproved, &appBy, &appAt, &n.CreatedAt, &n.UpdatedAt,
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
		       COALESCE(n.is_reserved, false), COALESCE(n.target_role, 'all'),
		       COALESCE(n.is_approved, true), COALESCE(n.approved_by::text, ''), n.approved_at, n.created_at, n.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, 'Docente') as teacher_name,
		       s.name as subject_name
		FROM student_notes n
		LEFT JOIN users u ON n.teacher_id = u.id
		LEFT JOIN subjects s ON n.subject_id = s.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.SchoolID != "" {
		if _, err := uuid.Parse(filter.SchoolID); err == nil {
			query += fmt.Sprintf(" AND n.school_id = $%d::uuid", argIdx)
			args = append(args, filter.SchoolID)
			argIdx++
		}
	}
	if filter.StudentID != "" {
		if _, err := uuid.Parse(filter.StudentID); err == nil {
			query += fmt.Sprintf(" AND (n.student_id = $%d::uuid OR n.student_id IN (SELECT user_id FROM students WHERE id = $%d::uuid))", argIdx, argIdx)
			args = append(args, filter.StudentID)
			argIdx++
		}
	}
	if filter.ClassID != "" {
		if _, err := uuid.Parse(filter.ClassID); err == nil {
			query += fmt.Sprintf(" AND n.class_id = $%d::uuid", argIdx)
			args = append(args, filter.ClassID)
			argIdx++
		}
	}
	if filter.TeacherID != "" {
		if _, err := uuid.Parse(filter.TeacherID); err == nil {
			query += fmt.Sprintf(" AND n.teacher_id = $%d::uuid", argIdx)
			args = append(args, filter.TeacherID)
			argIdx++
		}
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
	if filter.DateTo != "" {
		query += fmt.Sprintf(" AND n.date <= $%d", argIdx)
		args = append(args, filter.DateTo)
		argIdx++
	}

	// Reserved notes visibility logic
	if filter.ActorRole == "student" || filter.ActorRole == "parent" {
		query += " AND COALESCE(n.is_approved, true) = true AND COALESCE(n.is_reserved, false) = false"
	} else if filter.ActorRole == "teacher" && !filter.IsCoordinator {
		if _, err := uuid.Parse(filter.ActorID); err == nil {
			query += fmt.Sprintf(" AND (COALESCE(n.is_reserved, false) = false OR n.teacher_id = $%d::uuid)", argIdx)
			args = append(args, filter.ActorID)
			argIdx++
		} else {
			query += " AND COALESCE(n.is_reserved, false) = false"
		}
	}

	query += " ORDER BY n.date DESC, n.created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}
	if filter.Page > 1 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, offset)
	}

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
			&n.Type, &n.Note, &n.Date, &n.IsReserved, &n.TargetRole,
			&n.IsApproved, &appBy, &appAt, &n.CreatedAt, &n.UpdatedAt,
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

func (r *PostgresRepository) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	if _, err := uuid.Parse(teacherID); err != nil {
		return false, nil
	}
	if _, err := uuid.Parse(classID); err != nil {
		return false, nil
	}
	query := `
		SELECT EXISTS (
			SELECT 1 FROM class_subjects cs
			JOIN teachers t ON cs.teacher_id = t.id
			WHERE (t.user_id = $1::uuid OR t.id = $1::uuid) AND cs.class_id = $2::uuid
		) OR EXISTS (
			SELECT 1 FROM classes WHERE (coordinator_id = $1::uuid OR coordinator_id IN (SELECT user_id FROM teachers WHERE id = $1::uuid)) AND id = $2::uuid
		) OR EXISTS (
			SELECT 1 FROM users u
			JOIN classes c ON c.id = $2::uuid AND c.school_id = u.school_id
			WHERE (u.id = $1::uuid OR u.id IN (SELECT user_id FROM teachers WHERE id = $1::uuid)) AND COALESCE(u.is_staff, false) = true
		)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, teacherID, classID).Scan(&exists)
	return exists, err
}
