package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/reports"
	"registro-backend/internal/scrutiny"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupReportsRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	scrutinySvc := &scrutiny.Service{}
	reportsSvc := reports.NewService(scrutinySvc)
	handler := reports.NewHandler(reportsSvc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "secretary"
		}
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "sec-1"
		}
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_Reports_Sidi_Export_Workflow(t *testing.T) {
	r := setupReportsRouter()

	// 1. Export SIDI MIUR XML for class student registry
	reqStudents, _ := http.NewRequest("GET", "/api/v1/reports/sidi/students?class_id=class-2A", nil)
	reqStudents.Header.Set("X-Role", "secretary")
	wStudents := httptest.NewRecorder()
	r.ServeHTTP(wStudents, reqStudents)
	require.Equal(t, http.StatusOK, wStudents.Code)
	assert.Contains(t, wStudents.Header().Get("Content-Type"), "application/xml")
	assert.Contains(t, wStudents.Body.String(), "<FlussoAnagrafeSIDI>")
	assert.Contains(t, wStudents.Body.String(), "<CodiceFiscale>RSSMRA08A01H501U</CodiceFiscale>")

	// 2. Export SIDI MIUR XML for final scrutinies
	reqScrutini, _ := http.NewRequest("GET", "/api/v1/reports/sidi/scrutini?class_id=class-2A&semester=2", nil)
	reqScrutini.Header.Set("X-Role", "secretary")
	wScrutini := httptest.NewRecorder()
	r.ServeHTTP(wScrutini, reqScrutini)
	require.Equal(t, http.StatusOK, wScrutini.Code)
	assert.Contains(t, wScrutini.Header().Get("Content-Type"), "application/xml")
	assert.Contains(t, wScrutini.Body.String(), "<FlussoScrutiniSIDI>")

	// 3. Export SIDI Attendance CSV
	reqAttendance, _ := http.NewRequest("GET", "/api/v1/reports/sidi/attendance?class_id=class-2A", nil)
	reqAttendance.Header.Set("X-Role", "secretary")
	wAttendance := httptest.NewRecorder()
	r.ServeHTTP(wAttendance, reqAttendance)
	require.Equal(t, http.StatusOK, wAttendance.Code)
	assert.Contains(t, wAttendance.Header().Get("Content-Type"), "text/csv")
	assert.Contains(t, wAttendance.Body.String(), "CODICE_FISCALE")
	assert.Contains(t, wAttendance.Body.String(), "RSSMRA08A01H501U")

	// 4. Unauthorized student cannot export Excel grades matrix
	reqForbidden, _ := http.NewRequest("GET", "/api/v1/reports/grades/excel?class_id=class-2A", nil)
	reqForbidden.Header.Set("X-Role", "student")
	wForbidden := httptest.NewRecorder()
	r.ServeHTTP(wForbidden, reqForbidden)
	assert.Equal(t, http.StatusForbidden, wForbidden.Code)
}
