package textbooks

import (
	"bytes"
	"math"
	"testing"
)

func FuzzParseAIECatalog(f *testing.F) {
	seeds := [][]byte{
		[]byte("CODICE_ISBN;TITOLO;AUTORI;EDITORE;PREZZO;MATERIA;CLASSE;ANNO\n9788808123456;Matematica Blu;Bergamini;Zanichelli;28.50;Matematica;1;2025\n"),
		[]byte("ISBN,TITLE,AUTHOR,PUBLISHER,PRICE,SUBJECT,CLASS,YEAR\n1234567890,Fisica,Walker,Linx,32.00,Fisica,2,2025\n"),
		[]byte(""),
		[]byte("INVALID HEADER ONLY"),
		[]byte("A;B;C\n1;2;3\n4;5"),
		[]byte("ISBN;PREZZO\n97888;INVALID_PRICE\n"),
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		r := bytes.NewReader(data)
		books, err := ParseAIECatalog(r)
		if err == nil {
			if books == nil {
				t.Fatalf("expected non-nil books slice when error is nil")
			}
		}
	})
}

func FuzzCalculateSpendingReport(f *testing.F) {
	f.Add(28.50, 32.00, 15.00, 300.0, 10.0)
	f.Add(150.0, 200.0, 0.0, 300.0, 10.0)
	f.Add(0.0, 0.0, 0.0, 0.0, 0.0)
	f.Add(500.0, 500.0, 500.0, 100.0, 20.0)

	f.Fuzz(func(t *testing.T, p1, p2, p3, limitAmount, tol float64) {
		if math.IsNaN(p1) || math.IsNaN(p2) || math.IsNaN(p3) || math.IsNaN(limitAmount) || math.IsNaN(tol) ||
			math.IsInf(p1, 0) || math.IsInf(p2, 0) || math.IsInf(p3, 0) || math.IsInf(limitAmount, 0) || math.IsInf(tol, 0) {
			return
		}

		adoptions := []ClassAdoptionItem{
			{Price: math.Abs(p1), AdoptionType: "obbligatorio", IsAlreadyOwned: false},
			{Price: math.Abs(p2), AdoptionType: "consigliato", IsAlreadyOwned: false},
			{Price: math.Abs(p3), AdoptionType: "obbligatorio", IsAlreadyOwned: true},
		}

		limit := SpendingLimit{
			MaxAmount:           math.Abs(limitAmount),
			AllowedTolerancePct: math.Abs(tol),
		}

		report := CalculateSpendingReport(adoptions, limit)

		expectedTotal := math.Round(math.Abs(p1)*100) / 100
		if report.TotalSpending != expectedTotal {
			t.Fatalf("expected total spending %.2f, got %.2f", expectedTotal, report.TotalSpending)
		}

		switch report.Status {
		case StatusWithinLimit, StatusWarningTolerance, StatusOverLimit:
		default:
			t.Fatalf("invalid report status: %v", report.Status)
		}
	})
}
