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
		userID:   "user-999",
		email:    "old@school.com",
		role:     "student",
		schoolID: "school-1",
		exp:      time.Now().Add(-1 * time.Second),
	}
	store.mu.Unlock()

	_, _, _, _, ok := store.Consume("expired-ticket")
	if ok {
		t.Errorf("expected consume to fail for expired ticket")
	}
}
