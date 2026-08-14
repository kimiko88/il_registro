package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/attendance"
	"registro-backend/internal/documents"
	"registro-backend/internal/grades"
	"registro-backend/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRBACFixes_GetClassAverage_ForbiddenForStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := grades.NewHandler(nil, nil)

	r := gin.New()
	r.GET("/grades/classes/:classID/average", func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		handler.GetClassAverage(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/grades/classes/class-101/average", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRBACFixes_GetStudentSummaryForTeacher_ForbiddenForStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := attendance.NewHandler(nil)

	r := gin.New()
	r.GET("/attendance/students/:studentID/summary", func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		c.Set("school_id", "school-A")
		handler.GetStudentSummaryForTeacher(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/attendance/students/student-2/summary", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRBACFixes_GetMonthlyBreakdown_ForbiddenForOtherStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := attendance.NewHandler(nil)

	r := gin.New()
	r.GET("/attendance/students/:studentID/monthly", func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		handler.GetMonthlyBreakdown(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/attendance/students/student-999/monthly", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRBACFixes_ExportDocument_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := documents.NewHandler(nil, nil)

	r := gin.New()
	r.GET("/documents/:id/export", func(c *gin.Context) {
		c.Set("user_id", "user-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-A")
		handler.ExportDocument(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/documents/doc-1/export?format=exe", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWSFixes_SchoolIDRequiredForNonSuperadmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hub := ws.NewHub("")
	handler := ws.NewHandler(hub)

	r := gin.New()
	r.GET("/ws", func(c *gin.Context) {
		c.Set("user_id", "user-1")
		c.Set("role", "teacher")
		// school_id is NOT set
		handler.Listen(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ws", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
