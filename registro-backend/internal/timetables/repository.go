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
	UpdateTeacher(ctx context.Context, teacherID string, entries []TeacherScheduleEntry) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetStudentClassID(ctx context.Context, userID string) (string, error) {
	query := `
		SELECT class_id FROM students WHERE (user_id = $1::uuid OR id = $1::uuid) AND class_id IS NOT NULL
		UNION
		SELECT cs.class_id FROM class_students cs LEFT JOIN students s ON cs.student_id = s.id OR cs.student_id = s.user_id WHERE (cs.student_id = $1::uuid OR s.user_id = $1::uuid OR s.id = $1::uuid) AND cs.class_id IS NOT NULL
		LIMIT 1
	`
	var classID string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&classID)
	return classID, err
}

func (r *PostgresRepository) GetParentStudentClassID(ctx context.Context, userID string) (string, error) {
	query := `
		SELECT s.class_id
		FROM student_parents sp
		LEFT JOIN parents p ON sp.parent_id = p.id
		LEFT JOIN students s ON sp.student_id = s.id
		WHERE (sp.parent_id = $1::uuid OR p.user_id = $1::uuid OR p.id = $1::uuid)
		  AND s.class_id IS NOT NULL
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
		SELECT cs.id, cs.class_id, COALESCE(c.name || c.section, c.name, ''), cs.day_of_week, cs.hour_index, cs.subject_id, s.name, 
		       cs.teacher_id, COALESCE(u.last_name, tu.last_name, ''), COALESCE(u.first_name, tu.first_name, ''), COALESCE(cs.room, ''), cs.created_at, cs.updated_at
		FROM class_schedules cs
		JOIN subjects s ON cs.subject_id = s.id
		JOIN classes c ON cs.class_id = c.id
		LEFT JOIN users u ON cs.teacher_id = u.id
		LEFT JOIN teachers t ON t.id = cs.teacher_id OR t.user_id = cs.teacher_id
		LEFT JOIN users tu ON t.user_id = tu.id
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
			&cs.ID, &cs.ClassID, &cs.ClassName, &cs.DayOfWeek, &cs.HourIndex, &cs.SubjectID, &cs.SubjectName,
			&cs.TeacherID, &tLast, &tFirst, &cs.Room, &cs.CreatedAt, &cs.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if tLast.Valid && tFirst.Valid {
			cs.TeacherName = strings.TrimSpace(tLast.String + " " + tFirst.String)
		}
		results = append(results, cs)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if results == nil {
		results = []ClassSchedule{}
	}
	return results, nil
}

func (r *PostgresRepository) GetByClass(ctx context.Context, classID string) ([]ClassSchedule, error) {
	query := `
		SELECT cs.id, cs.class_id, cs.day_of_week, cs.hour_index, cs.subject_id, s.name, 
		       cs.teacher_id, COALESCE(u.last_name, tu.last_name, ''), COALESCE(u.first_name, tu.first_name, ''), COALESCE(cs.room, ''), cs.created_at, cs.updated_at
		FROM class_schedules cs
		JOIN subjects s ON cs.subject_id = s.id
		LEFT JOIN users u ON cs.teacher_id = u.id
		LEFT JOIN teachers t ON t.id = cs.teacher_id OR t.user_id = cs.teacher_id
		LEFT JOIN users tu ON t.user_id = tu.id
		WHERE cs.class_id = $1::uuid
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
		if tLast.Valid && tFirst.Valid {
			cs.TeacherName = strings.TrimSpace(tLast.String + " " + tFirst.String)
		}
		results = append(results, cs)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if results == nil {
		results = []ClassSchedule{}
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
	if _, err := uuid.Parse(tid); err != nil {
		return nil
	}
	var userExists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1::uuid)", tid).Scan(&userExists)
	if err == nil && userExists {
		return &tid
	}
	var userID string
	err = r.db.QueryRowContext(ctx, "SELECT user_id::text FROM teachers WHERE id = $1::uuid OR user_id = $1::uuid LIMIT 1", tid).Scan(&userID)
	if err == nil && userID != "" && userID != "<nil>" {
		if _, err := uuid.Parse(userID); err == nil {
			return &userID
		}
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

func (r *PostgresRepository) UpdateTeacher(ctx context.Context, teacherID string, entries []TeacherScheduleEntry) error {
	normTeacherIDPtr := r.normalizeTeacherID(ctx, &teacherID)
	if normTeacherIDPtr == nil {
		return fmt.Errorf("teacher_id non valido o docente non trovato: %s", teacherID)
	}
	normTeacherID := *normTeacherIDPtr

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 1. Clear teacher_id from existing slots assigned to this teacher
	_, err = tx.ExecContext(ctx, "UPDATE class_schedules SET teacher_id = NULL WHERE teacher_id = $1::uuid", normTeacherID)
	if err != nil {
		return err
	}

	// 2. Process new entries for this teacher
	for _, e := range entries {
		if strings.TrimSpace(e.ClassID) == "" || strings.TrimSpace(e.SubjectID) == "" {
			continue
		}
		// Check if a slot already exists for this class, day, and hour
		var existingID string
		err := tx.QueryRowContext(ctx,
			"SELECT id FROM class_schedules WHERE class_id = $1::uuid AND day_of_week = $2 AND hour_index = $3",
			e.ClassID, e.DayOfWeek, e.HourIndex,
		).Scan(&existingID)

		if err == nil && existingID != "" {
			// Update existing slot with teacher_id, subject_id, room
			_, err = tx.ExecContext(ctx, `
				UPDATE class_schedules 
				SET teacher_id = $1::uuid, subject_id = $2::uuid, room = $3, updated_at = NOW()
				WHERE id = $4::uuid
			`, normTeacherID, e.SubjectID, e.Room, existingID)
			if err != nil {
				return err
			}
		} else {
			// Insert new slot for this class
			id := uuid.New().String()
			_, err = tx.ExecContext(ctx, `
				INSERT INTO class_schedules (id, class_id, day_of_week, hour_index, subject_id, teacher_id, room)
				VALUES ($1::uuid, $2::uuid, $3, $4, $5::uuid, $6::uuid, $7)
			`, id, e.ClassID, e.DayOfWeek, e.HourIndex, e.SubjectID, normTeacherID, e.Room)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
