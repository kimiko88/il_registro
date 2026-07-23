package trips

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized action on educational trip")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("trips.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) CreateTrip(ctx context.Context, teacherID, schoolID string, req CreateTripRequest) (*EducationalTrip, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school_id required")
	}

	dep, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		dep, err = time.Parse(time.RFC3339, req.DepartureDate)
		if err != nil {
			return nil, fmt.Errorf("invalid departure_date format")
		}
	}

	ret, err := time.Parse("2006-01-02", req.ReturnDate)
	if err != nil {
		ret, err = time.Parse(time.RFC3339, req.ReturnDate)
		if err != nil {
			return nil, fmt.Errorf("invalid return_date format")
		}
	}

	t := &EducationalTrip{
		SchoolID:             schoolID,
		Title:                req.Title,
		Destination:          req.Destination,
		DepartureDate:        dep,
		ReturnDate:           ret,
		Description:          req.Description,
		AccompanyingTeachers: req.AccompanyingTeachers,
		ClassIDs:             req.ClassIDs,
	}

	if err := s.repo.CreateTrip(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTrips(ctx context.Context, schoolID, studentID string) ([]*EducationalTrip, error) {
	return s.repo.ListTrips(ctx, schoolID, studentID)
}

func (s *Service) SubmitConsent(ctx context.Context, actorID, actorRole, ipAddress string, req SubmitConsentRequest) error {
	var parentID *string
	if actorRole == "parent" {
		parentID = &actorID
	}

	c := &TripConsent{
		TripID:    req.TripID,
		StudentID: req.StudentID,
		ParentID:  parentID,
		Status:    req.Status,
		IPAddress: ipAddress,
	}

	return s.repo.SubmitConsent(ctx, c)
}

func (s *Service) ListConsents(ctx context.Context, tripID string) ([]*TripConsent, error) {
	return s.repo.ListConsents(ctx, tripID)
}
