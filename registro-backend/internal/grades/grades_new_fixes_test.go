package grades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGrades_Handler_GetStudentAverage_TeacherAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil) // nil validator

	r := gin.New()
	r.GET("/grades/student/:studentID/average", func(c *gin.Context) {
		c.Set("user_id", "t1")
		c.Set("role", "teacher")
		h.GetStudentAverage(c)
	})

	req := httptest.NewRequest("GET", "/grades/student/s1/average", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code, "Teacher requesting student average with unconfigured validator must return 403")
}

func TestGrades_Handler_GetChildGrades_RoleForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil)

	r := gin.New()
	r.GET("/grades/child/:studentID", func(c *gin.Context) {
		c.Set("user_id", "t1")
		c.Set("role", "teacher") // teacher calling parent endpoint
		h.GetChildGrades(c)
	})

	req := httptest.NewRequest("GET", "/grades/child/s1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGrades_Service_GetMyAverages_WeightedOverall(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil)

	// Math (sub1): 2 grades (val 10, weight 1)
	// PE (sub2): 10 grades (val 6, weight 1)
	// Unweighted average of averages = (10 + 6) / 2 = 8.0
	// True weighted overall = (2*10 + 10*6) / 12 = 80 / 12 = 6.67
	var grades []Grade
	for i := 0; i < 2; i++ {
		grades = append(grades, Grade{
			ID:            "g1",
			StudentID:     "s1",
			SubjectID:     "sub1",
			GradeValue:    10.0,
			Weight:        1.0,
			Semester:      1,
			IsPublished:   true,
			GradeCategory: GradeCategorySummative,
		})
	}
	for i := 0; i < 10; i++ {
		grades = append(grades, Grade{
			ID:            "g2",
			StudentID:     "s1",
			SubjectID:     "sub2",
			GradeValue:    6.0,
			Weight:        1.0,
			Semester:      1,
			IsPublished:   true,
			GradeCategory: GradeCategorySummative,
		})
	}

	mockRepo.On("FindByStudent", "s1").Return(grades, nil).Once()

	res, err := svc.GetMyAverages(context.Background(), "s1")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 6.67, res.Semester1.OverallAverage, "Overall average must be calculated across individual grades weighted sum")
}
