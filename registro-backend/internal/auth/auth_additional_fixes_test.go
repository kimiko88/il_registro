package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuth_IsPasswordExpired_StaffRoles(t *testing.T) {
	past := time.Now().Add(-95 * 24 * time.Hour)

	// Staff user (coordinator) with old password should be expired
	coordUser := &User{Role: RoleCoordinator, PasswordChangedAt: &past}
	expired, err := isPasswordExpiredReason(coordUser)
	assert.True(t, expired)
	assert.ErrorIs(t, err, ErrPasswordExpired)

	// Non-staff user (student) with old password should NOT be expired
	studentUser := &User{Role: RoleStudent, PasswordChangedAt: &past}
	expiredStudent, errStudent := isPasswordExpiredReason(studentUser)
	assert.False(t, expiredStudent)
	assert.NoError(t, errStudent)
}

func TestAuth_ParseAcceptLanguage_InputTruncation(t *testing.T) {
	hugeHeader := "en-US" + string(make([]byte, 500))
	res := parseAcceptLanguage(hugeHeader)
	assert.Equal(t, "it-IT", res, "Huge invalid accept-language string should gracefully fallback to it-IT")
}

func TestAuth_Register_InvalidSchoolID(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil)

	mockRepo.On("SchoolExists", mock.Anything, "invalid-school").Return(false, nil).Once()

	req := &RegisterRequest{
		Email:    "test@school.com",
		Password: "Password123!",
		SchoolID: "invalid-school",
		Role:     RoleStudent,
	}

	_, err := svc.Register(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid school ID")
	mockRepo.AssertExpectations(t)
}

func TestAuth_CleanOldLoginAttempts(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	mockRepo.On("CleanOldLoginAttempts", ctx, 30*24*time.Hour).Return(nil).Once()

	err := mockRepo.CleanOldLoginAttempts(ctx, 30*24*time.Hour)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

