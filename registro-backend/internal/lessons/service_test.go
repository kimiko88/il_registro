package lessons

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	lessons   []Lesson
	homeworks []Homework
	errLesson error
	errHW     error
}

func (m *mockRepository) CreateLesson(l *Lesson) error {
	if m.errLesson != nil {
		return m.errLesson
	}
	l.ID = "lesson-id"
	l.TeacherName = "Professor Snape"
	m.lessons = append(m.lessons, *l)
	return nil
}

func (m *mockRepository) GetLessonsByClass(classID string, date string) ([]Lesson, error) {
	if m.errLesson != nil {
		return nil, m.errLesson
	}
	return m.lessons, nil
}

func (m *mockRepository) GetLessonsByClassAndSubject(classID, subjectID string, date string) ([]Lesson, error) {
	if m.errLesson != nil {
		return nil, m.errLesson
	}
	var res []Lesson
	for _, l := range m.lessons {
		if l.ClassID == classID && l.SubjectID == subjectID {
			res = append(res, l)
		}
	}
	return res, nil
}

func (m *mockRepository) CreateHomework(h *Homework) error {
	if m.errHW != nil {
		return m.errHW
	}
	h.ID = "hw-id"
	h.TeacherName = "Professor Snape"
	m.homeworks = append(m.homeworks, *h)
	return nil
}

func (m *mockRepository) GetHomeworkByClass(classID string) ([]Homework, error) {
	if m.errHW != nil {
		return nil, m.errHW
	}
	return m.homeworks, nil
}

func TestCreateLesson(t *testing.T) {
	repo := &mockRepository{}
	s := NewService(repo)

	// Test invalid date
	req := CreateLessonRequest{
		ClassID:   "class-1",
		SubjectID: "potion-1",
		Date:      "invalid-date",
		Hour:      1,
		Duration:  2,
		Topic:     "Intro",
	}
	resp, err := s.CreateLesson("teacher-1", req)
	assert.Nil(t, resp)
	assert.Error(t, err)

	// Test successful creation
	req.Date = "2026-07-14"
	resp, err = s.CreateLesson("teacher-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "lesson-id", resp.ID)
	assert.Equal(t, "Professor Snape", resp.TeacherName)

	// Test repo error
	repo.errLesson = errors.New("db error")
	resp, err = s.CreateLesson("teacher-1", req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "db error")
}

func TestCreateHomework(t *testing.T) {
	repo := &mockRepository{}
	s := NewService(repo)

	// Test invalid date
	req := CreateHomeworkRequest{
		ClassID:     "class-1",
		SubjectID:   "potion-1",
		DueDate:     "invalid-date",
		Description: "Read page 394",
	}
	resp, err := s.CreateHomework("teacher-1", req)
	assert.Nil(t, resp)
	assert.Error(t, err)

	// Test successful creation
	req.DueDate = "2026-07-15"
	resp, err = s.CreateHomework("teacher-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "hw-id", resp.ID)
	assert.Equal(t, "Professor Snape", resp.TeacherName)
}

func TestGetLessonsAndHomeworks(t *testing.T) {
	repo := &mockRepository{
		lessons: []Lesson{
			{ID: "lesson-1", ClassID: "class-1", SubjectID: "potion-1", Date: time.Now()},
		},
		homeworks: []Homework{
			{ID: "hw-1", ClassID: "class-1", SubjectID: "potion-1", DueDate: time.Now()},
		},
	}
	s := NewService(repo)

	// Get lessons
	lessons, err := s.GetLessons("class-1", "", "")
	assert.NoError(t, err)
	assert.Len(t, lessons, 1)

	// Get lessons with subject
	lessonsSub, err := s.GetLessons("class-1", "potion-1", "")
	assert.NoError(t, err)
	assert.Len(t, lessonsSub, 1)

	// Get homeworks
	hws, err := s.GetHomeworks("class-1")
	assert.NoError(t, err)
	assert.Len(t, hws, 1)
}
