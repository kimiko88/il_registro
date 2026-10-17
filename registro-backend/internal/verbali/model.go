package verbali

import (
	"time"
)

type CouncilMeeting struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	ClassID     string    `json:"class_id,omitempty" db:"class_id"`
	MeetingType string    `json:"meeting_type" db:"meeting_type"`
	Title       string    `json:"title" db:"title"`
	Date        time.Time `json:"date" db:"date"`
	StartTime   string    `json:"start_time" db:"start_time"`
	EndTime     string    `json:"end_time" db:"end_time"`
	Agenda      string    `json:"agenda,omitempty" db:"agenda"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type MeetingVerbale struct {
	ID           string     `json:"id" db:"id"`
	MeetingID    string     `json:"meeting_id" db:"meeting_id"`
	Title        string     `json:"title" db:"title"`
	Content      string     `json:"content" db:"content"`
	SecretaryID  *string    `json:"secretary_id,omitempty" db:"secretary_id"`
	PresidentID  *string    `json:"president_id,omitempty" db:"president_id"`
	IsPublished  bool       `json:"is_published" db:"is_published"`
	IsSigned     bool       `json:"is_signed" db:"is_signed"`
	SignedAt     *time.Time `json:"signed_at,omitempty" db:"signed_at"`
	Status       string     `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	IsSignedByMe bool       `json:"is_signed_by_me"`
	CanEdit      bool       `json:"can_edit"`
}

type VerbaleSignature struct {
	ID        string    `json:"id" db:"id"`
	VerbaleID string    `json:"verbale_id" db:"verbale_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	UserName  string    `json:"user_name,omitempty"`
	SignedAt  time.Time `json:"signed_at" db:"signed_at"`
	IPAddress string    `json:"ip_address,omitempty" db:"ip_address"`
}

type MeetingVerbaleTemplate struct {
	ID              string    `json:"id" db:"id"`
	SchoolID        string    `json:"school_id" db:"school_id"`
	Title           string    `json:"title" db:"title"`
	MeetingType     string    `json:"meeting_type" db:"meeting_type"`
	Description     string    `json:"description" db:"description"`
	DefaultAgenda   string    `json:"default_agenda" db:"default_agenda"`
	TemplateContent string    `json:"template_content" db:"template_content"`
	CreatedBy       string    `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type CreateMeetingRequest struct {
	ClassID     string `json:"class_id"`
	MeetingType string `json:"meeting_type"`
	Title       string `json:"title" binding:"required"`
	Date        string `json:"date" binding:"required"`       // YYYY-MM-DD
	StartTime   string `json:"start_time" binding:"required"` // HH:MM
	EndTime     string `json:"end_time" binding:"required"`   // HH:MM
	Agenda      string `json:"agenda"`
}

type CreateVerbaleRequest struct {
	MeetingID   string  `json:"meeting_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Content     string  `json:"content" binding:"required"`
	SecretaryID *string `json:"secretary_id,omitempty"`
	PresidentID *string `json:"president_id,omitempty"`
	IsPublished bool    `json:"is_published"`
}

type UpdateVerbaleRequest struct {
	Title       string  `json:"title" binding:"required"`
	Content     string  `json:"content" binding:"required"`
	SecretaryID *string `json:"secretary_id,omitempty"`
	PresidentID *string `json:"president_id,omitempty"`
	IsPublished bool    `json:"is_published"`
}

type CreateTemplateRequest struct {
	Title           string `json:"title" binding:"required"`
	MeetingType     string `json:"meeting_type" binding:"required"`
	Description     string `json:"description"`
	DefaultAgenda   string `json:"default_agenda" binding:"required"`
	TemplateContent string `json:"template_content" binding:"required"`
}

type UpdateTemplateRequest struct {
	Title           string `json:"title" binding:"required"`
	MeetingType     string `json:"meeting_type" binding:"required"`
	Description     string `json:"description"`
	DefaultAgenda   string `json:"default_agenda" binding:"required"`
	TemplateContent string `json:"template_content" binding:"required"`
}
