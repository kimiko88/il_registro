package schoolcalendar

import "time"

// SchoolYearSettings contiene le date di inizio e fine dell'anno scolastico
// per una specifica scuola.
type SchoolYearSettings struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	YearLabel string    `json:"year_label" db:"year_label"` // Es: "2025/2026"
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NonTeachingDay rappresenta un singolo giorno non didattico (festività, ponti,
// sospensione delle lezioni).
type NonTeachingDay struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Date      time.Time `json:"date" db:"date"`
	Label     string    `json:"label" db:"label"` // Es: "Festa della Repubblica"
	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
