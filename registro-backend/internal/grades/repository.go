package grades

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	FindByStudent(studentID uint) ([]Grade, error)
	Create(grade *Grade) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByStudent(studentID uint) ([]Grade, error) {
	query := `SELECT id, student_id, subject_id, teacher_id, value, type, description, date, created_at, updated_at FROM grades WHERE student_id = $1 AND deleted_at IS NULL`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var grades []Grade
	for rows.Next() {
		var g Grade
		if err := rows.Scan(&g.ID, &g.StudentID, &g.SubjectID, &g.TeacherID, &g.Value, &g.Type, &g.Description, &g.Date, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		grades = append(grades, g)
	}
	return grades, nil
}

func (r *repository) Create(grade *Grade) error {
	query := `INSERT INTO grades (student_id, subject_id, teacher_id, value, type, description, date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id`

	err := r.db.QueryRow(query, grade.StudentID, grade.SubjectID, grade.TeacherID, grade.Value, grade.Type, grade.Description, grade.Date).Scan(&grade.ID)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	return nil
}
