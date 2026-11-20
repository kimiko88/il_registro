package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
)

// MFAService handles MFA operations
type MFAService struct {
	issuer      string
	usedOTPs    map[string]time.Time
	redisClient *redis.Client
	replayErr   error
	mu          sync.RWMutex
}

// NewMFAService creates a new MFA service
func NewMFAService(issuer string, redisURLs ...string) *MFAService {
	s := &MFAService{
		issuer:   issuer,
		usedOTPs: make(map[string]time.Time),
	}
	if len(redisURLs) > 0 && strings.TrimSpace(redisURLs[0]) != "" {
		opt, err := redis.ParseURL(redisURLs[0])
		if err != nil {
			s.replayErr = err
			return s
		}
		client := redis.NewClient(opt)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Ping(ctx).Err()
		cancel()
		if err != nil {
			_ = client.Close()
			s.replayErr = err
			return s
		}
		s.redisClient = client
		return s
	}
	// Process-local replay protection is permitted only when no shared Redis
	// backend was explicitly configured (development/test mode).
	go s.cleanupLoop()
	return s
}

// Close releases the shared Redis client used for TOTP replay protection.
func (m *MFAService) Close() error {
	if m.redisClient != nil {
		return m.redisClient.Close()
	}
	return nil
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

// VerifyTOTPUser verifies a TOTP code and atomically prevents replay for 90s per user.
func (m *MFAService) VerifyTOTPUser(userID, secret, code string) bool {
	if !totp.Validate(code, secret) || m.replayErr != nil {
		return false
	}
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%s", userID, code)))
	key := "mfa:otp:" + hex.EncodeToString(sum[:])
	if m.redisClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		accepted, err := m.redisClient.SetNX(ctx, key, "1", 90*time.Second).Result()
		return err == nil && accepted
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if exp, exists := m.usedOTPs[key]; exists && time.Now().Before(exp) {
		return false
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
