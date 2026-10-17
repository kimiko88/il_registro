package homework_submissions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHomeworkSubmissionsRepo struct {
	mock.Mock
}

func (m *mockHomeworkSubmissionsRepo) UpsertSubmission(s *HomeworkSubmission) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *mockHomeworkSubmissionsRepo) GetSubmission(id string) (*HomeworkSubmission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*HomeworkSubmission), args.Error(1)
}

func (m *mockHomeworkSubmissionsRepo) GetSubmissionByHomeworkAndStudent(homeworkID, studentID string) (*HomeworkSubmission, error) {
	args := m.Called(homeworkID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*HomeworkSubmission), args.Error(1)
}

func (m *mockHomeworkSubmissionsRepo) ListSubmissionsByHomework(homeworkID string) ([]HomeworkSubmission, error) {
	args := m.Called(homeworkID)
	return args.Get(0).([]HomeworkSubmission), args.Error(1)
}

func (m *mockHomeworkSubmissionsRepo) ListSubmissionsByStudent(studentID string) ([]HomeworkSubmission, error) {
	args := m.Called(studentID)
	return args.Get(0).([]HomeworkSubmission), args.Error(1)
}

func (m *mockHomeworkSubmissionsRepo) SaveGrading(id string, req GradeSubmissionRequest, teacherID string) error {
	args := m.Called(id, req, teacherID)
	return args.Error(0)
}

func (m *mockHomeworkSubmissionsRepo) SetClassroomSync(id, classroomSubmissionID string) error {
	args := m.Called(id, classroomSubmissionID)
	return args.Error(0)
}

func TestHomeworkSubmissionsService_SubmitAndGrade(t *testing.T) {
	repo := new(mockHomeworkSubmissionsRepo)
	svc := NewService(repo)

	repo.On("UpsertSubmission", mock.AnythingOfType("*homework_submissions.HomeworkSubmission")).Return(nil).Once()

	sub := &HomeworkSubmission{
		HomeworkID: "hw-1",
		StudentID:  "stud-1",
		FileURL:    "https://storage.scuola.edu.it/homeworks/compito1.pdf",
		FileType:   "pdf",
	}

	err := svc.SubmitHomework(sub)
	assert.NoError(t, err)

	// Missing fileURL
	err = svc.SubmitHomework(&HomeworkSubmission{
		HomeworkID: "hw-1",
		StudentID:  "stud-1",
	})
	assert.Error(t, err)

	// Grade submission
	grade := 8.5
	repo.On("SaveGrading", "sub-1", mock.AnythingOfType("homework_submissions.GradeSubmissionRequest"), "teach-1").Return(nil).Once()

	err = svc.GradeSubmission("sub-1", GradeSubmissionRequest{
		TeacherGrade:    &grade,
		TeacherFeedback: "Ottimo svolgimento del problema 3",
		AnnotationsJSON: `[{"type":"pen","x":100,"y":200}]`,
	}, "teach-1")
	assert.NoError(t, err)

	// Invalid grade > 10
	invalidGrade := 11.0
	err = svc.GradeSubmission("sub-1", GradeSubmissionRequest{
		TeacherGrade: &invalidGrade,
	}, "teach-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "compreso tra 1 e 10")

	repo.AssertExpectations(t)
}
