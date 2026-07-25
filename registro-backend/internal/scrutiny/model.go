package scrutiny

import (
	"time"
)

type ScrutinyRecord struct {
	ID            string    `json:"id" db:"id"`
	StudentID     string    `json:"student_id" db:"student_id"`
	ClassID       string    `json:"class_id" db:"class_id"`
	Semester      int       `json:"semester" db:"semester"`
	PeriodType    string    `json:"period_type,omitempty" db:"period_type"` // 'semester_1', 'semester_2', 'infraquadrimestrale_1', 'infraquadrimestrale_2'
	ConductGrade  int       `json:"conduct_grade" db:"conduct_grade"`
	FinalDecision string    `json:"final_decision" db:"final_decision"`
	Notes         string    `json:"notes" db:"notes"`
	CoordinatorID string    `json:"coordinator_id" db:"coordinator_id"`
	Status        string    `json:"status" db:"status"` // draft, in_progress, submitted, validated
	ValidatedBy   string    `json:"validated_by,omitempty" db:"validated_by"`
	ValidatedAt   *time.Time `json:"validated_at,omitempty" db:"validated_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	
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

type ScrutinyMatrix struct {
	ClassID    string               `json:"class_id"`
	Semester   int                  `json:"semester"`
	PeriodType string               `json:"period_type,omitempty"`
	Subjects   []SubjectInfo        `json:"subjects"`
	Students   []StudentScrutinyRow `json:"students"`
}

type SubjectInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StudentScrutinyRow struct {
	StudentID      string                       `json:"student_id"`
	StudentName    string                       `json:"student_name"`
	SubjectData    map[string]SubjectAverages   `json:"subject_data"` // SubjectID -> Data
	Record         *ScrutinyRecord              `json:"record,omitempty"`
	AttendanceStats AttendanceSummary           `json:"attendance_stats"`
}

type AttendanceSummary struct {
	Absences   int `json:"absences"`
	Lates      int `json:"lates"`
	EarlyExits int `json:"early_exits"`
}

type SubjectAverages struct {
	Average     float64 `json:"average"`
	GradeCount  int     `json:"grade_count"`
	Proposed    float64 `json:"proposed"` // Rounded average
}

type SaveScrutinyRequest struct {
	StudentID     string `json:"student_id" binding:"required"`
	ClassID       string `json:"class_id" binding:"required"`
	Semester      int    `json:"semester" binding:"required"`
	PeriodType    string `json:"period_type"`
	ConductGrade  int    `json:"conduct_grade"`
	FinalDecision string `json:"final_decision"`
	Notes         string `json:"notes"`
	Grades        []SaveScrutinyGradeRequest `json:"grades"`
}

type SaveScrutinyGradeRequest struct {
	SubjectID  string  `json:"subject_id" binding:"required"`
	FinalGrade float64 `json:"final_grade" binding:"required"`
	TeacherID  string  `json:"teacher_id"`
}
