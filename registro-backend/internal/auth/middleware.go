package auth

import (
	"context"
	"net/http"
	"strings"

	"registro-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Middleware provides authentication middleware
type Middleware struct {
	tokenManager *jwt.TokenManager
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(tokenManager *jwt.TokenManager) *Middleware {
	return &Middleware{
		tokenManager: tokenManager,
	}
}

// Authenticate validates JWT token and sets user context
func (m *Middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing authorization header"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := m.tokenManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or expired token"})
			c.Abort()
			return
		}

		// Set claims in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("school_id", claims.SchoolID)

		c.Next()
	}
}

// RequireRole checks if user has required role
func (m *Middleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "role not found in context"})
			c.Abort()
			return
		}

		roleStr := role.(string)
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole || roleStr == "superadmin" {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, ErrorResponse{Error: "insufficient permissions"})
		c.Abort()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	return role.(string), true
}

// GetSchoolID extracts school ID from context
func GetSchoolID(c *gin.Context) (string, bool) {
	schoolID, exists := c.Get("school_id")
	if !exists {
		return "", false
	}
	return schoolID.(string), true
}

// ContextKey type for context keys
type ContextKey string

const (
	UserIDKey   ContextKey = "user_id"
	EmailKey    ContextKey = "email"
	RoleKey     ContextKey = "role"
	SchoolIDKey ContextKey = "school_id"
)

// SetUserContext sets user information in context
func SetUserContext(ctx context.Context, userID, email, role, schoolID string) context.Context {
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, EmailKey, email)
	ctx = context.WithValue(ctx, RoleKey, role)
	ctx = context.WithValue(ctx, SchoolIDKey, schoolID)
	return ctx
}
