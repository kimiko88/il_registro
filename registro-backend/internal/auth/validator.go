package auth

import (
	"regexp"
	"unicode"
)

// Field length limits prevent oversized payloads from reaching the database
// or consuming excessive memory during hashing.
const (
	maxEmailLength    = 254 // RFC 5321 maximum
	maxNameLength     = 100
	maxPasswordLength = 128 // bcrypt silently truncates beyond 72 bytes; 128 is a safe practical limit
)

// Role constants — single source of truth used by validator, service and middleware.
const (
	RoleSuperAdmin    = "superadmin"
	RoleAdmin         = "admin"
	RolePrincipal     = "principal"
	RoleVicePrincipal = "vice_principal"
	RoleSecretary     = "secretary"
	RoleTeacher       = "teacher"
	RoleStudent       = "student"
	RoleParent        = "parent"
)

// allRoles is the exhaustive set of valid role strings.
var allRoles = map[string]bool{
	RoleSuperAdmin:    true,
	RoleAdmin:         true,
	RolePrincipal:     true,
	RoleVicePrincipal: true,
	RoleSecretary:     true,
	RoleTeacher:       true,
	RoleStudent:       true,
	RoleParent:        true,
}

// creatableRoles defines which roles each caller role is allowed to create.
var creatableRoles = map[string]map[string]bool{
	RoleSuperAdmin: {
		RoleSuperAdmin:    true,
		RoleAdmin:         true,
		RolePrincipal:     true,
		RoleVicePrincipal: true,
		RoleSecretary:     true,
		RoleTeacher:       true,
		RoleStudent:       true,
		RoleParent:        true,
	},
	RoleAdmin: {
		RoleAdmin:         true,
		RolePrincipal:     true,
		RoleVicePrincipal: true,
		RoleSecretary:     true,
		RoleTeacher:       true,
		RoleStudent:       true,
		RoleParent:        true,
	},
	RoleSecretary: {
		RolePrincipal:     true,
		RoleVicePrincipal: true,
		RoleTeacher:       true,
		RoleStudent:       true,
		RoleParent:        true,
	},
}

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

// ValidateRegisterRequest validates an authenticated registration request.
// callerRole is extracted from the JWT by the handler and must never come
// from the request body.
//
// Rules:
//   - Only superadmin, admin and segreteria can call this endpoint.
//   - Each caller can only create roles within their permission set
//     (see creatableRoles matrix above).
func ValidateRegisterRequest(callerRole string, req *RegisterRequest) error {
	// 1. Validate email and password first (cheap checks before DB hits)
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

	// 2. Requested role must be a known role string
	if !allRoles[req.Role] {
		return ErrInvalidRole
	}

	// 3. Caller must be allowed to register users at all
	allowed, callerCanRegister := creatableRoles[callerRole]
	if !callerCanRegister {
		return ErrInsufficientRole
	}

	// 4. Caller must be allowed to assign the requested role
	if !allowed[req.Role] {
		return ErrCannotCreateRole
	}

	return nil
}
