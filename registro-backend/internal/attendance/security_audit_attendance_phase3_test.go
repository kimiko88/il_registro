package attendance

import (
	"context"
	"testing"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepoPhase3 struct {
	mock.Mock
	users.Repository
}

func (m *mockUserRepoPhase3) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *mockUserRepoPhase3) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

type mockAttRepoPhase3 struct {
	mock.Mock
	Repository
}

func (m *mockAttRepoPhase3) FindByID(id string) (*Attendance, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Attendance), args.Error(1)
}

func (m *mockAttRepoPhase3) FindUnjustifiedByStudent(studentID string) ([]Attendance, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Attendance), args.Error(1)
}

func (m *mockAttRepoPhase3) JustifyAbsenceByParent(id, studentID, reason, notes string) error {
	args := m.Called(id, studentID, reason, notes)
	return args.Error(0)
}

func TestAttendance_JustifyChildAbsence_AttendanceStudentOwnershipValidation(t *testing.T) {
	uRepo := new(mockUserRepoPhase3)
	aRepo := new(mockAttRepoPhase3)

	svc := NewService(aRepo, uRepo, nil, nil)
	ctx := context.Background()

	parentID := "parent-1"
	studentID := "student-1"
	attID := "att-999"

	// Parent is guardian of student-1
	uRepo.On("GetByID", ctx, parentID).Return(&users.User{ID: parentID, Role: "parent"}, nil).Once()
	uRepo.On("IsGuardian", ctx, parentID, studentID).Return(true, nil).Once()

	// Attendance record belongs to student-2 (DIFFERENT student!)
	aRepo.On("FindByID", attID).Return(&Attendance{
		ID:        attID,
		StudentID: "student-2",
	}, nil).Once()

	err := svc.JustifyChildAbsence(ctx, parentID, studentID, attID, JustifyAbsenceRequest{
		Reason: "Malattia",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non appartiene allo studente")
	uRepo.AssertExpectations(t)
	aRepo.AssertExpectations(t)
}

func TestAttendance_AdminAccessToParentChildEndpoints(t *testing.T) {
	uRepo := new(mockUserRepoPhase3)
	aRepo := new(mockAttRepoPhase3)

	svc := NewService(aRepo, uRepo, nil, nil)
	ctx := context.Background()

	adminID := "admin-1"
	studentID := "student-1"
	schoolID := "school-1"

	// Caller is an Admin (role: admin) in school-1 and target student is in school-1
	uRepo.On("GetByID", ctx, adminID).Return(&users.User{ID: adminID, Role: "admin", SchoolID: &schoolID}, nil).Once()
	uRepo.On("GetByID", ctx, studentID).Return(&users.User{ID: studentID, Role: "student", SchoolID: &schoolID}, nil).Once()

	aRepo.On("FindUnjustifiedByStudent", studentID).Return([]Attendance{}, nil).Once()

	res, err := svc.GetChildUnjustified(ctx, adminID, studentID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	uRepo.AssertExpectations(t)
	aRepo.AssertExpectations(t)
}

func TestAttendance_CrossSchoolAdminAccess_Forbidden(t *testing.T) {
	uRepo := new(mockUserRepoPhase3)
	aRepo := new(mockAttRepoPhase3)

	svc := NewService(aRepo, uRepo, nil, nil)
	ctx := context.Background()

	adminID := "admin-1"
	studentID := "student-1"
	schoolA := "school-A"
	schoolB := "school-B"

	// Admin is in school-A, student is in school-B -> should be forbidden
	uRepo.On("GetByID", ctx, adminID).Return(&users.User{ID: adminID, Role: "admin", SchoolID: &schoolA}, nil).Once()
	uRepo.On("GetByID", ctx, studentID).Return(&users.User{ID: studentID, Role: "student", SchoolID: &schoolB}, nil).Once()

	_, err := svc.GetChildUnjustified(ctx, adminID, studentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
	uRepo.AssertExpectations(t)
}
