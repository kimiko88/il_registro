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

func (m *mockRepository) GetLessonByID(id string) (*Lesson, error) {
	if m.errLesson != nil {
		return nil, m.errLesson
	}
	if len(m.lessons) > 0 {
		return &m.lessons[0], nil
	}
	return &Lesson{ID: id, TeacherID: "teacher-1", Topic: "Test"}, nil
}

func (m *mockRepository) UpdateLesson(id string, req UpdateLessonRequest) (*Lesson, error) {
	if m.errLesson != nil {
		return nil, m.errLesson
	}
	return &Lesson{ID: id, TeacherID: "teacher-1", Topic: req.Topic}, nil
}

func (m *mockRepository) DeleteLesson(id string) error {
	return m.errLesson
}

func (m *mockRepository) IsTeacherAssignedToClass(teacherID, classID string) (bool, error) {
	return true, nil
}

func (m *mockRepository) HasApprovedSubstitution(teacherID, classID, date string, hour int) (bool, error) {
	return true, nil
}

func (m *mockRepository) GetHomeworkByID(id string) (*Homework, error) {
	if m.errHW != nil {
		return nil, m.errHW
	}
	if len(m.homeworks) > 0 {
		return &m.homeworks[0], nil
	}
	return &Homework{ID: id, TeacherID: "teacher-1"}, nil
}

func (m *mockRepository) UpdateHomework(id string, req UpdateHomeworkRequest) (*Homework, error) {
	if m.errHW != nil {
		return nil, m.errHW
	}
	return &Homework{ID: id, TeacherID: "teacher-1", Description: req.Description}, nil
}

func (m *mockRepository) DeleteHomework(id string) error {
	return m.errHW
}

func (m *mockRepository) GetLessonsByTeacher(teacherID string, fromDate, toDate string) ([]Lesson, error) {
	if m.errLesson != nil {
		return nil, m.errLesson
	}
	return m.lessons, nil
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

func (m *mockRepository) GetLessonsByGroup(groupID string, date string) ([]Lesson, error) {
	if m.errLesson != nil {
		return nil, m.errLesson
	}
	return m.lessons, nil
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
	req.DueDate = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
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

func TestLessonOwnershipChecks(t *testing.T) {
	repo := &mockRepository{
		lessons: []Lesson{
			{ID: "lesson-1", TeacherID: "teacher-owner"},
		},
		homeworks: []Homework{
			{ID: "hw-1", TeacherID: "teacher-owner"},
		},
	}
	s := NewService(repo)

	// Update lesson - owner succeeds
	res, err := s.UpdateLesson("teacher-owner", "teacher", "lesson-1", UpdateLessonRequest{Topic: "New Topic"})
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// Update lesson - non-owner fails
	_, err = s.UpdateLesson("other-teacher", "teacher", "lesson-1", UpdateLessonRequest{Topic: "New Topic"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	// Update lesson - admin succeeds even if non-owner
	_, err = s.UpdateLesson("admin-user", "admin", "lesson-1", UpdateLessonRequest{Topic: "New Topic"})
	assert.NoError(t, err)

	// Delete lesson - non-owner fails
	err = s.DeleteLesson("other-teacher", "teacher", "lesson-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	// Delete lesson - owner succeeds
	err = s.DeleteLesson("teacher-owner", "teacher", "lesson-1")
	assert.NoError(t, err)
}
