package attendance

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSecurity_MarkAttendance_EnrolmentCheck(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	ctx := context.Background()
	teacherID := "teacher-1"
	schoolID := "school-1"
	classID := "class-1"
	studentID := "student-outsider"
	today := time.Now().Format("2006-01-02")

	mockRepo.On("IsTeacherAssignedToClass", ctx, teacherID, classID).Return(true, nil).Once()
	// Student is NOT enrolled in class-1
	mockRepo.On("IsStudentInClass", ctx, studentID, classID).Return(false, nil).Once()

	err := svc.MarkAttendance(ctx, teacherID, schoolID, CreateAttendanceRequest{
		StudentID: studentID,
		ClassID:   classID,
		Date:      today,
		Hour:      1,
		Status:    StatusPresent,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lo studente non appartiene alla classe")
	mockRepo.AssertExpectations(t)
}

func TestSecurity_MarkBulk_SubstitutionSpoofingRejected(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	ctx := context.Background()
	teacherID := "teacher-1"
	schoolID := "school-1"
	classID := "class-1"
	today := time.Now().Format("2006-01-02")

	// Teacher is NOT assigned to class-1
	mockRepo.On("IsTeacherAssignedToClass", ctx, teacherID, classID).Return(false, nil).Once()
	mockRepo.On("IsClassInSchool", ctx, classID, schoolID).Return(true, nil).Once()
	// No substitute record in DB
	mockRepo.On("IsTeacherSubstitute", ctx, teacherID, classID, mock.Anything, 1).Return(false, nil).Once()

	// Client sends IsSubstitution: true attempting to bypass authorization
	err := svc.MarkBulk(ctx, teacherID, schoolID, BulkAttendanceRequest{
		ClassID:        classID,
		Date:           today,
		Hour:           1,
		IsSubstitution: true,
		Statuses: []StudentStatusRequest{
			{StudentID: "stu-1", Status: StatusPresent},
		},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docente non registrato come supplente")
	mockRepo.AssertExpectations(t)
}

func TestSecurity_JustifyChildAbsence_AttendanceOwnership(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	ctx := context.Background()
	parentID := "parent-1"
	studentID := "student-1"
	otherStudentAttID := "att-other-student"

	mockUserRepo.On("IsGuardian", ctx, parentID, studentID).Return(true, nil).Once()
	mockRepo.On("JustifyAbsenceByParent", otherStudentAttID, studentID, "Malattia", "").Return(assert.AnError).Once()

	err := svc.JustifyChildAbsence(ctx, parentID, studentID, otherStudentAttID, JustifyAbsenceRequest{
		Reason: "Malattia",
	})

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
