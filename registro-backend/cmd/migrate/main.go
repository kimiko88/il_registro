package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

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
			_ = os.Setenv(key, val)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("loadEnv scanner error: %v", err)
	}
}

func findMigrationsDir() string {
	candidates := []string{"migrations", "../migrations", "../../migrations"}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return "migrations"
}

func main() {
	loadEnv()

	user := getEnvOrDefault("DB_USER", "user")
	pass := getEnvOrDefault("DB_PASSWORD", "password")
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	dbname := getEnvOrDefault("DB_NAME", "registro")
	sslmode := getEnvOrDefault("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbname, sslmode)
	fmt.Printf("Connecting to PostgreSQL at %s:%s/%s...\n", host, port, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	fmt.Println("✓ Connected successfully!")

	// Ensure schema_migrations table exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Failed to create schema_migrations table: %v", err)
	}

	migrationsDir := findMigrationsDir()
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory '%s': %v", migrationsDir, err)
	}

	var migrationFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Single target file execution if argument provided
	if len(os.Args) >= 2 {
		targetArg := os.Args[1]
		targetFile := filepath.Base(targetArg)
		applyMigrationFile(db, migrationsDir, targetFile)
		fmt.Println("\n✅ Single migration finished!")
		return
	}

	fmt.Printf("Found %d migration files in '%s'. Checking pending migrations...\n", len(migrationFiles), migrationsDir)

	appliedCount := 0
	for _, fileName := range migrationFiles {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", fileName).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check migration status for %s: %v", fileName, err)
		}

		if exists {
			fmt.Printf("[%s] Already applied (skipping)\n", fileName)
			continue
		}

		fmt.Printf("[%s] Applying...", fileName)
		if applyMigrationFile(db, migrationsDir, fileName) {
			appliedCount++
		}
	}

	fmt.Printf("\n✅ Migrations completed! Applied %d new migration(s).\n", appliedCount)
}

func applyMigrationFile(db *sql.DB, dir, fileName string) bool {
	filePath := filepath.Join(dir, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("\n  ⚠ Error reading %s: %v\n", filePath, err)
		return false
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("\n  ⚠ Error starting transaction for %s: %v\n", fileName, err)
		return false
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(string(content)); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "duplicate") {
			fmt.Printf(" (already applied or table exists — recording version)\n")
			_, _ = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1) ON CONFLICT DO NOTHING", fileName)
			_ = tx.Commit()
			return true
		}
		log.Printf("\n  ⚠ Error executing %s: %v\n", fileName, err)
		return false
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1) ON CONFLICT DO NOTHING", fileName); err != nil {
		log.Printf("\n  ⚠ Error recording version for %s: %v\n", fileName, err)
		return false
	}

	if err := tx.Commit(); err != nil {
		log.Printf("\n  ⚠ Error committing transaction for %s: %v\n", fileName, err)
		return false
	}

	fmt.Println(" ✓ OK")
	return true
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
