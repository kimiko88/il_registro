package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// GenerateRSAKeys generates a new RSA key pair
func GenerateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// SaveRSAKeys saves RSA keys to files
func SaveRSAKeys(privateKey *rsa.PrivateKey, privatePath, publicPath string) error {
	// Encode Private Key
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	})
	if err := os.WriteFile(privatePath, privatePEM, 0600); err != nil {
		return err
	}

	// Encode Public Key
	publicBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	})
	if err := os.WriteFile(publicPath, publicPEM, 0644); err != nil {
		return err
	}

	return nil
}

// ParsePrivateKeyPEM parses a Private Key from PEM bytes (supports PKCS#1 and PKCS#8)
func ParsePrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to parse private key PEM block")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
	}
	return nil, errors.New("failed to parse RSA private key (unsupported format)")
}

// ParsePublicKeyPEM parses a Public Key from PEM bytes (supports PKCS#1 and PKIX/PKCS#8)
func ParsePublicKeyPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to parse public key PEM block")
	}
	if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PublicKey); ok {
			return rsaKey, nil
		}
	}
	return nil, errors.New("failed to parse RSA public key (unsupported format)")
}

// LoadRSAKeys loads RSA keys from files
func LoadRSAKeys(privatePath, publicPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privatePEM, err := os.ReadFile(privatePath)
	if err != nil {
		return nil, nil, fmt.Errorf("read private key file: %w", err)
	}
	privateKey, err := ParsePrivateKeyPEM(privatePEM)
	if err != nil {
		return nil, nil, fmt.Errorf("parse private key: %w", err)
	}

	publicPEM, err := os.ReadFile(publicPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read public key file: %w", err)
	}
	publicKey, err := ParsePublicKeyPEM(publicPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("parse public key: %w", err)
	}

	return privateKey, publicKey, nil
}

// GetOrGenerateKeys loads keys from environment variables, disk files, or generates new ones in memory if missing
func GetOrGenerateKeys(privatePath, publicPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// 1. Check environment variables first (RSA_PRIVATE_KEY or JWT_PRIVATE_KEY)
	envPriv := os.Getenv("RSA_PRIVATE_KEY")
	if envPriv == "" {
		envPriv = os.Getenv("JWT_PRIVATE_KEY")
	}

	if envPriv != "" {
		envPriv = strings.ReplaceAll(envPriv, "\\n", "\n")
		privKey, err := ParsePrivateKeyPEM([]byte(envPriv))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse RSA_PRIVATE_KEY / JWT_PRIVATE_KEY from env: %w", err)
		}

		envPub := os.Getenv("RSA_PUBLIC_KEY")
		if envPub == "" {
			envPub = os.Getenv("JWT_PUBLIC_KEY")
		}
		if envPub != "" {
			envPub = strings.ReplaceAll(envPub, "\\n", "\n")
			pubKey, err := ParsePublicKeyPEM([]byte(envPub))
			if err == nil {
				return privKey, pubKey, nil
			}
		}
		// If public key is not explicitly supplied in env, derive public key from private key
		return privKey, &privKey.PublicKey, nil
	}

	// 2. Check if private key exists on disk
	if _, err := os.Stat(privatePath); err == nil {
		return LoadRSAKeys(privatePath, publicPath)
	}

	// 3. Otherwise generate new keys in memory
	priv, pub, err := GenerateRSAKeys()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA keys: %w", err)
	}

	// Try to save to disk if path is writable; if permission is denied, log warning and use in-memory keys
	if err := SaveRSAKeys(priv, privatePath, publicPath); err != nil {
		log.Printf("Warning: could not save RSA keys to disk (%s/%s): %v (continuing with in-memory keys)", privatePath, publicPath, err)
	} else {
		log.Printf("Generated and saved new RSA keys to %s and %s", privatePath, publicPath)
	}

	return priv, pub, nil
}
