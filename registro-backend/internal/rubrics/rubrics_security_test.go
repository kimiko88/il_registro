package rubrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRubricsSecurity_StudentAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRubricsRepo)
	repo.On("ListAssessmentsByStudent", mock.Anything, "student-1").Return([]*RubricAssessment{}, nil)
	svc := NewService(repo)
	h := NewHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		c.Next()
	})
	h.RegisterRoutes(r.Group("/api"))

	// 1. Student accessing own assessments -> OK
	req := httptest.NewRequest(http.MethodGet, "/api/rubrics/assessments/student/student-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Student accessing another student's assessments -> Forbidden
	req = httptest.NewRequest(http.MethodGet, "/api/rubrics/assessments/student/student-2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 3. Student accessing whole class assessments -> Forbidden
	req = httptest.NewRequest(http.MethodGet, "/api/rubrics/assessments/class/class-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
