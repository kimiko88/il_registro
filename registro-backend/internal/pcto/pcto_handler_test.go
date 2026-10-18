package pcto

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPCTOHttpService struct {
	mock.Mock
}

func (m *MockPCTOHttpService) CreateProject(ctx context.Context, schoolID, actorRole, teacherID string, req CreateProjectRequest) error {
	args := m.Called(ctx, schoolID, actorRole, teacherID, req)
	return args.Error(0)
}

func (m *MockPCTOHttpService) GetProjects(ctx context.Context, schoolID string) ([]Project, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Project), args.Error(1)
}

func (m *MockPCTOHttpService) UpdateProject(ctx context.Context, actorID, actorRole, id string, req CreateProjectRequest) error {
	args := m.Called(ctx, actorID, actorRole, id, req)
	return args.Error(0)
}

func (m *MockPCTOHttpService) DeleteProject(ctx context.Context, actorID, actorRole, id string) error {
	args := m.Called(ctx, actorID, actorRole, id)
	return args.Error(0)
}

func (m *MockPCTOHttpService) AssignStudent(ctx context.Context, actorRole, projectID, studentID string) error {
	args := m.Called(ctx, actorRole, projectID, studentID)
	return args.Error(0)
}

func (m *MockPCTOHttpService) LogHours(ctx context.Context, studentID string, req LogHourRequest) error {
	args := m.Called(ctx, studentID, req)
	return args.Error(0)
}

func (m *MockPCTOHttpService) GetMyProjects(ctx context.Context, studentID string) ([]Project, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Project), args.Error(1)
}

func (m *MockPCTOHttpService) GetMyProjectDetails(ctx context.Context, studentID, projectID string) (*Participation, []HourLog, error) {
	args := m.Called(ctx, studentID, projectID)
	var p *Participation
	if args.Get(0) != nil {
		p = args.Get(0).(*Participation)
	}
	var logs []HourLog
	if args.Get(1) != nil {
		logs = args.Get(1).([]HourLog)
	}
	return p, logs, args.Error(2)
}

func (m *MockPCTOHttpService) CreateCompany(ctx context.Context, schoolID, actorRole string, c Company) error {
	args := m.Called(ctx, schoolID, actorRole, c)
	return args.Error(0)
}

func (m *MockPCTOHttpService) GetCompanies(ctx context.Context, schoolID string) ([]Company, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Company), args.Error(1)
}

func (m *MockPCTOHttpService) GetStats(ctx context.Context, schoolID string) (*PCTOStats, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PCTOStats), args.Error(1)
}

func (m *MockPCTOHttpService) ApproveHours(ctx context.Context, actorID, actorRole, logID string, approved bool) error {
	args := m.Called(ctx, actorID, actorRole, logID, approved)
	return args.Error(0)
}

