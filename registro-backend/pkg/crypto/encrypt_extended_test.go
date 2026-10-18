package crypto

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetEncryptionKey_Default verifies the fallback key is returned when no env var is set.
func TestGetEncryptionKey_Default(t *testing.T) {
	t.Setenv("ENCRYPTION_KEY", "")
	key := GetEncryptionKey()
	require.Len(t, key, 32, "default key must be exactly 32 bytes")
}

// TestGetEncryptionKey_HexEnv verifies that a valid 64-hex-char env key is decoded correctly.
func TestGetEncryptionKey_HexEnv(t *testing.T) {
	hexKey := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	t.Setenv("ENCRYPTION_KEY", hexKey)
	key := GetEncryptionKey()
	require.Len(t, key, 32)
	assert.Equal(t, byte(0x01), key[0])
	assert.Equal(t, byte(0x20), key[31])
}

// TestGetEncryptionKey_RawStringEnv verifies that a 32-byte raw string env key is used directly.
func TestGetEncryptionKey_RawStringEnv(t *testing.T) {
	rawKey := "12345678901234567890123456789012" // exactly 32 chars
	t.Setenv("ENCRYPTION_KEY", rawKey)
	key := GetEncryptionKey()
	require.Len(t, key, 32)
	assert.Equal(t, []byte(rawKey), key)
}

// TestGetEncryptionKey_InvalidHexFallsBackToRaw verifies invalid hex but correct-length
// string is treated as raw bytes.
func TestGetEncryptionKey_InvalidHexFallsBackToRaw(t *testing.T) {
	rawKey := "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"
	t.Setenv("ENCRYPTION_KEY", rawKey)
	key := GetEncryptionKey()
	require.Len(t, key, 32)
	assert.Equal(t, []byte(rawKey), key)
}

// TestEncryptDecrypt_RoundTrip verifies that encrypting and decrypting produces the original
// plaintext for various input sizes.
func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		input []byte
	}{
		{"empty", []byte{}},
		{"single null byte", []byte{0x00}},
		{"short string", []byte("hello world")},
		{"long string", []byte(strings.Repeat("a", 1000))},
		{"64KB", make([]byte, 65536)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ct, err := Encrypt(tc.input)
			require.NoError(t, err)
			pt, err := Decrypt(ct)
			require.NoError(t, err)
			if len(tc.input) == 0 {
				assert.Empty(t, pt)
			} else {
				assert.Equal(t, tc.input, pt)
			}
		})
	}
}

// TestDecrypt_TooShort verifies that a ciphertext shorter than the nonce size is rejected.
func TestDecrypt_TooShort(t *testing.T) {
	_, err := Decrypt([]byte{0x01, 0x02, 0x03}) // 3 bytes, nonce is 12 bytes
	require.ErrorIs(t, err, ErrCiphertextTooShort)
}

// TestDecrypt_Tampered verifies that altering a byte in the ciphertext causes decryption failure.
func TestDecrypt_Tampered(t *testing.T) {
	ct, err := Encrypt([]byte("important secret"))
	require.NoError(t, err)
	ct[len(ct)-1] ^= 0xFF // flip the last byte
	_, err = Decrypt(ct)
	require.ErrorIs(t, err, ErrDecryptionFailed)
}

// TestDecryptString_Empty verifies that an empty hex string is decoded to an empty string.
func TestDecryptString_Empty(t *testing.T) {
	result, err := DecryptString("")
	require.NoError(t, err)
	assert.Equal(t, "", result)
}

// TestDecryptString_InvalidHex verifies that a non-hex input returns an error.
func TestDecryptString_InvalidHex(t *testing.T) {
	_, err := DecryptString("not-a-hex-string!!")
	require.Error(t, err)
}

// TestEncryptString_Uniqueness verifies that two encryptions of the same plaintext
// produce different ciphertexts (random nonce).
func TestEncryptString_Uniqueness(t *testing.T) {
	plain := "deterministic nonce check"
	ct1, err := EncryptString(plain)
	require.NoError(t, err)
	ct2, err := EncryptString(plain)
	require.NoError(t, err)
	assert.NotEqual(t, ct1, ct2, "same plaintext must produce different ciphertexts (random nonce)")
}

// TestPrehashPassword_Consistency verifies that the same password + pepper always produces the
// same hash (deterministic HMAC-SHA256).
func TestPrehashPassword_Consistency(t *testing.T) {
	t.Setenv("PEPPER_SECRET", "test-pepper-value")
	h1 := PrehashPassword("MyP@ssw0rd")
	h2 := PrehashPassword("MyP@ssw0rd")
	assert.Equal(t, h1, h2, "same password must produce identical pre-hash")
	assert.Len(t, h1, 32, "HMAC-SHA256 must be 32 bytes")
}

// TestPrehashPassword_DifferentPasswords verifies different passwords produce different hashes.
func TestPrehashPassword_DifferentPasswords(t *testing.T) {
	t.Setenv("PEPPER_SECRET", "test-pepper-value")
	h1 := PrehashPassword("Password1!")
	h2 := PrehashPassword("Password2!")
	assert.NotEqual(t, h1, h2)
}

// TestPrehashPassword_DefaultPepper verifies that the default pepper is used when env is unset.
func TestPrehashPassword_DefaultPepper(t *testing.T) {
	os.Unsetenv("PEPPER_SECRET")
	h := PrehashPassword("test-password")
	assert.Len(t, h, 32)
}

// TestPrehashPassword_PepperChangesHash verifies that changing the pepper changes the hash.
func TestPrehashPassword_PepperChangesHash(t *testing.T) {
	t.Setenv("PEPPER_SECRET", "pepper-A")
	hA := PrehashPassword("same-password")
	t.Setenv("PEPPER_SECRET", "pepper-B")
	hB := PrehashPassword("same-password")
	assert.NotEqual(t, hA, hB, "different pepper must produce different hashes")
}

// TestEncryptDecrypt_KnownKey tests with an explicitly set 32-byte hex key.
func TestEncryptDecrypt_KnownKey(t *testing.T) {
	hexKey := "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899"
	t.Setenv("ENCRYPTION_KEY", hexKey)

	plaintext := "sensitive school data"
	ct, err := EncryptString(plaintext)
	require.NoError(t, err)
	require.NotEmpty(t, ct)
	assert.NotEqual(t, plaintext, ct)

	decrypted, err := DecryptString(ct)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}
