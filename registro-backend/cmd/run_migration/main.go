package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/run_migration/main.go <migration_file_path>")
	}
	path := os.Args[1]

	loadEnv()

	user := "user"
	pass := "password"
	host := "localhost"
	port := "5432"
	dbname := "registro"
	sslmode := "disable"

	if os.Getenv("DB_USER") != "" {
		user = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		pass = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_HOST") != "" {
		host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		port = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_NAME") != "" {
		dbname = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_SSLMODE") != "" {
		sslmode = os.Getenv("DB_SSLMODE")
	}

	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "http://")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbname, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not find migration file %s: %v", path, err)
	}

	// Split SQL statements by semicolon so lib/pq executes each statement individually
	// without wrapping multi-statement scripts in an implicit transaction block (required for CREATE INDEX CONCURRENTLY).
	rawStatements := strings.Split(string(contentBytes), ";")
	for _, stmt := range rawStatements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		if _, err := db.Exec(trimmed); err != nil {
			log.Fatalf("Migration failed on statement [%s]: %v", trimmed, err)
		}
	}

	fmt.Printf("✓ Migration %s applied successfully!\n", path)
}
