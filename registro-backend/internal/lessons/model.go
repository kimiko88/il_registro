package lessons

import "time"

// Lesson represents a class register entry (argomento della lezione)
type Lesson struct {
	ID          string    `json:"id" db:"id"`
	ClassID     string    `json:"class_id" db:"class_id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	TeacherName string    `json:"teacher_name" db:"teacher_name"`
	SubjectID   string    `json:"subject_id" db:"subject_id"`
	Date        time.Time `json:"date" db:"date"`
	Hour        int       `json:"hour" db:"hour"`         // Quale ora
	Duration    int       `json:"duration" db:"duration"` // Per quante ore
	Topic       string    `json:"topic" db:"topic"`       // Argomento
	Type        string    `json:"type" db:"type"`         // Frontale, Laboratorio, Verifica
	GroupID                *string   `json:"group_id,omitempty" db:"group_id"`
	IsSubstitution         bool      `json:"is_substitution" db:"is_substitution"`
	SubstitutedTeacherID   *string   `json:"substituted_teacher_id,omitempty" db:"substituted_teacher_id"`
	SubstitutedTeacherName string    `json:"substituted_teacher_name,omitempty" db:"substituted_teacher_name"`
	ActivityType           string    `json:"activity_type" db:"activity_type"` // e.g. standard, substitution, ptof, project, assembly, trip, lab, other
	Notes                  string    `json:"notes" db:"notes"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// Homework represents an assignment given to a class
type Homework struct {
	ID          string    `json:"id" db:"id"`
	LessonID    *string   `json:"lesson_id,omitempty" db:"lesson_id"` // Optional link to the lesson when it was assigned
	ClassID     string    `json:"class_id" db:"class_id"`
	SubjectID   string    `json:"subject_id" db:"subject_id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	TeacherName string    `json:"teacher_name" db:"teacher_name"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
