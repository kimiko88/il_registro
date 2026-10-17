package meals

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkRecordRollCallAndAggregate(b *testing.B) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo)

	// Prepopulate 100 students
	studentIDs := make([]string, 100)
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("std-%d", i)
		studentIDs[i] = id
		if i%5 == 0 {
			_ = repo.SaveDiet(ctx, &SpecialDiet{
				StudentID:    id,
				SpecificDiet: "senza_glutine",
				IsApproved:   true,
			})
		}
	}

	req := MealRollCallRequest{
		ClassID:    "class-1A",
		Date:       "2026-10-15",
		PresentIDs: studentIDs,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = svc.RecordRollCallAndAggregate(ctx, "school-1", req)
	}
}
