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

// SeedOrarioTest popola un dataset completo per testare la generazione orario scolastico.
// Scuola: "I.I.S. Giacomo Leopardi" (Roma, RMIS08800G)
// - 3 Plessi (Sede Centrale, Succursale Nord, Succursale Est)
// - 14 Laboratori e Palestre prenotabili con tipi aula
// - 14 Materie con vincoli di aula (EDF->gym, INF->lab_computer, FIS->lab_physics, SCI->lab_science, ART->lab_art)
// - 55 Classi (1A-5M, 11 sezioni x 5 anni) con monte ore ministeriale 26-30h/sett
// - 110 Docenti con abilitazioni e cattedre bilanciate (14-18h/sett)
// - Preferenze orarie realistiche per tutti i docenti (giorno libero, prime/ultime ore, 20% con spot)
// - Vincoli di generazione orario: limiti classe, ore consecutive, laboratori, distribuzione
// - 6 Gruppi associati inter-classe (Spagnolo E-F anni 3-5, INF G-H anni 3-5)
// Credenziali: admin.leopardi@scuola.it / password
// Docenti: [nome].[cognome].leopardi@scuola.it / password
func SeedOrarioTest(ctx context.Context, dbConn *sql.DB) error {
	log.Println("==================================================================")
	log.Println("SEEDING: I.I.S. Giacomo Leopardi (Test Generazione Orario)")
	log.Println("==================================================================")

	pwdHashBytes, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	defaultPwdHash := string(pwdHashBytes)

	schoolCode := "RMIS08800G"
	schoolName := "I.I.S. Giacomo Leopardi"
	var schoolID string
	qSchool := "SELECT id FROM schools WHERE code = $1"
	err = dbConn.QueryRowContext(ctx, qSchool, schoolCode).Scan(&schoolID)
	if err != nil {
		schoolID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx,
			"INSERT INTO schools (id, name, address, city, zip_code, type, code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name, updated_at = NOW()",
			schoolID, schoolName, "Viale dei Romantici 15", "Roma", "00196", "Istituto Superiore", schoolCode)
		if err != nil {
			return fmt.Errorf("failed to insert school: %w", err)
		}
		_ = dbConn.QueryRowContext(ctx, qSchool, schoolCode).Scan(&schoolID)
	}
	log.Printf("Scuola: %s [%s] ID=%s", schoolName, schoolCode, schoolID)

	var academicYearID string
	err = dbConn.QueryRowContext(ctx, "SELECT id FROM academic_years WHERE school_id = $1 AND name = $2", schoolID, "2025/2026").Scan(&academicYearID)
	if err != nil {
		academicYearID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx,
			"INSERT INTO academic_years (id, school_id, name, start_date, end_date, current, created_at) VALUES ($1, $2, $3, $4, $5, TRUE, NOW())",
			academicYearID, schoolID, "2025/2026", "2025-09-01", "2026-06-30")
		if err != nil {
			return fmt.Errorf("failed to insert academic year: %w", err)
		}
	}
	log.Printf("Anno Accademico: 2025/2026 ID=%s", academicYearID)

	createUser := func(email, firstName, lastName, role string) (string, error) {
		var uID string
		err := dbConn.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", email).Scan(&uID)
		if err == nil {
			_, _ = dbConn.ExecContext(ctx,
				"UPDATE users SET password_hash = $1, is_active = TRUE, role = $2, school_id = $3, first_name = $4, last_name = $5 WHERE id = $6",
				defaultPwdHash, role, schoolID, firstName, lastName, uID)
			return uID, nil
		}
		uID = uuid.New().String()
		_, err = dbConn.ExecContext(ctx,
			"INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE, TRUE, NOW(), NOW())",
			uID, email, defaultPwdHash, firstName, lastName, role, schoolID)
		if err != nil {
			return "", err
		}
		_, _ = dbConn.ExecContext(ctx,
			"INSERT INTO user_roles (id, user_id, role, school_id, created_at) VALUES ($1, $2, $3, $4, NOW()) ON CONFLICT (user_id, role, school_id) DO NOTHING",
			uuid.New().String(), uID, role, schoolID)
		return uID, nil
	}

	adminID, _ := createUser("admin.leopardi@scuola.it", "Admin", "Leopardi", "admin")
	dirigenteID, _ := createUser("dirigente.leopardi@scuola.it", "Giacomo", "Leopardi", "principal")
	vicarioID, _ := createUser("vicario.leopardi@scuola.it", "Teresa", "Verzino", "vice_principal")
	dsgaID, _ := createUser("dsga.leopardi@scuola.it", "Carmelo", "Amoruso", "dsga")
	secID, _ := createUser("segreteria.leopardi@scuola.it", "Nunzia", "Palermo", "secretary")
	log.Printf("Staff: admin=%s dirigente=%s vicario=%s dsga=%s sec=%s",
		adminID, dirigenteID, vicarioID, dsgaID, secID)

	buildingIDs := make(map[string]string)
	for _, bDef := range []struct{ Name, Addr string }{
		{"Sede Centrale", "Viale dei Romantici 15, Roma"},
		{"Succursale Nord", "Via Leopardi 8, Roma"},
		{"Succursale Est", "Via Recanati 22, Roma"},
	} {
		var bID string
		err := dbConn.QueryRowContext(ctx, "SELECT id FROM school_buildings WHERE school_id = $1 AND name = $2", schoolID, bDef.Name).Scan(&bID)
		if err != nil {
			bID = uuid.New().String()
			_, _ = dbConn.ExecContext(ctx,
				"INSERT INTO school_buildings (id, school_id, name, address, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, TRUE, NOW(), NOW())",
				bID, schoolID, bDef.Name, bDef.Addr)
		}
		buildingIDs[bDef.Name] = bID
	}
	log.Println("3 Plessi: Sede Centrale, Succursale Nord, Succursale Est")

	roomCount := 0
	for _, rDef := range []struct {
		Building, Name, RoomType string
		Capacity                 int
	}{
		{"Sede Centrale", "Palestra Grande Centrale", "gym", 60},
		{"Sede Centrale", "Palestra Piccola Centrale", "gym", 35},
		{"Sede Centrale", "Laboratorio Informatica 1", "lab_computer", 32},
		{"Sede Centrale", "Laboratorio Informatica 2", "lab_computer", 30},
		{"Sede Centrale", "Laboratorio di Fisica", "lab_physics", 28},
		{"Sede Centrale", "Laboratorio Scienze e Chimica", "lab_science", 28},
		{"Sede Centrale", "Laboratorio Linguistico Multimediale", "lab_language", 26},
		{"Sede Centrale", "Aula di Arte e Disegno", "lab_art", 28},
		{"Succursale Nord", "Palestra Succursale Nord", "gym", 48},
		{"Succursale Nord", "Laboratorio Informatica Nord", "lab_computer", 28},
		{"Succursale Nord", "Laboratorio Scienze Nord", "lab_science", 26},
		{"Succursale Est", "Palestra Succursale Est", "gym", 45},
		{"Succursale Est", "Laboratorio Informatica Est", "lab_computer", 30},
		{"Succursale Est", "Laboratorio Elettronica e Sistemi", "lab_electronics", 26},
	} {
		bID := buildingIDs[rDef.Building]
		var rID string
		err := dbConn.QueryRowContext(ctx, "SELECT id FROM bookable_rooms WHERE school_id = $1 AND name = $2", schoolID, rDef.Name).Scan(&rID)
		if err != nil {
			rID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx,
				"INSERT INTO bookable_rooms (id, school_id, building_id, name, room_type, capacity, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, TRUE, NOW(), NOW())",
				rID, schoolID, bID, rDef.Name, rDef.RoomType, rDef.Capacity)
			if err != nil {
				return fmt.Errorf("failed to create room %s: %w", rDef.Name, err)
			}
		}
		roomCount++
	}
	log.Printf("%d Laboratori e Palestre pronti", roomCount)

	subjectIDs := make(map[string]string)
	for _, sDef := range []struct{ Code, Name, Desc string }{
		{"ITA", "Italiano e Letteratura", "Lingua e letteratura italiana"},
		{"STO", "Storia", "Storia generale ed europea"},
		{"GEO", "Geografia", "Geografia generale ed economica"},
		{"FIL", "Filosofia", "Filosofia occidentale e contemporanea"},
		{"ING", "Lingua Inglese", "Lingua e cultura inglese"},
		{"SPA", "Lingua Spagnola", "Seconda lingua comunitaria spagnolo"},
		{"MAT", "Matematica", "Matematica algebra e analisi"},
		{"FIS", "Fisica", "Fisica classica moderna e laboratorio"},
		{"SCI", "Scienze Naturali", "Biologia Chimica Scienze della Terra"},
		{"INF", "Informatica", "Informatica algoritmi e reti"},
		{"ART", "Disegno e Storia Arte", "Disegno tecnico e storia arte"},
		{"EDF", "Scienze Motorie", "Educazione fisica e attivita motorie"},
		{"DIR", "Diritto ed Economia", "Diritto ed economia politica"},
		{"REL", "Religione Cattolica", "Religione cattolica o att alternativa"},
	} {
		var sID string
		err := dbConn.QueryRowContext(ctx, "SELECT id FROM subjects WHERE school_id = $1 AND code = $2", schoolID, sDef.Code).Scan(&sID)
		if err != nil {
			sID = uuid.New().String()
			_, err = dbConn.ExecContext(ctx,
				"INSERT INTO subjects (id, school_id, name, code, description, is_mandatory, created_at) VALUES ($1, $2, $3, $4, $5, TRUE, NOW())",
				sID, schoolID, sDef.Name, sDef.Code, sDef.Desc)
			if err != nil {
				return fmt.Errorf("failed to create subject %s: %w", sDef.Name, err)
			}
		}
		subjectIDs[sDef.Code] = sID
	}
	log.Println("14 Materie registrate")

	for _, req := range []struct {
		Code, RoomType string
		Hard           bool
	}{
		{"EDF", "gym", true},
		{"INF", "lab_computer", true},
		{"FIS", "lab_physics", false},
		{"SCI", "lab_science", false},
		{"ART", "lab_art", false},
	} {
		_, _ = dbConn.ExecContext(ctx,
			"INSERT INTO subject_room_requirements (id, school_id, subject_id, required_room_type, is_mandatory, created_at) VALUES ($1, $2, $3, $4, $5, NOW()) ON CONFLICT (school_id, subject_id) DO UPDATE SET required_room_type = EXCLUDED.required_room_type, is_mandatory = EXCLUDED.is_mandatory",
			uuid.New().String(), schoolID, subjectIDs[req.Code], req.RoomType, req.Hard)
	}
	log.Println("Vincoli aula/lab configurati")

	type classInfo struct {
		ID, Name, Section string
		Year              int
		BuildingID        string
		Specialty         string
		TotalHours        int
		SubjectsMap       map[string]int
	}

	_, _ = dbConn.ExecContext(ctx,
		"DELETE FROM class_subjects WHERE class_id IN (SELECT id FROM classes WHERE school_id = $1)", schoolID)

	allClasses := make([]*classInfo, 0, 55)
	for _, sec := range []struct {
		Letter, Building, Specialty string
	}{
		{"A", "Sede Centrale", "Liceo Scientifico"},
		{"B", "Sede Centrale", "Liceo Scientifico"},
		{"C", "Sede Centrale", "Scienze Applicate"},
		{"D", "Sede Centrale", "Scienze Applicate"},
		{"E", "Sede Centrale", "Liceo Linguistico"},
		{"F", "Sede Centrale", "Liceo Linguistico"},
		{"G", "Succursale Nord", "Tecnico Informatico"},
		{"H", "Succursale Nord", "Tecnico Informatico"},
		{"I", "Succursale Nord", "Tecnico Elettronico"},
		{"L", "Succursale Est", "Tecnico Economico AFM"},
		{"M", "Succursale Est", "Tecnico Economico AFM"},
	} {
		for yr := 1; yr <= 5; yr++ {
			cName := fmt.Sprintf("%d%s", yr, sec.Letter)
			bID := buildingIDs[sec.Building]
			hm := make(map[string]int)
			switch sec.Specialty {
			case "Liceo Scientifico":
				if yr <= 2 {
					hm["ITA"] = 4; hm["STO"] = 2; hm["GEO"] = 1; hm["ING"] = 3; hm["MAT"] = 5
					hm["FIS"] = 2; hm["SCI"] = 3; hm["ART"] = 2; hm["EDF"] = 2; hm["INF"] = 2; hm["REL"] = 1
				} else {
					hm["ITA"] = 4; hm["STO"] = 2; hm["FIL"] = 3; hm["ING"] = 3; hm["MAT"] = 4
					hm["FIS"] = 3; hm["SCI"] = 3; hm["ART"] = 2; hm["INF"] = 2; hm["EDF"] = 2; hm["DIR"] = 1; hm["REL"] = 1
				}
			case "Scienze Applicate":
				if yr <= 2 {
					hm["ITA"] = 4; hm["STO"] = 3; hm["ING"] = 3; hm["MAT"] = 5; hm["FIS"] = 2
					hm["SCI"] = 3; hm["INF"] = 2; hm["ART"] = 2; hm["EDF"] = 2; hm["REL"] = 1
				} else {
					hm["ITA"] = 4; hm["STO"] = 2; hm["FIL"] = 2; hm["ING"] = 3; hm["MAT"] = 4
					hm["FIS"] = 3; hm["SCI"] = 4; hm["INF"] = 3; hm["ART"] = 2; hm["EDF"] = 2; hm["REL"] = 1
				}
			case "Liceo Linguistico":
				if yr <= 2 {
					hm["ITA"] = 4; hm["STO"] = 3; hm["ING"] = 4; hm["SPA"] = 3; hm["MAT"] = 3
					hm["FIS"] = 2; hm["SCI"] = 2; hm["ART"] = 2; hm["EDF"] = 2; hm["REL"] = 1; hm["DIR"] = 1
				} else {
					hm["ITA"] = 4; hm["STO"] = 2; hm["FIL"] = 2; hm["ING"] = 4; hm["SPA"] = 4
					hm["MAT"] = 3; hm["FIS"] = 2; hm["SCI"] = 2; hm["ART"] = 2; hm["EDF"] = 2; hm["DIR"] = 2; hm["REL"] = 1
				}
			case "Tecnico Informatico":
				if yr <= 2 {
					hm["ITA"] = 4; hm["STO"] = 2; hm["ING"] = 3; hm["MAT"] = 4; hm["SCI"] = 3
					hm["INF"] = 4; hm["DIR"] = 2; hm["ART"] = 2; hm["EDF"] = 2; hm["REL"] = 1
				} else {
					hm["ITA"] = 4; hm["STO"] = 2; hm["ING"] = 3; hm["MAT"] = 4; hm["INF"] = 6
					hm["FIS"] = 3; hm["SCI"] = 2; hm["DIR"] = 2; hm["EDF"] = 2; hm["REL"] = 1
				}
			case "Tecnico Elettronico":
				if yr <= 2 {
					hm["ITA"] = 4; hm["STO"] = 2; hm["ING"] = 3; hm["MAT"] = 4; hm["FIS"] = 3
					hm["SCI"] = 3; hm["INF"] = 2; hm["DIR"] = 2; hm["EDF"] = 2; hm["REL"] = 1
				} else {
					hm["ITA"] = 4; hm["STO"] = 2; hm["ING"] = 3; hm["MAT"] = 4; hm["FIS"] = 5
					hm["INF"] = 4; hm["DIR"] = 2; hm["ART"] = 2; hm["EDF"] = 2; hm["REL"] = 1
				}
			case "Tecnico Economico AFM":
				if yr <= 2 {
					hm["ITA"] = 4; hm["STO"] = 2; hm["ING"] = 3; hm["SPA"] = 3; hm["MAT"] = 4
					hm["DIR"] = 3; hm["SCI"] = 2; hm["INF"] = 2; hm["EDF"] = 2; hm["REL"] = 1
				} else {
					hm["ITA"] = 4; hm["STO"] = 2; hm["ING"] = 3; hm["SPA"] = 3; hm["MAT"] = 3
					hm["DIR"] = 5; hm["INF"] = 3; hm["ART"] = 2; hm["EDF"] = 2; hm["REL"] = 1
				}
			}
			totH := 0
			for _, h := range hm {
				totH += h
			}
			var cID string
			err := dbConn.QueryRowContext(ctx,
				"SELECT id FROM classes WHERE school_id = $1 AND name = $2 AND academic_year = $3",
				schoolID, cName, "2025/2026").Scan(&cID)
			if err != nil {
				cID = uuid.New().String()
				_, err = dbConn.ExecContext(ctx,
					"INSERT INTO classes (id, school_id, name, section, academic_year, academic_year_id, building_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())",
					cID, schoolID, cName, sec.Letter, "2025/2026", academicYearID, bID)
				if err != nil {
					return fmt.Errorf("failed to create class %s: %w", cName, err)
				}
			} else {
				_, _ = dbConn.ExecContext(ctx,
					"UPDATE classes SET building_id = $1, academic_year_id = $2 WHERE id = $3",
					bID, academicYearID, cID)
			}
			allClasses = append(allClasses, &classInfo{
				ID: cID, Name: cName, Section: sec.Letter, Year: yr,
				BuildingID: bID, Specialty: sec.Specialty, TotalHours: totH, SubjectsMap: hm,
			})
		}
	}
	log.Printf("%d Classi create (26-30h/sett)", len(allClasses))

	type teacherCandidate struct{ First, Last, SubjectCode string }
	rawTeachers := []teacherCandidate{
		// Italiano (13)
		{"Giovanna", "Conti", "ITA"}, {"Marco", "Gentile", "ITA"}, {"Elisa", "Franzese", "ITA"},
		{"Roberto", "Amato", "ITA"}, {"Paola", "DeRosa", "ITA"}, {"Luca", "Fiore", "ITA"},
		{"Marta", "Longo", "ITA"}, {"Silvio", "Palumbo", "ITA"}, {"Cristina", "Ferretti", "ITA"},
		{"Aldo", "Giordano", "ITA"}, {"Nunzia", "Bassi", "ITA"}, {"Giorgio", "Mancuso", "ITA"},
		{"Beatrice", "EspositoIta", "ITA"},
		// Storia & Filosofia & Geografia (12)
		{"Fabrizio", "Rosati", "STO"}, {"Patrizia", "Martini", "STO"}, {"Giovanni", "Russo", "STO"},
		{"Daniela", "Fiorentino", "STO"}, {"Claudio", "Messina", "FIL"}, {"Silvia", "Agostini", "FIL"},
		{"Pietro", "Carbone", "FIL"}, {"Laura", "Ricci", "FIL"}, {"Eugenio", "Pellegrino", "STO"},
		{"Marianna", "Gatto", "FIL"}, {"Saverio", "Coppola", "GEO"}, {"Antonella", "Napoli", "GEO"},
		// Matematica (14)
		{"Simone", "Barbato", "MAT"}, {"Federica", "Mazza", "MAT"}, {"Enrico", "Pisano", "MAT"},
		{"Teresa", "Tozzi", "MAT"}, {"Giacomo", "Ferraro", "MAT"}, {"Valentina", "Sica", "MAT"},
		{"Matteo", "Cammarano", "MAT"}, {"Annarita", "Fusco", "MAT"}, {"Leonardo", "Vitiello", "MAT"},
		{"Carmela", "EspositoMat", "MAT"}, {"Raffaele", "Cinque", "MAT"}, {"Serena", "Guarino", "MAT"},
		{"Emanuele", "Scognamiglio", "MAT"}, {"Monica", "Iacono", "MAT"},
		// Fisica (9)
		{"Tommaso", "Albano", "FIS"}, {"Ginevra", "Romano", "FIS"}, {"Federico", "Marrone", "FIS"},
		{"Alessia", "Viti", "FIS"}, {"Edoardo", "Rizzo", "FIS"}, {"Ilaria", "Montagna", "FIS"},
		{"Armando", "Petrone", "FIS"}, {"Lucia", "Caiazzo", "FIS"}, {"Vincenzo", "Fabbri", "FIS"},
		// Inglese (11)
		{"Catherine", "Moore", "ING"}, {"James", "Harrison", "ING"}, {"Susan", "Clarke", "ING"},
		{"Robert", "Fletcher", "ING"}, {"Emma", "Whitfield", "ING"}, {"Michael", "Lawson", "ING"},
		{"Charlotte", "Spencer", "ING"}, {"William", "Owens", "ING"}, {"Hannah", "Griffiths", "ING"},
		{"Thomas", "Barker", "ING"}, {"Sophie", "Drummond", "ING"},
		// Spagnolo (5)
		{"Ana", "Vega", "SPA"}, {"Luis", "Morales", "SPA"}, {"Carmen", "Blanco", "SPA"},
		{"Jorge", "Reyes", "SPA"}, {"Pilar", "Castillo", "SPA"},
		// Scienze (9)
		{"Carmine", "FerrettiSci", "SCI"}, {"Rossana", "DAngelo", "SCI"}, {"Bruno", "Lombardi", "SCI"},
		{"Luisa", "Todaro", "SCI"}, {"Pasquale", "Gargiulo", "SCI"}, {"Immacolata", "Sorrentino", "SCI"},
		{"Antonio", "Caserta", "SCI"}, {"Francesca", "Vitale", "SCI"}, {"Nicola", "Amendola", "SCI"},
		// Informatica (9)
		{"Damiano", "Borrelli", "INF"}, {"Eleonora", "Trotta", "INF"}, {"Gianluca", "Capasso", "INF"},
		{"Roberta", "Marigliano", "INF"}, {"Marco", "Abagnale", "INF"}, {"Sara", "Castellano", "INF"},
		{"Ciro", "Imperato", "INF"}, {"Valeria", "Piccolo", "INF"}, {"Salvatore", "Carillo", "INF"},
		// Arte (6)
		{"Renato", "Ferragamo", "ART"}, {"Miriam", "Velasquez", "ART"}, {"Angelo", "Brunelleschi", "ART"},
		{"Giulia", "Caravaggio", "ART"}, {"Francesco", "Borromini", "ART"}, {"Concetta", "Raffaello", "ART"},
		// Scienze Motorie (7)
		{"Alessandro", "Totti", "EDF"}, {"Valentina", "Comaneci", "EDF"}, {"Marco", "Pantani", "EDF"},
		{"Raffaella", "Carra", "EDF"}, {"Massimo", "Stano", "EDF"}, {"Fiorella", "Izzo", "EDF"},
		{"Daniele", "Biagioni", "EDF"},
		// Diritto (8)
		{"Alberto", "Garibaldi", "DIR"}, {"Rosaria", "Monti", "DIR"}, {"Filippo", "Cristaudo", "DIR"},
		{"Angela", "Mazzini", "DIR"}, {"Gianni", "Cavour", "DIR"}, {"Loredana", "Depretis", "DIR"},
		{"Cesare", "Croce", "DIR"}, {"Elena", "Giolitti", "DIR"},
		// Religione (4)
		{"Eugenio", "Sacchi", "REL"}, {"Paola", "Orsini", "REL"}, {"Lorenzo", "Bellini", "REL"},
		{"Maria", "Colomba", "REL"},
	}

	type teacherInfo struct {
		UserID, TeacherProfID, FullName, SubjectCode string
		CurrentHours                                 int
	}
	teachersBySubject := make(map[string][]*teacherInfo)
	allCreatedTeachers := make([]*teacherInfo, 0, len(rawTeachers))
	usedEmails := make(map[string]bool)

	for idx, rt := range rawTeachers {
		baseEmail := strings.ToLower(fmt.Sprintf("%s.%s.leopardi@scuola.it", rt.First, rt.Last))
		email := baseEmail
		if usedEmails[email] {
			email = strings.ToLower(fmt.Sprintf("%s.%s%d.leopardi@scuola.it", rt.First, rt.Last, idx))
		}
		usedEmails[email] = true
		uID, err := createUser(email, rt.First, rt.Last, "teacher")
		if err != nil {
			return fmt.Errorf("failed to create teacher %s: %w", email, err)
		}
		var tProfID string
		err = dbConn.QueryRowContext(ctx,
			"SELECT id FROM teachers WHERE user_id = $1 AND school_id = $2", uID, schoolID).Scan(&tProfID)
		if err != nil {
			tProfID = uuid.New().String()
			hd := time.Date(2010+(idx%14), time.Month(1+(idx%12)), 1+(idx%25), 0, 0, 0, 0, time.UTC)
			_, err = dbConn.ExecContext(ctx,
				"INSERT INTO teachers (id, user_id, school_id, hiring_date, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())",
				tProfID, uID, schoolID, hd.Format("2006-01-02"))
			if err != nil {
				return fmt.Errorf("failed to create teacher profile %s: %w", email, err)
			}
		}
		subID := subjectIDs[rt.SubjectCode]
		if subID != "" {
			_, _ = dbConn.ExecContext(ctx,
				"INSERT INTO teacher_subjects (id, teacher_id, subject_id, created_at) VALUES ($1, $2, $3, NOW()) ON CONFLICT (teacher_id, subject_id) DO NOTHING",
				uuid.New().String(), tProfID, subID)
		}
		tObj := &teacherInfo{
			UserID:       uID,
			TeacherProfID: tProfID,
			FullName:     fmt.Sprintf("%s %s", rt.First, rt.Last),
			SubjectCode:  rt.SubjectCode,
		}
		teachersBySubject[rt.SubjectCode] = append(teachersBySubject[rt.SubjectCode], tObj)
		allCreatedTeachers = append(allCreatedTeachers, tObj)
	}
	log.Printf("%d Docenti creati con account, profili e abilitazioni", len(allCreatedTeachers))

	assignStmt, err := dbConn.PrepareContext(ctx,
		"INSERT INTO class_subjects (id, class_id, subject_id, teacher_id, hours_per_week, created_at) VALUES ($1, $2, $3, $4, $5, NOW()) ON CONFLICT (class_id, subject_id, teacher_id) DO UPDATE SET hours_per_week = EXCLUDED.hours_per_week")
	if err != nil {
		return fmt.Errorf("failed to prepare class_subjects: %w", err)
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
			avail := teachersBySubject[subCode]
			if len(avail) == 0 {
				switch subCode {
				case "STO", "GEO":
					avail = teachersBySubject["ITA"]
				case "FIL":
					avail = teachersBySubject["STO"]
				case "FIS":
					avail = teachersBySubject["MAT"]
				default:
					avail = allCreatedTeachers
				}
			}
			var best *teacherInfo
			minH := 999
			for _, cand := range avail {
				if cand.CurrentHours+hrs <= 18 && cand.CurrentHours < minH {
					minH = cand.CurrentHours
					best = cand
				}
			}
			if best == nil {
				minH = 999
				for _, cand := range avail {
					if cand.CurrentHours < minH {
						minH = cand.CurrentHours
						best = cand
					}
				}
			}
			if best != nil {
				best.CurrentHours += hrs
				totalAssignedHours += hrs
				_, err := assignStmt.ExecContext(ctx,
					uuid.New().String(), cl.ID, subID, best.TeacherProfID, float64(hrs))
				if err != nil {
					return fmt.Errorf("assign teacher %s to %s (%s): %w", best.FullName, cl.Name, subCode, err)
				}
				if coordinatorsByClass[cl.ID] == "" && (subCode == "ITA" || subCode == "MAT") {
					coordinatorsByClass[cl.ID] = best.UserID
				}
			}
		}
	}
	for cID, coordUID := range coordinatorsByClass {
		_, _ = dbConn.ExecContext(ctx,
			"UPDATE classes SET coordinator_id = $1 WHERE id = $2", coordUID, cID)
	}
	log.Printf("Assegnate %d ore/sett su %d classi (%d docenti)", totalAssignedHours, len(allClasses), len(allCreatedTeachers))

	_, _ = dbConn.ExecContext(ctx,
		"DELETE FROM teacher_schedule_preferences WHERE school_id = $1", schoolID)
	prefStmt, err := dbConn.PrepareContext(ctx,
		"INSERT INTO teacher_schedule_preferences (id, school_id, teacher_id, academic_year_id, day_of_week, hour_index, preference_type, reason, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()) ON CONFLICT (teacher_id, academic_year_id, day_of_week, hour_index) DO UPDATE SET preference_type = EXCLUDED.preference_type, reason = EXCLUDED.reason")
	if err != nil {
		return fmt.Errorf("failed to prepare teacher preferences: %w", err)
	}
	defer prefStmt.Close()

	for idx, t := range allCreatedTeachers {
		dayOff := 1 + (idx % 5)
		for h := 1; h <= 6; h++ {
			_, _ = prefStmt.ExecContext(ctx, uuid.New().String(), schoolID, t.UserID, academicYearID,
				dayOff, h, "unavailable", "Richiesta giorno libero settimanale")
		}
		prefType := "early"
		if idx%2 == 1 {
			prefType = "late"
		}
		for d := 1; d <= 5; d++ {
			if d == dayOff {
				continue
			}
			if prefType == "early" {
				_, _ = prefStmt.ExecContext(ctx, uuid.New().String(), schoolID, t.UserID, academicYearID, d, 1, "preferred", "Preferenza prime ore")
				_, _ = prefStmt.ExecContext(ctx, uuid.New().String(), schoolID, t.UserID, academicYearID, d, 2, "preferred", "Preferenza prime ore")
			} else {
				_, _ = prefStmt.ExecContext(ctx, uuid.New().String(), schoolID, t.UserID, academicYearID, d, 4, "preferred", "Preferenza ore tarde")
				_, _ = prefStmt.ExecContext(ctx, uuid.New().String(), schoolID, t.UserID, academicYearID, d, 5, "preferred", "Preferenza ore tarde")
			}
		}
		if idx%5 == 2 {
			extraDay := 2
			if extraDay == dayOff {
				extraDay = 3
			}
			_, _ = prefStmt.ExecContext(ctx, uuid.New().String(), schoolID, t.UserID, academicYearID,
				extraDay, 5, "unavailable", "Aggiornamento professionale")
		}
	}
	log.Println("Preferenze orarie configurate per tutti i docenti")

	_, _ = dbConn.ExecContext(ctx,
		"DELETE FROM timetable_constraints WHERE school_id = $1", schoolID)
	constraintStmt, err := dbConn.PrepareContext(ctx,
		"INSERT INTO timetable_constraints (id, school_id, constraint_type, target_type, target_id, parameters, is_hard, priority, is_active, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE, NOW())")
	if err != nil {
		return fmt.Errorf("failed to prepare timetable_constraints: %w", err)
	}
	defer constraintStmt.Close()

	for _, cl := range allClasses {
		p, _ := json.Marshal(map[string]int{"min_hours_per_day": 4, "max_hours_per_day": 6})
		cIDCopy := cl.ID
		_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "class_daily_hours", "class", cIDCopy, p, false, 8)
	}
	p1, _ := json.Marshal(map[string]int{"max_hours": 4})
	_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "max_daily_teacher_hours", nil, nil, p1, false, 7)
	p2, _ := json.Marshal(map[string]bool{"minimize_gaps": true})
	_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "no_consecutive_gaps", nil, nil, p2, false, 5)
	p3, _ := json.Marshal(map[string]interface{}{"avoid_last_hour": true, "subjects": []string{subjectIDs["MAT"], subjectIDs["FIS"], subjectIDs["ITA"]}})
	_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "avoid_heavy_last_hour", nil, nil, p3, false, 4)
	p4, _ := json.Marshal(map[string]interface{}{"subject_id": subjectIDs["EDF"], "required_type": "gym", "distribute": true})
	_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "lab_booking_required", nil, nil, p4, true, 10)
	p5, _ := json.Marshal(map[string]interface{}{"subject_id": subjectIDs["INF"], "required_type": "lab_computer", "distribute": true})
	_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "lab_booking_required", nil, nil, p5, true, 10)
	p6, _ := json.Marshal(map[string]bool{"is_open": true})
	_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "desiderata_window", nil, nil, p6, false, 1)
	p7, _ := json.Marshal(map[string]interface{}{"max_same_subject_per_day": 2, "spread_over_week": true})
	_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "subject_distribution", nil, nil, p7, false, 6)

	// Gruppi associati: Spagnolo classi E-F anni 3-5
	for yr := 3; yr <= 5; yr++ {
		var eID, fID string
		for _, c := range allClasses {
			if c.Name == fmt.Sprintf("%dE", yr) {
				eID = c.ID
			}
			if c.Name == fmt.Sprintf("%dF", yr) {
				fID = c.ID
			}
		}
		if eID != "" && fID != "" && len(teachersBySubject["SPA"]) > 0 {
			t := teachersBySubject["SPA"][(yr-3)%len(teachersBySubject["SPA"])]
			ap, _ := json.Marshal(map[string]interface{}{
				"name":           fmt.Sprintf("Lab Spagnolo %dE-%dF", yr, yr),
				"subject_id":     subjectIDs["SPA"],
				"teacher_id":     t.TeacherProfID,
				"class_ids":      []string{eID, fID},
				"hours_per_week": 1,
			})
			_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "associated_group", nil, nil, ap, true, 10)
			log.Printf("Gruppo associato: Spagnolo %dE-%dF (%s)", yr, yr, t.FullName)
		}
	}
	// Gruppi associati: Informatica classi G-H anni 3-5
	for yr := 3; yr <= 5; yr++ {
		var gID, hID string
		for _, c := range allClasses {
			if c.Name == fmt.Sprintf("%dG", yr) {
				gID = c.ID
			}
			if c.Name == fmt.Sprintf("%dH", yr) {
				hID = c.ID
			}
		}
		if gID != "" && hID != "" && len(teachersBySubject["INF"]) > 0 {
			t := teachersBySubject["INF"][(yr-3)%len(teachersBySubject["INF"])]
			ap, _ := json.Marshal(map[string]interface{}{
				"name":           fmt.Sprintf("Lab INF %dG-%dH", yr, yr),
				"subject_id":     subjectIDs["INF"],
				"teacher_id":     t.TeacherProfID,
				"class_ids":      []string{gID, hID},
				"hours_per_week": 2,
			})
			_, _ = constraintStmt.ExecContext(ctx, uuid.New().String(), schoolID, "associated_group", nil, nil, ap, true, 10)
			log.Printf("Gruppo associato: INF %dG-%dH (%s)", yr, yr, t.FullName)
		}
	}

	// Statistiche finali
	minHh, maxHh, totalHh := 999, 0, 0
	for _, t := range allCreatedTeachers {
		if t.CurrentHours < minHh { minHh = t.CurrentHours }
		if t.CurrentHours > maxHh { maxHh = t.CurrentHours }
		totalHh += t.CurrentHours
	}
	avgHh := 0
	if len(allCreatedTeachers) > 0 { avgHh = totalHh / len(allCreatedTeachers) }

	log.Println("==================================================================")
	log.Printf("SEEDING COMPLETATO: %s", schoolName)
	log.Printf("  Plessi: 3 | Lab/Palestre: %d | Materie: 14", roomCount)
	log.Printf("  Classi: %d (1A-5M) | Docenti: %d", len(allClasses), len(allCreatedTeachers))
	log.Printf("  Ore tot: %d/sett | Min: %dh | Max: %dh | Media: %dh/docente",
		totalAssignedHours, minHh, maxHh, avgHh)
	log.Println("  Credenziali:")
	log.Println("    Admin:     admin.leopardi@scuola.it / password")
	log.Println("    Dirigente: dirigente.leopardi@scuola.it / password")
	log.Println("    Vicario:   vicario.leopardi@scuola.it / password")
	log.Println("    DSGA:      dsga.leopardi@scuola.it / password")
	log.Println("    Segreteria: segreteria.leopardi@scuola.it / password")
	log.Println("    Docenti:   [nome].[cognome].leopardi@scuola.it / password")
	log.Println("==================================================================")
	return nil
}
