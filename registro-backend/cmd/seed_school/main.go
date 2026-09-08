package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"registro-backend/internal/config"
)

type teacherSeed struct {
	FirstName string
	LastName  string
	Email     string
	Subject   string
}

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

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	ctx := context.Background()
	log.Println("=========================================================")
	log.Println("🏫 SEED ANNO SCOLASTICO COMPLETO - LICEO GALILEO GALILEI")
	log.Println("=========================================================")

	pwdHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash default password: %v", err)
	}
	defaultPwd := string(pwdHash)

	// 1. Scuola
	schoolCode := "RMPS01000P"
	schoolName := "Liceo Scientifico Galileo Galilei"
	var schoolID string
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM schools WHERE code = $1`, schoolCode).Scan(&schoolID)
	if err != nil {
		schoolID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx, `
			INSERT INTO schools (id, name, address, city, zip_code, type, code, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, schoolID, schoolName, "Via Cassia 120", "Roma", "00191", "Liceo Scientifico", schoolCode)
		if err != nil {
			log.Fatalf("Failed to insert school: %v", err)
		}
	}
	log.Printf("✔ Scuola creata/verificata: %s [%s] (ID: %s)\n", schoolName, schoolCode, schoolID)

	// Helper for user creation
	createUser := func(email, first, last, role string) string {
		var uID string
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&uID)
		if err == nil {
			_, _ = dbConn.ExecContext(ctx, `UPDATE users SET password_hash = $1, is_active = TRUE, role = $2 WHERE id = $3`, defaultPwd, role, uID)
			return uID
		}
		uID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx, `
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE, TRUE, NOW(), NOW())
		`, uID, email, defaultPwd, first, last, role, schoolID)
		if err != nil {
			log.Fatalf("Failed to create user %s: %v", email, err)
		}
		return uID
	}

	// 2. Admin & Segreteria
	adminUID := createUser("admin.galilei@scuola.local", "Mario", "Rossi", "admin")
	_ = createUser("segreteria.galilei@scuola.local", "Anna", "Bruni", "secretary")
	log.Printf("✔ Personale direttivo e amministrativo configurato (Admin UID: %s)\n", adminUID)

	// 3. Materie
	subjectsList := []string{
		"Italiano", "Latino", "Matematica", "Fisica",
		"Scienze Naturali", "Storia", "Filosofia", "Inglese",
		"Disegno e Storia dell'Arte", "Scienze Motorie e Sportive",
		"Religione Cattolica", "Informatica", "Diritto ed Economia", "Spagnolo",
	}
	subjectMap := make(map[string]string) // name -> ID
	for _, sName := range subjectsList {
		var sID string
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM subjects WHERE school_id = $1 AND name = $2`, schoolID, sName).Scan(&sID)
		if err != nil {
			sID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO subjects (id, school_id, name, is_votable, default_grade_type)
				VALUES ($1, $2, $3, TRUE, 'numeric')
			`, sID, schoolID, sName)
			if err != nil {
				log.Fatalf("Failed to create subject %s: %v", sName, err)
			}
		}
		subjectMap[sName] = sID
	}
	log.Printf("✔ %d Materie scolastiche registrate\n", len(subjectsList))

	// 4. 20 Docenti con cattedre
	teacherSeeds := []teacherSeed{
		{"Marco", "Bianchi", "m.bianchi@scuola.local", "Italiano"},
		{"Laura", "Conti", "l.conti@scuola.local", "Latino"},
		{"Roberto", "Ferrari", "r.ferrari@scuola.local", "Matematica"},
		{"Elena", "Ricci", "e.ricci@scuola.local", "Fisica"},
		{"Alessandro", "De Luca", "a.deluca@scuola.local", "Scienze Naturali"},
		{"Silvia", "Moretti", "s.moretti@scuola.local", "Scienze Naturali"},
		{"Matteo", "Russo", "m.russo@scuola.local", "Storia"},
		{"Chiara", "Colombo", "c.colombo@scuola.local", "Filosofia"},
		{"Paolo", "Fontana", "p.fontana@scuola.local", "Inglese"},
		{"Giulia", "Santoro", "g.santoro@scuola.local", "Inglese"},
		{"Andrea", "Greco", "a.greco@scuola.local", "Disegno e Storia dell'Arte"},
		{"Valentina", "Marini", "v.marini@scuola.local", "Disegno e Storia dell'Arte"},
		{"Davide", "Barbieri", "d.barbieri@scuola.local", "Scienze Motorie e Sportive"},
		{"Francesca", "Lombardi", "f.lombardi@scuola.local", "Scienze Motorie e Sportive"},
		{"Stefano", "Villa", "s.villa@scuola.local", "Religione Cattolica"},
		{"Sara", "Galli", "s.galli@scuola.local", "Fisica"},
		{"Luca", "Ferri", "l.ferri@scuola.local", "Informatica"},
		{"Martina", "Gatti", "m.gatti@scuola.local", "Spagnolo"},
		{"Simone", "Monti", "s.monti@scuola.local", "Diritto ed Economia"},
		{"Federica", "Pellegrini", "f.pellegrini@scuola.local", "Matematica"},
	}

	teacherUserIDs := make([]string, 0, len(teacherSeeds))
	teacherProfileIDs := make([]string, 0, len(teacherSeeds))
	teacherSubjectMap := make(map[string]string) // teacherProfileID -> subjectID

	for _, ts := range teacherSeeds {
		uID := createUser(ts.Email, ts.FirstName, ts.LastName, "teacher")
		teacherUserIDs = append(teacherUserIDs, uID)

		var tProfID string
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM teachers WHERE user_id = $1 AND school_id = $2`, uID, schoolID).Scan(&tProfID)
		if err != nil {
			tProfID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO teachers (id, user_id, school_id, is_staff, created_at, updated_at)
				VALUES ($1, $2, $3, TRUE, NOW(), NOW())
			`, tProfID, uID, schoolID)
			if err != nil {
				log.Fatalf("Failed to create teacher profile for %s: %v", ts.Email, err)
			}
		}
		teacherProfileIDs = append(teacherProfileIDs, tProfID)
		if subID, ok := subjectMap[ts.Subject]; ok {
			teacherSubjectMap[tProfID] = subID
		}
	}
	log.Printf("✔ %d Docenti d'Istituto creati con profili e cattedre\n", len(teacherSeeds))

	// 5. 5 Classi (1A, 2A, 3A, 4A, 5A)
	classNames := []string{"1A", "2A", "3A", "4A", "5A"}
	classIDs := make([]string, 0, len(classNames))
	for idx, cName := range classNames {
		var cID string
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM classes WHERE school_id = $1 AND name = $2 AND academic_year = '2024/2025'`, schoolID, cName).Scan(&cID)
		if err != nil {
			cID = uuid.New().String()
			coordUID := teacherUserIDs[idx%len(teacherUserIDs)]
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO classes (id, school_id, name, section, academic_year, coordinator_id, created_at, updated_at)
				VALUES ($1, $2, $3, 'A', '2024/2025', $4, NOW(), NOW())
			`, cID, schoolID, cName, coordUID)
			if err != nil {
				log.Fatalf("Failed to create class %s: %v", cName, err)
			}
		}
		classIDs = append(classIDs, cID)
	}
	log.Printf("✔ %d Classi registrate per l'A.S. 2024/2025 con coordinatori\n", len(classNames))

	// 6. 120 Studenti (24 per classe) e Genitori
	firstNamesM := []string{"Leonardo", "Francesco", "Alessandro", "Lorenzo", "Mattia", "Andrea", "Gabriele", "Riccardo", "Tommaso", "Edoardo", "Federico", "Matteo", "Diego", "Giuseppe", "Niccolò"}
	firstNamesF := []string{"Sofia", "Aurora", "Giulia", "Ginevra", "Beatrice", "Alice", "Vittoria", "Emma", "Giorgia", "Martina", "Chiara", "Greta", "Ludovica", "Anna", "Camilla"}
	lastNames := []string{"Rossi", "Russo", "Ferrari", "Esposito", "Bianchi", "Romano", "Colombo", "Ricci", "Marino", "Greco", "Bruno", "Gallo", "Conti", "De Luca", "Mancini", "Costa", "Giordano", "Rizzo", "Lombardi", "Moretti", "Barbieri", "Fontana", "Santoro", "Mariani"}

	type studentRecord struct {
		UserID    string
		ProfileID string
		ClassID   string
		FullName  string
	}
	allStudents := make([]studentRecord, 0, 120)
	parentCount := 0

	studentCounter := 1
	for cIdx, cID := range classIDs {
		cName := classNames[cIdx]
		for sInClass := 1; sInClass <= 24; sInClass++ {
			isMale := (studentCounter % 2) == 1
			var first string
			if isMale {
				first = firstNamesM[(studentCounter+sInClass)%len(firstNamesM)]
			} else {
				first = firstNamesF[(studentCounter+sInClass)%len(firstNamesF)]
			}
			last := lastNames[(studentCounter*3+sInClass)%len(lastNames)]
			email := fmt.Sprintf("studente.%s.%02d@scuola.local", cName, sInClass)

			sUID := createUser(email, first, last, "student")

			// Check or create student profile
			var sProfID string
			err := dbConn.QueryRowContext(ctx, `SELECT id FROM students WHERE user_id = $1 AND school_id = $2`, sUID, schoolID).Scan(&sProfID)
			if err != nil {
				sProfID = uuid.New().String()
				matr := fmt.Sprintf("MAT-2024-%03d", studentCounter)
				_, err = dbConn.ExecContext(ctx, `
					INSERT INTO students (id, user_id, school_id, class_id, enrollment_number, enrollment_date, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, '2024-09-01', NOW(), NOW())
				`, sProfID, sUID, schoolID, cID, matr)
				if err != nil {
					log.Fatalf("Failed to insert student profile: %v", err)
				}
			} else {
				_, _ = dbConn.ExecContext(ctx, `UPDATE students SET class_id = $1 WHERE id = $2`, cID, sProfID)
			}

			allStudents = append(allStudents, studentRecord{
				UserID:    sUID,
				ProfileID: sProfID,
				ClassID:   cID,
				FullName:  first + " " + last,
			})

			// Genitore associato
			pEmail := fmt.Sprintf("genitore.%s.%02d@famiglia.local", cName, sInClass)
			pFirst := "Genitore"
			pLast := last
			pUID := createUser(pEmail, pFirst, pLast, "parent")

			var pProfID string
			err = dbConn.QueryRowContext(ctx, `SELECT id FROM parents WHERE user_id = $1 AND school_id = $2`, pUID, schoolID).Scan(&pProfID)
			if err != nil {
				pProfID = uuid.New().String()
				_, _ = dbConn.ExecContext(ctx, `
					INSERT INTO parents (id, user_id, school_id, created_at)
					VALUES ($1, $2, $3, NOW())
				`, pProfID, pUID, schoolID)
			}

			// Link parent to student
			_, _ = dbConn.ExecContext(ctx, `
				INSERT INTO student_parents (student_id, parent_id, relationship_type, can_sign_grades)
				VALUES ($1, $2, 'Guardian', TRUE)
				ON CONFLICT DO NOTHING
			`, sProfID, pProfID)
			parentCount++

			studentCounter++
		}
	}
	log.Printf("✔ %d Studenti (24 per classe) e %d Genitori collegati tramite student_parents\n", len(allStudents), parentCount)

	// 7. Calendario Presenze su 3 mesi (dal 16 settembre ad oggi, ~60 giorni scolastici)
	rng := rand.New(rand.NewSource(42))
	baseDate := time.Now().AddDate(0, -3, 0)
	totalAttendance := 0

	// Pre-generate school dates (Mon-Fri)
	var schoolDays []time.Time
	for d := baseDate; d.Before(time.Now()); d = d.AddDate(0, 0, 1) {
		if d.Weekday() >= time.Monday && d.Weekday() <= time.Friday {
			schoolDays = append(schoolDays, d)
		}
	}

	for _, day := range schoolDays {
		dateStr := day.Format("2006-01-02")
		for _, st := range allStudents {
			roll := rng.Float64()
			status := "present"
			isJustified := false
			minutesLate := 0

			if roll > 0.96 {
				status = "absent"
				isJustified = rng.Float64() > 0.3
			} else if roll > 0.94 {
				status = "late"
				minutesLate = 10 + rng.Intn(35)
				isJustified = true
			} else if roll > 0.93 {
				status = "early_exit"
				isJustified = true
			}

			teacherUID := teacherUserIDs[rng.Intn(len(teacherUserIDs))]
			_, err := dbConn.ExecContext(ctx, `
				INSERT INTO attendance (id, school_id, student_id, class_id, teacher_id, date, status, minutes_late, is_justified, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
				ON CONFLICT DO NOTHING
			`, uuid.New().String(), schoolID, st.UserID, st.ClassID, teacherUID, dateStr, status, minutesLate, isJustified)
			if err == nil {
				totalAttendance++
			}
		}
	}
	log.Printf("✔ %d Registrazioni presenze generate sui %d giorni scolastici (3 mesi)\n", totalAttendance, len(schoolDays))

	// 8. 500+ Voti con distribuzione gaussiana realistica
	totalGrades := 0
	evalTypes := []string{"Written", "Oral", "Practical"}
	weights := []float64{1.0, 1.0, 0.75, 1.25}

	for _, st := range allStudents {
		// Assegna tra 4 e 6 voti per studente in diverse materie
		numGrades := 4 + rng.Intn(3)
		for g := 0; g < numGrades; g++ {
			// Campiona voto con distribuzione normale (media 6.8, dev. standard 1.3)
			rawGrade := rng.NormFloat64()*1.3 + 6.8
			if rawGrade < 4.0 {
				rawGrade = 4.0
			}
			if rawGrade > 10.0 {
				rawGrade = 10.0
			}
			// Arrotonda al quarto di punto (es. 6.25, 6.5, 6.75, 7.0)
			roundedGrade := math.Round(rawGrade*4.0) / 4.0

			tProfID := teacherProfileIDs[rng.Intn(len(teacherProfileIDs))]
			subID := teacherSubjectMap[tProfID]
			if subID == "" {
				subID = subjectMap["Matematica"]
			}
			evalType := evalTypes[rng.Intn(len(evalTypes))]
			weight := weights[rng.Intn(len(weights))]
			gDate := schoolDays[rng.Intn(len(schoolDays))]

			_, err := dbConn.ExecContext(ctx, `
				INSERT INTO grades (
					id, school_id, student_id, subject_id, teacher_id,
					grade_value, grade_type, evaluation_type, semester, date,
					weight, is_published, published_at, created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, 'numeric', $7, 1, $8, $9, TRUE, NOW(), NOW(), NOW())
			`, uuid.New().String(), schoolID, st.ProfileID, subID, tProfID, roundedGrade, evalType, gDate, weight)
			if err == nil {
				totalGrades++
			}
		}
	}
	log.Printf("✔ %d Voti inseriti con distribuzione gaussiana realistica e tipologia valutativa\n", totalGrades)

	// 9. 10 Circolari d'Istituto con firme di presa visione
	circulars := []struct {
		Subject string
		Body    string
	}{
		{"Circolare n. 01 - Avvio a.s. 2024/2025 e disposizioni organizzative", "Si comunica l'orario delle lezioni per la prima settimana di attività scolastica e le modalità di ingresso."},
		{"Circolare n. 02 - Elezioni organi collegiali e rappresentanti di classe", "Indizione delle elezioni dei rappresentanti dei genitori e degli studenti nei consigli di classe."},
		{"Circolare n. 03 - Ricevimento antimeridiano e colloqui generali docenti", "Pubblicazione del calendario dei colloqui individuali con i docenti su piattaforma registro elettronico."},
		{"Circolare n. 04 - Somministrazione prove Invalsi e calendario simulazioni", "Indicazioni operative per lo svolgimento delle prove nazionali Invalsi a cura del dipartimento scientifico."},
		{"Circolare n. 05 - Viaggi di istruzione e uscite didattiche a.s. 2024/2025", "Presentazione delle mete approvate dal Collegio Docenti e termini per la presentazione delle adesioni."},
		{"Circolare n. 06 - Progetto PCTO e convenzioni aziendali triennio", "Modalità di avvio dei percorsi per le competenze trasversali e l'orientamento per le classi 3^, 4^ e 5^."},
		{"Circolare n. 07 - Norme di sicurezza, piano di evacuazione e divieto di fumo", "Richiamo alle norme di sicurezza d'istituto e calendarizzazione della prima prova di evacuazione."},
		{"Circolare n. 08 - Corsi di recupero e sportello didattico pomeridiano", "Attivazione degli interventi di recupero e consolidamento per Matematica, Fisica e Lingua Inglese."},
		{"Circolare n. 09 - Scrutini del primo periodo didattico e chiusura trimestre", "Calendario delle operazioni di scrutinio e criteri per l'attribuzione delle valutazioni intermedie."},
		{"Circolare n. 10 - Giornata della Scienza e conferenze magistrali d'Istituto", "Programma degli interventi laboratoriali e seminari con docenti universitari presso l'Aula Magna."},
	}

	totalSigs := 0
	for cIdx, circ := range circulars {
		cID := uuid.New().String()
		_, err := dbConn.ExecContext(ctx, `
			INSERT INTO communications (id, sender_id, receiver_ids, subject, body, type, created_at)
			VALUES ($1, $2, $3, $4, $5, 'circular', NOW() - INTERVAL '1 day' * $6)
			ON CONFLICT DO NOTHING
		`, cID, adminUID, "{student,parent,teacher}", circ.Subject, circ.Body, (10-cIdx)*7)
		if err != nil {
			continue
		}

		// Simulazione firme di presa visione per l'80% degli studenti/genitori
		for _, st := range allStudents {
			if rng.Float64() < 0.80 {
				_, sErr := dbConn.ExecContext(ctx, `
					INSERT INTO communication_signatures (id, communication_id, user_id, signed_at)
					VALUES ($1, $2, $3, NOW() - INTERVAL '1 hour' * $4)
					ON CONFLICT DO NOTHING
				`, uuid.New().String(), cID, st.UserID, rng.Intn(72)+1)
				if sErr == nil {
					totalSigs++
				}
			}
		}
	}
	log.Printf("✔ %d Circolari d'Istituto pubblicate con %d firme di presa visione registrate\n", len(circulars), totalSigs)

	log.Println("=========================================================")
	log.Println("🎉 POPOLAMENTO ANNO SCOLASTICO COMPLETATO CON SUCCESSO!")
	log.Println("Credenziali di test pronte per l'accesso:")
	log.Println(" - Dirigente / Admin: admin.galilei@scuola.local / password")
	log.Println(" - Segreteria:        segreteria.galilei@scuola.local / password")
	log.Println(" - Docenti (x20):     m.bianchi@scuola.local, r.ferrari@scuola.local, ... / password")
	log.Println(" - Studenti (x120):   studente.1a.01@scuola.local, ... / password")
	log.Println(" - Genitori (x120):   genitore.1a.01@famiglia.local, ... / password")
	log.Println("=========================================================")
}
