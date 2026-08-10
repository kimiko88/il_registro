package timetables

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Repository interface {
	GetByClass(ctx context.Context, classID string) ([]ClassSchedule, error)
	GetStudentClassID(ctx context.Context, userID string) (string, error)
	GetParentStudentClassID(ctx context.Context, userID string) (string, error)
	GetTeacherSchedule(ctx context.Context, userID string) ([]ClassSchedule, error)
	GetClassSchoolID(ctx context.Context, classID string) (string, error)
	Update(ctx context.Context, classID string, entries []ScheduleEntry) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetStudentClassID(ctx context.Context, userID string) (string, error) {
	query := `
		SELECT class_id FROM students WHERE user_id = $1::uuid OR id = $1::uuid
	`
	var classID string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&classID)
	return classID, err
}

func (r *PostgresRepository) GetParentStudentClassID(ctx context.Context, userID string) (string, error) {
	query := `
		SELECT s.class_id
		FROM parent_student ps
		JOIN students s ON ps.student_id = s.id
		WHERE ps.parent_id = $1::uuid OR ps.parent_id IN (SELECT id FROM parents WHERE user_id = $1::uuid)
		LIMIT 1
	`
	var classID string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&classID)
	return classID, err
}

func (r *PostgresRepository) GetClassSchoolID(ctx context.Context, classID string) (string, error) {
	query := `SELECT school_id FROM classes WHERE id = $1::uuid`
	var schoolID string
	err := r.db.QueryRowContext(ctx, query, classID).Scan(&schoolID)
	return schoolID, err
}

func (r *PostgresRepository) GetTeacherSchedule(ctx context.Context, userID string) ([]ClassSchedule, error) {
	query := `
		SELECT cs.id, cs.class_id, cs.day_of_week, cs.hour_index, cs.subject_id, s.name, 
		       cs.teacher_id, u.last_name, u.first_name, COALESCE(cs.room, ''), cs.created_at, cs.updated_at
		FROM class_schedules cs
		JOIN subjects s ON cs.subject_id = s.id
		LEFT JOIN users u ON cs.teacher_id = u.id
		LEFT JOIN teachers t ON t.id = cs.teacher_id OR t.user_id = cs.teacher_id
		WHERE cs.teacher_id = $1::uuid OR t.user_id = $1::uuid OR t.id = $1::uuid
		ORDER BY cs.day_of_week, cs.hour_index
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ClassSchedule
	for rows.Next() {
		var cs ClassSchedule
		var tLast, tFirst sql.NullString
		err := rows.Scan(
			&cs.ID, &cs.ClassID, &cs.DayOfWeek, &cs.HourIndex, &cs.SubjectID, &cs.SubjectName,
			&cs.TeacherID, &tLast, &tFirst, &cs.Room, &cs.CreatedAt, &cs.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if tLast.Valid {
			cs.TeacherName = tLast.String + " " + tFirst.String
		}
		results = append(results, cs)
	}
	return results, nil
}

func (r *PostgresRepository) GetByClass(ctx context.Context, classID string) ([]ClassSchedule, error) {
	query := `
		SELECT cs.id, cs.class_id, cs.day_of_week, cs.hour_index, cs.subject_id, s.name, 
		       cs.teacher_id, u.last_name, u.first_name, COALESCE(cs.room, ''), cs.created_at, cs.updated_at
		FROM class_schedules cs
		JOIN subjects s ON cs.subject_id = s.id
		LEFT JOIN users u ON cs.teacher_id = u.id
		WHERE cs.class_id = $1
		ORDER BY cs.day_of_week, cs.hour_index
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ClassSchedule
	for rows.Next() {
		var cs ClassSchedule
		var tLast, tFirst sql.NullString
		err := rows.Scan(
			&cs.ID, &cs.ClassID, &cs.DayOfWeek, &cs.HourIndex, &cs.SubjectID, &cs.SubjectName,
			&cs.TeacherID, &tLast, &tFirst, &cs.Room, &cs.CreatedAt, &cs.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if tLast.Valid {
			cs.TeacherName = tLast.String + " " + tFirst.String
		}
		results = append(results, cs)
	}
	return results, nil
}

func (r *PostgresRepository) normalizeTeacherID(ctx context.Context, teacherID *string) *string {
	if teacherID == nil {
		return nil
	}
	tid := strings.TrimSpace(*teacherID)
	if tid == "" || tid == "null" || tid == "undefined" {
		return nil
	}
	var userExists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1::uuid)", tid).Scan(&userExists)
	if err == nil && userExists {
		return &tid
	}
	var userID string
	err = r.db.QueryRowContext(ctx, "SELECT user_id::text FROM teachers WHERE id = $1::uuid OR user_id = $1::uuid LIMIT 1", tid).Scan(&userID)
	if err == nil && userID != "" {
		return &userID
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, classID string, entries []ScheduleEntry) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 1. Delete existing schedule for the class
	_, err = tx.ExecContext(ctx, "DELETE FROM class_schedules WHERE class_id = $1", classID)
	if err != nil {
		return err
	}

	// 2. Insert new entries
	if len(entries) > 0 {
		var values []string
		var args []interface{}
		argIdx := 1

		for _, e := range entries {
			id := uuid.New().String()
			normTeacherID := r.normalizeTeacherID(ctx, e.TeacherID)
			values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", 
				argIdx, argIdx+1, argIdx+2, argIdx+3, argIdx+4, argIdx+5, argIdx+6))
			args = append(args, id, classID, e.DayOfWeek, e.HourIndex, e.SubjectID, normTeacherID, e.Room)
			argIdx += 7
		}

		query := fmt.Sprintf(`
			INSERT INTO class_schedules (id, class_id, day_of_week, hour_index, subject_id, teacher_id, room)
			VALUES %s
		`, strings.Join(values, ","))

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
