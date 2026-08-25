package student_goals

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepoForGoals struct {
	users.Repository
	mock.Mock
}

func (m *mockUserRepoForGoals) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func TestStudentGoalsSecurity_ParentGuardianship(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockGoalRepo{
		goals: []*StudentGoal{
			{ID: "g1", StudentID: "child-1", Title: "Goal 1"},
		},
	}
	userRepo := new(mockUserRepoForGoals)
	userRepo.On("IsGuardian", mock.Anything, "parent-1", "child-1").Return(true, nil)
	userRepo.On("IsGuardian", mock.Anything, "parent-1", "other-student").Return(false, nil)

	svc := NewService(repo, userRepo)
	h := NewHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "parent-1")
		c.Set("role", "parent")
		c.Next()
	})
	h.RegisterRoutes(r.Group("/api"))

	// 1. Parent accesses own child's goals -> OK
	req := httptest.NewRequest(http.MethodGet, "/api/student-goals/student/child-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Parent accesses other student's goals -> Forbidden
	req = httptest.NewRequest(http.MethodGet, "/api/student-goals/student/other-student", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
