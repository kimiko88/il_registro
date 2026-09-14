package parents

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"registro-backend/internal/grades"
	"registro-backend/internal/users"
)

func setupParentsRouter(h *Handler, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1")
	if userID != "" {
		rg.Use(func(c *gin.Context) {
			c.Set("user_id", userID)
			c.Next()
		})
	}
	h.RegisterRoutes(rg)
	return r
}

func TestParentsHandler_GetDashboardStats(t *testing.T) {
	ctx := context.Background()

	t.Run("unauthorized when no user_id", func(t *testing.T) {
		h := NewHandler(nil)
		r := setupParentsRouter(h, "")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/dashboard/stats", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("internal server error when service fails", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		svc := NewService(nil, uRepo, nil, nil, nil)
		h := NewHandler(svc)
		r := setupParentsRouter(h, "parent-1")

		uRepo.On("GetChildren", ctx, "parent-1").Return([]users.StudentChild(nil), errors.New("db error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/dashboard/stats", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success returns 200 and stats", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		pRepo := new(mockParentRepoForTest)
		svc := NewService(pRepo, uRepo, nil, nil, nil)
		h := NewHandler(svc)
		r := setupParentsRouter(h, "parent-1")

		uRepo.On("GetChildren", ctx, "parent-1").Return([]users.StudentChild{
			{ID: "c1", UserID: "s1"},
		}, nil).Once()

		pRepo.On("GetDashboardStats", ctx, "parent-1", []string{"s1"}).Return(&ParentDashboardStatsResponse{
			ChildrenCount: 1,
		}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/dashboard/stats", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp ParentDashboardStatsResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.ChildrenCount)
	})
}

func TestParentsHandler_GetDashboard(t *testing.T) {
	ctx := context.Background()

	t.Run("unauthorized when no user_id", func(t *testing.T) {
		h := NewHandler(nil)
		r := setupParentsRouter(h, "")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/dashboard", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("internal server error when service fails", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		svc := NewService(nil, uRepo, nil, nil, nil)
		h := NewHandler(svc)
		r := setupParentsRouter(h, "parent-1")

		uRepo.On("GetChildren", ctx, "parent-1").Return([]users.StudentChild(nil), errors.New("repo fail")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/dashboard", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success returns 200 and dashboard", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		svc := NewService(nil, uRepo, nil, nil, nil)
		h := NewHandler(svc)
		r := setupParentsRouter(h, "parent-1")

		uRepo.On("GetChildren", ctx, "parent-1").Return([]users.StudentChild{}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/dashboard", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestParentsHandler_GetChildGradesAverage(t *testing.T) {
	ctx := context.Background()

	t.Run("unauthorized when no user_id", func(t *testing.T) {
		h := NewHandler(nil)
		r := setupParentsRouter(h, "")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/child/student-1/grades-average", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("forbidden when not guardian", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		svc := NewService(nil, uRepo, nil, nil, nil)
		h := NewHandler(svc)
		r := setupParentsRouter(h, "parent-1")

		uRepo.On("IsGuardian", ctx, "parent-1", "student-1").Return(false, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/child/student-1/grades-average", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("internal server error when db fails", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		gRepo := new(mockGradesRepoForParentTest)
		svc := NewService(nil, uRepo, gRepo, nil, nil)
		h := NewHandler(svc)
		r := setupParentsRouter(h, "parent-1")

		uRepo.On("IsGuardian", ctx, "parent-1", "student-1").Return(true, nil).Once()
		gRepo.On("FindByStudent", ctx, "student-1").Return(nil, errors.New("db error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/child/student-1/grades-average", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success returns average", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		gRepo := new(mockGradesRepoForParentTest)
		svc := NewService(nil, uRepo, gRepo, nil, nil)
		h := NewHandler(svc)
		r := setupParentsRouter(h, "parent-1")

		uRepo.On("IsGuardian", ctx, "parent-1", "student-1").Return(true, nil).Once()
		gRepo.On("FindByStudent", ctx, "student-1").Return([]grades.Grade{
			{GradeValue: 8.5, IsPublished: true},
		}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/parents/child/student-1/grades-average", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var body map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &body)
		assert.NoError(t, err)
		assert.Equal(t, "student-1", body["student_id"])
		assert.Equal(t, 8.5, body["average"])
	})
}
