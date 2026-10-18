package infanzia

import (
	"errors"
	"fmt"
)

type Service interface {
	CreateSection(s *InfanziaSection) error
	ListSections(schoolID, academicYear string) ([]InfanziaSection, error)
	GetSection(id string) (*InfanziaSection, error)

	SaveAttendances(sectionID, date string, attendances []InfanziaDailyAttendance) error
	GetAttendances(sectionID, date string) ([]InfanziaDailyAttendance, error)

	RecordObservation(obs *InfanziaExperienceObservation) error
	ListStudentObservations(studentID string) ([]InfanziaExperienceObservation, error)

	AddDiarioStory(d *InfanziaDiarioBordo) error
	ListDiarioStories(sectionID string, limit int) ([]InfanziaDiarioBordo, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateSection(sec *InfanziaSection) error {
	if sec == nil {
		return errors.New("dati sezione non validi")
	}
	if sec.SchoolID == "" || sec.SectionName == "" {
		return errors.New("schoolID e sectionName sono obbligatori")
	}
	if sec.ColorHex == "" {
		sec.ColorHex = "#EF4444"
	}
	if sec.AcademicYear == "" {
		sec.AcademicYear = "2025/2026"
	}
	return s.repo.CreateSection(sec)
}

func (s *service) ListSections(schoolID, academicYear string) ([]InfanziaSection, error) {
	if schoolID == "" {
		return nil, errors.New("schoolID obbligatorio")
	}
	return s.repo.ListSections(schoolID, academicYear)
}

func (s *service) GetSection(id string) (*InfanziaSection, error) {
	return s.repo.GetSection(id)
}

func (s *service) SaveAttendances(sectionID, date string, attendances []InfanziaDailyAttendance) error {
	if sectionID == "" || date == "" {
		return errors.New("sectionID e date sono obbligatori")
	}
	return s.repo.SaveDailyAttendancesBatch(sectionID, date, attendances)
}

func (s *service) GetAttendances(sectionID, date string) ([]InfanziaDailyAttendance, error) {
	if sectionID == "" || date == "" {
		return nil, errors.New("sectionID e date sono obbligatori")
	}
	return s.repo.GetDailyAttendances(sectionID, date)
}

func (s *service) RecordObservation(obs *InfanziaExperienceObservation) error {
	if obs == nil {
		return errors.New("osservazione non valida")
	}
	if obs.StudentID == "" || obs.RecordedBy == "" || obs.ObservationNotes == "" {
		return errors.New("studentID, recordedBy e observationNotes sono obbligatori")
	}
	validCampi := map[string]bool{
		"il_se_e_l_altro":         true,
		"il_corpo_e_il_movimento": true,
		"immagini_suoni_colori":   true,
		"i_discorsi_e_le_parole":  true,
		"la_conoscenza_del_mondo": true,
	}
	if !validCampi[obs.CampoEsperienza] {
		return fmt.Errorf("campo di esperienza non valido: %s", obs.CampoEsperienza)
	}
	if obs.LevelDescriptor == "" {
		obs.LevelDescriptor = "In evoluzione"
	}
	return s.repo.SaveObservation(obs)
}

func (s *service) ListStudentObservations(studentID string) ([]InfanziaExperienceObservation, error) {
	if studentID == "" {
		return nil, errors.New("studentID obbligatorio")
	}
	return s.repo.ListObservationsByStudent(studentID)
}

func (s *service) AddDiarioStory(d *InfanziaDiarioBordo) error {
	if d == nil {
		return errors.New("diario non valido")
	}
	if d.SectionID == "" || d.Title == "" || d.NarrativeStory == "" {
		return errors.New("sectionID, title e narrativeStory sono obbligatori")
	}
	return s.repo.CreateDiarioEntry(d)
}

func (s *service) ListDiarioStories(sectionID string, limit int) ([]InfanziaDiarioBordo, error) {
	if sectionID == "" {
		return nil, errors.New("sectionID obbligatorio")
	}
	if limit <= 0 {
		limit = 30
	} else if limit > 100 {
		limit = 100
	}
	return s.repo.ListDiarioEntries(sectionID, limit)
}
