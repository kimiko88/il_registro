package primaryeval

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func BenchmarkIsValidPrimaryLevel(b *testing.B) {
	levels := []string{
		LevelAvanzato,
		LevelIntermedio,
		LevelBase,
		LevelInViaPrimaAcquisiz,
		"invalid_level_here",
	}

	idx := 0
	for b.Loop() {
		_ = IsValidPrimaryLevel(levels[idx%len(levels)])
		idx++
	}
}

func BenchmarkPrimaryMatrixAggregation(b *testing.B) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	// Seed 10 objectives
	for i := 1; i <= 10; i++ {
		id := fmt.Sprintf("obj-%d", i)
		repo.objectives[id] = &LearningObjective{
			ID:        id,
			SchoolID:  "school-bench",
			Title:     fmt.Sprintf("Obiettivo %d", i),
			YearGrade: 3,
		}
	}

	// Seed 25 students x 10 objectives = 250 evaluations
	for s := 1; s <= 25; s++ {
		studentID := fmt.Sprintf("st-%d", s)
		studentName := fmt.Sprintf("Studente Numero %d", s)
		for o := 1; o <= 10; o++ {
			repo.evaluations = append(repo.evaluations, PrimaryEvaluation{
				ID:          fmt.Sprintf("eval-%d-%d", s, o),
				SchoolID:    "school-bench",
				StudentID:   studentID,
				StudentName: studentName,
				ClassID:     "cls-3a",
				SubjectID:   "sub-math",
				ObjectiveID: fmt.Sprintf("obj-%d", o),
				Level:       LevelIntermedio,
				Semester:    1,
				Date:        time.Now(),
			})
		}
	}

	for b.Loop() {
		_, err := svc.GetPrimaryMatrix(ctx, "school-bench", "cls-3a", "sub-math", 1)
		if err != nil {
			b.Fatalf("benchmark failed: %v", err)
		}
	}
}
