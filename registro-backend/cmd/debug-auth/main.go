package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Using the same credentials from .env
	DB_CONN = "user=user password=password dbname=registro sslmode=disable host=localhost port=5432"
)

func main() {
	// 1. Connect to DB
	db, err := sql.Open("postgres", DB_CONN)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	fmt.Println("Connected to Database.")

	// 2. Fetch User
	email := "testuser@example.com"
	password := "Password123!"

	var storedHash string
	err = db.QueryRow("SELECT password_hash FROM users WHERE email = $1", email).Scan(&storedHash)
	if err != nil {
		log.Fatalf("User not found or query error: %v", err)
	}
	fmt.Printf("Found user %s. Hash: %s\n", email, storedHash)

	// 3. Verify Password
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Printf("MATCH FAIL: %v\n", err)

		// Debug: Generate new hash and see if it looks similar (cost etc)
		newHash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
		fmt.Printf("Expected hash format (generated new): %s\n", newHash)
	} else {
		fmt.Println("MATCH SUCCESS! Password is correct.")
	}
}
