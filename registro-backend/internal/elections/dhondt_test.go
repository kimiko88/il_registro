package elections

import (
	"testing"
)

func TestDhondtSeatAllocation(t *testing.T) {
	// 3 Lists, 5 total seats to assign
	// List A: 100 votes
	// List B: 80 votes
	// List C: 30 votes
	// Quotients:
	// List A: 100/1=100 (1st), 100/2=50 (3rd), 100/3=33.3 (5th)
	// List B: 80/1=80 (2nd), 80/2=40 (4th), 80/3=26.6
	// List C: 30/1=30
	// Seats allocated: List A: 3 seats, List B: 2 seats, List C: 0 seats. Total: 5

	lists := []ListVoteCount{
		{ListID: "list-a", ListNumber: 1, Motto: "Insieme per la Scuola", TotalVotes: 100},
		{ListID: "list-b", ListNumber: 2, Motto: "Futuro e Innovazione", TotalVotes: 80},
		{ListID: "list-c", ListNumber: 3, Motto: "Partecipazione Attiva", TotalVotes: 30},
	}

	seatsToAllocate := 5
	results := AllocateSeatsDhondt(lists, seatsToAllocate)

	seatMap := make(map[string]int)
	totalAllocated := 0
	for _, r := range results {
		seatMap[r.ListID] = r.SeatsWon
		totalAllocated += r.SeatsWon
	}

	if totalAllocated != seatsToAllocate {
		t.Fatalf("expected %d total seats allocated, got %d", seatsToAllocate, totalAllocated)
	}

	if seatMap["list-a"] != 3 {
		t.Errorf("expected List A to win 3 seats, got %d", seatMap["list-a"])
	}
	if seatMap["list-b"] != 2 {
		t.Errorf("expected List B to win 2 seats, got %d", seatMap["list-b"])
	}
	if seatMap["list-c"] != 0 {
		t.Errorf("expected List C to win 0 seats, got %d", seatMap["list-c"])
	}
}

func TestGenerateAnonymousReceiptToken(t *testing.T) {
	token := GenerateReceiptToken()
	if len(token) != 32 {
		t.Errorf("expected 32 hex chars receipt token, got %d", len(token))
	}

	token2 := GenerateReceiptToken()
	if token == token2 {
		t.Errorf("tokens must be uniquely generated")
	}
}
