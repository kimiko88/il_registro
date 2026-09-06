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

type ParentDashboardStatsResponse struct {
	ChildrenCount        int `json:"children_count"`
	UpcomingColloqui     int `json:"upcoming_colloqui"`
	UnreadCommunications int `json:"unread_communications"`
	DocumentsCount       int `json:"documents_count"`
	TotalChildren        int `json:"total_children"`
	UnreadMessages       int `json:"unread_messages"`
	PendingPayments      int `json:"pending_payments"`
}
