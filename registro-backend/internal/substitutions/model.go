package substitutions

import "time"

type SubstitutionStatus string

const (
	StatusPending   SubstitutionStatus = "pending"
	StatusAssigned  SubstitutionStatus = "assigned"
	StatusCancelled SubstitutionStatus = "cancelled"
)

type Substitution struct {
	ID                 string             `json:"id" db:"id"`
	SchoolID           string             `json:"school_id" db:"school_id"`
	ClassID            string             `json:"class_id" db:"class_id"`
	AbsentTeacherID    string             `json:"absent_teacher_id" db:"absent_teacher_id"`
	SubstituteTeacherID *string           `json:"substitute_teacher_id,omitempty" db:"substitute_teacher_id"`
	Date               time.Time          `json:"date" db:"date"`
	Hour               int                `json:"hour" db:"hour"`
	SubjectID          string             `json:"subject_id,omitempty" db:"subject_id"`
	Notes              string             `json:"notes,omitempty" db:"notes"`
	Status             SubstitutionStatus `json:"status" db:"status"`
	CreatedAt          time.Time          `json:"created_at" db:"created_at"`
}

type CreateSubstitutionRequest struct {
	ClassID             string `json:"class_id" binding:"required"`
	AbsentTeacherID    string `json:"absent_teacher_id" binding:"required"`
	SubstituteTeacherID *string`json:"substitute_teacher_id,omitempty"`
	Date               string `json:"date" binding:"required"` // YYYY-MM-DD
	Hour               int    `json:"hour" binding:"required"`
	SubjectID          string `json:"subject_id"`
	Notes              string `json:"notes"`
}

type AssignSubstituteRequest struct {
	SubstituteTeacherID string `json:"substitute_teacher_id" binding:"required"`
	Notes               string `json:"notes"`
}
