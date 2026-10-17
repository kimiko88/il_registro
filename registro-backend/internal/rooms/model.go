package rooms

import (
	"encoding/json"
	"time"
)

// RoomType constants
const (
	RoomTypeClassroom      = "classroom"
	RoomTypeLabInformatica = "lab_informatica"
	RoomTypeLabScienze     = "lab_scienze"
	RoomTypeLabLingue      = "lab_lingue"
	RoomTypeLabChimica     = "lab_chimica"
	RoomTypeLabFisica      = "lab_fisica"
	RoomTypeLabArte        = "lab_arte"
	RoomTypePalestra       = "palestra"
	RoomTypeAulaMagna      = "aula_magna"
	RoomTypeBiblioteca     = "biblioteca"
	RoomTypeAltro          = "altro"
)

// BookingStatus constants
const (
	BookingStatusConfirmed = "confirmed"
	BookingStatusCancelled = "cancelled"
)

// RecurrencePattern constants
const (
	RecurrenceWeekly   = "weekly"
	RecurrenceBiweekly = "biweekly"
)

// SchoolBuilding represents a school facility/building (plesso)
type SchoolBuilding struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Name      string    `json:"name" db:"name"`
	Address   *string   `json:"address,omitempty" db:"address"`
	Notes     *string   `json:"notes,omitempty" db:"notes"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// BookableRoom represents a room that can be booked or assigned
type BookableRoom struct {
	ID              string          `json:"id" db:"id"`
	SchoolID        string          `json:"school_id" db:"school_id"`
	BuildingID      *string         `json:"building_id,omitempty" db:"building_id"`
	BuildingName    string          `json:"building_name,omitempty" db:"building_name"`
	Name            string          `json:"name" db:"name"`
	RoomType        string          `json:"room_type" db:"room_type"`
	Capacity        int             `json:"capacity" db:"capacity"`
	Equipment       json.RawMessage `json:"equipment" db:"equipment"`
	RequiresBooking bool            `json:"requires_booking" db:"requires_booking"`
	IsActive        bool            `json:"is_active" db:"is_active"`
	Notes           *string         `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// RoomBooking represents a reservation of a bookable room
type RoomBooking struct {
	ID                string    `json:"id" db:"id"`
	SchoolID          string    `json:"school_id" db:"school_id"`
	RoomID            string    `json:"room_id" db:"room_id"`
	RoomName          string    `json:"room_name,omitempty" db:"room_name"`
	BuildingID        *string   `json:"building_id,omitempty" db:"building_id"`
	BuildingName      string    `json:"building_name,omitempty" db:"building_name"`
	TeacherID         string    `json:"teacher_id" db:"teacher_id"`
	TeacherName       string    `json:"teacher_name,omitempty" db:"teacher_name"`
	ClassID           *string   `json:"class_id,omitempty" db:"class_id"`
	ClassName         string    `json:"class_name,omitempty" db:"class_name"`
	SubjectID         *string   `json:"subject_id,omitempty" db:"subject_id"`
	SubjectName       string    `json:"subject_name,omitempty" db:"subject_name"`
	BookingDate       string    `json:"booking_date" db:"booking_date"` // YYYY-MM-DD
	HourIndex         int       `json:"hour_index" db:"hour_index"`
	Status            string    `json:"status" db:"status"`
	Notes             *string   `json:"notes,omitempty" db:"notes"`
	IsRecurring       bool      `json:"is_recurring" db:"is_recurring"`
	RecurrencePattern *string   `json:"recurrence_pattern,omitempty" db:"recurrence_pattern"`
	RecurringUntil    *string   `json:"recurring_until,omitempty" db:"recurring_until"`
	ParentBookingID   *string   `json:"parent_booking_id,omitempty" db:"parent_booking_id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Request and Response DTOs
type CreateBuildingRequest struct {
	Name    string  `json:"name" binding:"required"`
	Address *string `json:"address"`
	Notes   *string `json:"notes"`
}

type UpdateBuildingRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Notes    *string `json:"notes"`
	IsActive *bool   `json:"is_active"`
}

type CreateRoomRequest struct {
	BuildingID      *string  `json:"building_id"`
	Name            string   `json:"name" binding:"required"`
	RoomType        string   `json:"room_type"`
	Capacity        int      `json:"capacity"`
	Equipment       []string `json:"equipment"`
	RequiresBooking *bool    `json:"requires_booking"`
	Notes           *string  `json:"notes"`
}

type UpdateRoomRequest struct {
	BuildingID      *string  `json:"building_id"`
	Name            *string  `json:"name"`
	RoomType        *string  `json:"room_type"`
	Capacity        *int     `json:"capacity"`
	Equipment       []string `json:"equipment"`
	RequiresBooking *bool    `json:"requires_booking"`
	IsActive        *bool    `json:"is_active"`
	Notes           *string  `json:"notes"`
}

type CreateBookingRequest struct {
	RoomID            string  `json:"room_id" binding:"required"`
	ClassID           *string `json:"class_id"`
	SubjectID         *string `json:"subject_id"`
	BookingDate       string  `json:"booking_date" binding:"required"` // YYYY-MM-DD
	HourIndex         int     `json:"hour_index" binding:"required,min=1,max=12"`
	Notes             *string `json:"notes"`
	IsRecurring       bool    `json:"is_recurring"`
	RecurrencePattern *string `json:"recurrence_pattern"` // weekly, biweekly
	RecurringUntil    *string `json:"recurring_until"`    // YYYY-MM-DD
}

type BookingFilter struct {
	SchoolID        string
	RoomID          *string
	BuildingID      *string
	TeacherID       *string
	ClassID         *string
	FromDate        *string
	ToDate          *string
	Status          *string
	ParentBookingID *string
}

type RoomAvailabilitySlot struct {
	HourIndex   int          `json:"hour_index"`
	IsAvailable bool         `json:"is_available"`
	Booking     *RoomBooking `json:"booking,omitempty"`
}

type DayAvailability struct {
	Date  string                 `json:"date"`
	Slots []RoomAvailabilitySlot `json:"slots"`
}
