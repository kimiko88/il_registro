package teacher_activities

import (
	"errors"
	"fmt"
	"time"
)

// validActivityTypes elenca i tipi di attività consentiti.
var validActivityTypes = map[string]bool{
	"disponibilita": true,
	"riunione":      true,
	"formazione":    true,
	"ptof":          true,
	"gita":          true,
	"altro":         true,
}

// Service definisce la logica di business per le attività libere del docente.
type Service interface {
	Create(teacherID string, req CreateTeacherActivityRequest) (*TeacherActivityResponse, error)
	GetByID(id string) (*TeacherActivityResponse, error)
	Update(teacherID, id string, req UpdateTeacherActivityRequest) (*TeacherActivityResponse, error)
	Delete(teacherID, id string) error
	GetByTeacher(teacherID, fromDate, toDate string) ([]TeacherActivityResponse, error)
}

type service struct {
	repo Repository
}

// NewService restituisce una nuova implementazione del Service.
func NewService(r Repository) Service {
	if r == nil {
		panic("teacher_activities.NewService: repo must not be nil")
	}
	return &service{repo: r}
}

func (s *service) Create(teacherID string, req CreateTeacherActivityRequest) (*TeacherActivityResponse, error) {
	if err := s.validateRequest(req.Date, req.StartHour, req.Duration, req.ActivityType, req.Description); err != nil {
		return nil, err
	}

	actType := req.ActivityType
	if actType == "" {
		actType = "disponibilita"
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("formato data non valido: %w", err)
	}

	act := &TeacherFreeActivity{
		TeacherID:    teacherID,
		Date:         date,
		StartHour:    req.StartHour,
		Duration:     req.Duration,
		ActivityType: actType,
		Description:  req.Description,
		Notes:        req.Notes,
	}

	if err := s.repo.Create(act); err != nil {
		return nil, err
	}
	return s.toResponse(act), nil
}

func (s *service) GetByID(id string) (*TeacherActivityResponse, error) {
	act, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(act), nil
}

func (s *service) Update(teacherID, id string, req UpdateTeacherActivityRequest) (*TeacherActivityResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing.TeacherID != teacherID {
		return nil, errors.New("unauthorized: non puoi modificare l'attività di un altro docente")
	}
	if req.StartHour != nil && (*req.StartHour < 1 || *req.StartHour > 10) {
		return nil, errors.New("start_hour deve essere compreso tra 1 e 10")
	}
	if req.Duration != nil && *req.Duration < 1 {
		return nil, errors.New("duration deve essere almeno 1")
	}
	if req.ActivityType != "" && !validActivityTypes[req.ActivityType] {
		return nil, fmt.Errorf("activity_type non valido: %q", req.ActivityType)
	}

	act, err := s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}
	return s.toResponse(act), nil
}

func (s *service) Delete(teacherID, id string) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing.TeacherID != teacherID {
		return errors.New("unauthorized: non puoi eliminare l'attività di un altro docente")
	}
	return s.repo.Delete(id)
}

func (s *service) GetByTeacher(teacherID, fromDate, toDate string) ([]TeacherActivityResponse, error) {
	if fromDate != "" && toDate != "" {
		from, err1 := time.Parse("2006-01-02", fromDate)
		to, err2 := time.Parse("2006-01-02", toDate)
		if err1 == nil && err2 == nil {
			if from.After(to) {
				return nil, errors.New("fromDate non può essere successivo a toDate")
			}
			if to.Sub(from) > 365*24*time.Hour {
				return nil, errors.New("l'intervallo di date non può superare 1 anno")
			}
		}
	}
	activities, err := s.repo.GetByTeacher(teacherID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	var res []TeacherActivityResponse
	for _, a := range activities {
		res = append(res, *s.toResponse(&a))
	}
	return res, nil
}

// validateRequest esegue le validazioni di base comuni a Create e Update.
func (s *service) validateRequest(dateStr string, startHour, duration int, actType, description string) error {
	if description == "" {
		return errors.New("description è obbligatoria")
	}
	if startHour < 1 || startHour > 10 {
		return errors.New("start_hour deve essere compreso tra 1 e 10")
	}
	if duration < 1 {
		return errors.New("duration deve essere almeno 1")
	}
	if actType != "" && !validActivityTypes[actType] {
		return fmt.Errorf("activity_type non valido: %q", actType)
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("formato data non valido: %w", err)
	}
	now := time.Now()
	minDate := now.AddDate(-2, 0, 0)
	maxDate := now.AddDate(1, 0, 0)
	if date.Before(minDate) || date.After(maxDate) {
		return fmt.Errorf("data (%s) fuori dall'intervallo consentito (ultimi 2 anni – prossimo anno)", dateStr)
	}
	return nil
}

func (s *service) toResponse(a *TeacherFreeActivity) *TeacherActivityResponse {
	return &TeacherActivityResponse{
		ID:           a.ID,
		TeacherID:    a.TeacherID,
		TeacherName:  a.TeacherName,
		Date:         a.Date,
		StartHour:    a.StartHour,
		Duration:     a.Duration,
		ActivityType: a.ActivityType,
		Description:  a.Description,
		Notes:        a.Notes,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
