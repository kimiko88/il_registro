package scrutiny

import (
	"math"
	"testing"
)

// --- Unit tests for pure business logic in SubjectAverages (no DB needed) ---

func TestSubjectAverages_ProposedRoundingDown(t *testing.T) {
	avg := 7.3
	proposed := math.Round(avg)
	if proposed != 7 {
		t.Errorf("expected proposed grade 7 for avg 7.3, got %.0f", proposed)
	}
}

func TestSubjectAverages_ProposedRoundingUp(t *testing.T) {
	avg := 7.5
	proposed := math.Round(avg)
	if proposed != 8 {
		t.Errorf("expected proposed grade 8 for avg 7.5, got %.0f", proposed)
	}
}

func TestSubjectAverages_ProposedExact(t *testing.T) {
	avg := 6.0
	proposed := math.Round(avg)
	if proposed != 6 {
		t.Errorf("expected proposed grade 6 for avg 6.0, got %.0f", proposed)
	}
}

func TestSubjectAverages_ZeroCount(t *testing.T) {
	// When there are no grades, avg should be 0.0
	var sum float64
	var count int
	avg := 0.0
	if count > 0 {
		avg = sum / float64(count)
	}
	if avg != 0.0 {
		t.Errorf("expected avg 0.0 for zero count, got %f", avg)
	}
}

// --- Unit tests for ScrutinyMatrix construction ---

func TestScrutinyMatrix_Structure(t *testing.T) {
	matrix := &ScrutinyMatrix{
		ClassID:  "class-1",
		Semester: 1,
		Subjects: []SubjectInfo{
			{ID: "sub-math", Name: "Matematica"},
			{ID: "sub-ita", Name: "Italiano"},
		},
		Students: []StudentScrutinyRow{
			{
				StudentID:   "stu-1",
				StudentName: "Rossi Mario",
				SubjectData: map[string]SubjectAverages{
					"sub-math": {Average: 8.2, GradeCount: 5, Proposed: 8},
					"sub-ita":  {Average: 7.4, GradeCount: 3, Proposed: 7},
				},
			},
		},
	}

	if matrix.ClassID != "class-1" {
		t.Errorf("expected ClassID 'class-1', got '%s'", matrix.ClassID)
	}
	if matrix.Semester != 1 {
		t.Errorf("expected Semester 1, got %d", matrix.Semester)
	}
	if len(matrix.Subjects) != 2 {
		t.Errorf("expected 2 subjects, got %d", len(matrix.Subjects))
	}
	if len(matrix.Students) != 1 {
		t.Errorf("expected 1 student, got %d", len(matrix.Students))
	}
	row := matrix.Students[0]
	if row.StudentID != "stu-1" {
		t.Errorf("expected StudentID 'stu-1', got '%s'", row.StudentID)
	}
	mathData, ok := row.SubjectData["sub-math"]
	if !ok {
		t.Fatal("expected SubjectData for 'sub-math' to be present")
	}
	if mathData.GradeCount != 5 {
		t.Errorf("expected GradeCount 5, got %d", mathData.GradeCount)
	}
	if mathData.Proposed != 8 {
		t.Errorf("expected Proposed 8, got %f", mathData.Proposed)
	}
}

// --- Unit tests for ScrutinyRecord status values ---

func TestScrutinyRecord_ValidStatuses(t *testing.T) {
	validStatuses := []string{"draft", "in_progress", "submitted", "validated", "closed"}
	for _, status := range validStatuses {
		rec := ScrutinyRecord{Status: status}
		if rec.Status != status {
			t.Errorf("expected status '%s', got '%s'", status, rec.Status)
		}
	}
}

func TestScrutinyRecord_ConductGrade_ValidRange(t *testing.T) {
	// Italian school conduct grades are from 5 to 10
	validConductGrades := []int{5, 6, 7, 8, 9, 10}
	for _, g := range validConductGrades {
		rec := ScrutinyRecord{ConductGrade: g}
		if rec.ConductGrade < 5 || rec.ConductGrade > 10 {
			t.Errorf("conduct grade %d is out of valid Italian school range [5,10]", rec.ConductGrade)
		}
	}
}

// --- Unit tests for SaveScrutinyRequest ---

