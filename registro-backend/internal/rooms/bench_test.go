package rooms

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func BenchmarkListRooms(b *testing.B) {
	repo := newMockRoomsRepo()
	svc := NewService(repo)
	ctx := context.Background()

	// Seed 30 rooms
	for i := 0; i < 30; i++ {
		repo.rooms = append(repo.rooms, BookableRoom{
			ID:              fmt.Sprintf("room-%d", i),
			SchoolID:        "school-1",
			Name:            fmt.Sprintf("Aula/Lab %d", i),
			RoomType:        RoomTypeLabInformatica,
			Capacity:        25,
			RequiresBooking: true,
			IsActive:        true,
		})
	}

	for b.Loop() {
		_, _ = svc.ListRooms(ctx, "school-1", nil, nil, false)
	}
}

func BenchmarkCreateBookingConflictCheck(b *testing.B) {
	repo := newMockRoomsRepo()
	svc := NewService(repo)
	ctx := context.Background()

	repo.rooms = append(repo.rooms, BookableRoom{
		ID:              "room-bench",
		SchoolID:        "school-1",
		Name:            "Lab Chimica",
		RoomType:        RoomTypeLabChimica,
		Capacity:        25,
		RequiresBooking: true,
		IsActive:        true,
	})

	dateStr := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	for b.Loop() {
		repo.bookings = nil
		_, _, _ = svc.CreateBooking(ctx, "school-1", "teacher-1", CreateBookingRequest{
			RoomID:      "room-bench",
			BookingDate: dateStr,
			HourIndex:   2,
		})
	}
}
