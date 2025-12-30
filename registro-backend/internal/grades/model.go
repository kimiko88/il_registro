package grades

import (
	"registro-backend/internal/models"
	"time"
)

type Grade struct {
	models.Model
	StudentID   uint      `json:"student_id"`
	SubjectID   uint      `json:"subject_id"`
	TeacherID   uint      `json:"teacher_id"`
	Value       float64   `json:"value"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}
