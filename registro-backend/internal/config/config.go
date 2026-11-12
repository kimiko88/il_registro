package config

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Supabase  SupabaseConfig
	SPID      SPIDConfig
	CIE       CIEConfig
	Mail      MailConfig
	Elearning ElearningConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type ElearningConfig struct {
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleRedirectURI     string
	MicrosoftClientID     string
	MicrosoftClientSecret string
	MicrosoftTenantID     string
	MicrosoftRedirectURI  string
}

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
}

type SupabaseConfig struct {
	URL    string
	Key    string
	Bucket string
}

type SPIDConfig struct {
	EntityID       string
	CertPath       string
	KeyPath        string
	IDPMetadataURL string
}

type CIEConfig struct {
	EntityID       string
	CertPath       string
	KeyPath        string
	IDPMetadataURL string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// If .env file exists, read it, otherwise ignore error (for docker/render env vars)
	if err := viper.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Printf("Warning: error reading config file: %v", err)
			}
		}
	}

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_MODE", "release")
	viper.SetDefault("DB_SSLMODE", "require")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("SUPABASE_STORAGE_BUCKET", "documents")

	// Production must never start with public/default cryptographic material.
	isProduction := strings.EqualFold(viper.GetString("SERVER_MODE"), "release") ||
		strings.EqualFold(viper.GetString("APP_ENV"), "production") ||
		strings.EqualFold(viper.GetString("GIN_MODE"), "release")
	if isProduction {
		if strings.TrimSpace(viper.GetString("ALLOWED_ORIGINS")) == "" {
			return nil, fmt.Errorf("ALLOWED_ORIGINS must be configured in production")
		}
		key := strings.TrimSpace(viper.GetString("ENCRYPTION_KEY"))
		validEncryptionKey := false
		if decoded, err := hex.DecodeString(key); err == nil && len(decoded) == 32 {
			validEncryptionKey = true
		} else if len(key) == 32 {
			validEncryptionKey = true
		}
		if !validEncryptionKey {
			return nil, fmt.Errorf("ENCRYPTION_KEY is required in release mode (32 raw bytes or 64 hex characters)")
		}
		pepper := strings.TrimSpace(viper.GetString("PEPPER_SECRET"))
		if len(pepper) < 32 || strings.HasPrefix(strings.ToUpper(pepper), "INSERISCI") || pepper == "default_pepper_secret_registrov2" {
			return nil, fmt.Errorf("PEPPER_SECRET must be a generated secret of at least 32 characters in production")
		}
		jwtSecret := strings.TrimSpace(viper.GetString("JWT_SECRET"))
		if len(jwtSecret) < 64 || strings.HasPrefix(strings.ToUpper(jwtSecret), "INSERISCI") {
			return nil, fmt.Errorf("JWT_SECRET must be a generated secret of at least 64 characters in production")
		}
	}

	// Crypto helpers read process environment directly. Promote values loaded
	// from .env only after validating them, so .env deployments work as intended.
	for _, name := range []string{"ENCRYPTION_KEY", "PEPPER_SECRET", "ALLOW_LEGACY_ENCRYPTION_KEY"} {
		if value := viper.GetString(name); value != "" && os.Getenv(name) == "" {
			if err := os.Setenv(name, value); err != nil {
				return nil, fmt.Errorf("failed to initialize %s: %w", name, err)
			}
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
			Mode: viper.GetString("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			Host:     strings.TrimPrefix(strings.TrimPrefix(viper.GetString("DB_HOST"), "https://"), "http://"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},

		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
		Supabase: SupabaseConfig{
			URL:    viper.GetString("SUPABASE_URL"),
			Key:    viper.GetString("SUPABASE_KEY"),
			Bucket: viper.GetString("SUPABASE_STORAGE_BUCKET"),
		},
		SPID: SPIDConfig{
			EntityID:       viper.GetString("SPID_ENTITY_ID"),
			CertPath:       viper.GetString("SPID_CERT_PATH"),
			KeyPath:        viper.GetString("SPID_KEY_PATH"),
			IDPMetadataURL: viper.GetString("SPID_IDP_METADATA_URL"),
		},
		CIE: CIEConfig{
			EntityID:       viper.GetString("CIE_ENTITY_ID"),
			CertPath:       viper.GetString("CIE_CERT_PATH"),
			KeyPath:        viper.GetString("CIE_KEY_PATH"),
			IDPMetadataURL: viper.GetString("CIE_IDP_METADATA_URL"),
		},
		Mail: MailConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			Username: viper.GetString("SMTP_USERNAME"),
			Password: viper.GetString("SMTP_PASSWORD"),
			From:     viper.GetString("SMTP_FROM"),
		},
		Elearning: ElearningConfig{
			GoogleClientID:        viper.GetString("GOOGLE_CLIENT_ID"),
			GoogleClientSecret:    viper.GetString("GOOGLE_CLIENT_SECRET"),
			GoogleRedirectURI:     viper.GetString("GOOGLE_REDIRECT_URI"),
			MicrosoftClientID:     viper.GetString("MICROSOFT_CLIENT_ID"),
			MicrosoftClientSecret: viper.GetString("MICROSOFT_CLIENT_SECRET"),
			MicrosoftTenantID:     viper.GetString("MICROSOFT_TENANT_ID"),
			MicrosoftRedirectURI:  viper.GetString("MICROSOFT_REDIRECT_URI"),
		},
	}

	return config, nil
}

// RedisURL resolves the Redis connection string, preferring REDIS_URL if set,
// or constructing redis://host:port from the host and port config.
func (c *Config) RedisURL() string {
	if redisURL := viper.GetString("REDIS_URL"); redisURL != "" {
		return redisURL
	}
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		return redisURL
	}
	if c.Redis.Host != "" {
		port := c.Redis.Port
		if port == "" {
			port = "6379"
		}
		if c.Redis.Password != "" {
			return (&url.URL{Scheme: "redis", Host: fmt.Sprintf("%s:%s", c.Redis.Host, port), User: url.UserPassword("", c.Redis.Password)}).String()
		}
		return fmt.Sprintf("redis://%s:%s", c.Redis.Host, port)
	}
	return ""
}
