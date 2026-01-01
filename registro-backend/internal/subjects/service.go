package subjects

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSubject(ctx context.Context, schoolID string, req CreateSubjectRequest) (*Subject, error) {
	subj := &Subject{
		SchoolID:    schoolID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}
	if req.IsMandatory != nil {
		subj.IsMandatory = *req.IsMandatory
	} else {
		subj.IsMandatory = true // Default
	}

	if err := s.repo.Create(ctx, subj); err != nil {
		return nil, err
	}
	return subj, nil
}

func (s *Service) ListSubjects(ctx context.Context, schoolID string) ([]Subject, error) {
	return s.repo.List(ctx, schoolID)
}

func (s *Service) GetSubject(ctx context.Context, id string) (*Subject, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) UpdateSubject(ctx context.Context, id string, req CreateSubjectRequest) (*Subject, error) {
	subj, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	subj.Name = req.Name
	subj.Code = req.Code
	subj.Description = req.Description
	if req.IsMandatory != nil {
		subj.IsMandatory = *req.IsMandatory
	}

	if err := s.repo.Update(ctx, subj); err != nil {
		return nil, err
	}
	return subj, nil
}

func (s *Service) DeleteSubject(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
