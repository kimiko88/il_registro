package teachers

import (
	"time"
)

type Teacher struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	SchoolID      string    `json:"school_id" db:"school_id"`
	HiringDate    *string   `json:"hiring_date,omitempty" db:"hiring_date"` // String or Time
	Qualification string    `json:"qualification,omitempty" db:"qualification"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
}

type TeacherSubject struct {
	ID          string    `json:"id" db:"id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	SubjectID   string    `json:"subject_id" db:"subject_id"`
	SubjectName string    `json:"subject_name,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type AssignSubjectRequest struct {
	SubjectID string `json:"subject_id" binding:"required"`
}
