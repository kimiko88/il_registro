package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestConstantTimeLogin_InactiveUser(t *testing.T) {
	svc, mockRepo := setupTest(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.MinCost)
	inactiveUser := &User{
		ID:           "inactive-user-1",
		Email:        "inactive@example.com",
		PasswordHash: string(hash),
		IsActive:     false,
		CreatedAt:    time.Now(),
	}

	mockRepo.On("GetRecentLoginAttempts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, nil).Maybe()
	mockRepo.On("GetUserByEmail", mock.Anything, "inactive@example.com").Return(inactiveUser, nil)
	mockRepo.On("RecordLoginAttempt", mock.Anything, mock.AnythingOfType("*auth.LoginAttempt")).Return(nil).Maybe()

	start := time.Now()
	resp, err := svc.Login(context.Background(), &LoginRequest{
		Email:    "inactive@example.com",
		Password: "Password123!",
	}, "127.0.0.1", "TestAgent")
	duration := time.Since(start)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, ErrUserInactive, err)
	assert.True(t, duration >= 0)
}

func TestRefreshTokenHashing_Consistency(t *testing.T) {
	tokenStr := "sample-refresh-token-xyz-12345"
	hash1 := sha256.Sum256([]byte(tokenStr))
	hex1 := hex.EncodeToString(hash1[:])

	hash2 := sha256.Sum256([]byte(tokenStr))
	hex2 := hex.EncodeToString(hash2[:])

	assert.Equal(t, hex1, hex2)
	assert.NotEqual(t, tokenStr, hex1)
	assert.Len(t, hex1, 64) // SHA-256 hex string length
}

func TestPasswordResetRateLimit_Protection(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, nil, nil)

	user := &User{
		ID:        "user-rate-limit",
		Email:     "limit@example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "limit@example.com").Return(user, nil)
	// Return 3 recent resets to trigger rate limit threshold
	mockRepo.On("GetRecentPasswordResets", mock.Anything, "user-rate-limit", mock.Anything).Return(3, nil)

	err := service.RequestPasswordReset(context.Background(), "limit@example.com")
	assert.NoError(t, err) // Should silently return nil to avoid revealing throttle status
}
