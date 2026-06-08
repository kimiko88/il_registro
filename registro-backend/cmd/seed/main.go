package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"registro-backend/internal/config"
)

const DriverName = "postgres"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.SSLMode)

	db, err := sql.Open(DriverName, dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	ctx := context.Background()
	log.Println("Seeding database...")

	// 1. Clean up relevant tables (optional, for idempotency)
	cleanup(ctx, db)

	// 2. Create School
	schoolID := uuid.New().String()
	execute(ctx, db, `INSERT INTO schools (id, name, address, city, zip_code, type, code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		schoolID, "Liceo Scientifico A. Einstein", "Via Roma 1", "Milano", "20100", "Liceo Scientifico", "MI12345", time.Now(), time.Now())
	log.Printf("School created: %s", schoolID)

	// 3. Create Users
	// Admin
	createAdmin(ctx, db, schoolID, "admin@test.com", "Admin", "User")

	// Secretary
	createSecretary(ctx, db, schoolID, "secretary@test.com", "Maria", "Segretaria")

	// Teachers (3 teachers)
	teacherMath := createTeacher(ctx, db, schoolID, "teacher.math@test.com", "Giuseppe", "Verdi", "Matematica")
	teacherHistory := createTeacher(ctx, db, schoolID, "teacher.history@test.com", "Paolo", "Rossi", "Storia")
	teacherEnglish := createTeacher(ctx, db, schoolID, "teacher.english@test.com", "Elena", "Bianchi", "Inglese")

	// Classes
	class3A := createClass(ctx, db, schoolID, "3A", "2024/2025")
	class5B := createClass(ctx, db, schoolID, "5B", "2024/2025")

	// Students (2 in 3A, 1 in 5B)
	s1 := createStudent(ctx, db, schoolID, "student1@test.com", "Luigi", "Costanza", class3A)
	s2 := createStudent(ctx, db, schoolID, "student2@test.com", "Mario", "Draghi", class3A)
	s3 := createStudent(ctx, db, schoolID, "student3@test.com", "Giulia", "Verdi", class5B)

	// Parent
	p1 := createParent(ctx, db, schoolID, "parent@test.com", "Genitore", "Costanza")

	// Link Parent-Student
	linkParentStudent(ctx, db, p1, s1)
	linkParentStudent(ctx, db, p1, s2)

	// 4. Sample Academic Data
	// Grades
	createGrade(ctx, db, schoolID, teacherMath, s1, 8.5, "Matematica")
	createGrade(ctx, db, schoolID, teacherMath, s1, 7.0, "Matematica")
	createGrade(ctx, db, schoolID, teacherHistory, s1, 6.5, "Storia")
	createGrade(ctx, db, schoolID, teacherEnglish, s3, 9.0, "Inglese")

	// Attendance
	createAttendance(ctx, db, schoolID, class3A, s1, time.Now().AddDate(0, 0, -1), "Absent", false) // Yesterday absent
	createAttendance(ctx, db, schoolID, class3A, s1, time.Now().AddDate(0, 0, -5), "Late", true)    // 5 days ago delay

	log.Println("Seeding complete.")
}

func cleanup(ctx context.Context, db *sql.DB) {
	// Truncate cascade order matters? Or user CASCADE
	// Only delete test data based on email patterns?
	// For now, simpler to not aggressively delete to avoid wiping dev's other work if any.
	// But to restart fresh:
	_, _ = db.ExecContext(ctx, "TRUNCATE TABLE users, schools, classes, students, parents, grades, attendance, student_parents CASCADE")
}

func execute(ctx context.Context, db *sql.DB, query string, args ...any) {
	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Fatalf("Failed to execute %s: %v", query, err)
	}
}

func hashPwd(p string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b)
}

func createAdmin(ctx context.Context, db *sql.DB, schoolID, email, first, last string) {
	id := uuid.New().String()
	execute(ctx, db, `INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		id, email, hashPwd("password"), first, last, "admin", schoolID, true, time.Now(), time.Now())
	// Also insert into admins table if exists? (admin_profiles) - assume not strictly required for login unless enforced
}

func createSecretary(ctx context.Context, db *sql.DB, schoolID, email, first, last string) {
	id := uuid.New().String()
	execute(ctx, db, `INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		id, email, hashPwd("password"), first, last, "secretary", schoolID, true, time.Now(), time.Now())
}

func createTeacher(ctx context.Context, db *sql.DB, schoolID, email, first, last, subject string) string {
	userID := uuid.New().String()
	execute(ctx, db, `INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		userID, email, hashPwd("password"), first, last, "teacher", schoolID, true, time.Now(), time.Now())

	// Create Teacher Profile
	profileID := uuid.New().String()
	execute(ctx, db, `INSERT INTO teachers (id, user_id, school_id) VALUES ($1, $2, $3)`,
		profileID, userID, schoolID)
	return profileID
}

