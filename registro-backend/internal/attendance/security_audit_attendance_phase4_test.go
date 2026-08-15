package attendance

import (
	"context"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepoPhase4 struct {
	mock.Mock
	users.Repository
}

func (m *mockUserRepoPhase4) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *mockUserRepoPhase4) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

type mockAttRepoPhase4 struct {
	mock.Mock
	Repository
}

func (m *mockAttRepoPhase4) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	args := m.Called(ctx, teacherID, classID)
	return args.Bool(0), args.Error(1)
}

func (m *mockAttRepoPhase4) IsClassInSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	args := m.Called(ctx, classID, schoolID)
	return args.Bool(0), args.Error(1)
}

func (m *mockAttRepoPhase4) DeleteByClassDateHour(schoolID, classID string, date time.Time, hour int) error {
	args := m.Called(schoolID, classID, date, hour)
	return args.Error(0)
}

func (m *mockAttRepoPhase4) GetStats(studentID string) (*SummaryResponse, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SummaryResponse), args.Error(1)
}

func TestAttendance_ProcessJustification_NilUserRepo(t *testing.T) {
	assert.Panics(t, func() {
		NewService(new(mockAttRepoPhase4), nil, nil, nil)
	})
}

func TestAttendance_DeleteClassAttendanceHour_CrossSchoolTeacherCheck(t *testing.T) {
	aRepo := new(mockAttRepoPhase4)
	uRepo := new(mockUserRepoPhase4)
	svc := NewService(aRepo, uRepo, nil, nil)
	ctx := context.Background()

	teacherID := "teacher-1"
	schoolID := "school-A"
	classID := "class-B" // Class belonging to school B!

	aRepo.On("IsClassInSchool", ctx, classID, schoolID).Return(false, nil).Once()

	err := svc.DeleteClassAttendanceHour(ctx, teacherID, "teacher", schoolID, classID, "2026-08-15", 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non appartiene alla scuola")
	aRepo.AssertExpectations(t)
}

func (m *mockAttRepoPhase4) AreStudentsInClass(ctx context.Context, studentIDs []string, classID string) (map[string]bool, error) {
	args := m.Called(ctx, studentIDs, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]bool), args.Error(1)
}

func TestAttendance_GetStudentSummary_TeachingDaysWeekdayFallback(t *testing.T) {
	aRepo := new(mockAttRepoPhase4)
	uRepo := new(mockUserRepoPhase4)
	svc := NewService(aRepo, uRepo, nil, nil)
	ctx := context.Background()

	aRepo.On("GetStats", "student-1").Return(&SummaryResponse{
		TotalAbsences: 5,
	}, nil).Once()

	res, err := svc.GetStudentSummary(ctx, "superadmin-1", "superadmin", "school-1", "student-1")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Greater(t, res.AbsenceRate, 0.0)
	aRepo.AssertExpectations(t)
}
