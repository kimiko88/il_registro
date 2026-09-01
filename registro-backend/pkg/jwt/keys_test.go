package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

func TestGetOrGenerateKeys_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	privPath := filepath.Join(tmpDir, "private.pem")
	pubPath := filepath.Join(tmpDir, "public.pem")

	// First call should generate and save keys
	priv1, pub1, err := GetOrGenerateKeys(privPath, pubPath)
	if err != nil {
		t.Fatalf("unexpected error generating keys: %v", err)
	}

	if _, err := os.Stat(privPath); os.IsNotExist(err) {
		t.Errorf("expected private key file to exist at %s", privPath)
	}

	// Second call should load existing keys
	priv2, pub2, err := GetOrGenerateKeys(privPath, pubPath)
	if err != nil {
		t.Fatalf("unexpected error loading keys: %v", err)
	}

	if priv1.D.Cmp(priv2.D) != 0 {
		t.Errorf("expected loaded private key to match generated private key")
	}

	if pub1.N.Cmp(pub2.N) != 0 {
		t.Errorf("expected loaded public key to match generated public key")
	}
}

func TestGetOrGenerateKeys_FromEnvVar(t *testing.T) {
	// Generate a valid private key PEM
	privKey, _, err := GenerateRSAKeys()
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}))

	t.Setenv("RSA_PRIVATE_KEY", privPEM)

	privLoaded, pubLoaded, err := GetOrGenerateKeys("non_existent_priv.pem", "non_existent_pub.pem")
	if err != nil {
		t.Fatalf("failed to load keys from env var: %v", err)
	}

	if privLoaded.D.Cmp(privKey.D) != 0 {
		t.Errorf("expected private key loaded from env to match")
	}

	if pubLoaded.N.Cmp(privKey.N) != 0 {
		t.Errorf("expected public key derived from env private key to match")
	}
}

func TestGetOrGenerateKeys_ReadOnlyDirectoryFallback(t *testing.T) {
	// Provide a path in an invalid/non-existent directory to simulate permission error on save
	privPath := filepath.Join("/non_existent_folder_xyz_123", "private.pem")
	pubPath := filepath.Join("/non_existent_folder_xyz_123", "public.pem")

	priv, pub, err := GetOrGenerateKeys(privPath, pubPath)
	if err != nil {
		t.Fatalf("expected GetOrGenerateKeys to return in-memory keys without error when file write fails, got: %v", err)
	}

	if priv == nil || pub == nil {
		t.Fatalf("expected non-nil in-memory keys")
	}
}
