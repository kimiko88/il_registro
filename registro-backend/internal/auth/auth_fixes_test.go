package auth

import (
	"context"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/pkg/crypto"
)

func TestAuthFixes_LoginCaseInsensitiveRateLimit(t *testing.T) {
	s, mockRepo := setupTest(t)

	mockRepo.On("GetRecentLoginAttempts", mock.Anything, "mario@school.it", "127.0.0.1", mock.Anything).Return(5, nil).Once()

	req := &LoginRequest{
		Email:    "Mario@School.it",
		Password: "Password123!",
	}

	res, err := s.Login(context.Background(), req, "127.0.0.1", "test-agent")
	assert.Nil(t, res)
	assert.ErrorIs(t, err, ErrTooManyAttempts)
	mockRepo.AssertExpectations(t)
}

func TestAuthFixes_VerifyMFA_AtomicTx(t *testing.T) {
	s, mockRepo := setupTest(t)

	rawSecret := "JBSWY3DPEHPK3PXP"
	encryptedSecret, _ := crypto.EncryptString(rawSecret)

	mockRepo.On("GetMFASecret", mock.Anything, "u123").Return(encryptedSecret, nil).Once()

	totpToken, err := totp.GenerateCode(rawSecret, time.Now())
	assert.NoError(t, err)

	mockRepo.On("ConfirmMFAAndSaveRecoveryCodesTx", mock.Anything, "u123", mock.MatchedBy(func(codes []string) bool {
		return len(codes) == 10
	})).Return(nil).Once()

	codes, err := s.VerifyMFA(context.Background(), "u123", totpToken)
	assert.NoError(t, err)
	assert.Len(t, codes, 10)
	mockRepo.AssertExpectations(t)
}

func TestAuthFixes_CoordinatorPasswordExpired(t *testing.T) {
	past := time.Now().Add(-91 * 24 * time.Hour)
	u := &User{
		Role:              RoleCoordinator,
		PasswordChangedAt: &past,
	}

	expired, err := isPasswordExpiredReason(u)
	assert.True(t, expired)
	assert.ErrorIs(t, err, ErrPasswordExpired)
}

func TestAuthFixes_ResetPassword_SingleAtomicQuery(t *testing.T) {
	s, mockRepo := setupTest(t)

	mockRepo.On("GetPasswordResetToken", mock.Anything, "invalid-or-expired-token").Return(nil, ErrInvalidToken).Once()

	err := s.ResetPassword(context.Background(), "invalid-or-expired-token", "NewPassword123!")
	assert.ErrorIs(t, err, ErrInvalidToken)
	mockRepo.AssertExpectations(t)
}

func TestAuthFixes_Register_TOCTOU(t *testing.T) {
	s, mockRepo := setupTest(t)

	req := &RegisterRequest{
		Email:    "newuser@school.it",
		Password: "Password123!",
		Role:     "teacher",
	}

	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *User) bool {
		return u.Email == "newuser@school.it"
	})).Return(nil).Once()
	mockRepo.On("AddPasswordHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	user, err := s.Register(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser@school.it", user.Email)
	mockRepo.AssertExpectations(t)
}
