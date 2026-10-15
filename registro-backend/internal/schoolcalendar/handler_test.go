package schoolcalendar

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCalendarService struct {
	mock.Mock
}

func (m *MockCalendarService) SetSchoolYear(ctx context.Context, actorID, actorRole, schoolID string, req SetYearRequest) (*SchoolYearResponse, error) {
	args := m.Called(ctx, actorID, actorRole, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SchoolYearResponse), args.Error(1)
}

func (m *MockCalendarService) GetSchoolYear(ctx context.Context, schoolID string) (*SchoolYearResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SchoolYearResponse), args.Error(1)
}

func (m *MockCalendarService) AddNonTeachingDay(ctx context.Context, actorID, actorRole, schoolID string, req AddNonTeachingDayRequest) (*NonTeachingDayResponse, error) {
	args := m.Called(ctx, actorID, actorRole, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*NonTeachingDayResponse), args.Error(1)
}

func (m *MockCalendarService) DeleteNonTeachingDay(ctx context.Context, actorRole, schoolID, id string) error {
	args := m.Called(ctx, actorRole, schoolID, id)
	return args.Error(0)
}

func (m *MockCalendarService) ListNonTeachingDays(ctx context.Context, schoolID string) ([]NonTeachingDayResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]NonTeachingDayResponse), args.Error(1)
}

func (m *MockCalendarService) CountTeachingDays(ctx context.Context, schoolID string, from, to time.Time) (int, error) {
	args := m.Called(ctx, schoolID, from, to)
	return args.Int(0), args.Error(1)
}

func (m *MockCalendarService) GetSchoolYearDates(ctx context.Context, schoolID string) (start, end time.Time, err error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).(time.Time), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockCalendarService) CreateAcademicPeriod(ctx context.Context, actorRole, schoolID string, req CreateAcademicPeriodRequest) (*AcademicPeriod, error) {
	args := m.Called(ctx, actorRole, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AcademicPeriod), args.Error(1)
}

func (m *MockCalendarService) ListAcademicPeriods(ctx context.Context, schoolID string) ([]AcademicPeriod, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]AcademicPeriod), args.Error(1)
}

func (m *MockCalendarService) GetCalendarEvents(ctx context.Context, schoolID string, year string, eventType string) ([]CalendarEvent, error) {
	args := m.Called(ctx, schoolID, year, eventType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]CalendarEvent), args.Error(1)
}

func (m *MockCalendarService) GetStudentCalendarEvents(ctx context.Context, schoolID string, studentID string, year string, month string) ([]CalendarEvent, error) {
	args := m.Called(ctx, schoolID, studentID, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]CalendarEvent), args.Error(1)
}

func setupCalendarRouter(svc Service, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestSchoolCalendarHandler_SetYearAndGetYear(t *testing.T) {
	svc := new(MockCalendarService)
	rSec := setupCalendarRouter(svc, "sec-1", "secretary", "school-1")

	// 1. SetYear
	reqBody := SetYearRequest{
		YearLabel: "2025/2026",
		StartDate: "2025-09-15",
		EndDate:   "2026-06-10",
	}
	expectedResp := &SchoolYearResponse{
		ID:        "sy-1",
		YearLabel: "2025/2026",
		StartDate: time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
	}
	svc.On("SetSchoolYear", mock.Anything, "sec-1", "secretary", "school-1", reqBody).Return(expectedResp, nil).Once()

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/school-calendar/year", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rSec.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. GetYear
	svc.On("GetSchoolYear", mock.Anything, "school-1").Return(expectedResp, nil).Once()
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/school-calendar/year", nil)
	rSec.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSchoolCalendarHandler_AddNonTeachingDay(t *testing.T) {
	svc := new(MockCalendarService)
	rSec := setupCalendarRouter(svc, "sec-1", "secretary", "school-1")

	reqBody := AddNonTeachingDayRequest{
		Date:  "2025-11-01",
		Label: "Ognissanti",
	}
	expectedResp := &NonTeachingDayResponse{
		ID:    "ntd-1",
		Date:  "2025-11-01",
		Label: "Ognissanti",
	}
	svc.On("AddNonTeachingDay", mock.Anything, "sec-1", "secretary", "school-1", reqBody).Return(expectedResp, nil).Once()

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/school-calendar/non-teaching-days", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rSec.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSchoolCalendarHandler_GetCalendarEvents(t *testing.T) {
	svc := new(MockCalendarService)
	rStudent := setupCalendarRouter(svc, "stud-1", "student", "school-1")

	events := []CalendarEvent{
		{
			ID:    "ev-1",
			Title: "Inizio lezioni",
			Date:  "2025-09-15",
			Type:  "event",
		},
	}
	svc.On("GetCalendarEvents", mock.Anything, "school-1", "2025/2026", "").Return(events, nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/school-calendar?year=2025/2026", nil)
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp []CalendarEvent
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Inizio lezioni", resp[0].Title)
}
