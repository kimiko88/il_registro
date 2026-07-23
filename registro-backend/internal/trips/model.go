package trips

import (
	"time"
)

type EducationalTrip struct {
	ID                   string    `json:"id" db:"id"`
	SchoolID             string    `json:"school_id" db:"school_id"`
	Title                string    `json:"title" db:"title"`
	Destination          string    `json:"destination" db:"destination"`
	DepartureDate        time.Time `json:"departure_date" db:"departure_date"`
	ReturnDate           time.Time `json:"return_date" db:"return_date"`
	Description          string    `json:"description,omitempty" db:"description"`
	AccompanyingTeachers string    `json:"accompanying_teachers,omitempty" db:"accompanying_teachers"`
	ClassIDs             []string  `json:"class_ids" db:"class_ids"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`

	// Status relative to student
	ConsentStatus string `json:"consent_status,omitempty"`
}

type TripConsent struct {
	ID        string    `json:"id" db:"id"`
	TripID    string    `json:"trip_id" db:"trip_id"`
	StudentID string    `json:"student_id" db:"student_id"`
	ParentID  *string   `json:"parent_id,omitempty" db:"parent_id"`
	Status    string    `json:"status" db:"status"` // 'granted', 'denied'
	SignedAt  time.Time `json:"signed_at" db:"signed_at"`
	IPAddress string    `json:"ip_address,omitempty" db:"ip_address"`

	// Joined
	StudentName string `json:"student_name,omitempty"`
	ParentName  string `json:"parent_name,omitempty"`
}

type CreateTripRequest struct {
	Title                string   `json:"title" binding:"required"`
	Destination          string   `json:"destination" binding:"required"`
	DepartureDate        string   `json:"departure_date" binding:"required"` // RFC3339 or YYYY-MM-DD
	ReturnDate           string   `json:"return_date" binding:"required"`    // RFC3339 or YYYY-MM-DD
	Description          string   `json:"description"`
	AccompanyingTeachers string   `json:"accompanying_teachers"`
	ClassIDs             []string `json:"class_ids" binding:"required"`
}

type SubmitConsentRequest struct {
	TripID    string `json:"trip_id" binding:"required"`
	StudentID string `json:"student_id" binding:"required"`
	Status    string `json:"status" binding:"required"` // 'granted', 'denied'
}
