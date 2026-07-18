package timetables

import (
	"time"
)

type ClassSchedule struct {
	ID          string    `json:"id"`
	ClassID     string    `json:"class_id"`
	DayOfWeek   int       `json:"day_of_week"`  // 1-7
	HourIndex   int       `json:"hour_index"`   // 1-12
	SubjectID   string    `json:"subject_id"`
	SubjectName string    `json:"subject_name,omitempty"`
	TeacherID   *string   `json:"teacher_id"`
	TeacherName string    `json:"teacher_name,omitempty"`
	Room        string    `json:"room"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateScheduleRequest struct {
	Entries []ScheduleEntry `json:"entries"`
}

type ScheduleEntry struct {
	DayOfWeek int     `json:"day_of_week"`
	HourIndex int     `json:"hour_index"`
	SubjectID string  `json:"subject_id"`
	TeacherID *string `json:"teacher_id"`
	Room      string  `json:"room"`
}
