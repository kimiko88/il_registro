package auth

import (
	"regexp"
	"unicode"
)

// Field length limits prevent oversized payloads from reaching the database
// or consuming excessive memory during hashing.
const (
	maxEmailLength     = 254 // RFC 5321 maximum
	maxNameLength      = 100
	maxPasswordLength  = 128 // bcrypt silently truncates beyond 72 bytes; 128 is a safe practical limit
)

// PasswordValidator validates password strength
type PasswordValidator struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
}

// NewPasswordValidator creates a new password validator with default rules
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:      10,
		RequireUpper:   true,
		RequireLower:   true,
		RequireNumber:  true,
		RequireSpecial: true,
	}
}

// Validate checks if password meets requirements
func (v *PasswordValidator) Validate(password string) error {
	if len(password) < v.MinLength {
		return ErrPasswordTooShort
	}
	if len(password) > maxPasswordLength {
		return ErrPasswordTooLong
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if v.RequireUpper && !hasUpper {
		return ErrPasswordNoUppercase
	}
	if v.RequireLower && !hasLower {
		return ErrPasswordNoLowercase
	}
	if v.RequireNumber && !hasNumber {
		return ErrPasswordNoNumber
	}
	if v.RequireSpecial && !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}

// EmailValidator validates email format
type EmailValidator struct {
	emailRegex *regexp.Regexp
}

// NewEmailValidator creates a new email validator
func NewEmailValidator() *EmailValidator {
	return &EmailValidator{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

// Validate checks if email is valid and within allowed length.
func (v *EmailValidator) Validate(email string) error {
	if len(email) > maxEmailLength {
		return ErrInvalidEmail
	}
	if !v.emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// publicRegistrationRoles are the only roles allowed via the public /auth/register endpoint.
// Privileged roles (admin, superadmin) must be created by a superadmin via the admin API.
var publicRegistrationRoles = map[string]bool{
	"student": true,
	"teacher": true,
	"parent":  true,
}

// ValidateRegisterRequest validates a public registration request.
// Admin and superadmin roles are rejected here; they must be created
// through the dedicated admin user management endpoint.
func ValidateRegisterRequest(req *RegisterRequest) error {
	emailValidator := NewEmailValidator()
	if err := emailValidator.Validate(req.Email); err != nil {
		return err
	}

	passwordValidator := NewPasswordValidator()
	if err := passwordValidator.Validate(req.Password); err != nil {
		return err
	}

	if req.FirstName == "" || req.LastName == "" {
		return ErrMissingRequiredFields
	}
	if len(req.FirstName) > maxNameLength || len(req.LastName) > maxNameLength {
		return ErrFieldTooLong
	}

	if !publicRegistrationRoles[req.Role] {
		return ErrInvalidRole
	}

	return nil
}
