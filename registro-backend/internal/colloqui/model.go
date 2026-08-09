package colloqui

import (
	"time"
)

type ColloquioType string

const (
	TypeIndividual ColloquioType = "Individual"
	TypeGeneral    ColloquioType = "General"
	TypeAssembly   ColloquioType = "Assembly"
)

type BookingStatus string

const (
	StatusConfirmed BookingStatus = "Confirmed"
	StatusCancelled BookingStatus = "Cancelled"
	StatusCompleted BookingStatus = "Completed"
)

type ColloquioSlot struct {
	ID           string        `json:"id" db:"id"`
	TeacherID    string        `json:"teacher_id" db:"teacher_id"`
	SchoolID     string        `json:"school_id" db:"school_id"`
	Date         time.Time     `json:"date" db:"date"`
	StartTime    string        `json:"start_time" db:"start_time"`
	EndTime      string        `json:"end_time" db:"end_time"`
	MaxBookings  int           `json:"max_bookings" db:"max_bookings"`
	BookingCount int           `json:"booking_count" db:"booking_count"`
	Type         ColloquioType `json:"type" db:"type"`
	Location     string        `json:"location,omitempty" db:"location"`
	IsCancelled  bool          `json:"is_cancelled" db:"is_cancelled"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`

	// Metadata
	TeacherName string `json:"teacher_name,omitempty" db:"teacher_name"`
}

type ColloquioBooking struct {
	ID        string        `json:"id" db:"id"`
	SlotID    string        `json:"slot_id" db:"slot_id"`
	ParentID  *string       `json:"parent_id,omitempty" db:"parent_id"`
	StudentID *string       `json:"student_id,omitempty" db:"student_id"`
	Status    BookingStatus `json:"status" db:"status"`
	Notes     string        `json:"notes,omitempty" db:"notes"`
	BookedAt  time.Time     `json:"booked_at" db:"booked_at"`

	// Joined metadata
	Slot        *ColloquioSlot `json:"slot,omitempty"`
	ParentName  string         `json:"parent_name,omitempty"`
	StudentName string         `json:"student_name,omitempty"`
}

type CreateSlotRequest struct {
	Date        string        `json:"date" binding:"required"`       // YYYY-MM-DD
	StartTime   string        `json:"start_time" binding:"required"` // HH:MM
	EndTime     string        `json:"end_time" binding:"required"`   // HH:MM
	MaxBookings int           `json:"max_bookings"`
	Type        ColloquioType `json:"type"` // Individual, General, Assembly
	Location    string        `json:"location"`
}

type CreateBookingRequest struct {
	SlotID    string  `json:"slot_id" binding:"required"`
	StudentID *string `json:"student_id,omitempty"`
	Notes     string  `json:"notes"`
}

type UpdateBookingStatusRequest struct {
	Status BookingStatus `json:"status" binding:"required"` // Confirmed, Cancelled, Completed
	Reason string        `json:"reason"`
}

type SlotFilter struct {
	TeacherID string
	SchoolID  string
	From      time.Time
	To        time.Time
	Available bool
}

type GeneralAssembly struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	ClassID     string    `json:"class_id" db:"class_id"`
	Title       string    `json:"title" db:"title"`
	Date        time.Time `json:"date" db:"date"`
	StartTime   string    `json:"start_time" db:"start_time"`
	EndTime     string    `json:"end_time" db:"end_time"`
	Location    string    `json:"location" db:"location"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateAssemblyRequest struct {
	ClassID     string `json:"class_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Date        string `json:"date" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	MaxBookings int    `json:"max_bookings" binding:"required,gt=0"`
	Location    string `json:"location"`
	Description string `json:"description"`
}
