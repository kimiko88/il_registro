package teacher_hours

import "time"

// TeacherHourObligation tracks Art. 29 CCNL obligations ("40+40 hours")
type TeacherHourObligation struct {
	ID                   string    `json:"id" db:"id"`
	SchoolID             string    `json:"school_id" db:"school_id"`
	TeacherID            string    `json:"teacher_id" db:"teacher_id"`
	AcademicYear         string    `json:"academic_year" db:"academic_year"`
	CollegioPlannedHours float64   `json:"collegio_planned_hours" db:"collegio_planned_hours"` // 40.00
	ConsigliPlannedHours float64   `json:"consigli_planned_hours" db:"consigli_planned_hours"` // 40.00
	EccedentiHours       float64   `json:"eccedenti_hours" db:"eccedenti_hours"`
	DisposizioneHours    float64   `json:"disposizione_hours" db:"disposizione_hours"`
	BancaOreBalance      float64   `json:"banca_ore_balance" db:"banca_ore_balance"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`

	// Computed summary fields
	CollegioDoneHours      float64 `json:"collegio_done_hours"`
	CollegioRemainingHours float64 `json:"collegio_remaining_hours"`
	ConsigliDoneHours      float64 `json:"consigli_done_hours"`
	ConsigliRemainingHours float64 `json:"consigli_remaining_hours"`
	HasOvertime            bool    `json:"has_overtime"`
	HasDeficit             bool    `json:"has_deficit"`
}

// TeacherHourLog records an individual hour deduction / addition
type TeacherHourLog struct {
	ID           string    `json:"id" db:"id"`
	ObligationID string    `json:"obligation_id" db:"obligation_id"`
	TeacherID    string    `json:"teacher_id" db:"teacher_id"`
	Category     string    `json:"category" db:"category"` // collegio, consigli, eccedenti, disposizione, permesso_breve, recupero_permesso
	Hours        float64   `json:"hours" db:"hours"`
	ActivityDate string    `json:"activity_date" db:"activity_date"` // YYYY-MM-DD
	MeetingID    *string   `json:"meeting_id,omitempty" db:"meeting_id"`
	Description  string    `json:"description" db:"description"`
	Status       string    `json:"status" db:"status"` // completed, debited, recuperato
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// TeacherShortLeave records short leaves up to 2 hours (permessi brevi) with 2-month recovery limit
type TeacherShortLeave struct {
	ID               string    `json:"id" db:"id"`
	SchoolID         string    `json:"school_id" db:"school_id"`
	TeacherID        string    `json:"teacher_id" db:"teacher_id"`
	RequestDate      string    `json:"request_date" db:"request_date"` // YYYY-MM-DD
	Hours            float64   `json:"hours" db:"hours"`               // <= 2.0
	Reason           string    `json:"reason" db:"reason"`
	Status           string    `json:"status" db:"status"` // requested, approved, rejected, recovered
	RecoveryDeadline string    `json:"recovery_deadline" db:"recovery_deadline"`
	RecoveredHours   float64   `json:"recovered_hours" db:"recovered_hours"`
	RecoveredDate    *string   `json:"recovered_date,omitempty" db:"recovered_date"`
	Notes            string    `json:"notes" db:"notes"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`

	IsExpired bool `json:"is_expired"`
}
