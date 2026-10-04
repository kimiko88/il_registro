package enrollment

import (
	"fmt"
	"testing"
)

func TestFormationSolver_DPR81Compliance(t *testing.T) {
	// Create sample pool of 42 students:
	// - 2 students with L.104
	// - 20 Males, 22 Females
	// - Divided into 2 target classes (1A, 1B)
	var apps []EnrollmentApplication
	for i := 1; i <= 42; i++ {
		gender := "M"
		if i > 20 {
			gender = "F"
		}
		grade := 6 + (i % 5) // grades 6..10
		hasL104 := i == 1 || i == 2 // 2 disabled students
		apps = append(apps, EnrollmentApplication{
			ID:                 fmt.Sprintf("app-%d", i),
			StudentTaxCode:     fmt.Sprintf("CF%04d", i),
			StudentFirstName:   fmt.Sprintf("Nome%d", i),
			StudentLastName:    fmt.Sprintf("Cognome%d", i),
			Gender:             gender,
			MiddleSchoolGrade:  grade,
			HasDisabilityL104:  hasL104,
			SecondLanguage:     "Spagnolo",
		})
	}

	params := FormationParams{
		TargetClassCount: 2,
		MaxL104PerClass:  1,
		BalanceGender:    true,
		BalanceGrades:    true,
	}

	solver := NewFormationSolver()
	draft, err := solver.Solve(apps, params)
	if err != nil {
		t.Fatalf("unexpected error solving class formation: %v", err)
	}

	if len(draft.Classes) != 2 {
		t.Fatalf("expected 2 classes formed, got %d", len(draft.Classes))
	}

	for classIndex, classResult := range draft.Classes {
		// Normative compliance D.P.R. 81/2009:
		// If L.104 is present, class size must not exceed 20 students!
		if classResult.L104Count > 0 && classResult.TotalStudents > 21 {
			t.Errorf("class %d violates DPR 81/2009: %d students with %d L.104 students",
				classIndex, classResult.TotalStudents, classResult.L104Count)
		}

		if classResult.L104Count > params.MaxL104PerClass {
			t.Errorf("class %d exceeded max L.104: got %d, max %d",
				classIndex, classResult.L104Count, params.MaxL104PerClass)
		}

		// Check gender balance: should not be completely skewed
		if classResult.MalesCount == 0 || classResult.FemalesCount == 0 {
			t.Errorf("class %d has unskewed gender distribution: M=%d, F=%d",
				classIndex, classResult.MalesCount, classResult.FemalesCount)
		}
	}
}

func TestFormationSolver_RequestedClassmates(t *testing.T) {
	apps := []EnrollmentApplication{
		{
			ID:                 "app-1",
			StudentTaxCode:     "CF0001",
			StudentFirstName:   "Marco",
			StudentLastName:    "Rossi",
			Gender:             "M",
			MiddleSchoolGrade:  8,
			SecondLanguage:     "Inglese",
			RequestedClassmates: []string{"CF0002"},
		},
		{
			ID:                 "app-2",
			StudentTaxCode:     "CF0002",
			StudentFirstName:   "Luca",
			StudentLastName:    "Bianchi",
			Gender:             "M",
			MiddleSchoolGrade:  8,
			SecondLanguage:     "Inglese",
			RequestedClassmates: []string{"CF0001"},
		},
		{
			ID:                 "app-3",
			StudentTaxCode:     "CF0003",
			StudentFirstName:   "Sara",
			StudentLastName:    "Verdi",
			Gender:             "F",
			MiddleSchoolGrade:  7,
			SecondLanguage:     "Inglese",
		},
		{
			ID:                 "app-4",
			StudentTaxCode:     "CF0004",
			StudentFirstName:   "Elena",
			StudentLastName:    "Neri",
			Gender:             "F",
			MiddleSchoolGrade:  9,
			SecondLanguage:     "Inglese",
		},
	}

	params := FormationParams{
		TargetClassCount: 2,
		MaxL104PerClass:  1,
		BalanceGender:    true,
		BalanceGrades:    true,
	}

	solver := NewFormationSolver()
	draft, err := solver.Solve(apps, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that CF0001 and CF0002 are placed together in the same class
	var classApp1, classApp2 int
	for idx, cls := range draft.Classes {
		for _, s := range cls.Students {
			if s.StudentTaxCode == "CF0001" {
				classApp1 = idx
			}
			if s.StudentTaxCode == "CF0002" {
				classApp2 = idx
			}
		}
	}

	if classApp1 != classApp2 {
		t.Errorf("mutual requested classmates CF0001 and CF0002 were separated into classes %d and %d", classApp1, classApp2)
	}
}
