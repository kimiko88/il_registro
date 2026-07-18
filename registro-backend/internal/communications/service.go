package communications

import (
	"context"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SendMessage(ctx context.Context, senderID string, req CreateMessageRequest) (*Message, error) {
	msg := &Message{
		SenderID:    senderID,
		ReceiverIDs: req.Recipients,
		Subject:     req.Subject,
		Body:        req.Body,
		Type:        req.Type,
	}
	if err := s.repo.Create(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *Service) ListMessages(ctx context.Context, userID string) ([]*Message, error) {
	return s.repo.List(ctx, userID)
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

func (s *Service) GetMessageSignatures(ctx context.Context, communicationID string) ([]string, error) {
	return s.repo.GetSignatures(ctx, communicationID)
}
