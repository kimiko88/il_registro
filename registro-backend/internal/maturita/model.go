package maturita

import (
	"encoding/xml"
	"time"
)

type Commissioner struct {
	Name         string `json:"name"`
	Role         string `json:"role"` // "presidente", "commissario_interno", "commissario_esterno"
	Subject      string `json:"subject"`
	OriginSchool string `json:"origin_school,omitempty"`
}

type CommissioneMaturita struct {
	ID                 string         `json:"id" db:"id"`
	SchoolID           string         `json:"school_id" db:"school_id"`
	ClassID            string         `json:"class_id" db:"class_id"`
	SchoolYear         string         `json:"school_year" db:"school_year"`
	CommissionCode     string         `json:"commission_code" db:"commission_code"`
	PresidentName      string         `json:"president_name" db:"president_name"`
	PresidentUSRDecree string         `json:"president_usr_decree" db:"president_usr_decree"`
	Commissioners      []Commissioner `json:"commissioners"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
}

type TrienniumCredits struct {
	Grade3rd     float64 `json:"grade_3rd"`     // Average (media voti) 3rd year
	Credit3rd    int     `json:"credit_3rd"`    // Max 12
	Grade4th     float64 `json:"grade_4th"`     // Average 4th year
	Credit4th    int     `json:"credit_4th"`    // Max 13
	Grade5th     float64 `json:"grade_5th"`     // Average 5th year
	Credit5th    int     `json:"credit_5th"`    // Max 15
	TotalCredits int     `json:"total_credits"` // Max 40
}

type ExamScores struct {
	Written1Score float64 `json:"written1_score"` // Max 20 (Italiano)
	Written2Score float64 `json:"written2_score"` // Max 20 (Indirizzo)
	OralScore     float64 `json:"oral_score"`     // Max 20 (Colloquio)
	ExamTotal     float64 `json:"exam_total"`     // Max 60
	BonusEligible bool    `json:"bonus_eligible"` // TotalCredits >= 30 AND ExamTotal >= 50
	BonusPoints   int     `json:"bonus_points"`   // 1 to 5 points
	FinalScore    int     `json:"final_score"`    // Max 100
	Lode          bool    `json:"lode"`           // TotalCredits == 40 AND ExamTotal == 60 AND Bonus == 0
}

type StudentMaturitaRecord struct {
	ID                 string              `json:"id" db:"id"`
	SchoolID           string              `json:"school_id" db:"school_id"`
	StudentID          string              `json:"student_id" db:"student_id"`
	ClassID            string              `json:"class_id" db:"class_id"`
	StudentName        string              `json:"student_name,omitempty"`
	SchoolYear         string              `json:"school_year" db:"school_year"`
	Credits            TrienniumCredits    `json:"credits"`
	Scores             ExamScores          `json:"scores"`
	CurriculumStudente *CurriculumStudente `json:"curriculum_studente,omitempty"`
	Status             string              `json:"status" db:"status"` // "in_corso", "approvato_commissione", "diplomato"
	UpdatedAt          time.Time           `json:"updated_at" db:"updated_at"`
}

// Curriculum dello Studente (D.M. 88/2020)
type CurriculumStudente struct {
	StudentCF             string              `json:"student_cf"`
	StudentFullName       string              `json:"student_full_name"`
	DiplomaTitle          string              `json:"diploma_title"`
	FinalScore            int                 `json:"final_score"`
	Lode                  bool                `json:"lode"`
	CreditsTriennium      int                 `json:"credits_triennium"`
	PCTOHoursTotal        int                 `json:"pcto_hours_total"`
	PCTOExperiences       []PCTOExperience    `json:"pcto_experiences"`
	CertificationsLang    []CertificationLang `json:"certifications_lang"`
	CertificationsDigital []CertificationDig  `json:"certifications_digital"`
	ExtraActivities       []ExtraActivity     `json:"extra_activities"`
}

type PCTOExperience struct {
	HostCompany    string `json:"host_company"`
	ActivityTitle  string `json:"activity_title"`
	HoursCompleted int    `json:"hours_completed"`
	SkillsAcquired string `json:"skills_acquired"`
}

type CertificationLang struct {
	Language string `json:"language"`
	Level    string `json:"level"`     // B1, B2, C1, C2
	CertBody string `json:"cert_body"` // Cambridge, DELF, Goethe, Cervantes
	Year     int    `json:"year"`
}

type CertificationDig struct {
	Title    string `json:"title"` // ICDL, Cisco, EIPASS, Python
	CertBody string `json:"cert_body"`
	Year     int    `json:"year"`
}

type ExtraActivity struct {
	Category    string `json:"category"` // "volontariato", "sport_agonistico", "artistica_musicale"
	Description string `json:"description"`
}

// Ministerial XML Schema for Curriculum dello Studente
type MinisterialCurriculumXML struct {
	XMLName         xml.Name `xml:"curriculum_dello_studente"`
	CodiceFiscale   string   `xml:"studente>codice_fiscale"`
	Nominativo      string   `xml:"studente>nominativo"`
	EsitoMaturita   int      `xml:"esito>voto_finale"`
	LodeConferita   bool     `xml:"esito>lode"`
	CreditiTriennio int      `xml:"esito>crediti_scolastici"`
	OrePCTO         int      `xml:"pcto>ore_totali"`
}
