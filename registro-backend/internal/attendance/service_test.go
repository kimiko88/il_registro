package attendance

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository
type MockRepository struct {
	mock.Mock
}

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

func TestMarkAttendance(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, nil, nil)

	ctx := context.Background()
	teacherID := "t1"

	t.Run("MarkSingle_Success", func(t *testing.T) {
		req := CreateAttendanceRequest{
			StudentID: "s1",
			ClassID:   "c1",
			Date:      "2025-10-10",
			Status:    StatusPresent,
		}

		mockRepo.On("Create", mock.MatchedBy(func(a *Attendance) bool {
			return a.StudentID == "s1" && a.Status == StatusPresent
		})).Return(nil).Once()

		err := service.MarkAttendance(ctx, teacherID, req)
		assert.NoError(t, err)
	})

	t.Run("MarkSingle_WithTimes_Success", func(t *testing.T) {
		req := CreateAttendanceRequest{
			StudentID: "s1",
			ClassID:   "c1",
			Date:      "2025-10-10",
			Status:    StatusLate,
			EntryTime: "08:30",
			ExitTime:  "13:00",
		}

		mockRepo.On("Create", mock.MatchedBy(func(a *Attendance) bool {
			return a.StudentID == "s1" && a.Status == StatusLate && a.EntryTime != nil && *a.EntryTime == "08:30" && a.ExitTime != nil && *a.ExitTime == "13:00"
		})).Return(nil).Once()

		err := service.MarkAttendance(ctx, teacherID, req)
		assert.NoError(t, err)
	})

	t.Run("MarkBulk_Success", func(t *testing.T) {
		req := BulkAttendanceRequest{
			ClassID: "c1",
			Date:    "2025-10-10",
			Statuses: []CreateAttendanceRequest{
				{StudentID: "s1", Status: StatusPresent},
				{StudentID: "s2", Status: StatusAbsent},
			},
		}

		mockRepo.On("BatchCreate", mock.MatchedBy(func(atts []*Attendance) bool {
			return len(atts) == 2
		})).Return(nil).Once()

		err := service.MarkBulk(ctx, teacherID, req)
		assert.NoError(t, err)
	})
}

func TestGetClassAttendance(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, nil, nil)
	ctx := context.Background()

	t.Run("ReturnsSummary", func(t *testing.T) {
		dateStr := "2025-10-10"
		date, _ := time.Parse("2006-01-02", dateStr)

		atts := []Attendance{
			{ID: "1", StudentID: "s1", Status: StatusPresent, Date: date},
			{ID: "2", StudentID: "s2", Status: StatusAbsent, Date: date},
			{ID: "3", StudentID: "s3", Status: StatusLate, Date: date},
		}

		mockRepo.On("FindByClassAndDate", "c1", date).Return(atts, nil).Once()

		resp, err := service.GetClassAttendance(ctx, "c1", dateStr)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Summary.Present)
		assert.Equal(t, 1, resp.Summary.Absent)
		assert.Equal(t, 1, resp.Summary.Late)
	})
}

func TestJustificationFlow(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, nil, nil)
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
		mockRepo.On("UpdateJustification", mock.MatchedBy(func(j *Justification) bool {
			return j.Status == JustificationApproved && j.ApprovedBy != nil
		})).Return(nil).Once()

		err := service.ProcessJustification(ctx, "t1", jID, true)
		assert.NoError(t, err)
	})
}
