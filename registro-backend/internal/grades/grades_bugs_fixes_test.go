package grades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGrades_Handler_TeacherNilValidator_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil) // validator is nil

	r := gin.New()
	r.GET("/grades/my-grades/trend", func(c *gin.Context) {
		c.Set("user_id", "t1")
		c.Set("role", "teacher")
		h.GetMyTrend(c)
	})

	req := httptest.NewRequest("GET", "/grades/my-grades/trend?student_id=s2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code, "Teacher accessing another student's trend with nil validator must be rejected with 403")
}

func TestGrades_Handler_GetStudentGrades_DeprecationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil)

	mockSvc.On("GetStudentGradesWithFilter", mock.Anything, "s1", "student", "s1", mock.Anything).Return([]GradeResponse{}, nil).Once()

	r := gin.New()
	r.GET("/grades/student/:studentID", func(c *gin.Context) {
		c.Set("user_id", "s1")
		c.Set("role", "student")
		h.GetStudentGrades(c)
	})

	req := httptest.NewRequest("GET", "/grades/student/s1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Header().Get("Deprecation"))
	assert.Contains(t, w.Header().Get("Link"), "paged")
}

func TestGrades_Service_CoordinatorRoleAccess(t *testing.T) {
	s := &service{validator: nil}

	// Coordinator role should pass permission check
	err := s.checkGradeAccessPermissions(context.Background(), "coord1", "coordinator", "s1")
	assert.NoError(t, err)
}
