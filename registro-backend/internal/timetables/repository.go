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
	Update(ctx context.Context, classID string, entries []ScheduleEntry) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
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

func (r *PostgresRepository) Update(ctx context.Context, classID string, entries []ScheduleEntry) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
			values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", 
				argIdx, argIdx+1, argIdx+2, argIdx+3, argIdx+4, argIdx+5, argIdx+6))
			args = append(args, id, classID, e.DayOfWeek, e.HourIndex, e.SubjectID, e.TeacherID, e.Room)
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
