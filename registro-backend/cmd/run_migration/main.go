package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/run_migration/main.go <migration_file_path>")
	}
	path := os.Args[1]

	connStr := "host=localhost port=5432 user=user password=password dbname=registro sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Ping failed: ", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Could not find migration file: ", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	fmt.Printf("Migration %s applied successfully!\n", path)
}
