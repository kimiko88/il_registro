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

	// PagoPA & OPI Enterprise Fields
	IUV                  string     `json:"iuv,omitempty" db:"iuv"`
	QRCodePayload        string     `json:"qr_code_payload,omitempty" db:"qr_code_payload"`
	CheckoutSessionToken string     `json:"checkout_session_token,omitempty" db:"checkout_session_token"`
	ReconciledAt         *time.Time `json:"reconciled_at,omitempty" db:"reconciled_at"`
	SollecitoCount       int        `json:"sollecito_count" db:"sollecito_count"`
	UltimoSollecitoAt    *time.Time `json:"ultimo_sollecito_at,omitempty" db:"ultimo_sollecito_at"`

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

type BollettinoNotice struct {
	PaymentID         string    `json:"payment_id"`
	IUV               string    `json:"iuv"`
	SchoolCF          string    `json:"school_cf"`
	SchoolName        string    `json:"school_name"`
	StudentName       string    `json:"student_name"`
	Title             string    `json:"title"`
	Amount            float64   `json:"amount"`
	DueDate           time.Time `json:"due_date"`
	QRCodePayload     string    `json:"qr_code_payload"`
	Barcode128        string    `json:"barcode_128"`
	CausaleVersamento string    `json:"causale_versamento"`
	GeneratedAt       time.Time `json:"generated_at"`
}

type CheckoutSession struct {
	PaymentID    string    `json:"payment_id"`
	SessionToken string    `json:"session_token"`
	RedirectURL  string    `json:"redirect_url"`
	ExpiresAt    time.Time `json:"expires_at"`
	Amount       float64   `json:"amount"`
	IUV          string    `json:"iuv"`
}

type OPIQuietanza struct {
	CodiceFlusso    string    `json:"codice_flusso" xml:"codiceFlusso"`
	NumeroQuietanza string    `json:"numero_quietanza" xml:"numeroQuietanza"`
	IUV             string    `json:"iuv" xml:"iuv"`
	Importo         float64   `json:"importo" xml:"importo"`
	DataAccredito   time.Time `json:"data_accredito" xml:"dataAccredito"`
	CodiceDebitore  string    `json:"codice_debitore" xml:"codiceDebitore"`
	NomeDebitore    string    `json:"nome_debitore" xml:"nomeDebitore"`
	Esito           string    `json:"esito" xml:"esito"`
}

type ReconciliationReport struct {
	CodiceFlusso       string           `json:"codice_flusso"`
	TotalProcessed     int              `json:"total_processed"`
	TotalReconciled    int              `json:"total_reconciled"`
	TotalUnmatched     int              `json:"total_unmatched"`
	TotalAmount        float64          `json:"total_amount"`
	ReconciledPayments []*SchoolPayment `json:"reconciled_payments"`
	UnmatchedQuietanze []OPIQuietanza   `json:"unmatched_quietanze"`
}

type PaymentSollecito struct {
	PaymentID      string    `json:"payment_id"`
	StudentID      string    `json:"student_id"`
	StudentName    string    `json:"student_name"`
	ParentEmail    string    `json:"parent_email"`
	Title          string    `json:"title"`
	Amount         float64   `json:"amount"`
	DueDate        time.Time `json:"due_date"`
	DaysOverdue    int       `json:"days_overdue"`
	SollecitoCount int       `json:"sollecito_count"`
	SollecitoText  string    `json:"sollecito_text"`
}
