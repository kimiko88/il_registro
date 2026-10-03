package attendance

import "time"

// AnalyticsResponse contiene le statistiche aggregate di assenza per la scuola.
type AnalyticsResponse struct {
	AverageAbsenceRate float64  `json:"average_absence_rate"`
	TopAbsentees       []string `json:"top_absentees"`
}

// SummaryResponse contiene il riepilogo delle assenze di uno studente.
type SummaryResponse struct {
	TotalAbsences   int     `json:"total_absences"`
	TotalLates      int     `json:"total_lates"`
	TotalEarlyExits int     `json:"total_early_exits"`
	JustifiedCount  int     `json:"justified_count"`
	AbsenceRate     float64 `json:"absence_rate"`
	RiskLevel       string  `json:"risk_level"`
}

// AttendanceResponse è la risposta DTO per un singolo record di presenza.
type AttendanceResponse struct {
	ID          string           `json:"id"`
	StudentID   string           `json:"student_id"`
	Date        string           `json:"date"`
	Hour        int              `json:"hour"`
	Status      AttendanceStatus `json:"status"`
	IsJustified bool             `json:"is_justified"`
	Notes       string           `json:"notes,omitempty"`
	EntryTime   string           `json:"entry_time,omitempty"`
	ExitTime    string           `json:"exit_time,omitempty"`
}

// ClassDailyAttendance è la risposta per le presenze di una classe in un giorno.
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

// JustificationResponse è la risposta DTO per una giustifica.
type JustificationResponse struct {
	ID          string `json:"id"`
	StudentID   string `json:"student_id,omitempty"`
	StudentName string `json:"student_name,omitempty"`
	Date        string `json:"date,omitempty"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
	DateRange   string `json:"date_range"`
}

// --- Request DTOs ---

type CreateAttendanceRequest struct {
	StudentID string           `json:"student_id" binding:"required"`
	ClassID   string           `json:"class_id" binding:"required"`
	Date      string           `json:"date" binding:"required"`
	Hour      int              `json:"hour"`
	SubjectID string           `json:"subject_id"`
	Status    AttendanceStatus `json:"status" binding:"required"`
	Notes     string           `json:"notes"`
	EntryTime string           `json:"entry_time"`
	ExitTime  string           `json:"exit_time"`
}

type StudentStatusRequest struct {
	StudentID string           `json:"student_id" binding:"required"`
	Status    AttendanceStatus `json:"status" binding:"required"`
	Notes     string           `json:"notes"`
	EntryTime string           `json:"entry_time"`
	ExitTime  string           `json:"exit_time"`
}

type BulkAttendanceRequest struct {
	ClassID        string                 `json:"class_id" binding:"required"`
	Date           string                 `json:"date" binding:"required"`
	Hour           int                    `json:"hour"`
	SubjectID      string                 `json:"subject_id"`
	IsSubstitution bool                   `json:"is_substitution"`
	Statuses       []StudentStatusRequest `json:"statuses" binding:"required"`
}

type UpdateAttendanceRequest struct {
	Status    *AttendanceStatus `json:"status"`
	Notes     *string           `json:"notes"`
	EntryTime *string           `json:"entry_time"`
	ExitTime  *string           `json:"exit_time"`
}

type JustificationRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
}

// Ensure time package is used (needed if time.Time appears in other DTOs).
var _ = time.Now

// --- School Meals & Parent PIN DTOs ---

type StudentMealItem struct {
	StudentID string `json:"student_id" binding:"required"`
	MealType  string `json:"meal_type" binding:"required"` // standard, bianco, dieta_sanitaria, dieta_etico_religiosa, nessuno
	MealNotes string `json:"meal_notes,omitempty"`
}

type SaveMealsBatchRequest struct {
	ClassID string            `json:"class_id" binding:"required"`
	Date    string            `json:"date" binding:"required"`
	Meals   []StudentMealItem `json:"meals" binding:"required"`
}

type ClassMealsSummary struct {
	ClassID         string `json:"class_id"`
	ClassName       string `json:"class_name"`
	TotalMeals      int    `json:"total_meals"`
	StandardCount   int    `json:"standard_count"`
	WhiteCount      int    `json:"white_count"`
	HealthDietCount int    `json:"health_diet_count"`
	EthicDietCount  int    `json:"ethic_diet_count"`
}

type StudentSpecialDietDetail struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	ClassName   string `json:"class_name"`
	MealType    string `json:"meal_type"`
	MealNotes   string `json:"meal_notes"`
}

type DailyMealsReportResponse struct {
	Date                string                     `json:"date"`
	SchoolID            string                     `json:"school_id"`
	TotalMeals          int                        `json:"total_meals"`
	StandardCount       int                        `json:"standard_count"`
	WhiteCount          int                        `json:"white_count"`
	HealthDietCount     int                        `json:"health_diet_count"`
	EthicDietCount      int                        `json:"ethic_diet_count"`
	ClassesBreakdown    []ClassMealsSummary        `json:"classes_breakdown"`
	SpecialDietsDetails []StudentSpecialDietDetail `json:"special_diets_details"`
}

type VerifyPinAndJustifyRequest struct {
	AttendanceID string `json:"attendance_id" binding:"required"`
	Pin          string `json:"pin" binding:"required"`
	Reason       string `json:"reason"`
	Notes        string `json:"notes"`
}

type AbsenceLimitStatusResponse struct {
	StudentID         string  `json:"student_id"`
	TotalSchoolHours  int     `json:"total_school_hours"`
	AbsentHours       int     `json:"absent_hours"`
	AbsenceRate       float64 `json:"absence_rate"`
	MaxLimitRate      float64 `json:"max_limit_rate"` // 25.0% (D.P.R. 122/2009)
	IsExceedingLimit  bool    `json:"is_exceeding_limit"`
	HealthDerogations int     `json:"health_derogations"`
}
