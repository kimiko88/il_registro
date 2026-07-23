package jwt

import (
	"testing"
	"time"
)

func TestTokenManager_GenerateAndValidateAccessToken(t *testing.T) {
	privKey, pubKey, err := GenerateRSAKeys()
	if err != nil {
		t.Fatalf("failed to generate RSA keys: %v", err)
	}

	tm := NewTokenManager(privKey, pubKey)

	tokenStr, err := tm.GenerateAccessToken("user-123", "test@school.it", "teacher", "school-abc")
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	claims, err := tm.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("expected UserID 'user-123', got '%s'", claims.UserID)
	}
	if claims.Email != "test@school.it" {
		t.Errorf("expected Email 'test@school.it', got '%s'", claims.Email)
	}
	if claims.Role != "teacher" {
		t.Errorf("expected Role 'teacher', got '%s'", claims.Role)
	}
	if claims.SchoolID != "school-abc" {
		t.Errorf("expected SchoolID 'school-abc', got '%s'", claims.SchoolID)
	}
}

func TestTokenManager_ValidateInvalidToken(t *testing.T) {
	privKey, pubKey, err := GenerateRSAKeys()
	if err != nil {
		t.Fatalf("failed to generate RSA keys: %v", err)
	}

	tm := NewTokenManager(privKey, pubKey)

	_, err = tm.ValidateToken("invalid.jwt.string")
	if err == nil {
		t.Error("expected error for invalid token string, got nil")
	}
}

func TestTokenManager_GenerateRefreshToken(t *testing.T) {
	privKey, pubKey, err := GenerateRSAKeys()
	if err != nil {
		t.Fatalf("failed to generate RSA keys: %v", err)
	}

	tm := NewTokenManager(privKey, pubKey)

	tokenStr, err := tm.GenerateRefreshToken("user-456")
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	if tokenStr == "" {
		t.Error("expected non-empty refresh token string")
	}
}

func TestTokenManager_TokenExpiration(t *testing.T) {
	privKey, pubKey, err := GenerateRSAKeys()
	if err != nil {
		t.Fatalf("failed to generate RSA keys: %v", err)
	}

	tm := NewTokenManager(privKey, pubKey)
	tm.accessTokenTTL = -1 * time.Second // Expired immediately

	tokenStr, err := tm.GenerateAccessToken("user-789", "exp@school.it", "student", "")
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	_, err = tm.ValidateToken(tokenStr)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestTokenManager_GetAccessTokenTTL(t *testing.T) {
	privKey, pubKey, err := GenerateRSAKeys()
	if err != nil {
		t.Fatalf("failed to generate RSA keys: %v", err)
	}

	tm := NewTokenManager(privKey, pubKey)
	ttl := tm.GetAccessTokenTTL()
	if ttl != 900 {
		t.Errorf("expected 900 seconds TTL, got %d", ttl)
	}
}
