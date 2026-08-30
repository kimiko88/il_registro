package student_goals

import (
	"context"
	"errors"
	"fmt"
	"time"

	"registro-backend/internal/users"
)

var (
	ErrUnauthorizedStudent = errors.New("unauthorized: cannot view goals of another student")
	ErrNotGuardian        = errors.New("unauthorized: non sei tutore legale di questo studente")
	ErrUnauthorized       = errors.New("unauthorized")
)

type Service struct {
	repo     Repository
	userRepo users.Repository
}


func NewService(repo Repository, uRepo ...users.Repository) *Service {
	if repo == nil {
		panic("student_goals.NewService: repo must not be nil")
	}
	svc := &Service{repo: repo}
	if len(uRepo) > 0 && uRepo[0] != nil {
		svc.userRepo = uRepo[0]
	}
	return svc
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

func (s *Service) ListByStudent(ctx context.Context, actorID, actorRole, studentID string) ([]*StudentGoal, error) {
	if actorRole == "student" && actorID != studentID {
		return nil, ErrUnauthorizedStudent
	}
	if actorRole == "parent" && s.userRepo != nil {
		isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, studentID)
		if err != nil {
			return nil, fmt.Errorf("errore verifica tutela: %w", err)
		}
		if !isGuardian {
			return nil, ErrNotGuardian
		}
	}
	return s.repo.ListByStudent(ctx, studentID)
}

func (s *Service) UpdateGoalStatus(ctx context.Context, actorID, actorRole, id string, status GoalStatus) error {
	g, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if g.TeacherID != actorID && g.StudentID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

