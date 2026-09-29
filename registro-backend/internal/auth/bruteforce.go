package auth

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// LoginRateLimiter defines the interface for tracking and enforcing login rate limits.
type LoginRateLimiter interface {
	Check(ctx context.Context, email, ipAddress string) error
	RecordFailure(ctx context.Context, email, ipAddress string)
	RecordSuccess(ctx context.Context, email, ipAddress string)
}

const (
	MaxLoginAttemptsPerIP    = 5
	WindowLoginAttemptsPerIP = 15 * time.Minute

	MaxLoginAttemptsPerEmail    = 20
	WindowLoginAttemptsPerEmail = 1 * time.Hour
)

// RedisLoginRateLimiter implements atomic, cluster-wide login brute force protection in memory.
// It shields PostgreSQL from thousands of queries per second during dictionary or credential stuffing attacks.
type RedisLoginRateLimiter struct {
	rdb *redis.Client
}

func NewRedisLoginRateLimiter(rdb *redis.Client) *RedisLoginRateLimiter {
	return &RedisLoginRateLimiter{rdb: rdb}
}

func (r *RedisLoginRateLimiter) Check(ctx context.Context, email, ipAddress string) error {
	if r == nil || r.rdb == nil {
		return nil
	}

	ipKey := fmt.Sprintf("login:fail:ip:%s", ipAddress)
	ipAttempts, err := r.rdb.Get(ctx, ipKey).Int()
	if err == nil && ipAttempts >= MaxLoginAttemptsPerIP {
		return ErrTooManyAttempts
	}

	if email != "" {
		emailKey := fmt.Sprintf("login:fail:email:%s", strings.ToLower(strings.TrimSpace(email)))
		emailAttempts, err := r.rdb.Get(ctx, emailKey).Int()
		if err == nil && emailAttempts >= MaxLoginAttemptsPerEmail {
			return ErrTooManyAttempts
		}
	}

	return nil
}

func (r *RedisLoginRateLimiter) RecordFailure(ctx context.Context, email, ipAddress string) {
	if r == nil || r.rdb == nil {
		return
	}

	// Increment IP counter and set TTL if first failure
	ipKey := fmt.Sprintf("login:fail:ip:%s", ipAddress)
	pipe := r.rdb.Pipeline()
	pipe.Incr(ctx, ipKey)
	pipe.Expire(ctx, ipKey, WindowLoginAttemptsPerIP)

	var emailKey string
	if email != "" {
		emailKey = fmt.Sprintf("login:fail:email:%s", strings.ToLower(strings.TrimSpace(email)))
		pipe.Incr(ctx, emailKey)
		pipe.Expire(ctx, emailKey, WindowLoginAttemptsPerEmail)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("loginLimiter: failed to record failure in Redis: %v", err)
	}
}

func (r *RedisLoginRateLimiter) RecordSuccess(ctx context.Context, email, ipAddress string) {
	if r == nil || r.rdb == nil {
		return
	}

	pipe := r.rdb.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("login:fail:ip:%s", ipAddress))
	if email != "" {
		pipe.Del(ctx, fmt.Sprintf("login:fail:email:%s", strings.ToLower(strings.TrimSpace(email))))
	}
	_, _ = pipe.Exec(ctx)
}

// NewLoginRateLimiter creates a Redis-backed login limiter or returns nil if redisURL is empty/unreachable.
func NewLoginRateLimiter(redisURL string) LoginRateLimiter {
	if strings.TrimSpace(redisURL) == "" {
		return nil
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil
	}

	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("loginLimiter: Redis ping failed (%v), falling back to DB rate limiting", err)
		_ = client.Close()
		return nil
	}

	log.Println("loginLimiter: connected to Redis for distributed brute-force protection")
	return NewRedisLoginRateLimiter(client)
}
