package orientamento

import (
	"time"
)

type Event struct {
	ID           string    `json:"id" db:"id"`
	SchoolID     string    `json:"school_id" db:"school_id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Category     string    `json:"category" db:"category"` // University, Work
	Date         time.Time `json:"date" db:"date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
	Location     string    `json:"location" db:"location"`
	Hours        float64   `json:"hours" db:"hours"`
	MaxAttendees int       `json:"max_attendees" db:"max_attendees"`
	CreatedBy    string    `json:"created_by" db:"created_by"`
}

type Participation struct {
	ID           string    `json:"id" db:"id"`
	EventID      string    `json:"event_id" db:"event_id"`
	StudentID    string    `json:"student_id" db:"student_id"`
	Status       string    `json:"status" db:"status"` // Registered, Attended
	Attended     bool      `json:"attended" db:"attended"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
	Event        *Event    `json:"event,omitempty"`
}

// Capolavoro represents the student's masterpiece according to MIM ministerial orientation guidelines
type Capolavoro struct {
	ID                   string    `json:"id" db:"id"`
	StudentID            string    `json:"student_id" db:"student_id"`
	SchoolYear           string    `json:"school_year" db:"school_year"`
	Title                string    `json:"title" db:"title"`
	Description          string    `json:"description" db:"description"`
	CompetenzeSviluppate []string  `json:"competenze_sviluppate" db:"competenze_sviluppate"`
	ReflectiveNotes      string    `json:"reflective_notes" db:"reflective_notes"`
	AttachmentURL        string    `json:"attachment_url" db:"attachment_url"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

// CurriculumStudenteSummary provides a full dossier for the state exam commission
type CurriculumStudenteSummary struct {
	StudentID                string       `json:"student_id"`
	StudentName              string       `json:"student_name"`
	SchoolName               string       `json:"school_name"`
	PercorsoScolastico       []string     `json:"percorso_scolastico"`
	Certificazioni           []string     `json:"certificazioni"`
	AttivitaExtracurriculari []string     `json:"attivita_extracurriculari"`
	PCTOHours                float64      `json:"pcto_hours"`
	OrientamentoHours        float64      `json:"orientamento_hours"`
	Capolavori               []Capolavoro `json:"capolavori"`
	GeneratedAt              time.Time    `json:"generated_at"`
}

// DTOs
type CreateEventRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Category     string  `json:"category" binding:"required"`
	Date         string  `json:"date" binding:"required"`
	EndDate      string  `json:"end_date" binding:"required"`
	Location     string  `json:"location"`
	Hours        float64 `json:"hours" binding:"required"`
	MaxAttendees int     `json:"max_attendees"`
}

type StudentPreference struct {
	ID             string    `json:"id" db:"id"`
	StudentID      string    `json:"student_id" db:"student_id"`
	PreferredTrack string    `json:"preferred_track" db:"preferred_track"`
	TargetField    string    `json:"target_field" db:"target_field"`
	Notes          string    `json:"notes" db:"notes"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
