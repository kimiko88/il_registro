package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/registro?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	rows, err := db.Query("SELECT id, email, role, school_id FROM users WHERE role = 'admin'")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = rows.Close() }()

	fmt.Println("Admin users in DB:")
	count := 0
	for rows.Next() {
		var id, email, role string
		var schoolID sql.NullString
		if err := rows.Scan(&id, &email, &role, &schoolID); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- %s (ID: %s) Role: %s, SchoolID: %v\n", email, id, role, schoolID)
		count++
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total found: %d\n", count)
}
