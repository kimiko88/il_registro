package documents

import "time"

// Requests

type CreateDocumentRequest struct {
	Title      string  `json:"title" binding:"required"`
	Type       DocType `json:"type" binding:"required"`
	Content    string  `json:"content" binding:"required"`
	StudentID  *string `json:"student_id,omitempty"`
	ClassID    *string `json:"class_id,omitempty"`
	TemplateID *string `json:"template_id,omitempty"` // Create from template
}

type UpdateDocumentRequest struct {
	Title     *string `json:"title,omitempty"`
	Content   *string `json:"content,omitempty"` // New version
	ChangeLog string  `json:"change_log,omitempty"`
}

type SignDocumentRequest struct {
	CertificateData string `json:"certificate_data" binding:"required"` // Base64
	SignatureData   string `json:"signature_data" binding:"required"`   // Base64
}

type WorkflowActionRequest struct {
	Action string `json:"action" binding:"required"` // submit, reject, approve
	Reason string `json:"reason,omitempty"`
}

type TemplateRequest struct {
	Name    string  `json:"name" binding:"required"`
	Type    DocType `json:"type" binding:"required"`
	Content string  `json:"content" binding:"required"`
}

// Responses

type DocumentListResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Type      DocType   `json:"type"`
	Status    DocStatus `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	IsSigned  bool      `json:"is_signed"`
}

type DocumentDetailResponse struct {
	DocumentListResponse
	Content        string             `json:"content"`
	CurrentVersion int                `json:"current_version"`
	Versions       []VersionSummary   `json:"versions,omitempty"`
	Signatures     []SignatureSummary `json:"signatures,omitempty"`
}

type VersionSummary struct {
	VersionNum int       `json:"version_num"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	ChangeLog  string    `json:"change_log"`
}

type SignatureSummary struct {
	SignedBy string    `json:"signed_by"`
	SignedAt time.Time `json:"signed_at"`
}
