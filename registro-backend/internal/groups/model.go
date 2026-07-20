package groups

import "time"

type Group struct {
	ID           string    `json:"id" db:"id"`
	SchoolID     string    `json:"school_id" db:"school_id"`
	Name         string    `json:"name" db:"name"`
	SubjectID    *string   `json:"subject_id,omitempty" db:"subject_id"`
	SubjectName  string    `json:"subject_name,omitempty" db:"subject_name"`
	TeacherID    *string   `json:"teacher_id,omitempty" db:"teacher_id"`
	TeacherName  string    `json:"teacher_name,omitempty" db:"teacher_name"`
	AcademicYear string    `json:"academic_year" db:"academic_year"`
	Description  string    `json:"description" db:"description"`
	StudentCount int       `json:"student_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	Students []GroupStudentInfo `json:"students,omitempty"`
}

type GroupStudentInfo struct {
	StudentID string `json:"student_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ClassName string `json:"class_name,omitempty"`
}

type CreateGroupRequest struct {
	SchoolID     string   `json:"school_id" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	SubjectID    *string  `json:"subject_id"`
	TeacherID    *string  `json:"teacher_id"`
	AcademicYear string   `json:"academic_year"`
	Description  string   `json:"description"`
	StudentIDs   []string `json:"student_ids"`
}

type UpdateGroupRequest struct {
	Name        string  `json:"name"`
	SubjectID   *string `json:"subject_id"`
	TeacherID   *string `json:"teacher_id"`
	Description string  `json:"description"`
}

type AddGroupStudentsRequest struct {
	StudentIDs []string `json:"student_ids" binding:"required"`
}
