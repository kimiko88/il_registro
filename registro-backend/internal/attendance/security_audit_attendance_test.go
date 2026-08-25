package attendance

import (
	"context"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepoAudit struct {
	mock.Mock
	users.Repository
}

func (m *mockUserRepoAudit) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepoAudit) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func TestAttendance_MarkBulk_SubstituteCrossSchoolRejection(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(mockUserRepoAudit)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	ctx := context.Background()
	teacherID := "substitute-teacher"
	schoolID := "school-A"
	otherSchoolClassID := "class-school-B"
	today := time.Now().Format("2006-01-02")

	mockRepo.On("IsTeacherAssignedToClass", ctx, teacherID, otherSchoolClassID).Return(false, nil).Once()
	// Class belongs to School B, not School A
	mockRepo.On("IsClassInSchool", ctx, otherSchoolClassID, schoolID).Return(false, nil).Once()

	err := svc.MarkBulk(ctx, teacherID, schoolID, BulkAttendanceRequest{
		ClassID:        otherSchoolClassID,
		Date:           today,
		Hour:           1,
		IsSubstitution: true,
		Statuses: []StudentStatusRequest{
			{StudentID: "stu-1", Status: StatusPresent},
		},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "classe non appartiene alla scuola")
	mockRepo.AssertExpectations(t)
}

func TestAttendance_GetStudentSummary_StudentSelfAccessOnly(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(mockUserRepoAudit)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	ctx := context.Background()

	// Student "student-1" tries to view summary of "student-2"
	_, err := svc.GetStudentSummary(ctx, "student-1", "student", "school-A", "student-2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")

	// Student "student-1" views own summary -> allowed
	mockRepo.On("GetStats", "student-1").Return(&SummaryResponse{}, nil).Once()

	res, err := svc.GetStudentSummary(ctx, "student-1", "student", "school-A", "student-1")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}
