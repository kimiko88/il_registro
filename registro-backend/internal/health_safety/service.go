package health_safety

import (
	"errors"
	"fmt"
	"time"
)

type Service interface {
	RegisterMedication(m *StudentMedication) error
	ListMedications(schoolID, classID string) ([]StudentMedication, error)
	GetMedication(id string) (*StudentMedication, error)
	LogAdministration(a *MedicationAdministration) error
	ListAdministrations(medicationID string) ([]MedicationAdministration, error)

	RecordAccident(a *SchoolAccident) error
	ListAccidents(schoolID string) ([]SchoolAccident, error)
	GetAccident(id string) (*SchoolAccident, error)
	TransmitInail(accidentID, protocol string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterMedication(m *StudentMedication) error {
	if m == nil {
		return errors.New("dati farmaco non validi")
	}
	if m.SchoolID == "" || m.StudentID == "" || m.MedicationName == "" {
		return errors.New("schoolID, studentID e medication_name sono obbligatori")
	}
	if m.Dosage == "" || m.AdministrationMode == "" {
		return errors.New("dosaggio e modalità di somministrazione sono obbligatori")
	}
	if m.Status == "" {
		m.Status = "attivo"
	}
	return s.repo.CreateMedication(m)
}

func (s *service) ListMedications(schoolID, classID string) ([]StudentMedication, error) {
	if schoolID == "" {
		return nil, errors.New("schoolID obbligatorio")
	}
	return s.repo.ListMedicationsBySchool(schoolID, classID)
}

func (s *service) GetMedication(id string) (*StudentMedication, error) {
	if id == "" {
		return nil, errors.New("id farmaco obbligatorio")
	}
	return s.repo.GetMedication(id)
}

func (s *service) LogAdministration(a *MedicationAdministration) error {
	if a == nil {
		return errors.New("dati somministrazione non validi")
	}
	if a.MedicationID == "" || a.AdministeredBy == "" || a.DoseGiven == "" {
		return errors.New("medicationID, administeredBy e doseGiven sono obbligatori")
	}
	return s.repo.LogMedicationAdministration(a)
}

func (s *service) ListAdministrations(medicationID string) ([]MedicationAdministration, error) {
	if medicationID == "" {
		return nil, errors.New("medicationID obbligatorio")
	}
	return s.repo.ListMedicationAdministrations(medicationID)
}

func (s *service) RecordAccident(a *SchoolAccident) error {
	if a == nil {
		return errors.New("dati infortunio non validi")
	}
	if a.SchoolID == "" || a.PersonID == "" || a.DynamicsDescription == "" {
		return errors.New("schoolID, personID e dynamicsDescription sono obbligatori")
	}
	if a.AccidentTime.IsZero() {
		a.AccidentTime = time.Now()
	}
	if a.AccidentNumber == "" {
		a.AccidentNumber = fmt.Sprintf("INF-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)
	}
	if a.Status == "" {
		a.Status = "aperto"
	}
	// INAIL legal limit: 48 hours
	a.InailDeadline = a.AccidentTime.Add(48 * time.Hour)
	return s.repo.CreateSchoolAccident(a)
}

func (s *service) ListAccidents(schoolID string) ([]SchoolAccident, error) {
	if schoolID == "" {
		return nil, errors.New("schoolID obbligatorio")
	}
	return s.repo.ListSchoolAccidents(schoolID)
}

func (s *service) GetAccident(id string) (*SchoolAccident, error) {
	if id == "" {
		return nil, errors.New("id infortunio obbligatorio")
	}
	return s.repo.GetSchoolAccident(id)
}

func (s *service) TransmitInail(accidentID, protocol string) error {
	if accidentID == "" || protocol == "" {
		return errors.New("accidentID e protocollo di trasmissione INAIL sono obbligatori")
	}
	return s.repo.UpdateAccidentInailTransmission(accidentID, protocol)
}
