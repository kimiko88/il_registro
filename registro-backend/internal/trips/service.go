package trips

import (
	"context"
	"errors"
	"fmt"
	"time"

	"registro-backend/internal/users"
)

var (
	ErrUnauthorized = errors.New("unauthorized action on educational trip")
)

type Service struct {
	repo     Repository
	userRepo users.Repository
}

// NewService creates the trips service.
// Bug 131: userRepo is optional but strongly recommended for parent guardianship checks.
// When userRepo is nil, guardianship verification for parents is skipped.
func NewService(repo Repository, uRepo ...users.Repository) *Service {
	if repo == nil {
		panic("trips.NewService: repo must not be nil")
	}
	svc := &Service{repo: repo}
	if len(uRepo) > 0 && uRepo[0] != nil {
		svc.userRepo = uRepo[0]
	}
	return svc
}

// CreateTrip creates a new educational trip.
// Bug 132: validates that return date is not before departure date.
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

	// Bug 132: return date must be >= departure date
	if ret.Before(dep) {
		return nil, fmt.Errorf("return_date (%s) cannot be before departure_date (%s)",
			req.ReturnDate, req.DepartureDate)
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

func (s *Service) GetTripByID(ctx context.Context, id string) (*EducationalTrip, error) {
	return s.repo.GetTripByID(ctx, id)
}

// SubmitConsent records a parent or student consent for an educational trip.
// Bug 131: when the actor is a parent, verifies guardianship before recording consent.
func (s *Service) SubmitConsent(ctx context.Context, actorID, actorRole, ipAddress string, req SubmitConsentRequest) error {
	var parentID *string
	if actorRole == "parent" {
		// Bug 131: verify guardianship
		if s.userRepo != nil {
			isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, req.StudentID)
			if err != nil {
				return fmt.Errorf("errore verifica tutela: %w", err)
			}
			if !isGuardian {
				return errors.New("unauthorized: non sei tutore legale di questo studente")
			}
		}
		// NOTE: when userRepo is nil, guardianship check is skipped.
		// Callers must ensure userRepo is always provided in production.
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
