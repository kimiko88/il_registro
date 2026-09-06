package orientamento

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

type MockOrientamentoService struct {
	mock.Mock
}

func (m *MockOrientamentoService) CreateEvent(ctx context.Context, teacherID, schoolID string, req CreateEventRequest) error {
	args := m.Called(ctx, teacherID, schoolID, req)
	return args.Error(0)
}
func (m *MockOrientamentoService) GetEvents(ctx context.Context, schoolID string) ([]Event, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Event), args.Error(1)
}
func (m *MockOrientamentoService) RegisterStudent(ctx context.Context, studentID, eventID string) error {
	args := m.Called(ctx, studentID, eventID)
	return args.Error(0)
}
func (m *MockOrientamentoService) GetMyEvents(ctx context.Context, studentID string) ([]Participation, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Participation), args.Error(1)
}
func (m *MockOrientamentoService) MarkAttendance(ctx context.Context, eventID, studentID string) error {
	args := m.Called(ctx, eventID, studentID)
	return args.Error(0)
}
func (m *MockOrientamentoService) SavePreference(ctx context.Context, studentID string, pref StudentPreference) error {
	args := m.Called(ctx, studentID, pref)
	return args.Error(0)
}
func (m *MockOrientamentoService) GetPreference(ctx context.Context, studentID string) (*StudentPreference, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*StudentPreference), args.Error(1)
}
func (m *MockOrientamentoService) SaveCapolavoro(ctx context.Context, studentID string, c Capolavoro) error {
	args := m.Called(ctx, studentID, c)
	return args.Error(0)
}
func (m *MockOrientamentoService) GetCapolavori(ctx context.Context, studentID string) ([]Capolavoro, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Capolavoro), args.Error(1)
}
func (m *MockOrientamentoService) GetCurriculumStudente(ctx context.Context, studentID string) (*CurriculumStudenteSummary, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CurriculumStudenteSummary), args.Error(1)
}

func setupOrientamentoRouter(svc Service, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if role != "" {
			c.Set("role", role)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	handler := NewHandler(svc)
	handler.RegisterRoutes(r.Group("/api/v1"))
	return r
}

func TestHandler_Orientamento(t *testing.T) {
	mockSvc := new(MockOrientamentoService)

	// 1. Unauthorized & Forbidden
	rAnon := setupOrientamentoRouter(mockSvc, "", "", "")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orientamento/events", nil)
	w := httptest.NewRecorder()
	rAnon.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	rStudent := setupOrientamentoRouter(mockSvc, "s-1", "student", "school-1")
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/orientamento/events", bytes.NewBufferString(`{}`))
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 2. Teacher routes
	rTeacher := setupOrientamentoRouter(mockSvc, "t-1", "teacher", "school-1")

	// CreateEvent
	reqEvent := CreateEventRequest{
		Title:    "Open Day",
		Date:     time.Now().Format(time.RFC3339),
		EndDate:  time.Now().Add(2 * time.Hour).Format(time.RFC3339),
		Category: "university",
		Hours:    4.0,
	}
	mockSvc.On("CreateEvent", mock.Anything, "t-1", "school-1", reqEvent).Return(nil).Once()
	body, _ := json.Marshal(reqEvent)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/orientamento/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// GetEvents
	mockSvc.On("GetEvents", mock.Anything, "school-1").Return([]Event{{ID: "ev-1", Title: "Open Day"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/orientamento/events", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// MarkAttendance
	mockSvc.On("MarkAttendance", mock.Anything, "ev-1", "s-1").Return(nil).Once()
	attBody, _ := json.Marshal(map[string]string{"student_id": "s-1"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/orientamento/events/ev-1/attendance", bytes.NewBuffer(attBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Student routes
	// RegisterStudent
	mockSvc.On("RegisterStudent", mock.Anything, "s-1", "ev-1").Return(nil).Once()
	regBody, _ := json.Marshal(map[string]string{"event_id": "ev-1"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/orientamento/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GetMyEvents
	mockSvc.On("GetMyEvents", mock.Anything, "s-1").Return([]Participation{}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/orientamento/my-events", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// SavePreference
	pref := StudentPreference{PreferredTrack: "universita", TargetField: "Informatica"}
	mockSvc.On("SavePreference", mock.Anything, "s-1", pref).Return(nil).Once()
	prefBody, _ := json.Marshal(pref)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/orientamento/preference", bytes.NewBuffer(prefBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GetPreference
	mockSvc.On("GetPreference", mock.Anything, "s-1").Return(&pref, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/orientamento/preference", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// SaveCapolavoro
	capo := Capolavoro{Title: "Progetto Robot"}
	mockSvc.On("SaveCapolavoro", mock.Anything, "s-1", capo).Return(nil).Once()
	capoBody, _ := json.Marshal(capo)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/orientamento/capolavoro", bytes.NewBuffer(capoBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// GetCapolavori
	mockSvc.On("GetCapolavori", mock.Anything, "s-1").Return([]Capolavoro{capo}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/orientamento/capolavori", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GetCurriculumStudente
	mockSvc.On("GetCurriculumStudente", mock.Anything, "s-1").Return(&CurriculumStudenteSummary{StudentID: "s-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/orientamento/curriculum-studente", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockSvc.AssertExpectations(t)
}
