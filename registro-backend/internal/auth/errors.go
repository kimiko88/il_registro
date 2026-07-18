package auth

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenRevoked       = errors.New("token has been revoked")
	ErrUnauthorized       = errors.New("unauthorized")

	// Registration errors
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidEmail       = errors.New("invalid email format")

	// Password errors
	ErrPasswordTooShort    = errors.New("password must be at least 10 characters")
	ErrPasswordNoUppercase = errors.New("password must contain uppercase letter")
	ErrPasswordNoLowercase = errors.New("password must contain lowercase letter")
	ErrPasswordNoNumber    = errors.New("password must contain number")
	ErrPasswordNoSpecial   = errors.New("password must contain special character")
	ErrPasswordExpired     = errors.New("la password è scaduta e deve essere cambiata")

	// MFA errors
	ErrMFARequired       = errors.New("MFA token required")
	ErrInvalidMFAToken   = errors.New("invalid MFA token")
	ErrMFANotEnabled     = errors.New("MFA is not enabled for this user")
	ErrMFAAlreadyEnabled = errors.New("MFA is already enabled")

	// Rate limiting errors
	ErrTooManyAttempts = errors.New("too many login attempts, please try again later")

	// Generic errors
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrInvalidRole           = errors.New("invalid user role")
	ErrPasswordReused        = errors.New("la nuova password non può essere uguale ad una delle ultime 5 utilizzate")
)
