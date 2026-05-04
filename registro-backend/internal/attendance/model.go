package attendance

import (
	"time"

	"github.com/google/uuid"
)

// Enums
type AttendanceStatus string

const (
	StatusPresent   AttendanceStatus = "Present"
	StatusAbsent    AttendanceStatus = "Absent"
	StatusLate      AttendanceStatus = "Late"
	StatusEarlyExit AttendanceStatus = "LeftEarly"
	StatusExempt    AttendanceStatus = "Exempt"
)

type JustificationStatus string

const (
	JustificationPending  JustificationStatus = "pending"
	JustificationApproved JustificationStatus = "approved"
	JustificationRejected JustificationStatus = "rejected"
)

type Attendance struct {
	ID        string `json:"id" db:"id"`
	SchoolID  string `json:"school_id" db:"school_id"`
	StudentID string `json:"student_id" db:"student_id"`
	ClassID   string `json:"class_id" db:"class_id"`

	Date      time.Time        `json:"date" db:"date"`
	Hour      *int             `json:"hour,omitempty" db:"hour"`
	SubjectID *string          `json:"subject_id,omitempty" db:"subject_id"`
	Status    AttendanceStatus `json:"status" db:"status"`

	Justified   bool       `json:"justified" db:"justified"`
	JustifiedBy *string    `json:"justified_by,omitempty" db:"justified_by"`
	JustifiedAt *time.Time `json:"justified_at,omitempty" db:"justified_at"`
	Notes       string     `json:"notes,omitempty" db:"notes"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// Justification represents a request to justify an absence
type Justification struct {
	ID        string `json:"id" db:"id"`
	StudentID string `json:"student_id" db:"student_id"`
	ParentID  string `json:"parent_id,omitempty" db:"parent_id"`

	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`

	Reason string              `json:"reason" db:"reason"`
	Status JustificationStatus `json:"status" db:"status"`

	ApprovedBy *string    `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at,omitempty" db:"approved_at"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (a *Attendance) IsValid() bool {
	if a.StudentID == "" || a.ClassID == "" {
		return false
	}
	if a.Date.IsZero() {
		return false
	}
	// UUID checks
	if _, err := uuid.Parse(a.StudentID); err != nil {
		return false
	}
	return true
}
