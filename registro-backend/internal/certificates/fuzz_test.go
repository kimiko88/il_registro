package certificates

import (
	"strings"
	"testing"
)

func FuzzIsValidCertificateType(f *testing.F) {
	seeds := []string{
		"iscrizione",
		"frequenza",
		"promozione",
		"condotta",
		"diploma",
		"certificate",
		"attendance",
		"unknown",
		"",
		"ISCRIZIONE",
		"frequenza ",
		"admin",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		isValid := IsValidCertificateType(CertificateType(input))

		// Known valid types
		trimmed := strings.TrimSpace(input)
		switch CertificateType(trimmed) {
		case CertIscrizione, CertFrequenza, CertPromozione, CertBehavior, "diploma", "certificate", "attendance":
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
