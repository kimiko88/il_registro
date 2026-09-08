package extracurricular

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupExtracurricularRouter(svc *Service, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if role != "" {
			c.Set("role", role)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	handler := NewHandler(svc)
	handler.RegisterRoutes(r.Group("/api/v1"))
	return r
}

func TestHandler_Extracurricular(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	// 1. CreateCourse - Forbidden for student
	rStudent := setupExtracurricularRouter(svc, "s-1", "student", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/extracurricular/courses", bytes.NewBufferString(`{}`))
	w := httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 2. CreateCourse - Success for teacher
	rTeacher := setupExtracurricularRouter(svc, "t-1", "teacher", "school-1")
	mockRepo.On("CreateCourse", mock.Anything, mock.Anything).Return(nil).Once()
	body, _ := json.Marshal(CreateCourseRequest{
		Title:           "Teatro",
		StartDate:       "2026-10-01",
		EndDate:         "2026-11-01",
		MaxParticipants: 20,
	})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/extracurricular/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 3. ListCourses - Success
	mockRepo.On("ListCourses", mock.Anything, "school-1", "t-1").Return([]*Course{{ID: "c-1", Title: "Teatro"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/extracurricular/courses", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. EnrollStudent - Success
	course := &Course{ID: "c-1", SchoolID: "school-1", MaxParticipants: 10, EnrolledCount: 2}
	mockRepo.On("GetCourseByID", mock.Anything, "c-1").Return(course, nil).Once()
	mockRepo.On("EnrollStudent", mock.Anything, "c-1", "s-1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/extracurricular/courses/c-1/enroll", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. ListEnrollments - Success
	mockRepo.On("GetCourseByID", mock.Anything, "c-1").Return(course, nil).Once()
	mockRepo.On("ListEnrollments", mock.Anything, "c-1").Return([]*Enrollment{{ID: "e-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/extracurricular/courses/c-1/enrollments", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. MarkAttendance - Success
	mockRepo.On("GetCourseByID", mock.Anything, "c-1").Return(&Course{ID: "c-1", TeacherID: "t-1"}, nil).Once()
	mockRepo.On("MarkAttendance", mock.Anything, mock.Anything).Return(nil).Once()
	attBody, _ := json.Marshal(MarkAttendanceRequest{
		CourseID:  "c-1",
		StudentID: "s-1",
		Date:      time.Now().Format("2006-01-02"),
		Status:    "present",
		Hours:     2.0,
	})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/extracurricular/attendance", bytes.NewBuffer(attBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockRepo.AssertExpectations(t)
}
