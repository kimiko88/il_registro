package communications

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("communications.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) SendMessage(ctx context.Context, senderID string, req CreateMessageRequest) (*Message, error) {
	msg := &Message{
		SchoolID:          req.SchoolID,
		SenderID:          senderID,
		ReceiverIDs:       req.Recipients,
		Subject:           req.Subject,
		Body:              req.Body,
		AttachmentURL:     req.AttachmentURL,
		Type:              req.Type,
		RequiresSignature: req.RequiresSignature,
	}

	if req.SignatureDeadline != nil && *req.SignatureDeadline != "" {
		if d, err := time.Parse("2006-01-02", *req.SignatureDeadline); err == nil {
			msg.SignatureDeadline = &d
		} else if d, err := time.Parse(time.RFC3339, *req.SignatureDeadline); err == nil {
			msg.SignatureDeadline = &d
		}
	}

	if err := s.repo.Create(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *Service) ListMessages(ctx context.Context, userID string) ([]*Message, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) ListBacheca(ctx context.Context, schoolID, userID string) ([]*Message, error) {
	return s.repo.ListBacheca(ctx, schoolID, userID)
}

func (s *Service) DeleteMessage(ctx context.Context, actorID string, actorRole string, id string) error {
	msg, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if msg.SenderID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return errors.New("unauthorized: cannot delete message of another user")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) SignMessage(ctx context.Context, communicationID string, userID string) error {
	return s.repo.Sign(ctx, communicationID, userID)
}

func (s *Service) SignMessageWithIP(ctx context.Context, communicationID string, userID string, ipAddress string) error {
	return s.repo.SignWithIP(ctx, communicationID, userID, ipAddress)
}

func (s *Service) GetMessageSignatures(ctx context.Context, communicationID string) ([]string, error) {
	return s.repo.GetSignatures(ctx, communicationID)
}

func (s *Service) GetSignatureReport(ctx context.Context, actorRole, communicationID string) (*SignatureReportResponse, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "teacher" {
		return nil, errors.New("forbidden: signature reports are restricted to staff")
	}
	return s.repo.GetSignatureReport(ctx, communicationID)
}

func (s *Service) GetMessageByID(ctx context.Context, userID, role, id string) (*Message, error) {
	msg, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == "admin" || role == "superadmin" || msg.Type == "bacheca" || msg.SenderID == userID {
		return msg, nil
	}
	for _, r := range msg.ReceiverIDs {
		if r == userID {
			return msg, nil
		}
	}
	return nil, errors.New("unauthorized: cannot access communication")
}

func (s *Service) UpdateMessage(ctx context.Context, actorID, actorRole, id, subject, body string) error {
	msg, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if msg.SenderID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return errors.New("unauthorized: cannot edit message of another user")
	}
	return s.repo.Update(ctx, id, subject, body)
}

func (s *Service) MarkAsRead(ctx context.Context, communicationID, userID, ipAddress string) error {
	return s.repo.MarkAsRead(ctx, communicationID, userID, ipAddress)
}

func (s *Service) GetUnreadUsers(ctx context.Context, communicationID string) ([]string, error) {
	return s.repo.GetUnreadUsers(ctx, communicationID)
}

func (s *Service) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}
