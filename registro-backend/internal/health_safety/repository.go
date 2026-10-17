package health_safety

import (
	"database/sql"
	"time"
)

type Repository interface {
	CreateMedication(m *StudentMedication) error
	ListMedicationsBySchool(schoolID, classID string) ([]StudentMedication, error)
	GetMedication(id string) (*StudentMedication, error)
	UpdateMedicationStatus(id, status string) error
	LogMedicationAdministration(a *MedicationAdministration) error
	ListMedicationAdministrations(medicationID string) ([]MedicationAdministration, error)

	CreateSchoolAccident(a *SchoolAccident) error
	ListSchoolAccidents(schoolID string) ([]SchoolAccident, error)
	GetSchoolAccident(id string) (*SchoolAccident, error)
	UpdateAccidentInailTransmission(id, inailProtocol string) error
	UpdateAccidentStatus(id, status string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateMedication(m *StudentMedication) error {
	query := `
		INSERT INTO student_medications (
			school_id, student_id, class_id, medication_name, dosage, administration_mode,
			diagnosis_reason, medical_prescription_url, prescription_expiry, parent_authorization_signed,
			authorized_personnel_json, emergency_protocol, medication_expiry_date, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query, m.SchoolID, m.StudentID, m.ClassID, m.MedicationName, m.Dosage, m.AdministrationMode,
		m.DiagnosisReason, m.MedicalPrescriptionURL, m.PrescriptionExpiry, m.ParentAuthorizationSigned,
		m.AuthorizedPersonnelJSON, m.EmergencyProtocol, m.MedicationExpiryDate, m.Status,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *repository) ListMedicationsBySchool(schoolID, classID string) ([]StudentMedication, error) {
	query := `
		SELECT id, school_id, student_id, class_id, medication_name, dosage, administration_mode,
		       diagnosis_reason, COALESCE(medical_prescription_url, ''),
		       TO_CHAR(prescription_expiry, 'YYYY-MM-DD'), parent_authorization_signed,
		       authorized_personnel_json::text, emergency_protocol,
		       TO_CHAR(medication_expiry_date, 'YYYY-MM-DD'), status, created_at, updated_at
		FROM student_medications
		WHERE school_id = $1 AND ($2 = '' OR class_id = $2::uuid)
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, schoolID, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meds []StudentMedication
	today := time.Now().Format("2006-01-02")
	for rows.Next() {
		var m StudentMedication
		if err := rows.Scan(
			&m.ID, &m.SchoolID, &m.StudentID, &m.ClassID, &m.MedicationName, &m.Dosage, &m.AdministrationMode,
			&m.DiagnosisReason, &m.MedicalPrescriptionURL, &m.PrescriptionExpiry, &m.ParentAuthorizationSigned,
			&m.AuthorizedPersonnelJSON, &m.EmergencyProtocol, &m.MedicationExpiryDate, &m.Status,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if m.MedicationExpiryDate < today {
			m.IsMedicationExpired = true
		}
		if m.PrescriptionExpiry < today {
			m.IsPrescriptionDue = true
		}
		meds = append(meds, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return meds, nil
}

func (r *repository) GetMedication(id string) (*StudentMedication, error) {
	query := `
		SELECT id, school_id, student_id, class_id, medication_name, dosage, administration_mode,
		       diagnosis_reason, COALESCE(medical_prescription_url, ''),
		       TO_CHAR(prescription_expiry, 'YYYY-MM-DD'), parent_authorization_signed,
		       authorized_personnel_json::text, emergency_protocol,
		       TO_CHAR(medication_expiry_date, 'YYYY-MM-DD'), status, created_at, updated_at
		FROM student_medications
		WHERE id = $1
	`
	var m StudentMedication
	err := r.db.QueryRow(query, id).Scan(
		&m.ID, &m.SchoolID, &m.StudentID, &m.ClassID, &m.MedicationName, &m.Dosage, &m.AdministrationMode,
		&m.DiagnosisReason, &m.MedicalPrescriptionURL, &m.PrescriptionExpiry, &m.ParentAuthorizationSigned,
		&m.AuthorizedPersonnelJSON, &m.EmergencyProtocol, &m.MedicationExpiryDate, &m.Status,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	today := time.Now().Format("2006-01-02")
	if m.MedicationExpiryDate < today {
		m.IsMedicationExpired = true
	}
	if m.PrescriptionExpiry < today {
		m.IsPrescriptionDue = true
	}
	return &m, nil
}

func (r *repository) UpdateMedicationStatus(id, status string) error {
	query := `UPDATE student_medications SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id, status)
	return err
}

func (r *repository) LogMedicationAdministration(a *MedicationAdministration) error {
	query := `
		INSERT INTO medication_administrations (
			medication_id, administered_by, dose_given, notes, emergency_services_called
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, administered_at, created_at
	`
	return r.db.QueryRow(
		query, a.MedicationID, a.AdministeredBy, a.DoseGiven, a.Notes, a.EmergencyServicesCalled,
	).Scan(&a.ID, &a.AdministeredAt, &a.CreatedAt)
}

func (r *repository) ListMedicationAdministrations(medicationID string) ([]MedicationAdministration, error) {
	query := `
		SELECT id, medication_id, administered_by, administered_at, dose_given, COALESCE(notes, ''),
		       emergency_services_called, created_at
		FROM medication_administrations
		WHERE medication_id = $1
		ORDER BY administered_at DESC
	`
	rows, err := r.db.Query(query, medicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []MedicationAdministration
	for rows.Next() {
		var a MedicationAdministration
		if err := rows.Scan(
			&a.ID, &a.MedicationID, &a.AdministeredBy, &a.AdministeredAt, &a.DoseGiven, &a.Notes,
			&a.EmergencyServicesCalled, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) CreateSchoolAccident(a *SchoolAccident) error {
	// Set INAIL deadline: 48h from accident time
	if a.InailDeadline.IsZero() {
		a.InailDeadline = a.AccidentTime.Add(48 * time.Hour)
	}
	query := `
		INSERT INTO school_accidents (
			school_id, accident_number, accident_time, location, person_type, person_id, class_id,
			activity_during_accident, dynamics_description, witnesses_json, teacher_present_id,
			reported_injuries, first_aid_given, ambulance_called, medical_certificate_number,
			prognosis_days, inail_reported, inail_protocol, inail_deadline, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		) RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query, a.SchoolID, a.AccidentNumber, a.AccidentTime, a.Location, a.PersonType, a.PersonID, a.ClassID,
		a.ActivityDuringAccident, a.DynamicsDescription, a.WitnessesJSON, a.TeacherPresentID,
		a.ReportedInjuries, a.FirstAidGiven, a.AmbulanceCalled, a.MedicalCertificateNumber,
		a.PrognosisDays, a.InailReported, a.InailProtocol, a.InailDeadline, a.Status,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *repository) ListSchoolAccidents(schoolID string) ([]SchoolAccident, error) {
	query := `
		SELECT id, school_id, accident_number, accident_time, location, person_type, person_id,
		       class_id, activity_during_accident, dynamics_description, witnesses_json::text,
		       teacher_present_id, reported_injuries, first_aid_given, ambulance_called,
		       COALESCE(medical_certificate_number, ''), prognosis_days, inail_reported,
		       COALESCE(inail_protocol, ''), inail_deadline, status, created_at, updated_at
		FROM school_accidents
		WHERE school_id = $1
		ORDER BY accident_time DESC
	`
	rows, err := r.db.Query(query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accidents []SchoolAccident
	now := time.Now()
	for rows.Next() {
		var a SchoolAccident
		if err := rows.Scan(
			&a.ID, &a.SchoolID, &a.AccidentNumber, &a.AccidentTime, &a.Location, &a.PersonType, &a.PersonID,
			&a.ClassID, &a.ActivityDuringAccident, &a.DynamicsDescription, &a.WitnessesJSON,
			&a.TeacherPresentID, &a.ReportedInjuries, &a.FirstAidGiven, &a.AmbulanceCalled,
			&a.MedicalCertificateNumber, &a.PrognosisDays, &a.InailReported, &a.InailProtocol,
			&a.InailDeadline, &a.Status, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		a.HoursRemainingInail = a.InailDeadline.Sub(now).Hours()
		if !a.InailReported && a.HoursRemainingInail > 0 && a.HoursRemainingInail <= 48.0 {
			a.IsUrgent48hAlert = true
		}
		accidents = append(accidents, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return accidents, nil
}

func (r *repository) GetSchoolAccident(id string) (*SchoolAccident, error) {
	query := `
		SELECT id, school_id, accident_number, accident_time, location, person_type, person_id,
		       class_id, activity_during_accident, dynamics_description, witnesses_json::text,
		       teacher_present_id, reported_injuries, first_aid_given, ambulance_called,
		       COALESCE(medical_certificate_number, ''), prognosis_days, inail_reported,
		       COALESCE(inail_protocol, ''), inail_deadline, status, created_at, updated_at
		FROM school_accidents
		WHERE id = $1
	`
	var a SchoolAccident
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.SchoolID, &a.AccidentNumber, &a.AccidentTime, &a.Location, &a.PersonType, &a.PersonID,
		&a.ClassID, &a.ActivityDuringAccident, &a.DynamicsDescription, &a.WitnessesJSON,
		&a.TeacherPresentID, &a.ReportedInjuries, &a.FirstAidGiven, &a.AmbulanceCalled,
		&a.MedicalCertificateNumber, &a.PrognosisDays, &a.InailReported, &a.InailProtocol,
		&a.InailDeadline, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	a.HoursRemainingInail = a.InailDeadline.Sub(now).Hours()
	if !a.InailReported && a.HoursRemainingInail > 0 && a.HoursRemainingInail <= 48.0 {
		a.IsUrgent48hAlert = true
	}
	return &a, nil
}

func (r *repository) UpdateAccidentInailTransmission(id, inailProtocol string) error {
	query := `
		UPDATE school_accidents
		SET inail_reported = TRUE, inail_protocol = $2, status = 'trasmesso_inail', updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, inailProtocol)
	return err
}

func (r *repository) UpdateAccidentStatus(id, status string) error {
	query := `UPDATE school_accidents SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id, status)
	return err
}
