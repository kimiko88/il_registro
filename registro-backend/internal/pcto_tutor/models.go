package pcto_tutor

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type CompanyTutor struct {
	ID             string     `json:"id"`
	SchoolID       string     `json:"school_id"`
	CompanyName    string     `json:"company_name"`
	TutorFirstName string     `json:"tutor_first_name"`
	TutorLastName  string     `json:"tutor_last_name"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone,omitempty"`
	AccessToken    string     `json:"access_token,omitempty"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type TutorAssignment struct {
	ID        string    `json:"id"`
	TutorID   string    `json:"tutor_id"`
	ProjectID string    `json:"project_id"`
	StudentID string    `json:"student_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AssignedStudentInfo struct {
	StudentID      string  `json:"student_id"`
	StudentName    string  `json:"student_name"`
	ProjectID      string  `json:"project_id"`
	ProjectTitle   string  `json:"project_title"`
	TotalHours     int     `json:"total_hours"`
	CompletedHours float64 `json:"completed_hours"`
	IsEvaluated    bool    `json:"is_evaluated"`
}

type TimesheetVerification struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	StudentID     string    `json:"student_id"`
	ActivityDate  string    `json:"activity_date"` // YYYY-MM-DD
	HoursDeclared float64   `json:"hours_declared"`
	HoursApproved float64   `json:"hours_approved"`
	TutorID       string    `json:"tutor_id"`
	SignedAt      time.Time `json:"signed_at"`
	TutorNotes    string    `json:"tutor_notes,omitempty"`
}

type CompanyEvaluation struct {
	ID               string    `json:"id"`
	TutorID          string    `json:"tutor_id"`
	StudentID        string    `json:"student_id"`
	ProjectID        string    `json:"project_id"`
	ReliabilityLevel int       `json:"reliability_level"` // 1-5
	TechnicalSkills  int       `json:"technical_skills"`  // 1-5
	TeamworkSkills   int       `json:"teamwork_skills"`   // 1-5
	FinalFeedback    string    `json:"final_feedback"`
	SubmittedAt      time.Time `json:"submitted_at"`
}

func GenerateMagicLinkToken(duration time.Duration) (string, time.Time, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}
	token := hex.EncodeToString(bytes)
	expiresAt := time.Now().Add(duration)
	return token, expiresAt, nil
}

func ValidateEvaluation(eval *CompanyEvaluation) error {
	if eval.ReliabilityLevel < 1 || eval.ReliabilityLevel > 5 {
		return errors.New("livello affidabilità deve essere tra 1 e 5")
	}
	if eval.TechnicalSkills < 1 || eval.TechnicalSkills > 5 {
		return errors.New("competenze tecniche devono essere tra 1 e 5")
	}
	if eval.TeamworkSkills < 1 || eval.TeamworkSkills > 5 {
		return errors.New("competenze di lavoro di gruppo devono essere tra 1 e 5")
	}
	return nil
}

func ValidateTimesheetEntry(entry *TimesheetVerification) error {
	if entry.HoursDeclared <= 0 {
		return errors.New("le ore dichiarate devono essere maggiori di zero")
	}
	if entry.HoursApproved < 0 {
		return errors.New("le ore approvate non possono essere negative")
	}
	if entry.HoursApproved > entry.HoursDeclared {
		return errors.New("le ore approvate non possono superare le ore dichiarate dallo studente")
	}
	return nil
}
