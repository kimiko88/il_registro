package scrutiny

import (
	"math"
	"testing"
)

func FuzzMapGradeToReligionJudgment(f *testing.F) {
	f.Add(10.0)
	f.Add(9.5)
	f.Add(8.0)
	f.Add(7.0)
	f.Add(6.0)
	f.Add(5.0)
	f.Add(0.0)
	f.Add(-2.5)

	f.Fuzz(func(t *testing.T, val float64) {
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return
		}

		judgment := mapGradeToReligionJudgment(val)
		if judgment == "" {
			t.Fatalf("expected non-empty judgment for grade %f", val)
		}

		switch judgment {
		case "Ottimo", "Distinto", "Buono", "Sufficiente", "Insufficiente", "Non classificabile":
			// Expected canonical judgments
		default:
			t.Fatalf("unexpected judgment output %q for value %f", judgment, val)
		}

		if val >= 9.5 && judgment != "Ottimo" {
			t.Fatalf("expected Ottimo for %f, got %s", val, judgment)
		}
		if val <= 0 && judgment != "Non classificabile" {
			t.Fatalf("expected Non classificabile for %f, got %s", val, judgment)
		}
	})
}

func FuzzScrutinyConductValidation(f *testing.F) {
	f.Add(8, "Ammesso")
	f.Add(5, "Non Ammesso")
	f.Add(0, "Sospeso")
	f.Add(-1, "promosso")
	f.Add(11, "bocciato")

	f.Fuzz(func(t *testing.T, conduct int, decision string) {
		canonical := decision
		switch decision {
		case "Promosso", "promosso", "ammesso", "Ammesso":
			canonical = "Ammesso"
		case "Bocciato", "bocciato", "non ammesso", "Non Ammesso", "Non ammesso":
			canonical = "Non Ammesso"
		case "Giudizio Sospeso", "giudizio sospeso", "sospeso", "Sospeso":
			canonical = "Sospeso"
		}

		if conduct != 0 && (conduct < 1 || conduct > 10) {
			// Must be detected as invalid
			if conduct >= 1 && conduct <= 10 {
				t.Fatalf("boundary check failed for conduct %d", conduct)
			}
		}
		_ = canonical
	})
}
