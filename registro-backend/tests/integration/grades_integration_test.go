package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/grades"
	"registro-backend/internal/users"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGradesIntegration_GetMyGrades(t *testing.T) {
	// Setup
	mockRepo := new(testhelpers.MockGradesRepository)
	mockUserRepo := new(testhelpers.MockUsersRepository)
	mockAnalytics := new(testhelpers.MockAnalyticsService)

	// NewService(r Repository, ur users.Repository, db *sql.DB, b EventBroadcaster)
	service := grades.NewService(mockRepo, mockUserRepo, nil, nil)

	// NewHandler(s Service, a AnalyticsService)
	handler := grades.NewHandler(service, mockAnalytics)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/grades/my", func(c *gin.Context) {
		c.Set("user_id", "student1")
		c.Set("role", "student")
		handler.GetMyGrades(c)
	})

	t.Run("success", func(t *testing.T) {
		// Mock Data
		repoGrades := []grades.Grade{
			{ID: "g1", StudentID: "student1", GradeValue: 8.5, SubjectID: "math", IsPublished: true, Date: time.Now()},
		}

		mockRepo.On("FindByStudent", "student1").Return(repoGrades, nil)

		req := httptest.NewRequest("GET", "/grades/my", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp, "semesters")
	})
}

func TestGradesIntegration_GetClassGrades(t *testing.T) {
	mockRepo := new(testhelpers.MockGradesRepository)
	mockUserRepo := new(testhelpers.MockUsersRepository)
	mockAnalytics := new(testhelpers.MockAnalyticsService)

	service := grades.NewService(mockRepo, mockUserRepo, nil, nil)
	handler := grades.NewHandler(service, mockAnalytics)

	router := gin.Default()
	router.GET("/grades/class/:classID", func(c *gin.Context) {
		c.Set("user_id", "admin1")
		c.Set("role", "admin")
		handler.GetClassGrades(c)
	})

	t.Run("teacher success", func(t *testing.T) {
		// Mock find by class
		repoGrades := []grades.Grade{
			{ID: "g2", StudentID: "s2", SubjectID: "math", IsPublished: true}, // Removed ClassID
		}
		// FindByClassAndSubject signature: (classID, subjectID string, semester int)
		// handler calls with subjectID="" if not provided in query.
		mockRepo.On("FindByClassAndSubject", "classA", "", 0).Return(repoGrades, nil)

		// Mock student lookup for names
		mockUserRepo.On("GetStudentsByClass", mock.Anything, "classA").Return([]users.User{}, nil)

		req := httptest.NewRequest("GET", "/grades/class/classA", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}
