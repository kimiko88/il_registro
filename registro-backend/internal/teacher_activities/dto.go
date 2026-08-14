package teacher_activities

import "time"

// CreateTeacherActivityRequest è il payload per creare una nuova attività libera.
type CreateTeacherActivityRequest struct {
	Date         string `json:"date" binding:"required"`       // YYYY-MM-DD
	StartHour    int    `json:"start_hour" binding:"required"` // 1–10
	Duration     int    `json:"duration" binding:"required"`   // >= 1
	ActivityType string `json:"activity_type"`                 // vedi valori validi in service.go
	Description  string `json:"description" binding:"required"`
	Notes        string `json:"notes"`
}

// UpdateTeacherActivityRequest è il payload per modificare un'attività libera.
type UpdateTeacherActivityRequest struct {
	Date         string `json:"date"`
	StartHour    *int   `json:"start_hour"`
	Duration     *int   `json:"duration"`
	ActivityType string `json:"activity_type"`
	Description  string `json:"description"`
	Notes        string `json:"notes"`
}

// TeacherActivityResponse è la risposta API per un'attività libera.
type TeacherActivityResponse struct {
	ID           string    `json:"id"`
	TeacherID    string    `json:"teacher_id"`
	TeacherName  string    `json:"teacher_name"`
	Date         time.Time `json:"date"`
	StartHour    int       `json:"start_hour"`
	Duration     int       `json:"duration"`
	ActivityType string    `json:"activity_type"`
	Description  string    `json:"description"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
