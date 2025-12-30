package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/url"

	"github.com/pquerna/otp/totp"
)

// MFAService handles MFA operations
type MFAService struct {
	issuer string
}

// NewMFAService creates a new MFA service
func NewMFAService(issuer string) *MFAService {
	return &MFAService{
		issuer: issuer,
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
		secret,
		url.QueryEscape(m.issuer),
	)
}

// VerifyTOTP verifies a TOTP code
func (m *MFAService) VerifyTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
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
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}
