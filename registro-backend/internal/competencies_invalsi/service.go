package competencies_invalsi

import (
	"errors"
)

type Service interface {
	SaveCertification(c *CompetencyCertificationDM14) error
	GetCertification(id string) (*CompetencyCertificationDM14, error)
	GetCertificationByStudent(studentID, schoolYear, schoolTier string) (*CompetencyCertificationDM14, error)
	ListCertificationsByClass(classID, schoolYear string) ([]CompetencyCertificationDM14, error)
	SignCertification(id, principalID string) error

	ImportInvalsiBatch(results []InvalsiResult) error
	ListInvalsiByClass(classID, schoolYear string) ([]InvalsiResult, error)
	GetSchoolAnalytics(schoolID, schoolYear string) ([]InvalsiResult, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SaveCertification(c *CompetencyCertificationDM14) error {
	if c == nil {
		return errors.New("dati certificazione non validi")
	}
	if c.SchoolID == "" || c.StudentID == "" || c.ClassID == "" {
		return errors.New("schoolID, studentID e classID sono obbligatori")
	}
	if c.SchoolTier != "primaria" && c.SchoolTier != "primo_ciclo" {
		return errors.New("school_tier non valido (ammessi: primaria, primo_ciclo ex D.M. 14/2024)")
	}
	if c.Status == "" {
		c.Status = "bozza"
	}
	return s.repo.UpsertCertification(c)
}

func (s *service) GetCertification(id string) (*CompetencyCertificationDM14, error) {
	if id == "" {
		return nil, errors.New("id certificazione obbligatorio")
	}
	return s.repo.GetCertification(id)
}

func (s *service) GetCertificationByStudent(studentID, schoolYear, schoolTier string) (*CompetencyCertificationDM14, error) {
	if studentID == "" {
		return nil, errors.New("studentID obbligatorio")
	}
	return s.repo.GetCertificationByStudent(studentID, schoolYear, schoolTier)
}

func (s *service) ListCertificationsByClass(classID, schoolYear string) ([]CompetencyCertificationDM14, error) {
	if classID == "" {
		return nil, errors.New("classID obbligatorio")
	}
	return s.repo.ListCertificationsByClass(classID, schoolYear)
}

func (s *service) SignCertification(id, principalID string) error {
	if id == "" || principalID == "" {
		return errors.New("id certificazione e principalID obbligatori per la firma")
	}
	return s.repo.SignCertification(id, principalID)
}

func (s *service) ImportInvalsiBatch(results []InvalsiResult) error {
	if len(results) == 0 {
		return errors.New("nessun risultato INVALSI fornito per l'importazione")
	}
	for i := range results {
		if results[i].Level < 1 || results[i].Level > 5 {
			return errors.New("livello INVALSI non valido (deve essere compreso tra 1 e 5)")
		}
	}
	return s.repo.InsertInvalsiBatch(results)
}

func (s *service) ListInvalsiByClass(classID, schoolYear string) ([]InvalsiResult, error) {
	if classID == "" {
		return nil, errors.New("classID obbligatorio")
	}
	return s.repo.ListInvalsiByClass(classID, schoolYear)
}

func (s *service) GetSchoolAnalytics(schoolID, schoolYear string) ([]InvalsiResult, error) {
	if schoolID == "" {
		return nil, errors.New("schoolID obbligatorio")
	}
	return s.repo.GetSchoolInvalsiAnalytics(schoolID, schoolYear)
}
