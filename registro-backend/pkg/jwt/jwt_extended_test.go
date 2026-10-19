package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Helpers ──────────────────────────────────────────────────────────────────

func marshalPrivKeyPEM(key *rsa.PrivateKey) string {
	privBytes := x509.MarshalPKCS1PrivateKey(key)
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}))
}

func marshalPubKeyPEM(key *rsa.PublicKey) string {
	pubBytes := x509.MarshalPKCS1PublicKey(key)
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes}))
}

func bgCtx() context.Context { return context.Background() }

// ─── TokenManager ─────────────────────────────────────────────────────────────

// TestTokenManager_ExtractJTI_ExpiredToken verifies that ExtractJTI fails on an expired token.
func TestTokenManager_ExtractJTI_ExpiredToken(t *testing.T) {
	priv, pub, err := GenerateRSAKeys()
	require.NoError(t, err)

	tm := NewTokenManager(priv, pub)
	tm.accessTokenTTL = -1 * time.Second // immediately expired

	token, err := tm.GenerateAccessToken("user-exp", "e@t.it", "teacher", "s1")
	require.NoError(t, err)

	_, _, err = tm.ExtractJTI(token)
	require.Error(t, err, "ExtractJTI must fail on an expired token")
}

// TestTokenManager_ExtractJTI_InvalidToken verifies ExtractJTI returns error for garbage input.
func TestTokenManager_ExtractJTI_InvalidToken(t *testing.T) {
	priv, pub, err := GenerateRSAKeys()
	require.NoError(t, err)
	tm := NewTokenManager(priv, pub)

	_, _, err = tm.ExtractJTI("not-a-jwt")
	require.Error(t, err)
}

// TestTokenManager_RefreshTokenCannotValidateAsAccess verifies that a refresh token
// (no audience) is rejected by ValidateToken which requires the registered audience.
func TestTokenManager_RefreshTokenCannotValidateAsAccess(t *testing.T) {
	priv, pub, err := GenerateRSAKeys()
	require.NoError(t, err)

	tm := NewTokenManager(priv, pub)
	refreshToken, err := tm.GenerateRefreshToken("user-refresh")
	require.NoError(t, err)

	_, err = tm.ValidateToken(refreshToken)
	require.Error(t, err, "refresh token must be rejected by ValidateToken (no audience)")
}

// TestTokenManager_AccessTokenHasJTI verifies that GenerateAccessToken sets a non-empty JTI.
func TestTokenManager_AccessTokenHasJTI(t *testing.T) {
	priv, pub, err := GenerateRSAKeys()
	require.NoError(t, err)
	tm := NewTokenManager(priv, pub)

	token, err := tm.GenerateAccessToken("user-1", "u@t.it", "teacher", "s1")
	require.NoError(t, err)

	jti, ttl, err := tm.ExtractJTI(token)
	require.NoError(t, err)
	assert.NotEmpty(t, jti)
	assert.True(t, ttl > 0)
	assert.True(t, ttl <= 15*time.Minute)
}

// TestTokenManager_AccessTokenEmptySchoolID verifies tokens without a schoolID are valid.
func TestTokenManager_AccessTokenEmptySchoolID(t *testing.T) {
	priv, pub, err := GenerateRSAKeys()
	require.NoError(t, err)
	tm := NewTokenManager(priv, pub)

	token, err := tm.GenerateAccessToken("user-2", "u@t.it", "parent", "")
	require.NoError(t, err)

	claims, err := tm.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "", claims.SchoolID)
	assert.Equal(t, "parent", claims.Role)
}

// TestTokenManager_JTIsUnique verifies that two access tokens have different JTIs.
func TestTokenManager_JTIsUnique(t *testing.T) {
	priv, pub, err := GenerateRSAKeys()
	require.NoError(t, err)
	tm := NewTokenManager(priv, pub)

	tok1, _ := tm.GenerateAccessToken("u", "u@t.it", "teacher", "s")
	tok2, _ := tm.GenerateAccessToken("u", "u@t.it", "teacher", "s")

	jti1, _, _ := tm.ExtractJTI(tok1)
	jti2, _, _ := tm.ExtractJTI(tok2)
	assert.NotEqual(t, jti1, jti2, "every access token must have a unique JTI")
}

// TestTokenManager_ClaimsRoles verifies all standard roles can be embedded and read back.
func TestTokenManager_ClaimsRoles(t *testing.T) {
	priv, pub, err := GenerateRSAKeys()
	require.NoError(t, err)
	tm := NewTokenManager(priv, pub)

	roles := []string{"admin", "teacher", "student", "parent", "ata"}
	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			tok, err := tm.GenerateAccessToken("uid", "u@t.it", role, "s1")
			require.NoError(t, err)
			claims, err := tm.ValidateToken(tok)
			require.NoError(t, err)
			assert.Equal(t, role, claims.Role)
		})
	}
}

// ─── MemoryRevocationStore ────────────────────────────────────────────────────

// TestMemoryRevocationStore_Revoke_EmptyJTI verifies that revoking an empty JTI is a no-op.
func TestMemoryRevocationStore_Revoke_EmptyJTI(t *testing.T) {
	store := NewMemoryRevocationStore()
	ctx := bgCtx()

	err := store.Revoke(ctx, "", 5*time.Minute)
	assert.NoError(t, err)

	revoked, err := store.IsRevoked(ctx, "")
	assert.NoError(t, err)
	assert.False(t, revoked)
}

