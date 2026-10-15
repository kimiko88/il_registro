package interpelli

import (
	"math"
	"testing"
)

// FuzzCalculateScore fuzz-tests scoring logic across arbitrary numeric values.
func FuzzCalculateScore(f *testing.F) {
	f.Add(110.0, true, true, 12, "C1", 4)
	f.Add(60.0, false, false, 0, "", 0)
	f.Add(75.5, false, true, 24, "B2", 1)
	f.Add(-10.0, false, false, -5, "INVALID", -1)
	f.Add(999.0, true, false, 1000, "C2", 100)

	f.Fuzz(func(t *testing.T, grade float64, lode bool, hab bool, service int, lang string, dig int) {
		if math.IsNaN(grade) || math.IsInf(grade, 0) {
			return
		}
		deg, srv, certs, total := CalculateScore(grade, lode, hab, service, lang, dig)

		// Assert invariant constraints:
		// Degree points >= 12 (minimum for 60)
		if deg < 12.0 {
			t.Errorf("expected deg score >= 12, got %f", deg)
		}
		// Service capped at 36
		if srv < 0 || srv > 36.0 {
			t.Errorf("expected service score in [0, 36], got %f", srv)
		}
		// Digital certs capped at 2, plus max language cert 6 = max 8
		if certs < 0 || certs > 8.0 {
			t.Errorf("expected certs score in [0, 8], got %f", certs)
		}
		// Total must equal sum
		if math.Abs(total-(deg+srv+certs)) > 1e-6 {
			t.Errorf("total %f does not match sum %f", total, deg+srv+certs)
		}
	})
}
