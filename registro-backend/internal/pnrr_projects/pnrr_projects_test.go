package pnrr_projects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPNRRRepo struct {
	mock.Mock
}

func (m *mockPNRRRepo) CreateProject(p *PNRRProject) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *mockPNRRRepo) ListProjectsBySchool(schoolID, academicYear string) ([]PNRRProject, error) {
	args := m.Called(schoolID, academicYear)
	return args.Get(0).([]PNRRProject), args.Error(1)
}

func (m *mockPNRRRepo) GetProject(id string) (*PNRRProject, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PNRRProject), args.Error(1)
}

func (m *mockPNRRRepo) EnrollStudent(ps *PNRRProjectStudent) error {
	args := m.Called(ps)
	return args.Error(0)
}

func (m *mockPNRRRepo) ListProjectStudents(projectID string) ([]PNRRProjectStudent, error) {
	args := m.Called(projectID)
	return args.Get(0).([]PNRRProjectStudent), args.Error(1)
}

func (m *mockPNRRRepo) RecordLesson(l *PNRRLessonAttendance, attendedStudentIDs []string) error {
	args := m.Called(l, attendedStudentIDs)
	return args.Error(0)
}

func (m *mockPNRRRepo) ListLessons(projectID string) ([]PNRRLessonAttendance, error) {
	args := m.Called(projectID)
	return args.Get(0).([]PNRRLessonAttendance), args.Error(1)
}

func TestPNRRProjectsService_ValidationAndWorkflow(t *testing.T) {
	repo := new(mockPNRRRepo)
	svc := NewService(repo)

	repo.On("CreateProject", mock.AnythingOfType("*pnrr_projects.PNRRProject")).Return(nil).Once()

	proj := &PNRRProject{
		SchoolID:     "sch-1",
		CUPCode:      "B12C23000450006",
		ProjectTitle: "Laboratorio Robotica e Coding STEM",
		NoticeType:   "DM_65_STEM",
		TotalHours:   30,
		AcademicYear: "2025/2026",
	}

	err := svc.CreateProject(proj)
	assert.NoError(t, err)
	assert.Equal(t, "attivo", proj.Status)

	// Invalid total hours <= 0
	err = svc.CreateProject(&PNRRProject{
		SchoolID:     "sch-1",
		CUPCode:      "CUP123",
		ProjectTitle: "Invalid",
		TotalHours:   0,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maggiore di zero")

	// Enroll student
	repo.On("EnrollStudent", mock.AnythingOfType("*pnrr_projects.PNRRProjectStudent")).Return(nil).Once()
	err = svc.EnrollStudent("proj-1", "stud-1")
	assert.NoError(t, err)

	// Record lesson
	repo.On("RecordLesson", mock.AnythingOfType("*pnrr_projects.PNRRLessonAttendance"), []string{"stud-1", "stud-2"}).Return(nil).Once()
	err = svc.RecordLesson(&PNRRLessonAttendance{
		ProjectID:  "proj-1",
		LessonDate: "2026-03-10",
		Hours:      3,
		Topic:      "Algoritmi e strutture dati",
	}, []string{"stud-1", "stud-2"})
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}
