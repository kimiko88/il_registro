package parents

import (
	"context"
	"testing"
)

func TestAccessGuard_CustodyRestrictions(t *testing.T) {
	guard := NewAccessGuard(newMockDualSignRepo())
	ctx := context.Background()

	// 1. Shared custody -> allowed
	allowed, err := guard.CanAccessStudent(ctx, "parent-user-1", "student-1")
	if err != nil || !allowed {
		t.Errorf("expected parent 1 to have access to student 1 with shared custody")
	}

	// 2. Restricted custody (e.g. court order per Tribunale per i Minorenni) -> blocked
	allowed, err = guard.CanAccessStudent(ctx, "parent-restricted", "student-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Errorf("expected restricted parent to be denied access to student data")
	}

	// 3. Can authorize activities check
	canAuth, err := guard.CanAuthorizeActivity(ctx, "parent-user-1", "student-1")
	if err != nil || !canAuth {
		t.Errorf("expected parent-user-1 to be able to authorize activities")
	}

	canAuth, err = guard.CanAuthorizeActivity(ctx, "parent-restricted", "student-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if canAuth {
		t.Errorf("expected restricted parent to NOT be able to authorize activities")
	}
}
