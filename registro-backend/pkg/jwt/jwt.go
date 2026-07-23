package jwt

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// jwtAudience is the expected audience for all access tokens issued by this service.
const jwtAudience = "registro-api"

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	SchoolID string `json:"school_id,omitempty"`
	jwt.RegisteredClaims
}

// TokenManager handles JWT token operations
type TokenManager struct {
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewTokenManager creates a new token manager
func NewTokenManager(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *TokenManager {
	return &TokenManager{
		privateKey:      privateKey,
		publicKey:       publicKey,
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 7 * 24 * time.Hour,
	}
}

// GenerateAccessToken creates a new access token.
// The token is signed with RS256 and includes Issuer and Audience claims
// so that ValidateToken can reject tokens issued by other services.
func (tm *TokenManager) GenerateAccessToken(userID, email, role, schoolID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		SchoolID: schoolID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "registro-backend",
			// Audience prevents cross-service token reuse.
			Audience: jwt.ClaimStrings{jwtAudience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(tm.privateKey)
}

// GenerateRefreshToken creates a new refresh token.
// Refresh tokens do not carry an audience — they are only ever validated
// against the database, not via ValidateToken.
func (tm *TokenManager) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "registro-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(tm.privateKey)
}

// ValidateToken validates and parses a JWT access token.
// It verifies the signing algorithm (RS256), the issuer, and the audience
// to prevent accepting tokens issued by other services.
func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, ErrInvalidToken
			}
			return tm.publicKey, nil
		},
		// Validate issuer: only accept tokens we signed.
		jwt.WithIssuer("registro-backend"),
		// Validate audience: only accept tokens destined for this API.
		jwt.WithAudience(jwtAudience),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetAccessTokenTTL returns access token TTL in seconds
func (tm *TokenManager) GetAccessTokenTTL() int64 {
	return int64(tm.accessTokenTTL.Seconds())
}
