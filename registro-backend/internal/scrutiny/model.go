package scrutiny

import (
	"time"
)

type ScrutinyRecord struct {
	ID            string     `json:"id" db:"id"`
	StudentID     string     `json:"student_id" db:"student_id"`
	ClassID       string     `json:"class_id" db:"class_id"`
	Semester      int        `json:"semester" db:"semester"`
	PeriodType    string     `json:"period_type,omitempty" db:"period_type"` // 'semester_1', 'semester_2', 'differito'
	ConductGrade  int        `json:"conduct_grade" db:"conduct_grade"`
	FinalDecision string     `json:"final_decision" db:"final_decision"`
	Notes         string     `json:"notes" db:"notes"`
	CoordinatorID string     `json:"coordinator_id" db:"coordinator_id"`
	Status        string     `json:"status" db:"status"` // draft, in_progress, submitted, validated
	ValidatedBy   string     `json:"validated_by,omitempty" db:"validated_by"`
	ValidatedAt   *time.Time `json:"validated_at,omitempty" db:"validated_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`

	// Grades for this record
	Grades []ScrutinyGrade `json:"grades"`
}

type ScrutinyGrade struct {
	ID               string  `json:"id" db:"id"`
	ScrutinyRecordID string  `json:"scrutiny_record_id" db:"scrutiny_record_id"`
	SubjectID        string  `json:"subject_id" db:"subject_id"`
	FinalGrade       float64 `json:"final_grade" db:"final_grade"`
	TeacherID        string  `json:"teacher_id" db:"teacher_id"`
}

type StudentDeficiency struct {
	ID               string     `json:"id" db:"id"`
	SchoolID         string     `json:"school_id" db:"school_id"`
	StudentID        string     `json:"student_id" db:"student_id"`
	StudentName      string     `json:"student_name,omitempty" db:"student_name"`
	ClassID          string     `json:"class_id" db:"class_id"`
	SubjectID        string     `json:"subject_id" db:"subject_id"`
	SubjectName      string     `json:"subject_name,omitempty" db:"subject_name"`
	ScrutinyRecordID string     `json:"scrutiny_record_id,omitempty" db:"scrutiny_record_id"`
	Semester         int        `json:"semester" db:"semester"`
	PeriodType       string     `json:"period_type" db:"period_type"`
	Topics           string     `json:"topics" db:"topics"` // Argomenti/lacune della carenza
	RecoveryMode     string     `json:"recovery_mode" db:"recovery_mode"`
	Status           string     `json:"status" db:"status"` // 'da_recuperare', 'in_corso', 'recuperato', 'non_recuperato'
	RecoveryGrade    *float64   `json:"recovery_grade,omitempty" db:"recovery_grade"`
	RecoveryDate     *time.Time `json:"recovery_date,omitempty" db:"recovery_date"`
	Notes            string     `json:"notes,omitempty" db:"notes"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type ScrutinyMatrix struct {
	ClassID      string               `json:"class_id"`
	Semester     int                  `json:"semester"`
	PeriodType   string               `json:"period_type,omitempty"`
	Subjects     []SubjectInfo        `json:"subjects"`
	Students     []StudentScrutinyRow `json:"students"`
	Deficiencies []StudentDeficiency  `json:"deficiencies,omitempty"`
}

type SubjectInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StudentScrutinyRow struct {
	StudentID       string                     `json:"student_id"`
	StudentName     string                     `json:"student_name"`
	SubjectData     map[string]SubjectAverages `json:"subject_data"` // SubjectID -> Data
	Record          *ScrutinyRecord            `json:"record,omitempty"`
	AttendanceStats AttendanceSummary          `json:"attendance_stats"`
}

type AttendanceSummary struct {
	Absences   int `json:"absences"`
	Lates      int `json:"lates"`
	EarlyExits int `json:"early_exits"`
}

type SubjectAverages struct {
	Average    float64 `json:"average"`
	GradeCount int     `json:"grade_count"`
	Proposed   float64 `json:"proposed"` // Rounded average
}

type SaveScrutinyRequest struct {
	StudentID     string                     `json:"student_id" binding:"required"`
	ClassID       string                     `json:"class_id" binding:"required"`
	Semester      int                        `json:"semester" binding:"required"`
	PeriodType    string                     `json:"period_type"`
	ConductGrade  int                        `json:"conduct_grade"`
	FinalDecision string                     `json:"final_decision"`
	Notes         string                     `json:"notes"`
	Grades        []SaveScrutinyGradeRequest `json:"grades"`
}

type SaveScrutinyGradeRequest struct {
	SubjectID  string  `json:"subject_id" binding:"required"`
	FinalGrade float64 `json:"final_grade" binding:"required"`
	TeacherID  string  `json:"teacher_id"`
}

type SaveDeficiencyRequest struct {
	ID            string   `json:"id,omitempty"`
	SchoolID      string   `json:"school_id"`
	StudentID     string   `json:"student_id" binding:"required"`
	ClassID       string   `json:"class_id" binding:"required"`
	SubjectID     string   `json:"subject_id" binding:"required"`
	Semester      int      `json:"semester" binding:"required"`
	PeriodType    string   `json:"period_type"`
	Topics        string   `json:"topics"`        // Argomenti della carenza
	RecoveryMode  string   `json:"recovery_mode"` // studio_individuale, corso_recupero, sportello_didattico
	Status        string   `json:"status"`        // da_recuperare, in_corso, recuperato, non_recuperato
	RecoveryGrade *float64 `json:"recovery_grade,omitempty"`
	RecoveryDate  *string  `json:"recovery_date,omitempty"`
	Notes         string   `json:"notes"`
}

type SaveDeferredScrutinyRequest struct {
	StudentID     string                          `json:"student_id" binding:"required"`
	ClassID       string                          `json:"class_id" binding:"required"`
	FinalDecision string                          `json:"final_decision" binding:"required"` // promosso_con_debiti_saldati, non_promosso
	Notes         string                          `json:"notes"`
	Deficiencies  []SaveDeferredDeficiencyItemReq `json:"deficiencies"`
}

type SaveDeferredDeficiencyItemReq struct {
	DeficiencyID  string   `json:"deficiency_id" binding:"required"`
	Status        string   `json:"status" binding:"required"` // recuperato, non_recuperato
	RecoveryGrade *float64 `json:"recovery_grade,omitempty"`
	RecoveryDate  *string  `json:"recovery_date,omitempty"`
}
