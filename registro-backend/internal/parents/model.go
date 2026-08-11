package parents

import (
	"registro-backend/internal/attendance"
	"registro-backend/internal/communications"
	"registro-backend/internal/grades"
	"registro-backend/internal/users"
)

type ChildOverview struct {
	Student          users.StudentChild          `json:"student"`
	Grades           []grades.Grade              `json:"grades"`
	AverageGrade     float64                     `json:"average_grade"`
	AttendanceStats  *attendance.SummaryResponse `json:"attendance_stats,omitempty"`
	PendingCirculars []*communications.Message   `json:"pending_circulars,omitempty"`
}

type ParentDashboardResponse struct {
	ParentID string          `json:"parent_id"`
	Children []ChildOverview `json:"children"`
}
