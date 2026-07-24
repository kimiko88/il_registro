package auth

import (
	"context"
	"net/http"
	"strings"

	"registro-backend/internal/users"
	"registro-backend/pkg/jwt"
	"registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Middleware provides authentication middleware
type Middleware struct {
	tokenManager *jwt.TokenManager
	userRepo     users.Repository
}

// NewMiddleware creates a new auth middleware.
// userRepo must not be nil — it is required to verify that accounts are still
// active on every request. Passing nil will panic at startup to prevent silent
// security bypasses where disabled accounts could keep making requests.
func NewMiddleware(tokenManager *jwt.TokenManager, userRepo users.Repository) *Middleware {
	if userRepo == nil {
		panic("auth.NewMiddleware: userRepo must not be nil — required for active-account checks")
	}
	return &Middleware{
		tokenManager: tokenManager,
		userRepo:     userRepo,
	}
}

// Authenticate validates JWT token and sets user context.
// It also verifies that the account is still active in the database;
// a valid JWT for a disabled account is rejected with 403.
func (m *Middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authHeader := c.GetHeader("Authorization")

		// 1. Try Header
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 2. Try WebSocket subprotocol header (Sec-WebSocket-Protocol: access_token, <token>)
		if token == "" {
			secProto := c.GetHeader("Sec-WebSocket-Protocol")
			if secProto != "" {
				parts := strings.Split(secProto, ",")
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if p != "" && p != "access_token" && p != "bearer" {
						// Verify p is a 3-part JWT token (header.payload.signature)
						if strings.Count(p, ".") == 2 {
							token = p
							c.Header("Sec-WebSocket-Protocol", p)
							break
						}
					}
				}
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing authorization token"})
			c.Abort()
			return
		}

		claims, err := m.tokenManager.ValidateToken(token)
		if err != nil {
			logger.Log.Debugf("Token Validation Failed: %v", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or expired token"})
			c.Abort()
			return
		}

		// Always verify the account is still active in the DB.
		// This ensures that disabling a user takes effect within one request,
		// not just after the JWT expires (up to 15 minutes later).
		isActive, err := m.userRepo.IsActive(c.Request.Context(), claims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "auth check failed"})
			c.Abort()
			return
		}
		if !isActive {
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "account is disabled"})
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

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "invalid role format in context"})
			c.Abort()
			return
		}
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
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
	s, ok := schoolID.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
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
