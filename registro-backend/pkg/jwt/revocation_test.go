package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRevocationStore(t *testing.T) {
	store := NewMemoryRevocationStore()
	defer func() { _ = store.Close() }()

	ctx := context.Background()
	jti := "test-jti-uuid-1234"

	// Initially not revoked
	revoked, err := store.IsRevoked(ctx, jti)
	require.NoError(t, err)
	assert.False(t, revoked)

	// Revoke with short TTL
	err = store.Revoke(ctx, jti, 500*time.Millisecond)
	require.NoError(t, err)

	// Now should be revoked
	revoked, err = store.IsRevoked(ctx, jti)
	require.NoError(t, err)
	assert.True(t, revoked)

	// Another JTI is still not revoked
	revokedOther, err := store.IsRevoked(ctx, "other-jti")
	require.NoError(t, err)
	assert.False(t, revokedOther)

	// After expiry
	time.Sleep(600 * time.Millisecond)
	revokedAfter, err := store.IsRevoked(ctx, jti)
	require.NoError(t, err)
	assert.False(t, revokedAfter)
}

func TestTokenManager_ExtractJTI(t *testing.T) {
	priv, pub, err := GetOrGenerateKeys("", "")
	require.NoError(t, err)

	tm := NewTokenManager(priv, pub)
	token, err := tm.GenerateAccessToken("user-1", "user@test.it", "teacher", "school-1")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	jti, remaining, err := tm.ExtractJTI(token)
	require.NoError(t, err)
	assert.NotEmpty(t, jti)
	assert.True(t, remaining > 0)
	assert.True(t, remaining <= 15*time.Minute)
}
