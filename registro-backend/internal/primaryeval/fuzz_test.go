package primaryeval

import (
	"context"
	"strings"
	"testing"
)

func FuzzIsValidPrimaryLevel(f *testing.F) {
	seeds := []string{
		LevelAvanzato,
		LevelIntermedio,
		LevelBase,
		LevelInViaPrimaAcquisiz,
		"ottimo",
		"insufficiente",
		"AVANZATO",
		" avanzato ",
		"\x00\xff",
		"'; DROP TABLE evaluations; --",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		isValid := IsValidPrimaryLevel(input)
		switch input {
		case LevelAvanzato, LevelIntermedio, LevelBase, LevelInViaPrimaAcquisiz:
			if !isValid {
				t.Fatalf("expected true for ministerial level %s", input)
			}
		default:
			if isValid {
				t.Fatalf("expected false for unexpected level %s", input)
			}
		}
	})
}

func FuzzCreateObjectiveValidation(f *testing.F) {
	f.Add("title", "desc", 1, "2025/2026")
	f.Add("", "desc", 3, "2025/2026")
	f.Add("Scienze", "", 0, "2025/2026")
	f.Add("Matematica", "numeri", 6, "")
	f.Add("Musica", "canto", 5, "anno")

	f.Fuzz(func(t *testing.T, title, desc string, yearGrade int, academicYear string) {
		repo := newMockRepository()
		svc := NewService(repo)

		req := &CreateObjectiveRequest{
			Title:        title,
			Description:  desc,
			YearGrade:    yearGrade,
			SubjectID:    "sub-fuzz",
			AcademicYear: academicYear,
		}

		res, err := svc.CreateObjective(context.Background(), "school-fuzz", req)
		if strings.TrimSpace(title) == "" && err == nil {
			t.Fatalf("expected error for empty title")
		}
		if (yearGrade < 1 || yearGrade > 5) && err == nil {
			t.Fatalf("expected error for invalid year grade %d", yearGrade)
		}
		if err == nil && res == nil {
			t.Fatalf("expected non-nil response when err is nil")
		}
	})
}
