package schoolcalendar

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func BenchmarkListNonTeachingDays(b *testing.B) {
	repo := &mockCalendarRepo{}
	for i := 1; i <= 30; i++ {
		d := &NonTeachingDay{
			ID:        fmt.Sprintf("ntd-%d", i),
			SchoolID:  "school-bench",
			Date:      time.Date(2025, 10, i%28+1, 0, 0, 0, 0, time.UTC),
			Label:     fmt.Sprintf("Holiday %d", i),
			CreatedBy: "sec-1",
		}
		_ = repo.AddNonTeachingDay(d)
	}

	svc := NewService(repo)
	ctx := context.Background()

	for b.Loop() {
		_, _ = svc.ListNonTeachingDays(ctx, "school-bench")
	}
}

func BenchmarkGetSchoolYearDates(b *testing.B) {
	repo := &mockCalendarRepo{
		year: &SchoolYearSettings{
			ID:        "sy-1",
			SchoolID:  "school-bench",
			YearLabel: "2025/2026",
			StartDate: time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	svc := NewService(repo)
	ctx := context.Background()

	for b.Loop() {
		_, _, _ = svc.GetSchoolYearDates(ctx, "school-bench")
	}
}
