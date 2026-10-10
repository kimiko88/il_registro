package elections

import (
	"testing"
)

func FuzzAllocateSeatsDhondt(f *testing.F) {
	seeds := [][4]int{
		{100, 80, 30, 5},
		{500, 250, 100, 8},
		{0, 0, 0, 5},
		{10, 0, 0, 1},
		{100, 100, 100, 3},
		{1000, 2000, 3000, 10},
		{50, 50, 50, 0},
	}

	for _, s := range seeds {
		f.Add(s[0], s[1], s[2], s[3])
	}

	f.Fuzz(func(t *testing.T, v1, v2, v3, seats int) {
		if v1 < 0 || v2 < 0 || v3 < 0 || seats < 0 || seats > 100 {
			return
		}

		lists := []ListVoteCount{
			{ListID: "list-1", ListNumber: 1, Motto: "Motto 1", TotalVotes: v1},
			{ListID: "list-2", ListNumber: 2, Motto: "Motto 2", TotalVotes: v2},
			{ListID: "list-3", ListNumber: 3, Motto: "Motto 3", TotalVotes: v3},
		}

		results := AllocateSeatsDhondt(lists, seats)

		totalAllocated := 0
		for _, r := range results {
			if r.SeatsWon < 0 {
				t.Fatalf("negative seats assigned to list %s: %d", r.ListID, r.SeatsWon)
			}
			totalAllocated += r.SeatsWon
		}

		totalVotes := v1 + v2 + v3
		if totalVotes == 0 || seats == 0 {
			if totalAllocated != 0 {
				t.Fatalf("expected 0 seats when votes or seats are 0, got %d", totalAllocated)
			}
		} else {
			if totalAllocated != seats {
				t.Fatalf("expected total allocated seats %d, got %d", seats, totalAllocated)
			}
		}
	})
}
