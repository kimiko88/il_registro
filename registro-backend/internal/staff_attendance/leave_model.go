package staff_attendance

import "time"

// LeaveType categorizza il tipo di assenza/permesso
type LeaveType string

const (
	LeaveTypeFerie          LeaveType = "ferie"
	LeaveTypePermesso       LeaveType = "permesso"
	LeaveTypeMalattia       LeaveType = "malattia"
	LeaveTypePermessoStudio LeaveType = "permesso_studio"
	LeaveTypeRecupero       LeaveType = "recupero"
	LeaveTypeAspettativa    LeaveType = "aspettativa"
	LeaveTypeCongedo        LeaveType = "congedo_parentale"
	LeaveTypeLegge104       LeaveType = "legge_104"
)

// LeaveStatus stati del workflow di una richiesta
type LeaveStatus string

const (
	LeaveStatusPending  LeaveStatus = "pending"
	LeaveStatusApproved LeaveStatus = "approved"
	LeaveStatusRejected LeaveStatus = "rejected"
)

// LeaveRequest rappresenta una richiesta di assenza/permesso del personale ATA
type LeaveRequest struct {
	ID          string      `json:"id" db:"id"`
	SchoolID    string      `json:"school_id" db:"school_id"`
	UserID      string      `json:"user_id" db:"user_id"`
	Type        LeaveType   `json:"type" db:"type"`
	StartDate   string      `json:"start_date" db:"start_date"` // YYYY-MM-DD
	EndDate     string      `json:"end_date" db:"end_date"`
	Days        float64     `json:"days" db:"days"`
	Hours       float64     `json:"hours" db:"hours"` // per permessi orari
	Notes       string      `json:"notes,omitempty" db:"notes"`
	Status      LeaveStatus `json:"status" db:"status"`
	ApprovedBy  *string     `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt  *time.Time  `json:"approved_at,omitempty" db:"approved_at"`
	RejectedAt  *time.Time  `json:"rejected_at,omitempty" db:"rejected_at"`
	RejectReason string     `json:"reject_reason,omitempty" db:"reject_reason"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`

	// Join fields
	UserFirstName string `json:"user_first_name,omitempty"`
	UserLastName  string `json:"user_last_name,omitempty"`
	UserRole      string `json:"user_role,omitempty"`
}

// MonthlyTimecard è il cartellino mensile aggregato per utente
type MonthlyTimecard struct {
	UserID        string  `json:"user_id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Role          string  `json:"role"`
	Month         string  `json:"month"` // YYYY-MM
	ContractHours float64 `json:"contract_hours"` // ore contrattuali mensili (da CCNL: 36h/sett)
	WorkedHours   float64 `json:"worked_hours"`   // ore effettive da badge
	OvertimeHours float64 `json:"overtime_hours"`
	AbsenceDays   int     `json:"absence_days"`
	LeaveDays     int     `json:"leave_days"`  // ferie approvate
	SickDays      int     `json:"sick_days"`
	PermitHours   float64 `json:"permit_hours"` // ore di permesso

	// Dettaglio giornaliero
	DailyEntries []DailyTimecardEntry `json:"daily_entries,omitempty"`
}

// DailyTimecardEntry riga giornaliera del cartellino
type DailyTimecardEntry struct {
	Date          string     `json:"date"`
	EntryTime     *time.Time `json:"entry_time,omitempty"`
	ExitTime      *time.Time `json:"exit_time,omitempty"`
	WorkedMinutes int        `json:"worked_minutes"`
	Status        string     `json:"status"` // present, absent, leave, sick, holiday
	Notes         string     `json:"notes,omitempty"`
}

// --- Request DTOs ---

type CreateLeaveRequest struct {
	Type      LeaveType `json:"type" binding:"required"`
	StartDate string    `json:"start_date" binding:"required"`
	EndDate   string    `json:"end_date" binding:"required"`
	Days      float64   `json:"days"`
	Hours     float64   `json:"hours"`
	Notes     string    `json:"notes"`
}

type ApproveLeaveRequest struct {
	Notes string `json:"notes"`
}

type RejectLeaveRequest struct {
	Reason string `json:"reason" binding:"required"`
}
