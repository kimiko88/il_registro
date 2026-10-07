package maturita

import (
	"testing"
)

func BenchmarkCalculateYearCredits(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = CalculateYearCredits(5, 8.4, true)
	}
}

func BenchmarkCalculateFinalExamScores(b *testing.B) {
	credits := TrienniumCredits{
		Credit3rd:    12,
		Credit4th:    13,
		Credit5th:    15,
		TotalCredits: 40,
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = CalculateFinalExamScores(credits, 18.0, 19.0, 20.0, 2, false)
	}
}
