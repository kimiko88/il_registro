package textbooks

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTextbook(ctx context.Context, schoolID string, req CreateTextbookRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}
	if req.Price < 0 {
		return errors.New("price cannot be negative")
	}
	t := &Textbook{
		SchoolID:  schoolID,
		Title:     strings.TrimSpace(req.Title),
		Author:    strings.TrimSpace(req.Author),
		Subject:   strings.TrimSpace(req.Subject),
		ISBN:      strings.TrimSpace(req.ISBN),
		Publisher: strings.TrimSpace(req.Publisher),
		Price:     req.Price,
	}
	return s.repo.Create(ctx, t)
}

func (s *Service) UpdateTextbook(ctx context.Context, id string, req CreateTextbookRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}
	if req.Price < 0 {
		return errors.New("price cannot be negative")
	}
	t := &Textbook{
		ID:        id,
		Title:     strings.TrimSpace(req.Title),
		Author:    strings.TrimSpace(req.Author),
		Subject:   strings.TrimSpace(req.Subject),
		ISBN:      strings.TrimSpace(req.ISBN),
		Publisher: strings.TrimSpace(req.Publisher),
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
