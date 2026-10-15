package privacy

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkEvaluateTrafficLight(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = EvaluateTrafficLight(i%2 == 0, i%3 == 0, i%5 == 0)
	}
}

func BenchmarkGetClassBadges(b *testing.B) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo)

	ids := make([]string, 30)
	for i := 0; i < 30; i++ {
		id := fmt.Sprintf("std-%d", i)
		ids[i] = id
		_ = repo.SaveConsent(ctx, &StudentConsent{
			StudentID:         id,
			SchoolYear:        "2025/2026",
			TrafficLightBadge: TrafficLightGreen,
		})
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = svc.GetClassBadges(ctx, ids, "2025/2026")
	}
}
