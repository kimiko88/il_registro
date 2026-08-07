package attendance

import (
	"context"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository
type MockRepository struct {
	mock.Mock
}

var _ Repository = (*MockRepository)(nil)

func (m *MockRepository) Create(att *Attendance) error {
	args := m.Called(att)
	return args.Error(0)
}
func (m *MockRepository) BatchCreate(atts []*Attendance) error {
	args := m.Called(atts)
	return args.Error(0)
}
func (m *MockRepository) Update(att *Attendance) error {
	args := m.Called(att)
	return args.Error(0)
}
func (m *MockRepository) DeleteByClassDateHour(classID string, date time.Time, hour int) error {
	args := m.Called(classID, date, hour)
	return args.Error(0)
}
func (m *MockRepository) ProcessJustificationTx(ctx context.Context, j *Justification, teacherID string, approve bool) error {
	args := m.Called(ctx, j, teacherID, approve)
	return args.Error(0)
}
func (m *MockRepository) FindByID(id string) (*Attendance, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Attendance), args.Error(1)
}
func (m *MockRepository) FindByClassAndDate(classID string, date time.Time) ([]Attendance, error) {
	args := m.Called(classID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Attendance), args.Error(1)
}
func (m *MockRepository) FindByStudent(studentID string, startDate, endDate time.Time) ([]Attendance, error) {
	args := m.Called(studentID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Attendance), args.Error(1)
}
func (m *MockRepository) GetStats(studentID string) (*SummaryResponse, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SummaryResponse), args.Error(1)
}
func (m *MockRepository) CreateJustification(j *Justification) error {
	args := m.Called(j)
	return args.Error(0)
}
func (m *MockRepository) UpdateJustification(j *Justification) error {
	args := m.Called(j)
	return args.Error(0)
}
func (m *MockRepository) FindJustificationByID(id string) (*Justification, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Justification), args.Error(1)
}
func (m *MockRepository) FindPendingJustifications(classID string) ([]Justification, error) {
	args := m.Called(classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Justification), args.Error(1)
}
func (m *MockRepository) CountDistinctDays(studentID string) (int, error) {
	args := m.Called(studentID)
	return args.Int(0), args.Error(1)
}
func (m *MockRepository) GetAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AnalyticsResponse), args.Error(1)
}
func (m *MockRepository) DeleteJustification(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepository) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	args := m.Called(ctx, teacherID, classID)
	return args.Bool(0), args.Error(1)
}
func (m *MockRepository) GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) ([]MonthlyBreakdownRow, error) {
	args := m.Called(ctx, studentID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MonthlyBreakdownRow), args.Error(1)
}
func (m *MockRepository) FindUnjustifiedByStudent(studentID string) ([]Attendance, error) {
	return nil, nil
}
func (m *MockRepository) JustifyAbsenceByParent(attendanceID string, reason string, notes string) error {
	return nil
}
func (m *MockRepository) GetStudentAttendanceStats(studentID string) (*AttendanceStats, error) {
	return &AttendanceStats{}, nil
}

type MockUserRepo struct {
	users.Repository
}

func (m *MockUserRepo) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	return true, nil
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*users.User, error) {
	return &users.User{ID: id, Role: "teacher"}, nil
}

func TestMarkAttendance(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUserRepo := new(MockUserRepo)
	service := NewService(mockRepo, mockUserRepo, nil, nil)

	ctx := context.Background()
	teacherID := "t1"
	mockRepo.On("IsTeacherAssignedToClass", mock.Anything, teacherID, "c1").Return(true, nil).Maybe()

	todayStr := time.Now().Format("2006-01-02")
	hourVal := 1

	t.Run("MarkSingle_Success", func(t *testing.T) {
		req := CreateAttendanceRequest{
			StudentID: "s1",
			ClassID:   "c1",
			Date:      todayStr,
			Status:    StatusPresent,
		}

		mockRepo.On("Create", mock.MatchedBy(func(a *Attendance) bool {
			return a.StudentID == "s1" && a.Status == StatusPresent
		})).Return(nil).Once()

		err := service.MarkAttendance(ctx, teacherID, "school-1", req)
		assert.NoError(t, err)
	})

	t.Run("MarkSingle_WithTimes_Success", func(t *testing.T) {
		req := CreateAttendanceRequest{
			StudentID: "s1",
			ClassID:   "c1",
			Date:      todayStr,
			Status:    StatusLate,
			Hour:      hourVal,
			EntryTime: "08:30",
			ExitTime:  "13:00",
		}

		mockRepo.On("Create", mock.MatchedBy(func(a *Attendance) bool {
			return a.StudentID == "s1" && a.Status == StatusLate && a.EntryTime != nil && *a.EntryTime == "08:30" && a.ExitTime != nil && *a.ExitTime == "13:00"
		})).Return(nil).Once()

		err := service.MarkAttendance(ctx, teacherID, "school-1", req)
		assert.NoError(t, err)
	})

	t.Run("MarkBulk_Success", func(t *testing.T) {
		req := BulkAttendanceRequest{
			ClassID: "c1",
			Date:    todayStr,
			Statuses: []StudentStatusRequest{
				{StudentID: "s1", Status: StatusPresent},
				{StudentID: "s2", Status: StatusAbsent},
			},
		}

		mockRepo.On("BatchCreate", mock.MatchedBy(func(atts []*Attendance) bool {
			return len(atts) == 2
		})).Return(nil).Once()

		err := service.MarkBulk(ctx, teacherID, "school-1", req)
		assert.NoError(t, err)
	})
}

func TestGetClassAttendance(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUserRepo := new(MockUserRepo)
	service := NewService(mockRepo, mockUserRepo, nil, nil)
	ctx := context.Background()

	t.Run("ReturnsSummary", func(t *testing.T) {
		dateStr := "2025-10-10"
		date, _ := time.Parse("2006-01-02", dateStr)

		atts := []Attendance{
			{ID: "1", StudentID: "s1", Status: StatusPresent, Date: date},
			{ID: "2", StudentID: "s2", Status: StatusAbsent, Date: date},
			{ID: "3", StudentID: "s3", Status: StatusLate, Date: date},
		}

		mockRepo.On("IsTeacherAssignedToClass", mock.Anything, "teacher-1", "c1").Return(true, nil).Maybe()
		mockRepo.On("FindByClassAndDate", "c1", date).Return(atts, nil).Once()

		resp, err := service.GetClassAttendance(ctx, "teacher-1", "teacher", "school-1", "c1", dateStr)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Summary.Present)
		assert.Equal(t, 1, resp.Summary.Absent)
		assert.Equal(t, 1, resp.Summary.Late)
	})
}

func TestJustificationFlow(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUserRepo := new(MockUserRepo)
	service := NewService(mockRepo, mockUserRepo, nil, nil)
	ctx := context.Background()

	t.Run("RequestJustification", func(t *testing.T) {
		mockRepo.On("CreateJustification", mock.Anything).Return(nil).Once()

		err := service.RequestJustification(ctx, "p1", JustificationRequest{
			StudentID: "s1", StartDate: "2025-10-10", EndDate: "2025-10-11", Reason: "Sick",
		})
		assert.NoError(t, err)
	})

	t.Run("ApproveJustification", func(t *testing.T) {
		jID := "j1"
		jPending := &Justification{ID: jID, Status: JustificationPending}

		mockRepo.On("FindJustificationByID", jID).Return(jPending, nil).Once()
		mockRepo.On("ProcessJustificationTx", mock.Anything, jPending, "t1", true).Return(nil).Once()

		err := service.ProcessJustification(ctx, "t1", jID, true)
		assert.NoError(t, err)
	})
}
