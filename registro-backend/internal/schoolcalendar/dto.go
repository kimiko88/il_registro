package schoolcalendar

import "time"

// SetYearRequest è il body per impostare/aggiornare l'anno scolastico.
type SetYearRequest struct {
	YearLabel string `json:"year_label" binding:"required"` // Es: "2025/2026"
	StartDate string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate   string `json:"end_date" binding:"required"`   // YYYY-MM-DD
}

// AddNonTeachingDayRequest è il body per aggiungere un giorno non didattico.
type AddNonTeachingDayRequest struct {
	Date  string `json:"date" binding:"required"` // YYYY-MM-DD
	Label string `json:"label" binding:"required"`
}

// SchoolYearResponse è la risposta DTO per i settings dell'anno scolastico.
type SchoolYearResponse struct {
	ID           string    `json:"id"`
	SchoolID     string    `json:"school_id"`
	YearLabel    string    `json:"year_label"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	TeachingDays int       `json:"teaching_days"` // Giorni scolastici effettivi
}

// NonTeachingDayResponse è la risposta DTO per un giorno non didattico.
type NonTeachingDayResponse struct {
	ID    string `json:"id"`
	Date  string `json:"date"`
	Label string `json:"label"`
}

type CreateAcademicPeriodRequest struct {
	AcademicYearID *string `json:"academic_year_id,omitempty"`
	Name           string  `json:"name" binding:"required"`       // '1° Quadrimestre', '2° Quadrimestre'
	Code           string  `json:"code"`                          // 'Q1', 'Q2'
	StartDate      string  `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate        string  `json:"end_date" binding:"required"`   // YYYY-MM-DD
	IsCurrent      bool    `json:"is_current"`
}
