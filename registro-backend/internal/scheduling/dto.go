package scheduling

import "time"

type CreateScheduleRequest struct {
	ClassID   uint         `json:"class_id" binding:"required"`
	SubjectID uint         `json:"subject_id" binding:"required"`
	TeacherID uint         `json:"teacher_id" binding:"required"`
	DayOfWeek time.Weekday `json:"day_of_week" binding:"required"`
	StartTime string       `json:"start_time" binding:"required"`
	EndTime   string       `json:"end_time" binding:"required"`
	Room      string       `json:"room"`
}

type ScheduleResponse struct {
	ID        uint         `json:"id"`
	ClassID   uint         `json:"class_id"`
	Subject   string       `json:"subject"`
	Teacher   string       `json:"teacher"`
	DayOfWeek time.Weekday `json:"day_of_week"`
	StartTime string       `json:"start_time"`
	EndTime   string       `json:"end_time"`
	Room      string       `json:"room"`
}
