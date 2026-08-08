package attendance

import (
	"context"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockRepo struct {
	mock.Mock
}

var _ Repository = (*MockRepo)(nil)

func (m *MockRepo) Create(a *Attendance) error {
	args := m.Called(a)
	return args.Error(0)
}
func (m *MockRepo) BatchCreate(atts []*Attendance) error {
	args := m.Called(atts)
	return args.Error(0)
}
func (m *MockRepo) Update(a *Attendance) error {
	args := m.Called(a)
	return args.Error(0)
}
func (m *MockRepo) DeleteByClassDateHour(classID string, date time.Time, hour int) error {
	args := m.Called(classID, date, hour)
	return args.Error(0)
}
func (m *MockRepo) ProcessJustificationTx(ctx context.Context, j *Justification, teacherID string, approve bool) error {
	args := m.Called(ctx, j, teacherID, approve)
	return args.Error(0)
}
func (m *MockRepo) FindByID(id string) (*Attendance, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Attendance), args.Error(1)
}
func (m *MockRepo) FindByClassAndDate(classID string, date time.Time) ([]Attendance, error) {
	args := m.Called(classID, date)
	return args.Get(0).([]Attendance), args.Error(1)
}
func (m *MockRepo) FindByStudent(studentID string, start, end time.Time) ([]Attendance, error) {
	args := m.Called(studentID, start, end)
	return args.Get(0).([]Attendance), args.Error(1)
}
func (m *MockRepo) GetStats(studentID string) (*SummaryResponse, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SummaryResponse), args.Error(1)
}
func (m *MockRepo) CreateJustification(j *Justification) error {
	args := m.Called(j)
	return args.Error(0)
}
func (m *MockRepo) UpdateJustification(j *Justification) error {
	args := m.Called(j)
	return args.Error(0)
}
func (m *MockRepo) FindJustificationByID(id string) (*Justification, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Justification), args.Error(1)
}
func (m *MockRepo) GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) ([]MonthlyBreakdownRow, error) {
	args := m.Called(ctx, studentID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MonthlyBreakdownRow), args.Error(1)
}
func (m *MockRepo) FindUnjustifiedByStudent(studentID string) ([]Attendance, error) {
	return nil, nil
}
func (m *MockRepo) JustifyAbsenceByParent(attendanceID string, reason string, notes string) error {
	return nil
}
func (m *MockRepo) GetStudentAttendanceStats(studentID string) (*AttendanceStats, error) {
	return &AttendanceStats{}, nil
}
func (m *MockRepo) FindPendingJustifications(classID string) ([]Justification, error) {
	args := m.Called(classID)
	return args.Get(0).([]Justification), args.Error(1)
}
func (m *MockRepo) CountDistinctDays(studentID string) (int, error) {
	args := m.Called(studentID)
	return args.Int(0), args.Error(1)
}
func (m *MockRepo) GetAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AnalyticsResponse), args.Error(1)
}
func (m *MockRepo) DeleteJustification(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	args := m.Called(ctx, teacherID, classID)
	return args.Bool(0), args.Error(1)
}
func (m *MockRepo) IsTeacherSubstitute(ctx context.Context, teacherID, classID string, date time.Time, hour int) (bool, error) {
	return false, nil
}
func (m *MockRepo) HasOverlappingJustification(ctx context.Context, studentID string, startDate, endDate time.Time) (bool, error) {
	return false, nil
}

// --- Tests ---

func TestValidator_ValidateEntry(t *testing.T) {
	v := NewValidator()

	t.Run("Valid Present", func(t *testing.T) {
		a := &Attendance{
			Date:   time.Now(),
			Status: StatusPresent,
		}
		assert.NoError(t, v.ValidateEntry(a))
	})

	t.Run("Future Date Error", func(t *testing.T) {
		a := &Attendance{
			Date:   time.Now().Add(48 * time.Hour),
			Status: StatusPresent,
		}
		assert.Error(t, v.ValidateEntry(a))
	})

	t.Run("Late without Time Error", func(t *testing.T) {
		a := &Attendance{
			Date:   time.Now(),
			Status: StatusLate,
			Hour:   nil,
		}
		assert.Error(t, v.ValidateEntry(a))
	})
}

func TestService_MarkAttendance(t *testing.T) {
	mockRepo := new(MockRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	req := CreateAttendanceRequest{
		StudentID: "S1", ClassID: "C1", Date: time.Now().Format("2006-01-02"), Status: StatusPresent,
	}

	mockRepo.On("IsTeacherAssignedToClass", mock.Anything, "T1", "C1").Return(true, nil).Maybe()
	mockRepo.On("Create", mock.Anything).Return(nil)

	err := svc.MarkAttendance(context.Background(), "T1", "school-1", req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_ProcessJustification(t *testing.T) {
	mockRepo := new(MockRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	jid := "J1"
	j := &Justification{ID: jid, StudentID: "S1", Status: JustificationPending}
	classID := "C1"

	mockUserRepo.On("GetByID", mock.Anything, "S1").Return(&users.User{ID: "S1", ClassID: &classID}, nil).Maybe()
	mockRepo.On("IsTeacherAssignedToClass", mock.Anything, "T1", "C1").Return(true, nil).Maybe()
	mockRepo.On("FindJustificationByID", jid).Return(j, nil)
	mockRepo.On("ProcessJustificationTx", mock.Anything, j, "T1", true).Return(nil)

	err := svc.ProcessJustification(context.Background(), "T1", jid, true)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
