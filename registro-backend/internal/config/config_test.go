package config

import (
	"os"
	"testing"
)

func TestLoadConfig_DefaultsAndEnv(t *testing.T) {
	// Set custom environment variables for testing
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_MODE", "debug")
	os.Setenv("DB_HOST", "https://db.example.com")
	os.Setenv("JWT_SECRET", "super-secret-key-123")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_MODE")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("JWT_SECRET")
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
	os.Setenv("DB_HOST", "http://localhost")
	defer os.Unsetenv("DB_HOST")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected Database.Host 'localhost', got '%s'", cfg.Database.Host)
	}
}
