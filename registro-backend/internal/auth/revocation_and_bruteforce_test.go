package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/users"
	"registro-backend/pkg/crypto"
	"registro-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepoForRevocation struct {
	users.Repository
}

func (m *mockUserRepoForRevocation) GetByID(ctx context.Context, id string) (*users.User, error) {
	return &users.User{
		ID:       id,
		Email:    "test@school.it",
		Role:     "teacher",
		IsActive: true,
	}, nil
}

func (m *mockUserRepoForRevocation) IsActive(ctx context.Context, userID string) (bool, error) {
	return true, nil
}

func TestAuthMiddleware_RevocationCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	priv, pub, err := jwt.GetOrGenerateKeys("", "")
	require.NoError(t, err)

	tm := jwt.NewTokenManager(priv, pub)
	token, err := tm.GenerateAccessToken("user-123", "teacher@school.it", "teacher", "school-1")
	require.NoError(t, err)

	jti, remaining, err := tm.ExtractJTI(token)
	require.NoError(t, err)
	require.NotEmpty(t, jti)

	revStore := jwt.NewMemoryRevocationStore()
	defer func() { _ = revStore.Close() }()

	userRepo := &mockUserRepoForRevocation{}
	middleware := NewMiddleware(tm, userRepo, revStore)

	r := gin.New()
	r.GET("/protected", middleware.Authenticate(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 1. Initial request with valid token -> 200 OK
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Revoke the token
	err = revStore.Revoke(context.Background(), jti, remaining)
	require.NoError(t, err)

	// 3. Subsequent request with revoked token -> 401 Unauthorized
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
	assert.Contains(t, w2.Body.String(), "token has been revoked")
}

func TestBcryptPrehash_LongPasswordAndMigration(t *testing.T) {
	// Password > 72 characters (testing beyond bcrypt's truncation limit)
	longPassword := "A_very_long_secure_passphrase_exceeding_seventy_two_characters_completely_safe_for_production_use_2026!"
	assert.True(t, len(longPassword) > 72)

	prehashed := crypto.PrehashPassword(longPassword)
	assert.Equal(t, 32, len(prehashed)) // HMAC-SHA256 is exactly 32 bytes

	hash, err := bcrypt.GenerateFromPassword(prehashed, bcrypt.MinCost)
	require.NoError(t, err)

	// Compare with identical password -> success
	err = bcrypt.CompareHashAndPassword(hash, crypto.PrehashPassword(longPassword))
	assert.NoError(t, err)

	// Compare with password differing only on the 73+ character -> fails (proving no 72-byte truncation)
	slightlyDifferent := "A_very_long_secure_passphrase_exceeding_seventy_two_characters_completely_safe_for_production_use_2027?"
	err = bcrypt.CompareHashAndPassword(hash, crypto.PrehashPassword(slightlyDifferent))
	assert.Error(t, err)
}

func TestLoginRateLimiter_FallbackOnNil(t *testing.T) {
	var limiter *RedisLoginRateLimiter = nil
	err := limiter.Check(context.Background(), "user@school.it", "127.0.0.1")
	assert.NoError(t, err)

	limiter = NewRedisLoginRateLimiter(nil)
	err = limiter.Check(context.Background(), "user@school.it", "127.0.0.1")
	assert.NoError(t, err)
	limiter.RecordFailure(context.Background(), "user@school.it", "127.0.0.1")
	limiter.RecordSuccess(context.Background(), "user@school.it", "127.0.0.1")
}

func TestNewLoginRateLimiter_InvalidURL(t *testing.T) {
	limiter := NewLoginRateLimiter("redis://invalid-host:6379")
	assert.Nil(t, limiter)
}
