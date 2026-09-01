package attendance

import (
	"context"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAttendanceRepo struct {
	mock.Mock
}

func (m *MockAttendanceRepo) Create(att *Attendance) error {
	args := m.Called(att)
	att.ID = "att-1"
	return args.Error(0)
}
func (m *MockAttendanceRepo) BatchCreate(atts []*Attendance) error {
	args := m.Called(atts)
	return args.Error(0)
}
func (m *MockAttendanceRepo) Update(att *Attendance) error {
	args := m.Called(att)
	return args.Error(0)
}
func (m *MockAttendanceRepo) DeleteByClassDateHour(schoolID, classID string, date time.Time, hour int) error {
	args := m.Called(schoolID, classID, date, hour)
	return args.Error(0)
}
func (m *MockAttendanceRepo) FindByID(id string) (*Attendance, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Attendance), args.Error(1)
}
func (m *MockAttendanceRepo) FindByClassAndDate(classID string, date time.Time) ([]Attendance, error) {
	args := m.Called(classID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Attendance), args.Error(1)
}
func (m *MockAttendanceRepo) FindByStudent(studentID string, startDate, endDate time.Time) ([]Attendance, error) {
	args := m.Called(studentID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Attendance), args.Error(1)
}
func (m *MockAttendanceRepo) GetStats(studentID string) (*SummaryResponse, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SummaryResponse), args.Error(1)
}
func (m *MockAttendanceRepo) GetStatsBatch(ctx context.Context, studentIDs []string) (map[string]*SummaryResponse, error) {
	return make(map[string]*SummaryResponse), nil
}
func (m *MockAttendanceRepo) CountDistinctDays(studentID string) (int, error) {
	args := m.Called(studentID)
	return args.Int(0), args.Error(1)
}
func (m *MockAttendanceRepo) GetAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AnalyticsResponse), args.Error(1)
}
func (m *MockAttendanceRepo) CreateJustification(j *Justification) error {
	args := m.Called(j)
	j.ID = "just-1"
	return args.Error(0)
}
func (m *MockAttendanceRepo) UpdateJustification(j *Justification) error {
	args := m.Called(j)
	return args.Error(0)
}
func (m *MockAttendanceRepo) ProcessJustificationTx(ctx context.Context, j *Justification, teacherID string, approve bool) error {
	args := m.Called(ctx, j, teacherID, approve)
	return args.Error(0)
}
func (m *MockAttendanceRepo) FindJustificationByID(id string) (*Justification, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Justification), args.Error(1)
}
func (m *MockAttendanceRepo) FindPendingJustifications(classID, schoolID string) ([]Justification, error) {
	args := m.Called(classID, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Justification), args.Error(1)
}
func (m *MockAttendanceRepo) FindPendingJustificationsForTeacher(ctx context.Context, teacherID, schoolID string) ([]Justification, error) {
	args := m.Called(ctx, teacherID, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Justification), args.Error(1)
}
func (m *MockAttendanceRepo) DeleteJustification(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockAttendanceRepo) DeletePendingJustification(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockAttendanceRepo) IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error) {
	args := m.Called(ctx, studentID, classID)
	return args.Bool(0), args.Error(1)
}
func (m *MockAttendanceRepo) IsClassInSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	args := m.Called(ctx, classID, schoolID)
	return args.Bool(0), args.Error(1)
}
func (m *MockAttendanceRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	args := m.Called(ctx, teacherID, classID)
	return args.Bool(0), args.Error(1)
}
func (m *MockAttendanceRepo) AreStudentsInClass(ctx context.Context, studentIDs []string, classID string) (map[string]bool, error) {
	args := m.Called(ctx, studentIDs, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]bool), args.Error(1)
}
func (m *MockAttendanceRepo) IsTeacherSubstitute(ctx context.Context, teacherID, classID string, date time.Time, hour int) (bool, error) {
	args := m.Called(ctx, teacherID, classID, date, hour)
	return args.Bool(0), args.Error(1)
}
func (m *MockAttendanceRepo) HasOverlappingJustification(ctx context.Context, studentID string, startDate, endDate time.Time) (bool, error) {
	args := m.Called(ctx, studentID, startDate, endDate)
	return args.Bool(0), args.Error(1)
}
func (m *MockAttendanceRepo) GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) ([]MonthlyBreakdownRow, error) {
	args := m.Called(ctx, studentID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MonthlyBreakdownRow), args.Error(1)
}
func (m *MockAttendanceRepo) FindUnjustifiedByStudent(studentID string) ([]Attendance, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Attendance), args.Error(1)
}
func (m *MockAttendanceRepo) JustifyAbsenceByParent(attendanceID string, studentID string, reason string, notes string) error {
	args := m.Called(attendanceID, studentID, reason, notes)
	return args.Error(0)
}
func (m *MockAttendanceRepo) GetStudentAttendanceStats(studentID string) (*AttendanceStats, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AttendanceStats), args.Error(1)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *users.User) error { return nil }
func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}
func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Update(ctx context.Context, user *users.User) error { return nil }
func (m *MockUserRepo) Delete(ctx context.Context, id string) error        { return nil }
func (m *MockUserRepo) Restore(ctx context.Context, id string) error       { return nil }
func (m *MockUserRepo) List(ctx context.Context, filter users.UserFilter) ([]users.User, int, error) {
	return nil, 0, nil
}
func (m *MockUserRepo) ListByIDs(ctx context.Context, ids []string) ([]users.User, error) {
	return nil, nil
}
func (m *MockUserRepo) LogAudit(ctx context.Context, log *users.AuditLog) error { return nil }
func (m *MockUserRepo) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]users.AuditLog, int, error) {
	return nil, 0, nil
}
func (m *MockUserRepo) BulkCreate(ctx context.Context, users []users.User) (int, []string, error) {
	return 0, nil, nil
}
func (m *MockUserRepo) HardDelete(ctx context.Context, id string) error              { return nil }
func (m *MockUserRepo) RevokeAllUserTokens(ctx context.Context, userID string) error { return nil }
func (m *MockUserRepo) ClearTempMFASecret(ctx context.Context, userID string) error  { return nil }
func (m *MockUserRepo) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserRepo) AddGuardian(ctx context.Context, studentID, parentID, relation string) error {
	return nil
}
func (m *MockUserRepo) AddPasswordHistory(ctx context.Context, userID, passwordHash string) error {
	return nil
}
func (m *MockUserRepo) BulkDelete(ctx context.Context, ids []string) (int, error) {
	return len(ids), nil
}
func (m *MockUserRepo) GetChildren(ctx context.Context, parentID string) ([]users.StudentChild, error) {
	return nil, nil
}
func (m *MockUserRepo) GetPasswordHistory(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}
func (m *MockUserRepo) GetStudentsByClass(ctx context.Context, classID string) ([]users.User, error) {
	return nil, nil
}
func (m *MockUserRepo) GetStudentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockUserRepo) GetParentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockUserRepo) RemoveGuardian(ctx context.Context, studentProfileID, parentProfileID string) error {
	return nil
}
func (m *MockUserRepo) GetGuardians(ctx context.Context, studentProfileID string) ([]users.GuardianInfo, error) {
	return nil, nil
}
func (m *MockUserRepo) GetFascicoloSummary(ctx context.Context, studentID string, isActive bool) (map[string]interface{}, error) {
	return nil, nil
}
func (m *MockUserRepo) IsActive(ctx context.Context, id string) (bool, error) {
	return true, nil
}
func (m *MockUserRepo) ChangePasswordTx(ctx context.Context, userID, newPasswordHash string) error {
	return nil
}

func TestAttendanceService_MarkAttendance(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	ctx := context.Background()
	teacherID := "teacher-1"
	schoolID := "school-1"
	classID := "class-1"
	today := time.Now().Format("2006-01-02")

	mockRepo.On("IsTeacherAssignedToClass", ctx, teacherID, classID).Return(true, nil).Once()
	mockRepo.On("IsStudentInClass", ctx, "student-1", classID).Return(true, nil).Once()
	mockRepo.On("Create", mock.Anything).Return(nil).Once()

	err := svc.MarkAttendance(ctx, teacherID, schoolID, CreateAttendanceRequest{
		StudentID: "student-1",
		ClassID:   classID,
		Date:      today,
		Hour:      1,
		Status:    StatusPresent,
	})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAttendanceService_UpdateJustifiedRejection(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	ctx := context.Background()
	attID := "att-1"
	mockRepo.On("FindByID", attID).Return(&Attendance{
		ID:        attID,
		SchoolID:  "school-1",
		ClassID:   "class-1",
		Justified: true,
	}, nil).Once()

	newStatus := StatusPresent
	err := svc.UpdateAttendance(ctx, "teacher-1", "school-1", attID, UpdateAttendanceRequest{
		Status: &newStatus,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "già stato giustificato")
}
