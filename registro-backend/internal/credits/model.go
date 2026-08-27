package credits

import "time"

type StudentSchoolCredit struct {
	ID                 string    `json:"id"`
	SchoolID           string    `json:"school_id"`
	StudentID          string    `json:"student_id"`
	StudentName        string    `json:"student_name,omitempty"`
	ClassID            string    `json:"class_id"`
	ClassName          string    `json:"class_name,omitempty"`
	AcademicYear       string    `json:"academic_year"`
	GradeLevel         int       `json:"grade_level"` // 3, 4, 5
	GradeAverage       float64   `json:"grade_average"`
	ConductGrade       int       `json:"conduct_grade"`
	BaseCreditRangeMin int       `json:"base_credit_range_min"`
	BaseCreditRangeMax int       `json:"base_credit_range_max"`
	AssignedCredit     int       `json:"assigned_credit"`
	PCTOHours          int       `json:"pcto_hours"`
	HasExtracurricular bool      `json:"has_extracurricular"`
	DeliberationNotes  string    `json:"deliberation_notes"`
	ValidatedBy        *string   `json:"validated_by,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreditCalculationResult struct {
	GradeLevel         int     `json:"grade_level"`
	GradeAverage       float64 `json:"grade_average"`
	ConductGrade       int     `json:"conduct_grade"`
	BaseCreditRangeMin int     `json:"base_credit_range_min"`
	BaseCreditRangeMax int     `json:"base_credit_range_max"`
	SuggestedCredit    int     `json:"suggested_credit"`
	MaxYearCredit      int     `json:"max_year_credit"`
	Motivation         string  `json:"motivation"`
}

type StudentCreditSummary struct {
	StudentID            string                `json:"student_id"`
	StudentName          string                `json:"student_name"`
	ClassID              string                `json:"class_id"`
	ClassName            string                `json:"class_name"`
	CreditsByYear        []StudentSchoolCredit `json:"credits_by_year"`
	TotalTrienniumCredit int                   `json:"total_triennium_credit"` // Max 40
	MaxPossibleCredit    int                   `json:"max_possible_credit"`
}

type AssignCreditRequest struct {
	StudentID          string  `json:"student_id" binding:"required"`
	ClassID            string  `json:"class_id" binding:"required"`
	AcademicYear       string  `json:"academic_year" binding:"required"`
	GradeLevel         int     `json:"grade_level" binding:"required"`
	GradeAverage       float64 `json:"grade_average"`
	ConductGrade       int     `json:"conduct_grade"`
	AssignedCredit     int     `json:"assigned_credit" binding:"required"`
	PCTOHours          int     `json:"pcto_hours"`
	HasExtracurricular bool    `json:"has_extracurricular"`
	DeliberationNotes  string  `json:"deliberation_notes"`
}
