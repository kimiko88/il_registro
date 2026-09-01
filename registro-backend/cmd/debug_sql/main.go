package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://user:password@localhost:5432/registro?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	fmt.Println("Testing /api/v1/users query...")
	usersQuery := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.role, u.school_id, u.is_active, u.created_at, u.deleted_at, u.pseudonymized_at,
		       s.class_id, c.name, c.section
		FROM users u
		LEFT JOIN students s ON u.id = s.user_id
		LEFT JOIN classes c ON s.class_id = c.id
		WHERE u.deleted_at IS NULL
		LIMIT 1
	`
	rows, err := db.Query(usersQuery)
	if err != nil {
		fmt.Printf("Users Query Error: %v\n", err)
	} else {
		fmt.Println("Users Query: OK")
		_ = rows.Close()
	}

	fmt.Println("\nTesting /api/v1/classes query...")
	classesQuery := `SELECT id, school_id, name, COALESCE(section, ''), academic_year, coordinator_id, created_at, updated_at 
	          FROM classes LIMIT 1`
	rows2, err := db.Query(classesQuery)
	if err != nil {
		fmt.Printf("Classes Query Error: %v\n", err)
	} else {
		fmt.Println("Classes Query: OK")
		_ = rows2.Close()
	}
}
