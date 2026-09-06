package recovery

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRecoveryRouter(svc Service, schoolID, teacherID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		if teacherID != "" {
			c.Set("teacher_id", teacherID)
			c.Set("user_id", teacherID)
		}
		c.Next()
	})
	handler := NewHandler(svc)
	handler.RegisterRoutes(r.Group("/api/v1"))
	return r
}

func TestHandler_RecoveryRoutes(t *testing.T) {
	mockRepo := new(MockRecoveryRepo)
	svc := NewService(mockRepo)
	router := setupRecoveryRouter(svc, "school-1", "teacher-1")

	// 1. CreateCourse - Success
	mockRepo.On("CreateCourse", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	mockRepo.On("GetCourseByID", mock.Anything, "course-123").Return(&RecoveryCourse{ID: "course-123", Title: "Corso Matematica"}, nil).Once()
	body, _ := json.Marshal(CreateCourseRequest{
		SubjectID:    "sub-1",
		TeacherID:    "teacher-1",
		Title:        "Corso Matematica",
		AcademicYear: "2026/2027",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/recovery/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. CreateCourse - Bad JSON
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/recovery/courses", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 3. GetCourse - Success
	mockRepo.On("GetCourseByID", mock.Anything, "course-1").Return(&RecoveryCourse{ID: "course-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/recovery/courses/course-1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. GetCourse - Not found
	mockRepo.On("GetCourseByID", mock.Anything, "course-404").Return(nil, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/recovery/courses/course-404", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 5. ListCourses - Success
	mockRepo.On("ListCourses", mock.Anything, "school-1", "2026/2027", "teacher-1").Return([]RecoveryCourse{{ID: "c-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/recovery/courses?academic_year=2026/2027&teacher_id=teacher-1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. UpdateStatus - Success
	mockRepo.On("UpdateCourseStatus", mock.Anything, "course-1", "completed").Return(nil).Once()
	statusBody, _ := json.Marshal(map[string]string{"status": "completed"})
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/recovery/courses/course-1/status", bytes.NewBuffer(statusBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. UpdateAttendance - Success
	mockRepo.On("UpdateStudentAttendance", mock.Anything, "course-1", "student-1", 4.0, "buona").Return(nil).Once()
	attBody, _ := json.Marshal(map[string]interface{}{"attendance_hours": 4.0, "notes": "buona"})
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/recovery/courses/course-1/students/student-1/attendance", bytes.NewBuffer(attBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. RecordTest - Success
	mockRepo.On("RecordRecoveryTest", mock.Anything, mock.Anything).Return(nil).Once()
	defID := "def-1"
	testBody, _ := json.Marshal(RecordTestOutcomeRequest{
		DeficiencyID: &defID,
		StudentID:    "student-1",
		SubjectID:    "sub-1",
		ClassID:      "class-1",
		TestDate:     "2026-09-01",
		TestType:     "written",
		Grade:        8.0,
	})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/recovery/tests", bytes.NewBuffer(testBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 9. ListTests - Success
	mockRepo.On("ListRecoveryTests", mock.Anything, "school-1", "class-1", "student-1").Return([]RecoveryTest{{ID: "test-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/recovery/tests?class_id=class-1&student_id=student-1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 10. ListTests - Error
	mockRepo.On("ListRecoveryTests", mock.Anything, "school-1", "class-err", "").Return(nil, errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/recovery/tests?class_id=class-err", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockRepo.AssertExpectations(t)
}
