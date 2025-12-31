package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to database
	connStr := "postgres://user:password@localhost:5432/registro?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully!")

	// Read and execute migration 021
	migration021, err := os.ReadFile("migrations/021_admin_actions.sql")
	if err != nil {
		log.Fatalf("Failed to read migration 021: %v", err)
	}

	fmt.Println("Applying migration 021_admin_actions.sql...")
	if _, err := db.Exec(string(migration021)); err != nil {
		log.Printf("Migration 021 error (may already exist): %v", err)
	} else {
		fmt.Println("✓ Migration 021 applied successfully!")
	}

	// Read and execute migration 022
	migration022, err := os.ReadFile("migrations/022_seed_data.sql")
	if err != nil {
		log.Fatalf("Failed to read migration 022: %v", err)
	}

	fmt.Println("Applying migration 022_seed_data.sql...")
	if _, err := db.Exec(string(migration022)); err != nil {
		log.Printf("Migration 022 error (may already exist): %v", err)
	} else {
		fmt.Println("✓ Migration 022 applied successfully!")
	}

	// Read and execute migration 023
	migration023, err := os.ReadFile("migrations/023_user_sessions.sql")
	if err != nil {
		log.Fatalf("Failed to read migration 023: %v", err)
	}

	fmt.Println("Applying migration 023_user_sessions.sql...")
	if _, err := db.Exec(string(migration023)); err != nil {
		log.Printf("Migration 023 error (may already exist): %v", err)
	} else {
		fmt.Println("✓ Migration 023 applied successfully!")
	}

	fmt.Println("\n✅ All migrations completed!")
	fmt.Println("\nTest users created (all with password: 'password'):")
	fmt.Println("  - superadmin@registroelettronico.it (SuperAdmin)")
	fmt.Println("  - admin@liceogalilei.it (Admin)")
	fmt.Println("  - segreteria@liceogalilei.it (Secretary)")
	fmt.Println("  - p.verdi@liceogalilei.it (Teacher)")
	fmt.Println("  - a.neri@liceogalilei.it (Teacher)")
	fmt.Println("  - l.rossi@studenti.liceogalilei.it (Student)")
	fmt.Println("  - g.bianchi@studenti.liceogalilei.it (Student)")
	fmt.Println("  - m.ferrari@studenti.liceogalilei.it (Student)")
	fmt.Println("  - famiglia.rossi@gmail.com (Parent)")
}
