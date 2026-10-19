package notifications

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestNotificationsRouter(svc Service, userID, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestNotificationsHandler_Endpoints(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	r := setupTestNotificationsRouter(svc, "user-123", "admin")

	// 1. GET /notifications - Unauthorized
	unauthRouter := setupTestNotificationsRouter(svc, "", "")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	w := httptest.NewRecorder()
	unauthRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// GET /notifications - Success
	mockRepo.On("ListDBNotifications", mock.Anything, "user-123", false, 50, 0).Return([]DBNotification{{ID: "n1", Title: "Avviso"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /notifications - DB error
	mockRepo.On("ListDBNotifications", mock.Anything, "user-123", true, 20, 5).Return(nil, errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/notifications?unread_only=true&limit=20&offset=5", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 2. PUT /notifications/:id/read - Success & Unauthorized
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/notifications/n1/read", nil)
	w = httptest.NewRecorder()
	unauthRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockRepo.On("MarkAsRead", mock.Anything, "user-123", "n1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/notifications/n1/read", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. PUT /notifications/read-all - Success
	mockRepo.On("MarkAllAsRead", mock.Anything, "user-123").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/notifications/read-all", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. POST /notifications/push-tokens (RegisterToken)
	mockRepo.On("SaveToken", mock.Anything, mock.Anything).Return(nil).Once()
	tokPayload := `{"device_token":"tok-abc-123","platform":"android"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/notifications/push-tokens", bytes.NewBufferString(tokPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/notifications/push-tokens", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 5. POST /notifications/register-device (RegisterDevice)
	mockRepo.On("SaveToken", mock.Anything, mock.Anything).Return(nil).Once()
	devPayload := `{"fcm_token":"fcm-xyz","platform":"ios"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/notifications/register-device", bytes.NewBufferString(devPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Invalid platform
	badDevPayload := `{"fcm_token":"fcm-xyz","platform":"windowsphone"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/notifications/register-device", bytes.NewBufferString(badDevPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 6. DELETE /notifications/push-tokens (UnregisterToken)
	mockRepo.On("DeleteToken", mock.Anything, "user-123", "tok-abc-123").Return(nil).Once()
	unregPayload := `{"device_token":"tok-abc-123"}`
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/notifications/push-tokens", bytes.NewBufferString(unregPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Missing token in delete
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/notifications/push-tokens", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 7. POST /notifications/send-push
	// Forbidden for non-admin
	teacherRouter := setupTestNotificationsRouter(svc, "teacher-1", "teacher")
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/notifications/send-push", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	teacherRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Success for admin
	mockRepo.On("GetUserTokens", mock.Anything, "student-99").Return([]PushToken{}, nil).Once()
	pushPayload := `{"user_id":"student-99","title":"Verifica","body":"Verifica domani"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/notifications/send-push", bytes.NewBufferString(pushPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. GET /manifest.json
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/manifest.json", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Registro Elettronico Scolastico")

	mockRepo.AssertExpectations(t)
}
