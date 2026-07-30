package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/classes"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestTeacherCoordinationIntegration tests that:
// 1. A secretary/admin can assign a teacher as coordinator of 1 or more classes
// 2. The backend correctly links classes to the coordinator_id
func TestTeacherCoordinationIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClassRepo := new(testhelpers.MockClassesRepository)
	classSvc := classes.NewService(mockClassRepo)
	classH := classes.NewHandler(classSvc)

	router := gin.Default()
	api := router.Group("/api/v1")

	// Middleware setting secretary auth context
	api.Use(func(c *gin.Context) {
		c.Set("user_id", "sec-1")
		c.Set("role", "secretary")
		c.Set("school_id", "school-galileo")
		c.Next()
	})

	api.POST("/classes", classH.Create)

	t.Run("Secretary assigns a teacher as coordinator of a class", func(t *testing.T) {
		mockClassRepo.On("Create", mock.Anything, mock.AnythingOfType("*classes.Class")).Return(nil)

		reqBody := `{"name": "5A", "section": "A", "academic_year": "2025/2026", "coordinator_id": "teacher-math-uuid"}`
		req := httptest.NewRequest("POST", "/api/v1/classes", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		var createdClass classes.Class
		err := json.Unmarshal(res.Body.Bytes(), &createdClass)
		assert.NoError(t, err)
		assert.Equal(t, "teacher-math-uuid", createdClass.CoordinatorID)
	})
}
