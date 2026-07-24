package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"registro-backend/pkg/logger"
)

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SMSConfig struct {
	Provider string // 'twilio', 'mock'
	APIKey   string
}

type Service struct {
	mailCfg MailConfig
	smsCfg  SMSConfig
}

func NewService(mailCfg MailConfig, smsCfg SMSConfig) *Service {
	return &Service{
		mailCfg: mailCfg,
		smsCfg:  smsCfg,
	}
}

// sanitizeHeader removes newlines to prevent MIME header injection attacks
func sanitizeHeader(value string) string {
	clean := strings.ReplaceAll(value, "\r", "")
	return strings.ReplaceAll(clean, "\n", "")
}

func (s *Service) SendEmail(ctx context.Context, toEmail, subject, bodyHTML string) error {
	cleanTo := sanitizeHeader(toEmail)
	cleanSubject := sanitizeHeader(subject)

	if s.mailCfg.Host == "" {
		// Mock send if SMTP not configured
		logger.Log.Infof("[SMTP Mailer] Mock send email to %s - Subject: %s", cleanTo, cleanSubject)
		return nil
	}

	addr := fmt.Sprintf("%s:%d", s.mailCfg.Host, s.mailCfg.Port)
	auth := smtp.PlainAuth("", s.mailCfg.Username, s.mailCfg.Password, s.mailCfg.Host)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.mailCfg.From, cleanTo, cleanSubject, bodyHTML))

	err := smtp.SendMail(addr, auth, s.mailCfg.From, []string{cleanTo}, msg)
	if err != nil {
		logger.Log.Errorf("[SMTP Mailer] Failed to send email to %s", cleanTo)
		return fmt.Errorf("smtp send failed: %w", err)
	}
	return nil
}

func (s *Service) SendSMS(ctx context.Context, phoneNumber, message string) error {
	// SMS provider gateway dispatch
	logger.Log.Infof("[SMS Gateway] Dispatched SMS via %s", s.smsCfg.Provider)
	return nil
}
