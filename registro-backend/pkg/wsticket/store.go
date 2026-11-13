package wsticket

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const ticketTTL = 30 * time.Second

// consumeLuaScript atomically reads and deletes a ticket in a single step (one-shot).
const consumeLuaScript = `
local val = redis.call('GET', KEYS[1])
if val then
    redis.call('DEL', KEYS[1])
end
return val
`

type entry struct {
	UserID   string    `json:"user_id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	SchoolID string    `json:"school_id"`
	Exp      time.Time `json:"exp"`
}

// Store is a thread-safe store for one-time WebSocket authentication tickets.
// It supports Redis-backed distributed tickets; memory mode is used only when no Redis URL is configured.
type Store struct {
	mu      sync.Mutex
	tickets map[string]entry
	rdb     *redis.Client
	initErr error
}

// NewStore creates a new WS ticket store.
// If redisURL is provided, Redis must be reachable. A configured but unavailable
// Redis backend fails closed rather than silently switching to process-local tickets.
// In-memory mode is reserved for an explicitly empty URL (local development).
func NewStore(redisURL ...string) *Store {
	if len(redisURL) > 0 && strings.TrimSpace(redisURL[0]) != "" {
		opt, err := redis.ParseURL(redisURL[0])
		if err != nil {
			log.Printf("wsticket: invalid REDIS_URL; ticket store unavailable: %v", err)
			return &Store{initErr: err}
		}
		client := redis.NewClient(opt)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if pingErr := client.Ping(ctx).Err(); pingErr != nil {
			log.Printf("wsticket: Redis ping failed; ticket store unavailable: %v", pingErr)
			_ = client.Close()
			return &Store{initErr: pingErr}
		}
		log.Println("wsticket: connected to Redis backend for distributed ticket validation")
		return &Store{rdb: client}
	}

	s := &Store{tickets: make(map[string]entry)}
	go s.sweep()
	return s
}

// Issue generates a 30-second opaque one-time ticket for the given authenticated claims.
func (s *Store) Issue(userID, email, role, schoolID string) (string, error) {
	if s.initErr != nil {
		return "", s.initErr
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	ticket := hex.EncodeToString(b)
	e := entry{
		UserID:   userID,
		Email:    email,
		Role:     role,
		SchoolID: schoolID,
		Exp:      time.Now().Add(ticketTTL),
	}

	if s.rdb != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		payload, err := json.Marshal(e)
		if err != nil {
			return "", err
		}
		if err := s.rdb.Set(ctx, "wsticket:"+ticket, payload, ticketTTL).Err(); err != nil {
			log.Printf("wsticket: failed to store ticket in Redis: %v", err)
			return "", err
		}
		return ticket, nil
	}

	s.mu.Lock()
	s.tickets[ticket] = e
	s.mu.Unlock()
	return ticket, nil
}

// Consume validates and consumes a ticket (one-shot: deleted on first use or if expired).
func (s *Store) Consume(ticket string) (userID, email, role, schoolID string, ok bool) {
	if s.initErr != nil || ticket == "" {
		return "", "", "", "", false
	}

	if s.rdb != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		res, err := s.rdb.Eval(ctx, consumeLuaScript, []string{"wsticket:" + ticket}).Result()
		if err != nil || res == nil {
			return "", "", "", "", false
		}

		rawJSON, isStr := res.(string)
		if !isStr || rawJSON == "" {
			return "", "", "", "", false
		}

		var e entry
		if err := json.Unmarshal([]byte(rawJSON), &e); err != nil {
			return "", "", "", "", false
		}
		return e.UserID, e.Email, e.Role, e.SchoolID, true
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	e, exists := s.tickets[ticket]
	if !exists || time.Now().After(e.Exp) {
		delete(s.tickets, ticket)
		return "", "", "", "", false
	}
	delete(s.tickets, ticket) // One-shot consumption
	return e.UserID, e.Email, e.Role, e.SchoolID, true
}

// Close gracefully closes the underlying Redis client if active.
func (s *Store) Close() error {
	if s.rdb != nil {
		return s.rdb.Close()
	}
	return nil
}

// sweep periodically purges expired tickets to prevent memory leaks in in-memory mode.
func (s *Store) sweep() {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		s.mu.Lock()
		for k, v := range s.tickets {
			if now.After(v.Exp) {
				delete(s.tickets, k)
			}
		}
		s.mu.Unlock()
	}
}
