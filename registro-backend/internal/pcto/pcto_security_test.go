package pcto_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/pcto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPctoRepo struct {
	mock.Mock
}

func (m *mockPctoRepo) CreateProject(ctx context.Context, p *pcto.Project) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *mockPctoRepo) GetProjects(ctx context.Context, schoolID string) ([]pcto.Project, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]pcto.Project), args.Error(1)
}
func (m *mockPctoRepo) GetProjectByID(ctx context.Context, id string) (*pcto.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pcto.Project), args.Error(1)
}
func (m *mockPctoRepo) UpdateProject(ctx context.Context, p *pcto.Project) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *mockPctoRepo) DeleteProject(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockPctoRepo) AssignStudent(ctx context.Context, participation *pcto.Participation) error {
	args := m.Called(ctx, participation)
	return args.Error(0)
}
func (m *mockPctoRepo) GetParticipationsByProject(ctx context.Context, projectID string) ([]pcto.Participation, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]pcto.Participation), args.Error(1)
}
func (m *mockPctoRepo) GetParticipationsByStudent(ctx context.Context, studentID string) ([]pcto.Participation, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]pcto.Participation), args.Error(1)
}
func (m *mockPctoRepo) GetParticipation(ctx context.Context, projectID, studentID string) (*pcto.Participation, error) {
	args := m.Called(ctx, projectID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pcto.Participation), args.Error(1)
}
func (m *mockPctoRepo) LogHours(ctx context.Context, h *pcto.HourLog) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}
func (m *mockPctoRepo) GetHours(ctx context.Context, participationID string) ([]pcto.HourLog, error) {
	args := m.Called(ctx, participationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]pcto.HourLog), args.Error(1)
}
func (m *mockPctoRepo) GetHourLogByID(ctx context.Context, id string) (*pcto.HourLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pcto.HourLog), args.Error(1)
}
func (m *mockPctoRepo) GetParticipationByID(ctx context.Context, id string) (*pcto.Participation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pcto.Participation), args.Error(1)
}
func (m *mockPctoRepo) VerifyHours(ctx context.Context, hourID, teacherID string) error {
	args := m.Called(ctx, hourID, teacherID)
	return args.Error(0)
}
func (m *mockPctoRepo) UpdateHourLogStatus(ctx context.Context, logID, status string) error {
	args := m.Called(ctx, logID, status)
	return args.Error(0)
}
func (m *mockPctoRepo) CreateCompany(ctx context.Context, c *pcto.Company) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *mockPctoRepo) GetCompanies(ctx context.Context, schoolID string) ([]pcto.Company, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]pcto.Company), args.Error(1)
}
func (m *mockPctoRepo) GetStats(ctx context.Context, schoolID string) (*pcto.PCTOStats, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pcto.PCTOStats), args.Error(1)
}

func buildPctoEngine(repo pcto.Repository, userID, role, schoolID string) *gin.Engine {
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

	svc := pcto.NewService(repo)
	h := pcto.NewHandler(svc)
	h.RegisterRoutes(r.Group("/api"))

	return r
}

func TestPcto_GetProjects_RBAC(t *testing.T) {
	t.Run("Student blocked from school-wide projects with 403", func(t *testing.T) {
		r := buildPctoEngine(nil, "student-1", "student", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/pcto/projects", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Parent blocked from school-wide projects with 403", func(t *testing.T) {
		r := buildPctoEngine(nil, "parent-1", "parent", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/pcto/projects", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Teacher allowed to view school-wide projects", func(t *testing.T) {
		mockRepo := new(mockPctoRepo)
		mockRepo.On("GetProjects", mock.Anything, "school-1").Return([]pcto.Project{}, nil).Once()

		r := buildPctoEngine(mockRepo, "teacher-1", "teacher", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/pcto/projects", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
