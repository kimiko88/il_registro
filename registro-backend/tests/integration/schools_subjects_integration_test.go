package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/schools"
	"registro-backend/internal/subjects"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSchoolsSubjectsIntegration(t *testing.T) {
	// Setup Mocks
	mockSchoolRepo := new(testhelpers.MockSchoolsRepository)
	mockSubjectRepo := new(testhelpers.MockSubjectsRepository)

	// Setup Services
	schoolSvc := schools.NewService(mockSchoolRepo)
	subjectSvc := subjects.NewService(mockSubjectRepo)

	// Setup Handlers
	schoolHandler := schools.NewHandler(schoolSvc)
	subjectHandler := subjects.NewHandler(subjectSvc)

	// Setup Router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Register Routes
	schoolsGroup := router.Group("/schools")
	schoolsGroup.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Set("role", "superadmin")
		c.Next()
	})
	{
		schoolsGroup.POST("", func(c *gin.Context) {
			schoolHandler.Create(c)
		})
		schoolsGroup.GET("", func(c *gin.Context) {
			schoolHandler.List(c)
		})
	}

	subjectsGroup := router.Group("/subjects")
	subjectsGroup.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Set("role", "superadmin")
		c.Set("school_id", "school-1")
		c.Next()
	})
	{
		subjectsGroup.POST("", func(c *gin.Context) {
			subjectHandler.Create(c)
		})
		subjectsGroup.GET("", func(c *gin.Context) {
			subjectHandler.List(c)
		})
	}

	// Test Case 1: Create School
	t.Run("Create School", func(t *testing.T) {
		mockSchoolRepo.On("Create", mock.Anything, mock.AnythingOfType("*schools.School")).Return(nil)

		reqBody := `{"name": "Test School", "code": "TS001", "address": "123 Test St", "city": "Test City", "phone": "123456789", "email": "test@school.com"}`
		req := httptest.NewRequest("POST", "/schools", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
	})

	// Test Case 2: Create Subject for School
	t.Run("Create Subject", func(t *testing.T) {
		// Expect subject repo create called
		mockSubjectRepo.On("Create", mock.Anything, mock.AnythingOfType("*subjects.Subject")).Return(nil)

		// Note: The service doesn't implicitly validate school ID yet, but this tests the handlers wiring
		reqBody := `{"name": "Math", "code": "MATH101", "description": "Mathematics", "school_id": "school-1"}`
		req := httptest.NewRequest("POST", "/subjects", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code) // Update to StatusCreated if handler returns that
	})

	// Test Case 3: List Subjects for School
	t.Run("List Subjects", func(t *testing.T) {
		mockList := []subjects.Subject{
			{ID: "sub-1", Name: "Math", SchoolID: "school-1"},
			{ID: "sub-2", Name: "Science", SchoolID: "school-1"},
		}
		mockSubjectRepo.On("List", mock.Anything, "school-1").Return(mockList, nil)

		req := httptest.NewRequest("GET", "/subjects?school_id=school-1", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var resp []subjects.Subject
		err := json.Unmarshal(res.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, "Math", resp[0].Name)
	})
}
