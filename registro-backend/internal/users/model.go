package users

import (
	"time"
)

// User represents the full user entity for management
type User struct {
	ID            string  `json:"id" db:"id"`
	Email         string  `json:"email" db:"email"`
	PasswordHash  string  `json:"-" db:"password_hash"`
	FirstName     string  `json:"first_name" db:"first_name"`
	LastName      string  `json:"last_name" db:"last_name"`
	FiscalCode    string  `json:"fiscal_code" db:"fiscal_code"` // Codice Fiscale
	Role          string  `json:"role" db:"role"`
	SchoolID      *string `json:"school_id,omitempty" db:"school_id"`
	IsActive      bool    `json:"is_active" db:"is_active"`
	EmailVerified bool    `json:"email_verified" db:"email_verified"`
	MFAEnabled    bool    `json:"mfa_enabled" db:"mfa_enabled"`
	MFASecret     string  `json:"-" db:"mfa_secret"`
	PhoneNumber   string  `json:"phone_number" db:"phone_number"`
	JobTitle      string  `json:"job_title" db:"job_title"`

	// Student specific (populated via joins)
	ClassID   *string `json:"class_id,omitempty" db:"class_id"`
	ClassName *string `json:"class_name,omitempty" db:"class_name"`

	// Audit & Lifecycle
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty" db:"last_login"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"` // Soft Delete

	// GDPR
	PseudonymizedAt *time.Time `json:"pseudonymized_at,omitempty" db:"pseudonymized_at"`
}

// AuditLog tracks specific actions on user records
type AuditLog struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`   // The user being affected
	ActorID   string    `json:"actor_id" db:"actor_id"` // Who performed the action
	Action    string    `json:"action" db:"action"`     // e.g., "create", "update_role"
	Details   string    `json:"details" db:"details"`   // JSON details of change
	IPAddress string    `json:"ip_address" db:"ip_address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserFilter defines criteria for searching/filtering users
type UserFilter struct {
	Query        string // Search by name, email, fiscal_code
	Role         string
	ExcludeRoles []string
	SchoolID     *string
	IsActive     *bool
	IsDeleted    bool // If true, include soft-deleted users
	Page         int
	PageSize     int
	SortBy       string
	SortOrder    string // ASC or DESC
}
type StudentChild struct {
	ID         string `json:"id"`      // Student Profile ID
	UserID     string `json:"user_id"` // Student User ID
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Class      string `json:"class"`
	SchoolName string `json:"school_name"`
}

type GuardianInfo struct {
	ID               string `json:"id"`
	ParentUserID     string `json:"parent_user_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	RelationshipType string `json:"relationship_type"`
}
