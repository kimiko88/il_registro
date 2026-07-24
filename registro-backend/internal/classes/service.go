package classes

import (
	"context"
	"fmt"
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
		Articolazione: req.Articolazione,
		AcademicYear:  req.AcademicYear,
		CoordinatorID: req.CoordinatorID,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListClasses(ctx context.Context, schoolID string, academicYear string) ([]Class, error) {
	return s.repo.List(ctx, schoolID, academicYear)
}

func (s *Service) GetTeacherClasses(ctx context.Context, teacherID string) ([]Class, error) {
	return s.repo.ListByTeacher(ctx, teacherID)
}

func (s *Service) GetClass(ctx context.Context, id string) (*Class, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) UpdateClass(ctx context.Context, schoolID, id string, req CreateClassRequest) (*Class, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.SchoolID != schoolID {
		return nil, fmt.Errorf("forbidden: class belongs to another school")
	}

	c.Name = req.Name
	c.Section = req.Section
	c.Articolazione = req.Articolazione
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

func (s *Service) AssignSubject(ctx context.Context, classID string, req AssignSubjectRequest) error {
	return s.repo.AssignSubject(ctx, classID, req.SubjectID, req.TeacherID, req.HoursPerWeek)
}

func (s *Service) RemoveSubject(ctx context.Context, assignmentID string) error {
	return s.repo.UnassignSubject(ctx, assignmentID)
}

func (s *Service) GetClassSubjects(ctx context.Context, classID string) ([]ClassSubject, error) {
	return s.repo.GetClassSubjects(ctx, classID)
}

func (s *Service) GetClassGuardians(ctx context.Context, classID string) ([]GuardianInfo, error) {
	return s.repo.GetClassGuardians(ctx, classID)
}
