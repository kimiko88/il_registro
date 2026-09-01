package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to localhost:5432
	// Assumes default credentials from docker-compose if .env not available
	connStr := "host=localhost port=5432 user=user password=password dbname=registro sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatal("Ping failed (makes sure DB is up and exposed on 5432): ", err)
	}

	// Read migration file
	// Assuming running from registro-backend root or cmd/migrate_notes
	// We'll try absolute path or relative
	path := "migrations/025_add_student_notes.sql"
	content, err := os.ReadFile(path)
	if err != nil {
		// Try stepping back if running from cmd
		path = "../../migrations/025_add_student_notes.sql"
		content, err = os.ReadFile(path)
		if err != nil {
			log.Fatal("Could not find migration file: ", err)
		}
	}

	// Execute
	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	fmt.Println("Migration applied successfully!")
}
