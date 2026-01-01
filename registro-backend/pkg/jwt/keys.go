package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
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

// LoadRSAKeys loads RSA keys from files
func LoadRSAKeys(privatePath, publicPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Read Private Key
	privatePEM, err := os.ReadFile(privatePath)
	if err != nil {
		return nil, nil, err
	}
	block, _ := pem.Decode(privatePEM)
	if block == nil {
		return nil, nil, errors.New("failed to parse private key PEM")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	// Read Public Key
	publicPEM, err := os.ReadFile(publicPath)
	if err != nil {
		return nil, nil, err
	}
	block, _ = pem.Decode(publicPEM)
	if block == nil {
		return nil, nil, errors.New("failed to parse public key PEM")
	}
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

// GetOrGenerateKeys loads keys from files or generates new ones if missing
func GetOrGenerateKeys(privatePath, publicPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if _, err := os.Stat(privatePath); os.IsNotExist(err) {
		priv, pub, err := GenerateRSAKeys()
		if err != nil {
			return nil, nil, err
		}
		if err := SaveRSAKeys(priv, privatePath, publicPath); err != nil {
			return nil, nil, err
		}
		return priv, pub, nil
	}
	return LoadRSAKeys(privatePath, publicPath)
}
