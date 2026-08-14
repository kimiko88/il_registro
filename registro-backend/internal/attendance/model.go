package attendance

import (
	"time"

	"github.com/google/uuid"
)

// Enums
type AttendanceStatus string

const (
	StatusPresent    AttendanceStatus = "Present"
	StatusAbsent     AttendanceStatus = "Absent"
	StatusLate       AttendanceStatus = "Late"
	StatusEarlyExit  AttendanceStatus = "LeftEarly"
	StatusOutOfClass AttendanceStatus = "OutOfClass"
	StatusExempt     AttendanceStatus = "Exempt"
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

	Justified           bool       `json:"justified" db:"justified"`
	JustifiedBy         *string    `json:"justified_by,omitempty" db:"justified_by"`
	JustifiedAt         *time.Time `json:"justified_at,omitempty" db:"justified_at"`
	ParentJustified     bool       `json:"parent_justified" db:"parent_justified"`
	ParentJustifiedAt   *time.Time `json:"parent_justified_at,omitempty" db:"parent_justified_at"`
	JustificationReason string     `json:"justification_reason,omitempty" db:"justification_reason"`
	Notes               string     `json:"notes,omitempty" db:"notes"`
	EntryTime           *string    `json:"entry_time,omitempty" db:"entry_time"`
	ExitTime            *string    `json:"exit_time,omitempty" db:"exit_time"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AttendanceStats struct {
	TotalSchoolDays   int                 `json:"total_school_days"`
	DaysPresent       int                 `json:"days_present"`
	DaysAbsent        int                 `json:"days_absent"`
	LateArrivals      int                 `json:"late_arrivals"`
	EarlyExits        int                 `json:"early_exits"`
	Justified         int                 `json:"justified"`
	Unjustified       int                 `json:"unjustified"`
	MonthlyBreakdown  []MonthlyAttendance `json:"monthly_breakdown"`
	AbsencePercentage float64             `json:"absence_percentage"`
}

type MonthlyAttendance struct {
	Month   string `json:"month"`
	Present int    `json:"present"`
	Absent  int    `json:"absent"`
}

type JustifyAbsenceRequest struct {
	Reason string `json:"reason"`
	Notes  string `json:"notes"`
}

// Justification represents a request to justify an absence
type Justification struct {
	ID          string `json:"id" db:"id"`
	SchoolID    string `json:"school_id,omitempty" db:"school_id"`
	StudentID   string `json:"student_id" db:"student_id"`
	StudentName string `json:"student_name,omitempty"`
	ParentID    string `json:"parent_id,omitempty" db:"parent_id"`

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

type MonthlyTrend struct {
	Month        string  `json:"month"`
	Absences     int     `json:"absences"`
	Lates        int     `json:"lates"`
	EarlyExits   int     `json:"early_exits"`
	PresenceRate float64 `json:"presence_rate"`
}

type TrendsResponse struct {
	StudentID string         `json:"student_id"`
	Trends    []MonthlyTrend `json:"trends"`
}

// MonthlyBreakdownRow represents attendance statistics for a single month.
type MonthlyBreakdownRow struct {
	Month             string  `json:"month"`              // YYYY-MM
	Absences          int     `json:"absences"`           // StatusAbsent count
	Lates             int     `json:"lates"`              // StatusLate count
	EarlyExits        int     `json:"early_exits"`        // StatusEarlyExit count
	JustifiedAbsences int     `json:"justified_absences"` // justified Absent
	TotalSchoolDays   int     `json:"total_school_days"`  // total records in month
	PresenceRate      float64 `json:"presence_rate"`      // % presenti
}

// MonthlyBreakdownResponse is the API response for the monthly breakdown endpoint.
type MonthlyBreakdownResponse struct {
	StudentID  string                `json:"student_id"`
	SchoolYear string                `json:"school_year"`
	Months     []MonthlyBreakdownRow `json:"months"`
}
