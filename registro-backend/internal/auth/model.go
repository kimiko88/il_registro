package auth

import (
	"time"
)

// User represents the authenticated user profile
type User struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"` // Never expose in JSON
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Role              string     `json:"role"`
	SchoolID          *string    `json:"school_id,omitempty"`
	IsActive          bool       `json:"is_active"`
	IsStaff           bool       `json:"is_staff"`
	EmailVerified     bool       `json:"email_verified"`
	MFAEnabled        bool       `json:"mfa_enabled"`
	MFASecret         string     `json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	LastLogin         *time.Time `json:"last_login,omitempty"`
	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty"`
}

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `json:"revoked"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

// LoginAttempt tracks failed login attempts for rate limiting
type LoginAttempt struct {
	Email       string    `json:"email"`
	IPAddress   string    `json:"ip_address"`
	Success     bool      `json:"success"`
	AttemptedAt time.Time `json:"attempted_at"`
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// MFARecoveryCode represents a backup code for MFA
type MFARecoveryCode struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Code      string     `json:"code"`
	Used      bool       `json:"used"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
}
