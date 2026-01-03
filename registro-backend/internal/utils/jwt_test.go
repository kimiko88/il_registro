package utils

import (
	"registro-backend/internal/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	cfg := config.JWTConfig{
		Secret: "test-secret-key-for-testing",
	}

	tests := []struct {
		name    string
		userID  uint
		role    string
		wantErr bool
	}{
		{
			name:    "successful token generation for admin",
			userID:  1,
			role:    "admin",
			wantErr: false,
		},
		{
			name:    "successful token generation for teacher",
			userID:  123,
			role:    "teacher",
			wantErr: false,
		},
		{
			name:    "successful token generation for student",
			userID:  456,
			role:    "student",
			wantErr: false,
		},
		{
			name:    "successful token generation for parent",
			userID:  789,
			role:    "parent",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.role, cfg)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify token can be parsed
				parsedToken, parseErr := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(cfg.Secret), nil
				})
				assert.NoError(t, parseErr)
				assert.True(t, parsedToken.Valid)

				// Verify claims
				if claims, ok := parsedToken.Claims.(*Claims); ok {
					assert.Equal(t, tt.userID, claims.UserID)
					assert.Equal(t, tt.role, claims.Role)
					assert.True(t, claims.ExpiresAt.After(time.Now()))
					assert.True(t, claims.IssuedAt.Before(time.Now().Add(time.Second)))
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	cfg := config.JWTConfig{
		Secret: "test-secret-key-for-testing",
	}

	tests := []struct {
		name      string
		tokenFunc func() string
		wantErr   bool
		wantUser  uint
		wantRole  string
	}{
		{
			name: "valid token",
			tokenFunc: func() string {
				token, _ := GenerateToken(123, "teacher", cfg)
				return token
			},
			wantErr:  false,
			wantUser: 123,
			wantRole: "teacher",
		},
		{
			name: "expired token",
			tokenFunc: func() string {
				claims := &Claims{
					UserID: 456,
					Role:   "student",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(cfg.Secret))
				return tokenString
			},
			wantErr: true,
		},
		{
			name: "invalid signature",
			tokenFunc: func() string {
				wrongCfg := config.JWTConfig{Secret: "wrong-secret"}
				token, _ := GenerateToken(789, "admin", wrongCfg)
				return token
			},
			wantErr: true,
		},
		{
			name: "malformed token",
			tokenFunc: func() string {
				return "not.a.valid.jwt.token"
			},
			wantErr: true,
		},
		{
			name: "empty token",
			tokenFunc: func() string {
				return ""
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.tokenFunc()
			claims, err := ValidateToken(token, cfg)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.wantUser, claims.UserID)
				assert.Equal(t, tt.wantRole, claims.Role)
			}
		})
	}
}

func TestClaims_Structure(t *testing.T) {
	// Test that Claims struct has required fields
	claims := &Claims{
		UserID: 123,
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	assert.Equal(t, uint(123), claims.UserID)
	assert.Equal(t, "admin", claims.Role)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
}

func TestToken_Expiration(t *testing.T) {
	cfg := config.JWTConfig{
		Secret: "test-secret",
	}

	// Generate token
	token, err := GenerateToken(100, "student", cfg)
	assert.NoError(t, err)

	// Validate immediately - should succeed
	claims, err := ValidateToken(token, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Check expiration time is ~24 hours from now
	expectedExpiry := time.Now().Add(24 * time.Hour)
	actualExpiry := claims.ExpiresAt.Time

	// Allow 5 second tolerance for test execution time
	timeDiff := actualExpiry.Sub(expectedExpiry).Abs()
	assert.Less(t, timeDiff, 5*time.Second, "Expiry time should be approximately 24 hours from now")
}

func TestToken_RolePreservation(t *testing.T) {
	cfg := config.JWTConfig{
		Secret: "test-secret",
	}

	roles := []string{"superadmin", "admin", "principal", "secretary", "teacher", "student", "parent"}

	for _, role := range roles {
		t.Run("role_"+role, func(t *testing.T) {
			token, err := GenerateToken(1, role, cfg)
			assert.NoError(t, err)

			claims, err := ValidateToken(token, cfg)
			assert.NoError(t, err)
			assert.Equal(t, role, claims.Role)
		})
	}
}
