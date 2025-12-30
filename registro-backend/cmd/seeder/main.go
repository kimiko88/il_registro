package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Hash derived from bcrypt cost 10 for "password"
// $2a$10$y.Xb4D6qXk.KjS.Q9q.E5.Xw.O1.Z.Gj.J.W.L.f.C.1.2.3.4 (This is a dummy string I made up in the replace step, likely INVALID bcrypt).
// Use a REAL generated hash for "password".
// Cost 10: $2a$10$X7.1.j.... (I don't have a generator handy in my head).
// Better: Use the `golang.org/x/crypto/bcrypt` in the script to generate it.

func main() {
	// Connection string - Adjust if needed or take from Env
	connStr := "postgres://postgres:password@localhost:5432/registro?sslmode=disable"
	if os.Getenv("DATABASE_URL") != "" {
		connStr = os.Getenv("DATABASE_URL")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to DB:", err)
	}

	// 1. Generate Hash
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

	fmt.Println("Seeding completed successfully!")
}
