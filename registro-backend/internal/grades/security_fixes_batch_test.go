package grades

import (
	"bytes"
	"context"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. Test BatchCreateGrades dedup key with distinct descriptions/testIDs
func TestBatchCreateGrades_DeduplicationKey(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := NewService(repo, userRepo, nil, nil)

	now := time.Now()
	testID1 := "test-1"
	testID2 := "test-2"

	gradesList := []*Grade{
		// Two grades on same date with same value but DIFFERENT testIDs
		{
			StudentID:  "student-1",
			SubjectID:  "subject-1",
			Date:       now,
			GradeValue: 8.0,
			GradeType:  "scritto",
			TestID:     &testID1,
		},
		{
			StudentID:  "student-1",
			SubjectID:  "subject-1",
			Date:       now,
			GradeValue: 8.0,
			GradeType:  "scritto",
			TestID:     &testID2,
		},
		// A third grade on same date with same value but DIFFERENT description
		{
			StudentID:   "student-1",
			SubjectID:   "subject-1",
			Date:        now,
			GradeValue:  8.0,
			GradeType:   "orale",
			Description: "Interrogazione capitolo 1",
		},
		// A duplicate of the third grade (should be deduplicated)
		{
			StudentID:   "student-1",
			SubjectID:   "subject-1",
			Date:        now,
			GradeValue:  8.0,
			GradeType:   "orale",
			Description: "Interrogazione capitolo 1",
		},
	}

	repo.On("BatchCreate", mock.Anything, mock.MatchedBy(func(g []*Grade) bool {
		return len(g) == 3
	})).Return(nil).Once()

	err := svc.BatchCreateGrades(context.Background(), "teacher-1", "teacher", "school-1", gradesList)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// 2. Test Export forces teacherID scoping
func TestExport_TeacherScopingSecurity(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := NewService(repo, userRepo, nil, nil)

	// Teacher provides filter trying to query another teacher's grades
	filter := GradeFilter{
		TeacherID: "other-teacher-victim",
	}

	// Service should overwrite filter.TeacherID with authenticated teacherID ("teacher-authenticated")
	repo.On("FindWithFilter", mock.Anything, mock.MatchedBy(func(f GradeFilter) bool {
		return f.TeacherID == "teacher-authenticated"
	})).Return([]Grade{}, nil).Once()

	_, _, err := svc.Export(context.Background(), "teacher-authenticated", "school-1", filter, "json")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// 3. Test GetMyTrend with unavailable class average
func TestGetMyTrend_UnavailableClassAveragePosition(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := NewService(repo, userRepo, nil, nil)

	now := time.Now()
	studentID := "student-1"
	subjectID := "subject-1"

	repo.On("FindByStudent", mock.Anything, studentID).Return([]Grade{
		{
			StudentID:   studentID,
			SubjectID:   subjectID,
			GradeValue:  7.5,
			Date:        now,
			IsPublished: true,
			Semester:    1,
		},
	}, nil).Once()

	resp, err := svc.GetMyTrend(context.Background(), studentID, "student", studentID, subjectID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Trends, 1)
	// When classAverage is unavailable (-1), Position must be "not_available"
	assert.Equal(t, "not_available", resp.Trends[0].Position)
}

// 4. Test BulkImport teacher assignment validation
func TestBulkImport_TeacherAssignmentEnforcement(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := NewService(repo, userRepo, nil, nil)

	schoolID := "school-1"
	teacherID := "teacher-1"
	userRepo.On("GetByID", mock.Anything, teacherID).Return(&users.User{
		ID:       teacherID,
		Role:     "teacher",
		SchoolID: &schoolID,
	}, nil).Once()

	csvContent := "student_id,subject_id,grade_value,grade_type,date\nstudent-1,subject-1,8.0,scritto,2026-03-01\n"
	r := bytes.NewReader([]byte(csvContent))

	repo.On("BatchCreate", mock.Anything, mock.Anything).Return(nil).Once()

	res, err := svc.BulkImport(context.Background(), teacherID, schoolID, r, 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Imported)
}

