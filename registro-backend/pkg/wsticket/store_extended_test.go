package wsticket

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWSTicketStore_Close_InMemory verifies Close is a no-op for in-memory store.
func TestWSTicketStore_Close_InMemory(t *testing.T) {
	store := NewStore()
	err := store.Close()
	assert.NoError(t, err)
}

// TestWSTicketStore_Consume_NeverIssued verifies consuming an unknown ticket returns false.
func TestWSTicketStore_Consume_NeverIssued(t *testing.T) {
	store := NewStore()
	_, _, _, _, ok := store.Consume("nonexistent-ticket-xyz")
	assert.False(t, ok)
}

// TestWSTicketStore_NewStore_NoURL verifies NewStore() without args creates in-memory store.
func TestWSTicketStore_NewStore_NoURL(t *testing.T) {
	store := NewStore()
	require.NotNil(t, store)
	assert.Nil(t, store.rdb, "in-memory store should not have a Redis client")
	assert.NotNil(t, store.tickets, "in-memory store should initialise the tickets map")
}

// TestWSTicketStore_NewStore_EmptyString verifies NewStore("") also creates in-memory store.
func TestWSTicketStore_NewStore_EmptyString(t *testing.T) {
	store := NewStore("")
	require.NotNil(t, store)
	assert.Nil(t, store.rdb)
}

// TestWSTicketStore_Concurrency verifies that concurrent Issue and Consume operations are race-free.
func TestWSTicketStore_Concurrency(t *testing.T) {
	store := NewStore()

	const n = 50
	tickets := make([]string, n)
	var mu sync.Mutex

	var wg sync.WaitGroup
	// Issue n tickets concurrently
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			ticket, err := store.Issue("user-concurrent", "c@test.it", "teacher", "s1")
			require.NoError(t, err)
			mu.Lock()
			tickets = append(tickets, ticket)
			mu.Unlock()
		}()
	}
	wg.Wait()

	// Consume all tickets concurrently
	wg.Add(n)
	for i := 0; i < n; i++ {
		idx := i
		go func() {
			defer wg.Done()
			mu.Lock()
			if idx < len(tickets) && tickets[idx] != "" {
				ticket := tickets[idx]
				mu.Unlock()
				store.Consume(ticket)
			} else {
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
}

// TestWSTicketStore_OneShotConsumption verifies a ticket can only be consumed once.
func TestWSTicketStore_OneShotConsumption(t *testing.T) {
	store := NewStore()
	ticket, err := store.Issue("user-1", "a@test.it", "student", "school-1")
	require.NoError(t, err)

	_, _, _, _, ok1 := store.Consume(ticket)
	assert.True(t, ok1, "first consume should succeed")

	_, _, _, _, ok2 := store.Consume(ticket)
	assert.False(t, ok2, "second consume of same ticket should fail")
}

// TestWSTicketStore_InvalidRedisURL_FallsBackToMemory verifies that an invalid Redis URL
// falls back to in-memory store and everything still works.
func TestWSTicketStore_InvalidRedisURL_FallsBackToMemory(t *testing.T) {
	store := NewStore("not-a-valid-url://???")
	require.NotNil(t, store)

	ticket, err := store.Issue("user-fallback", "fallback@test.it", "admin", "school-99")
	require.NoError(t, err)
	require.NotEmpty(t, ticket)

	uid, email, role, schoolID, ok := store.Consume(ticket)
	require.True(t, ok)
	assert.Equal(t, "user-fallback", uid)
	assert.Equal(t, "fallback@test.it", email)
	assert.Equal(t, "admin", role)
	assert.Equal(t, "school-99", schoolID)
}

// TestWSTicketStore_AllFieldsPreserved verifies all claim fields are preserved through issue+consume.
func TestWSTicketStore_AllFieldsPreserved(t *testing.T) {
	store := NewStore()

	cases := []struct {
		userID   string
		email    string
		role     string
		schoolID string
	}{
		{"u1", "a@b.com", "teacher", "s1"},
		{"u2", "b@c.com", "student", "s2"},
		{"u3", "c@d.com", "parent", ""},
		{"u4", "d@e.com", "admin", "s4"},
	}

	for _, tc := range cases {
		t.Run(tc.role, func(t *testing.T) {
			ticket, err := store.Issue(tc.userID, tc.email, tc.role, tc.schoolID)
			require.NoError(t, err)

			uid, email, role, schoolID, ok := store.Consume(ticket)
			require.True(t, ok)
			assert.Equal(t, tc.userID, uid)
			assert.Equal(t, tc.email, email)
			assert.Equal(t, tc.role, role)
			assert.Equal(t, tc.schoolID, schoolID)
		})
	}
}

// TestWSTicketStore_TicketLengthIs64 verifies the ticket is a 64-char hex string (32 random bytes).
func TestWSTicketStore_TicketLengthIs64(t *testing.T) {
	store := NewStore()
	ticket, err := store.Issue("u", "u@t.it", "teacher", "s")
	require.NoError(t, err)
	assert.Len(t, ticket, 64)
}

// TestWSTicketStore_MultipleTicketsUnique verifies that two issued tickets are different.
func TestWSTicketStore_MultipleTicketsUnique(t *testing.T) {
	store := NewStore()
	t1, _ := store.Issue("u1", "a@b.com", "teacher", "s1")
	t2, _ := store.Issue("u1", "a@b.com", "teacher", "s1")
	assert.NotEqual(t, t1, t2, "tickets must be unique")
}

// TestWSTicketStore_TicketJustExpired verifies that a ticket with TTL already elapsed is rejected.
func TestWSTicketStore_TicketJustExpired(t *testing.T) {
	store := NewStore()

	// Manually insert an already-expired ticket
	store.mu.Lock()
	store.tickets["expired-now"] = entry{
		UserID:   "u",
		Email:    "e@t.it",
		Role:     "teacher",
		SchoolID: "s",
		Exp:      time.Now().Add(-1 * time.Millisecond),
	}
	store.mu.Unlock()

	_, _, _, _, ok := store.Consume("expired-now")
	assert.False(t, ok, "expired ticket must not be consumed")
}
