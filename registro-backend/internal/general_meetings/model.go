package general_meetings

import (
	"time"
)

type GeneralMeeting struct {
	ID                   string     `json:"id" db:"id"`
	SchoolID             string     `json:"school_id" db:"school_id"`
	Title                string     `json:"title" db:"title"`
	Description          string     `json:"description,omitempty" db:"description"`
	Location             string     `json:"location,omitempty" db:"location"`
	MeetingDate          time.Time  `json:"meeting_date" db:"meeting_date"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty" db:"registration_deadline"`
	MaxParticipants      *int       `json:"max_participants,omitempty" db:"max_participants"`
	IsMandatory          bool       `json:"is_mandatory" db:"is_mandatory"`
	TargetRoles          []string   `json:"target_roles" db:"target_roles"`
	CreatedBy            string     `json:"created_by" db:"created_by"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`

	// Derived fields
	RegistrationsCount int  `json:"registrations_count"`
	IsRegistered       bool `json:"is_registered"`
}

type GeneralMeetingRegistration struct {
	ID           string    `json:"id" db:"id"`
	MeetingID    string    `json:"meeting_id" db:"meeting_id"`
	UserID       string    `json:"user_id" db:"user_id"`
	UserName     string    `json:"user_name,omitempty"`
	UserEmail    string    `json:"user_email,omitempty"`
	UserRole     string    `json:"user_role,omitempty"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
}

type CreateGeneralMeetingRequest struct {
	Title                string   `json:"title" binding:"required"`
	Description          string   `json:"description"`
	Location             string   `json:"location"`
	MeetingDate          string   `json:"meeting_date" binding:"required"` // ISO string
	RegistrationDeadline *string  `json:"registration_deadline"`
	MaxParticipants      *int     `json:"max_participants"`
	IsMandatory          bool     `json:"is_mandatory"`
	TargetRoles          []string `json:"target_roles"`
}
