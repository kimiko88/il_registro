package integration

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/auth"
	"registro-backend/internal/middleware"
	"registro-backend/internal/users"
	"registro-backend/pkg/crypto"
	"registro-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepoForSecurityIntegration struct {
	users.Repository
}

func (m *mockUserRepoForSecurityIntegration) GetByID(ctx context.Context, id string) (*users.User, error) {
	return &users.User{
		ID:       id,
		Email:    "security-test@school.it",
		Role:     "teacher",
		IsActive: true,
	}, nil
}

func (m *mockUserRepoForSecurityIntegration) IsActive(ctx context.Context, userID string) (bool, error) {
	return true, nil
}

func TestIntegration_AuthSecurityHardening_RevocationAndProtection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	priv, pub, err := jwt.GetOrGenerateKeys("", "")
	require.NoError(t, err)

	tm := jwt.NewTokenManager(priv, pub)
	revStore := jwt.NewMemoryRevocationStore()
	defer func() { _ = revStore.Close() }()

	userRepo := &mockUserRepoForSecurityIntegration{}
	authMid := auth.NewMiddleware(tm, userRepo, revStore)

	r := gin.New()
	r.Use(middleware.RequestBodyLimitMiddleware(1024)) // 1KB limit for testing

	// Protected endpoint
	r.GET("/api/v1/secure/resource", authMid.Authenticate(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "authorized", "user": c.GetString("user_id")})
	})

	// Logout endpoint that revokes token JTI
	r.POST("/api/v1/auth/logout", authMid.Authenticate(), func(c *gin.Context) {
		jti := c.GetString("jti")
		if jti != "" {
			_ = revStore.Revoke(c.Request.Context(), jti, 1*time.Hour)
		}
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	})

	// Test endpoint with payload limit
	r.POST("/api/v1/secure/data", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "payload too large"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"read": len(body)})
	})

	// 1. Issue a valid token
	token, err := tm.GenerateAccessToken("user-sec-1", "security-test@school.it", "teacher", "school-1")
	require.NoError(t, err)

	// 2. Access protected resource with valid token -> 200 OK
	req1, _ := http.NewRequest("GET", "/api/v1/secure/resource", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 3. User logs out -> token JTI is revoked
	reqLogout, _ := http.NewRequest("POST", "/api/v1/auth/logout", nil)
	reqLogout.Header.Set("Authorization", "Bearer "+token)
	wLogout := httptest.NewRecorder()
	r.ServeHTTP(wLogout, reqLogout)
	assert.Equal(t, http.StatusOK, wLogout.Code)

	// 4. Subsequent request with the same token is rejected with 401 Unauthorized
	req2, _ := http.NewRequest("GET", "/api/v1/secure/resource", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
	assert.Contains(t, w2.Body.String(), "token has been revoked")

	// 5. Oversized payload (>1KB) rejected by RequestBodyLimitMiddleware
	oversized := bytes.Repeat([]byte("X"), 2048)
	reqOver, _ := http.NewRequest("POST", "/api/v1/secure/data", bytes.NewReader(oversized))
	wOver := httptest.NewRecorder()
	r.ServeHTTP(wOver, reqOver)
	assert.Equal(t, http.StatusRequestEntityTooLarge, wOver.Code)

	// 6. Multipart request with oversized payload passes through (exempted from limit)
	reqMultipart, _ := http.NewRequest("POST", "/api/v1/secure/data", bytes.NewReader(oversized))
	reqMultipart.Header.Set("Content-Type", "multipart/form-data; boundary=my-boundary")
	wMultipart := httptest.NewRecorder()
	r.ServeHTTP(wMultipart, reqMultipart)
	assert.Equal(t, http.StatusOK, wMultipart.Code)
}

func TestIntegration_AuthSecurityHardening_PasswordPrehashingIntegration(t *testing.T) {
	// Verify that password prehashing prevents the 72-byte truncation attack
	passwordBase := "SecureMasterKey2026!WithLengthExceedingTheSeventyTwoByteBcryptTruncationLimitForEnterpriseSchoolRegistroV2"
	require.True(t, len(passwordBase) > 72)

	h1 := crypto.PrehashPassword(passwordBase)
	hash, err := bcrypt.GenerateFromPassword(h1, bcrypt.MinCost)
	require.NoError(t, err)

	// Correct password authentication
	err = bcrypt.CompareHashAndPassword(hash, crypto.PrehashPassword(passwordBase))
	assert.NoError(t, err)

	// Password changed at index 80 (beyond standard bcrypt 72 byte limit)
	tamperedPassword := passwordBase[:80] + "Z" + passwordBase[81:]
	err = bcrypt.CompareHashAndPassword(hash, crypto.PrehashPassword(tamperedPassword))
	assert.Error(t, err, "Bcrypt with prehashing must detect modifications beyond 72 bytes")
}
