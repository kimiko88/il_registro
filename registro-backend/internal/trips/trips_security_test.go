package trips_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/trips"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTripsRepo struct {
	mock.Mock
}

func (m *mockTripsRepo) CreateTrip(ctx context.Context, t *trips.EducationalTrip) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *mockTripsRepo) GetTripByID(ctx context.Context, id string) (*trips.EducationalTrip, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*trips.EducationalTrip), args.Error(1)
}
func (m *mockTripsRepo) ListTrips(ctx context.Context, schoolID, studentID string) ([]*trips.EducationalTrip, error) {
	args := m.Called(ctx, schoolID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*trips.EducationalTrip), args.Error(1)
}
func (m *mockTripsRepo) SubmitConsent(ctx context.Context, c *trips.TripConsent) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *mockTripsRepo) ListConsents(ctx context.Context, tripID string) ([]*trips.TripConsent, error) {
	args := m.Called(ctx, tripID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*trips.TripConsent), args.Error(1)
}

type mockUserRepoForTrips struct {
	mock.Mock
	users.Repository
}

func (m *mockUserRepoForTrips) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func buildTripsEngine(repo trips.Repository, userRepo users.Repository, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	if userID != "" {
		r.Use(func(c *gin.Context) {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
			c.Next()
		})
	}

	svc := trips.NewService(repo, userRepo)
	h := trips.NewHandler(svc)
	h.RegisterRoutes(r.Group("/api"))

	return r
}

func TestTrips_SubmitConsent_StudentCannotImpersonate(t *testing.T) {
	mockRepo := new(mockTripsRepo)
	// Even if student sends student_id="victim-student", the consent must be saved with student_id="student-me"
	mockRepo.On("SubmitConsent", mock.Anything, mock.MatchedBy(func(c *trips.TripConsent) bool {
		return c.StudentID == "student-me" && c.TripID == "trip-1"
	})).Return(nil).Once()

	r := buildTripsEngine(mockRepo, nil, "student-me", "student", "school-1")

	body := trips.SubmitConsentRequest{
		TripID:    "trip-1",
		StudentID: "victim-student",
		Status:    "Consented",
	}
	raw, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/trips/consent", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestTrips_ListTrips_ParentGuardianship(t *testing.T) {
	t.Run("Parent querying non-guardian student is blocked with 403", func(t *testing.T) {
		mockRepo := new(mockTripsRepo)
		mockUser := new(mockUserRepoForTrips)
		mockUser.On("IsGuardian", mock.Anything, "parent-1", "other-student").Return(false, nil).Once()

		r := buildTripsEngine(mockRepo, mockUser, "parent-1", "parent", "school-1")

		req := httptest.NewRequest(http.MethodGet, "/api/trips?student_id=other-student", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockUser.AssertExpectations(t)
	})

	t.Run("Parent querying own child is allowed", func(t *testing.T) {
		mockRepo := new(mockTripsRepo)
		mockUser := new(mockUserRepoForTrips)

		mockUser.On("IsGuardian", mock.Anything, "parent-1", "my-child").Return(true, nil).Once()
		mockRepo.On("ListTrips", mock.Anything, "school-1", "my-child").Return([]*trips.EducationalTrip{}, nil).Once()

		r := buildTripsEngine(mockRepo, mockUser, "parent-1", "parent", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/trips?student_id=my-child", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
		mockUser.AssertExpectations(t)
	})
}
