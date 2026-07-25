package teachers

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

func (s *Service) ListTeachers(ctx context.Context, schoolID string) ([]Teacher, error) {
	return s.repo.List(ctx, schoolID)
}

func (s *Service) GetTeacher(ctx context.Context, schoolID, id string) (*Teacher, error) {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if schoolID != "" && t.SchoolID != schoolID {
		return nil, errors.New("forbidden: teacher belongs to another school")
	}
	return t, nil
}

func (s *Service) GetTeacherSubjects(ctx context.Context, teacherID string) ([]TeacherSubject, error) {
	return s.repo.GetSubjects(ctx, teacherID)
}

func (s *Service) GetTeachersBySubject(ctx context.Context, subjectID string) ([]Teacher, error) {
	return s.repo.GetBySubject(ctx, subjectID)
}

func (s *Service) AssignSubject(ctx context.Context, actorRole, schoolID, teacherID, subjectID string) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" {
		return errors.New("unauthorized: missing permissions to assign subjects")
	}
	if schoolID != "" {
		t, err := s.repo.Get(ctx, teacherID)
		if err != nil {
			return err
		}
		if t.SchoolID != schoolID {
			return errors.New("forbidden: teacher belongs to another school")
		}
	}
	return s.repo.AssignSubject(ctx, teacherID, subjectID)
}

func (s *Service) RemoveSubject(ctx context.Context, actorRole, schoolID, teacherID, subjectID string) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" {
		return errors.New("unauthorized: missing permissions to remove subjects")
	}
	if schoolID != "" {
		t, err := s.repo.Get(ctx, teacherID)
		if err != nil {
			return err
		}
		if t.SchoolID != schoolID {
			return errors.New("forbidden: teacher belongs to another school")
		}
	}
	return s.repo.RemoveSubject(ctx, teacherID, subjectID)
}

func (s *Service) GetDashboardStats(ctx context.Context, authenticatedUserID, teacherUserID, actorRole string) (map[string]interface{}, error) {
	if authenticatedUserID != teacherUserID && actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" {
		return nil, errors.New("unauthorized: cannot access another teacher's dashboard")
	}
	return s.repo.GetDashboardStats(ctx, teacherUserID)
}
