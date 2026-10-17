package pdp

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"registro-backend/pkg/crypto"
)

func TestPDP_DiagnosisFieldLevelEncryption(t *testing.T) {
	rawDiagnosis := "Disturbo Specifico dell'Apprendimento (F81.0 Dislessia Evolutiva L.170/2010 - Note: àèìòù € %)"

	// 1. Encrypt diagnosis
	ciphertext, err := crypto.EncryptString(rawDiagnosis)
	assert.NoError(t, err)
	assert.NotEmpty(t, ciphertext)
	assert.NotEqual(t, rawDiagnosis, ciphertext, "Ciphertext must not match raw plaintext")

	// 2. Decrypt diagnosis
	decrypted, err := crypto.DecryptString(ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, rawDiagnosis, decrypted, "Decrypted text must match original plaintext with accents/symbols")

	// 3. Fallback for unencrypted legacy data
	legacyCleartext := "Diagnosi pregressa in chiaro non cifrata"
	failedDecrypted, err := crypto.DecryptString(legacyCleartext)
	assert.Error(t, err, "Attempting to decrypt raw plaintext should return an error")
	assert.Empty(t, failedDecrypted)

	// Scan fallback behavior simulation
	diagToScan := legacyCleartext
	if dec, err := crypto.DecryptString(diagToScan); err == nil {
		diagToScan = dec
	}
	assert.Equal(t, legacyCleartext, diagToScan, "ScanPlan fallback preserves legacy cleartext without data loss")

	// 4. Tampered ciphertext detection
	cipherBytes, err := hex.DecodeString(ciphertext)
	assert.NoError(t, err)
	if len(cipherBytes) > 5 {
		// Flip a bit in the ciphertext payload
		cipherBytes[len(cipherBytes)-1] ^= 0xFF
		tamperedHex := hex.EncodeToString(cipherBytes)
		_, errTampered := crypto.DecryptString(tamperedHex)
		assert.Error(t, errTampered, "Tampered ciphertext must fail GCM authentication check")
	}

	// 5. Empty string encryption handling
	emptyCipher, err := crypto.EncryptString("")
	assert.NoError(t, err)
	emptyDecrypted, err := crypto.DecryptString(emptyCipher)
	assert.NoError(t, err)
	assert.Empty(t, emptyDecrypted)
}
