package middleware

import (
	"registro-backend/internal/auth"
	"registro-backend/internal/users"
	"registro-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware wraps internal/auth middleware for convenience in the middleware package.
type AuthMiddleware struct {
	inner *auth.Middleware
}

// NewAuthMiddleware creates a new JWT authentication middleware instance.
func NewAuthMiddleware(tokenManager *jwt.TokenManager, userRepo users.Repository) *AuthMiddleware {
	return &AuthMiddleware{
		inner: auth.NewMiddleware(tokenManager, userRepo),
	}
}

// Authenticate returns the Gin HandlerFunc for JWT authentication.
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return m.inner.Authenticate()
}

// RequireRole returns the Gin HandlerFunc for role enforcement.
func (m *AuthMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return m.inner.RequireRole(allowedRoles...)
}
