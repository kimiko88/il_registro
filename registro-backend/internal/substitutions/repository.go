package substitutions

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, sub *Substitution) error
	GetByID(ctx context.Context, id string) (*Substitution, error)
	ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error)
	ListByTeacher(ctx context.Context, teacherID string, date string) ([]*Substitution, error)
	AssignSubstitute(ctx context.Context, id string, substituteTeacherID string, notes string) error
	ConfirmSubstitution(ctx context.Context, id string, substituteTeacherID string) error
	SignRegister(ctx context.Context, id string, sigHash string, notes string) error
	GetAvailableTeachers(ctx context.Context, schoolID string) ([]TeacherCandidate, error)
	GetTeacherProfileID(ctx context.Context, userID string) (string, error)
	IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error)
	IsTeacherAssignedToSubject(ctx context.Context, teacherID, subjectID string) (bool, error)
	GetWeeklySubstitutionCount(ctx context.Context, teacherID string) (int, error)
}

type TeacherCandidate struct {
	TeacherID   string `json:"teacher_id"`
	UserID      string `json:"user_id"`
	TeacherName string `json:"teacher_name"`
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	// Auto-migrate missing columns for substitutions table
	_, _ = db.Exec(`
		ALTER TABLE substitutions
		ADD COLUMN IF NOT EXISTS signed_by_substitute BOOLEAN NOT NULL DEFAULT FALSE,
		ADD COLUMN IF NOT EXISTS signature_timestamp TIMESTAMP WITH TIME ZONE,
		ADD COLUMN IF NOT EXISTS signature_hash VARCHAR(255),
		ADD COLUMN IF NOT EXISTS official_register_notes TEXT;
	`)
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, s *Substitution) error {
	s.ID = uuid.New().String()
	s.CreatedAt = time.Now()

	query := `
		INSERT INTO substitutions (id, school_id, class_id, absent_teacher_id, substitute_teacher_id, date, hour, subject_id, notes, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.SchoolID, s.ClassID, s.AbsentTeacherID, s.SubstituteTeacherID, s.Date,
		s.Hour, s.SubjectID, s.Notes, s.Status, s.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Substitution, error) {
	query := `
		SELECT id, school_id, class_id, absent_teacher_id, substitute_teacher_id, date, hour, COALESCE(subject_id::text, ''), COALESCE(notes, ''), status,
		       COALESCE(signed_by_substitute, FALSE), signature_timestamp, COALESCE(signature_hash, ''), COALESCE(official_register_notes, ''), created_at
		FROM substitutions WHERE id = $1::uuid
	`
	s := &Substitution{}
	var subTeacher sql.NullString
	var sigTs sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.SchoolID, &s.ClassID, &s.AbsentTeacherID, &subTeacher, &s.Date,
		&s.Hour, &s.SubjectID, &s.Notes, &s.Status,
		&s.SignedBySubstitute, &sigTs, &s.SignatureHash, &s.OfficialRegisterNotes, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if subTeacher.Valid {
		s.SubstituteTeacherID = &subTeacher.String
	}
	if sigTs.Valid {
		s.SignatureTimestamp = &sigTs.Time
	}
	return s, nil
}

func (r *PostgresRepository) ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error) {
	var validDate string
	if len(date) == 10 && date[4] == '-' && date[7] == '-' {
		validDate = date
	}

	query := `
		SELECT id, school_id, class_id, absent_teacher_id, substitute_teacher_id, date, hour, COALESCE(subject_id::text, ''), COALESCE(notes, ''), status,
		       COALESCE(signed_by_substitute, FALSE), signature_timestamp, COALESCE(signature_hash, ''), COALESCE(official_register_notes, ''), created_at
		FROM substitutions
		WHERE ($1 = '' OR school_id = $1::uuid)
		  AND ($2 = '' OR date = $2::date)
		ORDER BY date DESC, hour ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, validDate)
	if err != nil {
		return []*Substitution{}, nil
	}
	defer rows.Close()

	var list []*Substitution
	for rows.Next() {
		s := &Substitution{}
		var subTeacher sql.NullString
		var sigTs sql.NullTime
		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.ClassID, &s.AbsentTeacherID, &subTeacher, &s.Date,
			&s.Hour, &s.SubjectID, &s.Notes, &s.Status,
			&s.SignedBySubstitute, &sigTs, &s.SignatureHash, &s.OfficialRegisterNotes, &s.CreatedAt,
		); err != nil {
			return []*Substitution{}, nil
		}
		if subTeacher.Valid {
			s.SubstituteTeacherID = &subTeacher.String
		}
		if sigTs.Valid {
			s.SignatureTimestamp = &sigTs.Time
		}
		list = append(list, s)
	}
	if list == nil {
		list = []*Substitution{}
	}
	return list, nil
}

