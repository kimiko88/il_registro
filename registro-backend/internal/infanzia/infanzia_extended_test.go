package infanzia_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/infanzia"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateSection(s *infanzia.InfanziaSection) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *mockRepo) ListSections(schoolID, academicYear string) ([]infanzia.InfanziaSection, error) {
	args := m.Called(schoolID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]infanzia.InfanziaSection), args.Error(1)
}

func (m *mockRepo) GetSection(id string) (*infanzia.InfanziaSection, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*infanzia.InfanziaSection), args.Error(1)
}

func (m *mockRepo) SaveDailyAttendancesBatch(sectionID, date string, attendances []infanzia.InfanziaDailyAttendance) error {
	args := m.Called(sectionID, date, attendances)
	return args.Error(0)
}

func (m *mockRepo) GetDailyAttendances(sectionID, date string) ([]infanzia.InfanziaDailyAttendance, error) {
	args := m.Called(sectionID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]infanzia.InfanziaDailyAttendance), args.Error(1)
}

func (m *mockRepo) SaveObservation(obs *infanzia.InfanziaExperienceObservation) error {
	args := m.Called(obs)
	return args.Error(0)
}

func (m *mockRepo) ListObservationsByStudent(studentID string) ([]infanzia.InfanziaExperienceObservation, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]infanzia.InfanziaExperienceObservation), args.Error(1)
}

func (m *mockRepo) CreateDiarioEntry(d *infanzia.InfanziaDiarioBordo) error {
	args := m.Called(d)
	return args.Error(0)
}

func (m *mockRepo) ListDiarioEntries(sectionID string, limit int) ([]infanzia.InfanziaDiarioBordo, error) {
	args := m.Called(sectionID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]infanzia.InfanziaDiarioBordo), args.Error(1)
}

func setupTestRouter(svc infanzia.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "sch-test-123")
		c.Set("user_id", "teacher-test-456")
		c.Next()
	})
	h := infanzia.NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestInfanziaService_FullCoverage(t *testing.T) {
	repo := new(mockRepo)
	svc := infanzia.NewService(repo)

	// Nil safety checks
	assert.Error(t, svc.CreateSection(nil))
	assert.Error(t, svc.RecordObservation(nil))
	assert.Error(t, svc.AddDiarioStory(nil))

	// CreateSection validation
	assert.Error(t, svc.CreateSection(&infanzia.InfanziaSection{SchoolID: ""}))
	assert.Error(t, svc.CreateSection(&infanzia.InfanziaSection{SchoolID: "sch", SectionName: ""}))

	// ListSections
	assert.Error(t, func() error {
		_, err := svc.ListSections("", "2025/2026")
		return err
	}())
	repo.On("ListSections", "sch-1", "2025/2026").Return([]infanzia.InfanziaSection{{ID: "sec-1"}}, nil).Once()
	sections, err := svc.ListSections("sch-1", "2025/2026")
	assert.NoError(t, err)
	assert.Len(t, sections, 1)

	// GetSection
	repo.On("GetSection", "sec-1").Return(&infanzia.InfanziaSection{ID: "sec-1"}, nil).Once()
	sec, err := svc.GetSection("sec-1")
	assert.NoError(t, err)
	assert.Equal(t, "sec-1", sec.ID)

	// SaveAttendances
	assert.Error(t, svc.SaveAttendances("", "2026-03-10", nil))
	assert.Error(t, svc.SaveAttendances("sec-1", "", nil))
	repo.On("SaveDailyAttendancesBatch", "sec-1", "2026-03-10", mock.Anything).Return(nil).Once()
	assert.NoError(t, svc.SaveAttendances("sec-1", "2026-03-10", []infanzia.InfanziaDailyAttendance{}))

	// GetAttendances
	assert.Error(t, func() error {
		_, err := svc.GetAttendances("", "2026-03-10")
		return err
	}())
	assert.Error(t, func() error {
		_, err := svc.GetAttendances("sec-1", "")
		return err
	}())
	repo.On("GetDailyAttendances", "sec-1", "2026-03-10").Return([]infanzia.InfanziaDailyAttendance{{StudentID: "s1"}}, nil).Once()
	atts, err := svc.GetAttendances("sec-1", "2026-03-10")
	assert.NoError(t, err)
	assert.Len(t, atts, 1)

	// RecordObservation validation
	assert.Error(t, svc.RecordObservation(&infanzia.InfanziaExperienceObservation{StudentID: ""}))
	assert.Error(t, svc.RecordObservation(&infanzia.InfanziaExperienceObservation{StudentID: "s1", RecordedBy: ""}))
	assert.Error(t, svc.RecordObservation(&infanzia.InfanziaExperienceObservation{StudentID: "s1", RecordedBy: "t1", ObservationNotes: ""}))
	assert.Error(t, svc.RecordObservation(&infanzia.InfanziaExperienceObservation{StudentID: "s1", RecordedBy: "t1", ObservationNotes: "Notes", CampoEsperienza: "invalid"}))

	// ListStudentObservations
	assert.Error(t, func() error {
		_, err := svc.ListStudentObservations("")
		return err
	}())
	repo.On("ListObservationsByStudent", "s1").Return([]infanzia.InfanziaExperienceObservation{{ID: "o1"}}, nil).Once()
	obsList, err := svc.ListStudentObservations("s1")
	assert.NoError(t, err)
	assert.Len(t, obsList, 1)

	// AddDiarioStory validation
	assert.Error(t, svc.AddDiarioStory(&infanzia.InfanziaDiarioBordo{SectionID: ""}))
	assert.Error(t, svc.AddDiarioStory(&infanzia.InfanziaDiarioBordo{SectionID: "sec-1", Title: ""}))
	assert.Error(t, svc.AddDiarioStory(&infanzia.InfanziaDiarioBordo{SectionID: "sec-1", Title: "Title", NarrativeStory: ""}))
	repo.On("CreateDiarioEntry", mock.Anything).Return(nil).Once()
	assert.NoError(t, svc.AddDiarioStory(&infanzia.InfanziaDiarioBordo{SectionID: "sec-1", Title: "Title", NarrativeStory: "Story"}))

	// ListDiarioStories validation & limit clamping
	assert.Error(t, func() error {
		_, err := svc.ListDiarioStories("", 10)
		return err
	}())
	repo.On("ListDiarioEntries", "sec-1", 30).Return([]infanzia.InfanziaDiarioBordo{{ID: "d1"}}, nil).Once()
	dList, err := svc.ListDiarioStories("sec-1", -5)
	assert.NoError(t, err)
	assert.Len(t, dList, 1)

	repo.On("ListDiarioEntries", "sec-1", 100).Return([]infanzia.InfanziaDiarioBordo{{ID: "d1"}}, nil).Once()
	dList, err = svc.ListDiarioStories("sec-1", 500)
	assert.NoError(t, err)
	assert.Len(t, dList, 1)

	repo.AssertExpectations(t)
}

