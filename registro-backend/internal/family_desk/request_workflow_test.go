package family_desk

import (
	"testing"
)

func TestValidRequestTypes(t *testing.T) {
	validTypes := []string{
		TypeDelegaRitiro,
		TypeUscitaAutonomaUnder14,
		TypeEsoneroMotoria,
		TypeSomministrazioneFarmaci,
		TypeNullaOsta,
		TypeCertificatoIscrizioneFrequenza,
	}

	for _, vt := range validTypes {
		if !IsValidRequestType(vt) {
			t.Errorf("expected request type %s to be valid", vt)
		}
	}

	if IsValidRequestType("richiesta_falsa") {
		t.Errorf("expected fake request type to be invalid")
	}
}

func TestUnder14Compliance(t *testing.T) {
	// L. 172/2017 requires parental disclaimer and consent for under-14 autonomous exit
	req := FamilyRequest{
		RequestType: TypeUscitaAutonomaUnder14,
		FormData: map[string]interface{}{
			"consent_given":     true,
			"exempt_school_rep": true,
		},
	}

	if err := ValidateRequestFormData(&req); err != nil {
		t.Errorf("valid under-14 request should pass validation, got err: %v", err)
	}

	invalidReq := FamilyRequest{
		RequestType: TypeUscitaAutonomaUnder14,
		FormData: map[string]interface{}{
			"consent_given": false,
		},
	}
	if err := ValidateRequestFormData(&invalidReq); err == nil {
		t.Errorf("under-14 request without consent must fail validation")
	}
}

func TestBuildDelegateFromDelegaRequest(t *testing.T) {
	req := FamilyRequest{
		ID:          "req-123",
		StudentID:   "student-456",
		SchoolID:    "school-789",
		RequestType: TypeDelegaRitiro,
		FormData: map[string]interface{}{
			"first_name":      "Mario",
			"last_name":       "Rossi",
			"tax_code":        "RSSMRA80A01H501U",
			"relationship":    "nonno",
			"phone":           "+393331234567",
			"id_card_details": "CI CA12345AA",
		},
	}

	delegate, err := BuildDelegateFromRequest(&req)
	if err != nil {
		t.Fatalf("failed to build delegate: %v", err)
	}

	if delegate.FirstName != "Mario" || delegate.TaxCode != "RSSMRA80A01H501U" {
		t.Errorf("unexpected delegate fields: %+v", delegate)
	}
	if !delegate.IsValid {
		t.Errorf("delegate should be valid by default upon approval")
	}
	if delegate.ApprovedRequestID == nil || *delegate.ApprovedRequestID != "req-123" {
		t.Errorf("delegate should link to approved request id")
	}
}
