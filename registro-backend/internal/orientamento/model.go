package orientamento

import (
	"time"
)

type Event struct {
	ID           string    `json:"id" db:"id"`
	SchoolID     string    `json:"school_id" db:"school_id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Category     string    `json:"category" db:"category"` // University, Work
	Date         time.Time `json:"date" db:"date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
	Location     string    `json:"location" db:"location"`
	Hours        float64   `json:"hours" db:"hours"`
	MaxAttendees int       `json:"max_attendees" db:"max_attendees"`
	CreatedBy    string    `json:"created_by" db:"created_by"`
}

type Participation struct {
	ID           string    `json:"id" db:"id"`
	EventID      string    `json:"event_id" db:"event_id"`
	StudentID    string    `json:"student_id" db:"student_id"`
	Status       string    `json:"status" db:"status"` // Registered, Attended
	Attended     bool      `json:"attended" db:"attended"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
}

// DTOs
type CreateEventRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Category     string  `json:"category" binding:"required"`
	Date         string  `json:"date" binding:"required"`
	EndDate      string  `json:"end_date" binding:"required"`
	Location     string  `json:"location"`
	Hours        float64 `json:"hours" binding:"required"`
	MaxAttendees int     `json:"max_attendees"`
}
