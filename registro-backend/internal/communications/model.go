package communications

import (
	"time"
)

type Message struct {
	ID                string     `json:"id" db:"id"`
	SchoolID          *string    `json:"school_id,omitempty" db:"school_id"`
	SenderID          string     `json:"sender_id" db:"sender_id"`
	ReceiverIDs       []string   `json:"receiver_ids" db:"receiver_ids"`
	Subject           string     `json:"subject" db:"subject"`
	Body              string     `json:"body" db:"body"`
	AttachmentURL     *string    `json:"attachment_url,omitempty" db:"attachment_url"`
	Type              string     `json:"type" db:"type"` // 'circular', 'notice', 'internal', 'email'
	RequiresSignature bool       `json:"requires_signature" db:"requires_signature"`
	SignatureDeadline *time.Time `json:"signature_deadline,omitempty" db:"signature_deadline"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	ReadAt            *time.Time `json:"read_at,omitempty" db:"read_at"`
	IsSigned          bool       `json:"is_signed"`
}

type CommunicationSignature struct {
	ID              string    `json:"id" db:"id"`
	CommunicationID string    `json:"communication_id" db:"communication_id"`
	UserID          string    `json:"user_id" db:"user_id"`
	UserName        string    `json:"user_name,omitempty"`
	UserRole        string    `json:"user_role,omitempty"`
	SignedAt        time.Time `json:"signed_at" db:"signed_at"`
	IPAddress       string    `json:"ip_address,omitempty" db:"ip_address"`
}

type SignatureReportResponse struct {
	CommunicationID   string                   `json:"communication_id"`
	Subject           string                   `json:"subject"`
	TotalRecipients   int                      `json:"total_recipients"`
	SignedCount       int                      `json:"signed_count"`
	PendingCount      int                      `json:"pending_count"`
	SignatureDeadline *time.Time               `json:"signature_deadline,omitempty"`
	Signatures        []CommunicationSignature `json:"signatures"`
}

type CreateMessageRequest struct {
	SchoolID          *string `json:"school_id,omitempty"`
	Recipients        []string`json:"recipients"`
	Subject           string  `json:"subject" binding:"required"`
	Body              string  `json:"body" binding:"required"`
	AttachmentURL     *string `json:"attachment_url,omitempty"`
	Type              string  `json:"type" binding:"required"` // 'circular', 'notice', 'internal'
	RequiresSignature bool    `json:"requires_signature"`
	SignatureDeadline *string `json:"signature_deadline,omitempty"` // YYYY-MM-DD or RFC3339
}
