package auth

import (
	"context"
	"strings"
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

func TestAuthFixes_RefreshToken_UserAgentSanitization(t *testing.T) {
	s, mockRepo := setupTest(t)

	rawToken := "valid-refresh-token"
	hashedToken := hashToken(rawToken)

	rt := &RefreshToken{
		ID:        "rt-1",
		UserID:    "u-1",
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IPAddress: "127.0.0.1",
		UserAgent: "original-agent",
	}

	user := &User{
		ID:       "u-1",
		IsActive: true,
		Role:     RoleTeacher,
	}

	mockRepo.On("GetRefreshToken", mock.Anything, rawToken).Return(rt, nil).Once()
	mockRepo.On("GetUserByID", mock.Anything, "u-1").Return(user, nil).Once()
	mockRepo.On("RotateRefreshTokenTx", mock.Anything, "rt-1", mock.MatchedBy(func(newRt *RefreshToken) bool {
		// Verify control characters (such as \x1b or \x00) are stripped
		return !strings.ContainsAny(newRt.UserAgent, "\x1b\x00\r\n") && strings.Contains(newRt.UserAgent, "malicious-agent")
	})).Return(nil).Once()

	dirtyUserAgent := "malicious-agent\x1b[31m\x00\r\ntest"
	tokens, err := s.RefreshToken(context.Background(), rawToken, "127.0.0.1", dirtyUserAgent)
	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	mockRepo.AssertExpectations(t)
}

func TestAuthFixes_PasswordExpired_StudentParentExempt(t *testing.T) {
	past := time.Now().Add(-180 * 24 * time.Hour) // 6 months ago

	student := &User{
		Role:              RoleStudent,
		PasswordChangedAt: &past,
	}
	expiredStudent, errStudent := isPasswordExpiredReason(student)
	assert.False(t, expiredStudent)
	assert.NoError(t, errStudent)

	parent := &User{
		Role:              RoleParent,
		PasswordChangedAt: &past,
	}
	expiredParent, errParent := isPasswordExpiredReason(parent)
	assert.False(t, expiredParent)
	assert.NoError(t, errParent)
}
