package wsticket

import (
	"testing"
	"time"
)

func TestWSTicketStore_IssueAndConsume(t *testing.T) {
	store := NewStore()

	ticket, err := store.Issue("user-123", "test@school.com", "teacher", "school-1")
	if err != nil {
		t.Fatalf("unexpected error issuing ticket: %v", err)
	}
	if len(ticket) != 64 { // 32 bytes hex encoded = 64 chars
		t.Fatalf("expected ticket length 64, got %d", len(ticket))
	}

	// First consume should succeed
	userID, email, role, schoolID, ok := store.Consume(ticket)
	if !ok {
		t.Fatalf("expected consume to succeed")
	}
	if userID != "user-123" || email != "test@school.com" || role != "teacher" || schoolID != "school-1" {
		t.Errorf("consumed claims mismatch: got (%s, %s, %s, %s)", userID, email, role, schoolID)
	}

	// Second consume should fail (one-shot)
	_, _, _, _, ok2 := store.Consume(ticket)
	if ok2 {
		t.Errorf("expected second consume to fail for used ticket")
	}
}

func TestWSTicketStore_ExpiredTicket(t *testing.T) {
	store := NewStore()

	// Manually insert an expired ticket
	store.mu.Lock()
	store.tickets["expired-ticket"] = entry{
		UserID:   "user-999",
		Email:    "old@school.com",
		Role:     "student",
		SchoolID: "school-1",
		Exp:      time.Now().Add(-1 * time.Second),
	}
	store.mu.Unlock()

	_, _, _, _, ok := store.Consume("expired-ticket")
	if ok {
		t.Errorf("expected consume to fail for expired ticket")
	}
}

func TestWSTicketStore_FallbackOnInvalidURL(t *testing.T) {
	// Should log warning and fall back to in-memory store gracefully
	store := NewStore("redis://invalid-host:6379")
	defer func() { _ = store.Close() }()

	ticket, err := store.Issue("user-456", "fallback@school.com", "admin", "school-2")
	if err != nil {
		t.Fatalf("expected issue to work on fallback store: %v", err)
	}

	userID, email, role, schoolID, ok := store.Consume(ticket)
	if !ok {
		t.Fatalf("expected consume to succeed on fallback store")
	}
	if userID != "user-456" || email != "fallback@school.com" || role != "admin" || schoolID != "school-2" {
		t.Errorf("claims mismatch on fallback store")
	}

	// Empty ticket consume
	_, _, _, _, okEmpty := store.Consume("")
	if okEmpty {
		t.Errorf("expected empty ticket consume to return false")
	}
}
