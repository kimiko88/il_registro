package security_passes

import (
	"errors"
	"time"
)

type Service interface {
	IssueHallPass(p *StudentHallPass) error
	ReturnHallPass(id string) error
	GetActivePasses(schoolID string) ([]StudentHallPass, error)
	GetClassPasses(classID string, limit int) ([]StudentHallPass, error)

	AddPickupDelegate(d *StudentPickupDelegate) error
	GetStudentDelegates(studentID string) ([]StudentPickupDelegate, error)
	RevokePickupDelegate(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) IssueHallPass(p *StudentHallPass) error {
	if p.SchoolID == "" || p.ClassID == "" || p.StudentID == "" || p.TeacherID == "" {
		return errors.New("parametri obbligatori mancanti per il rilascio del pass")
	}
	if p.PassType == "" {
		p.PassType = "bagno"
	}
	p.Status = "fuori_aula"
	return s.repo.CreateHallPass(p)
}

func (s *service) ReturnHallPass(id string) error {
	if id == "" {
		return errors.New("id del pass obbligatorio")
	}
	return s.repo.CloseHallPass(id, time.Now())
}

func (s *service) GetActivePasses(schoolID string) ([]StudentHallPass, error) {
	if schoolID == "" {
		return nil, errors.New("schoolID obbligatorio")
	}
	return s.repo.ListActiveHallPassesBySchool(schoolID)
}

func (s *service) GetClassPasses(classID string, limit int) ([]StudentHallPass, error) {
	if classID == "" {
		return nil, errors.New("classID obbligatorio")
	}
	return s.repo.ListHallPassesByClass(classID, limit)
}

func (s *service) AddPickupDelegate(d *StudentPickupDelegate) error {
	if d.StudentID == "" || d.DelegateName == "" || d.FiscalCode == "" || d.Phone == "" {
		return errors.New("studentID, delegateName, fiscalCode e phone sono obbligatori")
	}
	if d.Status == "" {
		d.Status = "attivo"
	}
	return s.repo.CreatePickupDelegate(d)
}

func (s *service) GetStudentDelegates(studentID string) ([]StudentPickupDelegate, error) {
	if studentID == "" {
		return nil, errors.New("studentID obbligatorio")
	}
	return s.repo.ListDelegatesByStudent(studentID)
}

func (s *service) RevokePickupDelegate(id string) error {
	if id == "" {
		return errors.New("id delegato obbligatorio")
	}
	return s.repo.RevokeDelegate(id)
}
