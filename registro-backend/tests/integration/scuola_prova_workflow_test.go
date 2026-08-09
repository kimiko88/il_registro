package integration

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"registro-backend/internal/attendance"
	"registro-backend/internal/classes"
	"registro-backend/internal/config"
	dbPkg "registro-backend/internal/db"
	"registro-backend/internal/grades"
	"registro-backend/internal/scrutiny"
	"registro-backend/internal/users"
)

func getTestDB(t *testing.T) *sql.DB {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		if _, err := os.Stat("../../.env"); err == nil {
			_ = os.Chdir("../../")
		} else if _, err := os.Stat("../.env"); err == nil {
			_ = os.Chdir("../")
		}
	}

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
	if err != nil || db.Ping() != nil {
		t.Skip("Skipping integration test: PostgreSQL DB unavailable")
		return nil
	}

	var hasSchools bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'schools')").Scan(&hasSchools)
	if err != nil || !hasSchools {
		t.Skip("Skipping integration test: 'schools' table does not exist in DB")
		return nil
	}

	_, _ = db.Exec("ALTER TABLE classes ADD COLUMN IF NOT EXISTS location VARCHAR(150) DEFAULT ''")
	_, _ = db.Exec("ALTER TABLE grades ADD COLUMN IF NOT EXISTS evaluation_type VARCHAR(20) DEFAULT 'Written'")
	_, _ = db.Exec("ALTER TABLE grades ADD COLUMN IF NOT EXISTS test_id UUID")
	_, _ = db.Exec("ALTER TABLE grades ADD COLUMN IF NOT EXISTS compensative_measures JSONB DEFAULT '[]'::jsonb")

	return db
}

