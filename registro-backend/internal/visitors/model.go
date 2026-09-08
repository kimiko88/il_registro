package visitors

import "time"

// VisitorPurpose categorizza il motivo della visita
type VisitorPurpose string

const (
	PurposeParent       VisitorPurpose = "parent"        // genitore/tutore
	PurposeSupplier     VisitorPurpose = "supplier"      // fornitore/tecnico
	PurposeInstitution  VisitorPurpose = "institution"   // ente/istituzione
	PurposeOther        VisitorPurpose = "other"         // altro
)

// MaintenancePriority livelli di priorità per le segnalazioni guasti
type MaintenancePriority string

const (
	PriorityLow    MaintenancePriority = "bassa"
	PriorityMedium MaintenancePriority = "media"
	PriorityHigh   MaintenancePriority = "alta"
	PriorityUrgent MaintenancePriority = "urgente"
)

// MaintenanceStatus stati di avanzamento segnalazione
type MaintenanceStatus string

const (
	MaintenanceOpen       MaintenanceStatus = "aperto"
	MaintenanceInProgress MaintenanceStatus = "in_lavorazione"
	MaintenanceClosed     MaintenanceStatus = "chiuso"
)

// Visitor registra l'ingresso di un visitatore esterno
type Visitor struct {
	ID          string         `json:"id" db:"id"`
	SchoolID    string         `json:"school_id" db:"school_id"`
	Name        string         `json:"name" db:"name"`
	DocumentID  string         `json:"document_id,omitempty" db:"document_id"`
	Purpose     VisitorPurpose `json:"purpose" db:"purpose"`
	HostName    string         `json:"host_name,omitempty" db:"host_name"` // persona da incontrare
	BadgeNumber *string        `json:"badge_number,omitempty" db:"badge_number"`
	EntryTime   time.Time      `json:"entry_time" db:"entry_time"`
	ExitTime    *time.Time     `json:"exit_time,omitempty" db:"exit_time"`
	Notes       string         `json:"notes,omitempty" db:"notes"`
	RecordedBy  string         `json:"recorded_by" db:"recorded_by"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
}

// EarlyExit registra l'uscita anticipata di uno studente
type EarlyExit struct {
	ID             string     `json:"id" db:"id"`
	SchoolID       string     `json:"school_id" db:"school_id"`
	StudentID      string     `json:"student_id" db:"student_id"`
	StudentName    string     `json:"student_name,omitempty"` // join
	ClassName      string     `json:"class_name,omitempty"`   // join
	ExitTime       time.Time  `json:"exit_time" db:"exit_time"`
	ReturnTime     *time.Time `json:"return_time,omitempty" db:"return_time"`
	DelegateeName  string     `json:"delegatee_name" db:"delegatee_name"`     // nome delegato
	DelegateRel    string     `json:"delegate_rel,omitempty" db:"delegate_rel"` // genitore, nonno, tutore, ...
	ReasonCode     string     `json:"reason_code,omitempty" db:"reason_code"`   // medica, famiglia, altro
	Notes          string     `json:"notes,omitempty" db:"notes"`
	RecordedBy     string     `json:"recorded_by" db:"recorded_by"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// MaintenanceReport segnala guasti o anomalie strutturali
type MaintenanceReport struct {
	ID          string              `json:"id" db:"id"`
	SchoolID    string              `json:"school_id" db:"school_id"`
	Location    string              `json:"location" db:"location"`       // es. "Aula 12", "Bagni piano 1"
	Category    string              `json:"category" db:"category"`       // elettrico, idraulico, strutturale, pulizia, informatica
	Description string              `json:"description" db:"description"`
	Priority    MaintenancePriority `json:"priority" db:"priority"`
	Status      MaintenanceStatus   `json:"status" db:"status"`
	ReportedBy  string              `json:"reported_by" db:"reported_by"`
	AssignedTo  *string             `json:"assigned_to,omitempty" db:"assigned_to"`
	ClosedAt    *time.Time          `json:"closed_at,omitempty" db:"closed_at"`
	CreatedAt   time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" db:"updated_at"`
}

// --- Request DTOs ---

type RegisterVisitorRequest struct {
	Name        string         `json:"name" binding:"required"`
	DocumentID  string         `json:"document_id"`
	Purpose     VisitorPurpose `json:"purpose" binding:"required"`
	HostName    string         `json:"host_name"`
	BadgeNumber string         `json:"badge_number"`
	Notes       string         `json:"notes"`
}

type RecordExitRequest struct {
	Notes string `json:"notes"`
}

type RecordEarlyExitRequest struct {
	StudentID     string `json:"student_id" binding:"required"`
	DelegateeName string `json:"delegatee_name" binding:"required"`
	DelegateRel   string `json:"delegate_rel"`
	ReasonCode    string `json:"reason_code"`
	Notes         string `json:"notes"`
}

type RecordStudentReturnRequest struct {
	Notes string `json:"notes"`
}

type CreateMaintenanceReportRequest struct {
	Location    string              `json:"location" binding:"required"`
	Category    string              `json:"category" binding:"required"`
	Description string              `json:"description" binding:"required"`
	Priority    MaintenancePriority `json:"priority"`
}

type UpdateMaintenanceStatusRequest struct {
	Status     MaintenanceStatus `json:"status" binding:"required"`
	AssignedTo string            `json:"assigned_to"`
}

// RolesAllowed sono i ruoli che possono accedere al registro visitatori
var RolesAllowed = []string{
	"collaboratore_scolastico",
	"collaboratore_ds",
	"dsga",
	"assistente_amministrativo",
	"secretary",
	"principal",
	"vice_principal",
	"admin",
	"superadmin",
}

func CanAccessVisitorRegistry(role string) bool {
	for _, r := range RolesAllowed {
		if r == role {
			return true
		}
	}
	return false
}
