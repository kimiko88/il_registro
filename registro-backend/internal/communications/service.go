package communications

import (
	"context"
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

func (s *Service) DeleteMessage(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
