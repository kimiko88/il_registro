package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionDecryption(t *testing.T) {
	originalText := "Questo è un testo estremamente sensibile che deve essere protetto dal GDPR!"

	// Test string encryption
	ciphertext, err := EncryptString(originalText)
	assert.NoError(t, err)
	assert.NotEmpty(t, ciphertext)
	assert.NotEqual(t, originalText, ciphertext)

	// Test string decryption
	decryptedText, err := DecryptString(ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, originalText, decryptedText)
}

func TestEncryptionDecryptionBytes(t *testing.T) {
	originalBytes := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	ciphertext, err := Encrypt(originalBytes)
	assert.NoError(t, err)
	assert.NotEmpty(t, ciphertext)

	decryptedBytes, err := Decrypt(ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, originalBytes, decryptedBytes)
}
