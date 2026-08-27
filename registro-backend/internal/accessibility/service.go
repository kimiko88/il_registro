package accessibility

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"strings"
	"time"
)

type Service interface {
	SubmitFeedback(ctx context.Context, req *CreateFeedbackRequest, userID, schoolID, userAgent, ipAddress *string) (*FeedbackResponse, error)
	ListFeedbacks(ctx context.Context, schoolID, status string, limit, offset int) ([]AccessibilityFeedback, int, error)
	UpdateFeedbackStatus(ctx context.Context, id, status, responseNotes string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SubmitFeedback(ctx context.Context, req *CreateFeedbackRequest, userID, schoolID, userAgent, ipAddress *string) (*FeedbackResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("il nome e cognome sono obbligatori")
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("indirizzo email non valido")
	}

	barrierType := strings.TrimSpace(req.BarrierType)
	if barrierType == "" {
		return nil, errors.New("la tipologia di barriera è obbligatoria")
	}

	desc := strings.TrimSpace(req.Description)
	if desc == "" || len(desc) < 5 {
		return nil, errors.New("la descrizione della problematica deve contenere almeno 5 caratteri")
	}

	// Generate AgID compliant protocol number: A11Y-YYYY-MMDD-XXXX
	now := time.Now()
	randNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
	protocolNum := fmt.Sprintf("A11Y-%d-%02d%02d-%04d", now.Year(), int(now.Month()), now.Day(), randNum.Int64()+1000)

	feedback := &AccessibilityFeedback{
		ProtocolNumber: protocolNum,
		Name:           name,
		Email:          email,
		BarrierType:    barrierType,
		Description:    desc,
		UserID:         userID,
		SchoolID:       schoolID,
		UserAgent:      userAgent,
		IPAddress:      ipAddress,
		Status:         "open",
	}

	if err := s.repo.Create(ctx, feedback); err != nil {
		return nil, fmt.Errorf("errore durante il salvataggio della segnalazione: %w", err)
	}

	return &FeedbackResponse{
		ID:             feedback.ID,
		ProtocolNumber: protocolNum,
		Message:        fmt.Sprintf("Segnalazione registrata con successo con protocollo %s", protocolNum),
	}, nil
}

func (s *service) ListFeedbacks(ctx context.Context, schoolID, status string, limit, offset int) ([]AccessibilityFeedback, int, error) {
	return s.repo.List(ctx, schoolID, status, limit, offset)
}

func (s *service) UpdateFeedbackStatus(ctx context.Context, id, status, responseNotes string) error {
	status = strings.ToLower(strings.TrimSpace(status))
	if status != "open" && status != "in_progress" && status != "resolved" {
		return errors.New("stato non valido (consentiti: open, in_progress, resolved)")
	}
	return s.repo.UpdateStatus(ctx, id, status, responseNotes)
}
