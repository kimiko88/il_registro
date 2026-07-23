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
}
