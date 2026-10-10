package rooms

import (
	"context"
	"errors"
	"testing"
	"time"
)

func FuzzRoomBookingDateValidation(f *testing.F) {
	f.Add("2026-10-15")
	f.Add("2020-01-01")
	f.Add("2099-12-31")
	f.Add("invalid-date")
	f.Add("")
	f.Add("15/10/2026")
	f.Add("2026-02-31")

	f.Fuzz(func(t *testing.T, dateStr string) {
		repo := newMockRoomsRepo()
		svc := NewService(repo)

		// Seed room
		repo.rooms = append(repo.rooms, BookableRoom{
			ID:              "r-fuzz",
			SchoolID:        "school-1",
			Name:            "Lab",
			RoomType:        RoomTypeLabScienze,
			RequiresBooking: true,
			IsActive:        true,
		})

		req := CreateBookingRequest{
			RoomID:      "r-fuzz",
			BookingDate: dateStr,
			HourIndex:   2,
		}

		_, _, err := svc.CreateBooking(context.Background(), "school-1", "user-1", req)

		_, parseErr := time.Parse("2006-01-02", dateStr)
		if parseErr != nil {
			if !errors.Is(err, ErrInvalidDate) {
				t.Fatalf("expected ErrInvalidDate for %q, got %v", dateStr, err)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error for valid date format %s: %v", dateStr, err)
			}
		}
	})
}
