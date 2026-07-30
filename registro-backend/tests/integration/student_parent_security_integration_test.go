package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/attendance"
	"registro-backend/internal/classes"
	"registro-backend/internal/grades"
	"registro-backend/internal/lessons"
	"registro-backend/internal/schools"
	"registro-backend/internal/subjects"
	"registro-backend/internal/users"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestStudentParentSecurityIntegration verifies that student and parent roles
// are strictly forbidden (HTTP 403) from accessing administrative and teacher-only endpoints.
func TestStudentParentSecurityIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup Mocks
	mockSchoolRepo := new(testhelpers.MockSchoolsRepository)
	mockUserRepo := new(testhelpers.MockUsersRepository)
	mockSubjectRepo := new(testhelpers.MockSubjectsRepository)

	// Setup Services
	schoolSvc := schools.NewService(mockSchoolRepo)
	userSvc := users.NewService(mockUserRepo)
	subjectSvc := subjects.NewService(mockSubjectRepo)

	// Setup Handlers
	schoolH := schools.NewHandler(schoolSvc)
	userH := users.NewHandler(userSvc)
	subjectH := subjects.NewHandler(subjectSvc)
	gradesH := grades.NewHandler(nil, nil)
	lessonsH := lessons.NewHandler(nil)
	classesH := classes.NewHandler(nil)
	attH := attendance.NewHandler(nil)

	testRoles := []string{"student", "parent"}

	for _, role := range testRoles {
		t.Run("Forbidden endpoints for role: "+role, func(t *testing.T) {
			router := gin.Default()
			api := router.Group("/api/v1")

			// Middleware for current tested role
			api.Use(func(c *gin.Context) {
				c.Set("user_id", role+"-user-id")
				c.Set("role", role)
				c.Set("school_id", "school-1")
				c.Next()
			})

			// Register endpoints
			api.POST("/schools", schoolH.Create)
			api.POST("/users", userH.Create)
			api.POST("/classes", classesH.Create)
			api.POST("/subjects", subjectH.Create)
			api.POST("/lessons", lessonsH.CreateLesson)
			api.POST("/homeworks", lessonsH.CreateHomework)
			api.POST("/grades", gradesH.AddGrade)
			api.PUT("/attendance/:id", attH.UpdateAttendance)

			// 1. POST /schools
			t.Run("POST /schools forbidden", func(t *testing.T) {
				reqBody := `{"name": "Unauthorized School", "code": "UNAUTH01"}`
				req := httptest.NewRequest("POST", "/api/v1/schools", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				assert.Equal(t, http.StatusForbidden, res.Code)
			})

			// 2. POST /users
			t.Run("POST /users forbidden", func(t *testing.T) {
				reqBody := `{"email": "unauth@test.com", "password": "SecurePassword123!", "first_name": "Bad", "last_name": "Actor", "role": "admin"}`
				req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				assert.Equal(t, http.StatusForbidden, res.Code)
			})

			// 3. POST /classes
			t.Run("POST /classes forbidden", func(t *testing.T) {
				reqBody := `{"name": "5Z", "section": "Z"}`
				req := httptest.NewRequest("POST", "/api/v1/classes", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				assert.Equal(t, http.StatusForbidden, res.Code)
			})

			// 4. POST /subjects
			t.Run("POST /subjects forbidden", func(t *testing.T) {
				reqBody := `{"name": "Unauthorized Subject", "code": "UNAUTH"}`
				req := httptest.NewRequest("POST", "/api/v1/subjects", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				assert.Equal(t, http.StatusForbidden, res.Code)
			})

			// 5. POST /lessons
			t.Run("POST /lessons forbidden", func(t *testing.T) {
				reqBody := `{"class_id": "class-1a", "subject_id": "subj-math", "topic": "Hack Lesson"}`
				req := httptest.NewRequest("POST", "/api/v1/lessons", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				assert.Equal(t, http.StatusForbidden, res.Code)
			})

			// 6. POST /homeworks
			t.Run("POST /homeworks forbidden", func(t *testing.T) {
				reqBody := `{"class_id": "class-1a", "subject_id": "subj-math", "description": "Fake Homework"}`
				req := httptest.NewRequest("POST", "/api/v1/homeworks", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				assert.Equal(t, http.StatusForbidden, res.Code)
			})

			// 7. POST /grades
			t.Run("POST /grades forbidden", func(t *testing.T) {
				reqBody := `{"student_id": "student-1", "subject_id": "subj-math", "grade_value": 10}`
				req := httptest.NewRequest("POST", "/api/v1/grades", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				assert.Equal(t, http.StatusForbidden, res.Code)
			})

			// 8. PUT /attendance/:id
			t.Run("PUT /attendance forbidden", func(t *testing.T) {
				reqBody := `{"status": "Present"}`
				req := httptest.NewRequest("PUT", "/api/v1/attendance/att-1", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				assert.Equal(t, http.StatusForbidden, res.Code)
			})
		})
	}
}
