package orientamento

import (
	"context"
	"errors"
	"time"
)

// Service defines the orientamento domain service interface.
// Bug 145: CreateEvent and GetEvents now accept schoolID to replace the hardcoded "default-school" literal.
type Service interface {
	CreateEvent(ctx context.Context, teacherID, schoolID string, req CreateEventRequest) error
	GetEvents(ctx context.Context, schoolID string) ([]Event, error)
	RegisterStudent(ctx context.Context, studentID, eventID string) error
	GetMyEvents(ctx context.Context, studentID string) ([]Participation, error)
	MarkAttendance(ctx context.Context, eventID, studentID string) error
	SavePreference(ctx context.Context, studentID string, pref StudentPreference) error
	GetPreference(ctx context.Context, studentID string) (*StudentPreference, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateEvent creates a new orientamento event.
// Bug 145: schoolID is now a parameter instead of the hardcoded literal "default-school".
func (s *service) CreateEvent(ctx context.Context, teacherID, schoolID string, req CreateEventRequest) error {
	if schoolID == "" {
		return errors.New("school_id required")
	}
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
		SchoolID:     schoolID,
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

// GetEvents returns events for a specific school, or all events if schoolID is empty.
// Bug 145: schoolID is now a parameter instead of the hardcoded literal "default-school".
func (s *service) GetEvents(ctx context.Context, schoolID string) ([]Event, error) {
	return s.repo.GetEvents(ctx, schoolID)
}

// RegisterStudent registers a student for an orientamento event.
// Bug 146: checks for duplicate enrollment before inserting.
func (s *service) RegisterStudent(ctx context.Context, studentID, eventID string) error {
	// Bug 146: check for existing participation to prevent duplicate enrollment
	existing, err := s.repo.GetParticipations(ctx, studentID)
	if err == nil {
		for _, p := range existing {
			if p.EventID == eventID {
				return errors.New("studente già iscritto a questo evento di orientamento")
			}
		}
	}
	p := &Participation{EventID: eventID, StudentID: studentID}
	return s.repo.RegisterStudent(ctx, p)
}

func (s *service) GetMyEvents(ctx context.Context, studentID string) ([]Participation, error) {
	return s.repo.GetParticipations(ctx, studentID)
}

func (s *service) MarkAttendance(ctx context.Context, eventID, studentID string) error {
	return s.repo.MarkAttendance(ctx, eventID, studentID, true)
}

func (s *service) SavePreference(ctx context.Context, studentID string, pref StudentPreference) error {
	pref.StudentID = studentID
	return s.repo.SavePreference(ctx, &pref)
}

func (s *service) GetPreference(ctx context.Context, studentID string) (*StudentPreference, error) {
	return s.repo.GetPreference(ctx, studentID)
}
