package competencies_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/competencies"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepoForComp struct {
	mock.Mock
	users.Repository
}

func (m *mockUserRepoForComp) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func buildEngine(userRepo users.Repository, userID, role, schoolID string) *gin.Engine {
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

	svc := competencies.NewService(nil, userRepo)
	h := competencies.NewHandler(svc)
	h.RegisterRoutes(r.Group("/api"))

	return r
}

func TestCompetencies_GetClassEvaluations_RBAC(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		userID     string
		wantStatus int
	}{
		{
			name:       "Student blocked with 403 Forbidden",
			role:       "student",
			userID:     "student-1",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "Parent blocked with 403 Forbidden",
			role:       "parent",
			userID:     "parent-1",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "Teacher allowed (fails later on nil DB but passes auth with 200/empty)",
			role:       "teacher",
			userID:     "teacher-1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Admin allowed",
			role:       "admin",
			userID:     "admin-1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Coordinator allowed",
			role:       "coordinator",
			userID:     "coord-1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Unauthenticated blocked with 401 Unauthorized",
			role:       "",
			userID:     "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := buildEngine(nil, tt.userID, tt.role, "school-1")
			req := httptest.NewRequest(http.MethodGet, "/api/competencies?class_id=cls-1", nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCompetencies_GetStudentEvaluations_GuardianshipAndIDOR(t *testing.T) {
	t.Run("Student viewing another student's evaluation is blocked with 403", func(t *testing.T) {
		r := buildEngine(nil, "my-student-id", "student", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/competencies/student/other-student-id", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Parent who is NOT guardian is blocked with 403", func(t *testing.T) {
		mockU := new(mockUserRepoForComp)
		mockU.On("IsGuardian", mock.Anything, "parent-1", "student-target").Return(false, nil).Once()

		r := buildEngine(mockU, "parent-1", "parent", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/competencies/student/student-target", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockU.AssertExpectations(t)
	})

	t.Run("Parent who IS guardian is allowed", func(t *testing.T) {
		mockU := new(mockUserRepoForComp)
		mockU.On("IsGuardian", mock.Anything, "parent-1", "student-child").Return(true, nil).Once()

		r := buildEngine(mockU, "parent-1", "parent", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/competencies/student/student-child", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockU.AssertExpectations(t)
	})
}

func TestCompetencies_BatchSave_Validation(t *testing.T) {
	t.Run("Rejects invalid competency level", func(t *testing.T) {
		r := buildEngine(nil, "teacher-1", "teacher", "school-1")

		body := competencies.BatchSaveRequest{
			ClassID: "class-1",
			Evaluations: []competencies.StudentCompetencyEvaluation{
				{
					StudentID: "s-1",
					Evaluations: map[string]string{
						"DIGITAL": "INVALID_LEVEL_999",
					},
				},
			},
		}
		raw, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/competencies/batch", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
