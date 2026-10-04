package parents

import (
	"context"
	"testing"
	"time"
)

func TestDualSignatureStateMachine(t *testing.T) {
	svc := NewDualSignatureService(newMockDualSignRepo())
	ctx := context.Background()

	// 1. Create a request for student with shared custody (both parents required)
	authReq, err := svc.CreateAuthorization(ctx, CreateDualAuthParams{
		StudentID:     "student-1",
		DocumentType:  "trip_consent",
		DocumentRefID: "trip-101",
		Title:         "Gita Scolastica a Firenze",
		Parent1ID:     "parent-user-1",
		Parent2ID:     "parent-user-2",
		IsShared:      true,
		Deadline:      time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("unexpected error creating dual auth: %v", err)
	}

	if authReq.Status != StatusPendingFirst {
		t.Errorf("expected initial status %s, got %s", StatusPendingFirst, authReq.Status)
	}

	// 2. Parent 1 signs with valid PIN
	updated, err := svc.SignDocument(ctx, authReq.ID, "parent-user-1", "1234")
	if err != nil {
		t.Fatalf("unexpected error on parent 1 sign: %v", err)
	}
	if updated.Status != StatusPendingSecond {
		t.Errorf("expected status %s after parent 1 sign, got %s", StatusPendingSecond, updated.Status)
	}
	if !updated.Parent1PinVerified || updated.Parent1SignedAt == nil {
		t.Errorf("expected parent 1 pin verified and signed at set")
	}

	// 3. Parent 2 signs with invalid PIN (must fail)
	_, err = svc.SignDocument(ctx, authReq.ID, "parent-user-2", "9999") // wrong pin
	if err == nil {
		t.Errorf("expected error with invalid PIN")
	}

	// 4. Parent 2 signs with valid PIN -> Completed
	completed, err := svc.SignDocument(ctx, authReq.ID, "parent-user-2", "5678")
	if err != nil {
		t.Fatalf("unexpected error on parent 2 sign: %v", err)
	}
	if completed.Status != StatusCompleted {
		t.Errorf("expected status %s after both sign, got %s", StatusCompleted, completed.Status)
	}
}

func TestDualSignature_SoleCustody_SingleSignCompletes(t *testing.T) {
	svc := NewDualSignatureService(newMockDualSignRepo())
	ctx := context.Background()

	// Sole custody (only parent 1 required)
	authReq, err := svc.CreateAuthorization(ctx, CreateDualAuthParams{
		StudentID:     "student-2",
		DocumentType:  "pdp_approval",
		DocumentRefID: "pdp-202",
		Title:         "Approvazione Piano Didattico Personalizzato",
		Parent1ID:     "parent-user-1",
		Parent2ID:     "",
		IsShared:      false,
		Deadline:      time.Now().Add(5 * 24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("error creating sole custody auth: %v", err)
	}

	// Parent 1 signs -> completes directly
	completed, err := svc.SignDocument(ctx, authReq.ID, "parent-user-1", "1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if completed.Status != StatusCompleted {
		t.Errorf("expected status %s for sole custody after 1 sign, got %s", StatusCompleted, completed.Status)
	}
}

func TestDualSignature_Rejection(t *testing.T) {
	svc := NewDualSignatureService(newMockDualSignRepo())
	ctx := context.Background()

	authReq, _ := svc.CreateAuthorization(ctx, CreateDualAuthParams{
		StudentID:     "student-1",
		DocumentType:  "trip_consent",
		DocumentRefID: "trip-102",
		Title:         "Viaggio d'istruzione",
		Parent1ID:     "parent-user-1",
		Parent2ID:     "parent-user-2",
		IsShared:      true,
	})

	rejected, err := svc.RejectDocument(ctx, authReq.ID, "parent-user-2", "Non autorizzo il viaggio")
	if err != nil {
		t.Fatalf("unexpected error on reject: %v", err)
	}
	if rejected.Status != StatusRejected {
		t.Errorf("expected status %s, got %s", StatusRejected, rejected.Status)
	}
}
