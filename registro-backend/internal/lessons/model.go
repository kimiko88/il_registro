package lessons

import "time"

// Lesson represents a class register entry (argomento della lezione)
type Lesson struct {
	ID        string    `json:"id" db:"id"`
	ClassID   string    `json:"class_id" db:"class_id"`
	TeacherID string    `json:"teacher_id" db:"teacher_id"`
	SubjectID string    `json:"subject_id" db:"subject_id"`
	Date      time.Time `json:"date" db:"date"`
	Hour      int       `json:"hour" db:"hour"`         // Quale ora
	Duration  int       `json:"duration" db:"duration"` // Per quante ore
	Topic     string    `json:"topic" db:"topic"`       // Argomento
	Type      string    `json:"type" db:"type"`         // Frontale, Laboratorio, Verifica
	Notes     string    `json:"notes" db:"notes"`       // Note aggiuntive interne
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Homework represents an assignment given to a class
type Homework struct {
	ID          string    `json:"id" db:"id"`
	LessonID    *string   `json:"lesson_id,omitempty" db:"lesson_id"` // Optional link to the lesson when it was assigned
	ClassID     string    `json:"class_id" db:"class_id"`
	SubjectID   string    `json:"subject_id" db:"subject_id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
