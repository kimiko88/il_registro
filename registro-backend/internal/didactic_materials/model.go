package didactic_materials

import "time"

type DidacticMaterial struct {
	ID            string    `json:"id" db:"id"`
	SchoolID      string    `json:"school_id" db:"school_id"`
	ClassID       string    `json:"class_id" db:"class_id"`
	SubjectID     string    `json:"subject_id" db:"subject_id"`
	TeacherID     string    `json:"teacher_id" db:"teacher_id"`
	TeacherName   string    `json:"teacher_name,omitempty" db:"teacher_name"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	AttachmentURL string    `json:"attachment_url" db:"attachment_url"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateMaterialRequest struct {
	ClassID       string `json:"class_id" binding:"required"`
	SubjectID     string `json:"subject_id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	AttachmentURL string `json:"attachment_url"`
}

type MaterialResponse struct {
	ID            string    `json:"id"`
	ClassID       string    `json:"class_id"`
	SubjectID     string    `json:"subject_id"`
	TeacherID     string    `json:"teacher_id"`
	TeacherName   string    `json:"teacher_name"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	AttachmentURL string    `json:"attachment_url"`
	CreatedAt     time.Time `json:"created_at"`
}
