package personnel_desk

import (
	"time"
)

type DeskCategory string

const (
	CatPermessoBreve    DeskCategory = "permesso_breve"
	CatFerie            DeskCategory = "ferie"
	CatMalattia         DeskCategory = "malattia"
	CatAspettativa      DeskCategory = "aspettativa"
	CatCongedoParentale DeskCategory = "congedo_parentale"
	CatPermessoStudio   DeskCategory = "permesso_studio"
	CatAltro            DeskCategory = "altro"
)

type DeskStatus string

const (
	StatusDraft      DeskStatus = "draft"
	StatusSubmitted  DeskStatus = "submitted"
	StatusAAReview   DeskStatus = "aa_review"
	StatusDSGAReview DeskStatus = "dsga_review"
	StatusDSReview   DeskStatus = "ds_review"
	StatusApproved   DeskStatus = "approved"
	StatusRejected   DeskStatus = "rejected"
)

type DeskRequest struct {
	ID          string       `json:"id"`
	SchoolID    string       `json:"school_id"`
	ApplicantID string       `json:"applicant_id"`
	ApplicantName string     `json:"applicant_name,omitempty"`
	ApplicantRole string     `json:"applicant_role,omitempty"`
	Category    DeskCategory `json:"category"`
	SubCategory string       `json:"sub_category,omitempty"`
	StartDate   string       `json:"start_date"`
	EndDate     string       `json:"end_date"`
	Days        float64      `json:"days"`
	Hours       float64      `json:"hours"`
	Description string       `json:"description"`
	Attachments []string     `json:"attachments"`

	Status DeskStatus `json:"status"`

	// Step 1: Istruttoria Assistente Amministrativo
	AANote       string     `json:"aa_note,omitempty"`
	AAReviewedBy *string    `json:"aa_reviewed_by,omitempty"`
	AAReviewedAt *time.Time `json:"aa_reviewed_at,omitempty"`

	// Step 2: Visto DSGA
	DSGANote     string     `json:"dsga_note,omitempty"`
	DSGASignedBy *string    `json:"dsga_signed_by,omitempty"`
	DSGASignedAt *time.Time `json:"dsga_signed_at,omitempty"`

	// Step 3: Approvazione Dirigente Scolastico
	DSDecreeNum  *string    `json:"ds_decree_num,omitempty"`
	DSNote       string     `json:"ds_note,omitempty"`
	DSApprovedBy *string    `json:"ds_approved_by,omitempty"`
	DSApprovedAt *time.Time `json:"ds_approved_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request DTOs
type CreateDeskRequestInput struct {
	Category    DeskCategory `json:"category" binding:"required"`
	SubCategory string       `json:"sub_category"`
	StartDate   string       `json:"start_date" binding:"required"`
	EndDate     string       `json:"end_date" binding:"required"`
	Days        float64      `json:"days"`
	Hours       float64      `json:"hours"`
	Description string       `json:"description" binding:"required"`
	Attachments []string     `json:"attachments"`
	SubmitNow   bool         `json:"submit_now"`
}

type AAReviewInput struct {
	Note    string `json:"note"`
	Approve bool   `json:"approve"` // true: passa a dsga_review, false: rejected
}

type DSGASignInput struct {
	Note    string `json:"note"`
	Approve bool   `json:"approve"` // true: passa a ds_review, false: rejected
}

type DSApproveInput struct {
	DecreeNum string `json:"decree_num"`
	Note      string `json:"note"`
	Approve   bool   `json:"approve"` // true: approved, false: rejected
}