// TestMemoryRevocationStore_Revoke_NegativeTTL verifies that a non-positive TTL is a no-op.
func TestMemoryRevocationStore_Revoke_NegativeTTL(t *testing.T) {
	store := NewMemoryRevocationStore()
	ctx := bgCtx()

	err := store.Revoke(ctx, "jti-123", -1*time.Second)
	assert.NoError(t, err)

	revoked, err := store.IsRevoked(ctx, "jti-123")
	assert.NoError(t, err)
	assert.False(t, revoked, "negative TTL should not revoke the JTI")
}

// TestMemoryRevocationStore_Close verifies Close is a no-op and returns nil.
func TestMemoryRevocationStore_Close(t *testing.T) {
	store := NewMemoryRevocationStore()
	assert.NoError(t, store.Close())
}

// ─── PEM Parsing ──────────────────────────────────────────────────────────────

// TestParsePrivateKeyPEM_InvalidPEM verifies error when given invalid PEM bytes.
func TestParsePrivateKeyPEM_InvalidPEM(t *testing.T) {
	_, err := ParsePrivateKeyPEM([]byte("this is not PEM"))
	require.Error(t, err)
}

// TestParsePublicKeyPEM_InvalidPEM verifies error when given invalid PEM bytes.
func TestParsePublicKeyPEM_InvalidPEM(t *testing.T) {
	_, err := ParsePublicKeyPEM([]byte("this is not PEM"))
	require.Error(t, err)
}

// TestSaveAndLoadRSAKeys_RoundTrip verifies that keys saved to disk can be loaded correctly.
func TestSaveAndLoadRSAKeys_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	privPath := filepath.Join(tmpDir, "priv.pem")
	pubPath := filepath.Join(tmpDir, "pub.pem")

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	err = SaveRSAKeys(privKey, privPath, pubPath)
	require.NoError(t, err)

	loadedPriv, loadedPub, err := LoadRSAKeys(privPath, pubPath)
	require.NoError(t, err)

	assert.Equal(t, privKey.D, loadedPriv.D, "private key must match after save/load")
	assert.Equal(t, privKey.N, loadedPub.N, "public key must match after save/load")
}

// TestLoadRSAKeys_MissingFile verifies error when private key file is missing.
func TestLoadRSAKeys_MissingFile(t *testing.T) {
	_, _, err := LoadRSAKeys("/nonexistent/priv.pem", "/nonexistent/pub.pem")
	require.Error(t, err)
}

// ─── GetOrGenerateKeys ────────────────────────────────────────────────────────

// TestGetOrGenerateKeys_JWT_PrivKeyEnv verifies that JWT_PRIVATE_KEY env var is accepted.
func TestGetOrGenerateKeys_JWT_PrivKeyEnv(t *testing.T) {
	priv, _, err := GenerateRSAKeys()
	require.NoError(t, err)

	privPEM := marshalPrivKeyPEM(priv)

	t.Setenv("RSA_PRIVATE_KEY", "")
	t.Setenv("JWT_PRIVATE_KEY", privPEM)

	loaded, _, err := GetOrGenerateKeys("", "")
	require.NoError(t, err)
	assert.Equal(t, priv.D, loaded.D)
}

// TestGetOrGenerateKeys_WithPublicKeyEnv verifies that RSA_PUBLIC_KEY env var is accepted.
func TestGetOrGenerateKeys_WithPublicKeyEnv(t *testing.T) {
	priv, pub, err := GenerateRSAKeys()
	require.NoError(t, err)

	privPEM := marshalPrivKeyPEM(priv)
	pubPEM := marshalPubKeyPEM(pub)

	t.Setenv("RSA_PRIVATE_KEY", privPEM)
	t.Setenv("RSA_PUBLIC_KEY", pubPEM)

	loadedPriv, loadedPub, err := GetOrGenerateKeys("", "")
	require.NoError(t, err)
	assert.Equal(t, priv.D, loadedPriv.D)
	assert.Equal(t, pub.N, loadedPub.N)
}

// TestGetOrGenerateKeys_InvalidPrivKeyEnv verifies that a malformed private key env causes error.
func TestGetOrGenerateKeys_InvalidPrivKeyEnv(t *testing.T) {
	t.Setenv("RSA_PRIVATE_KEY", "not-valid-pem-data")
	t.Setenv("JWT_PRIVATE_KEY", "")

	_, _, err := GetOrGenerateKeys("", "")
	require.Error(t, err, "invalid RSA_PRIVATE_KEY must cause GetOrGenerateKeys to error")
}

// ─── NewRevocationStore fallback ──────────────────────────────────────────────

// TestNewRevocationStore_EmptyURL falls back to memory.
func TestNewRevocationStore_EmptyURL(t *testing.T) {
	store := NewRevocationStore("")
	require.NotNil(t, store)
	_, ok := store.(*MemoryRevocationStore)
	assert.True(t, ok)
	_ = store.Close()
}

// TestNewRevocationStore_UnreachableRedis falls back to memory.
func TestNewRevocationStore_UnreachableRedis(t *testing.T) {
	store := NewRevocationStore("redis://127.0.0.1:1/0")
	require.NotNil(t, store)
	_, ok := store.(*MemoryRevocationStore)
	assert.True(t, ok)
	_ = store.Close()
}

// TestNewRevocationStore_InvalidRedisURL falls back to memory.
func TestNewRevocationStore_InvalidRedisURL(t *testing.T) {
	store := NewRevocationStore("not-a-redis-url")
	require.NotNil(t, store)
	_, ok := store.(*MemoryRevocationStore)
	assert.True(t, ok)
	_ = store.Close()
}
