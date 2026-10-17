package primaryeval

import "time"

// MinisterialLevels according to O.M. 172/2020
const (
	LevelAvanzato           = "avanzato"
	LevelIntermedio         = "intermedio"
	LevelBase               = "base"
	LevelInViaPrimaAcquisiz = "in_via_di_prima_acquisizione"
)

// LearningObjective represents a ministerial learning objective for primary school
type LearningObjective struct {
	ID           string    `json:"id" db:"id"`
	SchoolID     string    `json:"school_id" db:"school_id"`
	ClassID      *string   `json:"class_id,omitempty" db:"class_id"`
	SubjectID    string    `json:"subject_id" db:"subject_id"`
	SubjectName  string    `json:"subject_name,omitempty" db:"subject_name"`
	YearGrade    int       `json:"year_grade" db:"year_grade"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	AcademicYear string    `json:"academic_year" db:"academic_year"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type CreateObjectiveRequest struct {
	ClassID      *string `json:"class_id"`
	SubjectID    string  `json:"subject_id" binding:"required"`
	YearGrade    int     `json:"year_grade" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	AcademicYear string  `json:"academic_year" binding:"required"`
}

// PrimaryEvaluation represents an evaluation on a single learning objective
type PrimaryEvaluation struct {
	ID                   string    `json:"id" db:"id"`
	SchoolID             string    `json:"school_id" db:"school_id"`
	StudentID            string    `json:"student_id" db:"student_id"`
	StudentName          string    `json:"student_name,omitempty" db:"student_name"`
	ClassID              string    `json:"class_id" db:"class_id"`
	SubjectID            string    `json:"subject_id" db:"subject_id"`
	TeacherID            string    `json:"teacher_id" db:"teacher_id"`
	ObjectiveID          string    `json:"objective_id" db:"objective_id"`
	ObjectiveTitle       string    `json:"objective_title,omitempty" db:"objective_title"`
	Level                string    `json:"level" db:"level"` // avanzato, intermedio, base, in_via_di_prima_acquisizione
	DimensionAutonomy    string    `json:"dimension_autonomy" db:"dimension_autonomy"`
	DimensionContinuity  string    `json:"dimension_continuity" db:"dimension_continuity"`
	DimensionFamiliarity string    `json:"dimension_familiarity" db:"dimension_familiarity"`
	DimensionResources   string    `json:"dimension_resources" db:"dimension_resources"`
	Date                 time.Time `json:"date" db:"date"`
	Semester             int       `json:"semester" db:"semester"` // 1 or 2
	Notes                string    `json:"notes" db:"notes"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type StudentObjectiveLevelItem struct {
	StudentID            string `json:"student_id" binding:"required"`
	Level                string `json:"level" binding:"required"` // avanzato, intermedio, base, in_via_di_prima_acquisizione
	DimensionAutonomy    string `json:"dimension_autonomy"`
	DimensionContinuity  string `json:"dimension_continuity"`
	DimensionFamiliarity string `json:"dimension_familiarity"`
	DimensionResources   string `json:"dimension_resources"`
	Notes                string `json:"notes"`
}

type SaveEvaluationsBatchRequest struct {
	ClassID     string                      `json:"class_id" binding:"required"`
	SubjectID   string                      `json:"subject_id" binding:"required"`
	ObjectiveID string                      `json:"objective_id" binding:"required"`
	Date        string                      `json:"date" binding:"required"`
	Semester    int                         `json:"semester" binding:"required"`
	Evaluations []StudentObjectiveLevelItem `json:"evaluations" binding:"required"`
}

type PrimaryEvaluationCell struct {
	Level                string    `json:"level"`
	DimensionAutonomy    string    `json:"dimension_autonomy"`
	DimensionContinuity  string    `json:"dimension_continuity"`
	DimensionFamiliarity string    `json:"dimension_familiarity"`
	DimensionResources   string    `json:"dimension_resources"`
	Notes                string    `json:"notes"`
	Date                 time.Time `json:"date"`
}

type PrimaryStudentMatrixRow struct {
	StudentID   string                           `json:"student_id"`
	StudentName string                           `json:"student_name"`
	Evaluations map[string]PrimaryEvaluationCell `json:"evaluations"` // objective_id -> PrimaryEvaluationCell
}

type PrimaryMatrixResponse struct {
	ClassID    string                    `json:"class_id"`
	SubjectID  string                    `json:"subject_id"`
	Semester   int                       `json:"semester"`
	Objectives []LearningObjective       `json:"objectives"`
	Students   []PrimaryStudentMatrixRow `json:"students"`
}
