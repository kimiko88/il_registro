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
	loadEnv()

	// Default connection string
	user := "user"
	pass := "password"
	host := "localhost"
	port := "5432"
	dbname := "registro"
	sslmode := "disable"

	// Try to get from environment variables
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

	fmt.Printf("Connecting to %s:%s/%s...\n", host, port, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database (is it running and port 5432 exposed?): %v", err)
	}

	fmt.Println("✓ Connected to database successfully!")

	// Read all files in migrations directory
	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var migrationFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}

	// Sort files by name
	sort.Strings(migrationFiles)

	fmt.Printf("Found %d migrations. Starting application...\n", len(migrationFiles))

	for _, fileName := range migrationFiles {
		fmt.Printf("[%s] Applying...", fileName)

		content, err := os.ReadFile(filepath.Join("migrations", fileName))
		if err != nil {
			log.Fatalf("\nFailed to read migration %s: %v", fileName, err)
		}

		// Execute the migration SQL
		_, err = db.Exec(string(content))
		if err != nil {
			// Check if it's just an "already exists" error which we can ignore for idempotency
			errMsg := err.Error()
			if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "duplicate") {
				fmt.Println(" (already applied or partially applied)")
			} else {
				fmt.Printf("\n  ⚠ Error applying %s: %v\n", fileName, err)
			}
		} else {
			fmt.Println(" ✓ OK")
		}
	}

	fmt.Println("\n✅ All migrations processed!")
}
