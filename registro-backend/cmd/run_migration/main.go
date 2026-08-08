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

// splitSQLStatements splits a multi-statement SQL script by semicolons,
// while correctly preserving semicolons inside single quotes, dollar-quoted blocks ($$ ... $$),
// and multi-line/single-line comments.
func splitSQLStatements(sqlText string) []string {
	var statements []string
	var current strings.Builder

	inSingleQuote := false
	inDollarQuote := false
	dollarTag := ""
	inLineComment := false
	inBlockComment := false

	runes := []rune(sqlText)
	n := len(runes)

	for i := 0; i < n; i++ {
		ch := runes[i]

		// Line comments (-- ...)
		if !inSingleQuote && !inDollarQuote && !inBlockComment {
			if !inLineComment && ch == '-' && i+1 < n && runes[i+1] == '-' {
				inLineComment = true
			}
		}

		if inLineComment {
			current.WriteRune(ch)
			if ch == '\n' {
				inLineComment = false
			}
			continue
		}

		// Block comments (/* ... */)
		if !inSingleQuote && !inDollarQuote && !inLineComment {
			if !inBlockComment && ch == '/' && i+1 < n && runes[i+1] == '*' {
				inBlockComment = true
				current.WriteRune(ch)
				current.WriteRune(runes[i+1])
				i++
				continue
			} else if inBlockComment && ch == '*' && i+1 < n && runes[i+1] == '/' {
				inBlockComment = false
				current.WriteRune(ch)
				current.WriteRune(runes[i+1])
				i++
				continue
			}
		}

		if inBlockComment {
			current.WriteRune(ch)
			continue
		}

		// Single quotes ('...')
		if !inDollarQuote {
			if ch == '\'' {
				if inSingleQuote && i+1 < n && runes[i+1] == '\'' {
					current.WriteRune(ch)
					current.WriteRune('\'')
					i++
					continue
				}
				inSingleQuote = !inSingleQuote
				current.WriteRune(ch)
				continue
			}
		}

		// Dollar-quoted strings ($tag$ ... $tag$)
		if !inSingleQuote {
			if ch == '$' {
				j := i + 1
				for j < n && (runes[j] == '_' || (runes[j] >= 'a' && runes[j] <= 'z') || (runes[j] >= 'A' && runes[j] <= 'Z') || (runes[j] >= '0' && runes[j] <= '9')) {
					j++
				}
				if j < n && runes[j] == '$' {
					tag := string(runes[i : j+1])
					if !inDollarQuote {
						inDollarQuote = true
						dollarTag = tag
						current.WriteString(tag)
						i = j
						continue
					} else if inDollarQuote && tag == dollarTag {
						inDollarQuote = false
						dollarTag = ""
						current.WriteString(tag)
						i = j
						continue
					}
				}
			}
		}

		// Statement terminator (;)
		if ch == ';' && !inSingleQuote && !inDollarQuote {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
			continue
		}

		current.WriteRune(ch)
	}

	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
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

	statements := splitSQLStatements(string(contentBytes))
	for _, trimmed := range statements {
		if _, err := db.Exec(trimmed); err != nil {
			log.Fatalf("Migration failed on statement [%s]: %v", trimmed, err)
		}
	}

	fmt.Printf("✓ Migration %s applied successfully!\n", path)
}
