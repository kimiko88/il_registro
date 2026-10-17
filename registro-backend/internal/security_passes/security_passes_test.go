package security_passes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSecurityPassesRepo struct {
	mock.Mock
}

func (m *mockSecurityPassesRepo) CreateHallPass(p *StudentHallPass) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *mockSecurityPassesRepo) CloseHallPass(id string, returnedAt time.Time) error {
	args := m.Called(id, returnedAt)
	return args.Error(0)
}

func (m *mockSecurityPassesRepo) ListActiveHallPassesBySchool(schoolID string) ([]StudentHallPass, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]StudentHallPass), args.Error(1)
}

func (m *mockSecurityPassesRepo) ListHallPassesByClass(classID string, limit int) ([]StudentHallPass, error) {
	args := m.Called(classID, limit)
	return args.Get(0).([]StudentHallPass), args.Error(1)
}

func (m *mockSecurityPassesRepo) CreatePickupDelegate(d *StudentPickupDelegate) error {
	args := m.Called(d)
	return args.Error(0)
}

func (m *mockSecurityPassesRepo) ListDelegatesByStudent(studentID string) ([]StudentPickupDelegate, error) {
	args := m.Called(studentID)
	return args.Get(0).([]StudentPickupDelegate), args.Error(1)
}

func (m *mockSecurityPassesRepo) RevokeDelegate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestSecurityPassesService_HallPassFlow(t *testing.T) {
	repo := new(mockSecurityPassesRepo)
	svc := NewService(repo)

	repo.On("CreateHallPass", mock.AnythingOfType("*security_passes.StudentHallPass")).Return(nil).Once()

	pass := &StudentHallPass{
		SchoolID:  "sch-1",
		ClassID:   "cls-1",
		StudentID: "stud-1",
		TeacherID: "teach-1",
		PassType:  "bagno",
	}

	err := svc.IssueHallPass(pass)
	assert.NoError(t, err)
	assert.Equal(t, "fuori_aula", pass.Status)

	// Close pass
	repo.On("CloseHallPass", "pass-1", mock.AnythingOfType("time.Time")).Return(nil).Once()
	err = svc.ReturnHallPass("pass-1")
	assert.NoError(t, err)

	// Validation: missing classID
	err = svc.IssueHallPass(&StudentHallPass{
		SchoolID: "sch-1",
	})
	assert.Error(t, err)

	repo.AssertExpectations(t)
}

func TestSecurityPassesService_Delegates(t *testing.T) {
	repo := new(mockSecurityPassesRepo)
	svc := NewService(repo)

	repo.On("CreatePickupDelegate", mock.AnythingOfType("*security_passes.StudentPickupDelegate")).Return(nil).Once()

	del := &StudentPickupDelegate{
		StudentID:        "stud-1",
		DelegateName:     "Mario Rossi (Nonno)",
		DelegateRelation: "nonno",
		FiscalCode:       "RSSMRA50A01H501Z",
		Phone:            "3331234567",
		IDCardNumber:     "CA12345AA",
	}

	err := svc.AddPickupDelegate(del)
	assert.NoError(t, err)
	assert.Equal(t, "attivo", del.Status)

	// Missing phone
	err = svc.AddPickupDelegate(&StudentPickupDelegate{
		StudentID:    "stud-1",
		DelegateName: "Nonno",
		FiscalCode:   "RSSMRA50A01H501Z",
	})
	assert.Error(t, err)

	repo.AssertExpectations(t)
}
