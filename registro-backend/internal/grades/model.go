package grades

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// GradeType defines the type of evaluation
type GradeType string

const (
	GradeTypeNumeric    GradeType = "numeric"    // 0-10
	GradeTypeJudgment   GradeType = "judgment"   // Insufficiente, Buono, etc.
	GradeTypeCredit     GradeType = "credit"     // CFU, PCTO
	GradeTypeCompetence GradeType = "competence" // Livello
)

// GradeCategory defines the purpose of the evaluation
type GradeCategory string

const (
	GradeCategoryFormative GradeCategory = "formative" // Valutazione formativa (in itinere)
	GradeCategorySummative GradeCategory = "summative" // Valutazione sommativa (verifiche finali)
	GradeCategoryPractical GradeCategory = "practical" // Laboratorio/Pratica
)

// Semester represents the academic period
type Semester int

const (
	SemesterFirst  Semester = 1
	SemesterSecond Semester = 2
	SemesterThird  Semester = 3 // Rare, but supported for trimesters
)

// EvaluationType defines the type of test
type EvaluationType string

const (
	EvaluationTypeWritten   EvaluationType = "Written"
	EvaluationTypeOral      EvaluationType = "Oral"
	EvaluationTypePractical EvaluationType = "Practical"
)

// Grade represents a single evaluation entry in the Italian school context
type Grade struct {
	ID string `json:"id" db:"id"`

	StudentID string `json:"student_id" db:"student_id"`
	SchoolID  string `json:"school_id" db:"school_id"`
	SubjectID string `json:"subject_id" db:"subject_id"`
	TeacherID string `json:"teacher_id" db:"teacher_id"`

	// GradeValue stores the numeric representation.
	// For judgments, this maps to an internal scale (e.g., Ottimo -> 10).
	GradeValue float64   `json:"grade_value" db:"grade_value"`
	GradeType  GradeType `json:"grade_type" db:"grade_type"`

	Semester Semester  `json:"semester" db:"semester"`
	Date     time.Time `json:"date" db:"date"`

	Description string  `json:"description" db:"description"`
	RubricID    *string `json:"rubric_id,omitempty" db:"rubric_id"` // Optional: link to specific rubric
	Weight      float64 `json:"weight" db:"weight"`                 // Default 1.0

	IsPublished bool       `json:"is_published" db:"is_published"`
	PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at"`

	GradeCategory  GradeCategory   `json:"grade_category" db:"grade_category"`
	EvaluationType *EvaluationType `json:"evaluation_type,omitempty" db:"evaluation_type"`

	// Audit fields
	CreatedBy string     `json:"created_by" db:"created_by"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	// ModifiedBy tracks the last editor
	ModifiedBy *string `json:"modified_by,omitempty" db:"modified_by"`

	TestID *string `json:"test_id,omitempty" db:"test_id"`
}

// GradeHistory maintains an audit trail of changes to grades
type GradeHistory struct {
	ID      string `json:"id" db:"id"`
	GradeID string `json:"grade_id" db:"grade_id"`

	OldValue *float64 `json:"old_value" db:"old_value"`
	NewValue *float64 `json:"new_value" db:"new_value"`

	OldDescription string `json:"old_description" db:"old_description"`
	NewDescription string `json:"new_description" db:"new_description"`

	ModifiedBy string    `json:"modified_by" db:"modified_by"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
	Reason     string    `json:"reason" db:"reason"`
}

// GradeScale defines the configuration for grades in a specific school context
// (e.g., Elementary vs High School scales)
type GradeScale struct {
	ID       string `json:"id" db:"id"`
	SchoolID string `json:"school_id" db:"school_id"`

	GradeType GradeType `json:"grade_type" db:"grade_type"`
	MinValue  float64   `json:"min_value" db:"min_value"` // e.g., 0 or 1
	MaxValue  float64   `json:"max_value" db:"max_value"` // e.g., 10 or 100

	// Labels stores the judgment mappings e.g. [{"value": 10, "label": "Ottimo"}]
	Labels json.RawMessage `json:"labels" db:"labels"`

	DefaultWeight float64 `json:"default_weight" db:"default_weight"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Subject represents the "Materia" in valid context for grading.
// This struct combines data from 'subjects' and 'class_subjects' tables.
type Subject struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	SchoolYearID string `json:"school_year_id" db:"school_year_id"`
	ClassID      string `json:"class_id" db:"class_id"`
	TeacherID    string `json:"teacher_id" db:"teacher_id"`

	Hours int `json:"hours" db:"hours"`

	// Flags for valid grade types
	IsVotable   bool `json:"is_votable" db:"is_votable"`     // Accepts numeric grades 0-10
	IsJudgeable bool `json:"is_judgeable" db:"is_judgeable"` // Accepts judgments (Ottimo, etc.)
	IsCredit    bool `json:"is_credit" db:"is_credit"`       // Accepts CFU/PCTO credits

	DefaultGradeType GradeType `json:"default_grade_type" db:"default_grade_type"`
}

// IsValid checks basic validity of the Grade
func (g *Grade) IsValid() bool {
	if g.ID == "" && g.SubjectID == "" { // ID check might be loose strictly for validation before save
		return false
	}
	if g.SubjectID == "" || g.TeacherID == "" || g.StudentID == "" {
		return false
	}
	if g.GradeValue < 0 {
		return false
	}
	// UUID validation
	if _, err := uuid.Parse(g.StudentID); err != nil {
		return false
	}
	return true
}

// ClassTest represents an assessment scheduled by a teacher
type ClassTest struct {
	ID             string    `json:"id" db:"id"`
	ClassID        string    `json:"class_id" db:"class_id"`
	SubjectID      string    `json:"subject_id" db:"subject_id"`
	TeacherID      string    `json:"teacher_id" db:"teacher_id"`
	Title          string    `json:"title" db:"title"`
	Date           time.Time `json:"date" db:"date"`
	TeacherNotes   string    `json:"teacher_notes" db:"teacher_notes"`
	ParentNotes    string    `json:"parent_notes" db:"parent_notes"`
	EvaluationType string    `json:"evaluation_type" db:"evaluation_type"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// GradeWeightConfig stores the configurable weight multiplier for a grade category/evaluation type
// scoped to a school, optionally a subject and class.
type GradeWeightConfig struct {
	ID             string  `json:"id" db:"id"`
	SchoolID       string  `json:"school_id" db:"school_id"`
	SubjectID      *string `json:"subject_id,omitempty" db:"subject_id"`
	ClassID        *string `json:"class_id,omitempty" db:"class_id"`
	GradeCategory  string  `json:"grade_category" db:"grade_category"`
	EvaluationType *string `json:"evaluation_type,omitempty" db:"evaluation_type"`
	Weight         float64 `json:"weight" db:"weight"`
	CreatedBy      string  `json:"created_by" db:"created_by"`
}

// UpsertWeightConfigRequest is the request body for creating/updating a weight config.
type UpsertWeightConfigRequest struct {
	SubjectID      *string `json:"subject_id"`
	ClassID        *string `json:"class_id"`
	GradeCategory  string  `json:"grade_category" binding:"required"`
	EvaluationType *string `json:"evaluation_type"`
	Weight         float64 `json:"weight" binding:"required,min=0,max=10"`
}
