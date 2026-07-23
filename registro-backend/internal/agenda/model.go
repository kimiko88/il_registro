package agenda

import (
	"time"
)

type AgendaType string

const (
	TypeHomework AgendaType = "compito"
	TypeTest     AgendaType = "verifica"
	TypeNotice   AgendaType = "avviso"
	TypeEvent    AgendaType = "evento"
)

type AgendaItem struct {
	ID          string     `json:"id" db:"id"`
	SchoolID    string     `json:"school_id" db:"school_id"`
	ClassID     string     `json:"class_id" db:"class_id"`
	SubjectID   *string    `json:"subject_id,omitempty" db:"subject_id"`
	TeacherID   string     `json:"teacher_id" db:"teacher_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description,omitempty" db:"description"`
	Type        AgendaType `json:"type" db:"type"`
	Date        time.Time  `json:"date" db:"date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// Joined metadata
	SubjectName string `json:"subject_name,omitempty" db:"subject_name"`
	TeacherName string `json:"teacher_name,omitempty" db:"teacher_name"`
	IsCompleted bool   `json:"is_completed" db:"is_completed"`
}

type StudentAgendaCompletion struct {
	ID           string    `json:"id" db:"id"`
	AgendaItemID string    `json:"agenda_item_id" db:"agenda_item_id"`
	StudentID    string    `json:"student_id" db:"student_id"`
	CompletedAt  time.Time `json:"completed_at" db:"completed_at"`
}

type CreateAgendaItemRequest struct {
	ClassID     string     `json:"class_id" binding:"required"`
	SubjectID   *string    `json:"subject_id,omitempty"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Type        AgendaType `json:"type" binding:"required"` // compito, verifica, avviso, evento
	Date        string     `json:"date" binding:"required"` // YYYY-MM-DD
}

type UpdateAgendaItemRequest struct {
	Title       *string     `json:"title,omitempty"`
	Description *string     `json:"description,omitempty"`
	Type        *AgendaType `json:"type,omitempty"`
	Date        *string     `json:"date,omitempty"`
}

type CalendarFilter struct {
	ClassID   string
	StudentID string
	SubjectID string
	Type      string
	From      time.Time
	To        time.Time
}