func TestInfanziaHandler_Routes(t *testing.T) {
	repo := new(mockRepo)
	svc := infanzia.NewService(repo)
	r := setupTestRouter(svc)

	// 1. POST /sections - Success & Failure
	repo.On("CreateSection", mock.Anything).Return(nil).Once()
	payload := `{"section_name":"Farfalle","color_hex":"#10B981"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/infanzia/sections", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/infanzia/sections", bytes.NewBufferString(`{invalid_json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 2. GET /sections - Success & Failure
	repo.On("ListSections", "sch-test-123", "2025/2026").Return([]infanzia.InfanziaSection{{ID: "sec-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/infanzia/sections?academic_year=2025/2026", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListSections", "sch-test-123", "2025/2026").Return(nil, errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/infanzia/sections?academic_year=2025/2026", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 3. GET /sections/:id - Found & NotFound
	repo.On("GetSection", "sec-1").Return(&infanzia.InfanziaSection{ID: "sec-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/infanzia/sections/sec-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetSection", "sec-missing").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/infanzia/sections/sec-missing", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 4. POST /sections/:id/attendances
	repo.On("SaveDailyAttendancesBatch", "sec-1", "2026-03-10", mock.Anything).Return(nil).Once()
	attPayload := `{"date":"2026-03-10","attendances":[{"student_id":"s1","present":true,"meal_eaten":"completo"}]}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/infanzia/sections/sec-1/attendances", bytes.NewBufferString(attPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/infanzia/sections/sec-1/attendances", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 5. GET /sections/:id/attendances
	repo.On("GetDailyAttendances", "sec-1", "2026-03-10").Return([]infanzia.InfanziaDailyAttendance{}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/infanzia/sections/sec-1/attendances?date=2026-03-10", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. POST /observations
	repo.On("SaveObservation", mock.Anything).Return(nil).Once()
	obsPayload := `{"student_id":"s1","campo_esperienza":"il_se_e_l_altro","observation_notes":"Molto socievole"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/infanzia/observations", bytes.NewBufferString(obsPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 7. GET /observations/by-student
	repo.On("ListObservationsByStudent", "s1").Return([]infanzia.InfanziaExperienceObservation{}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/infanzia/observations/by-student?student_id=s1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. POST /sections/:id/diario
	repo.On("CreateDiarioEntry", mock.Anything).Return(nil).Once()
	diarioPayload := `{"title":"Laboratorio Colori","narrative_story":"Oggi i bimbi hanno dipinto con le dita"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/infanzia/sections/sec-1/diario", bytes.NewBufferString(diarioPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 9. GET /sections/:id/diario
	repo.On("ListDiarioEntries", "sec-1", 30).Return([]infanzia.InfanziaDiarioBordo{}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/infanzia/sections/sec-1/diario?limit=30", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.AssertExpectations(t)
}
