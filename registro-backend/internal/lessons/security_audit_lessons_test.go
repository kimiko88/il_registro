package lessons

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLessonsRepoAudit struct {
	mock.Mock
	Repository
}

func (m *mockLessonsRepoAudit) GetLessonsByTeacher(teacherID string, fromDate, toDate string) ([]Lesson, error) {
	args := m.Called(teacherID, fromDate, toDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Lesson), args.Error(1)
}

func (m *mockLessonsRepoAudit) IsTeacherAssignedToClass(teacherID, classID string) (bool, error) {
	args := m.Called(teacherID, classID)
	return args.Bool(0), args.Error(1)
}

func (m *mockLessonsRepoAudit) GetLessonsByClass(classID string, date string) ([]Lesson, error) {
	args := m.Called(classID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Lesson), args.Error(1)
}

func TestLessons_GetTeacherDiary_InvalidDateParsingRejection(t *testing.T) {
	repo := new(mockLessonsRepoAudit)
	svc := NewService(repo)

	// Invalid date string MUST return an error instead of bypassing filters
	_, err := svc.GetTeacherDiary("teacher-1", "invalid-date", "2026-08-15")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid date format")
}

func TestLessons_CreateLesson_SymmetricSubstitutionOverlap(t *testing.T) {
	repo := new(mockLessonsRepoAudit)
	svc := NewService(repo)

	teacherID := "teacher-1"
	classID := "class-1"
	date := "2026-08-15"

	repo.On("IsTeacherAssignedToClass", teacherID, classID).Return(true, nil).Once()
	// Existing lesson is a substitution lesson
	repo.On("GetLessonsByClass", classID, date).Return([]Lesson{
		{
			ID:             "lesson-sub",
			ClassID:        classID,
			Hour:           1,
			Duration:       1,
			IsSubstitution: true,
		},
	}, nil).Once()

	// Regular teacher attempts to insert a normal lesson (IsSubstitution = false)
	_, err := svc.CreateLesson(teacherID, CreateLessonRequest{
		ClassID:        classID,
		SubjectID:      "sub-1",
		Date:           date,
		Hour:           1,
		Duration:       1,
		IsSubstitution: false,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "impossibile inserire più lezioni")
	repo.AssertExpectations(t)
}
