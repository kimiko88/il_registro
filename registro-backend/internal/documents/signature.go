package documents

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"
)

type SignatureProvider struct{}

func NewSignatureProvider() *SignatureProvider {
	return &SignatureProvider{}
}

// Verify parses a PEM-encoded X.509 certificate and verifies a base64-encoded digital signature
// generated over the provided raw document content.
func (s *SignatureProvider) Verify(data, signatureBase64, certPEM string) error {
	if signatureBase64 == "" || certPEM == "" {
		return errors.New("missing signature or certificate")
	}

	// 1. Decode PEM certificate block
	block, _ := pem.Decode([]byte(certPEM))
	var certBytes []byte
	if block == nil {
		// Fallback: try decoding raw base64 DER certificate if not in PEM format
		var err error
		certBytes, err = base64.StdEncoding.DecodeString(certPEM)
		if err != nil {
			return errors.New("invalid certificate encoding: must be PEM or base64 DER")
		}
	} else {
		certBytes = block.Bytes
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return fmt.Errorf("failed to parse X.509 certificate: %w", err)
	}

	// 2. Validate certificate validity period
	now := time.Now()
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		return errors.New("certificate is expired or not yet valid")
	}

	// 3. Decode base64 signature
	sigBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return errors.New("invalid signature encoding: base64 required")
	}

	// 4. Compute SHA-256 digest of input document payload
	digest := sha256.Sum256([]byte(data))

	// 5. Verify signature against certificate public key (RSA or ECDSA)
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sigBytes); err != nil {
			return errors.New("RSA digital signature verification failed")
		}
	case *ecdsa.PublicKey:
		if !ecdsa.VerifyASN1(pub, digest[:], sigBytes) {
			return errors.New("ECDSA digital signature verification failed")
		}
	default:
		return errors.New("unsupported public key algorithm in certificate")
	}

	return nil
}
