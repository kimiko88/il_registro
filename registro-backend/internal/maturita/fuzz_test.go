package maturita

import (
	"math"
	"testing"
)

// FuzzCalculateYearCredits verifies credit boundaries across arbitrary average grades and years.
func FuzzCalculateYearCredits(f *testing.F) {
	f.Add(3, 7.5, true)
	f.Add(4, 6.0, false)
	f.Add(5, 10.0, true)
	f.Add(2, 8.0, true)
	f.Add(3, 5.9, false)
	f.Add(5, 11.0, true)

	f.Fuzz(func(t *testing.T, year int, avg float64, high bool) {
		if math.IsNaN(avg) || math.IsInf(avg, 0) {
			return
		}
		credits, err := CalculateYearCredits(year, avg, high)
		if year < 3 || year > 5 || avg < 6.0 || avg > 10.0 {
			if err == nil {
				t.Errorf("expected error for year=%d, avg=%f", year, avg)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for year=%d, avg=%f: %v", year, avg, err)
			}
			maxCredit := 12
			if year == 4 {
				maxCredit = 13
			} else if year == 5 {
				maxCredit = 15
			}
			if credits < 7 || credits > maxCredit {
				t.Errorf("credits %d out of bounds for year %d", credits, year)
			}
		}
	})
}
