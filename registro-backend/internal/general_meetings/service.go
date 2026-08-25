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
	DeleteMeeting(ctx context.Context, id, actorID, actorRole, actorSchoolID string) error

	RegisterUser(ctx context.Context, meetingID, userID string) error
	UnregisterUser(ctx context.Context, meetingID, userID string) error
	ListRegistrations(ctx context.Context, meetingID, actorID, actorRole, actorSchoolID string) ([]*GeneralMeetingRegistration, error)
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

// DeleteMeeting deletes a meeting.
// Bug 143: restricted to the meeting creator or admin/superadmin of the same school.
func (s *service) DeleteMeeting(ctx context.Context, id, actorID, actorRole, actorSchoolID string) error {
	m, err := s.repo.GetByID(ctx, id, actorID)
	if err != nil {
		return err
	}
	if actorRole == "admin" && actorSchoolID != "" && m.SchoolID != actorSchoolID {
		return fmt.Errorf("unauthorized: non puoi eliminare una riunione di un'altra scuola")
	}
	if m.CreatedBy != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return fmt.Errorf("unauthorized: solo il creatore o un amministratore può eliminare questa riunione")
	}
	return s.repo.Delete(ctx, id)
}

// RegisterUser registers a user for a meeting.
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

func (s *service) ListRegistrations(ctx context.Context, meetingID, actorID, actorRole, actorSchoolID string) ([]*GeneralMeetingRegistration, error) {
	m, err := s.repo.GetByID(ctx, meetingID, actorID)
	if err != nil {
		return nil, err
	}
	if actorRole != "superadmin" && actorSchoolID != "" && m.SchoolID != actorSchoolID {
		return nil, fmt.Errorf("forbidden: riunione appartenente ad un'altra scuola")
	}
	return s.repo.ListRegistrations(ctx, meetingID)
}
