package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/users"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupUserAssignmentsRouter creates an isolated test router for user assignment endpoints.
func setupUserAssignmentsRouter(mockRepo *testhelpers.MockUsersRepository) (*gin.Engine, *string, *string) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userSvc := users.NewService(mockRepo)
	userH := users.NewHandler(userSvc)

	var currentRole string
	var currentUserID string
	schoolID := "school-test-123"

	authMiddleware := func(c *gin.Context) {
		c.Set("role", currentRole)
		c.Set("user_id", currentUserID)
		c.Set("school_id", schoolID)
		c.Next()
	}

	api := router.Group("/api/v1")
	api.Use(authMiddleware)
	{
		api.GET("/users/:id/assignments", userH.GetAssignments)
		api.POST("/users/:id/assignments", userH.AddAssignment)
		api.DELETE("/users/:id/assignments/:assignmentId", userH.DeleteAssignment)
		api.PUT("/users/:id/coordinated-classes", userH.SetCoordinatedClasses)
	}

	return router, &currentRole, &currentUserID
}

func TestUserAssignmentsLifecycle_HTTP(t *testing.T) {
	mockRepo := new(testhelpers.MockUsersRepository)
	router, currentRole, currentUserID := setupUserAssignmentsRouter(mockRepo)

	targetUserID := "teacher-uuid-456"
	targetSchoolID := "school-test-123"
	targetTeacher := users.User{
		ID:        targetUserID,
		Email:     "docente.test@scuola.it",
		Role:      "teacher",
		SchoolID:  &targetSchoolID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("1. GET /users/:id/assignments returns active duties", func(t *testing.T) {
		*currentRole = "principal"
		*currentUserID = "principal-1"

		existingAssignments := []users.UserAssignment{
			{
				ID:             "asgn-1",
				SchoolID:       targetSchoolID,
				UserID:         targetUserID,
				AssignmentType: "referente_inclusione",
				ScopeType:      "school",
				Title:          "Referente Inclusione, BES e DSA",
				IsActive:       true,
			},
		}

		mockRepo.On("GetAssignments", mock.Anything, targetUserID).Return(existingAssignments, nil).Once()

		req := httptest.NewRequest("GET", "/api/v1/users/"+targetUserID+"/assignments", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var resp []users.UserAssignment
		err := json.Unmarshal(res.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "referente_inclusione", resp[0].AssignmentType)
	})

	t.Run("2. POST /users/:id/assignments: Dirigente can assign educational duty", func(t *testing.T) {
		*currentRole = "principal"
		*currentUserID = "principal-1"

		mockRepo.On("GetByID", mock.Anything, targetUserID).Return(&targetTeacher, nil).Once()
		mockRepo.On("CreateAssignment", mock.Anything, mock.MatchedBy(func(a *users.UserAssignment) bool {
			return a.UserID == targetUserID && a.AssignmentType == "animatore_digitale"
		})).Return(nil).Once()

		payload := `{"assignment_type": "animatore_digitale", "scope_type": "school", "title": "Animatore Digitale PNSD"}`
		req := httptest.NewRequest("POST", "/api/v1/users/"+targetUserID+"/assignments", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
	})

	t.Run("3. POST /users/:id/assignments: DSGA can assign ATA service duty", func(t *testing.T) {
		*currentRole = "dsga"
		*currentUserID = "dsga-1"

		mockRepo.On("GetByID", mock.Anything, targetUserID).Return(&targetTeacher, nil).Once()
		mockRepo.On("CreateAssignment", mock.Anything, mock.MatchedBy(func(a *users.UserAssignment) bool {
			return a.UserID == targetUserID && a.AssignmentType == "responsabile_servizio"
		})).Return(nil).Once()

		payload := `{"assignment_type": "responsabile_servizio", "scope_type": "service", "scope_id": "Lab 1", "title": "Resp. Lab Informatica"}`
		req := httptest.NewRequest("POST", "/api/v1/users/"+targetUserID+"/assignments", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
	})

	t.Run("4. POST /users/:id/assignments: DSGA CANNOT assign educational duties (403)", func(t *testing.T) {
		*currentRole = "dsga"
		*currentUserID = "dsga-1"

		payload := `{"assignment_type": "coordinatore_classe", "scope_type": "class", "scope_id": "class-2a", "title": "Coordinatore 2A"}`
		req := httptest.NewRequest("POST", "/api/v1/users/"+targetUserID+"/assignments", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("5. POST /users/:id/assignments: Teacher CANNOT assign any duty (403)", func(t *testing.T) {
		*currentRole = "teacher"
		*currentUserID = "teacher-other"

		payload := `{"assignment_type": "coordinatore_classe", "scope_type": "class", "scope_id": "class-2a"}`
		req := httptest.NewRequest("POST", "/api/v1/users/"+targetUserID+"/assignments", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("6. PUT /users/:id/coordinated-classes: Segreteria can set multi-class coordination", func(t *testing.T) {
		*currentRole = "secretary"
		*currentUserID = "sec-1"

		mockRepo.On("GetByID", mock.Anything, targetUserID).Return(&targetTeacher, nil).Once()
		mockRepo.On("SetCoordinatedClasses", mock.Anything, targetSchoolID, targetUserID, []string{"class-2a", "class-2b"}).Return(nil).Once()

		payload := `{"class_ids": ["class-2a", "class-2b"]}`
		req := httptest.NewRequest("PUT", "/api/v1/users/"+targetUserID+"/coordinated-classes", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("7. PUT /users/:id/coordinated-classes: Teacher cannot update coordination (403)", func(t *testing.T) {
		*currentRole = "teacher"
		*currentUserID = "teacher-uuid-456"

		payload := `{"class_ids": ["class-2a"]}`
		req := httptest.NewRequest("PUT", "/api/v1/users/"+targetUserID+"/coordinated-classes", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("8. DELETE /users/:id/assignments/:assignmentId: Dirigente can revoke duty", func(t *testing.T) {
		*currentRole = "principal"
		*currentUserID = "principal-1"

		assignmentID := "asgn-to-delete"
		existingAssignments := []users.UserAssignment{
			{
				ID:             assignmentID,
				SchoolID:       targetSchoolID,
				UserID:         targetUserID,
				AssignmentType: "referente_progetto",
			},
		}

		mockRepo.On("GetByID", mock.Anything, targetUserID).Return(&targetTeacher, nil).Once()
		mockRepo.On("GetAssignments", mock.Anything, targetUserID).Return(existingAssignments, nil).Once()
		mockRepo.On("DeleteAssignment", mock.Anything, assignmentID).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/api/v1/users/"+targetUserID+"/assignments/"+assignmentID, nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("9. DELETE /users/:id/assignments/:assignmentId: Student cannot revoke duty (403)", func(t *testing.T) {
		*currentRole = "student"
		*currentUserID = "student-1"

		assignmentID := "asgn-to-delete"
		existingAssignments := []users.UserAssignment{
			{
				ID:             assignmentID,
				SchoolID:       targetSchoolID,
				UserID:         targetUserID,
				AssignmentType: "referente_progetto",
			},
		}

		mockRepo.On("GetByID", mock.Anything, targetUserID).Return(&targetTeacher, nil).Once()
		mockRepo.On("GetAssignments", mock.Anything, targetUserID).Return(existingAssignments, nil).Once()

		req := httptest.NewRequest("DELETE", "/api/v1/users/"+targetUserID+"/assignments/"+assignmentID, nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusForbidden, res.Code)
	})
}
