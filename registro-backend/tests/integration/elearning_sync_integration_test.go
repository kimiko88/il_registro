package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/config"
	"registro-backend/internal/elearning"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupElearningRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	cfg := config.ElearningConfig{
		GoogleClientID:        "test-google-client-id",
		GoogleClientSecret:    "test-google-secret",
		GoogleRedirectURI:     "http://localhost:5173/auth/google/callback",
		MicrosoftClientID:     "test-ms-client-id",
		MicrosoftClientSecret: "test-ms-secret",
		MicrosoftRedirectURI:  "http://localhost:5173/auth/microsoft/callback",
		MicrosoftTenantID:     "common",
	}
	svc := elearning.NewService(cfg)
	handler := elearning.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "teacher"
		}
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "teacher-1"
		}
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_Elearning_Classroom_And_Teams_Workflow(t *testing.T) {
	r := setupElearningRouter()

	// 1. Query Providers initial connection status
	reqProviders, _ := http.NewRequest("GET", "/api/v1/elearning/providers", nil)
	wProviders := httptest.NewRecorder()
	r.ServeHTTP(wProviders, reqProviders)
	require.Equal(t, http.StatusOK, wProviders.Code)

	var initialStatus elearning.ProvidersResponse
	err := json.Unmarshal(wProviders.Body.Bytes(), &initialStatus)
	require.NoError(t, err)
	assert.True(t, initialStatus.Google.Configured)
	assert.False(t, initialStatus.Google.Connected)

	// 2. Connect Google Classroom with OAuth code
	connectReq := elearning.ConnectRequest{
		Code: "oauth2-code-google-auth-12345",
	}
	bodyConn, _ := json.Marshal(connectReq)
	reqConnect, _ := http.NewRequest("POST", "/api/v1/elearning/google/connect", bytes.NewReader(bodyConn))
	reqConnect.Header.Set("Content-Type", "application/json")
	wConnect := httptest.NewRecorder()
	r.ServeHTTP(wConnect, reqConnect)
	require.Equal(t, http.StatusOK, wConnect.Code)

	// 3. Verify Google Classroom is now connected
	wProvidersAfter := httptest.NewRecorder()
	reqProvidersAfter, _ := http.NewRequest("GET", "/api/v1/elearning/providers", nil)
	r.ServeHTTP(wProvidersAfter, reqProvidersAfter)
	require.Equal(t, http.StatusOK, wProvidersAfter.Code)

	var afterStatus elearning.ProvidersResponse
	err = json.Unmarshal(wProvidersAfter.Body.Bytes(), &afterStatus)
	require.NoError(t, err)
	assert.True(t, afterStatus.Google.Connected)

	// 4. Sync Google Classroom Courses
	reqSyncCourses, _ := http.NewRequest("POST", "/api/v1/elearning/google/sync-courses", nil)
	wSyncCourses := httptest.NewRecorder()
	r.ServeHTTP(wSyncCourses, reqSyncCourses)
	require.Equal(t, http.StatusOK, wSyncCourses.Code)

	var syncResp elearning.SyncResponse
	err = json.Unmarshal(wSyncCourses.Body.Bytes(), &syncResp)
	require.NoError(t, err)
	assert.Greater(t, syncResp.Count, 0)
	assert.Contains(t, syncResp.Message, "google")

	// 5. Sync Grades with Classroom
	gradeSyncReq := elearning.SyncGradesRequest{
		ClassID: "class-2A",
	}
	gradeBody, _ := json.Marshal(gradeSyncReq)
	reqSyncGrades, _ := http.NewRequest("POST", "/api/v1/elearning/google/sync-grades", bytes.NewReader(gradeBody))
	reqSyncGrades.Header.Set("Content-Type", "application/json")
	wSyncGrades := httptest.NewRecorder()
	r.ServeHTTP(wSyncGrades, reqSyncGrades)
	require.Equal(t, http.StatusOK, wSyncGrades.Code)

	var gradeSyncResp elearning.SyncResponse
	err = json.Unmarshal(wSyncGrades.Body.Bytes(), &gradeSyncResp)
	require.NoError(t, err)
	assert.Equal(t, 24, gradeSyncResp.Count)
}
