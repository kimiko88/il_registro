package psychology

import (
	"time"
)

type ConsentStatus string

const (
	ConsentPending  ConsentStatus = "in_attesa"
	ConsentApproved ConsentStatus = "approvato_entrambi"
	ConsentRejected ConsentStatus = "rifiutato"
)

type PsychologyConsent struct {
	ID              string        `json:"id" db:"id"`
	SchoolID        string        `json:"school_id" db:"school_id"`
	StudentID       string        `json:"student_id" db:"student_id"`
	SchoolYear      string        `json:"school_year" db:"school_year"`
	Parent1ID       string        `json:"parent1_id,omitempty" db:"parent1_id"`
	Parent1SignedAt *time.Time    `json:"parent1_signed_at,omitempty" db:"parent1_signed_at"`
	Parent2ID       string        `json:"parent2_id,omitempty" db:"parent2_id"`
	Parent2SignedAt *time.Time    `json:"parent2_signed_at,omitempty" db:"parent2_signed_at"`
	Status          ConsentStatus `json:"status" db:"status"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
}

type PsychologySlot struct {
	ID              string    `json:"id" db:"id"`
	SchoolID        string    `json:"school_id" db:"school_id"`
	PsychologistID  string    `json:"psychologist_id" db:"psychologist_id"`
	SlotTime        time.Time `json:"slot_time" db:"slot_time"`
	DurationMinutes int       `json:"duration_minutes" db:"duration_minutes"`
	IsBooked        bool      `json:"is_booked" db:"is_booked"`
}

type PsychologySession struct {
	ID                     string    `json:"id" db:"id"`
	SchoolID               string    `json:"school_id" db:"school_id"`
	StudentID              string    `json:"student_id" db:"student_id"`
	PsychologistID         string    `json:"psychologist_id" db:"psychologist_id"`
	AnonymousAlias         string    `json:"anonymous_alias" db:"anonymous_alias"`
	SlotTime               time.Time `json:"slot_time" db:"slot_time"`
	DurationMinutes        int       `json:"duration_minutes" db:"duration_minutes"`
	Status                 string    `json:"status" db:"status"`                 // "prenotato", "svolto", "annullato"
	EncryptedClinicalNotes string    `json:"encrypted_clinical_notes,omitempty"` // VISIBILE SOLO DALLO PSICOLOGO (L. 56/1989)
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
}

type BookSessionRequest struct {
	SlotTime        string `json:"slot_time" binding:"required"` // YYYY-MM-DD HH:MM
	DurationMinutes int    `json:"duration_minutes"`
	PsychologistID  string `json:"psychologist_id" binding:"required"`
}

type SaveConsentRequest struct {
	StudentID  string `json:"student_id" binding:"required"`
	SchoolYear string `json:"school_year" binding:"required"`
}
