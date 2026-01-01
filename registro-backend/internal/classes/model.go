package classes

import "time"

type Class struct {
	ID            string    `json:"id"`
	SchoolID      string    `json:"school_id"`
	Name          string    `json:"name"`          // e.g., "1A", "5B"
	Section       string    `json:"section"`       // e.g., "A", "B"
	AcademicYear  string    `json:"academic_year"` // e.g., "2024/2025"
	CoordinatorID string    `json:"coordinator_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateClassRequest struct {
	Name          string `json:"name" binding:"required"`
	SchoolID      string `json:"school_id"`
	Section       string `json:"section"`
	AcademicYear  string `json:"academic_year" binding:"required"`
	CoordinatorID string `json:"coordinator_id"`
}
