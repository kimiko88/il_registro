package timetables

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ─── Mock Service ─────────────────────────────────────────────────────────────

type mockTimetableService struct {
	classSched       []ClassSchedule
	teacherSched     []ClassSchedule
	getByClassErr    error
	getMySchedErr    error
	getTeacherErr    error
	updateErr        error
	updateTeacherErr error
}

func (m *mockTimetableService) GetByClass(_ context.Context, _, _, _, _ string) ([]ClassSchedule, error) {
	return m.classSched, m.getByClassErr
}
func (m *mockTimetableService) GetMySchedule(_ context.Context, _, _ string) ([]ClassSchedule, error) {
	return m.teacherSched, m.getMySchedErr
}
func (m *mockTimetableService) GetByTeacher(_ context.Context, _, _, _, _ string) ([]ClassSchedule, error) {
	return m.teacherSched, m.getTeacherErr
}
func (m *mockTimetableService) Update(_ context.Context, _, _, _, _ string, _ []ScheduleEntry) error {
	return m.updateErr
}
func (m *mockTimetableService) UpdateTeacher(_ context.Context, _, _, _, _ string, _ []TeacherScheduleEntry) error {
	return m.updateTeacherErr
}

// ─── Router setup ─────────────────────────────────────────────────────────────

func setupTimetableRouter(svc Service, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewHandlerWithService(svc)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("school_id", schoolID)
		c.Next()
	})
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)
	return r
}

// ─── GetByClass ───────────────────────────────────────────────────────────────

func TestTimetableHandler_GetByClass_Unauthorized(t *testing.T) {
	svc := &mockTimetableService{classSched: []ClassSchedule{}}
	r := setupTimetableRouter(svc, "", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/classes/class-1/schedule", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTimetableHandler_GetByClass_Success(t *testing.T) {
	svc := &mockTimetableService{
		classSched: []ClassSchedule{
			{ID: "s-1", ClassID: "class-1", DayOfWeek: 1, HourIndex: 1, SubjectID: "math"},
		},
	}
	r := setupTimetableRouter(svc, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/classes/class-1/schedule", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var sched []ClassSchedule
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &sched))
	assert.Len(t, sched, 1)
}

