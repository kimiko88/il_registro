package mailer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMailerService(t *testing.T) {
	svc := NewService(MailConfig{}, SMSConfig{Provider: "mock"})

	err := svc.SendEmail(context.Background(), "parent@school.it", "Nuova Comunicazione", "<h1>Circolare pubblicata</h1>")
	assert.NoError(t, err)

	err = svc.SendSMS(context.Background(), "+393331234567", "Assenza non giustificata per Rossi Mario")
	assert.NoError(t, err)

	err = svc.SendPasswordReset(context.Background(), "parent@school.it", "secret-token-12345")
	assert.NoError(t, err)
}

func TestSanitizeHeader_CRLFInjection(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"user@example.com\r\nBcc: evil@hacker.com", "user@example.comBcc: evil@hacker.com"},
		{"Subject Line\nSecond line", "Subject LineSecond line"},
		{"Subject Line\rSecond line", "Subject LineSecond line"},
		{"Safe subject line", "Safe subject line"},
	}

	for _, tc := range cases {
		cleaned := sanitizeHeader(tc.input)
		assert.Equal(t, tc.expected, cleaned)
		assert.NotContains(t, cleaned, "\r")
		assert.NotContains(t, cleaned, "\n")
	}
}
