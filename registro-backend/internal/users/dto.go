package users

import (
	"time"
)

type CreateUserRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=10"`
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	FiscalCode  string  `json:"fiscal_code" binding:"omitempty,len=16"`
	Role        string  `json:"role" binding:"required"`
	SchoolID    *string `json:"school_id"`
	PhoneNumber string  `json:"phone_number"`
	JobTitle    string  `json:"job_title"`
	ClassID     *string `json:"class_id"`      // For students
	DateOfBirth string  `json:"date_of_birth"` // Format: YYYY-MM-DD
	IsStaff     *bool   `json:"is_staff"`
}

type UpdateUserRequest struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	FiscalCode  *string `json:"fiscal_code" binding:"omitempty,len=16"`
	PhoneNumber *string `json:"phone_number"`
	JobTitle    *string `json:"job_title"`
	IsActive    *bool   `json:"is_active"`
	IsStaff     *bool   `json:"is_staff"`
	Role        *string `json:"role"`
	SchoolID    *string `json:"school_id"`
	ClassID     *string `json:"class_id"`      // For students
	DateOfBirth *string `json:"date_of_birth"` // Format: YYYY-MM-DD
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=10"`
}

type AssignRolesRequest struct {
	Role string `json:"role" binding:"required"`
}

// UserAssignment rappresenta un incarico aggiuntivo/funzionale assegnato all'utente
type UserAssignment struct {
	ID             string     `json:"id" db:"id"`
	SchoolID       string     `json:"school_id" db:"school_id"`
	UserID         string     `json:"user_id" db:"user_id"`
	AssignmentType string     `json:"assignment_type" db:"assignment_type"`
	ScopeType      string     `json:"scope_type" db:"scope_type"`
	ScopeID        *string    `json:"scope_id,omitempty" db:"scope_id"`
	Title          string     `json:"title" db:"title"`
	AssignedBy     *string    `json:"assigned_by,omitempty" db:"assigned_by"`
	Metadata       string     `json:"metadata,omitempty" db:"metadata"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	ValidFrom      *time.Time `json:"valid_from,omitempty" db:"valid_from"`
	ValidTo        *time.Time `json:"valid_to,omitempty" db:"valid_to"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateAssignmentRequest struct {
	AssignmentType string  `json:"assignment_type" binding:"required"`
	ScopeType      string  `json:"scope_type"`
	ScopeID        *string `json:"scope_id"`
	Title          string  `json:"title"`
	Metadata       string  `json:"metadata"`
}

type SetCoordinatedClassesRequest struct {
	ClassIDs []string `json:"class_ids" binding:"required"`
}

type UserResponse struct {
	ID              string           `json:"id"`
	Email           string           `json:"email"`
	FirstName       string           `json:"first_name"`
	LastName        string           `json:"last_name"`
	FiscalCode      *string          `json:"fiscal_code"`
	Role            string           `json:"role"`
	SchoolID        *string          `json:"school_id,omitempty"`
	ClassID         *string          `json:"class_id,omitempty"`   // For students
	ClassName       *string          `json:"class_name,omitempty"` // For students
	IsActive        bool             `json:"is_active"`
	EmailVerified   bool             `json:"email_verified"`
	MFAEnabled      bool             `json:"mfa_enabled"`
	PhoneNumber     *string          `json:"phone_number,omitempty"`
	JobTitle        *string          `json:"job_title,omitempty"`
	DateOfBirth     *string          `json:"date_of_birth,omitempty"` // YYYY-MM-DD string for easy display
	CreatedAt       time.Time        `json:"created_at"`
	LastLogin       *time.Time       `json:"last_login,omitempty"`
	DeletedAt       *time.Time       `json:"deleted_at,omitempty"`
	PseudonymizedAt *time.Time       `json:"pseudonymized_at,omitempty"`
	Assignments     []UserAssignment `json:"assignments,omitempty"`
}

type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

type BulkImportRequest struct {
	Format string `json:"format" binding:"required,oneof=csv xlsx"`
	// This structure is mostly for API documentation;
	// actual file upload is handled via multipart form
}

type ImportResult struct {
	Total       int      `json:"total"`
	Created     int      `json:"created"`
	Updated     int      `json:"updated"`
	Failed      int      `json:"failed"`
	Errors      []string `json:"errors,omitempty"`
	DownloadURL string   `json:"download_url,omitempty"` // URL to download generic import template/report
}

type BulkDeleteRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1"`
}
