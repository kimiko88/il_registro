package pdp

import (
	"time"
)

// PlanType distinguishes PDP (Piano Didattico Personalizzato) for BES/DSA
// from PEI (Piano Educativo Individualizzato) for students with disabilities.
type PlanType string

const (
	PlanTypePDP PlanType = "pdp"
	PlanTypePEI PlanType = "pei"
)

// PdpContent is the structured content stored as JSONB.
type PdpContent struct {
	Objectives      []string `json:"objectives"`       // Obiettivi didattici
	Compensative    []string `json:"compensative"`     // Misure compensative
	Dispensative    []string `json:"dispensative"`     // Misure dispensative
	EvaluationTools []string `json:"evaluation_tools"` // Strumenti di valutazione
	Notes           string   `json:"notes"`            // Note aggiuntive
	ReviewDate      string   `json:"review_date"`      // Data revisione (YYYY-MM-DD)
}

// PdpPlan is the main domain model for a PDP or PEI plan.
type PdpPlan struct {
	ID           string `json:"id"            db:"id"`
	StudentID    string `json:"student_id"    db:"student_id"`
	ClassID      string `json:"class_id"      db:"class_id"`
	SchoolID     string `json:"school_id"     db:"school_id"`
	AcademicYear string `json:"academic_year" db:"academic_year"`

	PlanType  PlanType   `json:"plan_type"  db:"plan_type"`
	Diagnosis string     `json:"diagnosis"  db:"diagnosis"` // Riservato ai docenti
	Content   PdpContent `json:"content"    db:"content"`

	CoordinatorID *string `json:"coordinator_id,omitempty" db:"coordinator_id"`
	ReferenteID   *string `json:"referente_id,omitempty"   db:"referente_id"`

	// Sharing state
	SharedWithFamily bool       `json:"shared_with_family"  db:"shared_with_family"`
	FamilyApprovedAt *time.Time `json:"family_approved_at,omitempty" db:"family_approved_at"`
	FamilyApprovedBy *string    `json:"family_approved_by,omitempty" db:"family_approved_by"`

	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields (not in db columns)
	StudentName     string `json:"student_name,omitempty"     db:"-"`
	CoordinatorName string `json:"coordinator_name,omitempty" db:"-"`
}

// CreatePdpRequest is the request body for creating a new plan.
type CreatePdpRequest struct {
	StudentID     string     `json:"student_id"    binding:"required"`
	ClassID       string     `json:"class_id"      binding:"required"`
	AcademicYear  string     `json:"academic_year" binding:"required"`
	PlanType      PlanType   `json:"plan_type"`
	Diagnosis     string     `json:"diagnosis"`
	Content       PdpContent `json:"content"`
	CoordinatorID *string    `json:"coordinator_id"`
	ReferenteID   *string    `json:"referente_id"`
}

// UpdatePdpRequest is the request body for updating an existing plan.
type UpdatePdpRequest struct {
	PlanType      *PlanType   `json:"plan_type"`
	Diagnosis     *string     `json:"diagnosis"`
	Content       *PdpContent `json:"content"`
	CoordinatorID *string     `json:"coordinator_id"`
	ReferenteID   *string     `json:"referente_id"`
}

// ShareWithFamilyRequest controls the family sharing flag.
type ShareWithFamilyRequest struct {
	Share bool `json:"share"`
}

// ApprovePdpRequest is used by a parent/guardian to acknowledge the plan.
type ApprovePdpRequest struct {
	// Optionally include a comment from the family
	Comment string `json:"comment"`
}

// Standard compensative measure keys (Italian school vocabulary)
var StandardCompensativeMeasures = []string{
	"calcolatrice",
	"tavola_pitagorica",
	"tabelle_formule",
	"dizionario_ortografico",
	"sintesi_vocale",
	"registratore_audio",
	"tempo_aggiuntivo_30",
	"tempo_aggiuntivo_50",
	"prova_equipollente",
	"mappe_concettuali",
	"testo_ingrandito",
	"font_leggibilita",
}

// Standard dispensative measure keys
var StandardDispensativeMeasures = []string{
	"lettura_ad_alta_voce",
	"scrittura_veloce",
	"copiatura_lavagna",
	"studio_mnemonico",
	"effettuazione_piu_prove_contemporanee",
}
