package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"registro-backend/pkg/jwt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) EnableMFA(ctx context.Context, userID, secret string) error {
	args := m.Called(ctx, userID, secret)
	return args.Error(0)
}

func (m *MockRepository) DisableMFA(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) GetMFASecret(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) CreateRecoveryCodes(ctx context.Context, userID string, codes []string) error {
	args := m.Called(ctx, userID, codes)
	return args.Error(0)
}

func (m *MockRepository) GetRecoveryCodes(ctx context.Context, userID string) ([]*MFARecoveryCode, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MFARecoveryCode), args.Error(1)
}

func (m *MockRepository) UseRecoveryCode(ctx context.Context, userID, code string) error {
	args := m.Called(ctx, userID, code)
	return args.Error(0)
}

func (m *MockRepository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRepository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RefreshToken), args.Error(1)
}

func (m *MockRepository) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}

func (m *MockRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRepository) GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PasswordResetToken), args.Error(1)
}

func (m *MockRepository) UsePasswordResetToken(ctx context.Context, tokenID string) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}

func (m *MockRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}

func (m *MockRepository) RecordLoginAttempt(ctx context.Context, attempt *LoginAttempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockRepository) GetRecentLoginAttempts(ctx context.Context, email, ipAddress string, since time.Time) (int, error) {
	args := m.Called(ctx, email, ipAddress, since)
	return args.Int(0), args.Error(1)
}

func TestService_Register(t *testing.T) {
	tests := []struct {
		name    string
		req     *RegisterRequest
		mockFn  func(*MockRepository)
		wantErr error
	}{
		{
			name: "successful registration",
			req: &RegisterRequest{
				Email:     "test@example.com",
				Password:  "Password123!",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "student",
			},
			mockFn: func(m *MockRepository) {
				m.On("GetUserByEmail", mock.Anything, "test@example.com").Return(nil, ErrUserNotFound)
				m.On("CreateUser", mock.Anything, mock.AnythingOfType("*auth.User")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "email already exists",
			req: &RegisterRequest{
				Email:     "existing@example.com",
				Password:  "Password123!",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "student",
			},
			mockFn: func(m *MockRepository) {
				m.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(&User{}, nil)
			},
			wantErr: ErrEmailAlreadyExists,
		},
		{
			name: "weak password",
			req: &RegisterRequest{
				Email:     "test@example.com",
				Password:  "weak",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "student",
			},
			mockFn:  func(m *MockRepository) {},
			wantErr: ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := &Service{
				repo:       mockRepo,
				bcryptCost: bcrypt.MinCost, // Use minimum cost for faster tests
			}

			user, err := service.Register(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.req.Email, user.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Login(t *testing.T) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.MinCost)

	tests := []struct {
		name    string
		req     *LoginRequest
		mockFn  func(*MockRepository)
		wantErr error
	}{
		{
			name: "successful login",
			req: &LoginRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			mockFn: func(m *MockRepository) {
				m.On("GetRecentLoginAttempts", mock.Anything, "test@example.com", "127.0.0.1", mock.Anything).Return(0, nil)
				m.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: string(passwordHash),
					IsActive:     true,
					MFAEnabled:   false,
				}, nil)
				m.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*auth.RefreshToken")).Return(nil)
				m.On("UpdateLastLogin", mock.Anything, "user-123").Return(nil)
				m.On("RecordLoginAttempt", mock.Anything, mock.AnythingOfType("*auth.LoginAttempt")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid credentials",
			req: &LoginRequest{
				Email:    "test@example.com",
				Password: "WrongPassword",
			},
			mockFn: func(m *MockRepository) {
				m.On("GetRecentLoginAttempts", mock.Anything, "test@example.com", "127.0.0.1", mock.Anything).Return(0, nil)
				m.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&User{
					PasswordHash: string(passwordHash),
				}, nil)
				m.On("RecordLoginAttempt", mock.Anything, mock.AnythingOfType("*auth.LoginAttempt")).Return(nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "too many attempts",
			req: &LoginRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			mockFn: func(m *MockRepository) {
				m.On("GetRecentLoginAttempts", mock.Anything, "test@example.com", "127.0.0.1", mock.Anything).Return(5, nil)
			},
			wantErr: ErrTooManyAttempts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			// Generate keys for token manager
			privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
			tokenManager := jwt.NewTokenManager(privateKey, &privateKey.PublicKey)

			service := &Service{
				repo:         mockRepo,
				tokenManager: tokenManager,
				mfaService:   NewMFAService("Test"),
				bcryptCost:   bcrypt.MinCost,
			}

			// Note: This test would need proper JWT mock setup in production
			_, err := service.Login(context.Background(), tt.req, "127.0.0.1", "test-agent")

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_SetupMFA(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		mockFn  func(*MockRepository)
		wantErr error
	}{
		{
			name:   "successful MFA setup",
			userID: "user-123",
			mockFn: func(m *MockRepository) {
				m.On("GetUserByID", mock.Anything, "user-123").Return(&User{
					ID:         "user-123",
					Email:      "test@example.com",
					MFAEnabled: false,
				}, nil)
				m.On("CreateRecoveryCodes", mock.Anything, "user-123", mock.AnythingOfType("[]string")).Return(nil)
				m.On("EnableMFA", mock.Anything, "user-123", mock.AnythingOfType("string")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "MFA already enabled",
			userID: "user-123",
			mockFn: func(m *MockRepository) {
				m.On("GetUserByID", mock.Anything, "user-123").Return(&User{
					ID:         "user-123",
					MFAEnabled: true,
				}, nil)
			},
			wantErr: ErrMFAAlreadyEnabled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := &Service{
				repo:       mockRepo,
				mfaService: NewMFAService("Test Issuer"),
			}

			result, err := service.SetupMFA(context.Background(), tt.userID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.Secret)
				assert.NotEmpty(t, result.QRCodeURL)
				assert.Len(t, result.RecoveryCodes, 10)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
