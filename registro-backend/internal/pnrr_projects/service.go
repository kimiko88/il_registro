package pnrr_projects

import (
	"errors"
)

type Service interface {
	CreateProject(p *PNRRProject) error
	ListProjects(schoolID, academicYear string) ([]PNRRProject, error)
	GetProject(id string) (*PNRRProject, error)
	EnrollStudent(projectID, studentID string) error
	ListStudents(projectID string) ([]PNRRProjectStudent, error)
	RecordLesson(l *PNRRLessonAttendance, attendedStudentIDs []string) error
	ListLessons(projectID string) ([]PNRRLessonAttendance, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateProject(p *PNRRProject) error {
	if p == nil {
		return errors.New("dati progetto non validi")
	}
	if p.SchoolID == "" || p.CUPCode == "" || p.ProjectTitle == "" {
		return errors.New("schoolID, CUP e titolo progetto sono obbligatori")
	}
	if p.TotalHours <= 0 {
		return errors.New("il totale delle ore del progetto deve essere maggiore di zero")
	}
	if p.Status == "" {
		p.Status = "attivo"
	}
	return s.repo.CreateProject(p)
}

func (s *service) ListProjects(schoolID, academicYear string) ([]PNRRProject, error) {
	if schoolID == "" {
		return nil, errors.New("schoolID obbligatorio")
	}
	return s.repo.ListProjectsBySchool(schoolID, academicYear)
}

func (s *service) GetProject(id string) (*PNRRProject, error) {
	if id == "" {
		return nil, errors.New("id progetto obbligatorio")
	}
	return s.repo.GetProject(id)
}

func (s *service) EnrollStudent(projectID, studentID string) error {
	if projectID == "" || studentID == "" {
		return errors.New("projectID e studentID obbligatori")
	}
	ps := &PNRRProjectStudent{
		ProjectID: projectID,
		StudentID: studentID,
	}
	return s.repo.EnrollStudent(ps)
}

func (s *service) ListStudents(projectID string) ([]PNRRProjectStudent, error) {
	if projectID == "" {
		return nil, errors.New("projectID obbligatorio")
	}
	return s.repo.ListProjectStudents(projectID)
}

func (s *service) RecordLesson(l *PNRRLessonAttendance, attendedStudentIDs []string) error {
	if l == nil {
		return errors.New("dati lezione non validi")
	}
	if l.ProjectID == "" || l.Hours <= 0 || l.Topic == "" {
		return errors.New("projectID, ore e argomento lezione sono obbligatori")
	}
	return s.repo.RecordLesson(l, attendedStudentIDs)
}

func (s *service) ListLessons(projectID string) ([]PNRRLessonAttendance, error) {
	if projectID == "" {
		return nil, errors.New("projectID obbligatorio")
	}
	return s.repo.ListLessons(projectID)
}
