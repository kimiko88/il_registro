package grades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGradesServiceAudit struct {
	mock.Mock
	Service
}

func (m *mockGradesServiceAudit) GetMyAverages(ctx context.Context, studentID string) (*StudentAveragesResponse, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*StudentAveragesResponse), args.Error(1)
}

func (m *mockGradesServiceAudit) ValidateParentGuardian(ctx context.Context, parentID, studentID string) error {
	args := m.Called(ctx, parentID, studentID)
	return args.Error(0)
}

func (m *mockGradesServiceAudit) GetMyTrend(ctx context.Context, actorID, actorRole, studentID, subjectID string) (*TrendResponse, error) {
	args := m.Called(ctx, actorID, actorRole, studentID, subjectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TrendResponse), args.Error(1)
}

func (m *mockGradesServiceAudit) CreateTestWithGrades(teacherID string, req CreateClassTestRequest) (*ClassTest, error) {
	args := m.Called(teacherID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ClassTest), args.Error(1)
}

func TestGrades_GetMyAverages_RoleEnforcement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockGradesServiceAudit)
	h := NewHandler(svc, nil)

	r := gin.New()
	r.GET("/grades/my-grades/average", func(c *gin.Context) {
		c.Set("user_id", "teacher-123")
		c.Set("role", "teacher")
		h.GetMyAverages(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/grades/my-grades/average", nil)
	r.ServeHTTP(w, req)

	// Teacher calling /my-grades/average MUST be rejected with HTTP 403
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "forbidden")
}

func TestGrades_GetStudentProfile_DirectGuardianValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockGradesServiceAudit)
	h := NewHandler(svc, nil)

	parentID := "parent-1"
	studentID := "student-2"

	// ValidateParentGuardian returns ErrNotGuardian
	svc.On("ValidateParentGuardian", mock.Anything, parentID, studentID).Return(ErrNotGuardian).Once()

	r := gin.New()
	r.GET("/analytics/student/:studentID/profile", func(c *gin.Context) {
		c.Set("user_id", parentID)
		c.Set("role", "parent")
		h.GetStudentProfile(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/analytics/student/"+studentID+"/profile", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	svc.AssertExpectations(t)
}

func TestGrades_CreateTestWithGrades_ResponseBodyIncludesTestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockGradesServiceAudit)
	h := NewHandler(svc, nil)

	teacherID := "teacher-1"
	svc.On("CreateTestWithGrades", teacherID, mock.Anything).Return(&ClassTest{
		ID:      "test-uuid-999",
		ClassID: "class-1",
	}, nil).Once()

	r := gin.New()
	r.POST("/tests", func(c *gin.Context) {
		c.Set("user_id", teacherID)
		c.Set("role", "teacher")
		h.CreateTestWithGrades(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tests", strings.NewReader(`{
		"class_id": "class-1",
		"subject_id": "sub-1",
		"title": "Verifica Matematica",
		"date": "2026-08-15",
		"evaluation_type": "Written"
	}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "test-uuid-999")
	svc.AssertExpectations(t)
}
