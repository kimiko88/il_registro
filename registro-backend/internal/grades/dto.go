package grades

import "time"

type CreateGradeRequest struct {
	StudentID   uint      `json:"student_id" binding:"required"`
	SubjectID   uint      `json:"subject_id" binding:"required"`
	Value       float64   `json:"value" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
}

type GradeResponse struct {
	ID          uint      `json:"id"`
	StudentID   uint      `json:"student_id"`
	SubjectID   uint      `json:"subject_id"`
	TeacherID   uint      `json:"teacher_id"`
	Value       float64   `json:"value"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}
