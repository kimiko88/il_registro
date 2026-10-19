package competencies_invalsi

import "time"

// CompetencyCertificationDM14 implements the unified certification model for Primaria & Primo Ciclo (D.M. 14/2024)
type CompetencyCertificationDM14 struct {
	ID                      string     `json:"id" db:"id"`
	SchoolID                string     `json:"school_id" db:"school_id"`
	StudentID               string     `json:"student_id" db:"student_id"`
	ClassID                 string     `json:"class_id" db:"class_id"`
	SchoolTier              string     `json:"school_tier" db:"school_tier"` // primaria, primo_ciclo
	SchoolYear              string     `json:"school_year" db:"school_year"` // 2025/2026
	RubricEvaluationsJSON   string     `json:"rubric_evaluations_json" db:"rubric_evaluations_json"`
	FinalConsiderations     string     `json:"final_considerations" db:"final_considerations"`
	CouncilDeliberationDate string     `json:"council_deliberation_date" db:"council_deliberation_date"` // YYYY-MM-DD
	Status                  string     `json:"status" db:"status"`                                       // bozza, approvata, firmata_ds, rilasciata
	SignedByPrincipalID     *string    `json:"signed_by_principal_id,omitempty" db:"signed_by_principal_id"`
	SignedAt                *time.Time `json:"signed_at,omitempty" db:"signed_at"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`

	StudentFullName string `json:"student_full_name,omitempty"`
}

// InvalsiResult records standardized test scores and comparison with school report grades
type InvalsiResult struct {
	ID               string    `json:"id" db:"id"`
	SchoolID         string    `json:"school_id" db:"school_id"`
	ClassID          string    `json:"class_id" db:"class_id"`
	StudentID        *string   `json:"student_id,omitempty" db:"student_id"`
	SchoolYear       string    `json:"school_year" db:"school_year"`
	GradeLevel       int       `json:"grade_level" db:"grade_level"` // 8 (3a media), 10 (2a sup), 13 (5a sup)
	Subject          string    `json:"subject" db:"subject"`         // italiano, matematica, inglese_reading, inglese_listening
	ScoreWLE         float64   `json:"score_wle" db:"score_wle"`
	Level            int       `json:"level" db:"level"` // 1 to 5
	SchoolMean       float64   `json:"school_mean" db:"school_mean"`
	NationalMean     float64   `json:"national_mean" db:"national_mean"`
	ReportClassGrade float64   `json:"report_class_grade" db:"report_class_grade"` // Real council grade for RAV/PdM correlation
	CreatedAt        time.Time `json:"created_at" db:"created_at"`

	ScoreDeltaVsNational float64 `json:"score_delta_vs_national"`
}
