package recovery

import "time"

type RecoveryCourse struct {
	ID           string                  `json:"id"`
	SchoolID     string                  `json:"school_id"`
	SubjectID    string                  `json:"subject_id"`
	SubjectName  string                  `json:"subject_name,omitempty"`
	TeacherID    string                  `json:"teacher_id"`
	TeacherName  string                  `json:"teacher_name,omitempty"`
	Title        string                  `json:"title"`
	Description  string                  `json:"description"`
	AcademicYear string                  `json:"academic_year"`
	Period       string                  `json:"period"` // 'summer', 'intermedio', 'pomeridiano'
	TotalHours   int                     `json:"total_hours"`
	Room         string                  `json:"room"`
	Status       string                  `json:"status"` // 'scheduled', 'in_progress', 'completed', 'cancelled'
	Sessions     []RecoveryCourseSession `json:"sessions,omitempty"`
	Students     []RecoveryCourseStudent `json:"students,omitempty"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

type RecoveryCourseSession struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	SessionDate string    `json:"session_date"` // YYYY-MM-DD
	StartTime   string    `json:"start_time"`   // HH:MM
	EndTime     string    `json:"end_time"`     // HH:MM
	Room        string    `json:"room"`
	Topic       string    `json:"topic"`
	CreatedAt   time.Time `json:"created_at"`
}

type RecoveryCourseStudent struct {
	ID              string  `json:"id"`
	CourseID        string  `json:"course_id"`
	StudentID       string  `json:"student_id"`
	StudentName     string  `json:"student_name,omitempty"`
	ClassName       string  `json:"class_name,omitempty"`
	AttendanceHours float64 `json:"attendance_hours"`
	Notes           string  `json:"notes"`
}

type RecoveryTest struct {
	ID                string    `json:"id"`
	SchoolID          string    `json:"school_id"`
	DeficiencyID      *string   `json:"deficiency_id,omitempty"`
	StudentID         string    `json:"student_id"`
	StudentName       string    `json:"student_name,omitempty"`
	SubjectID         string    `json:"subject_id"`
	SubjectName       string    `json:"subject_name,omitempty"`
	ClassID           string    `json:"class_id"`
	ClassName         string    `json:"class_name,omitempty"`
	TeacherID         string    `json:"teacher_id"`
	TeacherName       string    `json:"teacher_name,omitempty"`
	TestDate          string    `json:"test_date"` // YYYY-MM-DD
	TestType          string    `json:"test_type"` // 'written', 'oral', 'practical', 'mixed'
	Grade             float64   `json:"grade"`
	Outcome           string    `json:"outcome"`            // 'recuperato', 'non_recuperato'
	FinalDeliberation string    `json:"final_deliberation"` // 'Ammesso', 'Non Ammesso', 'Ammesso con delibera'
	VerbaleNumber     string    `json:"verbale_number"`
	Notes             string    `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateCourseRequest struct {
	SubjectID    string                  `json:"subject_id" binding:"required"`
	TeacherID    string                  `json:"teacher_id" binding:"required"`
	Title        string                  `json:"title" binding:"required"`
	Description  string                  `json:"description"`
	AcademicYear string                  `json:"academic_year" binding:"required"`
	Period       string                  `json:"period"`
	TotalHours   int                     `json:"total_hours"`
	Room         string                  `json:"room"`
	Sessions     []RecoveryCourseSession `json:"sessions"`
	StudentIDs   []string                `json:"student_ids"`
}

type RecordTestOutcomeRequest struct {
	DeficiencyID      *string `json:"deficiency_id"`
	StudentID         string  `json:"student_id" binding:"required"`
	SubjectID         string  `json:"subject_id" binding:"required"`
	ClassID           string  `json:"class_id" binding:"required"`
	TestDate          string  `json:"test_date" binding:"required"`
	TestType          string  `json:"test_type" binding:"required"`
	Grade             float64 `json:"grade" binding:"required"`
	VerbaleNumber     string  `json:"verbale_number"`
	FinalDeliberation string  `json:"final_deliberation"` // optional override
	Notes             string  `json:"notes"`
}
