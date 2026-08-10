package textbooks

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTextbook(ctx context.Context, schoolID string, req CreateTextbookRequest) error {
	t := &Textbook{
		SchoolID:  schoolID,
		Title:     req.Title,
		Author:    req.Author,
		Subject:   req.Subject,
		ISBN:      req.ISBN,
		Publisher: req.Publisher,
		Price:     req.Price,
	}
	return s.repo.Create(ctx, t)
}

func (s *Service) UpdateTextbook(ctx context.Context, id string, req CreateTextbookRequest) error {
	t := &Textbook{
		ID:        id,
		Title:     req.Title,
		Author:    req.Author,
		Subject:   req.Subject,
		ISBN:      req.ISBN,
		Publisher: req.Publisher,
		Price:     req.Price,
	}
	return s.repo.Update(ctx, t)
}

func (s *Service) ListTextbooks(ctx context.Context, schoolID string) ([]Textbook, error) {
	return s.repo.List(ctx, schoolID)
}

func (s *Service) DeleteTextbook(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) AssignToClass(ctx context.Context, classID string, req AssignTextbookRequest) error {
	return s.repo.AssignToClass(ctx, classID, req.SubjectID, req.TextbookID, req.IsOptional)
}

func (s *Service) RemoveFromClass(ctx context.Context, assignmentID string) error {
	return s.repo.RemoveFromClass(ctx, assignmentID)
}

func (s *Service) ListByClass(ctx context.Context, classID string) ([]ClassTextbook, error) {
	return s.repo.ListByClass(ctx, classID)
}
