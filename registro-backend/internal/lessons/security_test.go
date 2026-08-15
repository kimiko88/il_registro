package lessons

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLessonsRepo struct {
	mock.Mock
	Repository
}

func (m *mockLessonsRepo) IsTeacherAssignedToClass(teacherID, classID string) (bool, error) {
	args := m.Called(teacherID, classID)
	return args.Bool(0), args.Error(1)
}

func (m *mockLessonsRepo) HasApprovedSubstitution(teacherID, classID, date string, hour int) (bool, error) {
	args := m.Called(teacherID, classID, date, hour)
	return args.Bool(0), args.Error(1)
}

func (m *mockLessonsRepo) GetLessonByID(id string) (*Lesson, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Lesson), args.Error(1)
}

func (m *mockLessonsRepo) UpdateLesson(id string, req UpdateLessonRequest) (*Lesson, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Lesson), args.Error(1)
}

func (m *mockLessonsRepo) DeleteLesson(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockLessonsRepo) GetLessonsByClass(classID string, date string) ([]Lesson, error) {
	args := m.Called(classID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Lesson), args.Error(1)
}

func TestSecurity_CreateLesson_SubstitutionErrorHandling(t *testing.T) {
	repo := new(mockLessonsRepo)
	svc := NewService(repo)

	teacherID := "teacher-1"
	classID := "class-1"
	today := time.Now().Format("2006-01-02")

	repo.On("IsTeacherAssignedToClass", teacherID, classID).Return(false, nil).Once()
	// DB error when checking substitution
	repo.On("HasApprovedSubstitution", teacherID, classID, today, 1).Return(false, assert.AnError).Once()

	_, err := svc.CreateLesson(teacherID, CreateLessonRequest{
		ClassID:        classID,
		SubjectID:      "sub-1",
		Date:           today,
		Hour:           1,
		Duration:       1,
		IsSubstitution: true,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
	repo.AssertExpectations(t)
}

func TestSecurity_UpdateLesson_PrincipalAndVicePrincipalAllowed(t *testing.T) {
	repo := new(mockLessonsRepo)
	svc := NewService(repo)

	lessonID := "lesson-1"
	existing := &Lesson{
		ID:        lessonID,
		TeacherID: "other-teacher",
		ClassID:   "class-1",
		Hour:      1,
		Duration:  1,
		Date:      time.Now(),
	}

	repo.On("GetLessonByID", lessonID).Return(existing, nil).Times(2)
	repo.On("GetLessonsByClass", "class-1", existing.Date.Format("2006-01-02")).Return([]Lesson{}, nil)
	repo.On("UpdateLesson", lessonID, mock.Anything).Return(existing, nil)

	// Principal role updating another teacher's lesson
	resPrincipal, errP := svc.UpdateLesson("principal-1", "principal", lessonID, UpdateLessonRequest{Topic: "Update by Principal"})
	assert.NoError(t, errP)
	assert.NotNil(t, resPrincipal)

	// Vice-Principal role updating another teacher's lesson
	resVice, errV := svc.UpdateLesson("vice-1", "vice_principal", lessonID, UpdateLessonRequest{Topic: "Update by Vice Principal"})
	assert.NoError(t, errV)
	assert.NotNil(t, resVice)

	repo.AssertExpectations(t)
}
