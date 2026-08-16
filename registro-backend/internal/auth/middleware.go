package auth

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"registro-backend/internal/users"
	"registro-backend/pkg/jwt"
	"registro-backend/pkg/logger"
	"registro-backend/pkg/wsticket"

	"github.com/gin-gonic/gin"
)

type activeCacheEntry struct {
	isActive  bool
	expiresAt time.Time
}

// Middleware provides authentication middleware
type Middleware struct {
	tokenManager *jwt.TokenManager
	userRepo     users.Repository
	activeCache  sync.Map
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

		// Try Header
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
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

		// Verify the account is still active in the DB with a short 10s TTL cache.
		var isActive bool
		var foundInCache bool
		if val, ok := m.activeCache.Load(claims.UserID); ok {
			entry := val.(activeCacheEntry)
			if time.Now().Before(entry.expiresAt) {
				isActive = entry.isActive
				foundInCache = true
			} else {
				m.activeCache.Delete(claims.UserID)
			}
		}
		if !foundInCache {
			var err error
			isActive, err = m.userRepo.IsActive(c.Request.Context(), claims.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "auth check failed"})
				c.Abort()
				return
			}
			m.activeCache.Store(claims.UserID, activeCacheEntry{
				isActive:  isActive,
				expiresAt: time.Now().Add(10 * time.Second),
			})
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
		c.Set("is_staff", IsStaffRole(claims.Role))

		c.Set("locale", parseAcceptLanguage(c.GetHeader("Accept-Language")))

		// Synchronize with stdlib request context
		c.Request = c.Request.WithContext(SetUserContext(c.Request.Context(), claims.UserID, claims.Email, claims.Role, claims.SchoolID))

		c.Next()
	}
}

// InvalidateUserActiveCache clears the cached active status for a user.
func (m *Middleware) InvalidateUserActiveCache(userID string) {
	m.activeCache.Delete(userID)
}

// CleanupExpiredEntries iterates through activeCache and deletes entries past their expiration time.
func (m *Middleware) CleanupExpiredEntries() {
	now := time.Now()
	m.activeCache.Range(func(key, value any) bool {
		if entry, ok := value.(activeCacheEntry); ok {
			if now.After(entry.expiresAt) {
				m.activeCache.Delete(key)
			}
		}
		return true
	})
}

// StartCacheCleaner launches a background goroutine to periodically clean up expired activeCache entries until context cancellation.
func (m *Middleware) StartCacheCleaner(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.CleanupExpiredEntries()
			}
		}
	}()
}

// AuthenticateWSTicket validates a single-use opaque WS ticket for WebSocket upgrades.
// It prevents JWT tokens from appearing in query parameters or Nginx/proxy access logs.
func (m *Middleware) AuthenticateWSTicket(store *wsticket.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if store == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "ws ticket store unconfigured"})
			c.Abort()
			return
		}
		ticket := c.Query("ticket")
		if ticket == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing ws ticket"})
			c.Abort()
			return
		}

		userID, email, role, schoolID, ok := store.Consume(ticket)
		if !ok {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or expired ws ticket"})
			c.Abort()
			return
		}

		// Verify the account is still active in DB
		isActive, err := m.userRepo.IsActive(c.Request.Context(), userID)
		if err != nil || !isActive {
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "account is disabled"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("role", role)
		c.Set("school_id", schoolID)
		c.Set("is_staff", IsStaffRole(role))
		c.Set("locale", parseAcceptLanguage(c.GetHeader("Accept-Language")))
		c.Request = c.Request.WithContext(SetUserContext(c.Request.Context(), userID, email, role, schoolID))

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
	val, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	s, ok := val.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}

var bcp47Regex = regexp.MustCompile(`^[a-zA-Z]{2,3}(-[a-zA-Z0-9]{2,8})*$`)

// parseAcceptLanguage normalizes Accept-Language header value into a safe BCP-47 locale.
func parseAcceptLanguage(raw string) string {
	if raw == "" {
		return "it-IT"
	}
	first := strings.Split(raw, ",")[0]
	first = strings.TrimSpace(strings.Split(first, ";")[0])
	first = strings.ReplaceAll(first, "_", "-")

	if bcp47Regex.MatchString(first) {
		return first
	}
	return "it-IT"
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	val, exists := c.Get("role")
	if !exists {
		return "", false
	}
	s, ok := val.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
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
