package config

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	appcrypto "registro-backend/pkg/crypto"
)

func TestLoadConfig_DefaultsAndEnv(t *testing.T) {
	// Set custom environment variables for testing
	_ = os.Setenv("SERVER_PORT", "9090")
	_ = os.Setenv("SERVER_MODE", "debug")
	_ = os.Setenv("DB_HOST", "https://db.example.com")
	_ = os.Setenv("JWT_SECRET", "super-secret-key-123")
	defer func() {
		_ = os.Unsetenv("SERVER_PORT")
		_ = os.Unsetenv("SERVER_MODE")
		_ = os.Unsetenv("DB_HOST")
		_ = os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected LoadConfig to succeed, got error: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("expected Server.Port to be '9090', got '%s'", cfg.Server.Port)
	}

	if cfg.Server.Mode != "debug" {
		t.Errorf("expected Server.Mode to be 'debug', got '%s'", cfg.Server.Mode)
	}

	// Verify URL prefix trimming for DB Host
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected Database.Host to be 'db.example.com' (trimmed prefix), got '%s'", cfg.Database.Host)
	}

	if cfg.JWT.Secret != "super-secret-key-123" {
		t.Errorf("expected JWT.Secret to be 'super-secret-key-123', got '%s'", cfg.JWT.Secret)
	}

	if cfg.Supabase.Bucket != "documents" {
		t.Errorf("expected default Supabase bucket 'documents', got '%s'", cfg.Supabase.Bucket)
	}
}

func TestLoadConfig_TrimsHttpPrefix(t *testing.T) {
	_ = os.Setenv("SERVER_MODE", "debug")
	defer func() { _ = os.Unsetenv("SERVER_MODE") }()
	_ = os.Setenv("DB_HOST", "http://localhost")
	defer func() { _ = os.Unsetenv("DB_HOST") }()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected Database.Host 'localhost', got '%s'", cfg.Database.Host)
	}
}

func setProductionConfigEnv(t *testing.T) {
	t.Helper()
	t.Setenv("SERVER_MODE", "release")
	t.Setenv("DB_PASSWORD", "db-secret-not-for-production-0123456789")
	t.Setenv("ALLOWED_ORIGINS", "https://registro.example.it")
	t.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	t.Setenv("PEPPER_SECRET", "pepper-secret-not-for-production-0123456789")
	t.Setenv("JWT_SECRET", "jwt-secret-not-for-production-0123456789-abcdefghijklmnopqrstuvwxyz-0123456789")
}

func TestLoadConfig_ProductionRequiresRedisAuthentication(t *testing.T) {
	setProductionConfigEnv(t)
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_URL", "")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected production config without Redis authentication to fail")
	}
	if !strings.Contains(err.Error(), "redis authentication is required") {
		t.Fatalf("expected Redis authentication error, got: %v", err)
	}
}

func TestLoadConfig_UsesPasswordEmbeddedInRedisURL(t *testing.T) {
	setProductionConfigEnv(t)
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_URL", "redis://:0123456789abcdef0123456789abcdef@localhost:6379/0")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected production config with authenticated Redis URL to succeed: %v", err)
	}
	if cfg.Redis.Password != "0123456789abcdef0123456789abcdef" {
		t.Fatalf("expected Redis password to be extracted from REDIS_URL")
	}
	if !strings.Contains(cfg.RedisURL(), "0123456789abcdef0123456789abcdef") {
		t.Fatalf("expected RedisURL to retain Redis authentication")
	}
}

func TestLoadConfig_AppliesDotEnvSecretsToCryptoRuntime(t *testing.T) {
	setProductionConfigEnv(t)
	pepper := "runtime-pepper-secret-0123456789abcdef"
	keyHex := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	t.Setenv("PEPPER_SECRET", pepper)
	t.Setenv("ENCRYPTION_KEY", keyHex)
	t.Setenv("REDIS_PASSWORD", "redis-secret-0123456789abcdef0123456789abcdef")
	t.Setenv("ALLOW_LEGACY_PASSWORD_PEPPER", "false")

	if _, err := LoadConfig(); err != nil {
		t.Fatalf("expected production config to load: %v", err)
	}
	expectedKey, _ := hex.DecodeString(keyHex)
	if !bytes.Equal(appcrypto.GetEncryptionKey(), expectedKey) {
		t.Fatal("crypto runtime did not receive the validated ENCRYPTION_KEY")
	}
	mac := hmac.New(sha256.New, []byte(pepper))
	_, _ = mac.Write([]byte("runtime-check"))
	if !bytes.Equal(appcrypto.PrehashPassword("runtime-check"), mac.Sum(nil)) {
		t.Fatal("password hashing did not receive the validated PEPPER_SECRET")
	}
	if appcrypto.LegacyPasswordMigrationEnabled() {
		t.Fatal("legacy password migration should remain disabled")
	}
}
