package schoolcalendar

import (
	"context"
	"testing"
	"time"
)

func FuzzSchoolYearValidation(f *testing.F) {
	seeds := [][3]string{
		{"2025/2026", "2025-09-15", "2026-06-10"},
		{"2025/2026", "2026-06-10", "2025-09-15"},
		{"invalid", "not-a-date", "also-not-a-date"},
		{"2024/2025", "2024-09-01", "2024-09-01"},
		{"", "2025-01-01", "2025-12-31"},
		{"2025", "2025-02-30", "2025-02-31"},
	}

	for _, s := range seeds {
		f.Add(s[0], s[1], s[2])
	}

	repo := &mockCalendarRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	f.Fuzz(func(t *testing.T, yearLabel, startDate, endDate string) {
		req := SetYearRequest{
			YearLabel: yearLabel,
			StartDate: startDate,
			EndDate:   endDate,
		}

		res, err := svc.SetSchoolYear(ctx, "sec-1", "secretary", "school-1", req)

		parsedStart, errStart := time.Parse("2006-01-02", startDate)
		parsedEnd, errEnd := time.Parse("2006-01-02", endDate)

		if errStart != nil || errEnd != nil || !parsedEnd.After(parsedStart) {
			if err == nil {
				t.Fatalf("expected error for invalid dates: start=%s end=%s", startDate, endDate)
			}
			return
		}

		// When valid and end is strictly after start
		if err != nil {
			t.Fatalf("unexpected error for valid dates: %v", err)
		}
		if res == nil {
			t.Fatalf("expected response not to be nil")
		}
	})
}
