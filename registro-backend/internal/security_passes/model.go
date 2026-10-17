package security_passes

import "time"

type StudentHallPass struct {
	ID              string     `json:"id" db:"id"`
	SchoolID        string     `json:"school_id" db:"school_id"`
	ClassID         string     `json:"class_id" db:"class_id"`
	StudentID       string     `json:"student_id" db:"student_id"`
	TeacherID       string     `json:"teacher_id" db:"teacher_id"`
	PassType        string     `json:"pass_type" db:"pass_type"` // bagno, infermeria, segreteria_presidenza, collaboratori_ata, altro
	DepartedAt      time.Time  `json:"departed_at" db:"departed_at"`
	ReturnedAt      *time.Time `json:"returned_at,omitempty" db:"returned_at"`
	DurationMinutes int        `json:"duration_minutes" db:"duration_minutes"`
	Status          string     `json:"status" db:"status"` // fuori_aula, rientrato, ritardo_anomalo
	Notes           string     `json:"notes" db:"notes"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`

	StudentFullName string `json:"student_full_name,omitempty"`
	ClassName       string `json:"class_name,omitempty"`
	ElapsedMinutes  int    `json:"elapsed_minutes"`
	IsAnomalyAlert  bool   `json:"is_anomaly_alert"` // e.g. out for > 15 minutes
}

type StudentPickupDelegate struct {
	ID               string    `json:"id" db:"id"`
	SchoolID         string    `json:"school_id" db:"school_id"`
	StudentID        string    `json:"student_id" db:"student_id"`
	DelegateName     string    `json:"delegate_name" db:"delegate_name"`
	DelegateRelation string    `json:"delegate_relation" db:"delegate_relation"` // nonno, zio, babysitter, altro
	FiscalCode       string    `json:"fiscal_code" db:"fiscal_code"`
	Phone            string    `json:"phone" db:"phone"`
	IDCardNumber     string    `json:"id_card_number" db:"id_card_number"`
	IDCardFileURL    string    `json:"id_card_file_url" db:"id_card_file_url"`
	IsPermanent      bool      `json:"is_permanent" db:"is_permanent"`
	ValidUntil       *string   `json:"valid_until,omitempty" db:"valid_until"` // YYYY-MM-DD
	Status           string    `json:"status" db:"status"`                     // attivo, revocato
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`

	StudentFullName string `json:"student_full_name,omitempty"`
}