func (r *PostgresRepository) ListByTeacher(ctx context.Context, teacherID string, date string) ([]*Substitution, error) {
	if teacherID == "" {
		return []*Substitution{}, nil
	}

	var validDate string
	if len(date) == 10 && date[4] == '-' && date[7] == '-' {
		validDate = date
	}

	query := `
		SELECT id, school_id, class_id, absent_teacher_id, substitute_teacher_id, date, hour, COALESCE(subject_id::text, ''), COALESCE(notes, ''), status,
		       COALESCE(signed_by_substitute, FALSE), signature_timestamp, COALESCE(signature_hash, ''), COALESCE(official_register_notes, ''), created_at
		FROM substitutions
		WHERE (NULLIF($1, '') IS NOT NULL AND (
			absent_teacher_id = NULLIF($1, '')::uuid 
			OR substitute_teacher_id = NULLIF($1, '')::uuid
			OR absent_teacher_id IN (SELECT id FROM teachers WHERE user_id = NULLIF($1, '')::uuid)
			OR substitute_teacher_id IN (SELECT id FROM teachers WHERE user_id = NULLIF($1, '')::uuid)
		))
		  AND ($2 = '' OR date = $2::date)
		ORDER BY date DESC, hour ASC
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID, validDate)
	if err != nil {
		return []*Substitution{}, nil
	}
	defer rows.Close()

	var list []*Substitution
	for rows.Next() {
		s := &Substitution{}
		var subTeacher sql.NullString
		var sigTs sql.NullTime
		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.ClassID, &s.AbsentTeacherID, &subTeacher, &s.Date,
			&s.Hour, &s.SubjectID, &s.Notes, &s.Status,
			&s.SignedBySubstitute, &sigTs, &s.SignatureHash, &s.OfficialRegisterNotes, &s.CreatedAt,
		); err != nil {
			return []*Substitution{}, nil
		}
		if subTeacher.Valid {
			s.SubstituteTeacherID = &subTeacher.String
		}
		if sigTs.Valid {
			s.SignatureTimestamp = &sigTs.Time
		}
		list = append(list, s)
	}
	if list == nil {
		list = []*Substitution{}
	}
	return list, nil
}

func (r *PostgresRepository) AssignSubstitute(ctx context.Context, id string, substituteTeacherID string, notes string) error {
	query := `
		UPDATE substitutions
		SET substitute_teacher_id = $2::uuid, status = 'assigned', notes = COALESCE(NULLIF($3, ''), notes)
		WHERE id = $1::uuid
	`
	_, err := r.db.ExecContext(ctx, query, id, substituteTeacherID, notes)
	return err
}

func (r *PostgresRepository) ConfirmSubstitution(ctx context.Context, id string, substituteTeacherID string) error {
	query := `
		UPDATE substitutions
		SET status = 'confirmed'
		WHERE id = $1::uuid AND (substitute_teacher_id = NULLIF($2, '')::uuid OR NULLIF($2, '') IS NULL)
	`
	_, err := r.db.ExecContext(ctx, query, id, substituteTeacherID)
	return err
}

func (r *PostgresRepository) SignRegister(ctx context.Context, id string, sigHash string, notes string) error {
	query := `
		UPDATE substitutions
		SET signed_by_substitute = TRUE,
		    signature_timestamp = NOW(),
		    signature_hash = $2,
		    official_register_notes = COALESCE(NULLIF($3, ''), official_register_notes),
		    status = 'confirmed'
		WHERE id = $1::uuid
	`
	_, err := r.db.ExecContext(ctx, query, id, sigHash, notes)
	return err
}

func (r *PostgresRepository) GetAvailableTeachers(ctx context.Context, schoolID string) ([]TeacherCandidate, error) {
	query := `
		SELECT t.id, t.user_id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		WHERE NULLIF($1, '') IS NULL OR t.school_id = NULLIF($1, '')::uuid`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teachers: %w", err)
	}
	defer rows.Close()

	var list []TeacherCandidate
	for rows.Next() {
		var tc TeacherCandidate
		if err := rows.Scan(&tc.TeacherID, &tc.UserID, &tc.TeacherName); err == nil {
			list = append(list, tc)
		}
	}
	return list, rows.Err()
}

func (r *PostgresRepository) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return userID, nil
	}
	query := `SELECT id FROM teachers WHERE user_id = $1::uuid OR id = $1::uuid LIMIT 1`
	var id string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return userID, nil
	}
	return id, err
}

func (r *PostgresRepository) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	if teacherID == "" || classID == "" {
		return false, nil
	}
	query := `
		SELECT EXISTS (
			SELECT 1 FROM class_subjects cs
			JOIN teachers t ON cs.teacher_id = t.id
			WHERE (t.id = $1::uuid OR t.user_id = $1::uuid) AND cs.class_id = $2::uuid
		) OR EXISTS (
			SELECT 1 FROM classes
			WHERE (coordinator_id = $1::uuid OR coordinator_id IN (SELECT user_id FROM teachers WHERE id = $1::uuid))
			  AND id = $2::uuid
		)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, teacherID, classID).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) IsTeacherAssignedToSubject(ctx context.Context, teacherID, subjectID string) (bool, error) {
	if teacherID == "" || subjectID == "" {
		return false, nil
	}
	query := `
		SELECT EXISTS (
			SELECT 1 FROM class_subjects cs
			JOIN teachers t ON cs.teacher_id = t.id
			WHERE (t.id = $1::uuid OR t.user_id = $1::uuid) AND cs.subject_id = $2::uuid
		)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, teacherID, subjectID).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) GetWeeklySubstitutionCount(ctx context.Context, teacherID string) (int, error) {
	if teacherID == "" {
		return 0, nil
	}
	query := `
		SELECT COUNT(*)
		FROM substitutions
		WHERE (substitute_teacher_id = $1::uuid OR substitute_teacher_id IN (SELECT id FROM teachers WHERE user_id = $1::uuid))
		  AND date >= CURRENT_DATE - INTERVAL '7 days'`
	var count int
	err := r.db.QueryRowContext(ctx, query, teacherID).Scan(&count)
	return count, err
}
