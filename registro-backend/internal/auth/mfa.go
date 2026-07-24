package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/pquerna/otp/totp"
)

// MFAService handles MFA operations
type MFAService struct {
	issuer   string
	usedOTPs map[string]time.Time
	mu       sync.RWMutex
}

// NewMFAService creates a new MFA service
func NewMFAService(issuer string) *MFAService {
	s := &MFAService{
		issuer:   issuer,
		usedOTPs: make(map[string]time.Time),
	}
	go s.cleanupLoop()
	return s
}

func (m *MFAService) cleanupLoop() {
	ticker := time.NewTicker(2 * time.Minute)
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for k, exp := range m.usedOTPs {
			if now.After(exp) {
				delete(m.usedOTPs, k)
			}
		}
		m.mu.Unlock()
	}
}

// GenerateSecret generates a new TOTP secret
func (m *MFAService) GenerateSecret() (string, error) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

// GenerateQRCodeURL generates a QR code URL for Google Authenticator
func (m *MFAService) GenerateQRCodeURL(email, secret string) string {
	return fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		url.QueryEscape(m.issuer),
		url.QueryEscape(email),
		url.QueryEscape(secret),
		url.QueryEscape(m.issuer),
	)
}

// VerifyTOTPUser verifies a TOTP code and prevents replay attacks by caching used codes for 90s per user
func (m *MFAService) VerifyTOTPUser(userID, secret, code string) bool {
	if !totp.Validate(code, secret) {
		return false
	}
	key := fmt.Sprintf("%s:%s", userID, code)
	m.mu.Lock()
	defer m.mu.Unlock()
	if exp, exists := m.usedOTPs[key]; exists && time.Now().Before(exp) {
		return false // OTP code was already used within TTL
	}
	m.usedOTPs[key] = time.Now().Add(90 * time.Second)
	return true
}

// VerifyTOTP verifies a TOTP code (legacy/global)
func (m *MFAService) VerifyTOTP(secret, code string) bool {
	return m.VerifyTOTPUser("global", secret, code)
}

// GenerateRecoveryCodes generates backup recovery codes
func (m *MFAService) GenerateRecoveryCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		code, err := generateRandomCode(8)
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}
	return codes, nil
}

func generateRandomCode(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	b := make([]byte, 1)
	for i := 0; i < length; {
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}
		// 252 is 36 * 7. Reject byte values >= 252 to ensure perfectly uniform distribution.
		if b[0] >= 252 {
			continue
		}
		result[i] = charset[int(b[0])%len(charset)]
		i++
	}
	return string(result), nil
}
