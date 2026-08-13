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
	ID               string     `json:"id"`
	SchoolID         string     `json:"school_id"`
	StudentID        string     `json:"student_id"`
	TeacherID        string     `json:"teacher_id"`
	ClassID          string     `json:"class_id"`
	SubjectID        *string    `json:"subject_id,omitempty"` // Nullable
	Type             NoteType   `json:"type"`
	Note             string     `json:"note"`
	Date             string     `json:"date"` // YYYY-MM-DD
	IsReserved       bool       `json:"is_reserved"`
	TargetRole       string     `json:"target_role"` // "coordinator" | "admin" | "all"
	IsApproved       bool       `json:"is_approved"`
	ApprovedBy       string     `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	DeletionReason   string     `json:"deletion_reason,omitempty"`
	IsViewedByParent bool       `json:"is_viewed_by_parent"`
	ParentViewedAt   *time.Time `json:"parent_viewed_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Joined fields for display
	TeacherName string `json:"teacher_name,omitempty"`
	SubjectName string `json:"subject_name,omitempty"`
}

type CreateNoteRequest struct {
	StudentID  string   `json:"student_id" binding:"required"`
	ClassID    string   `json:"class_id" binding:"required"`
	SubjectID  *string  `json:"subject_id"`
	Type       NoteType `json:"type" binding:"required,oneof=generic homework behavior disciplinary"`
	Note       string   `json:"note" binding:"required"`
	Date       string   `json:"date" binding:"required"` // YYYY-MM-DD
	IsReserved bool     `json:"is_reserved"`
	TargetRole string   `json:"target_role"`
}

type UpdateNoteRequest struct {
	Type       NoteType `json:"type" binding:"omitempty,oneof=generic homework behavior disciplinary"`
	Note       string   `json:"note"`
	Date       string   `json:"date"`
	IsReserved *bool    `json:"is_reserved,omitempty"`
	TargetRole string   `json:"target_role,omitempty"`
}

type DeleteNoteRequest struct {
	Reason string `json:"reason"`
}

type NoteFilter struct {
	SchoolID      string
	StudentID     string
	ClassID       string
	TeacherID     string
	Type          NoteType
	DateFrom      string
	DateTo        string
	ActorID       string
	ActorRole     string
	IsCoordinator bool
	Page          int
	Limit         int
}
