package helpdesk

import (
	"testing"
)

func TestCalculateFISHours(t *testing.T) {
	// 15:00 to 16:30 = 1.5 hours
	hours, err := CalculateSlotHours("15:00", "16:30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hours != 1.5 {
		t.Errorf("expected 1.5 hours, got %f", hours)
	}

	// 14:30 to 16:30 = 2.0 hours
	hours2, err := CalculateSlotHours("14:30", "16:30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hours2 != 2.0 {
		t.Errorf("expected 2.0 hours, got %f", hours2)
	}

	// Invalid times: start after end
	_, err = CalculateSlotHours("17:00", "16:00")
	if err == nil {
		t.Errorf("expected error when start is after end")
	}
}

func TestSlotCapacityCheck(t *testing.T) {
	slot := HelpDeskSlot{
		MaxCapacity: 4,
	}

	if !CanBookSlot(slot, 3) {
		t.Errorf("slot with 3 bookings and max capacity 4 should be bookable")
	}

	if CanBookSlot(slot, 4) {
		t.Errorf("slot with 4 bookings and max capacity 4 should NOT be bookable")
	}

	if CanBookSlot(slot, 5) {
		t.Errorf("slot with 5 bookings and max capacity 4 should NOT be bookable")
	}
}
