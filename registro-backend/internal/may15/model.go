package may15

import "time"

type DocumentStatus string

const (
	StatusBozza        DocumentStatus = "bozza"
	StatusApprovatoCdC DocumentStatus = "approvato_cdc"
	StatusPubblicato   DocumentStatus = "pubblicato"
)

type ClassMay15Document struct {
	ID                 string         `json:"id" db:"id"`
	SchoolID           string         `json:"school_id" db:"school_id"`
	ClassID            string         `json:"class_id" db:"class_id"`
	ClassName          string         `json:"class_name,omitempty" db:"class_name"`
	AcademicYear       string         `json:"academic_year" db:"academic_year"`
	Status             DocumentStatus `json:"status" db:"status"`
	ClassPresentation  string         `json:"class_presentation" db:"class_presentation"`
	TeachingContinuity string         `json:"teaching_continuity" db:"teaching_continuity"`
	PCTOPathways       string         `json:"pcto_pathways" db:"pcto_pathways"`
	ExamSimulations    string         `json:"exam_simulations" db:"exam_simulations"`
	EvaluationRubrics  string         `json:"evaluation_rubrics" db:"evaluation_rubrics"`
	CLILModules        string         `json:"clil_modules" db:"clil_modules"`
	ApprovedAt         *time.Time     `json:"approved_at,omitempty" db:"approved_at"`
	PublishedAt        *time.Time     `json:"published_at,omitempty" db:"published_at"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
}

type SaveMay15Request struct {
	AcademicYear       string         `json:"academic_year" binding:"required"`
	Status             DocumentStatus `json:"status"`
	ClassPresentation  string         `json:"class_presentation"`
	TeachingContinuity string         `json:"teaching_continuity"`
	PCTOPathways       string         `json:"pcto_pathways"`
	ExamSimulations    string         `json:"exam_simulations"`
	EvaluationRubrics  string         `json:"evaluation_rubrics"`
	CLILModules        string         `json:"clil_modules"`
}
