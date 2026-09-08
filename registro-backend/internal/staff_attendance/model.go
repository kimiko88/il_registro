package staff_attendance

import (
	"time"
)

// ATARoles contiene i ruoli del personale ATA
var ATARoles = []string{
	"dsga",
	"assistente_amministrativo",
	"collaboratore_ds",
	"collaboratore_scolastico",
}

// StaffRolesWithAttendanceWrite contiene i ruoli autorizzati a registrare presenze del personale
// (incluse le presenze docenti durante gli scioperi)
var StaffRolesWithAttendanceWrite = []string{
	"dsga",
	"collaboratore_ds",
	"assistente_amministrativo",
	"collaboratore_scolastico",
	"principal",
	"vice_principal",
	"secretary",
	"admin",
	"superadmin",
}

// StaffRolesWithAttendanceRead contiene i ruoli che possono visualizzare la dashboard presenze
var StaffRolesWithAttendanceRead = []string{
	"dsga",
	"collaboratore_ds",
	"assistente_amministrativo",
	"collaboratore_scolastico",
	"principal",
	"vice_principal",
	"secretary",
	"admin",
	"superadmin",
}

// IsATARole verifica se un ruolo è un ruolo ATA
func IsATARole(role string) bool {
	for _, r := range ATARoles {
		if r == role {
			return true
		}
	}
	return false
}

// CanWriteAttendance verifica se un ruolo può registrare presenze staff
func CanWriteAttendance(role string) bool {
	for _, r := range StaffRolesWithAttendanceWrite {
		if r == role {
			return true
		}
	}
	return false
}

// CanReadAttendance verifica se un ruolo può leggere le presenze staff
func CanReadAttendance(role string) bool {
	for _, r := range StaffRolesWithAttendanceRead {
		if r == role {
			return true
		}
	}
	return false
}

// AttendanceStatus rappresenta lo stato di presenza del personale
type AttendanceStatus string

const (
	StatusPresent   AttendanceStatus = "present"
	StatusAbsent    AttendanceStatus = "absent"
	StatusLate      AttendanceStatus = "late"
	StatusMission   AttendanceStatus = "mission"    // missione/trasferta
	StatusPermit    AttendanceStatus = "permit"     // permesso
	StatusSickLeave AttendanceStatus = "sick_leave" // malattia
	StatusOnStrike  AttendanceStatus = "on_strike"  // in sciopero
)

