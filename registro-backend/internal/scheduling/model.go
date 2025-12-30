package scheduling

import (
	"time"
)

// Enums
type SlotType string

const (
	SlotIndividual SlotType = "Individual"
	SlotGeneral    SlotType = "General"
	SlotAssembly   SlotType = "Assembly"
)

type BookingStatus string

const (
	StatusConfirmed BookingStatus = "Confirmed"
	StatusCancelled BookingStatus = "Cancelled" // By Parent or Teacher
	StatusNoShow    BookingStatus = "NoShow"
	StatusCompleted BookingStatus = "Completed"
)

// ColloquioSlot represents a teacher's availability
type ColloquioSlot struct {
	ID           string    `json:"id" db:"id"`
	TeacherID    string    `json:"teacher_id" db:"teacher_id"`
	SchoolID     string    `json:"school_id" db:"school_id"`
	Date         time.Time `json:"date" db:"date"`             // stored as DATE in DB, Time here
	StartTime    time.Time `json:"start_time" db:"start_time"` // Time part only usually, but Go uses Time
	EndTime      time.Time `json:"end_time" db:"end_time"`
	MaxBookings  int       `json:"max_bookings" db:"max_bookings"`
	BookingCount int       `json:"booking_count" db:"booking_count"`
	Type         SlotType  `json:"type" db:"type"`
	Location     string    `json:"location" db:"location"`
	IsCancelled  bool      `json:"is_cancelled" db:"is_cancelled"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ColloquioBooking represents a reservation
type ColloquioBooking struct {
	ID        string        `json:"id" db:"id"`
	SlotID    string        `json:"slot_id" db:"slot_id"`
	ParentID  *string       `json:"parent_id,omitempty" db:"parent_id"`
	StudentID *string       `json:"student_id,omitempty" db:"student_id"`
	Status    BookingStatus `json:"status" db:"status"`
	Notes     string        `json:"notes" db:"notes"`
	BookedAt  time.Time     `json:"booked_at" db:"booked_at"`
}

// ColloquioSettings global config
type ColloquioSettings struct {
	SchoolID           string     `json:"school_id" db:"school_id"`
	BookingWindowDays  int        `json:"booking_window_days" db:"booking_window_days"`
	BookingBufferHours int        `json:"booking_buffer_hours" db:"booking_buffer_hours"`
	GeneralWindowStart *time.Time `json:"general_window_start" db:"general_colloqui_window_start"`
	GeneralWindowEnd   *time.Time `json:"general_window_end" db:"general_colloqui_window_end"`
}

// Notification log
type Notification struct {
	ID             string    `json:"id" db:"id"`
	BookingID      string    `json:"booking_id" db:"booking_id"`
	Type           string    `json:"type" db:"type"`
	RecipientEmail string    `json:"recipient_email" db:"recipient_email"`
	SentAt         time.Time `json:"sent_at" db:"sent_at"`
	Status         string    `json:"status" db:"status"`
}
