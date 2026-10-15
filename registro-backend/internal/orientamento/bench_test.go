package orientamento

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkGetCurriculumStudente(b *testing.B) {
	repo := &mockRepo{}
	for i := 0; i < 5; i++ {
		repo.capolavori = append(repo.capolavori, Capolavoro{
			ID:          fmt.Sprintf("cap-%d", i),
			StudentID:   "student-bench",
			SchoolYear:  "2025/2026",
			Title:       fmt.Sprintf("Capolavoro %d", i),
			Description: "Descrizione dettagliata",
		})
	}

	svc := NewService(repo)
	ctx := context.Background()

	for b.Loop() {
		_, _ = svc.GetCurriculumStudente(ctx, "student-bench")
	}
}

func BenchmarkSaveCapolavoro(b *testing.B) {
	repo := &mockRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	c := Capolavoro{
		Title:       "Progetto Robotica",
		Description: "Costruzione braccio meccanico",
		SchoolYear:  "2025/2026",
	}

	for b.Loop() {
		_ = svc.SaveCapolavoro(ctx, "student-bench", c)
	}
}

func BenchmarkRegisterStudentDuplicateCheck(b *testing.B) {
	repo := &mockRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	for b.Loop() {
		_ = svc.RegisterStudent(ctx, "student-1", "event-new")
	}
}
