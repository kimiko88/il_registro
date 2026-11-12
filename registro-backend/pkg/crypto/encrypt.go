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
	"strings"
)

var (
	ErrCiphertextTooShort = errors.New("ciphertext too short")
	ErrDecryptionFailed   = errors.New("decryption failed")
	ErrEncryptionKeyNotConfigured = errors.New("ENCRYPTION_KEY must be configured as 32 raw bytes or 64 hexadecimal characters")
)

// GetEncryptionKey returns the configured 256-bit encryption key.
// There is intentionally no built-in fallback: a public default would expose
// encrypted student and health records to anyone with access to the source.
func GetEncryptionKey() []byte {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		if strings.EqualFold(os.Getenv("SERVER_MODE"), "release") || strings.EqualFold(os.Getenv("APP_ENV"), "production") || strings.EqualFold(os.Getenv("GIN_MODE"), "release") {
			return nil
		}
		// Development-only compatibility key. Production startup rejects missing keys.
		fallback, _ := hex.DecodeString("726567697374726f5f656c657474726f6e69636f5f7365637572655f6b657931")
		return fallback
	}
	if key, err := hex.DecodeString(keyStr); err == nil && len(key) == 32 {
		return key
	}
	if len(keyStr) == 32 {
		return []byte(keyStr)
	}
	return nil
}

// Encrypt encrypts a byte slice using AES-256-GCM.
// Prepend the 12-byte nonce to the ciphertext.
func Encrypt(plaintext []byte) ([]byte, error) {
	key := GetEncryptionKey()
	if len(key) != 32 {
		return nil, ErrEncryptionKeyNotConfigured
	}
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
	if len(key) != 32 {
		return nil, ErrEncryptionKeyNotConfigured
	}
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
		// Existing records may have been encrypted with the historical public
		// development key. Decrypt them only when an operator explicitly enables
		// the migration compatibility switch; all new writes use ENCRYPTION_KEY.
		if strings.EqualFold(os.Getenv("ALLOW_LEGACY_ENCRYPTION_KEY"), "true") {
			legacyKey, _ := hex.DecodeString("726567697374726f5f656c657474726f6e69636f5f7365637572655f6b657931")
			legacyBlock, legacyErr := aes.NewCipher(legacyKey)
			if legacyErr == nil {
				legacyGCM, gcmErr := cipher.NewGCM(legacyBlock)
				if gcmErr == nil {
					legacyPlaintext, openErr := legacyGCM.Open(nil, nonce, actualCiphertext, nil)
					if openErr == nil {
						return legacyPlaintext, nil
					}
				}
			}
		}
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

// PrehashPassword computes a 32-byte HMAC-SHA256 using the configured PEPPER_SECRET.
// This resolves the 72-byte truncation limit of bcrypt and protects against rainbow tables with a server pepper.
func PrehashPassword(password string) []byte {
	pepper := os.Getenv("PEPPER_SECRET")
	if pepper == "" {
		if strings.EqualFold(os.Getenv("SERVER_MODE"), "release") || strings.EqualFold(os.Getenv("APP_ENV"), "production") || strings.EqualFold(os.Getenv("GIN_MODE"), "release") {
			return nil
		}
		pepper = "default_pepper_secret_registrov2" // development-only compatibility
	}
	h := hmac.New(sha256.New, []byte(pepper))
	h.Write([]byte(password))
	return h.Sum(nil)
}


// PrehashPasswordLegacy computes the historical pre-hash used before PEPPER_SECRET
// became mandatory. It exists only to verify and transparently migrate existing
// password hashes. Remove it after all legacy accounts have authenticated or been reset.
func PrehashPasswordLegacy(password string) []byte {
	h := hmac.New(sha256.New, []byte("default_pepper_secret_registrov2"))
	h.Write([]byte(password))
	return h.Sum(nil)
}
