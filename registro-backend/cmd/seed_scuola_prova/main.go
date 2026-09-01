package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/lib/pq"

	"registro-backend/internal/config"
	"registro-backend/internal/db"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		cfg, err := config.LoadConfig()
		if err == nil && cfg.Database.Host != "" {
			dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				url.QueryEscape(cfg.Database.User),
				url.QueryEscape(cfg.Database.Password),
				cfg.Database.Host,
				cfg.Database.Port,
				cfg.Database.Name,
				cfg.Database.SSLMode,
			)
		} else {
			dbURL = "postgres://postgres:postgres@localhost:5432/registro_db?sslmode=disable"
		}
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer func() { _ = conn.Close() }()

	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	ctx := context.Background()
	log.Println("=== Seeding 'Scuola di Prova' ===")

	if err := db.SeedScuolaDiProva(ctx, conn); err != nil {
		log.Fatalf("Seeding error: %v", err)
	}

	log.Println("=== Seeding completed successfully! ===")
}