func TestTimetableHandler_GetByClass_Forbidden(t *testing.T) {
	svc := &mockTimetableService{getByClassErr: errors.New("forbidden: non puoi visualizzare l'orario di un'altra classe")}
	r := setupTimetableRouter(svc, "student-1", "student", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/classes/class-99/schedule", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestTimetableHandler_GetByClass_InternalError(t *testing.T) {
	svc := &mockTimetableService{getByClassErr: errors.New("database error")}
	r := setupTimetableRouter(svc, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/classes/class-1/schedule", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Update (class schedule) ──────────────────────────────────────────────────

func TestTimetableHandler_Update_Unauthorized(t *testing.T) {
	svc := &mockTimetableService{}
	r := setupTimetableRouter(svc, "", "admin", "school-1")
	body, _ := json.Marshal(UpdateScheduleRequest{Entries: []ScheduleEntry{{DayOfWeek: 1, HourIndex: 1, SubjectID: "math"}}})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/classes/class-1/schedule", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTimetableHandler_Update_Forbidden(t *testing.T) {
	svc := &mockTimetableService{updateErr: errors.New("forbidden: la classe appartiene a un'altra scuola")}
	r := setupTimetableRouter(svc, "admin-1", "admin", "school-1")
	body, _ := json.Marshal(UpdateScheduleRequest{Entries: []ScheduleEntry{{DayOfWeek: 1, HourIndex: 1, SubjectID: "math"}}})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/classes/class-1/schedule", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestTimetableHandler_Update_Success(t *testing.T) {
	svc := &mockTimetableService{}
	r := setupTimetableRouter(svc, "admin-1", "admin", "school-1")
	body, _ := json.Marshal(UpdateScheduleRequest{
		Entries: []ScheduleEntry{
			{DayOfWeek: 1, HourIndex: 1, SubjectID: "math"},
			{DayOfWeek: 2, HourIndex: 2, SubjectID: "science"},
		},
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/classes/class-1/schedule", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "updated")
}

func TestTimetableHandler_Update_PUT_Success(t *testing.T) {
	svc := &mockTimetableService{}
	r := setupTimetableRouter(svc, "admin-1", "admin", "school-1")
	body, _ := json.Marshal(UpdateScheduleRequest{Entries: []ScheduleEntry{}})
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/classes/class-1/schedule", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── GetMySchedule ────────────────────────────────────────────────────────────

func TestTimetableHandler_GetMySchedule_Success(t *testing.T) {
	svc := &mockTimetableService{
		teacherSched: []ClassSchedule{{ID: "s-2", DayOfWeek: 3, HourIndex: 2}},
	}
	r := setupTimetableRouter(svc, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/timetables/my-schedule", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTimetableHandler_GetMySchedule_NotFound(t *testing.T) {
	svc := &mockTimetableService{getMySchedErr: errors.New("teacher schedule not found")}
	r := setupTimetableRouter(svc, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/timetables/my-schedule", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTimetableHandler_GetMySchedule_Forbidden(t *testing.T) {
	svc := &mockTimetableService{getMySchedErr: errors.New("forbidden: role not allowed")}
	r := setupTimetableRouter(svc, "student-1", "student", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/timetables/my-schedule", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── GetByTeacher ─────────────────────────────────────────────────────────────

func TestTimetableHandler_GetByTeacher_Unauthorized(t *testing.T) {
	svc := &mockTimetableService{}
	r := setupTimetableRouter(svc, "", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/timetables/teacher/teacher-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTimetableHandler_GetByTeacher_Success(t *testing.T) {
	svc := &mockTimetableService{
		teacherSched: []ClassSchedule{{ID: "s-3", DayOfWeek: 4, HourIndex: 3}},
	}
	r := setupTimetableRouter(svc, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/timetables/teacher/teacher-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTimetableHandler_GetByTeacher_AlternativeRoute(t *testing.T) {
	svc := &mockTimetableService{
		teacherSched: []ClassSchedule{{ID: "s-4", DayOfWeek: 5, HourIndex: 1}},
	}
	r := setupTimetableRouter(svc, "principal-1", "principal", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/teachers/teacher-1/schedule", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── UpdateTeacher ────────────────────────────────────────────────────────────

func TestTimetableHandler_UpdateTeacher_Unauthorized(t *testing.T) {
	svc := &mockTimetableService{}
	r := setupTimetableRouter(svc, "", "admin", "school-1")
	body, _ := json.Marshal(UpdateTeacherScheduleRequest{Entries: []TeacherScheduleEntry{}})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/timetables/teacher/t-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTimetableHandler_UpdateTeacher_Forbidden(t *testing.T) {
	forbidErr := errors.New("forbidden: only administrative staff can update teacher schedules")
	svc := &mockTimetableService{updateTeacherErr: forbidErr}
	r := setupTimetableRouter(svc, "teacher-1", "teacher", "school-1")
	body, _ := json.Marshal(UpdateTeacherScheduleRequest{Entries: []TeacherScheduleEntry{}})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/timetables/teacher/t-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestTimetableHandler_UpdateTeacher_Success(t *testing.T) {
	svc := &mockTimetableService{}
	r := setupTimetableRouter(svc, "admin-1", "admin", "school-1")
	body, _ := json.Marshal(UpdateTeacherScheduleRequest{
		Entries: []TeacherScheduleEntry{
			{DayOfWeek: 1, HourIndex: 1, ClassID: "class-1", SubjectID: "math"},
		},
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/timetables/teacher/teacher-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "updated"))
}

func TestTimetableHandler_UpdateTeacher_InternalError(t *testing.T) {
	svc := &mockTimetableService{updateTeacherErr: errors.New("database error")}
	r := setupTimetableRouter(svc, "admin-1", "admin", "school-1")
	body, _ := json.Marshal(UpdateTeacherScheduleRequest{Entries: []TeacherScheduleEntry{}})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/timetables/teacher/teacher-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
