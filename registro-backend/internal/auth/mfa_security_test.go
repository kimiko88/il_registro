package auth

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMFAService_RejectsTOTPReplayInMemoryMode(t *testing.T) {
	service := NewMFAService("test")
	secret, err := service.GenerateSecret()
	require.NoError(t, err)
	code, err := totp.GenerateCode(secret, time.Now())
	require.NoError(t, err)

	assert.True(t, service.VerifyTOTPUser("user-1", secret, code))
	assert.False(t, service.VerifyTOTPUser("user-1", secret, code))
}

func TestMFAService_FailsClosedWhenConfiguredRedisIsInvalid(t *testing.T) {
	service := NewMFAService("test", "not-a-redis-url")
	defer func() { _ = service.Close() }()
	secret, err := service.GenerateSecret()
	require.NoError(t, err)
	code, err := totp.GenerateCode(secret, time.Now())
	require.NoError(t, err)

	assert.False(t, service.VerifyTOTPUser("user-1", secret, code))
}
