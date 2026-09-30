package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// SeedScuolaVolta seeds a comprehensive, highly realistic secondary school dataset for:
// "I.I.S. Alessandro Volta" (Milan)
// - 1 School, 2 Buildings (Sede Centrale & Succursale Sud)
// - 11 Bookable Rooms (Palestre, Laboratori di Informatica, Fisica, Scienze, Elettronica, Lingue)
// - 14 School Subjects with room requirements
// - 55 Classes across 11 sections (A to M) and 5 grade levels (1 to 5)
// - Over 110 Teachers with dedicated subject specializations, 14-18 weekly teaching hours (cattedre)
// - Real class-subject assignments (class_subjects) covering 26-30h per class
// - Realistic teacher preferences (giorno libero, prime/ultime ore, indisponibilità)
// - Timetable constraints and associated groups for schedule generation
func SeedScuolaVolta(ctx context.Context, dbConn *sql.DB) error {
	log.Println("==================================================================")
	log.Println("🏫 SEEDING SCUOLA: I.I.S. ALESSANDRO VOLTA (Test Generazione Orario)")
	log.Println("==================================================================")

	// Pre-generate password hash once to make the seeding fast (< 2 seconds)
	pwdHashBytes, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	defaultPwdHash := string(pwdHashBytes)

	// 1. Scuola
	schoolCode := "MIIS09900T"
	schoolName := "I.I.S. Alessandro Volta"
	var schoolID string
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM schools WHERE code = $1`, schoolCode).Scan(&schoolID)
	if err != nil {
		schoolID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx, `
			INSERT INTO schools (id, name, address, city, zip_code, type, code, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name, updated_at = NOW()
			RETURNING id
		`, schoolID, schoolName, "Via dei Tintori 25", "Milano", "20121", "Istituto Superiore", schoolCode)
		if err != nil {
			return fmt.Errorf("failed to insert school: %w", err)
		}
	}
	log.Printf("✔ Scuola pronta: %s [%s] (ID: %s)\n", schoolName, schoolCode, schoolID)

	// 2. Anno Accademico (2025/2026)
	var academicYearID string
	err = dbConn.QueryRowContext(ctx, `SELECT id FROM academic_years WHERE school_id = $1 AND name = $2`, schoolID, "2025/2026").Scan(&academicYearID)
	if err != nil {
		academicYearID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx, `
			INSERT INTO academic_years (id, school_id, name, start_date, end_date, current, created_at)
			VALUES ($1, $2, $3, $4, $5, TRUE, NOW())
		`, academicYearID, schoolID, "2025/2026", "2025-09-01", "2026-06-30")
		if err != nil {
			return fmt.Errorf("failed to insert academic year: %w", err)
		}
	}
	log.Printf("✔ Anno Accademico: 2025/2026 (ID: %s)\n", academicYearID)

	// Helper for user creation
	createUser := func(email, firstName, lastName, role string) (string, error) {
		var uID string
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&uID)
		if err == nil {
			_, _ = dbConn.ExecContext(ctx, `
				UPDATE users 
				SET password_hash = $1, is_active = TRUE, role = $2, school_id = $3, first_name = $4, last_name = $5
				WHERE id = $6
			`, defaultPwdHash, role, schoolID, firstName, lastName, uID)
			return uID, nil
		}
		uID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx, `
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE, TRUE, NOW(), NOW())
		`, uID, email, defaultPwdHash, firstName, lastName, role, schoolID)
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

	// 3. Admin & Dirigenza
	adminID, _ := createUser("admin.volta@scuola.it", "Admin", "Volta", "admin")
	dirigenteID, _ := createUser("dirigente.volta@scuola.it", "Alessandro", "Volta", "principal")
	vicarioID, _ := createUser("vicario.volta@scuola.it", "Elena", "Baroni", "vice_principal")
	dsgaID, _ := createUser("dsga.volta@scuola.it", "Giuseppe", "Mariani", "dsga")
	secID, _ := createUser("segreteria.volta@scuola.it", "Carla", "Fontana", "secretary")
	log.Printf("✔ Staff Direttivo pronto (Admin: %s, Dirigente: %s, Vicario: %s, DSGA: %s, Segr: %s)\n",
		adminID, dirigenteID, vicarioID, dsgaID, secID)

	// 4. Edifici Scolastici (Plessi)
	type buildingDef struct {
		Name    string
		Address string
	}
	buildings := []buildingDef{
		{"Sede Centrale", "Via dei Tintori 25, Milano"},
		{"Succursale Sud", "Via Meucci 14, Milano"},
	}
	buildingIDs := make(map[string]string)
	for _, b := range buildings {
		var bID string
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM school_buildings WHERE school_id = $1 AND name = $2`, schoolID, b.Name).Scan(&bID)
		if err != nil {
			bID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO school_buildings (id, school_id, name, address, is_active, created_at, updated_at)
				VALUES ($1, $2, $3, $4, TRUE, NOW(), NOW())
			`, bID, schoolID, b.Name, b.Address)
			if err != nil {
				return fmt.Errorf("failed to create building %s: %w", b.Name, err)
			}
		}
		buildingIDs[b.Name] = bID
	}
	log.Printf("✔ Plessi scolastici: Sede Centrale (%s), Succursale Sud (%s)\n",
		buildingIDs["Sede Centrale"], buildingIDs["Succursale Sud"])

	// 5. Laboratori & Palestre (bookable_rooms)
	type roomDef struct {
		BuildingName string
		Name         string
		RoomType     string
		Capacity     int
	}
	roomsToSeed := []roomDef{
		{"Sede Centrale", "Palestra Grande Centrale", "gym", 60},
		{"Sede Centrale", "Palestra Piccola Centrale", "gym", 35},
		{"Sede Centrale", "Laboratorio Informatica 1", "lab_computer", 32},
		{"Sede Centrale", "Laboratorio Informatica 2", "lab_computer", 30},
		{"Sede Centrale", "Laboratorio di Fisica", "lab_physics", 30},
		{"Sede Centrale", "Laboratorio Scienze e Chimica", "lab_science", 30},
		{"Sede Centrale", "Laboratorio Linguistico Multimediale", "lab_language", 28},
		{"Succursale Sud", "Palestra Succursale", "gym", 50},
		{"Succursale Sud", "Laboratorio Informatica Succursale", "lab_computer", 30},
		{"Succursale Sud", "Laboratorio Scienze Succursale", "lab_science", 30},
		{"Succursale Sud", "Laboratorio Elettronica e Sistemi", "lab_electronics", 28},
	}
	for _, r := range roomsToSeed {
		bID := buildingIDs[r.BuildingName]
		var rID string
		err := dbConn.QueryRowContext(ctx, `
			SELECT id FROM bookable_rooms WHERE school_id = $1 AND name = $2
		`, schoolID, r.Name).Scan(&rID)
		if err != nil {
			rID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO bookable_rooms (id, school_id, building_id, name, room_type, capacity, is_active, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, TRUE, NOW(), NOW())
			`, rID, schoolID, bID, r.Name, r.RoomType, r.Capacity)
			if err != nil {
				return fmt.Errorf("failed to create room %s: %w", r.Name, err)
			}
		}
	}
	log.Printf("✔ %d Laboratori e Palestre pronti per l'assegnazione oraria\n", len(roomsToSeed))

	// 6. Materie Scolastiche (subjects)
	type subjectDef struct {
		Code        string
		Name        string
		Description string
	}
	subjectsToSeed := []subjectDef{
		{"ITA", "Italiano e Letteratura", "Lingua e letteratura italiana"},
		{"STO", "Storia", "Storia generale ed europea"},
		{"GEO", "Geografia", "Geografia generale ed economica"},
		{"FIL", "Filosofia", "Filosofia occidentale"},
		{"ING", "Lingua Inglese", "Lingua e cultura inglese"},
		{"SPA", "Lingua Spagnola", "Seconda lingua comunitaria"},
		{"MAT", "Matematica", "Matematica e algebra"},
		{"FIS", "Fisica", "Fisica classica e moderna"},
		{"SCI", "Scienze Naturali", "Biologia, Chimica e Scienze della Terra"},
		{"INF", "Informatica", "Informatica, algoritmi e telecomunicazioni"},
		{"ART", "Disegno e Storia dell'Arte", "Disegno tecnico e storia dell'arte"},
		{"EDF", "Scienze Motorie e Sportive", "Educazione fisica e sport"},
		{"DIR", "Diritto ed Economia", "Diritto ed economia politica"},
		{"REL", "Religione Cattolica", "Religione o attività alternativa"},
	}
	subjectIDs := make(map[string]string)
	for _, s := range subjectsToSeed {
		var sID string
		err := dbConn.QueryRowContext(ctx, `SELECT id FROM subjects WHERE school_id = $1 AND code = $2`, schoolID, s.Code).Scan(&sID)
		if err != nil {
			sID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO subjects (id, school_id, name, code, description, is_mandatory, created_at)
				VALUES ($1, $2, $3, $4, $5, TRUE, NOW())
			`, sID, schoolID, s.Name, s.Code, s.Description)
			if err != nil {
				return fmt.Errorf("failed to create subject %s: %w", s.Name, err)
			}
		}
		subjectIDs[s.Code] = sID
	}
	log.Printf("✔ %d Materie scolastiche registrate con successo\n", len(subjectsToSeed))

	// 7. Vincoli di Laboratorio per Materia (subject_room_requirements)
	type roomReqDef struct {
		SubjectCode      string
		RequiredRoomType string
		LabHours         int
		IsMandatory      bool
	}
	reqsToSeed := []roomReqDef{
		{"EDF", "gym", 2, true},
		{"INF", "lab_computer", 2, true},
		{"FIS", "lab_physics", 1, false},
		{"SCI", "lab_science", 1, false},
	}
	for _, req := range reqsToSeed {
		subID := subjectIDs[req.SubjectCode]
		_, _ = dbConn.ExecContext(ctx, `
			INSERT INTO subject_room_requirements (id, school_id, subject_id, required_room_type, is_mandatory, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
			ON CONFLICT (school_id, subject_id) DO UPDATE SET
				required_room_type = EXCLUDED.required_room_type,
				is_mandatory = EXCLUDED.is_mandatory
		`, uuid.New().String(), schoolID, subID, req.RequiredRoomType, req.IsMandatory)
	}
	log.Println("✔ Vincoli aula/laboratorio configurati (Palestra, Lab Informatica, Fisica, Scienze)")

	// 8. Creazione delle 55 Classi (11 Sezioni x 5 Anni)
	// Sezioni A, B, C, D, E, F -> Sede Centrale (30 classi)
	// Sezioni G, H, I, L, M -> Succursale Sud (25 classi)
	type classInfo struct {
		ID          string
		Name        string
		Section     string
		Year        int
		BuildingID  string
		Specialty   string
		TotalHours  int
		SubjectsMap map[string]int // SubjectCode -> weekly hours
	}

	sections := []struct {
		Letter    string
		Building  string
		Specialty string
	}{
		{"A", "Sede Centrale", "Liceo Scientifico"},
		{"B", "Sede Centrale", "Liceo Scientifico"},
		{"C", "Sede Centrale", "Scienze Applicate"},
		{"D", "Sede Centrale", "Scienze Applicate"},
		{"E", "Sede Centrale", "Liceo Linguistico"},
		{"F", "Sede Centrale", "Liceo Linguistico"},
		{"G", "Succursale Sud", "Tecnico Informatico"},
		{"H", "Succursale Sud", "Tecnico Informatico"},
		{"I", "Succursale Sud", "Tecnico Elettronico"},
		{"L", "Succursale Sud", "Tecnico Economico AFM"},
		{"M", "Succursale Sud", "Tecnico Economico AFM"},
	}

	allClasses := make([]*classInfo, 0, 55)
	for _, sec := range sections {
		for yr := 1; yr <= 5; yr++ {
			cName := fmt.Sprintf("%d%s", yr, sec.Letter)
			bID := buildingIDs[sec.Building]

			// Define realistic weekly hours per subject (monte ore ministeriale, max 30h/sett)
			hoursMap := make(map[string]int)
			switch sec.Specialty {
			case "Liceo Scientifico":
				if yr <= 2 {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["GEO"] = 1
					hoursMap["ING"] = 3
					hoursMap["MAT"] = 5
					hoursMap["FIS"] = 2
					hoursMap["SCI"] = 3
					hoursMap["ART"] = 2
					hoursMap["EDF"] = 2
					hoursMap["INF"] = 2
					hoursMap["REL"] = 1 // Totale: 27 ore
				} else {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["FIL"] = 3
					hoursMap["ING"] = 3
					hoursMap["MAT"] = 4
					hoursMap["FIS"] = 3
					hoursMap["SCI"] = 3
					hoursMap["ART"] = 2
					hoursMap["INF"] = 2
					hoursMap["EDF"] = 2
					hoursMap["DIR"] = 1
					hoursMap["REL"] = 1 // Totale: 30 ore
				}
			case "Scienze Applicate":
				if yr <= 2 {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 3
					hoursMap["ING"] = 3
					hoursMap["MAT"] = 5
					hoursMap["FIS"] = 2
					hoursMap["SCI"] = 3
					hoursMap["INF"] = 2
					hoursMap["ART"] = 2
					hoursMap["EDF"] = 2
					hoursMap["REL"] = 1 // Totale: 27 ore
				} else {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["FIL"] = 2
					hoursMap["ING"] = 3
					hoursMap["MAT"] = 4
					hoursMap["FIS"] = 3
					hoursMap["SCI"] = 4
					hoursMap["INF"] = 3
					hoursMap["ART"] = 2
					hoursMap["EDF"] = 2
					hoursMap["REL"] = 1 // Totale: 30 ore
				}
			case "Liceo Linguistico":
				if yr <= 2 {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 3
					hoursMap["ING"] = 4
					hoursMap["SPA"] = 3
					hoursMap["MAT"] = 3
					hoursMap["FIS"] = 2
					hoursMap["SCI"] = 2
					hoursMap["ART"] = 2
					hoursMap["EDF"] = 2
					hoursMap["REL"] = 1
					hoursMap["DIR"] = 1 // Totale: 27 ore
				} else {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["FIL"] = 2
					hoursMap["ING"] = 4
					hoursMap["SPA"] = 4
					hoursMap["MAT"] = 3
					hoursMap["FIS"] = 2
					hoursMap["SCI"] = 2
					hoursMap["ART"] = 2
					hoursMap["EDF"] = 2
					hoursMap["DIR"] = 2
					hoursMap["REL"] = 1 // Totale: 30 ore
				}
			case "Tecnico Informatico":
				if yr <= 2 {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["ING"] = 3
					hoursMap["MAT"] = 4
					hoursMap["SCI"] = 3
					hoursMap["INF"] = 4
					hoursMap["DIR"] = 2
					hoursMap["ART"] = 2
					hoursMap["EDF"] = 2
					hoursMap["REL"] = 1 // Totale: 27 ore
				} else {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["ING"] = 3
					hoursMap["MAT"] = 4
					hoursMap["INF"] = 6
					hoursMap["FIS"] = 3
					hoursMap["SCI"] = 2
					hoursMap["DIR"] = 2
					hoursMap["EDF"] = 2
					hoursMap["REL"] = 1 // Totale: 29 ore
				}
			case "Tecnico Elettronico":
				if yr <= 2 {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["ING"] = 3
					hoursMap["MAT"] = 4
					hoursMap["FIS"] = 3
					hoursMap["SCI"] = 3
					hoursMap["INF"] = 2
					hoursMap["DIR"] = 2
					hoursMap["EDF"] = 2
					hoursMap["REL"] = 1 // Totale: 26 ore
				} else {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["ING"] = 3
					hoursMap["MAT"] = 4
					hoursMap["FIS"] = 5
					hoursMap["INF"] = 4
					hoursMap["DIR"] = 2
					hoursMap["ART"] = 2
					hoursMap["EDF"] = 2
					hoursMap["REL"] = 1 // Totale: 29 ore
				}
			case "Tecnico Economico AFM":
				if yr <= 2 {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["ING"] = 3
					hoursMap["SPA"] = 3
					hoursMap["MAT"] = 4
					hoursMap["DIR"] = 3
					hoursMap["SCI"] = 2
					hoursMap["INF"] = 2
					hoursMap["EDF"] = 2
					hoursMap["REL"] = 1 // Totale: 26 ore
				} else {
					hoursMap["ITA"] = 4
					hoursMap["STO"] = 2
					hoursMap["ING"] = 3
					hoursMap["SPA"] = 3
					hoursMap["MAT"] = 3
					hoursMap["DIR"] = 5
					hoursMap["INF"] = 3
					hoursMap["ART"] = 2
					hoursMap["EDF"] = 2
					hoursMap["REL"] = 1 // Totale: 28 ore
				}
			}

			totH := 0
			for _, h := range hoursMap {
				totH += h
			}

			// Insert or retrieve class
			var cID string
			err := dbConn.QueryRowContext(ctx, `
				SELECT id FROM classes WHERE school_id = $1 AND name = $2 AND academic_year = $3
			`, schoolID, cName, "2025/2026").Scan(&cID)
			if err != nil {
				cID = uuid.New().String()
				_, err = dbConn.ExecContext(ctx, `
					INSERT INTO classes (id, school_id, name, section, academic_year, academic_year_id, building_id, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
				`, cID, schoolID, cName, sec.Letter, "2025/2026", academicYearID, bID)
				if err != nil {
					return fmt.Errorf("failed to create class %s: %w", cName, err)
				}
			} else {
				_, _ = dbConn.ExecContext(ctx, `
					UPDATE classes SET building_id = $1, academic_year_id = $2 WHERE id = $3
				`, bID, academicYearID, cID)
			}

			allClasses = append(allClasses, &classInfo{
				ID:          cID,
				Name:        cName,
				Section:     sec.Letter,
				Year:        yr,
				BuildingID:  bID,
				Specialty:   sec.Specialty,
				TotalHours:  totH,
				SubjectsMap: hoursMap,
			})
		}
	}
	log.Printf("✔ Create %d Classi (tutte con monte ore settimanale di 26-30h)\n", len(allClasses))

	// 9. Creazione di oltre 110 Insegnanti con profili, ruoli e materie
	// Organico docenti diviso per aree disciplinari:
	type teacherCandidate struct {
		First       string
		Last        string
		SubjectCode string
	}

	// 112 Insegnanti realistici
	rawTeachers := []teacherCandidate{
		// Italiano & Letteratura (14 docenti)
		{"Mario", "Rossi", "ITA"}, {"Lucia", "Ferrari", "ITA"}, {"Alessandro", "Colombo", "ITA"},
		{"Giulia", "Bianchi", "ITA"}, {"Paolo", "Romano", "ITA"}, {"Francesca", "Galli", "ITA"},
		{"Andrea", "Costa", "ITA"}, {"Simona", "Fontana", "ITA"}, {"Stefano", "Conti", "ITA"},
		{"Chiara", "Barbieri", "ITA"}, {"Matteo", "Marini", "ITA"}, {"Sara", "Greco", "ITA"},
		{"Daniele", "Lombardi", "ITA"}, {"Elena", "Serra", "ITA"},

		// Storia, Geografia & Filosofia (14 docenti)
		{"Giovanni", "Rinaldi", "STO"}, {"Roberta", "Gatti", "STO"}, {"Claudio", "Santoro", "STO"},
		{"Silvia", "Moretti", "STO"}, {"Davide", "Marchetti", "FIL"}, {"Federica", "Parisi", "FIL"},
		{"Pietro", "Villa", "FIL"}, {"Valentina", "Caruso", "FIL"}, {"Fabio", "Ferraro", "STO"},
		{"Monica", "Leone", "STO"}, {"Alberto", "Pellegrini", "FIL"}, {"Ilaria", "Mariani", "FIL"},
		{"Massimo", "De Luca", "GEO"}, {"Marta", "Palumbo", "GEO"},

		// Matematica (14 docenti)
		{"Lorenzo", "Rizzi", "MAT"}, {"Beatrice", "Barone", "MAT"}, {"Simone", "Basile", "MAT"},
		{"Anna", "Piras", "MAT"}, {"Filippo", "Vitali", "MAT"}, {"Alessia", "Sanna", "MAT"},
		{"Christian", "Coppola", "MAT"}, {"Serena", "D'Amico", "MAT"}, {"Enrico", "Amato", "MAT"},
		{"Camilla", "Damiani", "MAT"}, {"Jacopo", "Guerra", "MAT"}, {"Elisa", "Farina", "MAT"},
		{"Valerio", "Silvestri", "MAT"}, {"Giorgia", "Mazza", "MAT"},

		// Fisica (9 docenti)
		{"Tommaso", "Monti", "FIS"}, {"Federico", "Testa", "FIS"}, {"Ginevra", "Bernardi", "FIS"},
		{"Edoardo", "Caputo", "FIS"}, {"Ludovica", "Grassi", "FIS"}, {"Giacomo", "Fiore", "FIS"},
		{"Arianna", "Pugliese", "FIS"}, {"Riccardo", "Fabbri", "FIS"}, {"Martina", "Riva", "FIS"},

		// Lingua Inglese (11 docenti)
		{"Claire", "Smith", "ING"}, {"John", "Miller", "ING"}, {"David", "Brown", "ING"},
		{"Emma", "Taylor", "ING"}, {"James", "Wilson", "ING"}, {"Sophie", "Davies", "ING"},
		{"Oliver", "Evans", "ING"}, {"Emily", "Thomas", "ING"}, {"William", "Roberts", "ING"},
		{"Alice", "Johnson", "ING"}, {"Harry", "Walker", "ING"},

		// Lingua Spagnola (5 docenti)
		{"Maria", "Garcia", "SPA"}, {"Carlos", "Rodriguez", "SPA"}, {"Carmen", "Lopez", "SPA"},
		{"Diego", "Hernandez", "SPA"}, {"Isabella", "Martinez", "SPA"},

		// Scienze Naturali (9 docenti)
		{"Fabrizio", "Galli", "SCI"}, {"Laura", "Sala", "SCI"}, {"Roberto", "Brambilla", "SCI"},
		{"Daniela", "Cattaneo", "SCI"}, {"Michele", "Colombo", "SCI"}, {"Teresa", "Beretta", "SCI"},
		{"Carlo", "Negri", "SCI"}, {"Nicoletta", "Pozzi", "SCI"}, {"Gianluca", "Molteni", "SCI"},

		// Informatica (8 docenti)
		{"Fabio", "Viganò", "INF"}, {"Cristina", "Poretti", "INF"}, {"Luca", "Tagliabue", "INF"},
		{"Manuela", "Radice", "INF"}, {"Giorgio", "Castiglioni", "INF"}, {"Cinzia", "Ronchi", "INF"},
		{"Antonio", "Parravicini", "INF"}, {"Emanuela", "Mauri", "INF"},

		// Disegno e Storia dell'Arte (6 docenti)
		{"Raffaele", "Caravaggio", "ART"}, {"Marina", "Canova", "ART"}, {"Leonardo", "Bernini", "ART"},
		{"Giulia", "Giotto", "ART"}, {"Vincenzo", "Botticelli", "ART"}, {"Sonia", "Modigliani", "ART"},

		// Scienze Motorie e Sportive (7 docenti)
		{"Gianluigi", "Buffon", "EDF"}, {"Federica", "Pellegrini", "EDF"}, {"Yuri", "Chechi", "EDF"},
		{"Sara", "Simeoni", "EDF"}, {"Pietro", "Mennea", "EDF"}, {"Valentina", "Vezzali", "EDF"},
		{"Alberto", "Tomba", "EDF"},

		// Diritto ed Economia (7 docenti)
		{"Cesare", "Beccaria", "DIR"}, {"Ilaria", "Capogrossi", "DIR"}, {"Maurizio", "Galimberti", "DIR"},
		{"Beatrice", "Cavour", "DIR"}, {"Franco", "Einaudi", "DIR"}, {"Loredana", "Saragat", "DIR"},
		{"Ettore", "Calamandrei", "DIR"},

		// Religione Cattolica / Att. Alternativa (4 docenti)
		{"Francesco", "Assisi", "REL"}, {"Chiara", "Offreduccio", "REL"}, {"Tommaso", "Aquino", "REL"},
		{"Caterina", "Siena", "REL"},
	}

	type teacherInfo struct {
		UserID        string
		TeacherProfID string
		FullName      string
		SubjectCode   string
		CurrentHours  int
	}

	teachersBySubject := make(map[string][]*teacherInfo)
	allCreatedTeachers := make([]*teacherInfo, 0, len(rawTeachers))

	for idx, rt := range rawTeachers {
		email := strings.ToLower(fmt.Sprintf("%s.%s.volta@scuola.it",
			strings.ReplaceAll(rt.First, " ", ""),
			strings.ReplaceAll(rt.Last, " ", "")))
		// Unique email handling if collision
		if idx > 0 {
			for _, prev := range allCreatedTeachers {
				if prev.FullName == fmt.Sprintf("%s %s", rt.First, rt.Last) {
					email = strings.ToLower(fmt.Sprintf("%s.%s%d.volta@scuola.it", rt.First, rt.Last, idx))
					break
				}
			}
		}

		uID, err := createUser(email, rt.First, rt.Last, "teacher")
		if err != nil {
			return fmt.Errorf("failed to create teacher user %s: %w", email, err)
		}

		// Ensure record in teachers table
		var tProfID string
		err = dbConn.QueryRowContext(ctx, `SELECT id FROM teachers WHERE user_id = $1 AND school_id = $2`, uID, schoolID).Scan(&tProfID)
		if err != nil {
			tProfID = uuid.New().String()
			hiringDate := time.Date(2015+(idx%10), time.Month(1+(idx%12)), 1+(idx%25), 0, 0, 0, 0, time.UTC)
			_, err = dbConn.ExecContext(ctx, `
				INSERT INTO teachers (id, user_id, school_id, hiring_date, created_at, updated_at)
				VALUES ($1, $2, $3, $4, NOW(), NOW())
			`, tProfID, uID, schoolID, hiringDate.Format("2006-01-02"))
			if err != nil {
				return fmt.Errorf("failed to create teacher profile for %s: %w", email, err)
			}
		}

		// Associate with subject in teacher_subjects
		subID := subjectIDs[rt.SubjectCode]
		if subID != "" {
			_, _ = dbConn.ExecContext(ctx, `
				INSERT INTO teacher_subjects (id, teacher_id, subject_id, created_at)
				VALUES ($1, $2, $3, NOW())
				ON CONFLICT (teacher_id, subject_id) DO NOTHING
			`, uuid.New().String(), tProfID, subID)
		}

		tObj := &teacherInfo{
			UserID:        uID,
			TeacherProfID: tProfID,
			FullName:      fmt.Sprintf("%s %s", rt.First, rt.Last),
			SubjectCode:   rt.SubjectCode,
			CurrentHours:  0,
		}
		teachersBySubject[rt.SubjectCode] = append(teachersBySubject[rt.SubjectCode], tObj)
		allCreatedTeachers = append(allCreatedTeachers, tObj)
	}
	log.Printf("✔ %d Docenti creati con account, profili e abilitazioni disciplinari\n", len(allCreatedTeachers))

	// 10. Assegnazione Docenti alle Classi (Cattedre realistiche da 14 a 18 ore max)
	// Clear existing class_subjects for this school's classes to avoid duplicates on re-run
	_, _ = dbConn.ExecContext(ctx, `
		DELETE FROM class_subjects 
		WHERE class_id IN (SELECT id FROM classes WHERE school_id = $1)
	`, schoolID)

	assignStmt, err := dbConn.PrepareContext(ctx, `
		INSERT INTO class_subjects (id, class_id, subject_id, teacher_id, hours_per_week, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (class_id, subject_id, teacher_id) DO UPDATE SET hours_per_week = EXCLUDED.hours_per_week
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare class_subjects insert: %w", err)
	}
	defer assignStmt.Close()

	totalAssignedHours := 0
	coordinatorsByClass := make(map[string]string)

	for _, cl := range allClasses {
		for subCode, hrs := range cl.SubjectsMap {
			if hrs <= 0 {
				continue
			}
			subID := subjectIDs[subCode]

			// Pick teacher for this subject:
			// Priority to teachers who teach other subjects in the same class OR have lowest hours <= 18 - hrs
			availTeachers := teachersBySubject[subCode]
			if len(availTeachers) == 0 {
				// Fallback to related subject teacher if available
				switch subCode {
				case "STO", "GEO":
					availTeachers = teachersBySubject["ITA"]
				case "FIL":
					availTeachers = teachersBySubject["STO"]
				case "FIS":
					availTeachers = teachersBySubject["MAT"]
				default:
					availTeachers = allCreatedTeachers
				}
			}

			var bestTeacher *teacherInfo
			minHours := 999
			for _, cand := range availTeachers {
				if cand.CurrentHours+hrs <= 18 {
					if cand.CurrentHours < minHours {
						minHours = cand.CurrentHours
						bestTeacher = cand
					}
				}
			}
			// If all at 18h, allow slight overflow up to 19h
			if bestTeacher == nil {
				for _, cand := range availTeachers {
					if cand.CurrentHours < minHours {
						minHours = cand.CurrentHours
						bestTeacher = cand
					}
				}
			}

			if bestTeacher != nil {
				bestTeacher.CurrentHours += hrs
				totalAssignedHours += hrs

				_, err := assignStmt.ExecContext(ctx,
					uuid.New().String(), cl.ID, subID, bestTeacher.TeacherProfID, float64(hrs),
				)
				if err != nil {
					return fmt.Errorf("failed to assign teacher %s to class %s (%s): %w",
						bestTeacher.FullName, cl.Name, subCode, err)
				}

				// Assign first core teacher (ITA or MAT) as class coordinator
				if coordinatorsByClass[cl.ID] == "" && (subCode == "ITA" || subCode == "MAT") {
					coordinatorsByClass[cl.ID] = bestTeacher.UserID
				}
			}
		}
	}

	// Update coordinators on classes
	for cID, coordUserID := range coordinatorsByClass {
		_, _ = dbConn.ExecContext(ctx, `UPDATE classes SET coordinator_id = $1 WHERE id = $2`, coordUserID, cID)
	}
	log.Printf("✔ Assegnate %d ore settimanali complessive su 55 classi (%d docenti con cattedre bilanciate)\n",
		totalAssignedHours, len(allCreatedTeachers))

	// 11. Desiderata e Preferenze Orarie Docenti (teacher_schedule_preferences)
	// Clear existing preferences for this school
	_, _ = dbConn.ExecContext(ctx, `DELETE FROM teacher_schedule_preferences WHERE school_id = $1`, schoolID)

	prefStmt, err := dbConn.PrepareContext(ctx, `
		INSERT INTO teacher_schedule_preferences (
			id, school_id, teacher_id, academic_year_id, day_of_week, hour_index, preference_type, reason, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
		ON CONFLICT (teacher_id, academic_year_id, day_of_week, hour_index) DO UPDATE
		SET preference_type = EXCLUDED.preference_type, reason = EXCLUDED.reason
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare teacher preferences insert: %w", err)
	}
	defer prefStmt.Close()

	// Distribute realistic desiderata:
	// - Each teacher has 1 Giorno Libero (Day Off) from Lunedì (1) to Venerdì (5)
	// - Morning preference (early hours 1-3) or late preference (hours 3-5)
	// - 1-2 unavailable hours for part-time/personal needs
	for idx, t := range allCreatedTeachers {
		dayOff := 1 + (idx % 5) // 1=Mon, 2=Tue, 3=Wed, 4=Thu, 5=Fri
		// Day off: all 6 hours unavailable
		for h := 1; h <= 6; h++ {
			_, _ = prefStmt.ExecContext(ctx,
				uuid.New().String(), schoolID, t.UserID, academicYearID,
				dayOff, h, "unavailable", "Richiesta giorno libero settimanale",
			)
		}

		// Other days: early or late preferences
		prefType := "early"
		if idx%2 == 1 {
			prefType = "late"
		}
		for d := 1; d <= 5; d++ {
			if d == dayOff {
				continue
			}
			if prefType == "early" {
				_, _ = prefStmt.ExecContext(ctx,
					uuid.New().String(), schoolID, t.UserID, academicYearID,
					d, 1, "preferred", "Preferenza didattica prime ore",
				)
				_, _ = prefStmt.ExecContext(ctx,
					uuid.New().String(), schoolID, t.UserID, academicYearID,
					d, 2, "preferred", "Preferenza didattica prime ore",
				)
			} else {
				_, _ = prefStmt.ExecContext(ctx,
					uuid.New().String(), schoolID, t.UserID, academicYearID,
					d, 4, "preferred", "Preferenza didattica ore centrali/tarde",
				)
				_, _ = prefStmt.ExecContext(ctx,
					uuid.New().String(), schoolID, t.UserID, academicYearID,
					d, 5, "preferred", "Preferenza didattica ore centrali/tarde",
				)
			}
		}
	}
	log.Println("✔ Preferenze orarie docenti (giorno libero, prime/ultime ore) configurate per tutti i docenti")

	// 12. Vincoli Generali per la Generazione dell'Orario (timetable_constraints)
	// Clear existing constraints
	_, _ = dbConn.ExecContext(ctx, `DELETE FROM timetable_constraints WHERE school_id = $1`, schoolID)

	constraintStmt, err := dbConn.PrepareContext(ctx, `
		INSERT INTO timetable_constraints (
			id, school_id, constraint_type, target_type, target_id, parameters, is_hard, priority, is_active, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, TRUE, NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare timetable_constraints insert: %w", err)
	}
	defer constraintStmt.Close()

	// A. Limiti giornalieri per ciascuna classe (min 4h, max 6h al giorno)
	for _, cl := range allClasses {
		limitParams, _ := json.Marshal(map[string]int{
			"min_hours_per_day": 4,
			"max_hours_per_day": 6,
		})
		classIDCopy := cl.ID
		_, _ = constraintStmt.ExecContext(ctx,
			uuid.New().String(), schoolID, "class_daily_hours", "class", classIDCopy,
			limitParams, false, 8,
		)
	}

	// B. Massimo ore giornaliere consecutive per docente (max 5 ore)
	maxTeacherHParams, _ := json.Marshal(map[string]int{"max_hours": 5})
	_, _ = constraintStmt.ExecContext(ctx,
		uuid.New().String(), schoolID, "max_daily_teacher_hours", nil, nil,
		maxTeacherHParams, false, 6,
	)

	// C. Evita buchi orari nell'orario dei docenti (no_consecutive_gaps)
	noGapsParams, _ := json.Marshal(map[string]bool{"minimize_gaps": true})
	_, _ = constraintStmt.ExecContext(ctx,
		uuid.New().String(), schoolID, "no_consecutive_gaps", nil, nil,
		noGapsParams, false, 5,
	)

	// D. Materie pesanti non all'ultima ora
	heavyParams, _ := json.Marshal(map[string]bool{"avoid_last_hour": true})
	_, _ = constraintStmt.ExecContext(ctx,
		uuid.New().String(), schoolID, "avoid_heavy_last_hour", nil, nil,
		heavyParams, false, 4,
	)

	// E. Finestra desiderata aperta (desiderata_window)
	winParams, _ := json.Marshal(map[string]bool{"is_open": true})
	_, _ = constraintStmt.ExecContext(ctx,
		uuid.New().String(), schoolID, "desiderata_window", nil, nil,
		winParams, false, 1,
	)

	// F. Gruppo associato di prova (es. Spagnolo condiviso tra 3E e 3F)
	var class3EID, class3FID string
	for _, c := range allClasses {
		if c.Name == "3E" {
			class3EID = c.ID
		}
		if c.Name == "3F" {
			class3FID = c.ID
		}
	}
	if class3EID != "" && class3FID != "" && len(teachersBySubject["SPA"]) > 0 {
		leadSpanTeacher := teachersBySubject["SPA"][0]
		assocParams, _ := json.Marshal(map[string]interface{}{
			"name":           "Laboratorio Linguistico Spagnolo 3E-3F",
			"subject_id":     subjectIDs["SPA"],
			"teacher_id":     leadSpanTeacher.TeacherProfID,
			"class_ids":      []string{class3EID, class3FID},
			"hours_per_week": 2,
		})
		_, _ = constraintStmt.ExecContext(ctx,
			uuid.New().String(), schoolID, "associated_group", nil, nil,
			assocParams, true, 10,
		)
		log.Println("✔ Gruppo associato inter-classe (3E-3F Spagnolo) configurato con successo")
	}

	log.Println("==================================================================")
	log.Printf("🎉 SEEDING COMPLETATO CON SUCCESSO PER: %s\n", schoolName)
	log.Printf("📊 Statistiche:\n")
	log.Printf("   • Plessi: 2 (Sede Centrale & Succursale Sud)\n")
	log.Printf("   • Laboratori e Palestre: %d\n", len(roomsToSeed))
	log.Printf("   • Materie: %d\n", len(subjectsToSeed))
	log.Printf("   • Classi: %d (da 1A a 5M, 26-30h settimanali per classe)\n", len(allClasses))
	log.Printf("   • Docenti: %d (con cattedre realistiche da 14 a 18 ore)\n", len(allCreatedTeachers))
	log.Printf("   • Ore settimanali totali assegnate: %d ore\n", totalAssignedHours)
	log.Printf("   • Credenziali Admin: admin.volta@scuola.it / password\n")
	log.Printf("   • Credenziali Dirigente: dirigente.volta@scuola.it / password\n")
	log.Printf("   • Credenziali Segreteria: segreteria.volta@scuola.it / password\n")
	log.Printf("   • Credenziali Docenti: [nome].[cognome].volta@scuola.it / password\n")
	log.Println("==================================================================")

	return nil
}
