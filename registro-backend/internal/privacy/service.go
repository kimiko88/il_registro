package privacy

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrConsentNotFound = errors.New("consenso privacy non trovato")
)

// EvaluateTrafficLight assigns the privacy badge according to granted parental consents.
func EvaluateTrafficLight(photoVideo, cloudWorkspace, walkingTrips bool) TrafficLightBadge {
	if photoVideo {
		return TrafficLightGreen
	}
	if cloudWorkspace || walkingTrips {
		return TrafficLightYellow
	}
	return TrafficLightRed
}

type Repository interface {
	SaveConsent(ctx context.Context, c *StudentConsent) error
	GetConsent(ctx context.Context, studentID, schoolYear string) (*StudentConsent, error)
	ListConsentsBySchool(ctx context.Context, schoolID, schoolYear string) ([]*StudentConsent, error)
	SaveTreatment(ctx context.Context, t *PrivacyTreatment) error
	ListTreatments(ctx context.Context, schoolID string) ([]*PrivacyTreatment, error)
}

type MemoryRepository struct {
	mu         sync.RWMutex
	consents   map[string]*StudentConsent     // key: studentID:year
	treatments map[string][]*PrivacyTreatment // key: schoolID
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		consents:   make(map[string]*StudentConsent),
		treatments: make(map[string][]*PrivacyTreatment),
	}
}

func (m *MemoryRepository) SaveConsent(_ context.Context, c *StudentConsent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s:%s", c.StudentID, c.SchoolYear)
	m.consents[key] = c
	return nil
}

func (m *MemoryRepository) GetConsent(_ context.Context, studentID, schoolYear string) (*StudentConsent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := fmt.Sprintf("%s:%s", studentID, schoolYear)
	c, ok := m.consents[key]
	if !ok {
		return nil, ErrConsentNotFound
	}
	return c, nil
}

func (m *MemoryRepository) ListConsentsBySchool(_ context.Context, schoolID, schoolYear string) ([]*StudentConsent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*StudentConsent
	for _, c := range m.consents {
		if c.SchoolID == schoolID && (schoolYear == "" || c.SchoolYear == schoolYear) {
			list = append(list, c)
		}
	}
	return list, nil
}

func (m *MemoryRepository) SaveTreatment(_ context.Context, t *PrivacyTreatment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.treatments[t.SchoolID] = append(m.treatments[t.SchoolID], t)
	return nil
}

func (m *MemoryRepository) ListTreatments(_ context.Context, schoolID string) ([]*PrivacyTreatment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.treatments[schoolID], nil
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		repo = NewMemoryRepository()
	}
	return &Service{repo: repo}
}

func (s *Service) SaveStudentConsent(ctx context.Context, schoolID, signedBy string, req SaveConsentRequest) (*StudentConsent, error) {
	if schoolID == "" || req.StudentID == "" {
		return nil, errors.New("school_id e student_id sono obbligatori")
	}

	badge := EvaluateTrafficLight(req.PhotoVideoSocialConsent, req.CloudWorkspaceConsent, req.WalkingTripsConsent)
	now := time.Now()

	consent := &StudentConsent{
		ID:                      fmt.Sprintf("cons-%s-%s", req.StudentID, req.SchoolYear),
		SchoolID:                schoolID,
		StudentID:               req.StudentID,
		SchoolYear:              req.SchoolYear,
		PhotoVideoSocialConsent: req.PhotoVideoSocialConsent,
		CloudWorkspaceConsent:   req.CloudWorkspaceConsent,
		WalkingTripsConsent:     req.WalkingTripsConsent,
		TrafficLightBadge:       badge,
		SignedBy:                signedBy,
		SignedAt:                &now,
		Notes:                   req.Notes,
		UpdatedAt:               now,
	}

	if err := s.repo.SaveConsent(ctx, consent); err != nil {
		return nil, err
	}
	return consent, nil
}

func (s *Service) GetStudentConsent(ctx context.Context, studentID, schoolYear string) (*StudentConsent, error) {
	c, err := s.repo.GetConsent(ctx, studentID, schoolYear)
	if err != nil {
		// Default to Red if unrecorded
		return &StudentConsent{
			StudentID:         studentID,
			SchoolYear:        schoolYear,
			TrafficLightBadge: TrafficLightRed,
		}, nil
	}
	return c, nil
}

func (s *Service) GetClassBadges(ctx context.Context, studentIDs []string, schoolYear string) (map[string]TrafficLightBadge, error) {
	result := make(map[string]TrafficLightBadge)
	for _, id := range studentIDs {
		c, err := s.GetStudentConsent(ctx, id, schoolYear)
		if err == nil && c != nil {
			result[id] = c.TrafficLightBadge
		} else {
			result[id] = TrafficLightRed
		}
	}
	return result, nil
}

func (s *Service) CreateTreatment(ctx context.Context, schoolID string, req TreatmentRequest) (*PrivacyTreatment, error) {
	t := &PrivacyTreatment{
		ID:                   fmt.Sprintf("treat-%d", time.Now().UnixNano()),
		SchoolID:             schoolID,
		ActivityName:         req.ActivityName,
		LegalBasis:           req.LegalBasis,
		InterestedCategories: req.InterestedCategories,
		RetentionPeriod:      req.RetentionPeriod,
		DPOContact:           req.DPOContact,
		SecurityMeasures:     req.SecurityMeasures,
		CreatedAt:            time.Now(),
	}
	if err := s.repo.SaveTreatment(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTreatments(ctx context.Context, schoolID string) ([]*PrivacyTreatment, error) {
	return s.repo.ListTreatments(ctx, schoolID)
}
