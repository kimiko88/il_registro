package grades

import (
	"context"
	"io"
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

func (m *mockGradesServiceAudit) CreateTestWithGrades(ctx context.Context, teacherID string, req CreateClassTestRequest) (*ClassTest, error) {
	args := m.Called(teacherID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ClassTest), args.Error(1)
}

func (m *mockGradesServiceAudit) GenerateSemesterReportPDF(ctx context.Context, actorID, actorRole, studentID string, semester int) ([]byte, error) {
	args := m.Called(ctx, actorID, actorRole, studentID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockGradesServiceAudit) BulkImport(teacherID, schoolID string, r io.Reader, semester int) (*ImportResult, error) {
	args := m.Called(teacherID, schoolID, r, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ImportResult), args.Error(1)
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

func TestGrades_DownloadSemesterReportPDF_ParentGuardianCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockGradesServiceAudit)
	h := NewHandler(svc, nil)

	parentID := "parent-1"
	targetStudentID := "student-99"

	svc.On("ValidateParentGuardian", mock.Anything, parentID, targetStudentID).Return(ErrNotGuardian).Once()

	r := gin.New()
	r.GET("/my-grades/semester/:semester/pdf", func(c *gin.Context) {
		c.Set("user_id", parentID)
		c.Set("role", "parent")
		h.DownloadSemesterReportPDF(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/my-grades/semester/1/pdf?student_id="+targetStudentID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "forbidden: not a guardian of this student")
	svc.AssertExpectations(t)
}

func TestGrades_GetMyTrend_ParentGuardianCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockGradesServiceAudit)
	h := NewHandler(svc, nil)

	parentID := "parent-1"
	targetStudentID := "student-99"

	svc.On("ValidateParentGuardian", mock.Anything, parentID, targetStudentID).Return(ErrNotGuardian).Once()

	r := gin.New()
	r.GET("/my-grades/trend", func(c *gin.Context) {
		c.Set("user_id", parentID)
		c.Set("role", "parent")
		h.GetMyTrend(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/my-grades/trend?student_id="+targetStudentID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "forbidden: not a guardian of this student")
	svc.AssertExpectations(t)
}

func TestGrades_InternalError_MaskedFromClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockGradesServiceAudit)
	h := NewHandler(svc, nil)

	studentID := "student-1"
	// Return a database/internal error
	svc.On("GetMyAverages", mock.Anything, studentID).Return(nil, assert.AnError).Once()

	r := gin.New()
	r.GET("/my-grades/average", func(c *gin.Context) {
		c.Set("user_id", studentID)
		c.Set("role", "student")
		h.GetMyAverages(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/my-grades/average", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// Internal error must be generic, not leaking db error string
	assert.Equal(t, `{"error":"internal server error"}`, w.Body.String())
	svc.AssertExpectations(t)
}
