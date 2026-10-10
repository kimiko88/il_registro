package family_desk

import (
	"strings"
	"testing"
)

func FuzzIsValidRequestType(f *testing.F) {
	seeds := []string{
		TypeDelegaRitiro,
		TypeUscitaAutonomaUnder14,
		TypeEsoneroMotoria,
		TypeSomministrazioneFarmaci,
		TypeNullaOsta,
		TypeCertificatoIscrizioneFrequenza,
		"unknown",
		"",
		"DELEGA_RITIRO",
		"uscita_autonoma",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		isValid := IsValidRequestType(input)
		trimmed := strings.TrimSpace(input)
		switch trimmed {
		case TypeDelegaRitiro, TypeUscitaAutonomaUnder14, TypeEsoneroMotoria,
			TypeSomministrazioneFarmaci, TypeNullaOsta, TypeCertificatoIscrizioneFrequenza:
			if trimmed == input && !isValid {
				t.Fatalf("expected valid for %s", input)
			}
		default:
			if isValid {
				t.Fatalf("expected invalid for %s", input)
			}
		}
	})
}

func FuzzBuildDelegateFromRequest(f *testing.F) {
	f.Add("Mario", "Rossi", "RSSMRA80A01H501U", "Zio", "3331234567")
	f.Add("", "Rossi", "RSSMRA80A01H501U", "", "")
	f.Add("Luigi", "", "", "Nonno", "")
	f.Add("Anna", "Verdi", "VRDNNA75C41F205Z", "Nonna", "+3906123456")

	f.Fuzz(func(t *testing.T, fn, ln, tc, rel, phone string) {
		req := &FamilyRequest{
			ID:          "req-fuzz",
			SchoolID:    "school-1",
			StudentID:   "student-1",
			RequestType: TypeDelegaRitiro,
			FormData: map[string]interface{}{
				"first_name":      fn,
				"last_name":       ln,
				"tax_code":        tc,
				"relationship":    rel,
				"phone":           phone,
				"id_card_details": "CI123456",
			},
		}

		delegate, err := BuildDelegateFromRequest(req)
		if fn == "" || ln == "" || tc == "" {
			if err == nil {
				t.Fatalf("expected error for incomplete delegate data (fn=%q ln=%q tc=%q)", fn, ln, tc)
			}
			return
		}

		if err != nil {
			t.Fatalf("unexpected error for valid delegate data: %v", err)
		}
		if delegate == nil {
			t.Fatalf("expected delegate not to be nil")
		}
		if delegate.FirstName != fn || delegate.LastName != ln || delegate.TaxCode != tc {
			t.Fatalf("delegate data mismatch")
		}
	})
}
