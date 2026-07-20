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

func (s *Service) SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error {
	v, err := s.repo.GetVerbaleByID(ctx, verbaleID, userID)
	if err != nil {
		return ErrNotFound
	}
	if !v.IsPublished {
		return errors.New("cannot sign an unpublished verbale")
	}
	return s.repo.SignVerbale(ctx, verbaleID, userID, ipAddress)
}

func (s *Service) GetSignatures(ctx context.Context, verbaleID string) ([]VerbaleSignature, error) {
	return s.repo.GetSignatures(ctx, verbaleID)
}
