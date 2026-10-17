package pnrr_projects

import "time"

type PNRRProject struct {
	ID                  string    `json:"id" db:"id"`
	SchoolID            string    `json:"school_id" db:"school_id"`
	CUPCode             string    `json:"cup_code" db:"cup_code"`
	ProjectTitle        string    `json:"project_title" db:"project_title"`
	NoticeType          string    `json:"notice_type" db:"notice_type"` // DM_65_STEM, DM_66_DIGITALE, DM_19_DISPERSIONE, SCUOLA_40
	TotalHours          int       `json:"total_hours" db:"total_hours"`
	AcademicYear        string    `json:"academic_year" db:"academic_year"`
	TutorID             *string   `json:"tutor_id,omitempty" db:"tutor_id"`
	ExpertID            *string   `json:"expert_id,omitempty" db:"expert_id"`
	TargetStudentsCount int       `json:"target_students_count" db:"target_students_count"`
	Status              string    `json:"status" db:"status"` // attivo, concluso, rendicontato
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`

	TutorName           string  `json:"tutor_name,omitempty"`
	ExpertName          string  `json:"expert_name,omitempty"`
	EnrolledStudents    int     `json:"enrolled_students"`
	Target75PassedCount int     `json:"target_75_passed_count"`
	OverallTargetRate   float64 `json:"overall_target_rate"`
}

type PNRRProjectStudent struct {
	ID                   string    `json:"id" db:"id"`
	ProjectID            string    `json:"project_id" db:"project_id"`
	StudentID            string    `json:"student_id" db:"student_id"`
	EnrolledAt           string    `json:"enrolled_at" db:"enrolled_at"` // YYYY-MM-DD
	AttendedHours        float64   `json:"attended_hours" db:"attended_hours"`
	AttendancePercentage float64   `json:"attendance_percentage" db:"attendance_percentage"`
	Target75Reached      bool      `json:"target_75_reached" db:"target_75_reached"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`

	StudentFullName string `json:"student_full_name,omitempty"`
}

type PNRRLessonAttendance struct {
	ID             string    `json:"id" db:"id"`
	ProjectID      string    `json:"project_id" db:"project_id"`
	LessonDate     string    `json:"lesson_date" db:"lesson_date"` // YYYY-MM-DD
	Hours          int       `json:"hours" db:"hours"`
	Topic          string    `json:"topic" db:"topic"`
	SignaturesJSON string    `json:"signatures_json" db:"signatures_json"` // student present minutes/signatures
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
