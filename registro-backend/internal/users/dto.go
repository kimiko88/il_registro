package users

import (
	"time"
)

type CreateUserRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=8"`
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	FiscalCode  string  `json:"fiscal_code" binding:"omitempty,len=16"`
	Role        string  `json:"role" binding:"required,oneof=admin principal secretary teacher student parent"`
	SchoolID    *string `json:"school_id"`
	PhoneNumber string  `json:"phone_number"`
	JobTitle    string  `json:"job_title"`
	ClassID     *string `json:"class_id"` // For students
}

type UpdateUserRequest struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	FiscalCode  *string `json:"fiscal_code" binding:"omitempty,len=16"`
	PhoneNumber *string `json:"phone_number"`
	JobTitle    *string `json:"job_title"`
	IsActive    *bool   `json:"is_active"`
	Role        *string `json:"role"`
	SchoolID    *string `json:"school_id"`
	ClassID     *string `json:"class_id"` // For students
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type AssignRolesRequest struct {
	Role string `json:"role" binding:"required"`
}

type UserResponse struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	FiscalCode      *string    `json:"fiscal_code"`
	Role            string     `json:"role"`
	SchoolID        *string    `json:"school_id,omitempty"`
	ClassID         *string    `json:"class_id,omitempty"`   // For students
	ClassName       *string    `json:"class_name,omitempty"` // For students
	IsActive        bool       `json:"is_active"`
	EmailVerified   bool       `json:"email_verified"`
	MFAEnabled      bool       `json:"mfa_enabled"`
	PhoneNumber     *string    `json:"phone_number,omitempty"`
	JobTitle        *string    `json:"job_title,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	PseudonymizedAt *time.Time `json:"pseudonymized_at,omitempty"`
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
