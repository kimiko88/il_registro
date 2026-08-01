package unit

import (
	"context"
	"errors"
	"testing"

	"registro-backend/internal/auth"
	"registro-backend/pkg/jwt"
	"registro-backend/tests/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(testhelpers.MockAuthRepository)

	priv, pub := testhelpers.GenerateRSAKeys()
	tokenManager := jwt.NewTokenManager(priv, pub)

	mfaService := auth.NewMFAService("issuer")
	service := auth.NewService(mockRepo, tokenManager, mfaService, nil)

	tests := []struct {
		name          string
		input         auth.RegisterRequest
		setupMock     func()
		expectedError error
	}{
		{
			name: "success",
			input: auth.RegisterRequest{
				Email:     "new@test.com",
				Password:  "Password123!",
				FirstName: "Mario",
				LastName:  "Rossi",
				Role:      "student",
			},
			setupMock: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "new@test.com").Return(nil, errors.New("user not found"))
				mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *auth.User) bool {
					return u.Email == "new@test.com"
				})).Return(nil)
				mockRepo.On("AddPasswordHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "email already exists",
			input: auth.RegisterRequest{
				Email:     "existing@test.com",
				Password:  "Password123!",
				FirstName: "Existing",
				LastName:  "User",
				Role:      "student",
			},
			setupMock: func() {
				// Validation passes first, then checks email
				existingUser := &auth.User{ID: "1", Email: "existing@test.com"}
				// We expect GetUserByEmail to be called and RETURN existing user (nil error means found)
				mockRepo.On("GetUserByEmail", mock.Anything, "existing@test.com").Return(existingUser, nil)
			},
			expectedError: auth.ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo = new(testhelpers.MockAuthRepository)
			service = auth.NewService(mockRepo, tokenManager, mfaService, nil)
			tt.setupMock()

			_, err := service.Register(context.Background(), &tt.input)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(testhelpers.MockAuthRepository)
	priv, pub := testhelpers.GenerateRSAKeys()
	tokenManager := jwt.NewTokenManager(priv, pub)

	mfaService := auth.NewMFAService("issuer")
	service := auth.NewService(mockRepo, tokenManager, mfaService, nil)

	tests := []struct {
		name          string
		input         auth.LoginRequest
		setupMock     func()
		expectedError error
	}{
		{
			name:  "user not found",
			input: auth.LoginRequest{Email: "unknown@test.com", Password: "pass"},
			setupMock: func() {
				mockRepo.On("GetRecentLoginAttempts", mock.Anything, "unknown@test.com", mock.Anything, mock.Anything).Return(0, nil)
				mockRepo.On("GetRecentLoginAttemptsByEmail", mock.Anything, "unknown@test.com", mock.Anything).Return(0, nil)
				mockRepo.On("GetUserByEmail", mock.Anything, "unknown@test.com").Return(nil, errors.New("not found"))
				mockRepo.On("RecordLoginAttempt", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: auth.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo = new(testhelpers.MockAuthRepository)
			service = auth.NewService(mockRepo, tokenManager, mfaService, nil)
			tt.setupMock()

			_, err := service.Login(context.Background(), &tt.input, "127.0.0.1", "agent")
			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
