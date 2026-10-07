package albopretorio

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkPublishAct(b *testing.B) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = svc.PublishAct(ctx, "school-1", "user-ds", PublishActRequest{
			Category:        CategoryDetermineDirigente,
			Subject:         fmt.Sprintf("Determina num %d", i),
			DocumentFileURL: "https://scuola.edu.it/doc.pdf",
			DocumentSHA256:  "sha256placeholder",
		})
	}
}

func BenchmarkGenerateANACXML(b *testing.B) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo)

	for i := 0; i < 50; i++ {
		_, _ = svc.PublishAct(ctx, "school-1", "user-ds", PublishActRequest{
			Category:              CategoryDetermineDirigente,
			Subject:               fmt.Sprintf("Gara CIG %d", i),
			DocumentFileURL:       "https://scuola.edu.it/doc.pdf",
			DocumentSHA256:        "sha256placeholder",
			IsTransparencySection: true,
			CIGCode:               fmt.Sprintf("CIG%07d", i),
			AwardedAmount:         1000.0 * float64(i+1),
		})
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = svc.GenerateANACXML(ctx, "school-1", 2026)
	}
}
