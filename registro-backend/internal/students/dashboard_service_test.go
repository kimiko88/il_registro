package students

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStudentDashboardService_NilDB(t *testing.T) {
	svc := NewDashboardService(nil)
	stats, err := svc.GetStats(context.Background(), "user-123")
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 0.0, stats.AverageGrade)
	assert.Equal(t, 100.0, stats.AttendanceRate)
	assert.Equal(t, 0, stats.HomeworkCount)
	assert.Equal(t, 0, stats.DocumentsCount)
	assert.Equal(t, 0, stats.TotalGrades)
	assert.Equal(t, 100.0, stats.PresenceRate)
	assert.Equal(t, 0, stats.UpcomingTests)
}

func TestStudentDashboardHandler_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := NewDashboardService(nil)
	h := NewDashboardHandler(svc)

	r := gin.New()
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/students/dashboard/stats", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestStudentDashboardHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := NewDashboardService(nil)
	h := NewDashboardHandler(svc)

	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", "student-user-1")
		c.Next()
	})
	h.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/students/dashboard/stats", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "attendance_rate")
	assert.Contains(t, w.Body.String(), "presence_rate")
	assert.Contains(t, w.Body.String(), "upcoming_tests")
}
