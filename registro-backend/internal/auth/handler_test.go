package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"registro-backend/pkg/jwt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func setupTestHandler() (*Handler, *MockRepository, *jwt.TokenManager) {
	mockRepo := new(MockRepository)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	tokenManager := jwt.NewTokenManager(privateKey, &privateKey.PublicKey)
	mfaService := NewMFAService("Test")

	service := NewService(mockRepo, tokenManager, mfaService)
	handler := NewHandler(service)

	return handler, mockRepo, tokenManager
}

func setupAuthContext(c *gin.Context, userID, email, role string) {
	c.Set("user_id", userID)
	c.Set("email", email)
	c.Set("role", role)
}

func TestHandler_GetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockRepository)
		setupAuth      bool
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful get current user",
			userID: "user-123",
			setupMock: func(m *MockRepository) {
				schoolID := "school-1"
				m.On("GetUserByID", mock.Anything, "user-123").Return(&User{
					ID:            "user-123",
					Email:         "test@example.com",
					FirstName:     "John",
					LastName:      "Doe",
					Role:          "teacher",
					SchoolID:      &schoolID,
					EmailVerified: true,
					MFAEnabled:    false,
				}, nil)
			},
			setupAuth:      true,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "user-123", response.ID)
				assert.Equal(t, "test@example.com", response.Email)
				assert.Equal(t, "John", response.FirstName)
				assert.Equal(t, "Doe", response.LastName)
				assert.Equal(t, "teacher", response.Role)
				assert.NotNil(t, response.SchoolID)
				assert.Equal(t, "school-1", *response.SchoolID)
				assert.True(t, response.EmailVerified)
				assert.False(t, response.MFAEnabled)
			},
		},
		{
			name:           "unauthorized - no user ID in context",
			userID:         "",
			setupMock:      func(m *MockRepository) {},
			setupAuth:      false,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, "user not authenticated")
			},
		},
		{
			name:   "user not found",
			userID: "non-existent",
			setupMock: func(m *MockRepository) {
				m.On("GetUserByID", mock.Anything, "non-existent").Return(nil, ErrUserNotFound)
			},
			setupAuth:      true,
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo, _ := setupTestHandler()
			tt.setupMock(mockRepo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/auth/me", nil)

			if tt.setupAuth && tt.userID != "" {
				setupAuthContext(c, tt.userID, "test@example.com", "teacher")
			}

			handler.GetCurrentUser(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]string
		setupMock      func(*MockRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful logout",
			requestBody: map[string]string{
				"refresh_token": "valid-refresh-token",
			},
			setupMock: func(m *MockRepository) {
				m.On("GetRefreshToken", mock.Anything, "valid-refresh-token").Return(&RefreshToken{
					ID:        "token-id-123",
					UserID:    "user-123",
					Token:     "valid-refresh-token",
					ExpiresAt: time.Now().Add(24 * time.Hour),
					Revoked:   false,
				}, nil)
				// Now revokes ALL sessions for the user, not just the single token
				m.On("RevokeAllUserTokens", mock.Anything, "user-123").Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response MessageResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "logged out successfully", response.Message)
			},
		},
		{
			name:        "invalid request body - missing refresh_token returns 200 (idempotent)",
			requestBody: map[string]string{},
			setupMock: func(m *MockRepository) {
				// Token not found → no-op, logout is idempotent
				m.On("GetRefreshToken", mock.Anything, "").Return(nil, ErrInvalidToken)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response MessageResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
			},
		},
		{
			name: "token not found returns 200 (idempotent)",
			requestBody: map[string]string{
				"refresh_token": "non-existent-token",
			},
			setupMock: func(m *MockRepository) {
				m.On("GetRefreshToken", mock.Anything, "non-existent-token").Return(nil, ErrInvalidToken)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response MessageResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo, _ := setupTestHandler()
			tt.setupMock(mockRepo)

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/auth/logout", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Logout(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_Login_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.MinCost)

	tests := []struct {
		name           string
		requestBody    LoginRequest
		setupMock      func(*MockRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful login with all user fields",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			setupMock: func(m *MockRepository) {
				schoolID := "school-1"
				m.On("GetRecentLoginAttempts", mock.Anything, "test@example.com", mock.Anything, mock.Anything).Return(0, nil)
				m.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&User{
					ID:            "user-123",
					Email:         "test@example.com",
					PasswordHash:  string(passwordHash),
					FirstName:     "John",
					LastName:      "Doe",
					Role:          "teacher",
					SchoolID:      &schoolID,
					IsActive:      true,
					EmailVerified: true,
					MFAEnabled:    false,
				}, nil)
				m.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*auth.RefreshToken")).Return(nil)
				m.On("UpdateLastLogin", mock.Anything, "user-123").Return(nil)
				m.On("RecordLoginAttempt", mock.Anything, mock.AnythingOfType("*auth.LoginAttempt")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response AuthResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check user data
				assert.NotNil(t, response.User)
				assert.Equal(t, "user-123", response.User.ID)
				assert.Equal(t, "test@example.com", response.User.Email)
				assert.Equal(t, "John", response.User.FirstName)
				assert.Equal(t, "Doe", response.User.LastName)
				assert.Equal(t, "teacher", response.User.Role)
				assert.NotNil(t, response.User.SchoolID)
				assert.Equal(t, "school-1", *response.User.SchoolID)
				assert.True(t, response.User.EmailVerified)
				assert.False(t, response.User.MFAEnabled)

				// Check tokens
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.Greater(t, response.ExpiresIn, int64(0))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo, _ := setupTestHandler()
			tt.setupMock(mockRepo)

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request.RemoteAddr = "127.0.0.1:1234"

			handler.Login(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_FieldConsistency(t *testing.T) {
	t.Run("UserResponse fields match database schema", func(t *testing.T) {
		handler, mockRepo, _ := setupTestHandler()

		schoolID := "school-1"
		mockRepo.On("GetUserByID", mock.Anything, "user-123").Return(&User{
			ID:            "user-123",
			Email:         "test@example.com",
			FirstName:     "John", // snake_case in JSON
			LastName:      "Doe",  // snake_case in JSON
			Role:          "teacher",
			SchoolID:      &schoolID, // snake_case in JSON
			EmailVerified: true,      // snake_case in JSON
			MFAEnabled:    false,     // snake_case in JSON
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/auth/me", nil)
		setupAuthContext(c, "user-123", "test@example.com", "teacher")

		handler.GetCurrentUser(c)

		assert.Equal(t, http.StatusOK, w.Code)

		// Unmarshal to map to check exact JSON field names
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify snake_case field names
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "email")
		assert.Contains(t, response, "first_name")
		assert.Contains(t, response, "last_name")
		assert.Contains(t, response, "role")
		assert.Contains(t, response, "school_id")
		assert.Contains(t, response, "email_verified")
		assert.Contains(t, response, "mfa_enabled")

		mockRepo.AssertExpectations(t)
	})
}

// Add CountAll method for mock to satisfy interface
func (m *MockRepository) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
