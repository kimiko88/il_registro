package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipEntry associates a rate limiter with its last access time for TTL eviction.
type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter tracks per-IP rate limiters with automatic cleanup to prevent memory leaks.
//
// ⚠️  MULTI-INSTANCE WARNING: This implementation is in-memory and is NOT shared
// across multiple server replicas. If the service runs with more than one instance
// (e.g. Docker Swarm, Kubernetes, Render with multiple workers) an attacker can
// round-robin requests across replicas and multiply the effective rate limit by N.
// To fix this for multi-instance deployments, replace this struct with a Redis-backed
// implementation using the go-redis/redis_rate package:
//   https://github.com/go-redis/redis_rate
type IPRateLimiter struct {
	ips map[string]*ipEntry
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*ipEntry),
		r:   r,
		b:   b,
	}

	// Cleanup goroutine: every 5 minutes, evict IPs not seen in the last 10 minutes.
	// This prevents unbounded map growth under scanner/bot traffic.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			i.mu.Lock()
			cutoff := time.Now().Add(-10 * time.Minute)
			for ip, entry := range i.ips {
				if entry.lastSeen.Before(cutoff) {
					delete(i.ips, ip)
				}
			}
			i.mu.Unlock()
		}
	}()

	return i
}

// GetLimiter returns (or creates) the rate limiter for the given IP, updating its lastSeen.
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	// Fast path: read lock
	i.mu.RLock()
	entry, exists := i.ips[ip]
	i.mu.RUnlock()

	if exists {
		i.mu.Lock()
		entry.lastSeen = time.Now()
		i.mu.Unlock()
		return entry.limiter
	}

	// Slow path: create new entry
	i.mu.Lock()
	defer i.mu.Unlock()

	// Double-check after acquiring write lock
	if entry, exists = i.ips[ip]; exists {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = &ipEntry{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	// 5 requests per second, burst of 10.
	// See IPRateLimiter doc comment for multi-instance limitations.
	limiter := NewIPRateLimiter(5, 10)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = "127.0.0.1"
		}
		if !limiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
