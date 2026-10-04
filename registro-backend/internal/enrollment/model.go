package enrollment

import (
	"encoding/json"
	"time"
)

type EnrollmentApplication struct {
	ID                    string    `json:"id" db:"id"`
	SchoolID              string    `json:"school_id" db:"school_id"`
	AcademicYear          string    `json:"academic_year" db:"academic_year"`
	SidiApplicationID     string    `json:"sidi_application_id" db:"sidi_application_id"`
	StudentFirstName      string    `json:"student_first_name" db:"student_first_name"`
	StudentLastName       string    `json:"student_last_name" db:"student_last_name"`
	StudentTaxCode        string    `json:"student_tax_code" db:"student_tax_code"`
	BirthDate             string    `json:"birth_date" db:"birth_date"`
	Gender                string    `json:"gender" db:"gender"` // M or F
	OriginSchool          string    `json:"origin_school" db:"origin_school"`
	MiddleSchoolGrade     int       `json:"middle_school_grade" db:"middle_school_grade"`
	TrackChosen           string    `json:"track_chosen" db:"track_chosen"`
	SecondLanguage        string    `json:"second_language" db:"second_language"`
	HasDisabilityL104     bool      `json:"has_disability_l104" db:"has_disability_l104"`
	HasDSA                bool      `json:"has_dsa" db:"has_dsa"`
	ReligionChoice        string    `json:"religion_choice" db:"religion_choice"`
	RequestedClassmates   []string  `json:"requested_classmates" db:"requested_classmates"`
	IncompatibleClassmates []string `json:"incompatible_classmates" db:"incompatible_classmates"`
	Parent1FirstName      string    `json:"parent1_first_name" db:"parent1_first_name"`
	Parent1LastName       string    `json:"parent1_last_name" db:"parent1_last_name"`
	Parent1Email          string    `json:"parent1_email" db:"parent1_email"`
	Parent1Phone          string    `json:"parent1_phone" db:"parent1_phone"`
	Status                string    `json:"status" db:"status"` // pending, assigned, rejected
	AssignedClassID       string    `json:"assigned_class_id,omitempty" db:"assigned_class_id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

type FormationParams struct {
	TargetClassCount int     `json:"target_class_count"` // e.g. 3 or 4 classes
	MaxStudents      int     `json:"max_students"`       // standard limit 27
	MaxL104PerClass  int     `json:"max_l104_per_class"` // e.g. 1
	BalanceGender    bool    `json:"balance_gender"`
	BalanceGrades    bool    `json:"balance_grades"`
}

type ClassGroupAssignment struct {
	ClassName     string                  `json:"class_name"` // "1A", "1B", etc.
	SecondLanguage string                 `json:"second_language"`
	TotalStudents int                     `json:"total_students"`
	MalesCount    int                     `json:"males_count"`
	FemalesCount  int                     `json:"females_count"`
	L104Count     int                     `json:"l104_count"`
	DSACount      int                     `json:"dsa_count"`
	AverageGrade  float64                 `json:"average_grade"`
	Students      []EnrollmentApplication `json:"students"`
}

type ClassFormationDraftResult struct {
	Classes []ClassGroupAssignment `json:"classes"`
}

type ClassFormationDraft struct {
	ID           string                    `json:"id" db:"id"`
	SchoolID     string                    `json:"school_id" db:"school_id"`
	AcademicYear string                    `json:"academic_year" db:"academic_year"`
	Title        string                    `json:"title" db:"title"`
	Parameters   FormationParams           `json:"parameters"`
	Assignments  ClassFormationDraftResult `json:"assignments"`
	IsFinalized  bool                      `json:"is_finalized" db:"is_finalized"`
	CreatedAt    time.Time                 `json:"created_at" db:"created_at"`
}

type RawDraftRow struct {
	ID          string          `db:"id"`
	SchoolID    string          `db:"school_id"`
	AcademicYear string         `db:"academic_year"`
	Title       string          `db:"title"`
	Parameters  json.RawMessage `db:"parameters"`
	Assignments json.RawMessage `db:"assignments"`
	IsFinalized bool            `db:"is_finalized"`
	CreatedAt   time.Time       `db:"created_at"`
}
