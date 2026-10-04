package parents

import "time"

type CustodyType string

const (
	CustodyShared     CustodyType = "shared"
	CustodySole       CustodyType = "sole"
	CustodyRestricted CustodyType = "restricted"
)

type DualAuthStatus string

const (
	StatusPendingFirst  DualAuthStatus = "pending_first"
	StatusPendingSecond DualAuthStatus = "pending_second"
	StatusCompleted     DualAuthStatus = "completed"
	StatusRejected      DualAuthStatus = "rejected"
)

type DualParentalAuthorization struct {
	ID                 string         `json:"id" db:"id"`
	StudentID          string         `json:"student_id" db:"student_id"`
	StudentName        string         `json:"student_name,omitempty" db:"student_name"`
	DocumentType       string         `json:"document_type" db:"document_type"`
	DocumentRefID      string         `json:"document_ref_id" db:"document_ref_id"`
	Title              string         `json:"title" db:"title"`
	Parent1ID          string         `json:"parent1_id" db:"parent1_id"`
	Parent1Name        string         `json:"parent1_name,omitempty" db:"parent1_name"`
	Parent1SignedAt    *time.Time     `json:"parent1_signed_at,omitempty" db:"parent1_signed_at"`
	Parent1PinVerified bool           `json:"parent1_pin_verified" db:"parent1_pin_verified"`
	Parent2ID          string         `json:"parent2_id,omitempty" db:"parent2_id"`
	Parent2Name        string         `json:"parent2_name,omitempty" db:"parent2_name"`
	Parent2SignedAt    *time.Time     `json:"parent2_signed_at,omitempty" db:"parent2_signed_at"`
	Parent2PinVerified bool           `json:"parent2_pin_verified" db:"parent2_pin_verified"`
	Status             DualAuthStatus `json:"status" db:"status"`
	RejectionReason    string         `json:"rejection_reason,omitempty" db:"rejection_reason"`
	Deadline           *time.Time     `json:"deadline,omitempty" db:"deadline"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
}

type CreateDualAuthParams struct {
	StudentID     string    `json:"student_id"`
	DocumentType  string    `json:"document_type"`
	DocumentRefID string    `json:"document_ref_id"`
	Title         string    `json:"title"`
	Parent1ID     string    `json:"parent1_id"`
	Parent2ID     string    `json:"parent2_id"`
	IsShared      bool      `json:"is_shared"`
	Deadline      time.Time `json:"deadline"`
}

type SignDualAuthRequest struct {
	PIN string `json:"pin" binding:"required"`
}

type RejectDualAuthRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type StudentCustodyInfo struct {
	StudentID               string      `json:"student_id"`
	ParentID                string      `json:"parent_id"`
	CustodyType             CustodyType `json:"custody_type"`
	CourtOrderDetails       string      `json:"court_order_details,omitempty"`
	CourtOrderDate          string      `json:"court_order_date,omitempty"`
	CanAuthorizeActivities  bool        `json:"can_authorize_activities"`
	IsMirrorNotified        bool        `json:"is_mirror_notified"`
}
