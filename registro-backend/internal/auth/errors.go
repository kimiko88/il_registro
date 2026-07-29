package auth

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrPasswordTooShort      = errors.New("password must be at least 10 characters")
	ErrPasswordTooLong       = errors.New("password must be at most 128 characters")
	ErrPasswordNoUppercase   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber      = errors.New("password must contain at least one number")
	ErrPasswordNoSpecial     = errors.New("password must contain at least one special character")
	ErrPasswordExpired       = errors.New("password has expired, please reset it")
	ErrPasswordReused        = errors.New("cannot reuse a recent password")
	// ErrInvalidRole is returned when the requested role is not a valid role string.
	ErrInvalidRole           = errors.New("invalid role")
	// ErrInsufficientRole is returned when the caller's role does not have
	// permission to create users at all (e.g. a teacher trying to register someone).
	ErrInsufficientRole      = errors.New("only superadmin, admin and secretary can register new users")
	// ErrCannotCreateRole is returned when the caller is allowed to register users
	// but is trying to assign a role that exceeds their own privileges.
	ErrCannotCreateRole      = errors.New("you do not have permission to create a user with this role")
	ErrMissingRequiredFields = errors.New("first name and last name are required")
	ErrFieldTooLong          = errors.New("first name and last name must not exceed 100 characters")
	ErrUserInactive          = errors.New("account is inactive")
	ErrTooManyAttempts       = errors.New("too many login attempts, please try again later")
	ErrMFARequired           = errors.New("MFA token required")
	ErrInvalidMFAToken       = errors.New("invalid MFA token")
	ErrMFAAlreadyEnabled     = errors.New("MFA is already enabled")
	ErrMFANotEnabled         = errors.New("MFA is not enabled")
	ErrInvalidToken          = errors.New("invalid or expired token")
	ErrTokenRevoked          = errors.New("token has been revoked")
	ErrUnauthorized          = errors.New("unauthorized")
)
