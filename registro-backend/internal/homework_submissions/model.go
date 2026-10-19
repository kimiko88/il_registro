package homework_submissions

import "time"

type HomeworkSubmission struct {
	ID                          string     `json:"id" db:"id"`
	HomeworkID                  string     `json:"homework_id" db:"homework_id"`
	StudentID                   string     `json:"student_id" db:"student_id"`
	FileURL                     string     `json:"file_url" db:"file_url"`
	FileType                    string     `json:"file_type" db:"file_type"` // pdf, image, document
	SubmittedAt                 time.Time  `json:"submitted_at" db:"submitted_at"`
	TeacherAnnotatedFileURL     string     `json:"teacher_annotated_file_url" db:"teacher_annotated_file_url"`
	AnnotationsJSON             string     `json:"annotations_json" db:"annotations_json"` // coordinates of pen/highlighter/text
	TeacherGrade                *float64   `json:"teacher_grade,omitempty" db:"teacher_grade"`
	TeacherFeedback             string     `json:"teacher_feedback" db:"teacher_feedback"`
	Status                      string     `json:"status" db:"status"` // consegnato, corretto, restituito
	GradedAt                    *time.Time `json:"graded_at,omitempty" db:"graded_at"`
	GradedBy                    *string    `json:"graded_by,omitempty" db:"graded_by"`
	GoogleClassroomSubmissionID string     `json:"google_classroom_submission_id" db:"google_classroom_submission_id"`
	SyncedWithClassroom         bool       `json:"synced_with_classroom" db:"synced_with_classroom"`
	CreatedAt                   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at" db:"updated_at"`

	StudentFullName string `json:"student_full_name,omitempty"`
	HomeworkTitle   string `json:"homework_title,omitempty"`
}

type GradeSubmissionRequest struct {
	AnnotationsJSON string   `json:"annotations_json"`
	AnnotatedURL    string   `json:"teacher_annotated_file_url"`
	TeacherGrade    *float64 `json:"teacher_grade"`
	TeacherFeedback string   `json:"teacher_feedback"`
	SyncClassroom   bool     `json:"sync_classroom"`
}
