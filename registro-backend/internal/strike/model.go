package strike

import (
	"time"
)

type StrikeIntention string

const (
	IntentionParticipates    StrikeIntention = "participates"
	IntentionNotParticipates StrikeIntention = "not_participates"
	IntentionUndecided       StrikeIntention = "undecided"
	IntentionUnanswered      StrikeIntention = "unanswered"
)

func IsValidIntention(val StrikeIntention) bool {
	switch val {
	case IntentionParticipates, IntentionNotParticipates, IntentionUndecided:
		return true
	default:
		return false
	}
}

type StrikeNotice struct {
	ID                  string             `json:"id" db:"id"`
	SchoolID            string             `json:"school_id" db:"school_id"`
	Title               string             `json:"title" db:"title"`
	ProclaimedBy        string             `json:"proclaimed_by" db:"proclaimed_by"`
	StrikeDate          string             `json:"strike_date" db:"strike_date"`
	DeclarationDeadline time.Time          `json:"declaration_deadline" db:"declaration_deadline"`
	Content             string             `json:"content" db:"content"`
	CreatedBy           string             `json:"created_by" db:"created_by"`
	CreatedByName       string             `json:"created_by_name,omitempty"`
	IsPublished         bool               `json:"is_published" db:"is_published"`
	CommunicationID     *string            `json:"communication_id,omitempty" db:"communication_id"`
	CreatedAt           time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" db:"updated_at"`
	IsExpired           bool               `json:"is_expired"`
	UserDeclaration     *StrikeDeclaration `json:"user_declaration,omitempty"`
}

type StrikeDeclaration struct {
	ID             string          `json:"id" db:"id"`
	StrikeNoticeID string          `json:"strike_notice_id" db:"strike_notice_id"`
	UserID         string          `json:"user_id" db:"user_id"`
	Intention      StrikeIntention `json:"intention" db:"intention"`
	DeclaredAt     time.Time       `json:"declared_at" db:"declared_at"`
	IPAddress      string          `json:"ip_address,omitempty" db:"ip_address"`
	Notes          string          `json:"notes,omitempty" db:"notes"`
}

type CreateStrikeNoticeRequest struct {
	Title               string `json:"title" binding:"required"`
	ProclaimedBy        string `json:"proclaimed_by" binding:"required"`
	StrikeDate          string `json:"strike_date" binding:"required"`
	DeclarationDeadline string `json:"declaration_deadline" binding:"required"`
	Content             string `json:"content" binding:"required"`
	PublishToBacheca    bool   `json:"publish_to_bacheca"`
}

type SubmitDeclarationRequest struct {
	Intention StrikeIntention `json:"intention" binding:"required"`
	Notes     string          `json:"notes"`
}

type RoleDeclarationSummary struct {
	Role            string  `json:"role"`
	RoleDisplay     string  `json:"role_display"`
	Total           int     `json:"total"`
	Participates    int     `json:"participates"`
	NotParticipates int     `json:"not_participates"`
	Undecided       int     `json:"undecided"`
	Unanswered      int     `json:"unanswered"`
	Percent         float64 `json:"percent"`
}

type NominativeDeclarationItem struct {
	UserID      string     `json:"user_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	RoleDisplay string     `json:"role_display"`
	Intention   string     `json:"intention"`
	DeclaredAt  *time.Time `json:"declared_at,omitempty"`
	IPAddress   string     `json:"ip_address,omitempty"`
	Notes       string     `json:"notes,omitempty"`
}

type StrikeNoticeSummaryResponse struct {
	Notice                 *StrikeNotice               `json:"notice"`
	TotalStaff             int                         `json:"total_staff"`
	TotalAnswered          int                         `json:"total_answered"`
	ParticipatesCount      int                         `json:"participates_count"`
	NotParticipatesCount   int                         `json:"not_participates_count"`
	UndecidedCount         int                         `json:"undecided_count"`
	UnansweredCount        int                         `json:"unanswered_count"`
	ParticipatesPercent    float64                     `json:"participates_percent"`
	NotParticipatesPercent float64                     `json:"not_participates_percent"`
	UndecidedPercent       float64                     `json:"undecided_percent"`
	UnansweredPercent      float64                     `json:"unanswered_percent"`
	ByRole                 []RoleDeclarationSummary    `json:"by_role"`
	Declarations           []NominativeDeclarationItem `json:"declarations"`
}
