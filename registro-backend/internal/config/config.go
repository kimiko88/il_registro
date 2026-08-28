package config

import (
	"log"
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
	Host string
	Port string
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
