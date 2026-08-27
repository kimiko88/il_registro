package colloqui

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockServiceForRangeTest struct {
	mock.Mock
}

func (m *mockServiceForRangeTest) CreateSlot(ctx context.Context, teacherUserID, schoolID string, req CreateSlotRequest) (*ColloquioSlot, error) {
	args := m.Called(ctx, teacherUserID, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioSlot), args.Error(1)
}

func (m *mockServiceForRangeTest) ListSlots(ctx context.Context, schoolID, teacherID string, from, to time.Time, availableOnly bool) ([]*ColloquioSlot, error) {
	args := m.Called(ctx, schoolID, teacherID, from, to, availableOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ColloquioSlot), args.Error(1)
}

func (m *mockServiceForRangeTest) CancelSlot(ctx context.Context, actorID, actorRole, actorSchoolID, slotID string) error {
	args := m.Called(ctx, actorID, actorRole, actorSchoolID, slotID)
	return args.Error(0)
}

func (m *mockServiceForRangeTest) BookSlot(ctx context.Context, parentUserID, parentSchoolID string, req CreateBookingRequest) (*ColloquioBooking, error) {
	args := m.Called(ctx, parentUserID, parentSchoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioBooking), args.Error(1)
}

func (m *mockServiceForRangeTest) ListMyBookings(ctx context.Context, userID, schoolID string) ([]*ColloquioBooking, error) {
	args := m.Called(ctx, userID, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ColloquioBooking), args.Error(1)
}

func (m *mockServiceForRangeTest) ListSlotBookings(ctx context.Context, actorID, actorRole, slotID string) ([]*ColloquioBooking, error) {
	args := m.Called(ctx, actorID, actorRole, slotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ColloquioBooking), args.Error(1)
}

func (m *mockServiceForRangeTest) UpdateBookingStatus(ctx context.Context, actorID, actorRole, bookingID string, req UpdateBookingStatusRequest) error {
	args := m.Called(ctx, actorID, actorRole, bookingID, req)
	return args.Error(0)
}

func (m *mockServiceForRangeTest) PatchSlot(ctx context.Context, actorID, actorRole, actorSchoolID, slotID, startTime, endTime string) error {
	args := m.Called(ctx, actorID, actorRole, actorSchoolID, slotID, startTime, endTime)
	return args.Error(0)
}

func (m *mockServiceForRangeTest) GetBookingByID(ctx context.Context, actorID, actorRole, bookingID string) (*ColloquioBooking, error) {
	args := m.Called(ctx, actorID, actorRole, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioBooking), args.Error(1)
}

func (m *mockServiceForRangeTest) CreateAssembly(ctx context.Context, actorRole, teacherID, schoolID string, req CreateAssemblyRequest) (*ColloquioSlot, error) {
	args := m.Called(ctx, actorRole, teacherID, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioSlot), args.Error(1)
}

func (m *mockServiceForRangeTest) CreateGeneralMeeting(ctx context.Context, schoolID string, req CreateGeneralParentMeetingRequest) (*GeneralParentMeeting, error) {
	return nil, nil
}
func (m *mockServiceForRangeTest) ListGeneralMeetings(ctx context.Context, schoolID string) ([]GeneralParentMeeting, error) {
	return nil, nil
}
func (m *mockServiceForRangeTest) GetGeneralMeeting(ctx context.Context, id string) (*GeneralParentMeeting, error) {
	return nil, nil
}
func (m *mockServiceForRangeTest) BookQueueTicket(ctx context.Context, parentID string, req BookQueueTicketRequest) (*GeneralMeetingQueueTicket, error) {
	return nil, nil
}
func (m *mockServiceForRangeTest) ListQueueTickets(ctx context.Context, meetingID, teacherID, parentID string) ([]GeneralMeetingQueueTicket, error) {
	return nil, nil
}
func (m *mockServiceForRangeTest) UpdateTicketStatus(ctx context.Context, id, status, notes string) error {
	return nil
}

func TestGetAvailabilityByTeacher_RangeLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForRangeTest)
	h := NewHandler(mockSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "parent-1")
		c.Set("role", "parent")
		c.Set("school_id", "school-1")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	// Request with a 2-year range (> 366 days) should return 400 Bad Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/colloqui/availability/teacher-1?from=2024-01-01&to=2026-06-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "range di date troppo ampio")
}

func TestGetBookingByID_NotFound_And_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForRangeTest)
	h := NewHandler(mockSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "parent-1")
		c.Set("role", "parent")
		c.Set("school_id", "school-1")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	// 1. Not found sentinel
	mockSvc.On("GetBookingByID", mock.Anything, "parent-1", "parent", "booking-404").Return(nil, ErrBookingNotFound).Once()
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/colloqui/bookings/booking-404", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusNotFound, w1.Code)

	// 2. Generic internal error
	mockSvc.On("GetBookingByID", mock.Anything, "parent-1", "parent", "booking-500").Return(nil, errors.New("db error")).Once()
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/colloqui/bookings/booking-500", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}
