package cache

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrCacheMiss indicates that the requested key does not exist or has expired.
var ErrCacheMiss = errors.New("cache: key not found")

// Cache defines the standard interface for key-value caching across the application.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeletePrefix(ctx context.Context, prefix string) error
	Close() error
}

// ─── Memory Cache ─────────────────────────────────────────────────────────────

type memoryItem struct {
	value     string
	expiresAt time.Time
}

type MemoryCache struct {
	mu      sync.RWMutex
	items   map[string]memoryItem
	stopJan chan struct{}
}

// NewMemoryCache creates an in-memory thread-safe cache with automatic expired item eviction.
func NewMemoryCache() *MemoryCache {
	mc := &MemoryCache{
		items:   make(map[string]memoryItem),
		stopJan: make(chan struct{}),
	}

	go mc.janitor(2 * time.Minute)
	return mc
}

func (m *MemoryCache) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			m.mu.Lock()
			for k, v := range m.items {
				if !v.expiresAt.IsZero() && v.expiresAt.Before(now) {
					delete(m.items, k)
				}
			}
			m.mu.Unlock()
		case <-m.stopJan:
			return
		}
	}
}

func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	item, ok := m.items[key]
	m.mu.RUnlock()

	if !ok {
		return "", ErrCacheMiss
	}

	if !item.expiresAt.IsZero() && item.expiresAt.Before(time.Now()) {
		m.mu.Lock()
		delete(m.items, key)
		m.mu.Unlock()
		return "", ErrCacheMiss
	}

	return item.value, nil
}

func (m *MemoryCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	m.mu.Lock()
	m.items[key] = memoryItem{
		value:     value,
		expiresAt: expiresAt,
	}
	m.mu.Unlock()
	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, keys ...string) error {
	m.mu.Lock()
	for _, k := range keys {
		delete(m.items, k)
	}
	m.mu.Unlock()
	return nil
}

func (m *MemoryCache) DeletePrefix(ctx context.Context, prefix string) error {
	m.mu.Lock()
	for k := range m.items {
		if strings.HasPrefix(k, prefix) {
			delete(m.items, k)
		}
	}
	m.mu.Unlock()
	return nil
}

func (m *MemoryCache) Close() error {
	close(m.stopJan)
	m.mu.Lock()
	m.items = make(map[string]memoryItem)
	m.mu.Unlock()
	return nil
}

// ─── Redis Cache ──────────────────────────────────────────────────────────────

type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a Redis-backed Cache implementation.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrCacheMiss
		}
		return "", err
	}
	return val, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisCache) DeletePrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	matchPattern := prefix + "*"

	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, matchPattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

// ─── Factory ──────────────────────────────────────────────────────────────────

// NewCache instantiates RedisCache if redisURL is provided and reachable;
// otherwise falls back to a thread-safe MemoryCache.
func NewCache(redisURL string) Cache {
	if strings.TrimSpace(redisURL) != "" {
		opt, err := redis.ParseURL(redisURL)
		if err == nil {
			client := redis.NewClient(opt)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			if err := client.Ping(ctx).Err(); err == nil {
				log.Println("cache: initialized with Redis backend")
				return NewRedisCache(client)
			}
			log.Printf("cache: Redis ping failed (%v), falling back to in-memory cache", err)
			_ = client.Close()
		} else {
			log.Printf("cache: invalid REDIS_URL (%v), falling back to in-memory cache", err)
		}
	}

	log.Println("cache: initialized with thread-safe in-memory backend")
	return NewMemoryCache()
}
