package agenda

import (
	"context"
	"errors"
	"fmt"
	"time"

	"registro-backend/internal/users"
)

var (
	ErrNotFound     = errors.New("agenda item not found")
	ErrUnauthorized = errors.New("unauthorized action on agenda item")
	ErrInvalidDate  = errors.New("invalid date format: use YYYY-MM-DD")
)

type Service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, uRepo ...users.Repository) *Service {
	if repo == nil {
		panic("agenda.NewService: repo must not be nil")
	}
	var userRepo users.Repository
	if len(uRepo) > 0 {
		userRepo = uRepo[0]
	}
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) CreateAgendaItem(ctx context.Context, teacherID, schoolID string, req CreateAgendaItemRequest) (*AgendaItem, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school_id required")
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	item := &AgendaItem{
		SchoolID:    schoolID,
		ClassID:     req.ClassID,
		SubjectID:   req.SubjectID,
		TeacherID:   teacherID,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Date:        d,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetAgendaItem(ctx context.Context, id string) (*AgendaItem, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateAgendaItem(ctx context.Context, actorID, actorRole, actorSchoolID, id string, req UpdateAgendaItemRequest) (*AgendaItem, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if actorRole != "superadmin" {
		if actorSchoolID == "" || item.SchoolID != actorSchoolID {
			return nil, ErrUnauthorized
		}
		if actorRole != "admin" && item.TeacherID != actorID {
			return nil, ErrUnauthorized
		}
	}

	if req.Title != nil {
		item.Title = *req.Title
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Type != nil {
		item.Type = *req.Type
	}
	if req.Date != nil {
		d, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, ErrInvalidDate
		}
		item.Date = d
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) DeleteAgendaItem(ctx context.Context, actorID, actorRole, actorSchoolID, id string) error {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if actorRole != "superadmin" {
		if actorSchoolID == "" || item.SchoolID != actorSchoolID {
			return ErrUnauthorized
		}
		if actorRole != "admin" && item.TeacherID != actorID {
			return ErrUnauthorized
		}
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetCalendar(ctx context.Context, schoolID, userID, role string, filter CalendarFilter) ([]*AgendaItem, error) {
	if filter.From.IsZero() {
		now := time.Now()
		filter.From = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	}
	if filter.To.IsZero() {
		filter.To = filter.From.AddDate(0, 3, 0)
	}

	if role == "student" {
		if filter.StudentID == "" {
			filter.StudentID = userID
		}
	} else if role == "parent" {
		if filter.StudentID != "" && s.userRepo != nil {
			isGuardian, err := s.userRepo.IsGuardian(ctx, userID, filter.StudentID)
			if err != nil || !isGuardian {
				return nil, ErrUnauthorized
			}
		}
	}
	return s.repo.ListCalendar(ctx, schoolID, filter)
}

func (s *Service) SetTaskCompletion(ctx context.Context, studentID, role, itemID string, completed bool) error {
	if role != "student" {
		return errors.New("unauthorized: task completion can only be updated by a student")
	}
	if studentID == "" {
		return errors.New("unauthorized: studentID required")
	}
	item, err := s.repo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}
	// Bug 108: TODO — verificare che studentID sia iscritto a item.ClassID.
	// Richiede StudentBelongsToClass(ctx, studentID, item.ClassID) nel Repository.
	// Al momento il check avviene lato DB tramite la JOIN su student_classes,
	// che il SetCompletion handler potrebbe rafforzare in futuro.
	_ = item // usato per l'esistenza dell'item
	return s.repo.SetCompletion(ctx, itemID, studentID, completed)
}
