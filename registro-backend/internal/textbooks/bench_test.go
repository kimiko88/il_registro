package textbooks

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func BenchmarkCalculateSpendingReport(b *testing.B) {
	adoptions := make([]ClassAdoptionItem, 20)
	for i := 0; i < 20; i++ {
		adoptions[i] = ClassAdoptionItem{
			BookID:         fmt.Sprintf("book-%d", i),
			Price:          15.50 + float64(i)*1.5,
			AdoptionType:   "obbligatorio",
			IsAlreadyOwned: i%4 == 0,
		}
	}

	limit := SpendingLimit{
		MaxAmount:           320.00,
		AllowedTolerancePct: 10.0,
	}

	for b.Loop() {
		_ = CalculateSpendingReport(adoptions, limit)
	}
}

func BenchmarkParseAIECatalog(b *testing.B) {
	var sb strings.Builder
	sb.WriteString("CODICE_ISBN;TITOLO;AUTORI;EDITORE;PREZZO;MATERIA;CLASSE;ANNO\n")
	for i := 0; i < 50; i++ {
		sb.WriteString(fmt.Sprintf("978880%07d;Libro di Test %d;Autore %d;Editore;%d.50;Materia;1;2025\n",
			i, i, i, 15+(i%20)))
	}
	catalogBytes := []byte(sb.String())

	for b.Loop() {
		r := bytes.NewReader(catalogBytes)
		_, _ = ParseAIECatalog(r)
	}
}
