package lessons

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLessonsRepository struct {
	mock.Mock
}

func (m *MockLessonsRepository) CreateLesson(lesson *Lesson) error {
	return m.Called(lesson).Error(0)
}
func (m *MockLessonsRepository) GetLessonByID(id string) (*Lesson, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Lesson), args.Error(1)
}
func (m *MockLessonsRepository) UpdateLesson(id string, req UpdateLessonRequest) (*Lesson, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Lesson), args.Error(1)
}
func (m *MockLessonsRepository) DeleteLesson(id string) error {
	return m.Called(id).Error(0)
}
func (m *MockLessonsRepository) GetLessonsByClass(classID string, date string) ([]Lesson, error) {
	args := m.Called(classID, date)
	return args.Get(0).([]Lesson), args.Error(1)
}
func (m *MockLessonsRepository) GetLessonsByClassAndSubject(classID, subjectID string, date string) ([]Lesson, error) {
	args := m.Called(classID, subjectID, date)
	return args.Get(0).([]Lesson), args.Error(1)
}
func (m *MockLessonsRepository) GetLessonsByGroup(groupID string, date string) ([]Lesson, error) {
	args := m.Called(groupID, date)
	return args.Get(0).([]Lesson), args.Error(1)
}
func (m *MockLessonsRepository) GetLessonsByTeacher(teacherID string, fromDate, toDate string) ([]Lesson, error) {
	args := m.Called(teacherID, fromDate, toDate)
	return args.Get(0).([]Lesson), args.Error(1)
}
func (m *MockLessonsRepository) CreateHomework(homework *Homework) error {
	return m.Called(homework).Error(0)
}
func (m *MockLessonsRepository) GetHomeworkByID(id string) (*Homework, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Homework), args.Error(1)
}
func (m *MockLessonsRepository) UpdateHomework(id string, req UpdateHomeworkRequest) (*Homework, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Homework), args.Error(1)
}
func (m *MockLessonsRepository) DeleteHomework(id string) error {
	return m.Called(id).Error(0)
}
func (m *MockLessonsRepository) GetHomeworkByClass(classID string, fromDate ...string) ([]Homework, error) {
	args := m.Called(classID, fromDate)
	return args.Get(0).([]Homework), args.Error(1)
}
func (m *MockLessonsRepository) IsTeacherAssignedToClass(teacherID, classID string) (bool, error) {
	args := m.Called(teacherID, classID)
	return args.Bool(0), args.Error(1)
}
func (m *MockLessonsRepository) HasApprovedSubstitution(teacherID, classID, date string, hour int) (bool, error) {
	args := m.Called(teacherID, classID, date, hour)
	return args.Bool(0), args.Error(1)
}

func TestLessons_UpdateLesson_SubstitutionCollision(t *testing.T) {
	mockRepo := new(MockLessonsRepository)
	svc := NewService(mockRepo)

	existingLesson := &Lesson{
		ID:             "les-1",
		ClassID:        "class-1",
		TeacherID:      "t-1",
		Date:           time.Now(),
		Hour:           1,
		Duration:       1,
		IsSubstitution: true,
	}

	// Another existing lesson in the same hour is NOT a substitution
	classLessons := []Lesson{
		{
			ID:             "les-2",
			ClassID:        "class-1",
			TeacherID:      "t-2",
			Date:           time.Now(),
			Hour:           1,
			Duration:       1,
			IsSubstitution: false,
		},
	}

	mockRepo.On("GetLessonByID", "les-1").Return(existingLesson, nil).Once()
	mockRepo.On("GetLessonsByClass", "class-1", existingLesson.Date.Format("2006-01-02")).Return(classLessons, nil).Once()

	// Update lesson les-1 (which is a substitution) against les-2 (which is NOT a substitution)
	// Must result in collision error because both are not substitution/co-teaching
	_, err := svc.UpdateLesson("t-1", "teacher", "les-1", UpdateLessonRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "impossibile registrare più lezioni")
	mockRepo.AssertExpectations(t)
}

func TestLessons_GetHomeworks_Default30DaysFilter(t *testing.T) {
	mockRepo := new(MockLessonsRepository)
	svc := NewService(mockRepo)

	mockRepo.On("GetHomeworkByClass", "class-1", mock.Anything).Return([]Homework{}, nil).Once()

	res, err := svc.GetHomeworks("class-1")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}
