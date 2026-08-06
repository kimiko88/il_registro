package timetables

import (
	"context"
	"errors"
	"fmt"
)

type Service interface {
	GetByClass(ctx context.Context, actorID, actorRole, schoolID, classID string) ([]ClassSchedule, error)
	GetMySchedule(ctx context.Context, actorID, actorRole string) ([]ClassSchedule, error)
	Update(ctx context.Context, actorID, actorRole, schoolID, classID string, entries []ScheduleEntry) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	if repo == nil {
		panic("timetables.NewService: repo must not be nil")
	}
	return &service{repo: repo}
}

func (s *service) GetByClass(ctx context.Context, actorID, actorRole, schoolID, classID string) ([]ClassSchedule, error) {
	if actorID == "" || actorRole == "" {
		return nil, errors.New("unauthorized: missing actor context")
	}
	if classID == "" {
		return nil, errors.New("class_id is required")
	}

	if actorRole == "student" || actorRole == "parent" {
		studentClassID, err := s.repo.GetStudentClassID(ctx, actorID)
		if err != nil {
			return nil, fmt.Errorf("impossibile verificare la classe dell'utente: %w", err)
		}
		if studentClassID != classID {
			return nil, errors.New("forbidden: non puoi visualizzare l'orario di un'altra classe")
		}
	}

	return s.repo.GetByClass(ctx, classID)
}

func (s *service) GetMySchedule(ctx context.Context, actorID, actorRole string) ([]ClassSchedule, error) {
	if actorID == "" {
		return nil, errors.New("unauthorized")
	}

	switch actorRole {
	case "student":
		classID, err := s.repo.GetStudentClassID(ctx, actorID)
		if err != nil {
			return []ClassSchedule{}, nil
		}
		return s.repo.GetByClass(ctx, classID)
	case "parent":
		classID, err := s.repo.GetParentStudentClassID(ctx, actorID)
		if err != nil {
			return []ClassSchedule{}, nil
		}
		return s.repo.GetByClass(ctx, classID)
	case "teacher":
		return s.repo.GetTeacherSchedule(ctx, actorID)
	default:
		return []ClassSchedule{}, nil
	}
}

func (s *service) Update(ctx context.Context, actorID, actorRole, schoolID, classID string, entries []ScheduleEntry) error {
	if actorID == "" {
		return errors.New("unauthorized")
	}
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" {
		return errors.New("forbidden: only administrative staff can update class schedules")
	}
	if classID == "" {
		return errors.New("class_id is required")
	}

	return s.repo.Update(ctx, classID, entries)
}
