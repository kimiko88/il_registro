package teacher_activities

import "time"

// TeacherFreeActivity rappresenta un'attività del docente non collegata ad una classe
// (ore a disposizione, riunione, formazione, gita, etc.)
type TeacherFreeActivity struct {
	ID           string    `json:"id" db:"id"`
	TeacherID    string    `json:"teacher_id" db:"teacher_id"`
	TeacherName  string    `json:"teacher_name" db:"teacher_name"`
	Date         time.Time `json:"date" db:"date"`
	StartHour    int       `json:"start_hour" db:"start_hour"`   // Ora di inizio (1–10)
	Duration     int       `json:"duration" db:"duration"`       // Numero di ore coperte
	ActivityType string    `json:"activity_type" db:"activity_type"`
	// disponibilita | riunione | formazione | ptof | gita | altro
	Description string    `json:"description" db:"description"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
