package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Supabase SupabaseConfig
	SPID     SPIDConfig
	CIE      CIEConfig
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
	URL string
	Key string
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

	// If .env file exists, read it, otherwise ignore error (for docker env vars)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: error reading config file: %s", err)
	}

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_MODE", "debug")
	viper.SetDefault("DB_SSLMODE", "disable")

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
			Mode: viper.GetString("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
		Supabase: SupabaseConfig{
			URL: viper.GetString("SUPABASE_URL"),
			Key: viper.GetString("SUPABASE_KEY"),
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
	}

	return config, nil
}
