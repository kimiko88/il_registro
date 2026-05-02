package lessons

import "time"

type CreateLessonRequest struct {
	ClassID   string `json:"class_id" binding:"required"`
	SubjectID string `json:"subject_id" binding:"required"`
	Date      string `json:"date" binding:"required"` // YYYY-MM-DD
	Topic     string `json:"topic" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Notes     string `json:"notes"`
}

type UpdateLessonRequest struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Notes string `json:"notes"`
}

type CreateHomeworkRequest struct {
	ClassID     string  `json:"class_id" binding:"required"`
	SubjectID   string  `json:"subject_id" binding:"required"`
	LessonID    *string `json:"lesson_id"`
	DueDate     string  `json:"due_date" binding:"required"` // YYYY-MM-DD
	Description string  `json:"description" binding:"required"`
}

type UpdateHomeworkRequest struct {
	DueDate     string `json:"due_date"`
	Description string `json:"description"`
}

type LessonResponse struct {
	ID        string    `json:"id"`
	ClassID   string    `json:"class_id"`
	TeacherID string    `json:"teacher_id"`
	SubjectID string    `json:"subject_id"`
	Date      time.Time `json:"date"`
	Topic     string    `json:"topic"`
	Type      string    `json:"type"`
	Notes     string    `json:"notes"`
}

type HomeworkResponse struct {
	ID          string    `json:"id"`
	LessonID    *string   `json:"lesson_id"`
	ClassID     string    `json:"class_id"`
	SubjectID   string    `json:"subject_id"`
	TeacherID   string    `json:"teacher_id"`
	DueDate     time.Time `json:"due_date"`
	Description string    `json:"description"`
}
