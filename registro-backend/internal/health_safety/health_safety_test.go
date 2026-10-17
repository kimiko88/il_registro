package health_safety

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHealthSafetyRepo struct {
	mock.Mock
}

func (m *mockHealthSafetyRepo) CreateMedication(med *StudentMedication) error {
	args := m.Called(med)
	return args.Error(0)
}

func (m *mockHealthSafetyRepo) ListMedicationsBySchool(schoolID, classID string) ([]StudentMedication, error) {
	args := m.Called(schoolID, classID)
	return args.Get(0).([]StudentMedication), args.Error(1)
}

func (m *mockHealthSafetyRepo) GetMedication(id string) (*StudentMedication, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*StudentMedication), args.Error(1)
}

func (m *mockHealthSafetyRepo) UpdateMedicationStatus(id, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *mockHealthSafetyRepo) LogMedicationAdministration(a *MedicationAdministration) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *mockHealthSafetyRepo) ListMedicationAdministrations(medicationID string) ([]MedicationAdministration, error) {
	args := m.Called(medicationID)
	return args.Get(0).([]MedicationAdministration), args.Error(1)
}

func (m *mockHealthSafetyRepo) CreateSchoolAccident(a *SchoolAccident) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *mockHealthSafetyRepo) ListSchoolAccidents(schoolID string) ([]SchoolAccident, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]SchoolAccident), args.Error(1)
}

func (m *mockHealthSafetyRepo) GetSchoolAccident(id string) (*SchoolAccident, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SchoolAccident), args.Error(1)
}

func (m *mockHealthSafetyRepo) UpdateAccidentInailTransmission(id, inailProtocol string) error {
	args := m.Called(id, inailProtocol)
	return args.Error(0)
}

func (m *mockHealthSafetyRepo) UpdateAccidentStatus(id, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func TestHealthSafetyService_MedicationRegistration(t *testing.T) {
	repo := new(mockHealthSafetyRepo)
	svc := NewService(repo)

	repo.On("CreateMedication", mock.AnythingOfType("*health_safety.StudentMedication")).Return(nil).Once()

	med := &StudentMedication{
		SchoolID:           "sch-1",
		StudentID:          "stud-1",
		ClassID:            "cls-1",
		MedicationName:     "Salbutamolo (Ventolin)",
		Dosage:             "2 puff al bisogno",
		AdministrationMode: "inalatore",
		DiagnosisReason:    "Asma bronchiale cronica",
		PrescriptionExpiry: "2026-06-30",
	}

	err := svc.RegisterMedication(med)
	assert.NoError(t, err)
	assert.Equal(t, "attivo", med.Status)

	// Missing medication name
	err = svc.RegisterMedication(&StudentMedication{
		SchoolID:  "sch-1",
		StudentID: "stud-1",
	})
	assert.Error(t, err)

	repo.AssertExpectations(t)
}

func TestHealthSafetyService_Accident48hDeadline(t *testing.T) {
	repo := new(mockHealthSafetyRepo)
	svc := NewService(repo)

	repo.On("CreateSchoolAccident", mock.AnythingOfType("*health_safety.SchoolAccident")).Return(nil).Once()

	now := time.Now()
	acc := &SchoolAccident{
		SchoolID:            "sch-1",
		PersonID:            "stud-1",
		PersonType:          "student",
		AccidentTime:        now,
		Location:            "Palestra plesso centrale",
		DynamicsDescription: "Slogatura caviglia durante salto",
	}

	err := svc.RecordAccident(acc)
	assert.NoError(t, err)
	assert.NotEmpty(t, acc.AccidentNumber)
	assert.Equal(t, "aperto", acc.Status)
	// Check INAIL 48h deadline calculation
	expectedDeadline := now.Add(48 * time.Hour)
	assert.WithinDuration(t, expectedDeadline, acc.InailDeadline, time.Second)

	// Test INAIL transmission
	repo.On("UpdateAccidentInailTransmission", "acc-1", "INAIL-PROT-2026-99").Return(nil).Once()
	err = svc.TransmitInail("acc-1", "INAIL-PROT-2026-99")
	assert.NoError(t, err)

	// Missing protocol
	err = svc.TransmitInail("acc-1", "")
	assert.Error(t, err)

	repo.AssertExpectations(t)
}
