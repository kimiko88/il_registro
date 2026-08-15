package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	redis_rate "github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// ---------------------------------------------------------------------------
// Backend interface
// ---------------------------------------------------------------------------

// rateLimiterBackend abstracts the rate-limiting strategy so that the HTTP
// middleware functions (RateLimitMiddleware, AuthRateLimitMiddleware) are
// decoupled from the underlying implementation.
type rateLimiterBackend interface {
	allowGlobal(ctx context.Context, ip string) bool
	allowAuth(ctx context.Context, ip string) bool
}

var (
	activeBackend     rateLimiterBackend
	activeBackendOnce sync.Once
)

// InitRateLimiter configures the shared backend used by all rate-limit
// middleware. It must be called once during application startup.
//
//   - If redisURL is a valid Redis URL and the server is reachable,
//     a distributed Redis-backed limiter is used (safe for multi-instance).
//   - Otherwise a warning is logged and the in-memory fallback is used.
//
// If InitRateLimiter is never called, the first middleware invocation will
// trigger a lazy in-memory init (backward-compatible, single-instance only).
func InitRateLimiter(redisURL string) {
	activeBackendOnce.Do(func() {
		if redisURL != "" {
			opt, err := redis.ParseURL(redisURL)
			if err != nil {
				log.Printf("ratelimit: invalid REDIS_URL (%v) — falling back to in-memory limiter", err)
			} else {
				rdb := redis.NewClient(opt)
				pingCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				if pingErr := rdb.Ping(pingCtx).Err(); pingErr != nil {
					log.Printf("ratelimit: Redis ping failed (%v) — falling back to in-memory limiter", pingErr)
					_ = rdb.Close()
				} else {
					activeBackend = &redisBackend{
						limiter:     redis_rate.NewLimiter(rdb),
						globalLimit: redis_rate.PerSecond(5),
						authLimit:   redis_rate.PerMinute(5),
					}
					log.Println("ratelimit: using Redis-backed distributed rate limiter")
					return
				}
			}
		} else {
			log.Println("ratelimit: REDIS_URL not set — using in-memory rate limiter " +
				"(not distributed; unsuitable for multi-instance deployments)")
		}
		activeBackend = newMemoryBackend()
	})
}

// getBackend returns the active backend, performing a lazy in-memory init if
// InitRateLimiter was never called (backward-compatible for tests).
func getBackend() rateLimiterBackend {
	InitRateLimiter("")
	return activeBackend
}

// ---------------------------------------------------------------------------
// Redis backend
// ---------------------------------------------------------------------------

type redisBackend struct {
	limiter     *redis_rate.Limiter
	globalLimit redis_rate.Limit
	authLimit   redis_rate.Limit
}

func (r *redisBackend) allowGlobal(ctx context.Context, ip string) bool {
	res, err := r.limiter.Allow(ctx, "rl:global:"+ip, r.globalLimit)
	if err != nil {
		// Fail open: Redis errors must not take down the API.
		log.Printf("ratelimit: redis error (global): %v — allowing request", err)
		return true
	}
	return res.Allowed > 0
}

func (r *redisBackend) allowAuth(ctx context.Context, ip string) bool {
	res, err := r.limiter.Allow(ctx, "rl:auth:"+ip, r.authLimit)
	if err != nil {
		log.Printf("ratelimit: redis error (auth): %v — allowing request", err)
		return true
	}
	return res.Allowed > 0
}

// ---------------------------------------------------------------------------
// In-memory backend (original implementation, unchanged)
// ---------------------------------------------------------------------------

// ipEntry associates a rate limiter with its last access time for TTL eviction.
type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter tracks per-IP rate limiters with automatic cleanup to prevent
// memory leaks. Used by the in-memory backend.
//
// ⚠️  MULTI-INSTANCE WARNING: This implementation is not shared across replicas.
// Use the Redis backend (InitRateLimiter with a valid REDIS_URL) for deployments
// with more than one server instance.
type IPRateLimiter struct {
	ips      map[string]*ipEntry
	mu       sync.RWMutex
	r        rate.Limit
	b        int
	stopChan chan struct{}
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips:      make(map[string]*ipEntry),
		r:        r,
		b:        b,
		stopChan: make(chan struct{}),
	}
	// Cleanup goroutine: every 5 minutes, evict IPs not seen in the last 10 minutes.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				i.mu.Lock()
				cutoff := time.Now().Add(-10 * time.Minute)
				for ip, entry := range i.ips {
					if entry.lastSeen.Before(cutoff) {
						delete(i.ips, ip)
					}
				}
				i.mu.Unlock()
			case <-i.stopChan:
				return
			}
		}
	}()
	return i
}

func (i *IPRateLimiter) Close() {
	select {
	case <-i.stopChan:
	default:
		close(i.stopChan)
	}
}

// GetLimiter returns (or creates) the rate limiter for the given IP.
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()
	if entry, exists := i.ips[ip]; exists {
		entry.lastSeen = time.Now()
		return entry.limiter
	}
	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = &ipEntry{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

type memoryBackend struct {
	global *IPRateLimiter
	auth   *IPRateLimiter
}

func newMemoryBackend() *memoryBackend {
	return &memoryBackend{
		global: NewIPRateLimiter(5, 10),
		auth:   NewIPRateLimiter(rate.Every(12*time.Second), 5),
	}
}

func (m *memoryBackend) allowGlobal(_ context.Context, ip string) bool {
	return m.global.GetLimiter(ip).Allow()
}

func (m *memoryBackend) allowAuth(_ context.Context, ip string) bool {
	return m.auth.GetLimiter(ip).Allow()
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func resolveClientIP(c *gin.Context) string {
	var ip string
	if os.Getenv("TRUST_PROXY_HEADERS") == "true" {
		ip = c.ClientIP()
	} else {
		ip = c.RemoteIP()
	}
	if ip == "" {
		ip = "unknown"
	}
	return ip
}

// ---------------------------------------------------------------------------
// Public middleware constructors (signatures unchanged)
// ---------------------------------------------------------------------------

// RateLimitMiddleware limits each IP to 5 req/s with a burst of 10.
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := resolveClientIP(c)
		if !getBackend().allowGlobal(c.Request.Context(), ip) {
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AuthRateLimitMiddleware applies stricter limits for authentication endpoints
// (5 requests per minute, burst of 5) to prevent brute-force attacks.
func AuthRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := resolveClientIP(c)
		if !getBackend().allowAuth(c.Request.Context(), ip) {
			c.Header("Retry-After", "60")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":  "AUTH_RATE_LIMIT_EXCEEDED",
				"error": "Troppi tentativi di accesso. Riprova tra un minuto.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