// StaffAttendance rappresenta la presenza di un membro del personale in un giorno
type StaffAttendance struct {
	ID       string `json:"id" db:"id"`
	SchoolID string `json:"school_id" db:"school_id"`
	UserID   string `json:"user_id" db:"user_id"`
	Date     string `json:"date" db:"date"` // YYYY-MM-DD

	Status AttendanceStatus `json:"status" db:"status"`

	// Timbratura badge
	BadgeEntryTime  *time.Time `json:"badge_entry_time,omitempty" db:"badge_entry_time"`
	BadgeExitTime   *time.Time `json:"badge_exit_time,omitempty" db:"badge_exit_time"`
	BadgeDeviceID   *string    `json:"badge_device_id,omitempty" db:"badge_device_id"`
	BadgeImportedAt *time.Time `json:"badge_imported_at,omitempty" db:"badge_imported_at"`

	// Sciopero (per docenti)
	IsStrikeDay       bool    `json:"is_strike_day" db:"is_strike_day"`
	StrikeConfirmedBy *string `json:"strike_confirmed_by,omitempty" db:"strike_confirmed_by"`

	Notes      string  `json:"notes,omitempty" db:"notes"`
	RecordedBy *string `json:"recorded_by,omitempty" db:"recorded_by"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Campi join (popolati via query)
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Role        string `json:"role,omitempty"`
	RoleDisplay string `json:"role_display,omitempty"` // nome leggibile del ruolo
}

// BadgeSwipe rappresenta una timbratura grezza dal dispositivo badge
type BadgeSwipe struct {
	ID          string     `json:"id" db:"id"`
	SchoolID    string     `json:"school_id" db:"school_id"`
	UserID      *string    `json:"user_id,omitempty" db:"user_id"`
	BadgeCode   string     `json:"badge_code" db:"badge_code"`
	DeviceID    string     `json:"device_id" db:"device_id"`
	SwipeTime   time.Time  `json:"swipe_time" db:"swipe_time"`
	SwipeType   string     `json:"swipe_type" db:"swipe_type"` // in, out, break_out, break_in
	RawData     *string    `json:"raw_data,omitempty" db:"raw_data"`
	Processed   bool       `json:"processed" db:"processed"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// UserBadge associa un badge fisico a un utente
type UserBadge struct {
	ID         string     `json:"id" db:"id"`
	SchoolID   string     `json:"school_id" db:"school_id"`
	UserID     string     `json:"user_id" db:"user_id"`
	BadgeCode  string     `json:"badge_code" db:"badge_code"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	AssignedAt time.Time  `json:"assigned_at" db:"assigned_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	Notes      string     `json:"notes,omitempty" db:"notes"`
}

// RoleSummary contiene il conteggio presenze per una categoria di personale
type RoleSummary struct {
	Role        string `json:"role"`
	RoleDisplay string `json:"role_display"`
	Total       int    `json:"total"`
	Present     int    `json:"present"`
	Absent      int    `json:"absent"`
	Late        int    `json:"late"`
	OnMission   int    `json:"on_mission"`
	OnLeave     int    `json:"on_leave"` // sick_leave + permit
	OnStrike    int    `json:"on_strike"`
}

// DailyStaffSummary è la risposta completa per la dashboard presenze giornaliera
type DailyStaffSummary struct {
	Date        string            `json:"date"`
	IsStrikeDay bool              `json:"is_strike_day"`
	ByRole      []RoleSummary     `json:"by_role"`
	TotalStaff  RoleSummary       `json:"total_staff"` // aggregato tutto il personale
	StaffList   []StaffAttendance `json:"staff_list"`  // lista nominativa con dettagli
}

// UpsertStaffAttendanceRequest è la richiesta per registrare/aggiornare una presenza
type UpsertStaffAttendanceRequest struct {
	UserID      string           `json:"user_id" binding:"required"`
	Date        string           `json:"date" binding:"required"` // YYYY-MM-DD
	Status      AttendanceStatus `json:"status" binding:"required"`
	Notes       string           `json:"notes"`
	IsStrikeDay bool             `json:"is_strike_day"`
}

// BadgeSwipeRequest è la richiesta per registrare una timbratura badge
type BadgeSwipeRequest struct {
	BadgeCode string `json:"badge_code" binding:"required"`
	DeviceID  string `json:"device_id" binding:"required"`
	SwipeTime string `json:"swipe_time"` // ISO8601, se vuoto usa NOW()
	SwipeType string `json:"swipe_type"` // in, out, break_out, break_in
	RawData   string `json:"raw_data"`
}

// BulkUpsertRequest permette di registrare presenze multiple in un'unica chiamata
type BulkUpsertRequest struct {
	Date        string                         `json:"date" binding:"required"`
	IsStrikeDay bool                           `json:"is_strike_day"`
	Attendances []UpsertStaffAttendanceRequest `json:"attendances" binding:"required"`
}

// RoleDisplayNames mappa i ruoli tecnici ai nomi leggibili in italiano
var RoleDisplayNames = map[string]string{
	"teacher":                   "Docente",
	"coordinator":               "Coordinatore di Classe",
	"dsga":                      "DSGA",
	"assistente_amministrativo": "Assistente Amministrativo (AA)",
	"collaboratore_ds":          "Collaboratore DS",
	"collaboratore_scolastico":  "Collaboratore Scolastico",
	"secretary":                 "Segreteria",
	"principal":                 "Dirigente Scolastico",
	"vice_principal":            "Vice Dirigente",
}
