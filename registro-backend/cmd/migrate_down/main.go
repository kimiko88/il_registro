package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func getDB() (*sql.DB, error) {
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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, sslmode)
	if os.Getenv("DATABASE_URL") != "" {
		dsn = os.Getenv("DATABASE_URL")
	}

	return sql.Open("postgres", dsn)
}

var (
	createIndexRegex = regexp.MustCompile(`(?i)CREATE\s+INDEX\s+(?:IF\s+NOT\s+EXISTS\s+)?([a-zA-Z0-9_]+)`)
	createTableRegex = regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([a-zA-Z0-9_]+)`)
	addColumnRegex   = regexp.MustCompile(`(?i)ALTER\s+TABLE\s+([a-zA-Z0-9_]+)\s+ADD\s+COLUMN\s+(?:IF\s+NOT\s+EXISTS\s+)?([a-zA-Z0-9_]+)`)
)

func generateDownStatements(sqlContent string) []string {
	var downStmts []string

	// Check for created indexes -> DROP INDEX IF EXISTS
	indexMatches := createIndexRegex.FindAllStringSubmatch(sqlContent, -1)
	for _, m := range indexMatches {
		if len(m) > 1 {
			downStmts = append(downStmts, fmt.Sprintf("DROP INDEX IF EXISTS %s;", m[1]))
		}
	}

	// Check for added columns -> ALTER TABLE ... DROP COLUMN IF EXISTS
	colMatches := addColumnRegex.FindAllStringSubmatch(sqlContent, -1)
	for _, m := range colMatches {
		if len(m) > 2 {
			downStmts = append(downStmts, fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s;", m[1], m[2]))
		}
	}

	// Check for created tables -> DROP TABLE IF EXISTS ... CASCADE
	tableMatches := createTableRegex.FindAllStringSubmatch(sqlContent, -1)
	for _, m := range tableMatches {
		if len(m) > 1 {
			downStmts = append(downStmts, fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", m[1]))
		}
	}

	return downStmts
}

func main() {
	dryRun := flag.Bool("dry-run", false, "Preview rollback SQL statements without executing")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: go run cmd/migrate_down/main.go [-dry-run] <migration_file.sql>")
		fmt.Println("\nRecent migrations in ./migrations/:")
		files, err := filepath.Glob("./migrations/*.sql")
		if err == nil && len(files) > 0 {
			sort.Strings(files)
			start := 0
			if len(files) > 8 {
				start = len(files) - 8
			}
			for _, f := range files[start:] {
				fmt.Printf("  - %s\n", filepath.Base(f))
			}
		}
		os.Exit(0)
	}

	migrationPath := args[0]
	data, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration file %s: %v", migrationPath, err)
	}

	downStatements := generateDownStatements(string(data))
	if len(downStatements) == 0 {
		fmt.Printf("ℹ️ No auto-reversibile DDL objects (indexes, columns, tables) detected in %s\n", filepath.Base(migrationPath))
		return
	}

	fmt.Printf("Generated %d rollback statement(s) for %s:\n", len(downStatements), filepath.Base(migrationPath))
	for i, stmt := range downStatements {
		fmt.Printf("  [%d] %s\n", i+1, stmt)
	}

	if *dryRun {
		fmt.Println("\n✓ Dry-run completed: no database modifications executed.")
		return
	}

	db, err := getDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	for _, stmt := range downStatements {
		if _, err := tx.Exec(stmt); err != nil {
			_ = tx.Rollback()
			log.Fatalf("Rollback statement failed [%s]: %v", stmt, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit rollback transaction: %v", err)
	}

	fmt.Printf("\n✓ Rollback executed successfully for %s!\n", filepath.Base(migrationPath))
}
