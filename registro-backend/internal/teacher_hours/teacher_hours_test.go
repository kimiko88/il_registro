package teacher_hours

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) GetOrCreateObligation(schoolID, teacherID, academicYear string) (*TeacherHourObligation, error) {
	args := m.Called(schoolID, teacherID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TeacherHourObligation), args.Error(1)
}

func (m *mockRepository) GetObligation(obligationID string) (*TeacherHourObligation, error) {
	args := m.Called(obligationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TeacherHourObligation), args.Error(1)
}

func (m *mockRepository) ListObligationsBySchool(schoolID, academicYear string) ([]TeacherHourObligation, error) {
	args := m.Called(schoolID, academicYear)
	return args.Get(0).([]TeacherHourObligation), args.Error(1)
}

func (m *mockRepository) LogHourActivity(log *TeacherHourLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *mockRepository) ListLogsByObligation(obligationID string) ([]TeacherHourLog, error) {
	args := m.Called(obligationID)
	return args.Get(0).([]TeacherHourLog), args.Error(1)
}

func (m *mockRepository) CreateShortLeave(leave *TeacherShortLeave) error {
	args := m.Called(leave)
	return args.Error(0)
}

func (m *mockRepository) ListShortLeaves(schoolID, teacherID string) ([]TeacherShortLeave, error) {
	args := m.Called(schoolID, teacherID)
	return args.Get(0).([]TeacherShortLeave), args.Error(1)
}

func (m *mockRepository) UpdateShortLeaveStatus(id, status string, notes string) error {
	args := m.Called(id, status, notes)
	return args.Error(0)
}

func (m *mockRepository) RecordShortLeaveRecovery(id string, recoveredHours float64, recoveredDate string) error {
	args := m.Called(id, recoveredHours, recoveredDate)
	return args.Error(0)
}

func TestTeacherHoursService_ShortLeave_CCNLMax2Hours(t *testing.T) {
	mockRepo := new(mockRepository)
	svc := NewService(mockRepo)

	// Valid 1.5h leave request
	mockRepo.On("CreateShortLeave", mock.AnythingOfType("*teacher_hours.TeacherShortLeave")).Return(nil).Once()
	leave, err := svc.RequestShortLeave("sch-1", "teach-1", ShortLeaveRequest{
		RequestDate: "2026-03-10",
		Hours:       1.5,
		Reason:      "Visita medica",
	})
	assert.NoError(t, err)
	assert.NotNil(t, leave)
	assert.Equal(t, 1.5, leave.Hours)

	// Invalid > 2h leave request (must be rejected by CCNL rules)
	_, err = svc.RequestShortLeave("sch-1", "teach-1", ShortLeaveRequest{
		RequestDate: "2026-03-10",
		Hours:       3.0,
		Reason:      "Superamento 2 ore CCNL",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "al massimo di 2 ore")

	// Invalid <= 0h
	_, err = svc.RequestShortLeave("sch-1", "teach-1", ShortLeaveRequest{
		RequestDate: "2026-03-10",
		Hours:       0,
		Reason:      "Zero ore",
	})
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestTeacherHoursService_LogActivity_Validations(t *testing.T) {
	mockRepo := new(mockRepository)
	svc := NewService(mockRepo)

	// Valid activity logging
	mockRepo.On("LogHourActivity", mock.AnythingOfType("*teacher_hours.TeacherHourLog")).Return(nil).Once()
	log, err := svc.LogHourActivity("teach-1", LogHourActivityRequest{
		ObligationID: "obl-1",
		Category:     "collegio",
		Hours:        2.0,
		ActivityDate: "2026-03-15",
		Description:  "Collegio Docenti n. 3",
	})
	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, "completed", log.Status)

	// Invalid category
	_, err = svc.LogHourActivity("teach-1", LogHourActivityRequest{
		ObligationID: "obl-1",
		Category:     "invalid_cat",
		Hours:        2.0,
		ActivityDate: "2026-03-15",
		Description:  "Invalid",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "categoria non valida")

	// Negative hours
	_, err = svc.LogHourActivity("teach-1", LogHourActivityRequest{
		ObligationID: "obl-1",
		Category:     "collegio",
		Hours:        -1.0,
		ActivityDate: "2026-03-15",
		Description:  "Invalid",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maggiori di zero")

	mockRepo.AssertExpectations(t)
}

func TestTeacherHoursService_DashboardEnrichment(t *testing.T) {
	mockRepo := new(mockRepository)
	svc := NewService(mockRepo)

	obl := &TeacherHourObligation{
		ID:                     "obl-123",
		SchoolID:               "sch-1",
		TeacherID:              "teach-1",
		AcademicYear:           "2025/2026",
		CollegioPlannedHours:   40.0,
		ConsigliPlannedHours:   40.0,
		CollegioDoneHours:      15.0,
		CollegioRemainingHours: 25.0,
		ConsigliDoneHours:      10.0,
		ConsigliRemainingHours: 30.0,
	}

	logs := []TeacherHourLog{
		{ID: "log-1", Category: "collegio", Hours: 3.0, Description: "Collegio di inizio anno"},
	}

	leaves := []TeacherShortLeave{
		{ID: "leave-1", Hours: 1.0, Reason: "Permesso breve", Status: "approved"},
	}

	mockRepo.On("GetOrCreateObligation", "sch-1", "teach-1", "2025/2026").Return(obl, nil).Once()
	mockRepo.On("ListLogsByObligation", "obl-123").Return(logs, nil).Once()
	mockRepo.On("ListShortLeaves", "sch-1", "teach-1").Return(leaves, nil).Once()

	resObl, resLogs, resLeaves, err := svc.GetTeacherDashboard("sch-1", "teach-1", "2025/2026")
	assert.NoError(t, err)
	assert.Equal(t, "obl-123", resObl.ID)
	assert.Len(t, resLogs, 1)
	assert.Len(t, resLeaves, 1)

	// Missing parameters
	_, _, _, err = svc.GetTeacherDashboard("", "teach-1", "2025/2026")
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestTeacherHoursService_ReviewAndRecovery(t *testing.T) {
	mockRepo := new(mockRepository)
	svc := NewService(mockRepo)

	mockRepo.On("UpdateShortLeaveStatus", "leave-1", "approved", "Autorizzato dal DS").Return(nil).Once()
	err := svc.ReviewShortLeave("leave-1", "approved", "Autorizzato dal DS")
	assert.NoError(t, err)

	// Invalid status
	err = svc.ReviewShortLeave("leave-1", "arbitrary_status", "")
	assert.Error(t, err)

	// Recovery
	mockRepo.On("RecordShortLeaveRecovery", "leave-1", 1.5, "2026-03-20").Return(nil).Once()
	err = svc.RecordLeaveRecovery("leave-1", RecoverLeaveRequest{
		RecoveredHours: 1.5,
		RecoveredDate:  "2026-03-20",
	})
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
