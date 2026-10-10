package family_desk

import (
	"testing"
)

func BenchmarkIsValidRequestType(b *testing.B) {
	types := []string{
		TypeDelegaRitiro,
		TypeUscitaAutonomaUnder14,
		TypeEsoneroMotoria,
		TypeSomministrazioneFarmaci,
		TypeNullaOsta,
		TypeCertificatoIscrizioneFrequenza,
		"invalid_type",
	}

	idx := 0
	for b.Loop() {
		_ = IsValidRequestType(types[idx%len(types)])
		idx++
	}
}

func BenchmarkBuildDelegateFromRequest(b *testing.B) {
	req := &FamilyRequest{
		ID:          "req-bench",
		SchoolID:    "school-1",
		StudentID:   "student-1",
		RequestType: TypeDelegaRitiro,
		FormData: map[string]interface{}{
			"first_name":      "Mario",
			"last_name":       "Rossi",
			"tax_code":        "RSSMRA80A01H501U",
			"relationship":    "Zio",
			"phone":           "3331234567",
			"id_card_details": "Carta Identita CA987654",
		},
	}

	for b.Loop() {
		_, _ = BuildDelegateFromRequest(req)
	}
}

func BenchmarkValidateRequestFormData(b *testing.B) {
	req := &FamilyRequest{
		RequestType: TypeUscitaAutonomaUnder14,
		FormData: map[string]interface{}{
			"consent_given": true,
		},
	}

	for b.Loop() {
		_ = ValidateRequestFormData(req)
	}
}
