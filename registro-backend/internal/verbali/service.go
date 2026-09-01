package verbali

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized action on verbali")
	ErrNotFound     = errors.New("verbale or meeting not found")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("verbali.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

// CreateMeeting creates a new council meeting.
// Bug 140: verifying ClassID belongs to schoolID requires ClassBelongsToSchool in Repository.
// TODO: add ClassBelongsToSchool(ctx, classID, schoolID) to Repository and call it here.
func (s *Service) CreateMeeting(ctx context.Context, actorID, schoolID string, req CreateMeetingRequest) (*CouncilMeeting, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school_id required")
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: use YYYY-MM-DD")
	}

	m := &CouncilMeeting{
		SchoolID:  schoolID,
		ClassID:   req.ClassID,
		Title:     req.Title,
		Date:      d,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Agenda:    req.Agenda,
		CreatedBy: actorID,
	}

	if err := s.repo.CreateMeeting(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) ListMeetings(ctx context.Context, schoolID, classID string) ([]*CouncilMeeting, error) {
	return s.repo.ListMeetings(ctx, schoolID, classID)
}

// CreateVerbale creates a new verbale for a council meeting.
// Bug 140: verifying ClassID school membership would require ClassBelongsToSchool (see CreateMeeting).
func (s *Service) CreateVerbale(ctx context.Context, actorID, actorRole string, req CreateVerbaleRequest) (*MeetingVerbale, error) {
	if actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin" {
		return nil, ErrUnauthorized
	}
	_, err := s.repo.GetMeetingByID(ctx, req.MeetingID)
	if err != nil {
		return nil, ErrNotFound
	}

	v := &MeetingVerbale{
		MeetingID:   req.MeetingID,
		Title:       req.Title,
		Content:     req.Content,
		SecretaryID: req.SecretaryID,
		PresidentID: req.PresidentID,
		IsPublished: req.IsPublished,
	}

	if err := s.repo.CreateVerbale(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) GetVerbale(ctx context.Context, verbaleID, userID string) (*MeetingVerbale, error) {
	return s.repo.GetVerbaleByID(ctx, verbaleID, userID)
}

func (s *Service) ListVerbali(ctx context.Context, meetingID, userID string) ([]*MeetingVerbale, error) {
	return s.repo.ListVerbali(ctx, meetingID, userID)
}

// SignVerbale allows a user to sign a published verbale.
// Bug 139: verifies that the signer is the secretary or president of the verbale.
// NOTA: verifica partecipanti al consiglio richiede un campo ParticipantIDs nel modello/repo.
// Al momento verifichiamo solo SecretaryID e PresidentID che sono campi espliciti del verbale.
func (s *Service) SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error {
	v, err := s.repo.GetVerbaleByID(ctx, verbaleID, userID)
	if err != nil {
		return ErrNotFound
	}
	if !v.IsPublished {
		return errors.New("cannot sign an unpublished verbale")
	}

	// Bug 139: verify the signer is secretary or president
	isAuthorized := (v.SecretaryID != nil && *v.SecretaryID == userID) ||
		(v.PresidentID != nil && *v.PresidentID == userID)
	// NOTE: for full participant list verification, add ParticipantIDs to MeetingVerbale
	// and verify: for _, pid := range v.ParticipantIDs { if pid == userID { isAuthorized = true } }
	if !isAuthorized {
		return errors.New("unauthorized: non sei il segretario o il presidente di questo verbale")
	}

	return s.repo.SignVerbale(ctx, verbaleID, userID, ipAddress)
}

func (s *Service) GetSignatures(ctx context.Context, verbaleID string) ([]VerbaleSignature, error) {
	return s.repo.GetSignatures(ctx, verbaleID)
}

// GetMeeting retrieves a council meeting by ID (used for PDF generation).
func (s *Service) GetMeeting(ctx context.Context, meetingID string) (*CouncilMeeting, error) {
	return s.repo.GetMeetingByID(ctx, meetingID)
}
