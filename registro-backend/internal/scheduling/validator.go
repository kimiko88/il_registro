package scheduling

import (
	"errors"
	"time"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateSlot(s *ColloquioSlot) error {
	if s.StartTime.After(s.EndTime) {
		return errors.New("start time must be before end time")
	}
	if s.Date.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("cannot create slots in the past")
	}
	return nil
}

func (v *Validator) ValidateBooking(slot *ColloquioSlot, settings *ColloquioSettings) error {
	now := time.Now()
	// Booking Buffer (e.g., must book 24h in advance)
	buffer := time.Duration(settings.BookingBufferHours) * time.Hour
	slotTime := time.Date(slot.Date.Year(), slot.Date.Month(), slot.Date.Day(),
		slot.StartTime.Hour(), slot.StartTime.Minute(), 0, 0, slot.Date.Location())

	if now.Add(buffer).After(slotTime) {
		return errors.New("too late to book this slot")
	}

	// Window (e.g., can only book 14 days ahead)
	windowEnd := now.AddDate(0, 0, settings.BookingWindowDays)
	if slot.Date.After(windowEnd) {
		return errors.New("slot is too far in the future")
	}

	return nil
}
