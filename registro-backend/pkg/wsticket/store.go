package wsticket

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const ticketTTL = 30 * time.Second

type entry struct {
	userID   string
	email    string
	role     string
	schoolID string
	exp      time.Time
}

// Store is a thread-safe singleton store for one-time WebSocket authentication tickets.
type Store struct {
	mu      sync.Mutex
	tickets map[string]entry
}

// NewStore creates a new WS ticket store with background expiration cleanup.
func NewStore() *Store {
	s := &Store{tickets: make(map[string]entry)}
	go s.sweep()
	return s
}

// Issue generates a 30-second opaque one-time ticket for the given authenticated claims.
func (s *Store) Issue(userID, email, role, schoolID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	ticket := hex.EncodeToString(b)
	s.mu.Lock()
	s.tickets[ticket] = entry{
		userID:   userID,
		email:    email,
		role:     role,
		schoolID: schoolID,
		exp:      time.Now().Add(ticketTTL),
	}
	s.mu.Unlock()
	return ticket, nil
}

// Consume validates and consumes a ticket (one-shot: deleted on first use or if expired).
func (s *Store) Consume(ticket string) (userID, email, role, schoolID string, ok bool) {
	if ticket == "" {
		return "", "", "", "", false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, exists := s.tickets[ticket]
	if !exists || time.Now().After(e.exp) {
		delete(s.tickets, ticket)
		return "", "", "", "", false
	}
	delete(s.tickets, ticket) // One-shot consumption
	return e.userID, e.email, e.role, e.schoolID, true
}

// sweep periodically purges expired tickets to prevent memory leaks.
func (s *Store) sweep() {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		s.mu.Lock()
		for k, v := range s.tickets {
			if now.After(v.exp) {
				delete(s.tickets, k)
			}
		}
		s.mu.Unlock()
	}
}
