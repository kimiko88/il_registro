package jwt

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RevocationStore provides instant, distributed blacklisting of revoked JWT tokens.
type RevocationStore interface {
	Revoke(ctx context.Context, jti string, ttl time.Duration) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
	Close() error
}

const revocationKeyPrefix = "jwt:revoked:"

// RedisRevocationStore stores revoked JTIs in Redis with an auto-expiring TTL.
type RedisRevocationStore struct {
	rdb *redis.Client
}

func NewRedisRevocationStore(client *redis.Client) *RedisRevocationStore {
	return &RedisRevocationStore{rdb: client}
}

func (s *RedisRevocationStore) Revoke(ctx context.Context, jti string, ttl time.Duration) error {
	if jti == "" || ttl <= 0 {
		return nil
	}
	return s.rdb.Set(ctx, revocationKeyPrefix+jti, "1", ttl).Err()
}

func (s *RedisRevocationStore) IsRevoked(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}
	val, err := s.rdb.Exists(ctx, revocationKeyPrefix+jti).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (s *RedisRevocationStore) Close() error {
	if s.rdb != nil {
		return s.rdb.Close()
	}
	return nil
}

// MemoryRevocationStore is a thread-safe in-memory store used when Redis is unavailable.
type MemoryRevocationStore struct {
	mu      sync.RWMutex
	revoked map[string]time.Time
}

func NewMemoryRevocationStore() *MemoryRevocationStore {
	s := &MemoryRevocationStore{
		revoked: make(map[string]time.Time),
	}
	go s.cleanupLoop()
	return s
}

func (s *MemoryRevocationStore) Revoke(_ context.Context, jti string, ttl time.Duration) error {
	if jti == "" || ttl <= 0 {
		return nil
	}
	s.mu.Lock()
	s.revoked[jti] = time.Now().Add(ttl)
	s.mu.Unlock()
	return nil
}

func (s *MemoryRevocationStore) IsRevoked(_ context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}
	s.mu.RLock()
	expiry, exists := s.revoked[jti]
	s.mu.RUnlock()

	if !exists {
		return false, nil
	}
	if time.Now().After(expiry) {
		return false, nil
	}
	return true, nil
}

func (s *MemoryRevocationStore) Close() error {
	return nil
}

func (s *MemoryRevocationStore) cleanupLoop() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for jti, exp := range s.revoked {
			if now.After(exp) {
				delete(s.revoked, jti)
			}
		}
		s.mu.Unlock()
	}
}

// NewRevocationStore creates a RevocationStore connected to Redis, falling back to memory if unreachable.
func NewRevocationStore(redisURL string) RevocationStore {
	if strings.TrimSpace(redisURL) != "" {
		opt, err := redis.ParseURL(redisURL)
		if err == nil {
			client := redis.NewClient(opt)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			pingErr := client.Ping(ctx).Err()
			if pingErr == nil {
				log.Println("jwt: connected to Redis backend for distributed JTI token revocation")
				return NewRedisRevocationStore(client)
			}
			log.Printf("jwt: Redis ping failed (%v), falling back to in-memory revocation store", pingErr)
			_ = client.Close()
		}
	}
	return NewMemoryRevocationStore()
}
