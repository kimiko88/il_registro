package users

import (
	"regexp"
	"strings"
)

// Validator handles custom validation rules
type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

// ValidatePassword checks password complexity (min 10 chars, upper, lower, digit, special char)
func (v *Validator) ValidatePassword(password string) bool {
	if len(password) < 10 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>/?~`", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// ValidateFiscalCode checks Italian Codice Fiscale format
func (v *Validator) ValidateFiscalCode(fc string) bool {
	// Simple regex check for now: 6 letters, 2 digits, 1 letter, 2 digits, 1 letter, 3 digits, 1 letter
	// Example: RSSMRA80A01H501U
	fc = strings.ToUpper(fc)
	regex := regexp.MustCompile(`^[A-Z]{6}\d{2}[A-Z]\d{2}[A-Z]\d{3}[A-Z]$`)
	return regex.MatchString(fc)
}

func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func SanitizeText(text string) string {
	return strings.TrimSpace(text)
}

func SanitizeTextPtr(text string) *string {
	s := strings.TrimSpace(text)
	if s == "" {
		return nil
	}
	return &s
}
