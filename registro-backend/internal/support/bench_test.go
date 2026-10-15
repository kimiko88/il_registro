package support

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkSupportDiaryFilter(b *testing.B) {
	svc := &mockSupportService{}
	for i := 0; i < 50; i++ {
		_, _ = svc.CreateDiaryEntry(context.Background(), "school-1", "teacher-1", CreateDiaryEntryRequest{
			StudentID:          fmt.Sprintf("stu-%d", i%5),
			ClassID:            "cls-1",
			EntryDate:          "2026-03-10",
			TimeSlot:           "1ª Ora (08:00-09:00)",
			ActivityType:       "individuale",
			TopicAndActivities: "Attività di potenziamento della lettura",
			IsSharedWithFamily: i%2 == 0,
		})
	}

	ctx := context.Background()
	for b.Loop() {
		_, _ = svc.ListDiaryEntries(ctx, "school-1", "stu-1", "", "cls-1", true)
	}
}
