package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"registro-backend/internal/config"
	"registro-backend/pkg/logger"

	_ "github.com/lib/pq"
)

func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := strings.TrimPrefix(strings.TrimPrefix(cfg.Host, "https://"), "http://")
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Connection Pool Settings (Tuned for PgBouncer & High Concurrency)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)


	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Log.Info("Connected to database successfully")
	return db, nil
}
