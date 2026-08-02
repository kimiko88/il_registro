package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"registro-backend/internal/config"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		cfg, err := config.LoadConfig()
		if err == nil && cfg.Database.Host != "" {
			dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				url.QueryEscape(cfg.Database.User),
				url.QueryEscape(cfg.Database.Password),
				cfg.Database.Host,
				cfg.Database.Port,
				cfg.Database.Name,
				cfg.Database.SSLMode,
			)
		} else {
			dbURL = "postgres://postgres:postgres@localhost:5432/registro_db?sslmode=disable"
		}
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	ctx := context.Background()
	log.Println("=== Seeding 'Scuola di Prova' ===")

	if err := SeedScuolaDiProva(ctx, db); err != nil {
		log.Fatalf("Seeding error: %v", err)
	}

	log.Println("=== Seeding completed successfully! ===")
}

func SeedScuolaDiProva(ctx context.Context, db *sql.DB) error {
	// Ensure migration 072 columns exist
	_, _ = db.ExecContext(ctx, `ALTER TABLE students ADD COLUMN IF NOT EXISTS is_representative BOOLEAN DEFAULT FALSE`)
	_, _ = db.ExecContext(ctx, `ALTER TABLE parents ADD COLUMN IF NOT EXISTS is_representative BOOLEAN DEFAULT FALSE`)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	pwdStr := string(passwordHash)

	// 1. Create or fetch "Scuola di Prova"
	var schoolID string
	err = db.QueryRowContext(ctx, `SELECT id FROM schools WHERE name = $1`, "Scuola di Prova").Scan(&schoolID)
	if err != nil {
		schoolID = uuid.New().String()
		_, err = db.ExecContext(ctx, `
			INSERT INTO schools (id, name, address, city, zip_code, type, code, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		`, schoolID, "Scuola di Prova", "Via delle Prove 10", "Roma", "00100", "Istituto Superiore", "PROVA123")
		if err != nil {
			return fmt.Errorf("failed to create school: %w", err)
		}
		fmt.Printf("[SEED] School created: %s (%s)\n", "Scuola di Prova", schoolID)
	}

	// Helper function for user creation
	createUser := func(email, firstName, lastName, role string) (string, error) {
		var uID string
		err := db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&uID)
		if err == nil {
			_, _ = db.ExecContext(ctx, `UPDATE users SET password_hash = $1, is_active = TRUE WHERE id = $2`, pwdStr, uID)
			return uID, nil
		}
		uID = uuid.New().String()
		_, err = db.ExecContext(ctx, `
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE, TRUE, NOW(), NOW())
		`, uID, email, pwdStr, firstName, lastName, role, schoolID)
		if err != nil {
			return "", err
		}
		_, _ = db.ExecContext(ctx, `
			INSERT INTO user_roles (id, user_id, role, school_id, created_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (user_id, role, school_id) DO NOTHING
		`, uuid.New().String(), uID, role, schoolID)
		return uID, nil
	}

	// 2. Admin & Segreteria
	adminID, err := createUser("admin.prova@scuola.it", "Admin", "ScuolaProva", "admin")
	if err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}
	secID, err := createUser("segreteria.prova@scuola.it", "Segreteria", "ScuolaProva", "secretary")
	if err != nil {
		return fmt.Errorf("failed to create secretary: %w", err)
	}
	fmt.Printf("[SEED] Admin ID: %s, Secretary ID: %s\n", adminID, secID)

	// 3. 4 Subjects
	subjectNames := []string{"Matematica", "Italiano", "Inglese", "Storia"}
	subjectIDs := make(map[string]string)
	for _, sName := range subjectNames {
		var sID string
		err := db.QueryRowContext(ctx, `SELECT id FROM subjects WHERE school_id = $1 AND name = $2`, schoolID, sName).Scan(&sID)
		if err != nil {
			sID = uuid.New().String()
			_, err = db.ExecContext(ctx, `
				INSERT INTO subjects (id, school_id, name, code, is_mandatory, created_at)
				VALUES ($1, $2, $3, $4, TRUE, NOW())
			`, sID, schoolID, sName, sName[:3])
			if err != nil {
				return fmt.Errorf("failed to create subject %s: %w", sName, err)
			}
		}
		subjectIDs[sName] = sID
	}
	fmt.Println("[SEED] 4 Subjects ready:", subjectNames)

	// 4. 4 Teachers
	teacherData := []struct {
		Email, First, Last string
	}{
		{"docente1@scuola.it", "Marco", "Rossi"},
		{"docente2@scuola.it", "Laura", "Bianchi"},
		{"docente3@scuola.it", "Giuseppe", "Verdi"},
		{"docente4@scuola.it", "Elena", "Neri"},
	}
	teacherUserIDs := make([]string, 4)
	teacherProfileIDs := make([]string, 4)

	for i, t := range teacherData {
		uID, err := createUser(t.Email, t.First, t.Last, "teacher")
		if err != nil {
			return fmt.Errorf("failed to create teacher user %s: %w", t.Email, err)
		}
		teacherUserIDs[i] = uID

		var tProfID string
		err = db.QueryRowContext(ctx, `SELECT id FROM teachers WHERE user_id = $1 AND school_id = $2`, uID, schoolID).Scan(&tProfID)
		if err != nil {
			tProfID = uuid.New().String()
			_, err = db.ExecContext(ctx, `
				INSERT INTO teachers (id, user_id, school_id, hiring_date, created_at, updated_at)
				VALUES ($1, $2, $3, CURRENT_DATE, NOW(), NOW())
			`, tProfID, uID, schoolID)
			if err != nil {
				return fmt.Errorf("failed to create teacher profile for %s: %w", t.Email, err)
			}
		}
		teacherProfileIDs[i] = tProfID
	}
	fmt.Println("[SEED] 4 Teachers ready. Teacher 1 User ID (Coordinator):", teacherUserIDs[0])

	// 5. 2 Classes (2A and 2B). Teacher 1 is coordinator for 2A
	createClass := func(name, section, academicYear string, coordinatorUserID *string) (string, error) {
		var cID string
		err := db.QueryRowContext(ctx, `SELECT id FROM classes WHERE school_id = $1 AND name = $2 AND section = $3 AND academic_year = $4`,
			schoolID, name, section, academicYear).Scan(&cID)
		if err == nil {
			// Update coordinator
			if coordinatorUserID != nil {
				_, _ = db.ExecContext(ctx, `UPDATE classes SET coordinator_id = $1 WHERE id = $2`, *coordinatorUserID, cID)
			}
			return cID, nil
		}
		cID = uuid.New().String()
		_, err = db.ExecContext(ctx, `
			INSERT INTO classes (id, school_id, name, section, academic_year, coordinator_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		`, cID, schoolID, name, section, academicYear, coordinatorUserID)
		return cID, err
	}

	coordID2A := teacherUserIDs[0]
	class2AID, err := createClass("2A", "A", "2024/2025", &coordID2A)
	if err != nil {
		return fmt.Errorf("failed to create class 2A: %w", err)
	}
	class2BID, err := createClass("2B", "B", "2024/2025", nil)
	if err != nil {
		return fmt.Errorf("failed to create class 2B: %w", err)
	}
	fmt.Printf("[SEED] Classes created: 2A (%s, Coord: %s), 2B (%s)\n", class2AID, coordID2A, class2BID)

	// 6. Assign Subjects to Teachers (class_subjects)
	// Docente 1: Matematica in 2A & 2B
	// Docente 2: Italiano in 2A & 2B
	// Docente 3: Inglese in 2A & 2B
	// Docente 4: Storia in 2A & 2B
	assignments := []struct {
		ClassID, SubjectName string
		TeacherProfID        string
	}{
		{class2AID, "Matematica", teacherProfileIDs[0]},
		{class2AID, "Italiano", teacherProfileIDs[1]},
		{class2AID, "Inglese", teacherProfileIDs[2]},
		{class2AID, "Storia", teacherProfileIDs[3]},

		{class2BID, "Matematica", teacherProfileIDs[0]},
		{class2BID, "Italiano", teacherProfileIDs[1]},
		{class2BID, "Inglese", teacherProfileIDs[2]},
		{class2BID, "Storia", teacherProfileIDs[3]},
	}

	for _, a := range assignments {
		subID := subjectIDs[a.SubjectName]
		_, _ = db.ExecContext(ctx, `
			INSERT INTO class_subjects (id, class_id, subject_id, teacher_id, hours_per_week, created_at)
			VALUES ($1, $2, $3, $4, 4.0, NOW())
			ON CONFLICT (class_id, subject_id, teacher_id) DO NOTHING
		`, uuid.New().String(), a.ClassID, subID, a.TeacherProfID)
	}
	fmt.Println("[SEED] Subject-Teacher assignments created for 2A and 2B.")

	// 7. Create 10 Students & 10 Parents for 2A, and 10 Students & 10 Parents for 2B
	seedClassPeople := func(classID, className string) error {
		cNameLower := strings.ToLower(className)
		for i := 1; i <= 10; i++ {
			// Student
			stEmail := fmt.Sprintf("studente%s_%d@scuola.it", cNameLower, i)
			stUserFirstName := fmt.Sprintf("Studente%s_%d", className, i)
			stUserLastName := "Test"
			isStudentRep := (className == "2A" && i == 1) // Student 1 in 2A is class representative

			stUserID, err := createUser(stEmail, stUserFirstName, stUserLastName, "student")
			if err != nil {
				return err
			}

			var stProfID string
			err = db.QueryRowContext(ctx, `SELECT id FROM students WHERE user_id = $1 AND school_id = $2`, stUserID, schoolID).Scan(&stProfID)
			if err != nil {
				stProfID = uuid.New().String()
				_, err = db.ExecContext(ctx, `
					INSERT INTO students (id, user_id, school_id, class_id, is_representative, enrollment_number, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
				`, stProfID, stUserID, schoolID, classID, isStudentRep, fmt.Sprintf("MATR-%s-%d", className, i))
				if err != nil {
					return err
				}
			} else {
				_, _ = db.ExecContext(ctx, `UPDATE students SET class_id = $1, is_representative = $2 WHERE id = $3`, classID, isStudentRep, stProfID)
			}

			// Parent
			pEmail := fmt.Sprintf("genitore%s_%d@scuola.it", cNameLower, i)
			pFirstName := fmt.Sprintf("Genitore%s_%d", className, i)
			pLastName := "Test"
			isParentRep := (i == 1) // Parent 1 in 2A and Parent 1 in 2B are parent representatives

			pUserID, err := createUser(pEmail, pFirstName, pLastName, "parent")
			if err != nil {
				return err
			}

			var pProfID string
			err = db.QueryRowContext(ctx, `SELECT id FROM parents WHERE user_id = $1 AND school_id = $2`, pUserID, schoolID).Scan(&pProfID)
			if err != nil {
				pProfID = uuid.New().String()
				_, err = db.ExecContext(ctx, `
					INSERT INTO parents (id, user_id, school_id, is_representative, created_at)
					VALUES ($1, $2, $3, $4, NOW())
				`, pProfID, pUserID, schoolID, isParentRep)
				if err != nil {
					return err
				}
			} else {
				_, _ = db.ExecContext(ctx, `UPDATE parents SET is_representative = $1 WHERE id = $2`, isParentRep, pProfID)
			}

			// Link Student and Parent
			var linkCount int
			_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM student_parents WHERE student_id = $1 AND parent_id = $2`, stProfID, pProfID).Scan(&linkCount)
			if linkCount == 0 {
				_, err = db.ExecContext(ctx, `
					INSERT INTO student_parents (id, student_id, parent_id, relationship_type, can_sign_grades, created_at)
					VALUES ($1, $2, $3, 'Genitore', TRUE, NOW())
				`, uuid.New().String(), stProfID, pProfID)
				if err != nil {
					fmt.Printf("[WARN] student_parents insert error: %v\n", err)
				}
			}

			// 8. Create First Term Provisional Report Card (Scrutinio 1° Semestre) for student
			var scrutinyRecID string
			err = db.QueryRowContext(ctx, `SELECT id FROM scrutiny_records WHERE student_id = $1 AND class_id = $2 AND semester = 1`, stUserID, classID).Scan(&scrutinyRecID)
			if err != nil {
				scrutinyRecID = uuid.New().String()
				_, err = db.ExecContext(ctx, `
					INSERT INTO scrutiny_records (id, student_id, class_id, semester, conduct_grade, final_decision, notes, coordinator_id, status, updated_at)
					VALUES ($1, $2, $3, 1, 8, 'Ammesso', 'Pagella provvisoria 1° Semestre regolare.', $4, 'validated', NOW())
				`, scrutinyRecID, stUserID, classID, teacherUserIDs[0])
				if err != nil {
					fmt.Printf("[WARN] scrutiny_records insert error: %v\n", err)
				}
			}

			// Add scrutiny grades for each subject
			gradeValues := []int{7, 8, 7, 9} // Grades for Math, Ita, Eng, Sto
			subIndex := 0
			for _, sName := range subjectNames {
				subID := subjectIDs[sName]
				gVal := gradeValues[subIndex%len(gradeValues)]
				tID := teacherUserIDs[subIndex%len(teacherUserIDs)]
				tProfID := teacherProfileIDs[subIndex%len(teacherProfileIDs)]
				subIndex++

				var sgID string
				err = db.QueryRowContext(ctx, `SELECT id FROM scrutiny_grades WHERE scrutiny_record_id = $1 AND subject_id = $2`, scrutinyRecID, subID).Scan(&sgID)
				if err != nil {
					_, _ = db.ExecContext(ctx, `
						INSERT INTO scrutiny_grades (id, scrutiny_record_id, subject_id, final_grade, teacher_id)
						VALUES ($1, $2, $3, $4, $5)
					`, uuid.New().String(), scrutinyRecID, subID, gVal, tID)
				}

				// Also add standard grade entry into grades table for student profile id
				var gID string
				err = db.QueryRowContext(ctx, `SELECT id FROM grades WHERE student_id = $1 AND subject_id = $2 AND grade_type = 'Voto Scrutinio'`, stProfID, subID).Scan(&gID)
				if err != nil {
					_, err = db.ExecContext(ctx, `
						INSERT INTO grades (id, student_id, school_id, subject_id, teacher_id, grade_value, grade_type, semester, is_published, description, date, created_at, updated_at)
						VALUES ($1, $2, $3, $4, $5, $6, 'Voto Scrutinio', 1, TRUE, 'Pagella Provvisoria 1° Semestre', CURRENT_DATE, NOW(), NOW())
					`, uuid.New().String(), stProfID, schoolID, subID, tProfID, float64(gVal))
					if err != nil {
						fmt.Printf("[WARN] grades insert error: %v\n", err)
					}
				}
			}
		}
		return nil
	}

	if err := seedClassPeople(class2AID, "2A"); err != nil {
		return fmt.Errorf("failed seeding 2A people: %w", err)
	}
	if err := seedClassPeople(class2BID, "2B"); err != nil {
		return fmt.Errorf("failed seeding 2B people: %w", err)
	}

	fmt.Println("[SEED] 20 Students, 20 Parents, and 1st Term Report Cards successfully created!")
	return nil
}
