package elearning

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupElearningRouter(svc *Service, userID, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("role", role)
		}
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestElearningHandler_GetProviders(t *testing.T) {
	svc := NewService(config.ElearningConfig{
		GoogleClientID:     "g-client",
		GoogleClientSecret: "g-secret",
	})

	// Unauthorized
	rNoAuth := setupElearningRouter(svc, "", "")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/elearning/providers", nil)
	rNoAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Authorized teacher
	rTeacher := setupElearningRouter(svc, "teacher-1", "teacher")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/elearning/providers", nil)
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp ProvidersResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Google.Configured)
	assert.False(t, resp.Google.Connected)
}

func TestElearningHandler_ConnectGoogle_RoleAndValidation(t *testing.T) {
	svc := NewService(config.ElearningConfig{
		GoogleClientID:     "g-client",
		GoogleClientSecret: "g-secret",
	})

	// Forbidden for student
	rStudent := setupElearningRouter(svc, "student-1", "student")
	body, _ := json.Marshal(ConnectRequest{Code: "auth-123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/elearning/google/connect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Bad Request for missing code
	rTeacher := setupElearningRouter(svc, "teacher-1", "teacher")
	emptyBody, _ := json.Marshal(ConnectRequest{Code: ""})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/elearning/google/connect", bytes.NewBuffer(emptyBody))
	req.Header.Set("Content-Type", "application/json")
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Success
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/elearning/google/connect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestElearningHandler_SyncEndpoints(t *testing.T) {
	svc := NewService(config.ElearningConfig{})

	// Student forbidden
	rStudent := setupElearningRouter(svc, "student-1", "student")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/elearning/google/sync-courses", nil)
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Teacher allowed for sync-courses, sync-assignments, sync-grades
	rTeacher := setupElearningRouter(svc, "teacher-1", "teacher")

	// sync courses
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/elearning/google/sync-courses", nil)
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// sync assignments
	syncReq, _ := json.Marshal(SyncAssignmentsRequest{ClassID: "class-1a"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/elearning/google/sync-assignments", bytes.NewBuffer(syncReq))
	req.Header.Set("Content-Type", "application/json")
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// sync grades
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/elearning/google/sync-grades", bytes.NewBuffer(syncReq))
	req.Header.Set("Content-Type", "application/json")
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
