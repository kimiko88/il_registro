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
	"registro-backend/pkg/logger"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) SchoolExists(ctx context.Context, schoolID string) (bool, error) {
	args := m.Called(ctx, schoolID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetPasswordHistory(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) AddPasswordHistory(ctx context.Context, userID, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
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

func (m *MockRepository) GetRecentLoginAttemptsByEmail(ctx context.Context, email string, since time.Time) (int, error) {
	args := m.Called(ctx, email, since)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetRecentPasswordResets(ctx context.Context, userID string, since time.Time) (int, error) {
	args := m.Called(ctx, userID, since)
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

func (m *MockRepository) RotateRefreshTokenTx(ctx context.Context, oldID string, newRt *RefreshToken) error {
	args := m.Called(ctx, oldID, newRt)
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

func (m *MockRepository) SaveTempMFASecret(ctx context.Context, userID, secret string) error {
	args := m.Called(ctx, userID, secret)
	return args.Error(0)
}

func (m *MockRepository) ConfirmMFA(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) ConfirmMFAAndSaveRecoveryCodesTx(ctx context.Context, userID string, codes []string) error {
	args := m.Called(ctx, userID, codes)
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

func (m *MockRepository) ResetPasswordTx(ctx context.Context, userID, hash, tokenID string) error {
	args := m.Called(ctx, userID, hash, tokenID)
	return args.Error(0)
}

func (m *MockRepository) ChangePasswordTx(ctx context.Context, userID, hash string) error {
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

func (m *MockRepository) CleanOldLoginAttempts(ctx context.Context, olderThan time.Duration) error {
	args := m.Called(ctx, olderThan)
	return args.Error(0)
}

// Helper to create service with mocks
func setupTest(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	// Initialize logger to avoid nil panics
	logger.Init("info")

	mockRepo := new(MockRepository)
	mockRepo.On("GetRecentLoginAttemptsByEmail", mock.Anything, mock.Anything, mock.Anything).Return(0, nil).Maybe()
	mockRepo.On("GetRecentPasswordResets", mock.Anything, mock.Anything, mock.Anything).Return(0, nil).Maybe()

	// Create a real token manager for testing
	// Generate keys for token manager
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	tokenManager := jwt.NewTokenManager(privateKey, &privateKey.PublicKey)

	// MFA Service also struct dependency
	mfaService := NewMFAService("TestIssuer")

	service := NewService(mockRepo, tokenManager, mfaService, nil)
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

		mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *User) bool {
			return u.Email == req.Email && u.FirstName == req.FirstName
		})).Return(nil).Once()
		mockRepo.On("AddPasswordHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

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

		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(ErrEmailAlreadyExists).Once()

		user, err := s.Register(context.Background(), req)

		assert.ErrorIs(t, err, ErrEmailAlreadyExists)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
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
		mockRepo.On("RecordLoginAttempt", mock.Anything, mock.Anything).Return(nil).Once()

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
			Token:     hashToken(token),
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   false,
			IPAddress: "127.0.0.1",
		}

		user := &User{
			ID:       userID,
			IsActive: true,
		}

		mockRepo.On("GetRefreshToken", mock.Anything, token).Return(rt, nil).Once()
		mockRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil).Once()
		mockRepo.On("RotateRefreshTokenTx", mock.Anything, "rt-1", mock.AnythingOfType("*auth.RefreshToken")).Return(nil).Once()

		pair, err := s.RefreshToken(context.Background(), token, "127.0.0.1", "TestAgent")

		assert.NoError(t, err)
		assert.NotEmpty(t, pair.AccessToken)
		assert.NotEmpty(t, pair.RefreshToken)
		assert.NotEqual(t, token, pair.RefreshToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RevokedToken", func(t *testing.T) {
		rt := &RefreshToken{
			ID:        "rt-1",
			UserID:    userID,
			Token:     hashToken(token),
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   true,
		}

		mockRepo.On("GetRefreshToken", mock.Anything, token).Return(rt, nil).Once()

		_, err := s.RefreshToken(context.Background(), token, "", "")

		assert.ErrorIs(t, err, ErrTokenRevoked)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		rt := &RefreshToken{
			ID:        "rt-1",
			UserID:    userID,
			Token:     hashToken(token),
			ExpiresAt: time.Now().Add(-24 * time.Hour),
			Revoked:   false,
		}

		mockRepo.On("GetRefreshToken", mock.Anything, token).Return(rt, nil).Once()

		_, err := s.RefreshToken(context.Background(), token, "", "")

		assert.ErrorIs(t, err, ErrInvalidToken)
		mockRepo.AssertExpectations(t)
	})
}

func TestPasswordReset(t *testing.T) {
	t.Run("RequestResetSuccess", func(t *testing.T) {
		s, mockRepo := setupTest(t)
		email := "test@example.com"
		user := &User{ID: "user-123", Email: email}

		mockRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil).Once()
		mockRepo.On("CreatePasswordResetToken", mock.Anything, mock.AnythingOfType("*auth.PasswordResetToken")).Return(nil).Once()

		err := s.RequestPasswordReset(context.Background(), email)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RequestResetNormalizedEmail", func(t *testing.T) {
		s, mockRepo := setupTest(t)
		rawEmail := "  Test.User@Example.COM "
		normalizedEmail := "test.user@example.com"
		user := &User{ID: "user-123", Email: normalizedEmail}

		mockRepo.On("GetUserByEmail", mock.Anything, normalizedEmail).Return(user, nil).Once()
		mockRepo.On("CreatePasswordResetToken", mock.Anything, mock.AnythingOfType("*auth.PasswordResetToken")).Return(nil).Once()

		err := s.RequestPasswordReset(context.Background(), rawEmail)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ResetSuccess", func(t *testing.T) {
		s, mockRepo := setupTest(t)
		token := "valid-reset-token"
		newPass := "NewPass123!Safe"

		prt := &PasswordResetToken{
			ID:        "prt-1",
			UserID:    "user-123",
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour),
			Used:      false,
		}

		mockRepo.On("GetPasswordResetToken", mock.Anything, token).Return(prt, nil).Once()
		mockRepo.On("GetUserByID", mock.Anything, "user-123").Return(&User{ID: "user-123", IsActive: true}, nil).Once()
		mockRepo.On("GetPasswordHistory", mock.Anything, "user-123").Return([]string{}, nil).Once()
		mockRepo.On("ResetPasswordTx", mock.Anything, "user-123", mock.Anything, "prt-1").Return(nil).Once()
		mockRepo.On("RevokeAllUserTokens", mock.Anything, "user-123").Return(nil).Once()

		err := s.ResetPassword(context.Background(), token, newPass)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ChangePasswordSuccess", func(t *testing.T) {
		s, mockRepo := setupTest(t)
		userID := "user-change-123"
		currentPass := "OldPass123!Valid"
		newPass := "NewSuperPass123!"

		hashedCurrent, _ := bcrypt.GenerateFromPassword([]byte(currentPass), bcrypt.MinCost)
		user := &User{ID: userID, IsActive: true, PasswordHash: string(hashedCurrent)}

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil).Once()
		mockRepo.On("GetPasswordHistory", mock.Anything, userID).Return([]string{}, nil).Once()
		mockRepo.On("ChangePasswordTx", mock.Anything, userID, mock.Anything).Return(nil).Once()
		mockRepo.On("RevokeAllUserTokens", mock.Anything, userID).Return(nil).Once()

		err := s.ChangePassword(context.Background(), userID, currentPass, newPass)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ChangePasswordInvalidCurrent", func(t *testing.T) {
		s, mockRepo := setupTest(t)
		userID := "user-change-123"
		currentPass := "WrongOldPass123!"
		newPass := "NewSuperPass123!"

		hashedCurrent, _ := bcrypt.GenerateFromPassword([]byte("ActualOldPass123!"), bcrypt.MinCost)
		user := &User{ID: userID, IsActive: true, PasswordHash: string(hashedCurrent)}

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil).Once()

		err := s.ChangePassword(context.Background(), userID, currentPass, newPass)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})
}
