package unit

import (
	"html"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Helper sanitization matching backend security functions
func sanitizeCSVField(val string) string {
	if val == "" {
		return val
	}
	trimmed := strings.TrimSpace(val)
	if strings.HasPrefix(trimmed, "=") ||
		strings.HasPrefix(trimmed, "+") ||
		strings.HasPrefix(trimmed, "-") ||
		strings.HasPrefix(trimmed, "@") ||
		strings.HasPrefix(trimmed, "\t") ||
		strings.HasPrefix(trimmed, "\r") {
		return "'" + val
	}
	return val
}

func sanitizeHTMLInput(val string) string {
	return html.EscapeString(val)
}

func validateJWTToken(tokenStr string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

func TestSecurity_CSVFormulaInjectionDefense(t *testing.T) {
	t.Run("Sanitizes leading formula characters in CSV fields", func(t *testing.T) {
		assert.Equal(t, "'=SUM(A1:A10)", sanitizeCSVField("=SUM(A1:A10)"))
		assert.Equal(t, "'+cmd|' /C calc'!A0", sanitizeCSVField("+cmd|' /C calc'!A0"))
		assert.Equal(t, "'-100", sanitizeCSVField("-100"))
		assert.Equal(t, "'@SUM(B1)", sanitizeCSVField("@SUM(B1)"))
	})

	t.Run("Leaves safe strings untouched in CSV fields", func(t *testing.T) {
		assert.Equal(t, "Mario Rossi", sanitizeCSVField("Mario Rossi"))
		assert.Equal(t, "mario.rossi@scuola.it", sanitizeCSVField("mario.rossi@scuola.it"))
		assert.Equal(t, "Matematica", sanitizeCSVField("Matematica"))
	})
}

func TestSecurity_XSSAndHTMLSanitization(t *testing.T) {
	t.Run("Escapes malicious HTML script tags in user inputs", func(t *testing.T) {
		input := "<script>alert('xss')</script>"
		sanitized := sanitizeHTMLInput(input)
		assert.NotContains(t, sanitized, "<script>")
		assert.Contains(t, sanitized, "&lt;script&gt;")
	})

	t.Run("Escapes img onerror XSS payloads", func(t *testing.T) {
		input := "<img src=x onerror=alert(1)>"
		sanitized := sanitizeHTMLInput(input)
		assert.NotContains(t, sanitized, "<img")
		assert.Contains(t, sanitized, "&lt;img")
	})
}

func TestSecurity_JWTValidationIntegrity(t *testing.T) {
	secret := []byte("ci-test-secret-32-chars-long!!")

	t.Run("Validates valid signed JWT token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "user-123",
			"role":    "teacher",
			"exp":     time.Now().Add(1 * time.Hour).Unix(),
		})
		tokenStr, err := token.SignedString(secret)
		assert.NoError(t, err)

		claims, err := validateJWTToken(tokenStr, secret)
		assert.NoError(t, err)
		assert.Equal(t, "user-123", claims["user_id"])
		assert.Equal(t, "teacher", claims["role"])
	})

	t.Run("Rejects expired JWT token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "user-123",
			"exp":     time.Now().Add(-1 * time.Hour).Unix(),
		})
		tokenStr, err := token.SignedString(secret)
		assert.NoError(t, err)

		_, err = validateJWTToken(tokenStr, secret)
		assert.Error(t, err)
	})

	t.Run("Rejects JWT signed with wrong secret", func(t *testing.T) {
		wrongSecret := []byte("wrong-secret-32-chars-long!!")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "user-123",
			"exp":     time.Now().Add(1 * time.Hour).Unix(),
		})
		tokenStr, err := token.SignedString(wrongSecret)
		assert.NoError(t, err)

		_, err = validateJWTToken(tokenStr, secret)
		assert.Error(t, err)
	})
}
