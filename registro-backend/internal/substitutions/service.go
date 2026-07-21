package substitutions

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("substitutions.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) CreateSubstitution(ctx context.Context, schoolID string, req CreateSubstitutionRequest) (*Substitution, error) {
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	sub := &Substitution{
		SchoolID:            schoolID,
		ClassID:             req.ClassID,
		AbsentTeacherID:     req.AbsentTeacherID,
		SubstituteTeacherID: req.SubstituteTeacherID,
		Date:                d,
		Hour:                req.Hour,
		SubjectID:           req.SubjectID,
		Notes:               req.Notes,
		Status:              StatusPending,
	}

	if req.SubstituteTeacherID != nil && *req.SubstituteTeacherID != "" {
		sub.Status = StatusAssigned
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error) {
	return s.repo.ListBySchool(ctx, schoolID, date)
}

func (s *Service) ListByTeacher(ctx context.Context, teacherID string) ([]*Substitution, error) {
	return s.repo.ListByTeacher(ctx, teacherID)
}

func (s *Service) AssignSubstitute(ctx context.Context, id string, req AssignSubstituteRequest) error {
	return s.repo.AssignSubstitute(ctx, id, req.SubstituteTeacherID, req.Notes)
}
