package auth

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"strconv"
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
	activeCache  map[string]activeCacheEntry
	cacheMu      sync.RWMutex
}

const maxActiveCacheSize = 10000

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
		activeCache:  make(map[string]activeCacheEntry),
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

		isActive, err := m.isAccountActive(c.Request.Context(), claims.UserID)
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
		c.Set("is_staff", IsStaffRole(claims.Role))

		c.Set("locale", parseAcceptLanguage(c.GetHeader("Accept-Language")))

		// Synchronize with stdlib request context
		c.Request = c.Request.WithContext(SetUserContext(c.Request.Context(), claims.UserID, claims.Email, claims.Role, claims.SchoolID))

		c.Next()
	}
}

// InvalidateUserActiveCache clears the cached active status for a user.
func (m *Middleware) InvalidateUserActiveCache(userID string) {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()
	if m.activeCache != nil {
		delete(m.activeCache, userID)
	}
}

// CleanupExpiredEntries iterates through activeCache and deletes entries past their expiration time.
func (m *Middleware) CleanupExpiredEntries() {
	now := time.Now()
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()
	for key, entry := range m.activeCache {
		if now.After(entry.expiresAt) {
			delete(m.activeCache, key)
		}
	}
}

func getActiveCacheTTL() time.Duration {
	ttlStr := os.Getenv("ACTIVE_ACCOUNT_CACHE_TTL_SECONDS")
	if ttlStr != "" {
		if sec, err := strconv.Atoi(ttlStr); err == nil && sec > 0 {
			return time.Duration(sec) * time.Second
		}
	}
	return 10 * time.Second
}

// isAccountActive checks whether a user account is active, utilizing a configurable TTL in-memory cache (default 10s).
func (m *Middleware) isAccountActive(ctx context.Context, userID string) (bool, error) {
	m.cacheMu.RLock()
	if m.activeCache != nil {
		if entry, ok := m.activeCache[userID]; ok {
			if time.Now().Before(entry.expiresAt) {
				m.cacheMu.RUnlock()
				return entry.isActive, nil
			}
		}
	}
	m.cacheMu.RUnlock()

	isActive, err := m.userRepo.IsActive(ctx, userID)
	if err != nil {
		return false, err
	}

	m.cacheMu.Lock()
	if m.activeCache == nil {
		m.activeCache = make(map[string]activeCacheEntry)
	}
	if len(m.activeCache) >= maxActiveCacheSize {
		now := time.Now()
		for k, v := range m.activeCache {
			if now.After(v.expiresAt) {
				delete(m.activeCache, k)
			}
		}
	}
	if len(m.activeCache) < maxActiveCacheSize {
		m.activeCache[userID] = activeCacheEntry{
			isActive:  isActive,
			expiresAt: time.Now().Add(getActiveCacheTTL()),
		}
	}
	m.cacheMu.Unlock()

	return isActive, nil
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

		isActive, err := m.isAccountActive(c.Request.Context(), userID)
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

// bcp47Regex is pre-compiled at package level for high performance on incoming requests.
var bcp47Regex = regexp.MustCompile(`^[a-zA-Z]{2,3}(-[a-zA-Z0-9]{2,8})*$`)

// parseAcceptLanguage normalizes Accept-Language header value into a safe BCP-47 locale.
func parseAcceptLanguage(raw string) string {
	if raw == "" {
		return "it-IT"
	}
	first := raw
	if commaIdx := strings.IndexByte(first, ','); commaIdx != -1 {
		first = first[:commaIdx]
	}
	if semiIdx := strings.IndexByte(first, ';'); semiIdx != -1 {
		first = first[:semiIdx]
	}
	first = strings.TrimSpace(first)
	first = strings.ReplaceAll(first, "_", "-")
	if len(first) > 35 {
		first = first[:35]
	}

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
