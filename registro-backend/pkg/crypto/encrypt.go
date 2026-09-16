package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

var (
	ErrCiphertextTooShort = errors.New("ciphertext too short")
	ErrDecryptionFailed   = errors.New("decryption failed")
)

// GetEncryptionKey returns the 32-byte key from ENCRYPTION_KEY env var
// or a secure default key for development purposes.
func GetEncryptionKey() []byte {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr != "" {
		key, err := hex.DecodeString(keyStr)
		if err == nil && len(key) == 32 {
			return key
		}
		if len(keyStr) == 32 {
			return []byte(keyStr)
		}
	}
	// Fallback 32-byte development key (hex of "registro_elettronico_secure_aes_key")
	// "726567697374726f5f656c657474726f6e69636f5f7365637572655f6b657931"
	fallbackKey, _ := hex.DecodeString("726567697374726f5f656c657474726f6e69636f5f7365637572655f6b657931")
	return fallbackKey
}

// Encrypt encrypts a byte slice using AES-256-GCM.
// Prepend the 12-byte nonce to the ciphertext.
func Encrypt(plaintext []byte) ([]byte, error) {
	key := GetEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt and append to nonce
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts a byte slice using AES-256-GCM.
func Decrypt(ciphertext []byte) ([]byte, error) {
	key := GetEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns a hex-encoded string.
func EncryptString(plaintext string) (string, error) {
	cipherBytes, err := Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(cipherBytes), nil
}

// DecryptString decrypts a hex-encoded string and returns the plaintext string.
func DecryptString(ciphertextHex string) (string, error) {
	if ciphertextHex == "" {
		return "", nil
	}
	cipherBytes, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	plainBytes, err := Decrypt(cipherBytes)
	if err != nil {
		return "", err
	}

	return string(plainBytes), nil
}

// PrehashPassword computes a 32-byte HMAC-SHA256 of the password using the PEPPER_SECRET env var.
// This resolves the 72-byte truncation limit of bcrypt and protects against rainbow tables with a server pepper.
func PrehashPassword(password string) []byte {
	pepper := os.Getenv("PEPPER_SECRET")
	if pepper == "" {
		pepper = "default_pepper_secret_registrov2"
	}
	h := hmac.New(sha256.New, []byte(pepper))
	h.Write([]byte(password))
	return h.Sum(nil)
}
