package lessons

import "time"

type CreateLessonRequest struct {
	ClassID              string  `json:"class_id" binding:"required"`
	SubjectID            string  `json:"subject_id" binding:"required"`
	Date                 string  `json:"date" binding:"required"` // YYYY-MM-DD
	Hour                 int     `json:"hour" binding:"required"`
	Duration             int     `json:"duration" binding:"required"`
	Topic                string  `json:"topic" binding:"required"`
	Type                 string  `json:"type" binding:"required"`
	GroupID              *string `json:"group_id"`
	IsSubstitution       bool    `json:"is_substitution"`
	SubstitutedTeacherID *string `json:"substituted_teacher_id"`
	ActivityType         string  `json:"activity_type"` // e.g. standard, substitution, ptof, project, assembly, trip, lab, other
	IsCoTeaching         bool    `json:"is_co_teaching"`
	Notes                string  `json:"notes"`
}

type UpdateLessonRequest struct {
	Topic                string  `json:"topic"`
	Type                 string  `json:"type"`
	Hour                 *int    `json:"hour"`
	Duration             *int    `json:"duration"`
	GroupID              *string `json:"group_id"`
	IsSubstitution       *bool   `json:"is_substitution"`
	SubstitutedTeacherID *string `json:"substituted_teacher_id"`
	ActivityType         string  `json:"activity_type"`
	IsCoTeaching         *bool   `json:"is_co_teaching"`
	Notes                string  `json:"notes"`
}

type CreateHomeworkRequest struct {
	ClassID     string  `json:"class_id" binding:"required"`
	SubjectID   string  `json:"subject_id" binding:"required"`
	LessonID    *string `json:"lesson_id"`
	DueDate     string  `json:"due_date" binding:"required"` // YYYY-MM-DD
	Description string  `json:"description" binding:"required"`
	Type        string  `json:"type"` // compito, verifica, avviso, interrogazione
}

type UpdateHomeworkRequest struct {
	DueDate     string `json:"due_date"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type LessonResponse struct {
	ID                     string    `json:"id"`
	ClassID                string    `json:"class_id"`
	TeacherID              string    `json:"teacher_id"`
	TeacherName            string    `json:"teacher_name"`
	SubjectID              string    `json:"subject_id"`
	Date                   time.Time `json:"date"`
	Hour                   int       `json:"hour"`
	Duration               int       `json:"duration"`
	Topic                  string    `json:"topic"`
	Type                   string    `json:"type"`
	GroupID                *string   `json:"group_id,omitempty"`
	IsSubstitution         bool      `json:"is_substitution"`
	SubstitutedTeacherID   *string   `json:"substituted_teacher_id,omitempty"`
	SubstitutedTeacherName string    `json:"substituted_teacher_name,omitempty"`
	ActivityType           string    `json:"activity_type"`
	IsCoTeaching           bool      `json:"is_co_teaching"`
	Notes                  string    `json:"notes"`
}

type HomeworkResponse struct {
	ID          string    `json:"id"`
	LessonID    *string   `json:"lesson_id,omitempty"`
	ClassID     string    `json:"class_id"`
	SubjectID   string    `json:"subject_id"`
	TeacherID   string    `json:"teacher_id"`
	TeacherName string    `json:"teacher_name"`
	DueDate     time.Time `json:"due_date"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
}
