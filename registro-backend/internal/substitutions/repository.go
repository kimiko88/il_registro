package substitutions

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, sub *Substitution) error
	GetByID(ctx context.Context, id string) (*Substitution, error)
	ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error)
	ListByTeacher(ctx context.Context, teacherID string) ([]*Substitution, error)
	AssignSubstitute(ctx context.Context, id string, substituteTeacherID string, notes string) error
	ConfirmSubstitution(ctx context.Context, id string, substituteTeacherID string) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
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
		SELECT id, school_id, class_id, absent_teacher_id, substitute_teacher_id, date, hour, subject_id, notes, status, created_at
		FROM substitutions WHERE id = $1::uuid
	`
	s := &Substitution{}
	var subTeacher, subj sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.SchoolID, &s.ClassID, &s.AbsentTeacherID, &subTeacher, &s.Date,
		&s.Hour, &subj, &s.Notes, &s.Status, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if subTeacher.Valid {
		s.SubstituteTeacherID = &subTeacher.String
	}
	if subj.Valid {
		s.SubjectID = subj.String
	}
	return s, nil
}

func (r *PostgresRepository) ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error) {
	query := `
		SELECT id, school_id, class_id, absent_teacher_id, substitute_teacher_id, date, hour, subject_id, notes, status, created_at
		FROM substitutions
		WHERE ($1 = '' OR school_id = $1::uuid)
		  AND ($2 = '' OR date = $2::date)
		ORDER BY date DESC, hour ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Substitution
	for rows.Next() {
		s := &Substitution{}
		var subTeacher, subj sql.NullString
		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.ClassID, &s.AbsentTeacherID, &subTeacher, &s.Date,
			&s.Hour, &subj, &s.Notes, &s.Status, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		if subTeacher.Valid {
			s.SubstituteTeacherID = &subTeacher.String
		}
		if subj.Valid {
			s.SubjectID = subj.String
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *PostgresRepository) ListByTeacher(ctx context.Context, teacherID string) ([]*Substitution, error) {
	if teacherID == "" {
		return []*Substitution{}, nil
	}
	query := `
		SELECT id, school_id, class_id, absent_teacher_id, substitute_teacher_id, date, hour, subject_id, notes, status, created_at
		FROM substitutions
		WHERE (NULLIF($1, '') IS NOT NULL AND (
			absent_teacher_id = NULLIF($1, '')::uuid 
			OR substitute_teacher_id = NULLIF($1, '')::uuid
			OR absent_teacher_id IN (SELECT id FROM teachers WHERE user_id = NULLIF($1, '')::uuid)
			OR substitute_teacher_id IN (SELECT id FROM teachers WHERE user_id = NULLIF($1, '')::uuid)
		))
		ORDER BY date DESC, hour ASC
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Substitution
	for rows.Next() {
		s := &Substitution{}
		var subTeacher, subj sql.NullString
		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.ClassID, &s.AbsentTeacherID, &subTeacher, &s.Date,
			&s.Hour, &subj, &s.Notes, &s.Status, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		if subTeacher.Valid {
			s.SubstituteTeacherID = &subTeacher.String
		}
		if subj.Valid {
			s.SubjectID = subj.String
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
		WHERE id = $1::uuid AND (substitute_teacher_id = $2::uuid OR $2 = '')
	`
	_, err := r.db.ExecContext(ctx, query, id, substituteTeacherID)
	return err
}
