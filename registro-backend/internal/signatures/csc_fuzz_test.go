package signatures

import (
	"testing"
)

// FuzzVerifyDigitalStamp tests that VerifyDigitalStamp never panics on arbitrary string inputs.
func FuzzVerifyDigitalStamp(f *testing.F) {
	key := []byte("fuzz-test-key-seal")
	sha := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"

	f.Add("CAD-GLIFO|01|doc1|a591a6d40bf42040|Dirigente|1700000000|0123456789abcdef", sha)
	f.Add("MALFORMED", "EMPTY")
	f.Add("||||||", "")
	f.Add("CAD-GLIFO|01|doc1|short|role|notanumber|sig", "anysha")

	f.Fuzz(func(t *testing.T, payload string, originalSHA string) {
		_, _ = VerifyDigitalStamp(payload, originalSHA, key)
	})
}
