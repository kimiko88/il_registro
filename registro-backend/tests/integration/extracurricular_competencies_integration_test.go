package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/competencies"
	"registro-backend/internal/extracurricular"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Extracurricular_CreateCourse_InvalidDates_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := extracurricular.NewService(&mockExtraRepo{})
	handler := extracurricular.NewHandler(svc)

	r := gin.New()
	r.POST("/extracurricular/courses", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.CreateCourse(c)
	})

	body := extracurricular.CreateCourseRequest{
		Title:     "Corso di Robotica",
		StartDate: "2026-06-01",
		EndDate:   "2026-05-01", // End date before start date
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/extracurricular/courses", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIntegration_Competencies_SaveEvaluation_InvalidLevel_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := competencies.NewService(nil)
	handler := competencies.NewHandler(svc)

	r := gin.New()
	r.POST("/competencies/evaluations", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.SaveEvaluation(c)
	})

	body := competencies.SaveEvaluationRequest{
		StudentID:      "student-1",
		ClassID:        "class-1",
		CompetenceCode: "COMP_1",
		Level:          "INVALID_LEVEL",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/competencies/evaluations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

type mockExtraRepo struct{}

func (m *mockExtraRepo) CreateCourse(ctx context.Context, c *extracurricular.Course) error {
	c.ID = "course-1"
	return nil
}
func (m *mockExtraRepo) GetCourseByID(ctx context.Context, id string) (*extracurricular.Course, error) {
	return &extracurricular.Course{ID: id, SchoolID: "school-1"}, nil
}
func (m *mockExtraRepo) ListCourses(ctx context.Context, schoolID, studentID string) ([]*extracurricular.Course, error) {
	return []*extracurricular.Course{}, nil
}
func (m *mockExtraRepo) EnrollStudent(ctx context.Context, courseID, studentID string) error {
	return nil
}
func (m *mockExtraRepo) ListEnrollments(ctx context.Context, courseID string) ([]*extracurricular.Enrollment, error) {
	return []*extracurricular.Enrollment{}, nil
}
func (m *mockExtraRepo) MarkAttendance(ctx context.Context, att *extracurricular.AttendanceRecord) error {
	return nil
}
func (m *mockExtraRepo) ListAttendance(ctx context.Context, courseID string, date time.Time) ([]*extracurricular.AttendanceRecord, error) {
	return []*extracurricular.AttendanceRecord{}, nil
}
