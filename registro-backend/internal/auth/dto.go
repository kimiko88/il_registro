package auth

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=10"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Role      string `json:"role" validate:"required,oneof=student teacher parent admin"`
	SchoolID  string `json:"school_id,omitempty"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	MFAToken string `json:"mfa_token,omitempty"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetConfirmRequest represents password reset confirmation
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=10"`
}

// MFASetupRequest represents MFA setup initiation
type MFASetupRequest struct {
	UserID string `json:"user_id"`
}

// MFAVerifyRequest represents MFA token verification
type MFAVerifyRequest struct {
	Token string `json:"token" validate:"required,len=6"`
}

// AuthResponse represents authentication response with tokens
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Role          string  `json:"role"`
	SchoolID      *string `json:"school_id,omitempty"`
	EmailVerified bool    `json:"email_verified"`
	MFAEnabled    bool    `json:"mfa_enabled"`
}

// MFASetupResponse represents MFA setup response with QR code
type MFASetupResponse struct {
	Secret        string   `json:"secret"`
	QRCodeURL     string   `json:"qr_code_url"`
	RecoveryCodes []string `json:"recovery_codes"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// MessageResponse represents a generic message response
type MessageResponse struct {
	Message string `json:"message"`
}
