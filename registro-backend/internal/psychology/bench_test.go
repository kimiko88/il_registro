package psychology

import (
	"context"
	"testing"
	"time"
)

func BenchmarkBookSession(b *testing.B) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo)

	// Pre-approve consent
	_, _ = svc.SignParentConsent(ctx, "school-1", "std-bench", "2025/2026", "parent-1")
	_, _ = svc.SignParentConsent(ctx, "school-1", "std-bench", "2025/2026", "parent-2")

	req := BookSessionRequest{
		SlotTime:        time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04"),
		PsychologistID:  "psycho-1",
		DurationMinutes: 45,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = svc.BookSession(ctx, "school-1", "std-bench", "2025/2026", true, req)
	}
}
