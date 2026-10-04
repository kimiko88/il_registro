package textbooks

import (
	"testing"
)

func TestCalculateSpendingStatus(t *testing.T) {
	limit := SpendingLimit{
		ClassYear:           1,
		SchoolOrder:         "secondaria_2",
		MaxAmount:           300.00,
		AllowedTolerancePct: 10.00, // max 330.00
		AcademicYear:        "2026/2027",
	}

	tests := []struct {
		name           string
		adoptions      []ClassAdoptionItem
		expectedTotal  float64
		expectedStatus SpendingStatus
	}{
		{
			name: "Within limit - standard adoptions",
			adoptions: []ClassAdoptionItem{
				{BookTitle: "Italiano", Price: 25.00, AdoptionType: "nuova_adozione", IsAlreadyOwned: false},
				{BookTitle: "Matematica", Price: 30.00, AdoptionType: "scorrimento", IsAlreadyOwned: false},
				{BookTitle: "Inglese", Price: 20.00, AdoptionType: "nuova_adozione", IsAlreadyOwned: false},
			},
			expectedTotal:  75.00,
			expectedStatus: StatusWithinLimit,
		},
		{
			name: "Excludes consigliato and already owned books",
			adoptions: []ClassAdoptionItem{
				{BookTitle: "Italiano", Price: 200.00, AdoptionType: "nuova_adozione", IsAlreadyOwned: false},
				{BookTitle: "Dizionario Consigliato", Price: 80.00, AdoptionType: "consigliato", IsAlreadyOwned: false},
				{BookTitle: "Libro Triennale Posseduto", Price: 60.00, AdoptionType: "scorrimento", IsAlreadyOwned: true},
			},
			expectedTotal:  200.00,
			expectedStatus: StatusWithinLimit,
		},
		{
			name: "Warning tolerance (+5% within +10%)",
			adoptions: []ClassAdoptionItem{
				{BookTitle: "Pacchetto Libri", Price: 315.00, AdoptionType: "nuova_adozione", IsAlreadyOwned: false},
			},
			expectedTotal:  315.00,
			expectedStatus: StatusWarningTolerance,
		},
		{
			name: "Over limit (+15% > +10%)",
			adoptions: []ClassAdoptionItem{
				{BookTitle: "Pacchetto Libri Costosi", Price: 345.00, AdoptionType: "nuova_adozione", IsAlreadyOwned: false},
			},
			expectedTotal:  345.00,
			expectedStatus: StatusOverLimit,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			report := CalculateSpendingReport(tc.adoptions, limit)
			if report.TotalSpending != tc.expectedTotal {
				t.Errorf("expected total %.2f, got %.2f", tc.expectedTotal, report.TotalSpending)
			}
			if report.Status != tc.expectedStatus {
				t.Errorf("expected status %s, got %s", tc.expectedStatus, report.Status)
			}
		})
	}
}
