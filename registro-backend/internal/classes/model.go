package classes

import "time"

type Class struct {
	ID            string    `json:"id"`
	SchoolID      string    `json:"school_id"`
	Name          string    `json:"name"`          // e.g., "1A", "5B"
	Section       string    `json:"section"`       // e.g., "A", "B"
	Articolazione string    `json:"articolazione"` // e.g., "Informatica" (Optional)
	AcademicYear  string    `json:"academic_year"` // e.g., "2024/2025"
	CoordinatorID string    `json:"coordinator_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateClassRequest struct {
	Name          string `json:"name" binding:"required"`
	SchoolID      string `json:"school_id"`
	Section       string `json:"section"`
	Articolazione string `json:"articolazione"`
	AcademicYear  string `json:"academic_year" binding:"required"`
	CoordinatorID string `json:"coordinator_id"`
}

type ClassSubject struct {
	ID           string  `json:"id" db:"id"`
	ClassID      string  `json:"class_id" db:"class_id"`
	SubjectID    string  `json:"subject_id" db:"subject_id"`
	SubjectName  string  `json:"subject_name,omitempty"`     // Joined
	TeacherID    *string `json:"teacher_id" db:"teacher_id"` // Nullable
	TeacherName  string  `json:"teacher_name,omitempty"`     // Joined
	HoursPerWeek float64 `json:"hours_per_week" db:"hours_per_week"`
}

type AssignSubjectRequest struct {
	SubjectID    string  `json:"subject_id" binding:"required"`
	TeacherID    *string `json:"teacher_id"`
	HoursPerWeek float64 `json:"hours_per_week"`
}

type GuardianInfo struct {
	GuardianID  string `json:"guardian_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
}
