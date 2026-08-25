package grades

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- Validator Tests ---

func TestValidateGradeValue(t *testing.T) {
	v := NewValidator(nil) // DB not needed for value check

	tests := []struct {
		name       string
		value      float64
		gradeType  string
		wantErr    bool
		wantErrMsg string
	}{
		{"Valid numeric 7.5", 7.5, string(GradeTypeNumeric), false, ""},
		{"Invalid numeric 0", 0, string(GradeTypeNumeric), true, "Voto deve essere tra 1 e 10"},
		{"Valid numeric 10", 10, string(GradeTypeNumeric), false, ""},
		{"Invalid numeric > 10", 10.1, string(GradeTypeNumeric), true, "Voto deve essere tra 1 e 10"},
		{"Invalid numeric < 0 (except -1)", -2, string(GradeTypeNumeric), true, "Voto deve essere tra 1 e 10"},
		{"Valid numeric -1 (absence)", -1, string(GradeTypeNumeric), false, ""},
		{"Invalid judgment value", -1, string(GradeTypeJudgment), true, "Valore giudizio fuori range"},
		{"Valid judgment value", 6, string(GradeTypeJudgment), false, ""},
		{"Valid credit", 5, string(GradeTypeCredit), false, ""},
		{"Invalid credit", 0, string(GradeTypeCredit), true, "Credito"},
		{"Invalid credit > 25", 26, string(GradeTypeCredit), true, "Credito"},
		{"Valid competence", 3, string(GradeTypeCompetence), false, ""},
		{"Invalid competence", 5, string(GradeTypeCompetence), true, "Livello competenza non valido"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateGradeValue(tt.value, tt.gradeType)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	v := NewValidator(nil)

	tests := []struct {
		name    string
		desc    string
		wantErr bool
	}{
		{"Valid description", "Compito in classe ben fatto", false},
		{"Too long", strings.Repeat("a", 501), true},
		{"SQL Injection Attempt", "'; DROP TABLE grades; --", false},
		{"Another SQLi", "SELECT * FROM users", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateDescription(tt.desc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateGradeDate(t *testing.T) {
	v := NewValidator(nil)
	// Semester 1: 2025-09-01 -> 2026-01-31

	tests := []struct {
		name     string
		dateStr  string
		semester int
		wantErr  bool
	}{
		{"Valid Date Sem 1", "2025-10-15", 1, false},
		{"Too Old", "2010-01-01", 1, true},
		{"Future Warning (Simulated)", "2030-01-01", 1, true}, // Assuming run in 2025
		{"Wrong Semester", "2026-03-01", 1, true},             // March is Sem 2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, _ := time.Parse("2006-01-02", tt.dateStr)
			err := v.ValidateGradeDate(d, tt.semester)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- Helper Tests ---

func TestConvertJudgmentToNumeric(t *testing.T) {
	tests := []struct {
		judgment string
		want     float64
	}{
		{"Insufficiente", 4.0},
		{"Mediocre", 5.0},
		{"Sufficiente", 6.0},
		{"Buono", 8.0},
		{"Ottimo", 10.0},
		{"Invalid", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.judgment, func(t *testing.T) {
			got := ConvertJudgmentToNumeric(tt.judgment)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertNumericToJudgment(t *testing.T) {
	tests := []struct {
		val  float64
		want string
	}{
		{3.0, "Gravemente Insufficiente"},
		{5.0, "Insufficiente"},
		{6.5, "Sufficiente"},
		{7.5, "Discreto"},
		{8.5, "Buono"},
		{9.5, "Distinto"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.1f", tt.val), func(t *testing.T) {
			got := ConvertNumericToJudgment(tt.val)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- Importer Tests ---

func TestParseCSVGrades(t *testing.T) {
	csvContent := `StudentID,SubjectID,GradeValue,Date,GradeCategory,Description
S1,SUB1,7.5,2025-10-15,sommativo,Test
S2,SUB1,6.0,2025-10-15,formativo,Quiz
S3,SUB1,invalid,,`

	r := strings.NewReader(csvContent)
	reqs, err := ParseCSVGrades(r, 1)

	assert.NoError(t, err)
	assert.Len(t, reqs, 2, "Invalid rows should be safely skipped")
}
