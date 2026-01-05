package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"registro-backend/pkg/jwt"
)

// MockRepository is a mock implementation of the Repository interface
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

func (m *MockRepository) GetRecentLoginAttempts(ctx context.Context, email, ip string, since time.Time) (int, error) {
	args := m.Called(ctx, email, ip, since)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) RecordLoginAttempt(ctx context.Context, attempt *LoginAttempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRepository) GetMFASecret(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) EnableMFA(ctx context.Context, userID, secret string) error {
	args := m.Called(ctx, userID, secret)
	return args.Error(0)
}

func (m *MockRepository) DisableMFA(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
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

func (m *MockRepository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RefreshToken), args.Error(1)
}

func (m *MockRepository) RevokeRefreshToken(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
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

func (m *MockRepository) UpdatePassword(ctx context.Context, userID, hash string) error {
	args := m.Called(ctx, userID, hash)
	return args.Error(0)
}

func (m *MockRepository) UsePasswordResetToken(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Helper to create service with mocks
func setupTest(t *testing.T) (*Service, *MockRepository) {
	mockRepo := new(MockRepository)

	// Create a real token manager for testing
	// Generate keys for token manager
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	tokenManager := jwt.NewTokenManager(privateKey, &privateKey.PublicKey)

	// MFA Service also struct dependency
	mfaService := NewMFAService("TestIssuer")

	service := NewService(mockRepo, tokenManager, mfaService)
	// Lower bcrypt cost for faster tests
	service.bcryptCost = bcrypt.MinCost

	return service, mockRepo
}

func TestRegister(t *testing.T) {
	s, mockRepo := setupTest(t)

	t.Run("Success", func(t *testing.T) {
		req := &RegisterRequest{
			Email:     "test@example.com",
			Password:  "Password123!",
			FirstName: "John",
			LastName:  "Doe",
			Role:      "student",
		}

		mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(nil, errors.New("not found")).Once()
		mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *User) bool {
			return u.Email == req.Email && u.FirstName == req.FirstName
		})).Return(nil).Once()

		user, err := s.Register(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Email, user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmailAlreadyExists", func(t *testing.T) {
		req := &RegisterRequest{
			Email:     "existing@example.com",
			Password:  "Password123!",
			FirstName: "John",
			LastName:  "Doe",
			Role:      "student",
		}

		existingUser := &User{Email: req.Email}
		mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(existingUser, nil).Once()

		user, err := s.Register(context.Background(), req)

		assert.ErrorIs(t, err, ErrEmailAlreadyExists)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("WeakPassword", func(t *testing.T) {
		req := &RegisterRequest{
			Email:    "test@example.com",
			Password: "weak",
		}

		_, err := s.Register(context.Background(), req)
		assert.ErrorIs(t, err, ErrPasswordTooShort)
	})
}

func TestLogin(t *testing.T) {
	s, mockRepo := setupTest(t)

	password := "Password123!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	user := &User{
		ID:            "user-123",
		Email:         "test@example.com",
		PasswordHash:  string(hashedPassword),
		Role:          "student",
		IsActive:      true,
		EmailVerified: true,
	}

	t.Run("Success", func(t *testing.T) {
		req := &LoginRequest{
			Email:    "test@example.com",
			Password: password,
		}

		mockRepo.On("GetRecentLoginAttempts", mock.Anything, req.Email, mock.Anything, mock.Anything).Return(0, nil).Once()
		mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(user, nil).Once()
		mockRepo.On("CreateRefreshToken", mock.Anything, mock.Anything).Return(nil).Once()
		mockRepo.On("UpdateLastLogin", mock.Anything, user.ID).Return(nil).Once()
		mockRepo.On("RecordLoginAttempt", mock.Anything, mock.MatchedBy(func(a *LoginAttempt) bool {
			return a.Success == true
		})).Return(nil).Once()

		resp, err := s.Login(context.Background(), req, "127.0.0.1", "test-agent")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, user.Email, resp.User.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		req := &LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		mockRepo.On("GetRecentLoginAttempts", mock.Anything, req.Email, mock.Anything, mock.Anything).Return(0, nil).Once()
		mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(user, nil).Once()
		mockRepo.On("RecordLoginAttempt", mock.Anything, mock.MatchedBy(func(a *LoginAttempt) bool {
			return a.Success == false
		})).Return(nil).Once()

		resp, err := s.Login(context.Background(), req, "127.0.0.1", "test-agent")

		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserInactive", func(t *testing.T) {
		inactiveUser := &User{
			ID:           "user-inactive",
			Email:        "inactive@example.com",
			PasswordHash: string(hashedPassword),
			IsActive:     false,
		}

		req := &LoginRequest{Email: "inactive@example.com", Password: password}

		mockRepo.On("GetRecentLoginAttempts", mock.Anything, req.Email, mock.Anything, mock.Anything).Return(0, nil).Once()
		mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(inactiveUser, nil).Once()

		resp, err := s.Login(context.Background(), req, "127.0.0.1", "test-agent")

		assert.ErrorIs(t, err, ErrUserInactive)
		assert.Nil(t, resp)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		req := &LoginRequest{
			Email:    "unknown@example.com",
			Password: "password",
		}

		mockRepo.On("GetRecentLoginAttempts", mock.Anything, req.Email, mock.Anything, mock.Anything).Return(0, nil).Once()
		mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(nil, errors.New("not found")).Once()
		mockRepo.On("RecordLoginAttempt", mock.Anything, mock.MatchedBy(func(a *LoginAttempt) bool {
			return a.Success == false
		})).Return(nil).Once()

		resp, err := s.Login(context.Background(), req, "127.0.0.1", "test-agent")

		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})
}

func TestRefreshToken(t *testing.T) {
	s, mockRepo := setupTest(t)

	userID := "user-123"
	token := "refresh-token-123"

	t.Run("Success", func(t *testing.T) {
		rt := &RefreshToken{
			ID:        "rt-1",
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}

		user := &User{
			ID:    userID,
			Email: "test@example.com",
			Role:  "student",
		}

		mockRepo.On("GetRefreshToken", mock.Anything, token).Return(rt, nil).Once()
		mockRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil).Once()

		pair, err := s.RefreshToken(context.Background(), token)

		assert.NoError(t, err)
		assert.NotNil(t, pair)
		assert.Equal(t, token, pair.RefreshToken)
		assert.NotEmpty(t, pair.AccessToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Revoked", func(t *testing.T) {
		rt := &RefreshToken{
			ID:      "rt-1",
			Token:   token,
			Revoked: true,
		}

		mockRepo.On("GetRefreshToken", mock.Anything, token).Return(rt, nil).Once()

		_, err := s.RefreshToken(context.Background(), token)
		assert.ErrorIs(t, err, ErrTokenRevoked)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Expired", func(t *testing.T) {
		rt := &RefreshToken{
			ID:        "rt-1",
			Token:     token,
			ExpiresAt: time.Now().Add(-time.Hour),
		}

		mockRepo.On("GetRefreshToken", mock.Anything, token).Return(rt, nil).Once()

		_, err := s.RefreshToken(context.Background(), token)
		assert.ErrorIs(t, err, ErrInvalidToken)
		mockRepo.AssertExpectations(t)
	})
}

func TestPasswordReset(t *testing.T) {
	s, mockRepo := setupTest(t)

	t.Run("RequestSuccess", func(t *testing.T) {
		email := "test@example.com"
		user := &User{ID: "user-123", Email: email}

		mockRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil).Once()
		mockRepo.On("CreatePasswordResetToken", mock.Anything, mock.Anything).Return(nil).Once()

		err := s.RequestPasswordReset(context.Background(), email)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ResetSuccess", func(t *testing.T) {
		token := "reset-token"
		newPass := "NewPassword123!"

		prt := &PasswordResetToken{
			ID:        "prt-1",
			UserID:    "user-123",
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour),
			Used:      false,
		}

		mockRepo.On("GetPasswordResetToken", mock.Anything, token).Return(prt, nil).Once()
		mockRepo.On("UpdatePassword", mock.Anything, "user-123", mock.Anything).Return(nil).Once()
		mockRepo.On("UsePasswordResetToken", mock.Anything, "prt-1").Return(nil).Once()
		mockRepo.On("RevokeAllUserTokens", mock.Anything, "user-123").Return(nil).Once()

		err := s.ResetPassword(context.Background(), token, newPass)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
