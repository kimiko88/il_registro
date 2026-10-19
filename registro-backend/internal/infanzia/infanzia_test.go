package infanzia

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockInfanziaRepo struct {
	mock.Mock
}

func (m *mockInfanziaRepo) CreateSection(s *InfanziaSection) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *mockInfanziaRepo) ListSections(schoolID, academicYear string) ([]InfanziaSection, error) {
	args := m.Called(schoolID, academicYear)
	return args.Get(0).([]InfanziaSection), args.Error(1)
}

func (m *mockInfanziaRepo) GetSection(id string) (*InfanziaSection, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*InfanziaSection), args.Error(1)
}

func (m *mockInfanziaRepo) SaveDailyAttendancesBatch(sectionID, date string, attendances []InfanziaDailyAttendance) error {
	args := m.Called(sectionID, date, attendances)
	return args.Error(0)
}

func (m *mockInfanziaRepo) GetDailyAttendances(sectionID, date string) ([]InfanziaDailyAttendance, error) {
	args := m.Called(sectionID, date)
	return args.Get(0).([]InfanziaDailyAttendance), args.Error(1)
}

func (m *mockInfanziaRepo) SaveObservation(obs *InfanziaExperienceObservation) error {
	args := m.Called(obs)
	return args.Error(0)
}

func (m *mockInfanziaRepo) ListObservationsByStudent(studentID string) ([]InfanziaExperienceObservation, error) {
	args := m.Called(studentID)
	return args.Get(0).([]InfanziaExperienceObservation), args.Error(1)
}

func (m *mockInfanziaRepo) CreateDiarioEntry(d *InfanziaDiarioBordo) error {
	args := m.Called(d)
	return args.Error(0)
}

func (m *mockInfanziaRepo) ListDiarioEntries(sectionID string, limit int) ([]InfanziaDiarioBordo, error) {
	args := m.Called(sectionID, limit)
	return args.Get(0).([]InfanziaDiarioBordo), args.Error(1)
}

func TestInfanziaService_SectionAndCampiEsperienza(t *testing.T) {
	repo := new(mockInfanziaRepo)
	svc := NewService(repo)

	repo.On("CreateSection", mock.AnythingOfType("*infanzia.InfanziaSection")).Return(nil).Once()

	sec := &InfanziaSection{
		SchoolID:    "sch-1",
		SectionName: "Sezione Farfalle",
		ColorHex:    "#3B82F6",
	}

	err := svc.CreateSection(sec)
	assert.NoError(t, err)

	// Record observation on valid campo di esperienza
	repo.On("SaveObservation", mock.AnythingOfType("*infanzia.InfanziaExperienceObservation")).Return(nil).Once()

	obs := &InfanziaExperienceObservation{
		StudentID:        "stud-infanzia-1",
		CampoEsperienza:  "il_corpo_e_il_movimento",
		ObservationDate:  "2026-03-10",
		ObservationNotes: "Sviluppa ottima coordinazione oculo-manuale nei percorsi motori",
		RecordedBy:       "teach-1",
	}

	err = svc.RecordObservation(obs)
	assert.NoError(t, err)
	assert.Equal(t, "In evoluzione", obs.LevelDescriptor)

	// Invalid campo di esperienza
	err = svc.RecordObservation(&InfanziaExperienceObservation{
		StudentID:        "stud-1",
		CampoEsperienza:  "materia_inesistente",
		ObservationNotes: "Test",
		RecordedBy:       "teach-1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "campo di esperienza non valido")

	repo.AssertExpectations(t)
}
