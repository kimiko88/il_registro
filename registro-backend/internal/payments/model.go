package payments

import (
	"time"
)

type SchoolPayment struct {
	ID            string     `json:"id" db:"id"`
	SchoolID      string     `json:"school_id" db:"school_id"`
	StudentID     string     `json:"student_id" db:"student_id"`
	Title         string     `json:"title" db:"title"`
	Description   string     `json:"description" db:"description"`
	Amount        float64    `json:"amount" db:"amount"`
	DueDate       time.Time  `json:"due_date" db:"due_date"`
	Status        string     `json:"status" db:"status"` // 'pending', 'paid', 'cancelled'
	PaidAt        *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	PaymentMethod string     `json:"payment_method,omitempty" db:"payment_method"`
	TransactionID string     `json:"transaction_id,omitempty" db:"transaction_id"`
	PayerUserID   *string    `json:"payer_user_id,omitempty" db:"payer_user_id"`
	ReceiptNumber string     `json:"receipt_number,omitempty" db:"receipt_number"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`

	// Enriched fields for responses
	StudentName string `json:"student_name,omitempty"`
}

type CreatePaymentRequest struct {
	StudentID   string  `json:"student_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	DueDate     string  `json:"due_date" binding:"required"` // YYYY-MM-DD or RFC3339
}

type PayRequest struct {
	PaymentMethod string `json:"payment_method"` // 'PagoPA', 'credit_card', etc.
}

type PaymentSummaryResponse struct {
	PendingTotal float64          `json:"pending_total"`
	PendingCount int              `json:"pending_count"`
	PaidCount    int              `json:"paid_count"`
	Payments     []*SchoolPayment `json:"payments"`
}
