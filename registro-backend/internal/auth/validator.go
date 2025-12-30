package auth

import (
	"regexp"
	"unicode"
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
		MinLength:      8,
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

// Validate checks if email is valid
func (v *EmailValidator) Validate(email string) error {
	if !v.emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateRegisterRequest validates registration request
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

	validRoles := map[string]bool{
		"student":    true,
		"teacher":    true,
		"parent":     true,
		"admin":      true,
		"superadmin": true,
	}

	if !validRoles[req.Role] {
		return ErrInvalidRole
	}

	return nil
}
