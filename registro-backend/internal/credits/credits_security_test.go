package credits_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/credits"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreditsRepo struct {
	mock.Mock
}

func (m *mockCreditsRepo) SaveCredit(ctx context.Context, credit *credits.StudentSchoolCredit) error {
	args := m.Called(ctx, credit)
	return args.Error(0)
}

func (m *mockCreditsRepo) GetCreditByStudentAndYear(ctx context.Context, studentID, academicYear string, gradeLevel int) (*credits.StudentSchoolCredit, error) {
	args := m.Called(ctx, studentID, academicYear, gradeLevel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*credits.StudentSchoolCredit), args.Error(1)
}

func (m *mockCreditsRepo) ListCreditsByClass(ctx context.Context, classID, academicYear string) ([]credits.StudentSchoolCredit, error) {

	args := m.Called(ctx, classID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]credits.StudentSchoolCredit), args.Error(1)
}

func (m *mockCreditsRepo) GetStudentCreditSummary(ctx context.Context, studentID string) (*credits.StudentCreditSummary, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*credits.StudentCreditSummary), args.Error(1)
}

type mockUserRepoForCredits struct {
	mock.Mock
	users.Repository
}

func (m *mockUserRepoForCredits) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func buildCreditsEngine(repo credits.Repository, userRepo users.Repository, userID, role, schoolID string) *gin.Engine {
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

	svc := credits.NewService(repo, userRepo)
	h := credits.NewHandler(svc)
	h.RegisterRoutes(r.Group("/api"))

	return r
}

func TestCredits_Calculate_Auth(t *testing.T) {
	t.Run("Unauthenticated is blocked with 401", func(t *testing.T) {
		r := buildCreditsEngine(nil, nil, "", "", "")
		req := httptest.NewRequest(http.MethodGet, "/api/credits/calculate?grade_level=3&average=8.5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Authenticated user is allowed", func(t *testing.T) {
		r := buildCreditsEngine(nil, nil, "user-1", "student", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/credits/calculate?grade_level=3&average=8.5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCredits_Assign_RBAC(t *testing.T) {
	body := credits.AssignCreditRequest{
		StudentID:      "student-1",
		ClassID:        "class-1",
		AcademicYear:   "2025/2026",
		GradeLevel:     4,
		GradeAverage:   8.0,
		ConductGrade:   9,
		AssignedCredit: 11,
	}
	raw, _ := json.Marshal(body)

	t.Run("Student blocked from assigning credits with 403", func(t *testing.T) {
		r := buildCreditsEngine(nil, nil, "student-1", "student", "school-1")
		req := httptest.NewRequest(http.MethodPost, "/api/credits/assign", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Parent blocked from assigning credits with 403", func(t *testing.T) {
		r := buildCreditsEngine(nil, nil, "parent-1", "parent", "school-1")
		req := httptest.NewRequest(http.MethodPost, "/api/credits/assign", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Teacher is allowed to assign credits", func(t *testing.T) {
		mockRepo := new(mockCreditsRepo)
		mockRepo.On("SaveCredit", mock.Anything, mock.Anything).Return(nil).Once()

		r := buildCreditsEngine(mockRepo, nil, "teacher-1", "teacher", "school-1")
		req := httptest.NewRequest(http.MethodPost, "/api/credits/assign", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestCredits_ListByClass_RBAC(t *testing.T) {
	t.Run("Student blocked with 403", func(t *testing.T) {
		r := buildCreditsEngine(nil, nil, "student-1", "student", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/credits/class/cls-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Teacher allowed", func(t *testing.T) {
		mockRepo := new(mockCreditsRepo)
		mockRepo.On("ListCreditsByClass", mock.Anything, "cls-1", "").Return([]credits.StudentSchoolCredit{}, nil).Once()

		r := buildCreditsEngine(mockRepo, nil, "teacher-1", "teacher", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/credits/class/cls-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestCredits_GetStudentSummary_IDOR(t *testing.T) {
	t.Run("Student blocked from viewing another student's summary with 403", func(t *testing.T) {
		r := buildCreditsEngine(nil, nil, "student-1", "student", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/credits/student/other-student/summary", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Student allowed to view their own summary", func(t *testing.T) {
		mockRepo := new(mockCreditsRepo)
		mockRepo.On("GetStudentCreditSummary", mock.Anything, "student-1").Return(&credits.StudentCreditSummary{StudentID: "student-1"}, nil).Once()

		r := buildCreditsEngine(mockRepo, nil, "student-1", "student", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/credits/student/student-1/summary", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Parent who is NOT guardian is blocked with 403", func(t *testing.T) {
		mockUser := new(mockUserRepoForCredits)
		mockUser.On("IsGuardian", mock.Anything, "parent-1", "target-student").Return(false, nil).Once()

		r := buildCreditsEngine(nil, mockUser, "parent-1", "parent", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/credits/student/target-student/summary", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
		mockUser.AssertExpectations(t)
	})

	t.Run("Parent who IS guardian is allowed", func(t *testing.T) {
		mockRepo := new(mockCreditsRepo)
		mockUser := new(mockUserRepoForCredits)

		mockUser.On("IsGuardian", mock.Anything, "parent-1", "child-student").Return(true, nil).Once()
		mockRepo.On("GetStudentCreditSummary", mock.Anything, "child-student").Return(&credits.StudentCreditSummary{StudentID: "child-student"}, nil).Once()

		r := buildCreditsEngine(mockRepo, mockUser, "parent-1", "parent", "school-1")
		req := httptest.NewRequest(http.MethodGet, "/api/credits/student/child-student/summary", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
		mockUser.AssertExpectations(t)
	})
}
