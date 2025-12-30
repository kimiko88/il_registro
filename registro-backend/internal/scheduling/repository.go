package scheduling

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	FindByClass(classID uint) ([]Schedule, error)
	Create(schedule *Schedule) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByClass(classID uint) ([]Schedule, error) {
	query := `SELECT id, class_id, subject_id, teacher_id, day_of_week, start_time, end_time, room, created_at, updated_at FROM schedules WHERE class_id = $1 AND deleted_at IS NULL`
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var scheds []Schedule
	for rows.Next() {
		var s Schedule
		if err := rows.Scan(&s.ID, &s.ClassID, &s.SubjectID, &s.TeacherID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.Room, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		scheds = append(scheds, s)
	}
	return scheds, nil
}

func (r *repository) Create(schedule *Schedule) error {
	query := `INSERT INTO schedules (class_id, subject_id, teacher_id, day_of_week, start_time, end_time, room, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id`

	err := r.db.QueryRow(query, schedule.ClassID, schedule.SubjectID, schedule.TeacherID, schedule.DayOfWeek, schedule.StartTime, schedule.EndTime, schedule.Room).Scan(&schedule.ID)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	return nil
}
