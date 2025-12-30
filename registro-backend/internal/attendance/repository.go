package attendance

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	FindByClassAndDate(classID uint, date string) ([]Attendance, error)
	Create(attendance *Attendance) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByClassAndDate(classID uint, date string) ([]Attendance, error) {
	query := `SELECT id, student_id, class_id, date, status, justified, created_at, updated_at FROM attendance WHERE class_id = $1 AND date::date = $2::date AND deleted_at IS NULL`
	rows, err := r.db.Query(query, classID, date)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var atts []Attendance
	for rows.Next() {
		var a Attendance
		if err := rows.Scan(&a.ID, &a.StudentID, &a.ClassID, &a.Date, &a.Status, &a.Justified, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		atts = append(atts, a)
	}
	return atts, nil
}

func (r *repository) Create(attendance *Attendance) error {
	query := `INSERT INTO attendance (student_id, class_id, date, status, justified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id`

	err := r.db.QueryRow(query, attendance.StudentID, attendance.ClassID, attendance.Date, attendance.Status, attendance.Justified).Scan(&attendance.ID)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	return nil
}
