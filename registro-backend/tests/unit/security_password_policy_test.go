package unit

import (
	"errors"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func ValidatePasswordEntropyPolicy(password string, userEmail string, firstName string, lastName string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("password must contain uppercase, lowercase, digit, and special character")
	}

	lowerPass := strings.ToLower(password)
	if userEmail != "" {
		parts := strings.Split(userEmail, "@")
		if len(parts) > 0 && len(parts[0]) >= 3 && strings.Contains(lowerPass, strings.ToLower(parts[0])) {
			return errors.New("password cannot contain email address prefix")
		}
	}

	if firstName != "" && len(firstName) >= 3 && strings.Contains(lowerPass, strings.ToLower(firstName)) {
		return errors.New("password cannot contain first name")
	}

	if lastName != "" && len(lastName) >= 3 && strings.Contains(lowerPass, strings.ToLower(lastName)) {
		return errors.New("password cannot contain last name")
	}

	return nil
}

func TestSecurity_PasswordEntropyAndPersonalAttributePolicy(t *testing.T) {
	email := "mario.rossi@scuola.it"
	firstName := "Mario"
	lastName := "Rossi"

	t.Run("Validates strong password adhering to all entropy rules", func(t *testing.T) {
		err := ValidatePasswordEntropyPolicy("SecureP@ss2026!", email, firstName, lastName)
		assert.NoError(t, err)
	})

	t.Run("Rejects password shorter than 8 characters", func(t *testing.T) {
		err := ValidatePasswordEntropyPolicy("P@s1!", email, firstName, lastName)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "8 characters")
	})

	t.Run("Rejects password missing special character", func(t *testing.T) {
		err := ValidatePasswordEntropyPolicy("SecurePass2026", email, firstName, lastName)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "special character")
	})

	t.Run("Rejects password containing email prefix", func(t *testing.T) {
		err := ValidatePasswordEntropyPolicy("Mario.rossi@2026!", email, firstName, lastName)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email address prefix")
	})

	t.Run("Rejects password containing first name", func(t *testing.T) {
		err := ValidatePasswordEntropyPolicy("SuperMario#2026!", email, firstName, lastName)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first name")
	})
}
