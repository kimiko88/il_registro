package attendance

import (
	"registro-backend/internal/models"
	"time"
)

type Attendance struct {
	models.Model
	StudentID uint      `json:"student_id"`
	ClassID   uint      `json:"class_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"` // Present, Absent, Late
	Justified bool      `json:"justified"`
}
