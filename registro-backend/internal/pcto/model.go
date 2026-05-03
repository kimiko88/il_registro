package pcto

import (
	"time"
)

type Company struct {
	ID            string    `json:"id" db:"id"`
	SchoolID      string    `json:"school_id" db:"school_id"`
	Name          string    `json:"name" db:"name"`
	VatNumber     string    `json:"vat_number" db:"vat_number"`
	Address       string    `json:"address" db:"address"`
	ContactPerson string    `json:"contact_person" db:"contact_person"`
	Email         string    `json:"email" db:"email"`
	AgreementDate time.Time `json:"agreement_date" db:"agreement_date"`
}

type Project struct {
	ID            string    `json:"id" db:"id"`
	SchoolID      string    `json:"school_id" db:"school_id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	Type          string    `json:"type" db:"type"` // Internal, External
	StartDate     time.Time `json:"start_date" db:"start_date"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
	TotalHours    int       `json:"total_hours" db:"total_hours"`
	CompanyID     *string   `json:"company_id,omitempty" db:"company_id"`
	SchoolTutorID *string   `json:"school_tutor_id,omitempty" db:"school_tutor_id"`
	CompanyTutor  string    `json:"company_tutor_name" db:"company_tutor_name"`
	CreatedBy     string    `json:"created_by" db:"created_by"`
}

type Participation struct {
	ID                string  `json:"id" db:"id"`
	ProjectID         string  `json:"project_id" db:"project_id"`
	StudentID         string  `json:"student_id" db:"student_id"`
	Status            string  `json:"status" db:"status"`
	HoursCompleted    float64 `json:"hours_completed" db:"hours_completed"`
	FinalEvaluation   string  `json:"final_evaluation" db:"final_evaluation"`
	RiskAssessmentAck bool    `json:"risk_assessment_ack" db:"risk_assessment_ack"`
}

type HourLog struct {
	ID              string    `json:"id" db:"id"`
	ParticipationID string    `json:"participation_id" db:"participation_id"`
	Date            time.Time `json:"date" db:"date"`
	Hours           float64   `json:"hours" db:"hours"`
	Activity        string    `json:"activity_description" db:"activity_description"`
	Verified        bool      `json:"verified" db:"verified"`
	VerifiedBy      *string   `json:"verified_by,omitempty" db:"verified_by"`
}

// Request DTOs
type CreateProjectRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Type         string  `json:"type" binding:"required"`
	StartDate    string  `json:"start_date" binding:"required"`
	EndDate      string  `json:"end_date" binding:"required"`
	TotalHours   int     `json:"total_hours" binding:"required"`
	CompanyID    *string `json:"company_id"`
	CompanyTutor string  `json:"company_tutor_name"`
}

type LogHourRequest struct {
	ProjectID string  `json:"project_id" binding:"required"`
	Date      string  `json:"date" binding:"required"`
	Hours     float64 `json:"hours" binding:"required"`
	Activity  string  `json:"activity" binding:"required"`
}

type PCTOStats struct {
	TotalProjects   int64   `json:"total_projects"`
	TotalStudents   int64   `json:"total_students"`
	TotalHours      float64 `json:"total_hours"`
	ActiveCompanies int64   `json:"active_companies"`
}
