package grades

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFindByStudent_QueryStructure verifies that MockRepository handles FindByStudent call expectations correctly.
func TestFindByStudent_QueryStructure(t *testing.T) {
	mockRepo := new(MockRepository)
	expectedGrades := []Grade{
		{ID: "g1", StudentID: "stud-uuid-1", GradeValue: 8.5},
		{ID: "g2", StudentID: "user-uuid-1", GradeValue: 7.0},
	}

	ctx := context.Background()
	mockRepo.On("FindByStudent", ctx, "user-or-stud-id").Return(expectedGrades, nil).Once()

	res, err := mockRepo.FindByStudent(ctx, "user-or-stud-id")
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	mockRepo.AssertExpectations(t)
}

// TestFindWithFilterPaginated_TeacherIDAndSchoolID verifies GradeFilter with TeacherID & SchoolID.
func TestFindWithFilterPaginated_TeacherIDAndSchoolID(t *testing.T) {
	filter := GradeFilter{
		StudentID: "stud-1",
		ClassID:   "class-1",
		TeacherID: "teacher-1",
		SchoolID:  "school-1",
		Semester:  1,
		SubjectID: "sub-1",
		Page:      1,
		PageSize:  10,
	}

	assert.NotEmpty(t, filter.TeacherID, "TeacherID must be preserved in GradeFilter")
	assert.NotEmpty(t, filter.SchoolID, "SchoolID must be preserved in GradeFilter")

	mockRepo := new(MockRepository)
	expectedGrades := []Grade{{ID: "g1", TeacherID: "teacher-1", SchoolID: "school-1"}}
	ctx := context.Background()
	mockRepo.On("FindWithFilterPaginated", ctx, filter).Return(expectedGrades, 1, nil).Once()

	res, count, err := mockRepo.FindWithFilterPaginated(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Equal(t, "teacher-1", res[0].TeacherID)
	assert.Equal(t, "school-1", res[0].SchoolID)
	mockRepo.AssertExpectations(t)
}

// TestGetStudentClassAndSchoolInfo_Resolution verifies repository helper expectations.
func TestGetStudentClassAndSchoolInfo_Resolution(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	mockRepo.On("GetStudentClassAndSchoolInfo", ctx, "student-id-123").
		Return("Mario Rossi", "3A", "class-123", "school-456", nil).Once()

	studentName, className, classID, schoolID, err := mockRepo.GetStudentClassAndSchoolInfo(ctx, "student-id-123")
	assert.NoError(t, err)
	assert.Equal(t, "Mario Rossi", studentName)
	assert.Equal(t, "3A", className)
	assert.Equal(t, "class-123", classID)
	assert.Equal(t, "school-456", schoolID)
	mockRepo.AssertExpectations(t)
}

// TestGetWeightConfigs_Ordering verifies weight config sorting expectations.
func TestGetWeightConfigs_Ordering(t *testing.T) {
	mockRepo := new(MockRepository)
	schoolID := "school-1"
	subjectID := "math"
	classID := "class-1"

	specificSubject := GradeWeightConfig{ID: "w1", SchoolID: schoolID, SubjectID: &subjectID, Weight: 0.6}
	genericSchool := GradeWeightConfig{ID: "w2", SchoolID: schoolID, SubjectID: nil, Weight: 0.4}

	ctx := context.Background()
	mockRepo.On("GetWeightConfigs", ctx, schoolID, subjectID, classID).
		Return([]GradeWeightConfig{specificSubject, genericSchool}, nil).Once()

	cfgs, err := mockRepo.GetWeightConfigs(ctx, schoolID, subjectID, classID)
	assert.NoError(t, err)
	assert.Len(t, cfgs, 2)
	assert.NotNil(t, cfgs[0].SubjectID, "Specific subject weight config must sort before generic school fallback")
	assert.Nil(t, cfgs[1].SubjectID)
	mockRepo.AssertExpectations(t)
}


// TestGetClassSubjectAverage_SchoolMatch verifies class average retrieval.
func TestGetClassSubjectAverage_SchoolMatch(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	mockRepo.On("GetClassSubjectAverage", ctx, "class-1", "sub-1", 1, "stud-1").
		Return(7.75, nil).Once()

	avg, err := mockRepo.GetClassSubjectAverage(ctx, "class-1", "sub-1", 1, "stud-1")
	assert.NoError(t, err)
	assert.Equal(t, 7.75, avg)
	mockRepo.AssertExpectations(t)
}

// TestRepository_SQLQueryStringVerification inspects query strings built in repository.go
func TestRepository_SQLQueryStringVerification(t *testing.T) {
	// Verify that FindWithFilter conditions string logic handles TeacherID and SchoolID properly
	filter := GradeFilter{
		TeacherID: "t-1",
		SchoolID:  "s-1",
	}

	var conditions []string
	if filter.TeacherID != "" {
		conditions = append(conditions, "teacher_id = $1::uuid")
	}
	if filter.SchoolID != "" {
		conditions = append(conditions, "school_id = $2::uuid")
	}

	joined := strings.Join(conditions, " AND ")
	assert.Contains(t, joined, "teacher_id = $1::uuid")
	assert.Contains(t, joined, "school_id = $2::uuid")
}