func setupTestPCTORouter(svc Service, userID, role, schoolID string) *gin.Engine {
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
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestPCTOHandlers_AllRoutes(t *testing.T) {
	svc := new(MockPCTOHttpService)

	rTeacher := setupTestPCTORouter(svc, "teacher-1", "teacher", "sch-1")

	// 1. GET /pcto/stats
	svc.On("GetStats", mock.Anything, "sch-1").Return(&PCTOStats{TotalProjects: 5, ActiveCompanies: 3}, nil).Once()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/pcto/stats", nil)
	w := httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. POST /pcto/projects - Create Project
	svc.On("CreateProject", mock.Anything, "sch-1", "teacher", "teacher-1", mock.AnythingOfType("pcto.CreateProjectRequest")).Return(nil).Once()
	projectPayload := `{"title":"Stage Robotica","type":"External","start_date":"2026-06-01","end_date":"2026-06-30","total_hours":80}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pcto/projects", bytes.NewBufferString(projectPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 3. GET /pcto/projects - List Projects
	svc.On("GetProjects", mock.Anything, "sch-1").Return([]Project{{ID: "p1", Title: "Stage Robotica"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pcto/projects", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. PUT /pcto/projects/:id - Update Project
	svc.On("UpdateProject", mock.Anything, "teacher-1", "teacher", "p1", mock.AnythingOfType("pcto.CreateProjectRequest")).Return(nil).Once()
	updatePayload := `{"title":"Stage Robotica Avanzata","type":"External","start_date":"2026-06-01","end_date":"2026-06-30","total_hours":100}`
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/pcto/projects/p1", bytes.NewBufferString(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. POST /pcto/projects/:id/students - Assign Student
	svc.On("AssignStudent", mock.Anything, "teacher", "p1", "std-1").Return(nil).Once()
	assignPayload := `{"student_id":"std-1"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pcto/projects/p1/students", bytes.NewBufferString(assignPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. DELETE /pcto/projects/:id - Delete Project
	svc.On("DeleteProject", mock.Anything, "teacher-1", "teacher", "p1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/pcto/projects/p1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. Companies: POST & GET /pcto/companies
	svc.On("CreateCompany", mock.Anything, "sch-1", "teacher", mock.AnythingOfType("pcto.Company")).Return(nil).Once()
	compPayload := `{"name":"Tech Solutions SRL","address":"Via Roma 10"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pcto/companies", bytes.NewBufferString(compPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	svc.On("GetCompanies", mock.Anything, "sch-1").Return([]Company{{ID: "c1", Name: "Tech Solutions SRL"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pcto/companies", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. Student endpoints
	rStudent := setupTestPCTORouter(svc, "std-1", "student", "sch-1")

	// GET /pcto/my-projects
	svc.On("GetMyProjects", mock.Anything, "std-1").Return([]Project{{ID: "p1", Title: "Stage Robotica"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pcto/my-projects", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// POST /pcto/hours
	svc.On("LogHours", mock.Anything, "std-1", mock.AnythingOfType("pcto.LogHourRequest")).Return(nil).Once()
	hoursPayload := `{"project_id":"p1","date":"2026-06-10","hours":4.5,"activity":"Programmazione PLC"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pcto/hours", bytes.NewBufferString(hoursPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// GET /pcto/my-projects/:id
	svc.On("GetMyProjectDetails", mock.Anything, "std-1", "p1").Return(&Participation{ID: "part-1", HoursCompleted: 4.5}, []HourLog{{ID: "hl-1", Hours: 4.5}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pcto/my-projects/p1", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 9. Tutor / Admin: POST /pcto/hours/:id/approve
	svc.On("ApproveHours", mock.Anything, "teacher-1", "teacher", "hl-1", true).Return(nil).Once()
	approvePayload := `{"approved":true}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pcto/hours/hl-1/approve", bytes.NewBufferString(approvePayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 10. Unauthorized checks
	// Student attempting staff endpoint
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pcto/projects", bytes.NewBufferString(projectPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Error handling
	svc.On("GetProjects", mock.Anything, "sch-1").Return(nil, errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pcto/projects", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	svc.AssertExpectations(t)
}

func TestPCTOService_ExtendedGuards(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)
	ctx := context.Background()

	// 1. CreateProject guards
	err := svc.CreateProject(ctx, "", "teacher", "t1", CreateProjectRequest{})
	assert.ErrorContains(t, err, "school_id is required")

	err = svc.CreateProject(ctx, "sch-1", "teacher", "t1", CreateProjectRequest{Title: ""})
	assert.ErrorContains(t, err, "title is required")

	err = svc.CreateProject(ctx, "sch-1", "teacher", "t1", CreateProjectRequest{
		Title: "P1", StartDate: "2026-06-10", EndDate: "2026-06-01",
	})
	assert.ErrorContains(t, err, "cannot be before")

	// 2. AssignStudent guards
	err = svc.AssignStudent(ctx, "teacher", "", "std-1")
	assert.ErrorContains(t, err, "required")

	// 3. LogHours guards
	err = svc.LogHours(ctx, "", LogHourRequest{})
	assert.ErrorContains(t, err, "required")

	// 4. GetMyProjects guards
	_, err = svc.GetMyProjects(ctx, "")
	assert.ErrorContains(t, err, "required")

	// 5. GetMyProjectDetails guards
	_, _, err = svc.GetMyProjectDetails(ctx, "", "p1")
	assert.ErrorContains(t, err, "required")

	// 6. CreateCompany guards
	err = svc.CreateCompany(ctx, "", "teacher", Company{})
	assert.ErrorContains(t, err, "school_id is required")

	err = svc.CreateCompany(ctx, "sch-1", "teacher", Company{})
	assert.ErrorContains(t, err, "company name is required")

	// 7. UpdateProject & DeleteProject guards
	err = svc.UpdateProject(ctx, "t1", "teacher", "", CreateProjectRequest{})
	assert.ErrorContains(t, err, "project id is required")

	err = svc.DeleteProject(ctx, "t1", "teacher", "")
	assert.ErrorContains(t, err, "project id is required")

	// 8. ApproveHours guards
	err = svc.ApproveHours(ctx, "t1", "teacher", "", true)
	assert.ErrorContains(t, err, "log_id is required")

	// 9. UpdateProject & DeleteProject execution
	err = svc.UpdateProject(ctx, "other-teacher", "teacher", "p1", CreateProjectRequest{Title: "Title"})
	assert.ErrorContains(t, err, "unauthorized")

	err = svc.UpdateProject(ctx, "t1", "admin", "p1", CreateProjectRequest{Title: "Title Admin"})
	assert.NoError(t, err)

	err = svc.DeleteProject(ctx, "other-teacher", "teacher", "p1")
	assert.ErrorContains(t, err, "unauthorized")

	err = svc.DeleteProject(ctx, "t1", "admin", "p1")
	assert.NoError(t, err)

	// 10. GetMyProjectDetails & GetStats & GetProjects & Companies
	part, logs, err := svc.GetMyProjectDetails(ctx, "std-1", "p1")
	assert.NoError(t, err)
	assert.NotNil(t, part)
	assert.Empty(t, logs)

	projs, err := svc.GetProjects(ctx, "sch-1")
	assert.NoError(t, err)
	assert.Empty(t, projs)

	stats, err := svc.GetStats(ctx, "sch-1")
	assert.NoError(t, err)
	assert.Nil(t, stats)

	compList, err := svc.GetCompanies(ctx, "sch-1")
	assert.NoError(t, err)
	assert.Empty(t, compList)
}
