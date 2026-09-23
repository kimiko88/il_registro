package subjects

import "time"

type Subject struct {
	ID             string    `json:"id" db:"id"`
	SchoolID       string    `json:"school_id" db:"school_id"`
	Name           string    `json:"name" db:"name"`
	Code           string    `json:"code,omitempty" db:"code"`
	Description    string    `json:"description,omitempty" db:"description"`
	IsMandatory    bool      `json:"is_mandatory" db:"is_mandatory"`
	// IsReligion indica che la materia è IRC (Insegnamento Religione Cattolica).
	// Quando true, accetta solo giudizi: Non classificabile, Insufficiente,
	// Sufficiente, Buono, Distinto, Ottimo.
	IsReligion     bool      `json:"is_religion" db:"is_religion"`
	// IsJudgmentOnly indica che la materia non accetta voti numerici.
	// Viene impostato automaticamente a true quando IsReligion è true.
	IsJudgmentOnly bool      `json:"is_judgment_only" db:"is_judgment_only"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type CreateSubjectRequest struct {
	Name           string `json:"name" binding:"required"`
	Code           string `json:"code"`
	Description    string `json:"description"`
	IsMandatory    *bool  `json:"is_mandatory"`    // Pointer to distinguish false vs nil (optional)
	IsReligion     *bool  `json:"is_religion"`     // Se true, materia IRC a soli giudizi
	IsJudgmentOnly *bool  `json:"is_judgment_only"` // Se true, no voti numerici standard
	SchoolID       string `json:"school_id"`       // Optional override for superadmin
}
