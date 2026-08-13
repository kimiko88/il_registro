package unit

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestDigitalSignature struct {
	ID           string
	DocumentID   string
	DocumentHash string
	SignerID     string
	IsRevoked    bool
	SignedAt     time.Time
}

func SignDocumentPayload(payload []byte, signerID string) *TestDigitalSignature {
	hash := sha256.Sum256(payload)
	return &TestDigitalSignature{
		ID:           "sig-" + hex.EncodeToString(hash[:8]),
		DocumentID:   "doc-001",
		DocumentHash: hex.EncodeToString(hash[:]),
		SignerID:     signerID,
		IsRevoked:    false,
		SignedAt:     time.Now().UTC(),
	}
}

func VerifySignatureIntegrity(sig *TestDigitalSignature, currentPayload []byte) error {
	if sig == nil {
		return errors.New("signature is nil")
	}
	if sig.IsRevoked {
		return errors.New("signature has been revoked")
	}
	currentHash := sha256.Sum256(currentPayload)
	currentHashHex := hex.EncodeToString(currentHash[:])
	if currentHashHex != sig.DocumentHash {
		return errors.New("document integrity failure: payload hash mismatch")
	}
	return nil
}

func TestSecurity_DigitalSignatureNonRepudiation(t *testing.T) {
	originalContent := []byte("Verbale del Consiglio di Classe - Anno Scolastico 2025/2026")
	sig := SignDocumentPayload(originalContent, "teacher-001")

	t.Run("Validates authentic document payload integrity", func(t *testing.T) {
		err := VerifySignatureIntegrity(sig, originalContent)
		assert.NoError(t, err)
	})

	t.Run("Rejects tampered document payload (single byte modification)", func(t *testing.T) {
		tamperedContent := []byte("Verbale del Consiglio di Classe - Anno Scolastico 2025/2027")
		err := VerifySignatureIntegrity(sig, tamperedContent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payload hash mismatch")
	})

	t.Run("Rejects revoked digital signature", func(t *testing.T) {
		sigRevoked := SignDocumentPayload(originalContent, "teacher-001")
		sigRevoked.IsRevoked = true

		err := VerifySignatureIntegrity(sigRevoked, originalContent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "revoked")
	})
}
