package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"registro-backend/internal/scrutiny"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/registro_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB open error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}

	service := scrutiny.NewService(db)

	classID := "162737ff-081f-436c-8874-11cd57bc60f1"
	actorID := "550e8400-e29b-41d4-a716-446655440001"
	actorRole := "teacher"
	semester := 1

	fmt.Println("=== RUNNING GetMatrix DEBUG ===")
	matrix, err := service.GetMatrix(context.Background(), actorID, actorRole, classID, semester)
	if err != nil {
		fmt.Printf("GetMatrix FAILED: %v\n", err)
	} else {
		fmt.Printf("GetMatrix SUCCESS! Matrix class=%s, students=%d, subjects=%d\n", matrix.ClassID, len(matrix.Students), len(matrix.Subjects))
	}
}
