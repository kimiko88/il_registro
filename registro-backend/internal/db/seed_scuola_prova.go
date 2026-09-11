package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// SeedScuolaDiProva seeds or updates all mock entities for "Scuola di Prova".
func SeedScuolaDiProva(ctx context.Context, dbConn *sql.DB) error {
	// Ensure migration 072 columns exist
	_, _ = dbConn.ExecContext(ctx, `ALTER TABLE students ADD COLUMN IF NOT EXISTS is_representative BOOLEAN DEFAULT FALSE`)
	_, _ = dbConn.ExecContext(ctx, `ALTER TABLE parents ADD COLUMN IF NOT EXISTS is_representative BOOLEAN DEFAULT FALSE`)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	pwdStr := string(passwordHash)

	// 1. Create or fetch "Scuola di Prova"
	var schoolID string
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM schools WHERE name = $1`, "Scuola di Prova").Scan(&schoolID)
	if err != nil {
		schoolID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx, `
			INSERT INTO schools (id, name, address, city, zip_code, type, code, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
		`, schoolID, "Scuola di Prova", "Via delle Prove 10", "Roma", "00100", "Istituto Superiore", "PROVA123")
		if err != nil {
			return fmt.Errorf("failed to create school: %w", err)
		}
		_ = dbConn.QueryRowContext(ctx, `SELECT id FROM schools WHERE code = $1`, "PROVA123").Scan(&schoolID)
		fmt.Printf("[SEED] School created: %s (%s)\n", "Scuola di Prova", schoolID)
	}

	// Helper function for user creation
	createUser := func(email, firstName, lastName, role string) (string, error) {
		var uID string
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&uID)
		if err == nil {
			_, _ = dbConn.ExecContext(ctx, `UPDATE users SET password_hash = $1, is_active = TRUE WHERE id = $2`, pwdStr, uID)
			return uID, nil
		}
		uID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx, `
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE, TRUE, NOW(), NOW())
		`, uID, email, pwdStr, firstName, lastName, role, schoolID)
		if err != nil {
			return "", err
		}
		_, _ = dbConn.ExecContext(ctx, `
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
	dirigenteID, err := createUser("dirigente.prova@scuola.it", "Laura", "Dirigente", "principal")
	if err != nil {
		return fmt.Errorf("failed to create dirigente: %w", err)
	}
	fmt.Printf("[SEED] Admin ID: %s, Secretary ID: %s, Dirigente ID: %s\n", adminID, secID, dirigenteID)

	// 2b. Personale ATA (1 DSGA + almeno 2 per ciascun nuovo ruolo ATA: assistente_amministrativo, collaboratore_ds, collaboratore_scolastico)
	type ataUserData struct {
		email     string
		firstName string
		lastName  string
		role      string
		badgeCode string
	}
	ataStaff := []ataUserData{
		// 1 DSGA
		{"dsga.prova@scuola.it", "Giovanna", "Conti", "dsga", "BADGE-DSGA-001"},
		// 2 Assistenti Amministrativi
		{"assistente1.prova@scuola.it", "Marco", "Ferrari", "assistente_amministrativo", "BADGE-AA-001"},
		{"assistente2.prova@scuola.it", "Lucia", "Romano", "assistente_amministrativo", "BADGE-AA-002"},
		// 2 Collaboratori del Dirigente Scolastico (Collaboratore DS)
		{"collaboratore_ds1.prova@scuola.it", "Roberto", "Mancini", "collaboratore_ds", "BADGE-CDS-001"},
		{"collaboratore_ds2.prova@scuola.it", "Elena", "Galli", "collaboratore_ds", "BADGE-CDS-002"},
		// 2 Collaboratori Scolastici (Portineria / Personale ausiliario)
		{"collaboratore_scolastico1.prova@scuola.it", "Salvatore", "Esposito", "collaboratore_scolastico", "BADGE-CS-001"},
		{"collaboratore_scolastico2.prova@scuola.it", "Carmela", "Russo", "collaboratore_scolastico", "BADGE-CS-002"},
	}

	var hasUserBadgesTable bool
	_ = dbConn.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'user_badges')").Scan(&hasUserBadgesTable)

	for _, ata := range ataStaff {
		ataID, err := createUser(ata.email, ata.firstName, ata.lastName, ata.role)
		if err != nil {
			return fmt.Errorf("failed to create ATA user %s (%s): %w", ata.email, ata.role, err)
		}

		if hasUserBadgesTable {
			_, _ = dbConn.ExecContext(ctx, `
				INSERT INTO user_badges (id, school_id, user_id, badge_code, notes, is_active, assigned_at, created_at)
				VALUES ($1, $2, $3, $4, $5, TRUE, NOW(), NOW())
				ON CONFLICT (school_id, badge_code) DO NOTHING
			`, uuid.New().String(), schoolID, ataID, ata.badgeCode, "Badge di prova ATA")
		}
		fmt.Printf("[SEED] ATA User created: %s (%s, ID: %s)\n", ata.email, ata.role, ataID)
	}

	// 3. 4 Subjects
	subjectNames := []string{"Matematica", "Italiano", "Inglese", "Storia"}
	subjectIDs := make(map[string]string)
	for _, sName := range subjectNames {
		var sID string
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM subjects WHERE school_id = $1 AND name = $2`, schoolID, sName).Scan(&sID)
		if err != nil {
			sID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx, `
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
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM teachers WHERE user_id = $1 AND school_id = $2`, uID, schoolID).Scan(&tProfID)
		if err != nil {
			tProfID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx, `
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
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM classes WHERE school_id = $1 AND name = $2 AND section = $3 AND academic_year = $4`,
			schoolID, name, section, academicYear).Scan(&cID)
		if err == nil {
			// Update coordinator
			if coordinatorUserID != nil {
				_, _ = dbConn.ExecContext(ctx, `UPDATE classes SET coordinator_id = $1 WHERE id = $2`, *coordinatorUserID, cID)
			}
			return cID, nil
		}
		cID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx, `
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
		_, _ = dbConn.ExecContext(ctx, `
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
			err = dbConn.QueryRowContext(ctx, `SELECT id FROM students WHERE user_id = $1 AND school_id = $2`, stUserID, schoolID).Scan(&stProfID)
			if err != nil {
				stProfID = uuid.New().String()
				_, err = dbConn.ExecContext(ctx, `
					INSERT INTO students (id, user_id, school_id, class_id, is_representative, enrollment_number, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
				`, stProfID, stUserID, schoolID, classID, isStudentRep, fmt.Sprintf("MATR-%s-%d", className, i))
				if err != nil {
					return err
				}
			} else {
				_, _ = dbConn.ExecContext(ctx, `UPDATE students SET class_id = $1, is_representative = $2 WHERE id = $3`, classID, isStudentRep, stProfID)
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
			err = dbConn.QueryRowContext(ctx, `SELECT id FROM parents WHERE user_id = $1 AND school_id = $2`, pUserID, schoolID).Scan(&pProfID)
			if err != nil {
				pProfID = uuid.New().String()
				_, err = dbConn.ExecContext(ctx, `
					INSERT INTO parents (id, user_id, school_id, is_representative, created_at)
					VALUES ($1, $2, $3, $4, NOW())
				`, pProfID, pUserID, schoolID, isParentRep)
				if err != nil {
					return err
				}
			} else {
				_, _ = dbConn.ExecContext(ctx, `UPDATE parents SET is_representative = $1 WHERE id = $2`, isParentRep, pProfID)
			}

			// Link Student and Parent
			var linkCount int
			_ = dbConn.QueryRowContext(ctx, `SELECT COUNT(*) FROM student_parents WHERE student_id = $1 AND parent_id = $2`, stProfID, pProfID).Scan(&linkCount)
			if linkCount == 0 {
				_, err = dbConn.ExecContext(ctx, `
					INSERT INTO student_parents (id, student_id, parent_id, relationship_type, can_sign_grades, created_at)
					VALUES ($1, $2, $3, 'Genitore', TRUE, NOW())
				`, uuid.New().String(), stProfID, pProfID)
				if err != nil {
					fmt.Printf("[WARN] student_parents insert error: %v\n", err)
				}
			}

			// 8. Create First Term Provisional Report Card (Scrutinio 1° Semestre) for student
			var scrutinyRecID string
			err = dbConn.QueryRowContext(ctx, `SELECT id FROM scrutiny_records WHERE student_id = $1 AND class_id = $2 AND semester = 1`, stUserID, classID).Scan(&scrutinyRecID)
			if err != nil {
				scrutinyRecID = uuid.New().String()
				_, err = dbConn.ExecContext(ctx, `
					INSERT INTO scrutiny_records (id, student_id, class_id, semester, conduct_grade, final_decision, notes, coordinator_id, status, updated_at)
					VALUES ($1, $2, $3, 1, 8, 'Ammesso', 'Pagella provvisoria 1° Semestre regolare.', $4, 'validated', NOW())
				`, scrutinyRecID, stUserID, classID, teacherUserIDs[0])
				if err != nil {
					fmt.Printf("[WARN] scrutiny_records insert error: %v\n", err)
				}
			}

			// Add scrutiny grades for each subject
			gradeValues := []int{7, 8, 7, 9}
			subIndex := 0
			for _, sName := range subjectNames {
				subID := subjectIDs[sName]
				gVal := gradeValues[subIndex%len(gradeValues)]
				tID := teacherUserIDs[subIndex%len(teacherUserIDs)]
				tProfID := teacherProfileIDs[subIndex%len(teacherProfileIDs)]
				subIndex++

				var sgID string
				err = dbConn.QueryRowContext(ctx, `SELECT id FROM scrutiny_grades WHERE scrutiny_record_id = $1 AND subject_id = $2`, scrutinyRecID, subID).Scan(&sgID)
				if err != nil {
					_, _ = dbConn.ExecContext(ctx, `
						INSERT INTO scrutiny_grades (id, scrutiny_record_id, subject_id, final_grade, teacher_id)
						VALUES ($1, $2, $3, $4, $5)
					`, uuid.New().String(), scrutinyRecID, subID, gVal, tID)
				}

				// Standard grade entry
				var gID string
				err = dbConn.QueryRowContext(ctx, `SELECT id FROM grades WHERE student_id = $1 AND subject_id = $2 AND grade_type = 'Voto Scrutinio'`, stProfID, subID).Scan(&gID)
				if err != nil {
					_, err = dbConn.ExecContext(ctx, `
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

	// 9. Fetch user IDs and student profile IDs for specific test students and parents
	var st2A2_UserID, st2B3_UserID, p2A2_UserID, st2A1_UserID, p2A1_UserID string
	var st2A1_ProfID, st2A2_ProfID string
	_ = dbConn.QueryRowContext(ctx, `SELECT id FROM users WHERE email = 'studente2a_2@scuola.it'`).Scan(&st2A2_UserID)
	_ = dbConn.QueryRowContext(ctx, `SELECT id FROM users WHERE email = 'studente2b_3@scuola.it'`).Scan(&st2B3_UserID)
	_ = dbConn.QueryRowContext(ctx, `SELECT id FROM users WHERE email = 'genitore2a_2@scuola.it'`).Scan(&p2A2_UserID)
	_ = dbConn.QueryRowContext(ctx, `SELECT id FROM users WHERE email = 'studente2a_1@scuola.it'`).Scan(&st2A1_UserID)
	_ = dbConn.QueryRowContext(ctx, `SELECT id FROM users WHERE email = 'genitore2a_1@scuola.it'`).Scan(&p2A1_UserID)

	if st2A1_UserID != "" {
		_ = dbConn.QueryRowContext(ctx, `SELECT id FROM students WHERE user_id = $1`, st2A1_UserID).Scan(&st2A1_ProfID)
	}
	if st2A2_UserID != "" {
		_ = dbConn.QueryRowContext(ctx, `SELECT id FROM students WHERE user_id = $1`, st2A2_UserID).Scan(&st2A2_ProfID)
	}

	// 10. Certificates BES / DSA and PDP Plans
	if st2A2_UserID != "" {
		var cert1ID string
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM certificates WHERE student_id = $1 AND type = 'DSA'`, st2A2_UserID).Scan(&cert1ID)
		if err != nil {
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO certificates (id, school_id, student_id, type, issued_by, issued_at, academic_year, notes, pdf_url, protocol_no)
				VALUES ($1, $2, $3, 'DSA', $4, NOW(), '2024/2025', 'Diagnosi DSA (Dislessia e Disgrafia) accertata dalla ASL', 'https://storage.school.it/cert/dsa_2a2.pdf', 'PROT-2024-00123')
			`, uuid.New().String(), schoolID, st2A2_UserID, secID)
			if err != nil {
				fmt.Printf("[WARN] certificates insert DSA error: %v\n", err)
			}
		}

		var pdp1ID string
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM pdp_plans WHERE student_id = $1`, st2A2_UserID).Scan(&pdp1ID)
		if err != nil {
			pdpContent := `{"compensative": ["tempo_aggiuntivo", "calcolatrice", "mappe_concettuali"], "dispensative": ["lettura_ad_alta_voce"], "notes": "Piano Didattico Personalizzato concordato con la famiglia."}`
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO pdp_plans (id, student_id, class_id, school_id, academic_year, plan_type, diagnosis, content, coordinator_id, shared_with_family, family_approved_at, family_approved_by, created_by, created_at, updated_at)
				VALUES ($1, $2, $3, $4, '2024/2025', 'pdp', 'DSA - Dislessia ed Disgrafia', $5::jsonb, $6, TRUE, NOW(), $7, $6, NOW(), NOW())
			`, uuid.New().String(), st2A2_UserID, class2AID, schoolID, pdpContent, teacherUserIDs[0], p2A2_UserID)
			if err != nil {
				fmt.Printf("[WARN] pdp_plans insert error: %v\n", err)
			}
		}
	}

	if st2B3_UserID != "" {
		var cert2ID string
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM certificates WHERE student_id = $1 AND type = 'BES'`, st2B3_UserID).Scan(&cert2ID)
		if err != nil {
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO certificates (id, school_id, student_id, type, issued_by, issued_at, academic_year, notes, pdf_url, protocol_no)
				VALUES ($1, $2, $3, 'BES', $4, NOW(), '2024/2025', 'Certificazione BES per svantaggio socio-linguistico', 'https://storage.school.it/cert/bes_2b3.pdf', 'PROT-2024-00124')
			`, uuid.New().String(), schoolID, st2B3_UserID, secID)
			if err != nil {
				fmt.Printf("[WARN] certificates insert BES error: %v\n", err)
			}
		}
	}
	fmt.Println("[SEED] BES/DSA Certificates and PDP Plans created.")

	// 11. Class Lessons (Registro di classe)
	lessonsSeed := []struct {
		ClassID, TeacherID, SubjectName, Date, Topic, Type, Notes string
	}{
		{class2AID, teacherUserIDs[0], "Matematica", "2024-10-10", "Equazioni e Disequazioni di 2° grado", "Frontale", "Svolti esercizi da pag. 120 a pag. 125"},
		{class2AID, teacherUserIDs[1], "Italiano", "2024-10-10", "Introduzione al Capitolo III dei Promessi Sposi", "Frontale", "Analisi del personaggio di Don Abbondio"},
		{class2AID, teacherUserIDs[2], "Inglese", "2024-10-11", "Present Perfect vs Past Simple & Conversation Practice", "Laboratorio", "Esercitazione a coppie in laboratorio linguistico"},
		{class2AID, teacherUserIDs[3], "Storia", "2024-10-11", "La Prima Rivoluzione Industriale in Inghilterra", "Frontale", "Visione documentario storico"},
		{class2BID, teacherUserIDs[0], "Matematica", "2024-10-10", "Geometria analitica: la retta nel piano cartesiano", "Frontale", "Spiegazione formula della distanza punto-retta"},
		{class2BID, teacherUserIDs[1], "Italiano", "2024-10-11", "Il Decameron di Boccaccio: Struttura e Tematiche", "Frontale", "Lettura della novella di Andreuccio da Perugia"},
	}

	for _, l := range lessonsSeed {
		subID := subjectIDs[l.SubjectName]
		var lID string
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM class_lessons WHERE class_id = $1 AND teacher_id = $2 AND subject_id = $3 AND date = $4::date`, l.ClassID, l.TeacherID, subID, l.Date).Scan(&lID)
		if err != nil {
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO class_lessons (id, class_id, teacher_id, subject_id, date, topic, type, notes, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5::date, $6, $7, $8, NOW(), NOW())
			`, uuid.New().String(), l.ClassID, l.TeacherID, subID, l.Date, l.Topic, l.Type, l.Notes)
			if err != nil {
				fmt.Printf("[WARN] class_lessons insert error: %v\n", err)
			}
		}
	}
	fmt.Println("[SEED] Class Lessons created.")

	// 12. Attendance Records & Justifications
	if st2A1_ProfID != "" && st2A2_ProfID != "" {
		// Attendance for Student 1 (Present)
		var att1ID string
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM attendance WHERE student_id = $1 AND date = '2024-10-10'::date`, st2A1_ProfID).Scan(&att1ID)
		if err != nil {
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO attendance (id, school_id, student_id, class_id, date, status, justified, created_at, updated_at)
				VALUES ($1, $2, $3, $4, '2024-10-10'::date, 'Present', TRUE, NOW(), NOW())
			`, uuid.New().String(), schoolID, st2A1_ProfID, class2AID)
			if err != nil {
				fmt.Printf("[WARN] attendance 1 insert error: %v\n", err)
			}
		}

		// Attendance for Student 2 (Absent on 2024-10-10, justified)
		var att2ID string
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM attendance WHERE student_id = $1 AND date = '2024-10-10'::date`, st2A2_ProfID).Scan(&att2ID)
		if err != nil {
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO attendance (id, school_id, student_id, class_id, date, status, justified, notes, created_at, updated_at)
				VALUES ($1, $2, $3, $4, '2024-10-10'::date, 'Absent', TRUE, 'Motivi di salute', NOW(), NOW())
			`, uuid.New().String(), schoolID, st2A2_ProfID, class2AID)
			if err != nil {
				fmt.Printf("[WARN] attendance 2 insert error: %v\n", err)
			}
		}

		if p2A2_UserID != "" {
			var justID string
			err = dbConn.QueryRowContext(ctx, `SELECT id FROM justifications WHERE student_id = $1 AND start_date = '2024-10-10'::date`, st2A2_UserID).Scan(&justID)
			if err != nil {
				_, err = dbConn.ExecContext(ctx, `
					INSERT INTO justifications (id, student_id, parent_id, start_date, end_date, reason, status, approved_by, approved_at, created_at, updated_at)
					VALUES ($1, $2, $3, '2024-10-10'::date, '2024-10-10'::date, 'Visita medica specialistica allegata', 'approved', $4, NOW(), NOW(), NOW())
				`, uuid.New().String(), st2A2_UserID, p2A2_UserID, teacherUserIDs[0])
				if err != nil {
					fmt.Printf("[WARN] justifications insert error: %v\n", err)
				}
			}
		}
	}

	fmt.Println("[SEED] Attendance and Justifications created.")

	// 13. Agenda Items & Student Completions
	agendaSeed := []struct {
		ClassID, SubjectName, TeacherID, Title, Description, Type, Date string
	}{
		{class2AID, "Matematica", teacherUserIDs[0], "Verifica Scritta di Matematica", "Verifica su equazioni di secondo grado e problemi di secondo grado", "verifica", "2024-10-25"},
		{class2AID, "Italiano", teacherUserIDs[1], "Compito a casa: Analisi del testo", "Leggere e commentare i capitoli 4 e 5 dei Promessi Sposi", "compito", "2024-10-22"},
		{class2AID, "Inglese", teacherUserIDs[2], "Uscita Didattica Teatro in Lingua", "Spettacolo teatrale in lingua inglese presso il Teatro Olimpico", "evento", "2024-10-28"},
		{class2BID, "Storia", teacherUserIDs[3], "Verifica Scritta di Storia", "Verifica sulla Rivoluzione Industriale ed Illuminismo", "verifica", "2024-10-26"},
	}

	var sampleAgendaItemID string
	for idx, ag := range agendaSeed {
		subID := subjectIDs[ag.SubjectName]
		var agID string
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM agenda_items WHERE class_id = $1 AND title = $2`, ag.ClassID, ag.Title).Scan(&agID)
		if err != nil {
			agID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO agenda_items (id, school_id, class_id, subject_id, teacher_id, title, description, type, date, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::date, NOW(), NOW())
			`, agID, schoolID, ag.ClassID, subID, ag.TeacherID, ag.Title, ag.Description, ag.Type, ag.Date)
			if err != nil {
				fmt.Printf("[WARN] agenda_items insert error: %v\n", err)
			}
		}
		if idx == 1 {
			sampleAgendaItemID = agID
		}
	}

	if st2A1_UserID != "" && sampleAgendaItemID != "" {
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO student_agenda_completions (id, agenda_item_id, student_id, completed_at)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (agenda_item_id, student_id) DO NOTHING
		`, uuid.New().String(), sampleAgendaItemID, st2A1_UserID)
	}
	fmt.Println("[SEED] Agenda Items and Student Completions created.")

	// 14. Communications in Bacheca & Presa d'Atto
	comm1ID := uuid.New().String()
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM communications WHERE subject = 'Circolare 01 - Avvio Anno Scolastico e Regolamento'`).Scan(&comm1ID)
	if err != nil {
		comm1ID = uuid.New().String()
		receivers := []string{st2A1_UserID, p2A1_UserID, st2A2_UserID, p2A2_UserID}
		_, err = dbConn.ExecContext(ctx, `
			INSERT INTO communications (id, sender_id, receiver_ids, subject, body, type, requires_acknowledgment, created_at)
			VALUES ($1, $2, $3, 'Circolare 01 - Avvio Anno Scolastico e Regolamento', 'Si richiede a tutte le famiglie e studenti di visionare ed accettare il nuovo regolamento di disciplina e sicurezza.', 'circolare', TRUE, NOW())
		`, comm1ID, adminID, receivers)
		if err != nil {
			fmt.Printf("[WARN] communications 1 insert error: %v\n", err)
		}
	}

	comm2ID := uuid.New().String()
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM communications WHERE subject = 'Avviso Bacheca - Convocazione Assemblea Genitori'`).Scan(&comm2ID)
	if err != nil {
		comm2ID = uuid.New().String()
		receivers := []string{p2A1_UserID, p2A2_UserID}
		_, err = dbConn.ExecContext(ctx, `
			INSERT INTO communications (id, sender_id, receiver_ids, subject, body, type, requires_acknowledgment, created_at)
			VALUES ($1, $2, $3, 'Avviso Bacheca - Convocazione Assemblea Genitori', 'L assemblea di classe dei genitori è convocata per il prossimo giovedì alle ore 17:00.', 'bacheca', FALSE, NOW())
		`, comm2ID, secID, receivers)
		if err != nil {
			fmt.Printf("[WARN] communications 2 insert error: %v\n", err)
		}
	}

	if p2A1_UserID != "" {
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO communication_acks (id, communication_id, user_id, acknowledged_at)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (communication_id, user_id) DO NOTHING
		`, uuid.New().String(), comm1ID, p2A1_UserID)
	}
	if st2A1_UserID != "" {
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO communication_acks (id, communication_id, user_id, acknowledged_at)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (communication_id, user_id) DO NOTHING
		`, uuid.New().String(), comm1ID, st2A1_UserID)
	}
	fmt.Println("[SEED] Communications and Acknowledgments created.")

	// 15. Orientamento Events, Participations & Preferences
	var orEvt1ID, orEvt2ID, orEvt3ID string
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM orientamento_events WHERE school_id = $1 AND title = 'Open Day Ingegneria e Nuove Tecnologie'`, schoolID).Scan(&orEvt1ID)
	if err != nil {
		orEvt1ID = uuid.New().String()
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO orientamento_events (id, school_id, title, description, category, date, end_date, location, hours, max_attendees, created_by, created_at)
			VALUES ($1, $2, 'Open Day Ingegneria e Nuove Tecnologie', 'Presentazione dei corsi di laurea triennale e magistrale in Ingegneria Informatica e Gestionale', 'University', '2024-11-15 09:00:00+01', '2024-11-15 13:00:00+01', 'Aula Magna - Politecnico', 4.0, 100, $3, NOW())
		`, orEvt1ID, schoolID, adminID)
	}

	err = dbConn.QueryRowContext(ctx, `SELECT id FROM orientamento_events WHERE school_id = $1 AND title = 'Salone dello Studente e Orientamento Post-Diploma'`, schoolID).Scan(&orEvt2ID)
	if err != nil {
		orEvt2ID = uuid.New().String()
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO orientamento_events (id, school_id, title, description, category, date, end_date, location, hours, max_attendees, created_by, created_at)
			VALUES ($1, $2, 'Salone dello Studente e Orientamento Post-Diploma', 'Incontro con università, accademie e aziende per la scelta del percorso futuro', 'Work', '2024-11-20 09:30:00+01', '2024-11-20 13:30:00+01', 'Fiera di Roma - Padiglione 3', 4.0, 250, $3, NOW())
		`, orEvt2ID, schoolID, adminID)
	}

	err = dbConn.QueryRowContext(ctx, `SELECT id FROM orientamento_events WHERE school_id = $1 AND title = 'Workshop Soft Skills & Colloqui di Lavoro'`, schoolID).Scan(&orEvt3ID)
	if err != nil {
		orEvt3ID = uuid.New().String()
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO orientamento_events (id, school_id, title, description, category, date, end_date, location, hours, max_attendees, created_by, created_at)
			VALUES ($1, $2, 'Workshop Soft Skills & Colloqui di Lavoro', 'Simulazione di colloqui di selezione e redazione del curriculum vitae efficace', 'SoftSkills', '2024-12-05 15:00:00+01', '2024-12-05 17:00:00+01', 'Laboratorio Multimediale', 2.0, 30, $3, NOW())
		`, orEvt3ID, schoolID, adminID)
	}

	if st2A1_ProfID != "" {
		if orEvt1ID != "" {
			_, _ = dbConn.ExecContext(ctx, `
				INSERT INTO orientamento_participations (id, event_id, student_id, status, attended, registered_at)
				VALUES ($1, $2, $3, 'Attended', TRUE, NOW() - INTERVAL '10 days')
				ON CONFLICT (event_id, student_id) DO NOTHING
			`, uuid.New().String(), orEvt1ID, st2A1_ProfID)
		}
		if orEvt2ID != "" {
			_, _ = dbConn.ExecContext(ctx, `
				INSERT INTO orientamento_participations (id, event_id, student_id, status, attended, registered_at)
				VALUES ($1, $2, $3, 'Registered', FALSE, NOW() - INTERVAL '2 days')
				ON CONFLICT (event_id, student_id) DO NOTHING
			`, uuid.New().String(), orEvt2ID, st2A1_ProfID)
		}
	}

	if st2A1_UserID != "" {
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO orientamento_preferences (id, student_id, preferred_track, target_field, notes, updated_at)
			VALUES ($1, $2, 'Università / Laurea Triennale', 'Ingegneria Informatica & AI', 'Interessato allo sviluppo software e all intelligenza artificiale.', NOW())
			ON CONFLICT (student_id) DO NOTHING
		`, uuid.New().String(), st2A1_UserID)
	}
	fmt.Println("[SEED] Orientamento Events, Participations & Preferences created.")

	// 16. PCTO Companies, Projects, Participations & Hours
	var companyID string
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM pcto_companies WHERE school_id = $1 AND name = 'Tech Innovators S.r.l.'`, schoolID).Scan(&companyID)
	if err != nil {
		companyID = uuid.New().String()
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO pcto_companies (id, school_id, name, vat_number, address, contact_person, contact_person_first_name, contact_person_last_name, contact_person_phone, email, agreement_date, created_at)
			VALUES ($1, $2, 'Tech Innovators S.r.l.', 'IT12345678901', 'Via dell Innovazione 42, Roma', 'Ing. Mario Rossi', 'Mario', 'Rossi', '+39 06 1234567', 'tutor@techinnovators.it', '2024-09-01', NOW())
		`, companyID, schoolID)
	}

	var pctoProjectID string
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM pcto_projects WHERE school_id = $1 AND title = 'Sviluppo Web e Applicazioni Didattiche Cloud'`, schoolID).Scan(&pctoProjectID)
	if err != nil {
		pctoProjectID = uuid.New().String()
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO pcto_projects (id, school_id, title, description, type, start_date, end_date, total_hours, company_id, school_tutor_id, company_tutor_name, created_by, created_at, updated_at)
			VALUES ($1, $2, 'Sviluppo Web e Applicazioni Didattiche Cloud', 'Progetto di tirocinio aziendale per la progettazione e implementazione di moduli web per la didattica digitale.', 'External', '2024-10-01', '2024-12-20', 60, $3, $4, 'Ing. Mario Rossi', $4, NOW(), NOW())
		`, pctoProjectID, schoolID, companyID, teacherUserIDs[0])
	}

	if st2A1_ProfID != "" && pctoProjectID != "" {
		var pctoPartID string
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM pcto_participations WHERE project_id = $1 AND student_id = $2`, pctoProjectID, st2A1_ProfID).Scan(&pctoPartID)
		if err != nil {
			pctoPartID = uuid.New().String()
			_, _ = dbConn.ExecContext(ctx, `
				INSERT INTO pcto_participations (id, project_id, student_id, status, hours_completed, risk_assessment_ack, created_at)
				VALUES ($1, $2, $3, 'Active', 20.0, TRUE, NOW())
			`, pctoPartID, pctoProjectID, st2A1_ProfID)

			// Seed hour logs for this participation
			hourLogs := []struct {
				Date     string
				Hours    float64
				Activity string
			}{
				{"2024-10-01", 5.0, "Configurazione ambiente di sviluppo e introduzione all architettura software"},
				{"2024-10-08", 5.0, "Analisi dei requisiti funzionali e wireframing interfaccia utente"},
				{"2024-10-15", 5.0, "Implementazione componenti frontend e integrazione API REST"},
				{"2024-10-22", 5.0, "Testing di integrazione, debug e stesura della documentazione tecnica"},
			}
			for _, hl := range hourLogs {
				_, _ = dbConn.ExecContext(ctx, `
					INSERT INTO pcto_hours (id, participation_id, date, hours, activity_description, verified, verified_by, verified_at, created_at)
					VALUES ($1, $2, $3::date, $4, $5, TRUE, $6, NOW(), NOW())
				`, uuid.New().String(), pctoPartID, hl.Date, hl.Hours, hl.Activity, teacherUserIDs[0])
			}
		}
	}
	fmt.Println("[SEED] PCTO Companies, Projects, Participations & Hours created.")

	// 17. Student Goals (Obiettivi Formativi & Badges)
	if st2A1_ProfID != "" {
		goalsSeed := []struct {
			Title, Description, BadgeName, BadgeIcon, Category, Status string
			Points                                                     int
			DueDate                                                    string
		}{
			{
				Title:       "Completare tutti gli esercizi di Matematica sulle Equazioni",
				Description: "Svolgimento completo degli esercizi assegnati per il consolidamento delle equazioni di secondo grado.",
				BadgeName:   "Matematico Brillante",
				BadgeIcon:   "🧮",
				Category:    "academic",
				Status:      "completed",
				Points:      25,
				DueDate:     "2024-10-20",
			},
			{
				Title:       "Firma e visione delle circolari scolastiche entro 48 ore",
				Description: "Costanza nella consultazione della bacheca e firma tempestiva delle circolari informative.",
				BadgeName:   "Lettore Attento",
				BadgeIcon:   "📚",
				Category:    "social",
				Status:      "completed",
				Points:      15,
				DueDate:     "2024-10-25",
			},
			{
				Title:       "Preparazione relazione sul laboratorio scientifico",
				Description: "Stesura accurata del report sperimentale con analisi dati e conclusioni.",
				BadgeName:   "Sperimentatore Provvetto",
				BadgeIcon:   "🧪",
				Category:    "academic",
				Status:      "in_progress",
				Points:      30,
				DueDate:     "2024-11-10",
			},
			{
				Title:       "Partecipazione attiva e diario di bordo PCTO",
				Description: "Compilazione puntuale del registro attività e raggiungimento delle prime 20 ore di tirocinio.",
				BadgeName:   "Futuro Professionista",
				BadgeIcon:   "💼",
				Category:    "behavioral",
				Status:      "in_progress",
				Points:      30,
				DueDate:     "2024-11-30",
			},
		}

		for _, g := range goalsSeed {
			var gID string
			err = dbConn.QueryRowContext(ctx, `SELECT id FROM student_goals WHERE student_id = $1 AND title = $2`, st2A1_ProfID, g.Title).Scan(&gID)
			if err != nil {
				_, _ = dbConn.ExecContext(ctx, `
					INSERT INTO student_goals (id, student_id, teacher_id, title, description, badge_name, badge_icon, category, status, points, due_date, created_at, completed_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::date, NOW(), CASE WHEN $9 = 'completed' THEN NOW() ELSE NULL END)
				`, uuid.New().String(), st2A1_ProfID, teacherUserIDs[0], g.Title, g.Description, g.BadgeName, g.BadgeIcon, g.Category, g.Status, g.Points, g.DueDate)
			}
		}
		fmt.Println("[SEED] Student Goals & Badges created.")
	}

	// 18. Verbali & Modelli Riunioni (Templates con ODG, Verbali Firmati e Bozze Riservate)
	_, _ = dbConn.ExecContext(ctx, `ALTER TABLE council_meetings ALTER COLUMN class_id DROP NOT NULL`)
	_, _ = dbConn.ExecContext(ctx, `ALTER TABLE council_meetings ADD COLUMN IF NOT EXISTS meeting_type VARCHAR(100) DEFAULT 'consiglio_classe'`)
	_, _ = dbConn.ExecContext(ctx, `ALTER TABLE meeting_verbali ADD COLUMN IF NOT EXISTS is_signed BOOLEAN DEFAULT FALSE`)
	_, _ = dbConn.ExecContext(ctx, `ALTER TABLE meeting_verbali ADD COLUMN IF NOT EXISTS signed_at TIMESTAMP WITH TIME ZONE`)
	_, _ = dbConn.ExecContext(ctx, `ALTER TABLE meeting_verbali ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'draft'`)
	_, _ = dbConn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS meeting_verbale_templates (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			meeting_type VARCHAR(100) NOT NULL DEFAULT 'consiglio_classe',
			description TEXT,
			default_agenda TEXT NOT NULL,
			template_content TEXT NOT NULL,
			created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)

	templatesSeed := []struct {
		Title, MeetingType, Description, Agenda, Content string
	}{
		{
			Title:       "Consiglio di Classe — Valutazione Intermedia e Andamento Didattico",
			MeetingType: "consiglio_classe",
			Description: "Modello standard per le sedute periodiche dei consigli di classe con analisi andamento e monitoraggio BES/DSA.",
			Agenda:      "1. Approvazione verbale seduta precedente\n2. Andamento didattico e disciplinare generale della classe\n3. Verifica esiti valutativi intermedi ed eventuali interventi di recupero\n4. Monitoraggio casi particolari, studenti con BES/DSA e verifica PDP/PEI\n5. Varie ed eventuali",
			Content:     "L'anno scolastico 2024/2025, in data odierna, nei locali dell'istituto si è riunito il Consiglio di Classe per discutere l'Ordine del Giorno stabilito.\n\nPresiede la seduta il coordinatore/presidente. Svolge le funzioni di segretario verbalista il docente designato.\n\nPunto 1: Il verbale della seduta precedente viene approvato all'unanimità.\nPunto 2: I docenti relazionano sull'andamento didattico generale. Il clima di classe risulta positivo e partecipativo.\nPunto 3: Vengono concordate strategie di supporto e recupero per gli studenti con lievi incertezze disciplinari.\nPunto 4: Si conferma la piena attuazione dei piani personalizzati (PDP/PEI) concordati.\n\nEsauriti i punti all'ODG, la seduta è tolta.",
		},
		{
			Title:       "Collegio dei Docenti — Delibere Organizzative e Aggiornamento PTOF",
			MeetingType: "collegio_docenti",
			Description: "Schema per il Collegio Docenti plenario presieduto dal Dirigente Scolastico.",
			Agenda:      "1. Approvazione verbale della seduta precedente\n2. Comunicazioni del Dirigente Scolastico\n3. Approvazione aggiornamento annuale Piano Triennale Offerta Formativa (PTOF)\n4. Criteri generali per la valutazione e assegnazione ore di potenziamento\n5. Delibere su viaggi di istruzione ed uscite didattiche",
			Content:     "Nell'Aula Magna dell'istituto, si riunisce il Collegio dei Docenti presieduto dal Dirigente Scolastico.\nSvolge le funzioni di segretario verbalista il docente designato.\n\nConstatata la validità del numero legale, il Presidente apre la seduta trattando i punti all'Ordine del Giorno.\nIl Collegio all'unanimità delibera l'approvazione delle proposte illustrative presentate.\nLa seduta è tolta al termine dei lavori.",
		},
		{
			Title:       "Riunione di Dipartimento Disciplinare — Programmazione e Prove Comuni",
			MeetingType: "dipartimento",
			Description: "Schema per le riunioni per assi culturali e dipartimenti disciplinari.",
			Agenda:      "1. Definizione obiettivi minimi e competenze trasversali\n2. Calendario e struttura delle prove parallele\n3. Monitoraggio adozioni libri di testo e proposte sussidi\n4. Proposte corsi di recupero e progetti di potenziamento",
			Content:     "Nei locali dell'istituto si riunisce il Dipartimento Disciplinare per esaminare i punti all'ODG.\nI docenti presenti concordano all'unanimità le griglie valutative e le tipologie di prove comuni da somministrare.",
		},
	}

	for _, tpl := range templatesSeed {
		var existingTplID string
		_ = dbConn.QueryRowContext(ctx, `SELECT id FROM meeting_verbale_templates WHERE school_id = $1 AND title = $2`, schoolID, tpl.Title).Scan(&existingTplID)
		if existingTplID == "" {
			_, _ = dbConn.ExecContext(ctx, `
				INSERT INTO meeting_verbale_templates (id, school_id, title, meeting_type, description, default_agenda, template_content, created_by, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			`, uuid.New().String(), schoolID, tpl.Title, tpl.MeetingType, tpl.Description, tpl.Agenda, tpl.Content, dirigenteID)
		}
	}
	fmt.Println("[SEED] 3 Verbali & ODG Templates created by Dirigente Scolastica.")

	// Meeting 1: Consiglio di Classe 2A - Verbale UFFICIALE FIRMATO (visibile alla Dirigente e bloccato in sola lettura)
	var meeting1ID string
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM council_meetings WHERE school_id = $1 AND title = $2`, schoolID, "Consiglio di Classe 2A - Periodo Intermedio").Scan(&meeting1ID)
	if err != nil {
		meeting1ID = uuid.New().String()
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO council_meetings (id, school_id, class_id, meeting_type, title, date, start_time, end_time, agenda, created_by, created_at)
			VALUES ($1, $2, $3, 'consiglio_classe', $4, '2024-11-15'::date, '15:00', '16:30', $5, $6, NOW())
		`, meeting1ID, schoolID, class2AID, "Consiglio di Classe 2A - Periodo Intermedio", templatesSeed[0].Agenda, teacherUserIDs[0])

		verbale1ID := uuid.New().String()
		coordUser := teacherUserIDs[0]
		secUser := teacherUserIDs[1]
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO meeting_verbali (id, meeting_id, title, content, secretary_id, president_id, is_published, is_signed, signed_at, status, created_at, updated_at)
			VALUES ($1, $2, 'Verbale n. 1 - Consiglio di Classe 2A (Approvato e Firmato)', $3, $4, $5, TRUE, TRUE, NOW(), 'signed', NOW(), NOW())
		`, verbale1ID, meeting1ID, templatesSeed[0].Content, secUser, coordUser)

		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO verbale_signatures (id, verbale_id, user_id, signed_at, ip_address)
			VALUES ($1, $2, $3, NOW() - INTERVAL '1 hour', '192.168.1.50'),
			       ($4, $2, $5, NOW(), '192.168.1.51')
			ON CONFLICT DO NOTHING
		`, uuid.New().String(), verbale1ID, secUser, uuid.New().String(), coordUser)
	}

	// Meeting 2: Consiglio di Classe Straordinario 2A - Verbale IN BOZZA (modificabile SOLO da Coordinatore e Verbalista, NASCOSTO alla Dirigente)
	var meeting2ID string
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM council_meetings WHERE school_id = $1 AND title = $2`, schoolID, "Consiglio di Classe Straordinario 2A").Scan(&meeting2ID)
	if err != nil {
		meeting2ID = uuid.New().String()
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO council_meetings (id, school_id, class_id, meeting_type, title, date, start_time, end_time, agenda, created_by, created_at)
			VALUES ($1, $2, $3, 'consiglio_classe', $4, '2024-12-05'::date, '16:30', '17:30', '1. Verifica andamento didattico e provvedimenti disciplinari', $5, NOW())
		`, meeting2ID, schoolID, class2AID, "Consiglio di Classe Straordinario 2A", teacherUserIDs[0])

		verbale2ID := uuid.New().String()
		coordUser := teacherUserIDs[0]
		secUser := teacherUserIDs[1]
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO meeting_verbali (id, meeting_id, title, content, secretary_id, president_id, is_published, is_signed, status, created_at, updated_at)
			VALUES ($1, $2, 'Bozza Verbale n. 2 - Consiglio Straordinario 2A', 'Bozza provvisoria in corso di stesura da parte del verbalista...', $3, $4, FALSE, FALSE, 'draft', NOW(), NOW())
		`, verbale2ID, meeting2ID, secUser, coordUser)
	}
	fmt.Println("[SEED] Sample Signed Verbale & Confidential Draft Verbale created.")

	fmt.Println("[SEED] Complete seeding for 'Scuola di Prova' finished successfully!")
	return nil
}
