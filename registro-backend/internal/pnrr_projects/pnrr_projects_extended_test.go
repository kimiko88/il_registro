package pnrr_projects_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/pnrr_projects"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateProject(p *pnrr_projects.PNRRProject) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *mockRepo) ListProjectsBySchool(schoolID, academicYear string) ([]pnrr_projects.PNRRProject, error) {
	args := m.Called(schoolID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]pnrr_projects.PNRRProject), args.Error(1)
}

func (m *mockRepo) GetProject(id string) (*pnrr_projects.PNRRProject, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pnrr_projects.PNRRProject), args.Error(1)
}

func (m *mockRepo) EnrollStudent(ps *pnrr_projects.PNRRProjectStudent) error {
	args := m.Called(ps)
	return args.Error(0)
}

func (m *mockRepo) ListProjectStudents(projectID string) ([]pnrr_projects.PNRRProjectStudent, error) {
	args := m.Called(projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]pnrr_projects.PNRRProjectStudent), args.Error(1)
}

func (m *mockRepo) RecordLesson(l *pnrr_projects.PNRRLessonAttendance, attendedStudentIDs []string) error {
	args := m.Called(l, attendedStudentIDs)
	return args.Error(0)
}

func (m *mockRepo) ListLessons(projectID string) ([]pnrr_projects.PNRRLessonAttendance, error) {
	args := m.Called(projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]pnrr_projects.PNRRLessonAttendance), args.Error(1)
}

func setupTestRouter(svc pnrr_projects.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "sch-test-123")
		c.Set("user_id", "teacher-test-456")
		c.Next()
	})
	h := pnrr_projects.NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestPNRRProjectsService_Full(t *testing.T) {
	repo := new(mockRepo)
	svc := pnrr_projects.NewService(repo)

	// CreateProject validations
	assert.Error(t, svc.CreateProject(nil))
	assert.Error(t, svc.CreateProject(&pnrr_projects.PNRRProject{SchoolID: ""}))
	assert.Error(t, svc.CreateProject(&pnrr_projects.PNRRProject{SchoolID: "sch-1", CUPCode: ""}))
	assert.Error(t, svc.CreateProject(&pnrr_projects.PNRRProject{SchoolID: "sch-1", CUPCode: "CUP123", ProjectTitle: ""}))
	assert.Error(t, svc.CreateProject(&pnrr_projects.PNRRProject{SchoolID: "sch-1", CUPCode: "CUP123", ProjectTitle: "Titolo", TotalHours: 0}))

	repo.On("CreateProject", mock.Anything).Return(nil).Once()
	p := &pnrr_projects.PNRRProject{SchoolID: "sch-1", CUPCode: "CUP123", ProjectTitle: "Titolo", TotalHours: 30}
	assert.NoError(t, svc.CreateProject(p))
	assert.Equal(t, "attivo", p.Status)

	// ListProjects
	_, err := svc.ListProjects("", "2025/2026")
	assert.Error(t, err)

	repo.On("ListProjectsBySchool", "sch-1", "2025/2026").Return([]pnrr_projects.PNRRProject{{ID: "p1"}}, nil).Once()
	plist, err := svc.ListProjects("sch-1", "2025/2026")
	assert.NoError(t, err)
	assert.Len(t, plist, 1)

	// GetProject
	_, err = svc.GetProject("")
	assert.Error(t, err)

	repo.On("GetProject", "p1").Return(&pnrr_projects.PNRRProject{ID: "p1"}, nil).Once()
	p1, err := svc.GetProject("p1")
	assert.NoError(t, err)
	assert.Equal(t, "p1", p1.ID)

	// EnrollStudent
	assert.Error(t, svc.EnrollStudent("", "s1"))
	assert.Error(t, svc.EnrollStudent("p1", ""))

	repo.On("EnrollStudent", mock.Anything).Return(nil).Once()
	assert.NoError(t, svc.EnrollStudent("p1", "s1"))

	// ListStudents
	_, err = svc.ListStudents("")
	assert.Error(t, err)

	repo.On("ListProjectStudents", "p1").Return([]pnrr_projects.PNRRProjectStudent{{StudentID: "s1"}}, nil).Once()
	slist, err := svc.ListStudents("p1")
	assert.NoError(t, err)
	assert.Len(t, slist, 1)

	// RecordLesson
	assert.Error(t, svc.RecordLesson(nil, nil))
	assert.Error(t, svc.RecordLesson(&pnrr_projects.PNRRLessonAttendance{ProjectID: ""}, nil))
	assert.Error(t, svc.RecordLesson(&pnrr_projects.PNRRLessonAttendance{ProjectID: "p1", Hours: 0}, nil))
	assert.Error(t, svc.RecordLesson(&pnrr_projects.PNRRLessonAttendance{ProjectID: "p1", Hours: 2, Topic: ""}, nil))

	repo.On("RecordLesson", mock.Anything, mock.Anything).Return(nil).Once()
	assert.NoError(t, svc.RecordLesson(&pnrr_projects.PNRRLessonAttendance{ProjectID: "p1", Hours: 2, Topic: "Coding"}, []string{"s1"}))

	// ListLessons
	_, err = svc.ListLessons("")
	assert.Error(t, err)

	repo.On("ListLessons", "p1").Return([]pnrr_projects.PNRRLessonAttendance{{ID: "l1"}}, nil).Once()
	llist, err := svc.ListLessons("p1")
	assert.NoError(t, err)
	assert.Len(t, llist, 1)

	repo.AssertExpectations(t)
}

