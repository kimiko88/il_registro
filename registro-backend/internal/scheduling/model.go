package scheduling

import (
	"registro-backend/internal/models"
	"time"
)

type Schedule struct {
	models.Model
	ClassID   uint         `json:"class_id"`
	SubjectID uint         `json:"subject_id"`
	TeacherID uint         `json:"teacher_id"`
	DayOfWeek time.Weekday `json:"day_of_week"`
	StartTime string       `json:"start_time"` // HH:MM
	EndTime   string       `json:"end_time"`   // HH:MM
	Room      string       `json:"room"`
}
