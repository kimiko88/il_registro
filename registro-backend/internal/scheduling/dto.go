package scheduling

import "time"

// Requests

type CreateSlotRequest struct {
	Dates          []string `json:"dates" binding:"required"`      // YYYY-MM-DD
	StartTime      string   `json:"start_time" binding:"required"` // HH:MM
	EndTime        string   `json:"end_time" binding:"required"`
	Duration       int      `json:"duration"`     // Minutes per slot, if > 0 split range
	MaxBookings    int      `json:"max_bookings"` // Default 1
	Type           SlotType `json:"type" binding:"required"`
	Location       string   `json:"location"`
	IsRecurring    bool     `json:"is_recurring"`
	RecurringUntil string   `json:"recurring_until"` // YYYY-MM-DD
}

type UpdateSlotRequest struct {
	Location    *string `json:"location,omitempty"`
	IsCancelled *bool   `json:"is_cancelled,omitempty"`
}

type BookSlotRequest struct {
	SlotID    string  `json:"slot_id" binding:"required"`
	StudentID *string `json:"student_id,omitempty"` // Which child
}

type UpdateBookingRequest struct {
	Notes  *string        `json:"notes,omitempty"`
	Status *BookingStatus `json:"status,omitempty"` // For cancellations or completion
}

type GeneralScheduleRequest struct {
	StartDate          string `json:"start_date" binding:"required"`
	EndDate            string `json:"end_date" binding:"required"`
	BookingWindowDays  int    `json:"booking_window_days"`
	BookingBufferHours int    `json:"booking_buffer_hours"`
}

// Responses

type SlotResponse struct {
	ID        string   `json:"id"`
	Date      string   `json:"date"`
	TimeRange string   `json:"time_range"`
	Type      SlotType `json:"type"`
	Available bool     `json:"available"`
	TeacherID string   `json:"teacher_id"`
	Location  string   `json:"location"`
}

type BookingResponse struct {
	ID          string        `json:"id"`
	SlotInfo    SlotResponse  `json:"slot_info"`
	Status      BookingStatus `json:"status"`
	BookedAt    time.Time     `json:"booked_at"`
	Notes       string        `json:"notes"`
	ParentName  string        `json:"parent_name"`
	StudentName string        `json:"student_name"`
}

type AnalyticsResponse struct {
	TotalSlots       int     `json:"total_slots"`
	UtilizationRate  float64 `json:"utilization_rate"`
	CancellationRate float64 `json:"cancellation_rate"`
}
