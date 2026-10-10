package may15

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

func BenchmarkGetDocument(b *testing.B) {
	repo := new(MockMay15Repo)
	svc := NewService(repo)
	ctx := context.Background()

	doc := &ClassMay15Document{
		ID:           "doc-bench",
		ClassID:      "class-5a",
		AcademicYear: "2025/2026",
		Status:       StatusPubblicato,
	}
	repo.On("GetByClassAndYear", mock.Anything, "class-5a", "2025/2026").Return(doc, nil)

	for b.Loop() {
		_, _ = svc.GetDocument(ctx, "class-5a", "2025/2026")
	}
}

func BenchmarkSaveDocument(b *testing.B) {
	repo := new(MockMay15Repo)
	svc := NewService(repo)
	ctx := context.Background()

	repo.On("Upsert", mock.Anything, mock.Anything).Return(nil)

	req := SaveMay15Request{
		AcademicYear:       "2025/2026",
		Status:             StatusApprovatoCdC,
		ClassPresentation:  "Presentazione approfondita della classe quinta",
		TeachingContinuity: "Continuità su tutte le materie principali",
		PCTOPathways:       "Percorsi per le competenze trasversali e l'orientamento",
	}

	for b.Loop() {
		_, _ = svc.SaveDocument(ctx, "school-1", "class-5a", req)
	}
}
