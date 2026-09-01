package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Hash derived from bcrypt cost 10 for "password"
// $2a$10$y.Xb4D6qXk.KjS.Q9q.E5.Xw.O1.Z.Gj.J.W.L.f.C.1.2.3.4 (This is a dummy string I made up in the replace step, likely INVALID bcrypt).
// Use a REAL generated hash for "password".
// Cost 10: $2a$10$X7.1.j.... (I don't have a generator handy in my head).
// Better: Use the `golang.org/x/crypto/bcrypt` in the script to generate it.

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
}

func main() {
	loadEnv()

	user := "postgres"
	pass := "password"
	host := "localhost"
	port := "5432"
	dbname := "postgres"
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
	if os.Getenv("DATABASE_URL") != "" {
		connStr = os.Getenv("DATABASE_URL")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to DB:", err)
	}

	// 1. Create schools table if not exists
	schoolsTable := `
	CREATE TABLE IF NOT EXISTS schools (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		code VARCHAR(50) NOT NULL UNIQUE,
		address VARCHAR(255),
		city VARCHAR(100),
		province VARCHAR(50),
		zip_code VARCHAR(10),
		phone VARCHAR(50),
		email VARCHAR(255),
		principal VARCHAR(255),
		type VARCHAR(50),
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	_, err = db.Exec(schoolsTable)
	if err != nil {
		log.Printf("Warning: schools table creation: %v", err)
	}

	// Add code column if not exists
	addCodeColumn := `ALTER TABLE schools ADD COLUMN IF NOT EXISTS code VARCHAR(50);`
	_, err = db.Exec(addCodeColumn)
	if err != nil {
		log.Printf("Warning: adding code column: %v", err)
	}

	// 2. Generate Hash
	pwd := "password"
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated Hash for '%s': %s\n", pwd, hash)

	// 2. Insert Users
	// Upsert to avoid dupes
	query := `
	INSERT INTO users (id, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at)
	VALUES 
	($1, $2, $3, 'teacher', 'Mario', 'Verdi', true, NOW(), NOW()),
	($4, $5, $3, 'student', 'Luigi', 'Rossi', true, NOW(), NOW()),
	($6, $7, $3, 'parent', 'Giulia', 'Rossi', true, NOW(), NOW())
	ON CONFLICT (email) DO UPDATE 
	SET password_hash = EXCLUDED.password_hash, is_active = true;
	`
	// Note: IDs should be UUIDs usually. If schema uses UUID, we need to generate them.
	// Looking at `migrations/001_initial_schema.sql` might tell us.
	// Assuming text or uuid. I'll use simple string IDs if allowed, or random UUIDs.
	// Better to use UUIDs to be safe.

	tID := "00000000-0000-0000-0000-000000000001"
	sID := "00000000-0000-0000-0000-000000000002"
	pID := "00000000-0000-0000-0000-000000000003"
	aID := "00000000-0000-0000-0000-000000000004"

	// Check if ID column is UUID type?
	// I'll try execute. If it fails on UUID format, I'll fix.

	_, err = db.Exec(query,
		tID, "teacher@school.it", string(hash),
		sID, "student@school.it",
		pID, "parent@school.it",
	)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	// Admin Query
	adminQuery := `
	INSERT INTO users (id, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, 'admin', 'Super', 'Admin', true, NOW(), NOW())
	ON CONFLICT (email) DO UPDATE 
	SET password_hash = EXCLUDED.password_hash, is_active = true, role = 'admin';
	`
	_, err = db.Exec(adminQuery, aID, "admin@school.it", string(hash))
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	fmt.Println("Seeding completed successfully!")
}
