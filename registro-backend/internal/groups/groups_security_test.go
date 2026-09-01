package groups_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/groups"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGroupsRepo struct {
	mock.Mock
}

func (m *mockGroupsRepo) Create(ctx context.Context, g *groups.Group) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}
func (m *mockGroupsRepo) Update(ctx context.Context, g *groups.Group) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}
func (m *mockGroupsRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockGroupsRepo) GetByID(ctx context.Context, id string) (*groups.Group, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*groups.Group), args.Error(1)
}
func (m *mockGroupsRepo) ListBySchool(ctx context.Context, schoolID string) ([]groups.Group, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]groups.Group), args.Error(1)
}
func (m *mockGroupsRepo) ListByTeacher(ctx context.Context, teacherID string) ([]groups.Group, error) {
	args := m.Called(ctx, teacherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]groups.Group), args.Error(1)
}
func (m *mockGroupsRepo) ListByStudent(ctx context.Context, studentID string) ([]groups.Group, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]groups.Group), args.Error(1)
}
func (m *mockGroupsRepo) AddStudents(ctx context.Context, groupID string, studentIDs []string) error {
	args := m.Called(ctx, groupID, studentIDs)
	return args.Error(0)
}
func (m *mockGroupsRepo) RemoveStudent(ctx context.Context, groupID, studentID string) error {
	args := m.Called(ctx, groupID, studentID)
	return args.Error(0)
}
func (m *mockGroupsRepo) GetStudentsInGroup(ctx context.Context, groupID string) ([]groups.GroupStudentInfo, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]groups.GroupStudentInfo), args.Error(1)
}

func buildGroupsEngine(repo groups.Repository, userID, role, schoolID string) *gin.Engine {
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

	svc := groups.NewService(repo)
	h := groups.NewHandler(svc)
	h.RegisterRoutes(r.Group("/api"))

	return r
}

func TestGroups_ListGroups_Security(t *testing.T) {
	t.Run("Student querying another teacher's groups is blocked with 403", func(t *testing.T) {
		r := buildGroupsEngine(nil, "student-1", "student", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/groups?teacher_id=teacher-target", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Teacher querying another teacher's groups is blocked with 403", func(t *testing.T) {
		r := buildGroupsEngine(nil, "teacher-1", "teacher", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/groups?teacher_id=teacher-other", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Teacher querying their own groups is allowed", func(t *testing.T) {
		mockRepo := new(mockGroupsRepo)
		mockRepo.On("ListByTeacher", mock.Anything, "teacher-1").Return([]groups.Group{}, nil).Once()

		r := buildGroupsEngine(mockRepo, "teacher-1", "teacher", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/groups?teacher_id=teacher-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Student querying another student's groups is blocked with 403", func(t *testing.T) {
		r := buildGroupsEngine(nil, "student-1", "student", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/groups?student_id=student-other", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Student querying school-wide groups without params is blocked from school-wide list and scoped to self", func(t *testing.T) {
		mockRepo := new(mockGroupsRepo)
		mockRepo.On("ListByStudent", mock.Anything, "student-1").Return([]groups.Group{}, nil).Once()

		r := buildGroupsEngine(mockRepo, "student-1", "student", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/groups", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
