package enrollment

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func BenchmarkParseSIDIExport(b *testing.B) {
	// Generate a 200-row SIDI CSV in memory
	var sb strings.Builder
	sb.WriteString("ID_DOMANDA;NOME;COGNOME;CODICE_FISCALE;DATA_NASCITA;SESSO;VOTO_LICENZA;INDIRIZZO;SECONDA_LINGUA;L104;DSA;RELIGIONE;COMPAGNI_RICHIESTI;GENITORE1_NOME;GENITORE1_COGNOME;GENITORE1_EMAIL;GENITORE1_TEL\n")
	for i := 0; i < 200; i++ {
		gender := "M"
		if i%2 == 1 {
			gender = "F"
		}
		fmt.Fprintf(&sb, "SIDI-%04d;Nome%d;Cognome%d;TAXCODE%09d;2012-05-15;%s;%d;Scientifico;Francese;false;false;irc;;GenNome;GenCognome;gen@example.com;3331234567\n",
			i, i, i, i, gender, 6+(i%5))
	}
	csvData := []byte(sb.String())

	for b.Loop() {
		r := bytes.NewReader(csvData)
		_, err := ParseSIDIExport(r)
		if err != nil {
			b.Fatalf("unexpected parse error: %v", err)
		}
	}
}

func BenchmarkFormationSolver(b *testing.B) {
	// Generate 120 students pool
	pool := make([]EnrollmentApplication, 120)
	for i := 0; i < 120; i++ {
		gender := "M"
		if i%2 == 1 {
			gender = "F"
		}
		pool[i] = EnrollmentApplication{
			SidiApplicationID: fmt.Sprintf("APP-%d", i),
			StudentFirstName:  fmt.Sprintf("Studente%d", i),
			StudentLastName:   fmt.Sprintf("Cognome%d", i),
			StudentTaxCode:    fmt.Sprintf("CF%014d", i),
			Gender:            gender,
			MiddleSchoolGrade: 6 + (i % 5),
			HasDisabilityL104: i%15 == 0,
			HasDSA:            i%10 == 0,
		}
	}

	solver := NewFormationSolver()
	params := FormationParams{
		TargetClassCount: 5,
		BalanceGender:    true,
		BalanceGrades:    true,
		MaxL104PerClass:  1,
	}

	for b.Loop() {
		_, err := solver.Solve(pool, params)
		if err != nil {
			b.Fatalf("unexpected solve error: %v", err)
		}
	}
}
