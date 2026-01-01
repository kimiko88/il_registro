package subjects

import "time"

type Subject struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	Name        string    `json:"name" db:"name"`
	Code        string    `json:"code,omitempty" db:"code"`
	Description string    `json:"description,omitempty" db:"description"`
	IsMandatory bool      `json:"is_mandatory" db:"is_mandatory"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateSubjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsMandatory *bool  `json:"is_mandatory"` // Pointer to distinguish false vs nil (optional)
	SchoolID    string `json:"school_id"`    // Optional override for superadmin
}
