package general_meetings

import (
	"context"
	"fmt"
	"time"
)

type Service interface {
	CreateMeeting(ctx context.Context, actorID, schoolID string, req CreateGeneralMeetingRequest) (*GeneralMeeting, error)
	ListMeetings(ctx context.Context, schoolID, userID, role string) ([]*GeneralMeeting, error)
	GetMeeting(ctx context.Context, id, userID string) (*GeneralMeeting, error)
	DeleteMeeting(ctx context.Context, id string) error

	RegisterUser(ctx context.Context, meetingID, userID string) error
	UnregisterUser(ctx context.Context, meetingID, userID string) error
	ListRegistrations(ctx context.Context, meetingID string) ([]*GeneralMeetingRegistration, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateMeeting(ctx context.Context, actorID, schoolID string, req CreateGeneralMeetingRequest) (*GeneralMeeting, error) {
	meetingDate, err := time.Parse(time.RFC3339, req.MeetingDate)
	if err != nil {
		// Fallback date-only format
		meetingDate, err = time.Parse("2006-01-02", req.MeetingDate)
		if err != nil {
			return nil, fmt.Errorf("invalid meeting_date format: %w", err)
		}
	}

	var regDeadline *time.Time
	if req.RegistrationDeadline != nil && *req.RegistrationDeadline != "" {
		t, err := time.Parse(time.RFC3339, *req.RegistrationDeadline)
		if err != nil {
			t, err = time.Parse("2006-01-02", *req.RegistrationDeadline)
		}
		if err == nil {
			regDeadline = &t
		}
	}

	m := &GeneralMeeting{
		SchoolID:             schoolID,
		Title:                req.Title,
		Description:          req.Description,
		Location:             req.Location,
		MeetingDate:          meetingDate,
		RegistrationDeadline: regDeadline,
		MaxParticipants:      req.MaxParticipants,
		IsMandatory:          req.IsMandatory,
		TargetRoles:          req.TargetRoles,
		CreatedBy:            actorID,
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) ListMeetings(ctx context.Context, schoolID, userID, role string) ([]*GeneralMeeting, error) {
	return s.repo.List(ctx, schoolID, userID, role)
}

func (s *service) GetMeeting(ctx context.Context, id, userID string) (*GeneralMeeting, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *service) DeleteMeeting(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) RegisterUser(ctx context.Context, meetingID, userID string) error {
	m, err := s.repo.GetByID(ctx, meetingID, userID)
	if err != nil {
		return err
	}
	if m.RegistrationDeadline != nil && time.Now().After(*m.RegistrationDeadline) {
		return fmt.Errorf("registration deadline has passed")
	}
	if m.MaxParticipants != nil && m.RegistrationsCount >= *m.MaxParticipants {
		return fmt.Errorf("meeting is fully booked")
	}
	return s.repo.RegisterUser(ctx, meetingID, userID)
}

func (s *service) UnregisterUser(ctx context.Context, meetingID, userID string) error {
	return s.repo.UnregisterUser(ctx, meetingID, userID)
}

func (s *service) ListRegistrations(ctx context.Context, meetingID string) ([]*GeneralMeetingRegistration, error) {
	return s.repo.ListRegistrations(ctx, meetingID)
}
