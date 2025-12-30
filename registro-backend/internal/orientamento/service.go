package orientamento

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	CreateEvent(ctx context.Context, teacherID string, req CreateEventRequest) error
	GetEvents(ctx context.Context) ([]Event, error)
	RegisterStudent(ctx context.Context, studentID, eventID string) error
	GetMyEvents(ctx context.Context, studentID string) ([]Participation, error)
	MarkAttendance(ctx context.Context, eventID, studentID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateEvent(ctx context.Context, teacherID string, req CreateEventRequest) error {
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return errors.New("invalid date format (expected RFC3339)")
	}
	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return errors.New("invalid end_date format (expected RFC3339)")
	}

	if end.Before(date) {
		return errors.New("end_date cannot be before start date")
	}

	e := &Event{
		SchoolID:     "default-school",
		Title:        req.Title,
		Description:  req.Description,
		Category:     req.Category,
		Date:         date,
		EndDate:      end,
		Location:     req.Location,
		Hours:        req.Hours,
		MaxAttendees: req.MaxAttendees,
		CreatedBy:    teacherID,
	}
	return s.repo.CreateEvent(ctx, e)
}

func (s *service) GetEvents(ctx context.Context) ([]Event, error) {
	return s.repo.GetEvents(ctx, "default-school")
}

func (s *service) RegisterStudent(ctx context.Context, studentID, eventID string) error {
	p := &Participation{EventID: eventID, StudentID: studentID}
	return s.repo.RegisterStudent(ctx, p)
}

func (s *service) GetMyEvents(ctx context.Context, studentID string) ([]Participation, error) {
	return s.repo.GetParticipations(ctx, studentID)
}

func (s *service) MarkAttendance(ctx context.Context, eventID, studentID string) error {
	return s.repo.MarkAttendance(ctx, eventID, studentID, true)
}
