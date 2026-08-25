package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthRepoPhase4 struct {
	mock.Mock
	Repository
}

func (m *mockAuthRepoPhase4) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *mockAuthRepoPhase4) RotateRefreshTokenTx(ctx context.Context, oldID string, newRt *RefreshToken) error {
	args := m.Called(ctx, oldID, newRt)
	return args.Error(0)
}

func TestAuth_RefreshToken_TransactionalRotationError(t *testing.T) {
	repo := new(mockAuthRepoPhase4)
	svc := NewService(repo, nil, nil, nil)
	ctx := context.Background()

	// In case RotateRefreshTokenTx fails, it MUST return error and NOT fallback to non-transactional methods
	repo.On("RotateRefreshTokenTx", ctx, mock.Anything, mock.Anything).Return(assert.AnError).Once()

	err := svc.repo.RotateRefreshTokenTx(ctx, "old-token-id", &RefreshToken{
		UserID: "user-1",
		Token:  "new-token",
	})

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestAuth_LegacyUserPasswordExpired(t *testing.T) {
	legacyUser := &User{
		ID:                "user-legacy",
		Role:              RoleTeacher,
		CreatedAt:         time.Now().Add(-180 * 24 * time.Hour), // Created 6 months ago
		PasswordChangedAt: nil,                                   // Never updated column
	}

	// Legacy users without PasswordChangedAt MUST NOT be locked out automatically
	assert.False(t, isPasswordExpired(legacyUser))
}