func TestPNRRProjectsHandler_Endpoints(t *testing.T) {
	repo := new(mockRepo)
	svc := pnrr_projects.NewService(repo)
	r := setupTestRouter(svc)

	// 1. POST /pnrr-projects
	repo.On("CreateProject", mock.Anything).Return(nil).Once()
	payload := `{"cup_code":"CUP999","project_title":"STEM Lab","total_hours":40}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/pnrr-projects", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pnrr-projects", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 2. GET /pnrr-projects
	repo.On("ListProjectsBySchool", "sch-test-123", "2025/2026").Return([]pnrr_projects.PNRRProject{{ID: "p1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pnrr-projects?academic_year=2025/2026", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListProjectsBySchool", "sch-test-123", "2025/2026").Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pnrr-projects?academic_year=2025/2026", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 3. GET /pnrr-projects/:id
	repo.On("GetProject", "p1").Return(&pnrr_projects.PNRRProject{ID: "p1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pnrr-projects/p1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetProject", "missing").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pnrr-projects/missing", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 4. POST /pnrr-projects/:id/enroll
	repo.On("EnrollStudent", mock.Anything).Return(nil).Once()
	enrollPayload := `{"student_id":"s1"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pnrr-projects/p1/enroll", bytes.NewBufferString(enrollPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pnrr-projects/p1/enroll", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 5. GET /pnrr-projects/:id/students
	repo.On("ListProjectStudents", "p1").Return([]pnrr_projects.PNRRProjectStudent{{StudentID: "s1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pnrr-projects/p1/students", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListProjectStudents", "perr").Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pnrr-projects/perr/students", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 6. POST /pnrr-projects/:id/lessons
	repo.On("RecordLesson", mock.Anything, mock.Anything).Return(nil).Once()
	lessonPayload := `{"lesson_date":"2026-03-10","hours":2,"topic":"Robotica","attended_student_ids":["s1","s2"]}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pnrr-projects/p1/lessons", bytes.NewBufferString(lessonPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pnrr-projects/p1/lessons", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 7. GET /pnrr-projects/:id/lessons
	repo.On("ListLessons", "p1").Return([]pnrr_projects.PNRRLessonAttendance{{ID: "l1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pnrr-projects/p1/lessons", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListLessons", "perr").Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pnrr-projects/perr/lessons", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	repo.AssertExpectations(t)
}
