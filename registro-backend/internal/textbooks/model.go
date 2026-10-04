package textbooks

import "time"

type Textbook struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Title     string    `json:"title" db:"title"`
	Author    string    `json:"author" db:"author"`
	Subject   string    `json:"subject" db:"subject"`
	ISBN      string    `json:"isbn" db:"isbn"`
	Publisher string    `json:"publisher" db:"publisher"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ClassTextbook struct {
	ID         string `json:"id" db:"id"`
	ClassID    string `json:"class_id" db:"class_id"`
	SubjectID  string `json:"subject_id" db:"subject_id"`
	TextbookID string `json:"textbook_id" db:"textbook_id"`
	IsOptional bool   `json:"is_optional" db:"is_optional"`

	// Joined fields
	Title       string `json:"title,omitempty"`
	Author      string `json:"author,omitempty"`
	SubjectName string `json:"subject_name,omitempty"`
}

type CreateTextbookRequest struct {
	Title     string  `json:"title" binding:"required"`
	Author    string  `json:"author"`
	Subject   string  `json:"subject"`
	ISBN      string  `json:"isbn"`
	Publisher string  `json:"publisher"`
	Price     float64 `json:"price"`
}

type AssignTextbookRequest struct {
	SubjectID  string `json:"subject_id" binding:"required"`
	TextbookID string `json:"textbook_id" binding:"required"`
	IsOptional bool   `json:"is_optional"`
}

// AIE Catalog & Spending Limits Models

type AIECatalogBook struct {
	ID            string    `json:"id" db:"id"`
	ISBN          string    `json:"isbn" db:"isbn"`
	Title         string    `json:"title" db:"title"`
	Authors       string    `json:"authors" db:"authors"`
	Publisher     string    `json:"publisher" db:"publisher"`
	Subject       string    `json:"subject" db:"subject"`
	Price         float64   `json:"price" db:"price"`
	Volume        string    `json:"volume" db:"volume"`
	EditionYear   int       `json:"edition_year" db:"edition_year"`
	SchoolOrder   string    `json:"school_order" db:"school_order"`
	IsDigitalOnly bool      `json:"is_digital_only" db:"is_digital_only"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type SpendingLimit struct {
	ID                  string    `json:"id" db:"id"`
	SchoolID            string    `json:"school_id" db:"school_id"`
	ClassYear           int       `json:"class_year" db:"class_year"`
	SchoolOrder         string    `json:"school_order" db:"school_order"`
	MaxAmount           float64   `json:"max_amount" db:"max_amount"`
	AllowedTolerancePct float64   `json:"allowed_tolerance_pct" db:"allowed_tolerance_pct"`
	AcademicYear        string    `json:"academic_year" db:"academic_year"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

type SpendingStatus string

const (
	StatusWithinLimit      SpendingStatus = "WITHIN_LIMIT"
	StatusWarningTolerance SpendingStatus = "WARNING_TOLERANCE"
	StatusOverLimit        SpendingStatus = "OVER_LIMIT"
)

type ClassAdoptionItem struct {
	ID             string     `json:"id" db:"id"`
	ClassID        string     `json:"class_id" db:"class_id"`
	SubjectID      string     `json:"subject_id" db:"subject_id"`
	SubjectName    string     `json:"subject_name,omitempty" db:"subject_name"`
	BookID         string     `json:"book_id" db:"book_id"`
	BookTitle      string     `json:"book_title,omitempty" db:"book_title"`
	Authors        string     `json:"authors,omitempty" db:"authors"`
	Publisher      string     `json:"publisher,omitempty" db:"publisher"`
	ISBN           string     `json:"isbn,omitempty" db:"isbn"`
	Price          float64    `json:"price" db:"price"`
	AdoptionType   string     `json:"adoption_type" db:"adoption_type"` // nuova_adozione, scorrimento, consigliato
	IsAlreadyOwned bool       `json:"is_already_owned" db:"is_already_owned"`
	IsMonographic  bool       `json:"is_monographic" db:"is_monographic"`
	DeliberatedAt  *time.Time `json:"deliberated_at,omitempty" db:"deliberated_at"`
	Notes          string     `json:"notes,omitempty" db:"notes"`
}

type SpendingReport struct {
	TotalSpending       float64           `json:"total_spending"`
	SpendingLimit       float64           `json:"spending_limit"`
	ToleranceThreshold  float64           `json:"tolerance_threshold"`
	AllowedTolerancePct float64           `json:"allowed_tolerance_pct"`
	Difference          float64           `json:"difference"`
	Percentage          float64           `json:"percentage"`
	Status              SpendingStatus    `json:"status"`
	AdoptionsCount      int               `json:"adoptions_count"`
	Items               []ClassAdoptionItem `json:"items"`
}

type AIEExportRecord struct {
	SchoolCode   string  `json:"school_code"`
	AcademicYear string  `json:"academic_year"`
	ClassName    string  `json:"class_name"`
	SubjectCode  string  `json:"subject_code"`
	ISBN         string  `json:"isbn"`
	Title        string  `json:"title"`
	Authors      string  `json:"authors"`
	Publisher    string  `json:"publisher"`
	Price        float64 `json:"price"`
	Volume       string  `json:"volume"`
	AdoptionType string  `json:"adoption_type"`
	AlreadyOwned bool    `json:"already_owned"`
}
