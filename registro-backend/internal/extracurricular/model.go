package extracurricular

import (
	"time"
)

type Course struct {
	ID              string    `json:"id" db:"id"`
	SchoolID        string    `json:"school_id" db:"school_id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description,omitempty" db:"description"`
	TeacherID       string    `json:"teacher_id" db:"teacher_id"`
	StartDate       time.Time `json:"start_date" db:"start_date"`
	EndDate         time.Time `json:"end_date" db:"end_date"`
	MaxParticipants int       `json:"max_participants" db:"max_participants"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`

	// Joined fields
	TeacherName     string `json:"teacher_name,omitempty"`
	EnrolledCount   int    `json:"enrolled_count"`
	IsStudentEnrolled bool `json:"is_student_enrolled"`
}

type Enrollment struct {
	ID         string    `json:"id" db:"id"`
	CourseID   string    `json:"course_id" db:"course_id"`
	StudentID  string    `json:"student_id" db:"student_id"`
	EnrolledAt time.Time `json:"enrolled_at" db:"enrolled_at"`

	// Joined
	StudentName string `json:"student_name,omitempty"`
}

type AttendanceRecord struct {
	ID        string    `json:"id" db:"id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	StudentID string    `json:"student_id" db:"student_id"`
	Date      time.Time `json:"date" db:"date"`
	Status    string    `json:"status" db:"status"` // 'present', 'absent', 'excused'
	Hours     float64   `json:"hours" db:"hours"`
	SignedBy  *string   `json:"signed_by,omitempty" db:"signed_by"`
	SignedAt  time.Time `json:"signed_at" db:"signed_at"`

	// Joined
	StudentName string `json:"student_name,omitempty"`
}

type CreateCourseRequest struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	StartDate       string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate         string `json:"end_date" binding:"required"`   // YYYY-MM-DD
	MaxParticipants int    `json:"max_participants"`
}

type MarkAttendanceRequest struct {
	CourseID  string  `json:"course_id" binding:"required"`
	StudentID string  `json:"student_id" binding:"required"`
	Date      string  `json:"date" binding:"required"` // YYYY-MM-DD
	Status    string  `json:"status" binding:"required"` // 'present', 'absent', 'excused'
	Hours     float64 `json:"hours"`
}
