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
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Type == "" {
		return nil, errors.New("type is required")
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	if req.StartTime == "" {
		req.StartTime = "09:00"
	}
	if req.EndTime == "" {
		req.EndTime = "10:00"
	}

	if !req.AllDay {
		if _, err := time.Parse("15:04", req.StartTime); err != nil {
			return nil, errors.New("invalid start_time format: expected HH:MM")
		}
		if _, err := time.Parse("15:04", req.EndTime); err != nil {
			return nil, errors.New("invalid end_time format: expected HH:MM")
		}
	}

	item := &AgendaItem{
		SchoolID:    schoolID,
		ClassID:     req.ClassID,
		SubjectID:   req.SubjectID,
		TeacherID:   teacherID,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		AllDay:      req.AllDay,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Date:        d,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetAgendaItem(ctx context.Context, actorRole, actorSchoolID, id string) (*AgendaItem, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if actorRole != "superadmin" && actorSchoolID != "" && item.SchoolID != actorSchoolID {
		return nil, ErrUnauthorized
	}
	return item, nil
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
	if req.AllDay != nil {
		item.AllDay = *req.AllDay
	}
	if req.StartTime != nil {
		if *req.StartTime != "" {
			if _, err := time.Parse("15:04", *req.StartTime); err != nil {
				return nil, errors.New("invalid start_time format: expected HH:MM")
			}
		}
		item.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		if *req.EndTime != "" {
			if _, err := time.Parse("15:04", *req.EndTime); err != nil {
				return nil, errors.New("invalid end_time format: expected HH:MM")
			}
		}
		item.EndTime = *req.EndTime
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

	switch role {
	case "student":
		if filter.StudentID == "" {
			filter.StudentID = userID
		}
	case "parent":
		if filter.StudentID == "" {
			return nil, errors.New("student_id is required for parent role")
		}
		if s.userRepo == nil {
			return nil, ErrUnauthorized
		}
		isGuardian, err := s.userRepo.IsGuardian(ctx, userID, filter.StudentID)
		if err != nil || !isGuardian {
			return nil, ErrUnauthorized
		}
	}
	return s.repo.ListCalendar(ctx, schoolID, filter)
}

func (s *Service) SetTaskCompletion(ctx context.Context, actorID, role, targetStudentID, itemID string, completed bool) error {
	// Only students can mark their own tasks; parents can mark on behalf of their children
	// after guardianship validation.
	switch role {
	case "student":
		if actorID == "" {
			return errors.New("unauthorized: studentID required")
		}
		// Student marks their own record.
		targetStudentID = actorID
	case "parent":
		if targetStudentID == "" {
			return errors.New("unauthorized: studentID required for parent role")
		}
		if s.userRepo == nil {
			return ErrUnauthorized
		}
		isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, targetStudentID)
		if err != nil {
			return fmt.Errorf("failed to verify guardianship: %w", err)
		}
		if !isGuardian {
			return ErrUnauthorized
		}
	default:
		return errors.New("unauthorized: task completion can only be updated by a student or parent")
	}

	item, err := s.repo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item.ClassID != "" {
		isMember, err := s.repo.IsStudentInClass(ctx, targetStudentID, item.ClassID)
		if err != nil || !isMember {
			return errors.New("forbidden: lo studente non appartiene alla classe dell'agenda item")
		}
	}
	return s.repo.SetCompletion(ctx, itemID, targetStudentID, completed)
}
