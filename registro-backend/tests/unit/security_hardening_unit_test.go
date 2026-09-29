package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"registro-backend/pkg/crypto"
	"registro-backend/pkg/jwt"
)

func TestSecurity_JTIRevocationStoreLifecycle(t *testing.T) {
	store := jwt.NewMemoryRevocationStore()
	defer func() { _ = store.Close() }()

	ctx := context.Background()
	jti := "test-jti-uuid-123456"

	// 1. Initially not revoked
	revoked, err := store.IsRevoked(ctx, jti)
	assert.NoError(t, err)
	assert.False(t, revoked)

	// 2. Revoke with positive TTL
	err = store.Revoke(ctx, jti, 500*time.Millisecond)
	assert.NoError(t, err)

	revoked, err = store.IsRevoked(ctx, jti)
	assert.NoError(t, err)
	assert.True(t, revoked, "Token should be marked as revoked")

	// 3. Querying empty JTI returns false
	revoked, err = store.IsRevoked(ctx, "")
	assert.NoError(t, err)
	assert.False(t, revoked)

	// 4. Revoking empty JTI or non-positive TTL is a no-op
	assert.NoError(t, store.Revoke(ctx, "", 1*time.Minute))
	assert.NoError(t, store.Revoke(ctx, "valid-jti", 0))
	assert.NoError(t, store.Revoke(ctx, "valid-jti", -1*time.Second))

	// 5. Wait for TTL to expire
	time.Sleep(550 * time.Millisecond)
	revoked, err = store.IsRevoked(ctx, jti)
	assert.NoError(t, err)
	assert.False(t, revoked, "Token revocation should naturally expire after TTL")
}

func TestSecurity_BcryptPrehashing_NoTruncation(t *testing.T) {
	// Bcrypt has a hard limit of 72 bytes. With PrehashPassword,
	// arbitrary length passphrases are hashed via HMAC-SHA256 (32 bytes).
	basePass := "QuestoTestDimostraCheIlPrehashingDellePasswordSuperaIlLimiteDei72ByteDiBcryptInModoTrasparenteEdEstremamenteSicuro!"
	require.True(t, len(basePass) > 72)

	h1 := crypto.PrehashPassword(basePass)
	assert.Equal(t, 32, len(h1))

	// Generate bcrypt hash from the 32-byte digest
	bcryptHash, err := bcrypt.GenerateFromPassword(h1, bcrypt.MinCost)
	require.NoError(t, err)

	// Verify exact match
	err = bcrypt.CompareHashAndPassword(bcryptHash, crypto.PrehashPassword(basePass))
	assert.NoError(t, err)

	// Single character modification at index 100
	modifiedPass := basePass[:100] + "?" + basePass[101:]
	err = bcrypt.CompareHashAndPassword(bcryptHash, crypto.PrehashPassword(modifiedPass))
	assert.Error(t, err, "Modification beyond 72 bytes must be detected by HMAC digest comparison")
}

func TestSecurity_DiagnosisFieldEncryption_RoundTrip(t *testing.T) {
	diagnosisPlain := "Diagnosi Funzionale: DSA con discalculia e disgrafia (Codice ICD-10 F81.2)"

	encrypted, err := crypto.EncryptString(diagnosisPlain)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, diagnosisPlain, encrypted)

	decrypted, err := crypto.DecryptString(encrypted)
	require.NoError(t, err)
	assert.Equal(t, diagnosisPlain, decrypted)

	// Legacy unencrypted text fallback test
	legacyPlain := "Diagnosi ante-crittografia non cifrata"
	failedDec, err := crypto.DecryptString(legacyPlain)
	assert.Error(t, err)
	assert.Empty(t, failedDec)
}
