package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/grades"
	"registro-backend/internal/lessons"
	"registro-backend/internal/schools"
	"registro-backend/internal/users"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestValidationAndEdgeCasesIntegration tests input validation and error handling across endpoints
func TestValidationAndEdgeCasesIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSchoolRepo := new(testhelpers.MockSchoolsRepository)
	mockUserRepo := new(testhelpers.MockUsersRepository)

	schoolSvc := schools.NewService(mockSchoolRepo)
	userSvc := users.NewService(mockUserRepo)

	schoolH := schools.NewHandler(schoolSvc)
	userH := users.NewHandler(userSvc)
	lessonsH := lessons.NewHandler(nil)
	gradesH := grades.NewHandler(nil, nil)

	router := gin.Default()
	api := router.Group("/api/v1")

	// Middleware setting privileged admin context
	api.Use(func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "superadmin")
		c.Set("school_id", "school-1")
		c.Next()
	})

	api.POST("/schools", schoolH.Create)
	api.POST("/users", userH.Create)
	api.POST("/lessons", lessonsH.CreateLesson)
	api.POST("/grades", gradesH.AddGrade)

	// Test 1: User Creation - Short password
	t.Run("User Creation - Short Password returns HTTP 400", func(t *testing.T) {
		reqBody := `{"email": "test@test.com", "password": "123", "first_name": "Test", "last_name": "User", "role": "teacher"}`
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	// Test 2: User Creation - Invalid email format
	t.Run("User Creation - Invalid Email returns HTTP 400", func(t *testing.T) {
		reqBody := `{"email": "not-an-email", "password": "ValidPassword123!", "first_name": "Test", "last_name": "User", "role": "teacher"}`
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	// Test 3: User Creation - Invalid role
	t.Run("User Creation - Invalid Role returns HTTP 400", func(t *testing.T) {
		reqBody := `{"email": "valid@test.com", "password": "ValidPassword123!", "first_name": "Test", "last_name": "User", "role": "superhero"}`
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	// Test 4: School Creation - Missing name
	t.Run("School Creation - Missing Name returns HTTP 400", func(t *testing.T) {
		reqBody := `{"code": "MISSINGNAME"}`
		req := httptest.NewRequest("POST", "/api/v1/schools", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	// Test 5: Lesson Creation - Missing topic
	t.Run("Lesson Creation - Missing Topic returns HTTP 400", func(t *testing.T) {
		reqBody := `{"class_id": "class-1a", "subject_id": "subj-1"}`
		req := httptest.NewRequest("POST", "/api/v1/lessons", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	// Test 6: Grade Creation - Invalid JSON
	t.Run("Grade Creation - Malformed JSON returns HTTP 400", func(t *testing.T) {
		reqBody := `{"student_id": "student-1", "grade_value":`
		req := httptest.NewRequest("POST", "/api/v1/grades", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	// Test 7: Successful School Creation mock assertion
	t.Run("School Creation - Valid Payload returns HTTP 201", func(t *testing.T) {
		mockSchoolRepo.On("Create", mock.Anything, mock.AnythingOfType("*schools.School")).Return(nil)

		reqBody := `{"name": "Scuola Test Valid", "code": "STV001", "address": "Via Test 1", "city": "Roma", "phone": "06123456", "email": "info@test.it"}`
		req := httptest.NewRequest("POST", "/api/v1/schools", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
	})
}