func TestScuolaDiProvaWorkflow(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	ctx := context.Background()

	// Seed / Update Scuola di Prova data
	err := dbPkg.SeedScuolaDiProva(ctx, db)
	require.NoError(t, err, "Seeding Scuola di Prova must succeed")

	// 1. Verify "Scuola di Prova"
	var schoolID string
	err = db.QueryRowContext(ctx, `SELECT id FROM schools WHERE name = $1`, "Scuola di Prova").Scan(&schoolID)
	require.NoError(t, err, "Scuola di Prova must exist in DB")

	assert.NotEmpty(t, schoolID)

	// 2. Verify Admin and Segreteria
	var adminID, secID string
	err = db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1 AND school_id = $2`, "admin.prova@scuola.it", schoolID).Scan(&adminID)
	require.NoError(t, err, "Admin user must exist")
	err = db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1 AND school_id = $2`, "segreteria.prova@scuola.it", schoolID).Scan(&secID)
	require.NoError(t, err, "Segreteria user must exist")

	// 3. Verify 4 Subjects
	var subjectCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM subjects WHERE school_id = $1`, schoolID).Scan(&subjectCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, subjectCount, 4, "Should have at least 4 subjects")

	// 4. Verify 4 Teachers & Coordinator assignment
	var teacher1ID, teacher2ID string
	err = db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1`, "docente1@scuola.it").Scan(&teacher1ID)
	require.NoError(t, err, "Docente 1 must exist")
	err = db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1`, "docente2@scuola.it").Scan(&teacher2ID)
	require.NoError(t, err, "Docente 2 must exist")

	// 5. Verify Classes 2A and 2B
	var class2AID, class2BCoord string
	err = db.QueryRowContext(ctx, `SELECT id, coordinator_id FROM classes WHERE school_id = $1 AND name = '2A' AND section = 'A'`, schoolID).Scan(&class2AID, &class2BCoord)
	require.NoError(t, err, "Class 2A must exist")
	assert.Equal(t, teacher1ID, class2BCoord, "Docente 1 must be Coordinator of 2A")

	var class2BID string
	err = db.QueryRowContext(ctx, `SELECT id FROM classes WHERE school_id = $1 AND name = '2B' AND section = 'B'`, schoolID).Scan(&class2BID)
	require.NoError(t, err, "Class 2B must exist")

	// 6. Verify Students & Class Representative flag
	var count2A, count2B int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM students WHERE class_id = $1`, class2AID).Scan(&count2A)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count2A, 10, "Class 2A should have at least 10 students")

	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM students WHERE class_id = $1`, class2BID).Scan(&count2B)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count2B, 10, "Class 2B should have at least 10 students")

	var studentRepCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM students WHERE class_id = $1 AND is_representative = TRUE`, class2AID).Scan(&studentRepCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, studentRepCount, 1, "Class 2A should have at least 1 student representative")

	// 7. Verify Parents & Parent Representatives
	var parentRepCount2A, parentRepCount2B int
	err = db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT p.id) 
		FROM parents p 
		JOIN student_parents sp ON sp.parent_id = p.id 
		JOIN students s ON sp.student_id = s.id 
		WHERE s.class_id = $1 AND p.is_representative = TRUE
	`, class2AID).Scan(&parentRepCount2A)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, parentRepCount2A, 1, "Class 2A should have at least 1 parent representative")

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT p.id) 
		FROM parents p 
		JOIN student_parents sp ON sp.parent_id = p.id 
		JOIN students s ON sp.student_id = s.id 
		WHERE s.class_id = $1 AND p.is_representative = TRUE
	`, class2BID).Scan(&parentRepCount2B)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, parentRepCount2B, 1, "Class 2B should have at least 1 parent representative")

	// 8. Verify Scrutinio Authorization & Matrix (Coordinator vs Non-Coordinator)
	scrutinyRepo := scrutiny.NewRepository(db)
	gradesRepo := grades.NewRepository(db)
	classesRepo := classes.NewRepository(db)
	usersRepo := users.NewRepository(db)
	attRepo := attendance.NewRepository(db)

	scrutinySvc := scrutiny.NewService(scrutinyRepo, gradesRepo, classesRepo, usersRepo, attRepo)

	// Teacher 1 (Coordinator of 2A) MUST succeed in getting matrix
	matrix, err := scrutinySvc.GetMatrix(ctx, teacher1ID, "teacher", class2AID, 1)
	require.NoError(t, err, "Coordinator (Teacher 1) should be authorized to view 2A scrutiny matrix")
	assert.NotNil(t, matrix)
	assert.GreaterOrEqual(t, len(matrix.Students), 10, "Matrix for 2A should contain at least 10 students")
	assert.GreaterOrEqual(t, len(matrix.Subjects), 4, "Matrix for 2A should contain 4 subjects")

	// Teacher 2 (Not Coordinator of 2A) MUST get ErrUnauthorizedScrutiny on modification/saving
	err = scrutinySvc.SaveScrutiny(ctx, teacher2ID, "teacher", scrutiny.SaveScrutinyRequest{ClassID: class2AID, Semester: 1, StudentID: matrix.Students[0].StudentID})
	assert.ErrorIs(t, err, scrutiny.ErrUnauthorizedScrutiny, "Non-coordinator (Teacher 2) must be forbidden from saving 2A scrutiny")

	// 9. Verify Parent access to child's provisional report card (Pagella provvisoria)
	var parent1UserID string
	err = db.QueryRowContext(ctx, `SELECT id FROM users WHERE LOWER(email) = LOWER($1)`, "genitore2a_1@scuola.it").Scan(&parent1UserID)
	require.NoError(t, err, "Parent 1 of 2A must exist")

	children, err := usersRepo.GetChildren(ctx, parent1UserID)
	require.NoError(t, err, "GetChildren must succeed for parent")
	assert.GreaterOrEqual(t, len(children), 1, "Parent 1 should be linked to at least 1 child")

	// Check child's provisional report card (Pagella provvisoria) in scrutiny records
	scRec, err := scrutinyRepo.GetRecord(ctx, children[0].UserID, class2AID, 1)
	require.NoError(t, err, "GetRecord must succeed for child's scrutiny")
	require.NotNil(t, scRec, "Child should have a scrutiny record for 1° Semestre")
	assert.GreaterOrEqual(t, len(scRec.Grades), 4, "Child should have provisional report card grades for all subjects")

	// 10. Verify BES / DSA Certificates and PDP Plans
	var certCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM certificates WHERE school_id = $1`, schoolID).Scan(&certCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, certCount, 2, "Should have at least 2 BES/DSA certificates in Scuola di Prova")

	var pdpCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pdp_plans WHERE school_id = $1 AND shared_with_family = TRUE`, schoolID).Scan(&pdpCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, pdpCount, 1, "Should have at least 1 shared PDP plan in Scuola di Prova")

	// 11. Verify Class Lessons
	var lessonCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM class_lessons WHERE class_id IN ($1, $2)`, class2AID, class2BID).Scan(&lessonCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, lessonCount, 6, "Should have at least 6 class lessons recorded")

	// 12. Verify Attendance & Justifications
	var attCount, justCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM attendance WHERE school_id = $1`, schoolID).Scan(&attCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, attCount, 2, "Should have attendance records")

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM justifications j 
		JOIN users u ON j.student_id = u.id 
		WHERE u.school_id = $1 AND j.status = 'approved'
	`, schoolID).Scan(&justCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, justCount, 1, "Should have at least 1 approved justification")

	// 13. Verify Agenda Items & Student Completions
	var agendaCount, completionCount int
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM agenda_items WHERE school_id = $1`, schoolID).Scan(&agendaCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, agendaCount, 4, "Should have at least 4 agenda items")

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM student_agenda_completions sac
		JOIN agenda_items ai ON sac.agenda_item_id = ai.id
		WHERE ai.school_id = $1
	`, schoolID).Scan(&completionCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, completionCount, 1, "Should have at least 1 student agenda completion")

	// 14. Verify Communications & Acknowledgments (Bacheca)
	var commCount, ackCount int
	err = db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM communications c
		JOIN users u ON c.sender_id = u.id
		WHERE u.school_id = $1
	`, schoolID).Scan(&commCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, commCount, 2, "Should have at least 2 bacheca/circolare communications")

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM communication_acks ca
		JOIN communications c ON ca.communication_id = c.id
		JOIN users u ON c.sender_id = u.id
		WHERE u.school_id = $1
	`, schoolID).Scan(&ackCount)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, ackCount, 2, "Should have at least 2 communication acknowledgments (presa d'atto)")

	fmt.Println("=== ALL INTEGRATION CHECKS PASSED PERFECTLY ===")
}

