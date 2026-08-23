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

func TestCalculator_DetectOutliers_UnclampedThresholds(t *testing.T) {
	calc := NewCalculator()

	// High mean cluster with low outlier
	// Mean around 8.8, stdDev around 2.4 -> low bound around 4.0, high bound around 13.6
	// Grade 2.0 should be detected as outlier (< low bound)
	grades := []Grade{
		{ID: "g1", GradeValue: 9.0, GradeType: GradeTypeNumeric},
		{ID: "g2", GradeValue: 9.5, GradeType: GradeTypeNumeric},
		{ID: "g3", GradeValue: 9.0, GradeType: GradeTypeNumeric},
		{ID: "g4", GradeValue: 8.5, GradeType: GradeTypeNumeric},
		{ID: "g5", GradeValue: 9.0, GradeType: GradeTypeNumeric},
		{ID: "g6", GradeValue: 10.0, GradeType: GradeTypeNumeric},
		{ID: "g-outlier", GradeValue: 2.0, GradeType: GradeTypeNumeric},
	}

	outliers := calc.DetectOutliers(grades)
	assert.Contains(t, outliers, "g-outlier")
	assert.Len(t, outliers, 1)
}

func TestCalculator_CalculatePercentile_InvalidScore(t *testing.T) {
	calc := NewCalculator()
	grades := []Grade{
		{ID: "g1", GradeValue: 6.0, GradeType: GradeTypeNumeric},
		{ID: "g2", GradeValue: 8.0, GradeType: GradeTypeNumeric},
	}

	assert.Equal(t, 0.0, calc.CalculatePercentile(grades, 0.0))
	assert.Equal(t, 0.0, calc.CalculatePercentile(grades, -1.0))
	assert.Equal(t, 0.0, calc.CalculatePercentile(grades, 11.0))
	assert.InDelta(t, 50.0, calc.CalculatePercentile(grades, 7.0), 0.01)
}
