package student_goals

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("student_goals.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) CreateGoal(ctx context.Context, teacherID string, req CreateGoalRequest) (*StudentGoal, error) {
	g := &StudentGoal{
		StudentID:   req.StudentID,
		TeacherID:   teacherID,
		Title:       req.Title,
		Description: req.Description,
		BadgeName:   req.BadgeName,
		BadgeIcon:   req.BadgeIcon,
		Category:    req.Category,
		Points:      req.Points,
		Status:      StatusPending,
	}

	if req.DueDate != "" {
		if d, err := time.Parse("2006-01-02", req.DueDate); err == nil {
			g.DueDate = &d
		}
	}

	if err := s.repo.Create(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *Service) ListByStudent(ctx context.Context, studentID string) ([]*StudentGoal, error) {
	return s.repo.ListByStudent(ctx, studentID)
}

func (s *Service) UpdateGoalStatus(ctx context.Context, actorID, actorRole, id string, status GoalStatus) error {
	g, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if g.TeacherID != actorID && g.StudentID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return errors.New("unauthorized")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}
