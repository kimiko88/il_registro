package enrollment

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func FuzzParseSIDIExport(f *testing.F) {
	// Seed valid and partial CSV streams
	validCSV := `ID_DOMANDA;NOME;COGNOME;CODICE_FISCALE;DATA_NASCITA;SESSO;VOTO_LICENZA;INDIRIZZO;SECONDA_LINGUA;L104;DSA;RELIGIONE;COMPAGNI_RICHIESTI;GENITORE1_NOME;GENITORE1_COGNOME;GENITORE1_EMAIL;GENITORE1_TEL
SIDI-01;Lorenzo;Galli;GLLLNZ10A01H501K;2012-04-10;M;9;Scientifico;Spagnolo;false;false;irc;;Paolo;Galli;paolo@example.com;3330001111
SIDI-02;Chiara;Fontana;FNTCHR10A41H501W;2012-06-25;F;8;Scientifico;Spagnolo;false;false;irc;;Marco;Fontana;marco@example.com;3330002222
`
	f.Add([]byte(validCSV))
	f.Add([]byte("invalid header only\n"))
	f.Add([]byte("NOME,COGNOME,CODICE_FISCALE\nMario,Rossi,RSSMRA12A01H501Z\n"))
	f.Add([]byte(""))
	f.Add([]byte("\x00\xff\xfe\xfd;NOME;COGNOME"))

	f.Fuzz(func(t *testing.T, data []byte) {
		r := bytes.NewReader(data)
		apps, err := ParseSIDIExport(r)
		if err == nil {
			// If parsing succeeded, verify each returned application has non-empty TaxCode
			for _, app := range apps {
				if app.StudentTaxCode == "" {
					t.Fatalf("parsed application has empty tax code")
				}
				if app.Gender != "M" && app.Gender != "F" {
					t.Fatalf("unexpected gender %s", app.Gender)
				}
			}
		}
	})
}

func FuzzFormationSolver(f *testing.F) {
	f.Add(10, 2, true, true)
	f.Add(30, 3, false, true)
	f.Add(50, 2, true, false)
	f.Add(5, 0, false, false)
	f.Add(0, 1, true, true)

	f.Fuzz(func(t *testing.T, studentCount, targetClasses int, balanceGender, balanceGrades bool) {
		if studentCount < 0 || studentCount > 200 {
			return
		}
		if targetClasses < -5 || targetClasses > 20 {
			return
		}

		pool := make([]EnrollmentApplication, studentCount)
		for i := 0; i < studentCount; i++ {
			gender := "M"
			if i%2 == 1 {
				gender = "F"
			}
			pool[i] = EnrollmentApplication{
				SidiApplicationID: fmt.Sprintf("APP-%d", i),
				StudentFirstName:  fmt.Sprintf("Nome%d", i),
				StudentLastName:   fmt.Sprintf("Cognome%d", i),
				StudentTaxCode:    fmt.Sprintf("CF%014d", i),
				Gender:            gender,
				MiddleSchoolGrade: 6 + (i % 5),
				HasDisabilityL104: i%7 == 0,
				HasDSA:            i%5 == 0,
			}
		}

		solver := NewFormationSolver()
		params := FormationParams{
			TargetClassCount: targetClasses,
			BalanceGender:    balanceGender,
			BalanceGrades:    balanceGrades,
		}

		res, err := solver.Solve(pool, params)
		if studentCount == 0 {
			if err == nil {
				t.Fatalf("expected error for empty student pool")
			}
			return
		}

		if err != nil && strings.Contains(err.Error(), "empty student pool") {
			return
		}

		if err == nil {
			if len(res.Classes) == 0 {
				t.Fatalf("expected at least one class formed")
			}
			totalAssigned := 0
			for _, cls := range res.Classes {
				totalAssigned += len(cls.Students)
			}
			if totalAssigned != studentCount {
				t.Fatalf("student count mismatch: pool=%d, assigned=%d", studentCount, totalAssigned)
			}
		}
	})
}
