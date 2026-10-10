package pdp

import (
	"testing"

	"registro-backend/pkg/crypto"
)

func FuzzPDPDiagnosisEncryptionRoundtrip(f *testing.F) {
	seeds := []string{
		"Disturbo Specifico dell'Apprendimento (F81.0 Dislessia Evolutiva)",
		"Diagnosi funzionale PEI ai sensi della L. 104/1992",
		"",
		"Caratteri speciali: àèìòù, €100, <script>alert(1)</script>, %20, \n\t\r",
		"1234567890",
		"✨ Diagnosi con emoji e simboli grafici 🎯",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, plaintext string) {
		ciphertext, err := crypto.EncryptString(plaintext)
		if err != nil {
			t.Fatalf("encryption failed for plaintext %q: %v", plaintext, err)
		}

		decrypted, err := crypto.DecryptString(ciphertext)
		if err != nil {
			t.Fatalf("decryption failed for ciphertext %q: %v", ciphertext, err)
		}

		if decrypted != plaintext {
			t.Fatalf("roundtrip mismatch: got %q, want %q", decrypted, plaintext)
		}
	})
}

func FuzzPDPPlanTypeValidation(f *testing.F) {
	seeds := []string{
		string(PlanTypePDP),
		string(PlanTypePEI),
		"pdp",
		"pei",
		"bes",
		"dsa",
		"",
		"PDP",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		pt := PlanType(input)
		switch pt {
		case PlanTypePDP, PlanTypePEI:
			if pt != "pdp" && pt != "pei" {
				t.Fatalf("unexpected value for valid plan type: %s", pt)
			}
		}
	})
}
