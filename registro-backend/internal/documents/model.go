package documents

import (
	"time"

	"github.com/google/uuid"
)

// Enums
type DocStatus string

const (
	StatusDraft     DocStatus = "draft"
	StatusSubmitted DocStatus = "submitted"
	StatusReview    DocStatus = "review"
	StatusApproved  DocStatus = "approved"
	StatusSigned    DocStatus = "signed"
	StatusArchived  DocStatus = "archived"
	StatusRejected  DocStatus = "rejected"
)

type DocType string

const (
	TypePDP         DocType = "pdp"
	TypePFI         DocType = "pfi"
	TypePFP         DocType = "pfp"
	TypeMay15       DocType = "may15"
	TypePCTO        DocType = "pcto"
	TypeOrientation DocType = "orientation"
	TypeCertificate DocType = "certificate"
	TypeGeneric     DocType = "generic"
)

// Document represents the metadata header
type Document struct {
	ID       string `json:"id" db:"id"`
	SchoolID string `json:"school_id" db:"school_id"`

	Title string  `json:"title" db:"title"`
	Type  DocType `json:"type" db:"type"`

	StudentID *string `json:"student_id,omitempty" db:"student_id"`
	ClassID   *string `json:"class_id,omitempty" db:"class_id"`

	Status         DocStatus `json:"status" db:"status"`
	CurrentVersion int       `json:"current_version" db:"current_version"`

	IsSigned bool       `json:"is_signed" db:"is_signed"`
	SignedBy *string    `json:"signed_by,omitempty" db:"signed_by"`
	SignedAt *time.Time `json:"signed_at,omitempty" db:"signed_at"`

	CreatedBy string     `json:"created_by" db:"created_by"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// DocumentVersion stores historical content
type DocumentVersion struct {
	ID         string `json:"id" db:"id"`
	DocumentID string `json:"document_id" db:"document_id"`

	VersionNum int    `json:"version_num" db:"version_num"`
	Content    string `json:"content" db:"content"` // Rich Text / JSON

	ChangeLog string `json:"change_log,omitempty" db:"change_log"`

	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DocumentSignature stores digital signature data
type DocumentSignature struct {
	ID         string `json:"id" db:"id"`
	DocumentID string `json:"document_id" db:"document_id"`
	VersionID  string `json:"version_id" db:"version_id"`

	SignerID      string `json:"signer_id" db:"signer_id"`
	SignatureData string `json:"-" db:"signature_data"` // Verify only, don't expose raw often
	Certificate   string `json:"-" db:"certificate_data"`

	SignedAt time.Time `json:"signed_at" db:"signed_at"`
}

// DocumentTemplate represents reusable templates
type DocumentTemplate struct {
	ID       string  `json:"id" db:"id"`
	SchoolID string  `json:"school_id" db:"school_id"`
	Name     string  `json:"name" db:"name"`
	Type     DocType `json:"type" db:"type"`
	Content  string  `json:"content" db:"content"`
	IsActive bool    `json:"is_active" db:"is_active"`
}

func (d *Document) IsValid() bool {
	if d.Title == "" || d.SchoolID == "" || d.CreatedBy == "" {
		return false
	}
	if d.StudentID != nil {
		if _, err := uuid.Parse(*d.StudentID); err != nil {
			return false
		}
	}
	return true
}
