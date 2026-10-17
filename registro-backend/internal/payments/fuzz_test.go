package payments

import (
	"strings"
	"testing"
)

// FuzzCalculateMod97 fuzz-tests the ISO 7064 Mod 97-10 check digits algorithm against unexpected numeric inputs.
func FuzzCalculateMod97(f *testing.F) {
	// Seed inputs
	f.Add("001000000000001")
	f.Add("000000000000000")
	f.Add("999999999999999")
	f.Add("123456789012345")
	f.Add("021984719283748")

	f.Fuzz(func(t *testing.T, input string) {
		check, err := CalculateMod97(input)
		// If input is strictly numeric and non-empty, it must succeed and return a 2-digit check
		isNumeric := len(input) > 0
		for _, ch := range input {
			if ch < '0' || ch > '9' {
				isNumeric = false
				break
			}
		}

		if isNumeric {
			if err != nil {
				t.Errorf("expected success for numeric input %q, got: %v", input, err)
			}
			if len(check) != 2 {
				t.Errorf("expected 2-digit check, got %q", check)
			}
		} else {
			if err == nil {
				t.Errorf("expected error for non-numeric/empty input %q", input)
			}
		}
	})
}

// FuzzParsePagoPAQRCodePayload tests robust parsing of untrusted QR code payload strings.
func FuzzParsePagoPAQRCodePayload(f *testing.F) {
	f.Add("PAGOPA|002|00100000000000142|97123456789|4500")
	f.Add("MALFORMED|DATA")
	f.Add("PAGOPA|002|SHORT|123|notanumber")
	f.Add("||||")

	f.Fuzz(func(t *testing.T, payload string) {
		// Ensure parser never panics on arbitrary string payloads
		_, _, _, _ = ParsePagoPAQRCodePayload(payload)
	})
}

// FuzzOPIStreamParser tests robust handling of arbitrary XML/CSV quietanza files.
func FuzzOPIStreamParser(f *testing.F) {
	f.Add([]byte("<flusso_giornale_di_cassa><identificativo_flusso>FL1</identificativo_flusso></flusso_giornale_di_cassa>"), "OPI_XML")
	f.Add([]byte("IUV,Importo,Data,Quietanza\n00100000000000142,45.00,2026-10-05,Q1"), "CSV")
	f.Add([]byte("random junk content that is neither xml nor csv"), "XML")

	f.Fuzz(func(t *testing.T, data []byte, format string) {
		// Must never crash or panic
		_, _ = ParseOPIStream(strings.NewReader(string(data)), format)
	})
}
