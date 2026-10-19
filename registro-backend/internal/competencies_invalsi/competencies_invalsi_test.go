package competencies_invalsi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCompInvalsiRepo struct {
	mock.Mock
}

func (m *mockCompInvalsiRepo) UpsertCertification(c *CompetencyCertificationDM14) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *mockCompInvalsiRepo) GetCertification(id string) (*CompetencyCertificationDM14, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CompetencyCertificationDM14), args.Error(1)
}

func (m *mockCompInvalsiRepo) GetCertificationByStudent(studentID, schoolYear, schoolTier string) (*CompetencyCertificationDM14, error) {
	args := m.Called(studentID, schoolYear, schoolTier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CompetencyCertificationDM14), args.Error(1)
}

func (m *mockCompInvalsiRepo) ListCertificationsByClass(classID, schoolYear string) ([]CompetencyCertificationDM14, error) {
	args := m.Called(classID, schoolYear)
	return args.Get(0).([]CompetencyCertificationDM14), args.Error(1)
}

func (m *mockCompInvalsiRepo) SignCertification(id, principalID string) error {
	args := m.Called(id, principalID)
	return args.Error(0)
}

func (m *mockCompInvalsiRepo) InsertInvalsiBatch(results []InvalsiResult) error {
	args := m.Called(results)
	return args.Error(0)
}

func (m *mockCompInvalsiRepo) ListInvalsiByClass(classID, schoolYear string) ([]InvalsiResult, error) {
	args := m.Called(classID, schoolYear)
	return args.Get(0).([]InvalsiResult), args.Error(1)
}

func (m *mockCompInvalsiRepo) GetSchoolInvalsiAnalytics(schoolID, schoolYear string) ([]InvalsiResult, error) {
	args := m.Called(schoolID, schoolYear)
	return args.Get(0).([]InvalsiResult), args.Error(1)
}

func TestCompetenciesInvalsiService_DM14Validation(t *testing.T) {
	repo := new(mockCompInvalsiRepo)
	svc := NewService(repo)

	repo.On("UpsertCertification", mock.AnythingOfType("*competencies_invalsi.CompetencyCertificationDM14")).Return(nil).Once()

	cert := &CompetencyCertificationDM14{
		SchoolID:                "sch-1",
		StudentID:               "stud-1",
		ClassID:                 "cls-3a",
		SchoolTier:              "primo_ciclo",
		SchoolYear:              "2025/2026",
		CouncilDeliberationDate: "2026-06-12",
	}

	err := svc.SaveCertification(cert)
	assert.NoError(t, err)
	assert.Equal(t, "bozza", cert.Status)

	// Invalid school tier
	err = svc.SaveCertification(&CompetencyCertificationDM14{
		SchoolID:   "sch-1",
		StudentID:  "stud-1",
		ClassID:    "cls-3a",
		SchoolTier: "universita", // invalid
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "school_tier non valido")

	repo.AssertExpectations(t)
}

func TestCompetenciesInvalsiService_InvalsiBatchValidation(t *testing.T) {
	repo := new(mockCompInvalsiRepo)
	svc := NewService(repo)

	repo.On("InsertInvalsiBatch", mock.AnythingOfType("[]competencies_invalsi.InvalsiResult")).Return(nil).Once()

	results := []InvalsiResult{
		{
			SchoolID:         "sch-1",
			ClassID:          "cls-1",
			SchoolYear:       "2025/2026",
			GradeLevel:       8,
			Subject:          "italiano",
			ScoreWLE:         215.5,
			Level:            4,
			ReportClassGrade: 7.5,
		},
	}

	err := svc.ImportInvalsiBatch(results)
	assert.NoError(t, err)

	// Invalid Level (> 5)
	invalidResults := []InvalsiResult{
		{Level: 6},
	}
	err = svc.ImportInvalsiBatch(invalidResults)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "livello INVALSI non valido")

	repo.AssertExpectations(t)
}
