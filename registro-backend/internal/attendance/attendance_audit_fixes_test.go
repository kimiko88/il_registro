package attendance

import (
	"context"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAttendance_UpdateAttendance_ValidationOrderBeforeFindByID(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	// Missing ID should fail validation before any repo FindByID call
	errID := svc.UpdateAttendance(context.Background(), "teacher-1", "school-1", "", UpdateAttendanceRequest{})
	assert.Error(t, errID)
	assert.Contains(t, errID.Error(), "id mancante")

	// Missing schoolID should fail validation before repo call
	errSchool := svc.UpdateAttendance(context.Background(), "teacher-1", "", "att-1", UpdateAttendanceRequest{})
	assert.Error(t, errSchool)
	assert.Contains(t, errSchool.Error(), "school_id mancante")

	// Missing teacherID should fail validation before repo call
	errTeacher := svc.UpdateAttendance(context.Background(), "", "school-1", "att-1", UpdateAttendanceRequest{})
	assert.Error(t, errTeacher)
	assert.Contains(t, errTeacher.Error(), "teacherID mancante")

	// Verify FindByID was never called
	mockRepo.AssertNotCalled(t, "FindByID", mock.Anything)
}

func TestAttendance_ProcessJustification_UsesActorRole(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	j := &Justification{
		ID:        "j-1",
		StudentID: "student-1",
		Status:    JustificationPending,
	}
	mockRepo.On("FindJustificationByID", "j-1").Return(j, nil)
	mockUserRepo.On("GetByID", mock.Anything, "student-1").Return(&users.User{ID: "student-1"}, nil)
	mockUserRepo.On("GetByID", mock.Anything, "admin-1").Return(&users.User{ID: "admin-1", Role: "admin"}, nil)
	mockRepo.On("ProcessJustificationTx", mock.Anything, j, "admin-1", true).Return(nil)

	// Passing "admin" role directly from verified token context should allow processing
	err := svc.ProcessJustification(context.Background(), "admin-1", "admin", "j-1", true)
	assert.NoError(t, err)
	assert.Equal(t, JustificationApproved, j.Status)
}

func TestAttendance_GetStudentAttendance_LogsAuditWarningWhenClassIDMissing(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	schID := "school-1"
	mockUserRepo.On("GetByID", mock.Anything, "teacher-1").Return(&users.User{ID: "teacher-1", Role: "teacher", SchoolID: &schID}, nil)
	// Student user has nil ClassID
	mockUserRepo.On("GetByID", mock.Anything, "student-1").Return(&users.User{ID: "student-1", SchoolID: &schID, ClassID: nil}, nil)

	_, err := svc.GetStudentAttendance(context.Background(), "teacher-1", "teacher", "school-1", "student-1", time.Now(), time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non è assegnato ad alcuna classe")
}

func TestAttendance_MarkBulk_DeduplicatesStudentIDs(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	schID := "school-1"
	mockUserRepo.On("GetByID", mock.Anything, "teacher-1").Return(&users.User{ID: "teacher-1", Role: "teacher", SchoolID: &schID}, nil)
	mockRepo.On("IsTeacherAssignedToClass", mock.Anything, "teacher-1", "class-1").Return(true, nil)

	// AreStudentsInClass should receive deduplicated student IDs slice ["s-1", "s-2"]
	mockRepo.On("AreStudentsInClass", mock.Anything, []string{"s-1", "s-2"}, "class-1").Return(map[string]bool{"s-1": true, "s-2": true}, nil)
	mockRepo.On("BatchCreate", mock.Anything).Return(nil)

	req := BulkAttendanceRequest{
		ClassID: "class-1",
		Hour:    1,
		Date:    "2026-03-01",
		Statuses: []StudentStatusRequest{
			{StudentID: "s-1", Status: StatusPresent},
			{StudentID: "s-1", Status: StatusPresent}, // Duplicate ID
			{StudentID: "s-2", Status: StatusAbsent},
		},
	}

	err := svc.MarkBulk(context.Background(), "teacher-1", "school-1", req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAttendance_GetChildAttendanceTrends_CountsEarlyExits(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	schID := "school-1"
	mockUserRepo.On("GetByID", mock.Anything, "parent-1").Return(&users.User{ID: "parent-1", Role: "parent"}, nil)
	mockUserRepo.On("IsGuardian", mock.Anything, "parent-1", "student-1").Return(true, nil)
	mockUserRepo.On("GetByID", mock.Anything, "student-1").Return(&users.User{ID: "student-1", SchoolID: &schID}, nil)

	d1, _ := time.Parse("2006-01-02", "2026-03-01")
	d2, _ := time.Parse("2006-01-02", "2026-03-02")
	atts := []Attendance{
		{StudentID: "student-1", Date: d1, Status: StatusEarlyExit},
		{StudentID: "student-1", Date: d2, Status: StatusPresent},
	}
	mockRepo.On("FindByStudent", "student-1", mock.Anything, mock.Anything).Return(atts, nil)

	trends, err := svc.GetChildAttendanceTrends(context.Background(), "parent-1", "student-1")
	assert.NoError(t, err)
	assert.NotNil(t, trends)
	assert.NotEmpty(t, trends.Trends)
	assert.Equal(t, 1, trends.Trends[0].EarlyExits)
	assert.Equal(t, 100.0, trends.Trends[0].PresenceRate)
}
