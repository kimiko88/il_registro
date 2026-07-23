package trips

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateTrip(ctx context.Context, t *EducationalTrip) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *MockRepository) GetTripByID(ctx context.Context, id string) (*EducationalTrip, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EducationalTrip), args.Error(1)
}
func (m *MockRepository) ListTrips(ctx context.Context, schoolID, studentID string) ([]*EducationalTrip, error) {
	args := m.Called(ctx, schoolID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*EducationalTrip), args.Error(1)
}
func (m *MockRepository) SubmitConsent(ctx context.Context, c *TripConsent) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockRepository) ListConsents(ctx context.Context, tripID string) ([]*TripConsent, error) {
	args := m.Called(ctx, tripID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*TripConsent), args.Error(1)
}

func TestCreateAndSubmitTripConsent(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	reqTrip := CreateTripRequest{
		Title:         "Gita a Firenze",
		Destination:   "Firenze, Italia",
		DepartureDate: "2025-12-01",
		ReturnDate:    "2025-12-03",
		ClassIDs:      []string{"class-1"},
	}

	mockRepo.On("CreateTrip", mock.Anything, mock.MatchedBy(func(tr *EducationalTrip) bool {
		return tr.Title == "Gita a Firenze"
	})).Return(nil).Once()

	tr, err := svc.CreateTrip(context.Background(), "t-1", "school-1", reqTrip)
	assert.NoError(t, err)
	assert.NotNil(t, tr)

	mockRepo.On("SubmitConsent", mock.Anything, mock.MatchedBy(func(c *TripConsent) bool {
		return c.TripID == "trip-1" && c.Status == "granted"
	})).Return(nil).Once()

	err = svc.SubmitConsent(context.Background(), "parent-1", "parent", "192.168.1.100", SubmitConsentRequest{
		TripID:    "trip-1",
		StudentID: "student-1",
		Status:    "granted",
	})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
