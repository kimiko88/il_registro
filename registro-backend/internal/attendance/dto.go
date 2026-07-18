package attendance

// Requests

type CreateAttendanceRequest struct {
	StudentID string           `json:"student_id" binding:"required"`
	ClassID   string           `json:"class_id" binding:"required"`
	Date      string           `json:"date" binding:"required"` // YYYY-MM-DD
	Hour      int              `json:"hour" binding:"required"`
	SubjectID string           `json:"subject_id"`
	Status    AttendanceStatus `json:"status" binding:"required"`
	EntryTime string           `json:"entry_time,omitempty"` // HH:MM
	ExitTime  string           `json:"exit_time,omitempty"`  // HH:MM
	Notes     string           `json:"notes,omitempty"`
}

type UpdateAttendanceRequest struct {
	Status    *AttendanceStatus `json:"status,omitempty"`
	EntryTime *string           `json:"entry_time,omitempty"`
	ExitTime  *string           `json:"exit_time,omitempty"`
	Notes     *string           `json:"notes,omitempty"`
}

type BulkAttendanceRequest struct {
	ClassID   string                    `json:"class_id" binding:"required"`
	Date      string                    `json:"date" binding:"required"`
	Hour      int                       `json:"hour" binding:"required"`
	SubjectID string                    `json:"subject_id" binding:"required"`
	Statuses  []CreateAttendanceRequest `json:"statuses" binding:"required"`
}

type JustificationRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
}

// Responses

type AttendanceResponse struct {
	ID          string           `json:"id"`
	StudentID   string           `json:"student_id"`
	StudentName string           `json:"student_name,omitempty"` // Enriched
	Date        string           `json:"date"`
	Status      AttendanceStatus `json:"status"`
	EntryTime   string           `json:"entry_time,omitempty"`
	ExitTime    string           `json:"exit_time,omitempty"`
	IsJustified bool             `json:"is_justified"`
	Notes       string           `json:"notes,omitempty"`
}

type SummaryResponse struct {
	TotalAbsences   int     `json:"total_absences"`
	TotalLates      int     `json:"total_lates"`
	TotalEarlyExits int     `json:"total_early_exits"`
	JustifiedCount  int     `json:"justified_count"`
	AbsenceRate     float64 `json:"absence_rate"` // %
	RiskLevel       string  `json:"risk_level"`   // Normal, Warning, Critical
	PendingRequests int     `json:"pending_requests"`
}

type JustificationResponse struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Requestor string `json:"requestor"`
	DateRange string `json:"date_range"`
	Reason    string `json:"reason"`
}

type ClassDailyAttendance struct {
	ClassID string               `json:"class_id"`
	Date    string               `json:"date"`
	Records []AttendanceResponse `json:"records"`
	Summary struct {
		Present int `json:"present"`
		Absent  int `json:"absent"`
		Late    int `json:"late"`
	} `json:"summary"`
}

type AnalyticsResponse struct {
	AverageAbsenceRate float64  `json:"average_absence_rate"`
	TopAbsentees       []string `json:"top_absentees"`
}