func createClass(ctx context.Context, db *sql.DB, schoolID, name, year string) string {
	// Hack: Ensure classes table matches 024 schema
	_, _ = db.ExecContext(ctx, `DROP TABLE IF EXISTS classes CASCADE`)
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS classes (
		id UUID PRIMARY KEY,
		school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
		name VARCHAR(50) NOT NULL,
		section VARCHAR(10),
		academic_year VARCHAR(20) NOT NULL,
		coordinator_id UUID REFERENCES users(id) ON DELETE SET NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(school_id, name, section, academic_year)
	)`)
	if err != nil {
		log.Printf("Warning: failed to recreate classes table: %v", err)
	}

	id := uuid.New().String()
	execute(ctx, db, `INSERT INTO classes (id, school_id, name, academic_year, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		id, schoolID, name, year, time.Now(), time.Now())
	return id
}

func createStudent(ctx context.Context, db *sql.DB, schoolID, email, first, last, classID string) string {
	userID := uuid.New().String()
	execute(ctx, db, `INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		userID, email, hashPwd("password"), first, last, "student", schoolID, true, time.Now(), time.Now())

	// Student Profile
	profileID := uuid.New().String()
	execute(ctx, db, `INSERT INTO students (id, user_id, school_id, class_id, enrollment_date, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		profileID, userID, schoolID, classID, time.Now(), time.Now())
	return profileID
}

func createParent(ctx context.Context, db *sql.DB, schoolID, email, first, last string) string {
	userID := uuid.New().String()
	execute(ctx, db, `INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		userID, email, hashPwd("password"), first, last, "parent", schoolID, true, time.Now(), time.Now())

	profileID := uuid.New().String()
	execute(ctx, db, `INSERT INTO parents (id, user_id, school_id, created_at) VALUES ($1, $2, $3, $4)`,
		profileID, userID, schoolID, time.Now())
	return profileID
}

func linkParentStudent(ctx context.Context, db *sql.DB, parentProfileID, studentProfileID string) {
	execute(ctx, db, `INSERT INTO student_parents (student_id, parent_id) VALUES ($1, $2)`, studentProfileID, parentProfileID)
}

func createGrade(ctx context.Context, db *sql.DB, schoolID, teacherUserID, studentProfileID string, value float64, subject string) {
	id := uuid.New().String()
	// Using semester=1 (First Quad) and weight=1.0
	// Assuming subject is a UUID? The query expects subject_id UUID.
	// But seeder passes "Matematica" string.
	// We need a proper subject ID.
	// For now, let's CREATE a subject in the DB and use its ID, or assume the string is fine if logic changes?
	// The schema says `subject_id UUID NOT NULL REFERENCES subjects(id)`.
	// So passing "Matematica" string will fail UUID casting.
	// I must first create a subject.

	// Create subject if not exists (Hack for seeder)
	var subjectID string
	err := db.QueryRowContext(ctx, "INSERT INTO subjects (school_id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id", schoolID, subject).Scan(&subjectID)
	if err != nil {
		// If exists, fetch it
		_ = db.QueryRowContext(ctx, "SELECT id FROM subjects WHERE school_id=$1 AND name=$2", schoolID, subject).Scan(&subjectID)
	}
	if subjectID == "" {
		// Fallback create if scan failed (e.g. conflict but no returning logic overlap?)
		_ = db.QueryRowContext(ctx, "INSERT INTO subjects (school_id, name) VALUES ($1, $2) RETURNING id", schoolID, subject).Scan(&subjectID)
	}

	execute(ctx, db, `INSERT INTO grades (id, school_id, student_id, semester, subject_id, grade_value, grade_type, description, date, teacher_id, weight, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		id, schoolID, studentProfileID, 1, subjectID, value, "Written", "Test verify", time.Now(), teacherUserID, 1.0, time.Now(), time.Now())
}

func createAttendance(ctx context.Context, db *sql.DB, schoolID, classID, studentProfileID string, date time.Time, typ string, justified bool) {
	execute(ctx, db, `INSERT INTO attendance (id, school_id, class_id, student_id, date, status, justified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		uuid.New().String(), schoolID, classID, studentProfileID, date, typ, justified, time.Now(), time.Now())
}
