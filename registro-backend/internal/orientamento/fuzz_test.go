package orientamento

import (
	"context"
	"testing"
	"time"
)

func FuzzCreateEventValidation(f *testing.F) {
	seeds := [][4]string{
		{"Open Day 2026", "2026-03-20T10:00:00Z", "2026-03-20T13:00:00Z", "school-1"},
		{"Workshop", "invalid-date", "also-invalid", "school-1"},
		{"Reverse Dates", "2026-03-20T13:00:00Z", "2026-03-20T10:00:00Z", "school-1"},
		{"No School", "2026-03-20T10:00:00Z", "2026-03-20T12:00:00Z", ""},
		{"Same Start End", "2026-03-20T10:00:00Z", "2026-03-20T10:00:00Z", "school-1"},
	}

	for _, s := range seeds {
		f.Add(s[0], s[1], s[2], s[3])
	}

	repo := &mockRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	f.Fuzz(func(t *testing.T, title, startDate, endDate, schoolID string) {
		req := CreateEventRequest{
			Title:   title,
			Date:    startDate,
			EndDate: endDate,
		}

		err := svc.CreateEvent(ctx, "teacher-1", schoolID, req)

		if schoolID == "" {
			if err == nil {
				t.Fatalf("expected error for empty schoolID")
			}
			return
		}

		start, errStart := time.Parse(time.RFC3339, startDate)
		end, errEnd := time.Parse(time.RFC3339, endDate)

		if errStart != nil || errEnd != nil {
			if err == nil {
				t.Fatalf("expected error for malformed RFC3339 dates")
			}
			return
		}

		if end.Before(start) {
			if err == nil {
				t.Fatalf("expected error when end is before start")
			}
			return
		}

		if err != nil {
			t.Fatalf("unexpected error for valid event: %v", err)
		}
	})
}

func FuzzSaveCapolavoro(f *testing.F) {
	seeds := []string{
		"Modello di Intelligenza Artificiale per la Diagnostica",
		"",
		"   ",
		"Progetto E-Portfolio 2026",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	repo := &mockRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	f.Fuzz(func(t *testing.T, title string) {
		c := Capolavoro{
			Title:       title,
			Description: "Descrizione del capolavoro",
			SchoolYear:  "2025/2026",
		}

		err := svc.SaveCapolavoro(ctx, "student-1", c)
		if title == "" {
			if err == nil {
				t.Fatalf("expected error for empty title")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error for valid title: %v", err)
			}
		}
	})
}
