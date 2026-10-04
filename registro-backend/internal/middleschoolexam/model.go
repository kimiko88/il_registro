package middleschoolexam

import "time"

type MiddleSchoolExam struct {
	ID                  string    `json:"id" db:"id"`
	SchoolID            string    `json:"school_id" db:"school_id"`
	ClassID             string    `json:"class_id" db:"class_id"`
	AcademicYear        string    `json:"academic_year" db:"academic_year"`
	SubcommissionNumber int       `json:"subcommission_number" db:"subcommission_number"`
	PresidentName       string    `json:"president_name" db:"president_name"`
	Status              string    `json:"status" db:"status"` // admission, in_progress, deliberated, closed
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

type MiddleSchoolExamCandidate struct {
	ID              string     `json:"id" db:"id"`
	ExamID          string     `json:"exam_id" db:"exam_id"`
	StudentID       string     `json:"student_id" db:"student_id"`
	StudentName     string     `json:"student_name,omitempty" db:"student_name"`
	AdmissionGrade  int        `json:"admission_grade" db:"admission_grade"`
	AdmissionJudgment string   `json:"admission_judgment" db:"admission_judgment"`
	IsAdmitted      bool       `json:"is_admitted" db:"is_admitted"`
	GradeItalian    float64    `json:"grade_italian" db:"grade_italian"`
	GradeMath       float64    `json:"grade_math" db:"grade_math"`
	GradeEnglish    float64    `json:"grade_english" db:"grade_english"`
	GradeSecondLang float64    `json:"grade_second_lang" db:"grade_second_lang"`
	GradeInterview  float64    `json:"grade_interview" db:"grade_interview"`
	ExamMean        float64    `json:"exam_mean" db:"exam_mean"`
	FinalGrade      int        `json:"final_grade" db:"final_grade"`
	HasHonors       bool       `json:"has_honors" db:"has_honors"`
	Outcome         string     `json:"outcome" db:"outcome"` // licenziato, non_licenziato
	DeliberatedAt   *time.Time `json:"deliberated_at,omitempty" db:"deliberated_at"`
	Notes           string     `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

type ExamCandidateGrades struct {
	AdmissionGrade  int     `json:"admission_grade"`
	GradeItalian    float64 `json:"grade_italian"`
	GradeMath       float64 `json:"grade_math"`
	GradeEnglish    float64 `json:"grade_english"`
	GradeSecondLang float64 `json:"grade_second_lang"`
	GradeInterview  float64 `json:"grade_interview"`
	ProposedHonors  bool    `json:"proposed_honors"`
}

type ExamOutcomeResult struct {
	ExamMean   float64 `json:"exam_mean"`
	FinalGrade int     `json:"final_grade"`
	HasHonors  bool    `json:"has_honors"`
	Outcome    string  `json:"outcome"`
}
