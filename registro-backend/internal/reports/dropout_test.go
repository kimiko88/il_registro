package reports

import (
	"context"
	"strings"
	"testing"
)

func TestCalculateStudentRisk(t *testing.T) {
	tests := []struct {
		name         string
		absenceRate  float64
		failingCount int
		delays       int
		wantLevel    string
	}{
		{
			name:         "Critical over 25% DPR 122/2009",
			absenceRate:  25.5,
			failingCount: 0,
			delays:       2,
			wantLevel:    "Critico",
		},
		{
			name:         "High risk near 20% and 3 failing subjects",
			absenceRate:  21.0,
			failingCount: 3,
			delays:       4,
			wantLevel:    "Critico",
		},
		{
			name:         "High risk with 3 failing subjects",
			absenceRate:  10.0,
			failingCount: 3,
			delays:       2,
			wantLevel:    "Alto",
		},
		{
			name:         "Moderate risk with 1 failing and 16% absence",
			absenceRate:  16.0,
			failingCount: 1,
			delays:       1,
			wantLevel:    "Moderato",
		},
		{
			name:         "Low risk normal student",
			absenceRate:  4.0,
			failingCount: 0,
			delays:       0,
			wantLevel:    "Basso",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, level, action := CalculateStudentRisk(tt.absenceRate, tt.failingCount, tt.delays)
			if level != tt.wantLevel {
				t.Errorf("CalculateStudentRisk(%v, %v, %v) level = %v, want %v (score: %.1f)", tt.absenceRate, tt.failingCount, tt.delays, level, tt.wantLevel, score)
			}
			if action == "" {
				t.Errorf("CalculateStudentRisk action should not be empty")
			}
		})
	}
}

func TestExportDropoutRiskCSV_EmptyDB(t *testing.T) {
	svc := NewService(nil, nil)
	csvBytes, err := svc.ExportDropoutRiskCSV(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ExportDropoutRiskCSV error: %v", err)
	}
	content := string(csvBytes)
	if !strings.Contains(content, "ID Studente") || !strings.Contains(content, "Tasso Assenze") {
		t.Errorf("expected CSV headers, got: %s", content)
	}
}
