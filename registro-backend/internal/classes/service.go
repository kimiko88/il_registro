package classes

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateClass(ctx context.Context, schoolID string, req CreateClassRequest) (*Class, error) {
	c := &Class{
		SchoolID:      schoolID,
		Name:          req.Name,
		Section:       req.Section,
		AcademicYear:  req.AcademicYear,
		CoordinatorID: req.CoordinatorID,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListClasses(ctx context.Context, schoolID string) ([]Class, error) {
	return s.repo.List(ctx, schoolID)
}

func (s *Service) GetClass(ctx context.Context, id string) (*Class, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) UpdateClass(ctx context.Context, id string, req CreateClassRequest) (*Class, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	c.Name = req.Name
	c.Section = req.Section
	c.AcademicYear = req.AcademicYear
	c.CoordinatorID = req.CoordinatorID

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteClass(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