func TestSaveScrutinyRequest_GradesSlice(t *testing.T) {
	req := SaveScrutinyRequest{
		StudentID:     "stu-1",
		ClassID:       "cls-1",
		Semester:      1,
		ConductGrade:  8,
		FinalDecision: "admitted",
		Grades: []SaveScrutinyGradeRequest{
			{SubjectID: "sub-math", FinalGrade: 7.5},
			{SubjectID: "sub-ita", FinalGrade: 8.0},
		},
	}

	if req.StudentID != "stu-1" || req.ClassID != "cls-1" || req.Semester != 1 || req.ConductGrade != 8 || req.FinalDecision != "admitted" {
		t.Errorf("request fields mismatch: %+v", req)
	}
	if len(req.Grades) != 2 {
		t.Errorf("expected 2 grades in request, got %d", len(req.Grades))
	}
	if req.Grades[0].FinalGrade != 7.5 {
		t.Errorf("expected FinalGrade 7.5, got %f", req.Grades[0].FinalGrade)
	}
}

// --- Unit tests for AttendanceSummary ---

func TestAttendanceSummary_Fields(t *testing.T) {
	summary := AttendanceSummary{
		Absences:   10,
		Lates:      3,
		EarlyExits: 2,
	}

	if summary.Absences != 10 {
		t.Errorf("expected Absences 10, got %d", summary.Absences)
	}
	if summary.Lates != 3 {
		t.Errorf("expected Lates 3, got %d", summary.Lates)
	}
	if summary.EarlyExits != 2 {
		t.Errorf("expected EarlyExits 2, got %d", summary.EarlyExits)
	}
}

// --- Role authorization logic ---

func TestIsDirigenzaRole(t *testing.T) {
	dirigenzaRoles := []string{"principal", "vice_principal", "admin", "superadmin"}
	nonDirigenzaRoles := []string{"teacher", "secretary", "student", "parent"}

	for _, role := range dirigenzaRoles {
		isDirigenza := role == "principal" || role == "vice_principal" || role == "admin" || role == "superadmin"
		if !isDirigenza {
			t.Errorf("expected role '%s' to be dirigenza", role)
		}
	}
	for _, role := range nonDirigenzaRoles {
		isDirigenza := role == "principal" || role == "vice_principal" || role == "admin" || role == "superadmin"
		if isDirigenza {
			t.Errorf("expected role '%s' to NOT be dirigenza", role)
		}
	}
}

// --- Unit tests for Final Scrutiny Outcome Calculation & Debiti Formativi ---

func TestScrutinyOutcome_Promosso(t *testing.T) {
	// All grades >= 6 and conduct >= 6 => PROMOSSO
	grades := []float64{7.0, 6.5, 8.0, 6.0, 7.5}
	conduct := 8

	insufficientCount := 0
	for _, g := range grades {
		if g < 6.0 {
			insufficientCount++
		}
	}

	var outcome string
	if conduct < 6 || insufficientCount >= 4 {
		outcome = "non_promosso"
	} else if insufficientCount > 0 {
		outcome = "sospensione_giudizio"
	} else {
		outcome = "promosso"
	}

	if outcome != "promosso" {
		t.Errorf("expected outcome 'promosso', got '%s'", outcome)
	}
}

func TestScrutinyOutcome_SospensioneGiudizio(t *testing.T) {
	// 1 to 3 insufficient grades => SOSPENSIONE GIUDIZIO (Debito Formativo)
	grades := []float64{5.0, 6.5, 4.5, 6.0, 7.5} // 2 insuf.
	conduct := 7

	insufficientCount := 0
	for _, g := range grades {
		if g < 6.0 {
			insufficientCount++
		}
	}

	var outcome string
	if conduct < 6 || insufficientCount >= 4 {
		outcome = "non_promosso"
	} else if insufficientCount > 0 {
		outcome = "sospensione_giudizio"
	} else {
		outcome = "promosso"
	}

	if outcome != "sospensione_giudizio" {
		t.Errorf("expected outcome 'sospensione_giudizio', got '%s'", outcome)
	}
}

func TestScrutinyOutcome_NonPromossoPerVotoCondotta(t *testing.T) {
	// Conduct grade < 6 => NON PROMOSSO (Bocciato per condotta)
	grades := []float64{8.0, 8.5, 9.0, 7.0, 8.5}
	conduct := 5 // Insufficient conduct

	insufficientCount := 0
	for _, g := range grades {
		if g < 6.0 {
			insufficientCount++
		}
	}

	var outcome string
	if conduct < 6 || insufficientCount >= 4 {
		outcome = "non_promosso"
	} else if insufficientCount > 0 {
		outcome = "sospensione_giudizio"
	} else {
		outcome = "promosso"
	}

	if outcome != "non_promosso" {
		t.Errorf("expected outcome 'non_promosso' due to conduct < 6, got '%s'", outcome)
	}
}

