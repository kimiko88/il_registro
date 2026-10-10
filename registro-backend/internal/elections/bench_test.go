package elections

import (
	"fmt"
	"testing"
)

func BenchmarkAllocateSeatsDhondtSmall(b *testing.B) {
	lists := []ListVoteCount{
		{ListID: "list-a", ListNumber: 1, Motto: "Insieme per la Scuola", TotalVotes: 120},
		{ListID: "list-b", ListNumber: 2, Motto: "Futuro e Innovazione", TotalVotes: 95},
		{ListID: "list-c", ListNumber: 3, Motto: "Partecipazione Attiva", TotalVotes: 40},
	}
	seats := 5

	for b.Loop() {
		_ = AllocateSeatsDhondt(lists, seats)
	}
}

func BenchmarkAllocateSeatsDhondtLarge(b *testing.B) {
	lists := make([]ListVoteCount, 15)
	for i := 0; i < 15; i++ {
		lists[i] = ListVoteCount{
			ListID:     fmt.Sprintf("list-%d", i),
			ListNumber: i + 1,
			Motto:      fmt.Sprintf("Motto %d", i),
			TotalVotes: 50 + (i * 25),
		}
	}
	seats := 25

	for b.Loop() {
		_ = AllocateSeatsDhondt(lists, seats)
	}
}

func BenchmarkGenerateReceiptToken(b *testing.B) {
	for b.Loop() {
		_ = GenerateReceiptToken()
	}
}
