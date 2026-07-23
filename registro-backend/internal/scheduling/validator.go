package scheduling

import (
	"fmt"
	"time"
)

// ConflictType defines the nature of the schedule conflict
type ConflictType string

const (
	ConflictTeacherOverlap ConflictType = "TEACHER_OVERLAP"
	ConflictClassOverlap   ConflictType = "CLASS_OVERLAP"
	ConflictRoomOverlap    ConflictType = "ROOM_OVERLAP"
)

// Conflict represents a detected scheduling issue
type Conflict struct {
	Type        ConflictType `json:"type"`
	Description string       `json:"description"`
	SlotA       Slot         `json:"slot_a"`
	SlotB       Slot         `json:"slot_b"`
}

// Validator handles logic for checking schedule validity
type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateSlot(slot *ColloquioSlot) error {
	if slot.StartTime.After(slot.EndTime) {
		return fmt.Errorf("start time must be before end time")
	}
	// Check for past dates? The test expects an error for past date.
	// We need to combine Date + StartTime or just check Date if it's strictly day-based.
	// Assuming Date is midnight.
	slotDateTime := time.Date(slot.Date.Year(), slot.Date.Month(), slot.Date.Day(),
		slot.StartTime.Hour(), slot.StartTime.Minute(), 0, 0, slot.Date.Location())

	// Use a small buffer or strictly Now()
	if slotDateTime.Before(time.Now().Add(-1 * time.Hour)) {
		return fmt.Errorf("cannot create slot in the past")
	}

	return nil
}

func (v *Validator) ValidateBooking(slot *ColloquioSlot, settings *ColloquioSettings) error {
	if settings == nil {
		return nil // loose validation if no settings
	}

	// Calculate slot start time
	slotStart := time.Date(slot.Date.Year(), slot.Date.Month(), slot.Date.Day(),
		slot.StartTime.Hour(), slot.StartTime.Minute(), 0, 0, slot.Date.Location())

	// Check if "too late" to book (Buffer check)
	// Example: BufferHours = 24. If Now + 24h > SlotStart, then it's too late.
	minBookTime := time.Now().Add(time.Duration(settings.BookingBufferHours) * time.Hour)
	if minBookTime.After(slotStart) {
		return fmt.Errorf("too late to book this slot (buffer %d hours)", settings.BookingBufferHours)
	}

	return nil
}

// Stub for now since we removed dependency
func (v *Validator) ValidateEntry(att interface{}) error {
	return nil
}

// Slot represents a time slot in the schedule (simplified for validation)
type Slot struct {
	ID        string       `json:"id"`
	TeacherID string       `json:"teacher_id"`
	ClassID   string       `json:"class_id"`
	SubjectID string       `json:"subject_id"`
	RoomID    string       `json:"room_id"`
	DayOfWeek time.Weekday `json:"day"`
	Hour      int          `json:"hour"` // 1-8 (e.g., 1st period, 2nd period)
}

// ValidateSchedule checks a list of slots for internal conflicts
func (v *Validator) ValidateSchedule(slots []Slot) []Conflict {
	var conflicts []Conflict

	for i := 0; i < len(slots); i++ {
		for j := i + 1; j < len(slots); j++ {
			a := slots[i]
			b := slots[j]

			if a.DayOfWeek != b.DayOfWeek || a.Hour != b.Hour {
				continue
			}

			if a.TeacherID == b.TeacherID && a.TeacherID != "" {
				conflicts = append(conflicts, Conflict{
					Type:        ConflictTeacherOverlap,
					Description: fmt.Sprintf("Teacher %s is double booked", a.TeacherID),
					SlotA:       a,
					SlotB:       b,
				})
			}

			if a.ClassID == b.ClassID && a.ClassID != "" {
				conflicts = append(conflicts, Conflict{
					Type:        ConflictClassOverlap,
					Description: fmt.Sprintf("Class %s is double booked", a.ClassID),
					SlotA:       a,
					SlotB:       b,
				})
			}

			if a.RoomID == b.RoomID && a.RoomID != "" {
				conflicts = append(conflicts, Conflict{
					Type:        ConflictRoomOverlap,
					Description: fmt.Sprintf("Room %s is double booked", a.RoomID),
					SlotA:       a,
					SlotB:       b,
				})
			}
		}
	}

	return conflicts
}
