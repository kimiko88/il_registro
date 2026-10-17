package interpelli

import (
	"time"
)

type InterpelloNotice struct {
	ID             string    `json:"id" db:"id"`
	SchoolID       string    `json:"school_id" db:"school_id"`
	ProtocolNumber string    `json:"protocol_number" db:"protocol_number"`
	Title          string    `json:"title" db:"title"`
	ConcorsoClass  string    `json:"concorso_class" db:"concorso_class"` // e.g. "A026", "ADSS", "A012"
	PostType       string    `json:"post_type" db:"post_type"`           // "comune", "sostegno", "potenziamento"
	WeeklyHours    int       `json:"weekly_hours" db:"weekly_hours"`
	StartDate      time.Time `json:"start_date" db:"start_date"`
	EndDate        time.Time `json:"end_date" db:"end_date"`
	Deadline       time.Time `json:"deadline" db:"deadline"` // Termine perentorio O.M. 88/2024
	Status         string    `json:"status" db:"status"`     // "aperto", "chiuso", "assegnato", "annullato"
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type InterpelloCandidatura struct {
	ID                  string     `json:"id" db:"id"`
	NoticeID            string     `json:"notice_id" db:"notice_id"`
	CandidateName       string     `json:"candidate_name" db:"candidate_name"`
	CandidateSurname    string     `json:"candidate_surname" db:"candidate_surname"`
	FiscalCode          string     `json:"fiscal_code" db:"fiscal_code"`
	Email               string     `json:"email" db:"email"`
	PEC                 string     `json:"pec,omitempty" db:"pec"`
	Phone               string     `json:"phone" db:"phone"`
	GraduationGrade     float64    `json:"graduation_grade" db:"graduation_grade"` // out of 110
	GraduationLode      bool       `json:"graduation_lode" db:"graduation_lode"`
	HasHabilitation     bool       `json:"has_habitation" db:"has_habitation"` // Abilitazione all'insegnamento / TFA
	ScoreService        float64    `json:"score_service" db:"score_service"`   // Mesi di servizio x punteggio
	ScoreCerts          float64    `json:"score_certs" db:"score_certs"`       // Lingua (B2/C1/C2), informatica
	TotalScore          float64    `json:"total_score" db:"total_score"`
	CVUrl               string     `json:"cv_url" db:"cv_url"`
	DPR445Declared      bool       `json:"dpr445_declared" db:"dpr445_declared"` // Autocertificazione veridicità
	Status              string     `json:"status" db:"status"`                   // "ricevuta", "valutata", "convocata", "accettata", "rifiutata"
	ConvocationSentAt   *time.Time `json:"convocation_sent_at,omitempty" db:"convocation_sent_at"`
	ConvocationDeadline *time.Time `json:"convocation_deadline,omitempty" db:"convocation_deadline"`
	ResponseAt          *time.Time `json:"response_at,omitempty" db:"response_at"`
	ResponseNotes       string     `json:"response_notes,omitempty" db:"response_notes"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
}

type CreateNoticeRequest struct {
	ProtocolNumber string `json:"protocol_number"`
	Title          string `json:"title" binding:"required"`
	ConcorsoClass  string `json:"concorso_class" binding:"required"`
	PostType       string `json:"post_type"`
	WeeklyHours    int    `json:"weekly_hours" binding:"required,gt=0"`
	StartDate      string `json:"start_date" binding:"required"`
	EndDate        string `json:"end_date" binding:"required"`
	Deadline       string `json:"deadline" binding:"required"`
	Description    string `json:"description"`
}

type SubmitCandidaturaRequest struct {
	CandidateName    string  `json:"candidate_name" binding:"required"`
	CandidateSurname string  `json:"candidate_surname" binding:"required"`
	FiscalCode       string  `json:"fiscal_code" binding:"required"`
	Email            string  `json:"email" binding:"required,email"`
	PEC              string  `json:"pec"`
	Phone            string  `json:"phone" binding:"required"`
	GraduationGrade  float64 `json:"graduation_grade" binding:"required,gt=0"`
	GraduationLode   bool    `json:"graduation_lode"`
	HasHabilitation  bool    `json:"has_habitation"`
	MonthsOfService  int     `json:"months_of_service"`
	CertLanguage     string  `json:"cert_language"` // "B2", "C1", "C2", ""
	CertDigitalCount int     `json:"cert_digital_count"`
	CVUrl            string  `json:"cv_url"`
	DPR445Declared   bool    `json:"dpr445_declared" binding:"required"`
}

type ConvocaRequest struct {
	HoursToRespond int `json:"hours_to_respond"` // Default 24h as per O.M. 88/2024
}

type RispondiRequest struct {
	Risposta string `json:"risposta" binding:"required"` // "accettata", "rifiutata"
	Notes    string `json:"notes"`
}
