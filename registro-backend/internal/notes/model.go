package notes

import (
	"time"
)

type NoteType string

const (
	NoteTypeGeneric      NoteType = "generic"
	NoteTypeHomework     NoteType = "homework"
	NoteTypeBehavior     NoteType = "behavior"
	NoteTypeDisciplinary NoteType = "disciplinary"
)

type StudentNote struct {
	ID        string    `json:"id"`
	SchoolID  string    `json:"school_id"`
	StudentID string    `json:"student_id"`
	TeacherID string    `json:"teacher_id"`
	ClassID   string    `json:"class_id"`
	SubjectID *string   `json:"subject_id,omitempty"` // Nullable
	Type      NoteType  `json:"type"`
	Note      string    `json:"note"`
	Date       string     `json:"date"` // YYYY-MM-DD
	IsApproved bool       `json:"is_approved"`
	ApprovedBy string     `json:"approved_by,omitempty"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	// Joined fields for display
	TeacherName string `json:"teacher_name,omitempty"`
	SubjectName string `json:"subject_name,omitempty"`
}

type CreateNoteRequest struct {
	StudentID string   `json:"student_id" binding:"required"`
	ClassID   string   `json:"class_id" binding:"required"`
	SubjectID *string  `json:"subject_id"`
	Type      NoteType `json:"type" binding:"required,oneof=generic homework behavior disciplinary"`
	Note      string   `json:"note" binding:"required"`
	Date      string   `json:"date" binding:"required"` // YYYY-MM-DD
}

type UpdateNoteRequest struct {
	Type NoteType `json:"type" binding:"omitempty,oneof=generic homework behavior disciplinary"`
	Note string   `json:"note"`
	Date string   `json:"date"`
}

type NoteFilter struct {
	StudentID string
	ClassID   string
	TeacherID string
	Type      NoteType
	DateFrom  string
	DateTo    string
	ActorID   string
	ActorRole string
	Page      int
	Limit     int
}
