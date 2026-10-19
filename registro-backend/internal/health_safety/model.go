package health_safety

import "time"

type StudentMedication struct {
	ID                        string    `json:"id" db:"id"`
	SchoolID                  string    `json:"school_id" db:"school_id"`
	StudentID                 string    `json:"student_id" db:"student_id"`
	ClassID                   string    `json:"class_id" db:"class_id"`
	MedicationName            string    `json:"medication_name" db:"medication_name"`
	Dosage                    string    `json:"dosage" db:"dosage"`
	AdministrationMode        string    `json:"administration_mode" db:"administration_mode"`
	DiagnosisReason           string    `json:"diagnosis_reason" db:"diagnosis_reason"`
	MedicalPrescriptionURL    string    `json:"medical_prescription_url" db:"medical_prescription_url"`
	PrescriptionExpiry        string    `json:"prescription_expiry" db:"prescription_expiry"` // YYYY-MM-DD
	ParentAuthorizationSigned bool      `json:"parent_authorization_signed" db:"parent_authorization_signed"`
	AuthorizedPersonnelJSON   string    `json:"authorized_personnel_json" db:"authorized_personnel_json"`
	EmergencyProtocol         string    `json:"emergency_protocol" db:"emergency_protocol"`
	MedicationExpiryDate      string    `json:"medication_expiry_date" db:"medication_expiry_date"` // YYYY-MM-DD
	Status                    string    `json:"status" db:"status"`                                 // attivo, sospeso, scaduto
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`

	IsMedicationExpired bool `json:"is_medication_expired"`
	IsPrescriptionDue   bool `json:"is_prescription_due"`
}

type MedicationAdministration struct {
	ID                      string    `json:"id" db:"id"`
	MedicationID            string    `json:"medication_id" db:"medication_id"`
	AdministeredBy          string    `json:"administered_by" db:"administered_by"`
	AdministeredAt          time.Time `json:"administered_at" db:"administered_at"`
	DoseGiven               string    `json:"dose_given" db:"dose_given"`
	Notes                   string    `json:"notes" db:"notes"`
	EmergencyServicesCalled bool      `json:"emergency_services_called" db:"emergency_services_called"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
}

type SchoolAccident struct {
	ID                       string    `json:"id" db:"id"`
	SchoolID                 string    `json:"school_id" db:"school_id"`
	AccidentNumber           string    `json:"accident_number" db:"accident_number"`
	AccidentTime             time.Time `json:"accident_time" db:"accident_time"`
	Location                 string    `json:"location" db:"location"`
	PersonType               string    `json:"person_type" db:"person_type"` // student, teacher, ata
	PersonID                 string    `json:"person_id" db:"person_id"`
	ClassID                  *string   `json:"class_id,omitempty" db:"class_id"`
	ActivityDuringAccident   string    `json:"activity_during_accident" db:"activity_during_accident"`
	DynamicsDescription      string    `json:"dynamics_description" db:"dynamics_description"`
	WitnessesJSON            string    `json:"witnesses_json" db:"witnesses_json"`
	TeacherPresentID         *string   `json:"teacher_present_id,omitempty" db:"teacher_present_id"`
	ReportedInjuries         string    `json:"reported_injuries" db:"reported_injuries"`
	FirstAidGiven            string    `json:"first_aid_given" db:"first_aid_given"`
	AmbulanceCalled          bool      `json:"ambulance_called" db:"ambulance_called"`
	MedicalCertificateNumber string    `json:"medical_certificate_number" db:"medical_certificate_number"`
	PrognosisDays            int       `json:"prognosis_days" db:"prognosis_days"`
	InailReported            bool      `json:"inail_reported" db:"inail_reported"`
	InailProtocol            string    `json:"inail_protocol" db:"inail_protocol"`
	InailDeadline            time.Time `json:"inail_deadline" db:"inail_deadline"` // accident_time + 48h
	Status                   string    `json:"status" db:"status"`                 // aperto, certificato_ricevuto, trasmesso_inail, archiviato
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`

	HoursRemainingInail float64 `json:"hours_remaining_inail"`
	IsUrgent48hAlert    bool    `json:"is_urgent_48h_alert"`
}
