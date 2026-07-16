package teachers

import (
	"context"
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

func (s *Service) GetTeacher(ctx context.Context, id string) (*Teacher, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetTeacherSubjects(ctx context.Context, teacherID string) ([]TeacherSubject, error) {
	return s.repo.GetSubjects(ctx, teacherID)
}

func (s *Service) GetTeachersBySubject(ctx context.Context, subjectID string) ([]Teacher, error) {
	return s.repo.GetBySubject(ctx, subjectID)
}

func (s *Service) AssignSubject(ctx context.Context, teacherID, subjectID string) error {
	return s.repo.AssignSubject(ctx, teacherID, subjectID)
}

func (s *Service) RemoveSubject(ctx context.Context, teacherID, subjectID string) error {
	return s.repo.RemoveSubject(ctx, teacherID, subjectID)
}

func (s *Service) GetDashboardStats(ctx context.Context, teacherUserID string) (map[string]interface{}, error) {
	return s.repo.GetDashboardStats(ctx, teacherUserID)
}
