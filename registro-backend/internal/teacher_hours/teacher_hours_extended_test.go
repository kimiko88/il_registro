package teacher_hours_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/teacher_hours"
)

type mockTeacherHoursRepo struct {
	mock.Mock
}

func (m *mockTeacherHoursRepo) GetOrCreateObligation(schoolID, teacherID, academicYear string) (*teacher_hours.TeacherHourObligation, error) {
	args := m.Called(schoolID, teacherID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*teacher_hours.TeacherHourObligation), args.Error(1)
}

func (m *mockTeacherHoursRepo) GetObligation(obligationID string) (*teacher_hours.TeacherHourObligation, error) {
	args := m.Called(obligationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*teacher_hours.TeacherHourObligation), args.Error(1)
}

func (m *mockTeacherHoursRepo) ListLogsByObligation(obligationID string) ([]teacher_hours.TeacherHourLog, error) {
	args := m.Called(obligationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]teacher_hours.TeacherHourLog), args.Error(1)
}

func (m *mockTeacherHoursRepo) ListShortLeaves(schoolID, teacherID string) ([]teacher_hours.TeacherShortLeave, error) {
	args := m.Called(schoolID, teacherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]teacher_hours.TeacherShortLeave), args.Error(1)
}

func (m *mockTeacherHoursRepo) LogHourActivity(log *teacher_hours.TeacherHourLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *mockTeacherHoursRepo) CreateShortLeave(leave *teacher_hours.TeacherShortLeave) error {
	args := m.Called(leave)
	return args.Error(0)
}

func (m *mockTeacherHoursRepo) UpdateShortLeaveStatus(leaveID, status, notes string) error {
	args := m.Called(leaveID, status, notes)
	return args.Error(0)
}

func (m *mockTeacherHoursRepo) RecordShortLeaveRecovery(leaveID string, recoveredHours float64, recoveredDate string) error {
	args := m.Called(leaveID, recoveredHours, recoveredDate)
	return args.Error(0)
}

func (m *mockTeacherHoursRepo) ListObligationsBySchool(schoolID, academicYear string) ([]teacher_hours.TeacherHourObligation, error) {
	args := m.Called(schoolID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]teacher_hours.TeacherHourObligation), args.Error(1)
}

func setupTestRouter(svc teacher_hours.Service, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "sch-test-123")
		c.Set("user_id", "teacher-test-456")
		c.Set("user_role", role)
		c.Next()
	})
	h := teacher_hours.NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestTeacherHoursService_EdgeCases(t *testing.T) {
	repo := new(mockTeacherHoursRepo)
	svc := teacher_hours.NewService(repo)

	// ReviewShortLeave
	assert.Error(t, svc.ReviewShortLeave("", "approved", ""))
	assert.Error(t, svc.ReviewShortLeave("l1", "invalid_status", ""))
	repo.On("UpdateShortLeaveStatus", "l1", "approved", "ok").Return(nil).Once()
	assert.NoError(t, svc.ReviewShortLeave("l1", "approved", "ok"))

	// RecordLeaveRecovery
	assert.Error(t, svc.RecordLeaveRecovery("", teacher_hours.RecoverLeaveRequest{RecoveredHours: 1.0}))
	assert.Error(t, svc.RecordLeaveRecovery("l1", teacher_hours.RecoverLeaveRequest{RecoveredHours: 0}))
	repo.On("RecordShortLeaveRecovery", "l1", 1.5, "2026-03-12").Return(nil).Once()
	assert.NoError(t, svc.RecordLeaveRecovery("l1", teacher_hours.RecoverLeaveRequest{RecoveredHours: 1.5, RecoveredDate: "2026-03-12"}))

	// ListSchoolObligations
	_, err := svc.ListSchoolObligations("", "2025/2026")
	assert.Error(t, err)
	repo.On("ListObligationsBySchool", "sch-1", "2025/2026").Return([]teacher_hours.TeacherHourObligation{{ID: "o1"}}, nil).Once()
	obls, err := svc.ListSchoolObligations("sch-1", "2025/2026")
	assert.NoError(t, err)
	assert.Len(t, obls, 1)

	repo.AssertExpectations(t)
}

func TestTeacherHoursHandler_Endpoints(t *testing.T) {
	repo := new(mockTeacherHoursRepo)
	svc := teacher_hours.NewService(repo)
	r := setupTestRouter(svc, "teacher")

	// 1. GET /teacher-hours/dashboard
	repo.On("GetOrCreateObligation", "sch-test-123", "teacher-test-456", "2025/2026").Return(&teacher_hours.TeacherHourObligation{ID: "obl-1"}, nil).Once()
	repo.On("ListLogsByObligation", "obl-1").Return([]teacher_hours.TeacherHourLog{}, nil).Once()
	repo.On("ListShortLeaves", "sch-test-123", "teacher-test-456").Return([]teacher_hours.TeacherShortLeave{}, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/teacher-hours/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Admin viewing other teacher dashboard
	adminRouter := setupTestRouter(svc, "admin")
	repo.On("GetOrCreateObligation", "sch-test-123", "other-teacher", "2025/2026").Return(&teacher_hours.TeacherHourObligation{ID: "obl-2"}, nil).Once()
	repo.On("ListLogsByObligation", "obl-2").Return([]teacher_hours.TeacherHourLog{}, nil).Once()
	repo.On("ListShortLeaves", "sch-test-123", "other-teacher").Return([]teacher_hours.TeacherShortLeave{}, nil).Once()

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/teacher-hours/dashboard?teacher_id=other-teacher", nil)
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. POST /teacher-hours/logs
	repo.On("LogHourActivity", mock.Anything).Return(nil).Once()
	logPayload := `{"obligation_id":"obl-1","category":"collegio","hours":2.0,"activity_date":"2026-03-10","description":"Collegio docenti marzo"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/teacher-hours/logs", bytes.NewBufferString(logPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/teacher-hours/logs", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 3. POST /teacher-hours/leaves
	repo.On("CreateShortLeave", mock.Anything).Return(nil).Once()
	leavePayload := `{"hours":1.5,"reason":"Visita dentistica","request_date":"2026-03-15"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/teacher-hours/leaves", bytes.NewBufferString(leavePayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/teacher-hours/leaves", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 4. PATCH /teacher-hours/leaves/:id/review
	repo.On("UpdateShortLeaveStatus", "leave-1", "approved", "Concesso").Return(nil).Once()
	reviewPayload := `{"status":"approved","notes":"Concesso"}`
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/teacher-hours/leaves/leave-1/review", bytes.NewBufferString(reviewPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/teacher-hours/leaves/leave-1/review", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 5. POST /teacher-hours/leaves/:id/recover
	repo.On("RecordShortLeaveRecovery", "leave-1", 1.5, "2026-03-20").Return(nil).Once()
	recPayload := `{"recovered_hours":1.5,"recovered_date":"2026-03-20"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/teacher-hours/leaves/leave-1/recover", bytes.NewBufferString(recPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/teacher-hours/leaves/leave-1/recover", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 6. GET /teacher-hours/school-summary
	repo.On("ListObligationsBySchool", "sch-test-123", "2025/2026").Return([]teacher_hours.TeacherHourObligation{{ID: "obl-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/teacher-hours/school-summary", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListObligationsBySchool", "sch-test-123", "2025/2026").Return(nil, errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/teacher-hours/school-summary", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	repo.AssertExpectations(t)
}
