package substitutions

import "time"

type SubstitutionStatus string

const (
	StatusPending   SubstitutionStatus = "pending"
	StatusConfirmed SubstitutionStatus = "confirmed"
	StatusAssigned  SubstitutionStatus = "assigned"
	StatusCancelled SubstitutionStatus = "cancelled"
)

type Substitution struct {
	ID                    string             `json:"id" db:"id"`
	SchoolID              string             `json:"school_id" db:"school_id"`
	ClassID               string             `json:"class_id" db:"class_id"`
	AbsentTeacherID       string             `json:"absent_teacher_id" db:"absent_teacher_id"`
	SubstituteTeacherID   *string            `json:"substitute_teacher_id,omitempty" db:"substitute_teacher_id"`
	Date                  time.Time          `json:"date" db:"date"`
	Slot                  int                `json:"slot" db:"slot"`
	Hour                  int                `json:"hour" db:"hour"`
	Subject               string             `json:"subject,omitempty" db:"subject"`
	SubjectID             string             `json:"subject_id,omitempty" db:"subject_id"`
	Notes                 string             `json:"notes,omitempty" db:"notes"`
	Status                SubstitutionStatus `json:"status" db:"status"`
	SignedBySubstitute    bool               `json:"signed_by_substitute" db:"signed_by_substitute"`
	SignatureTimestamp    *time.Time         `json:"signature_timestamp,omitempty" db:"signature_timestamp"`
	SignatureHash         string             `json:"signature_hash,omitempty" db:"signature_hash"`
	OfficialRegisterNotes string             `json:"official_register_notes,omitempty" db:"official_register_notes"`
	CreatedBy             string             `json:"created_by" db:"created_by"`
	CreatedAt             time.Time          `json:"created_at" db:"created_at"`
}

type CreateSubstitutionRequest struct {
	ClassID             string  `json:"class_id" binding:"required"`
	AbsentTeacherID     string  `json:"absent_teacher_id" binding:"required"`
	SubstituteTeacherID *string `json:"substitute_teacher_id,omitempty"`
	Date                string  `json:"date" binding:"required"` // YYYY-MM-DD
	Slot                int     `json:"slot"`
	Hour                int     `json:"hour"`
	Subject             string  `json:"subject"`
	SubjectID           string  `json:"subject_id"`
	Notes               string  `json:"notes"`
}

type AssignSubstituteRequest struct {
	SubstituteTeacherID string `json:"substitute_teacher_id" binding:"required"`
	Notes               string `json:"notes"`
}

type SubstituteRecommendation struct {
	TeacherID      string `json:"teacher_id"`
	TeacherName    string `json:"teacher_name"`
	Score          int    `json:"score"`
	Reason         string `json:"reason"`
	IsFree         bool   `json:"is_free"`
	TeachesClass   bool   `json:"teaches_class"`
	TeachesSubject bool   `json:"teaches_subject"`
}
