package teacher_hours

import (
	"errors"
	"fmt"
	"time"
)

type Service interface {
	GetTeacherDashboard(schoolID, teacherID, academicYear string) (*TeacherHourObligation, []TeacherHourLog, []TeacherShortLeave, error)
	LogHourActivity(teacherID string, req LogHourActivityRequest) (*TeacherHourLog, error)
	RequestShortLeave(schoolID, teacherID string, req ShortLeaveRequest) (*TeacherShortLeave, error)
	ReviewShortLeave(leaveID, status, notes string) error
	RecordLeaveRecovery(leaveID string, req RecoverLeaveRequest) error
	ListSchoolObligations(schoolID, academicYear string) ([]TeacherHourObligation, error)
}

type LogHourActivityRequest struct {
	ObligationID string  `json:"obligation_id" binding:"required"`
	Category     string  `json:"category" binding:"required"` // collegio, consigli, eccedenti, disposizione
	Hours        float64 `json:"hours" binding:"required"`
	ActivityDate string  `json:"activity_date" binding:"required"`
	MeetingID    *string `json:"meeting_id"`
	Description  string  `json:"description" binding:"required"`
}

type ShortLeaveRequest struct {
	RequestDate string  `json:"request_date" binding:"required"`
	Hours       float64 `json:"hours" binding:"required"`
	Reason      string  `json:"reason" binding:"required"`
}

type RecoverLeaveRequest struct {
	RecoveredHours float64 `json:"recovered_hours" binding:"required"`
	RecoveredDate  string  `json:"recovered_date" binding:"required"`
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTeacherDashboard(schoolID, teacherID, academicYear string) (*TeacherHourObligation, []TeacherHourLog, []TeacherShortLeave, error) {
	if schoolID == "" || teacherID == "" {
		return nil, nil, nil, errors.New("schoolID e teacherID sono obbligatori")
	}
	obl, err := s.repo.GetOrCreateObligation(schoolID, teacherID, academicYear)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("dashboard obligation: %w", err)
	}

	logs, err := s.repo.ListLogsByObligation(obl.ID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("dashboard logs: %w", err)
	}

	leaves, err := s.repo.ListShortLeaves(schoolID, teacherID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("dashboard leaves: %w", err)
	}

	return obl, logs, leaves, nil
}

func (s *service) LogHourActivity(teacherID string, req LogHourActivityRequest) (*TeacherHourLog, error) {
	if req.Hours <= 0 {
		return nil, errors.New("le ore devono essere maggiori di zero")
	}
	validCats := map[string]bool{
		"collegio": true, "consigli": true, "eccedenti": true,
		"disposizione": true, "permesso_breve": true, "recupero_permesso": true,
	}
	if !validCats[req.Category] {
		return nil, fmt.Errorf("categoria non valida: %s", req.Category)
	}

	log := &TeacherHourLog{
		ObligationID: req.ObligationID,
		TeacherID:    teacherID,
		Category:     req.Category,
		Hours:        req.Hours,
		ActivityDate: req.ActivityDate,
		MeetingID:    req.MeetingID,
		Description:  req.Description,
		Status:       "completed",
	}

	if err := s.repo.LogHourActivity(log); err != nil {
		return nil, fmt.Errorf("log activity: %w", err)
	}
	return log, nil
}

func (s *service) RequestShortLeave(schoolID, teacherID string, req ShortLeaveRequest) (*TeacherShortLeave, error) {
	if req.Hours <= 0 || req.Hours > 2.0 {
		return nil, errors.New("ai sensi del CCNL, i permessi brevi possono essere al massimo di 2 ore giornaliere")
	}
	if req.Reason == "" {
		return nil, errors.New("la motivazione è obbligatoria")
	}
	if req.RequestDate == "" {
		req.RequestDate = time.Now().Format("2006-01-02")
	}

	leave := &TeacherShortLeave{
		SchoolID:    schoolID,
		TeacherID:   teacherID,
		RequestDate: req.RequestDate,
		Hours:       req.Hours,
		Reason:      req.Reason,
		Status:      "requested",
	}

	if err := s.repo.CreateShortLeave(leave); err != nil {
		return nil, fmt.Errorf("create short leave: %w", err)
	}
	return leave, nil
}

func (s *service) ReviewShortLeave(leaveID, status, notes string) error {
	if leaveID == "" {
		return errors.New("leaveID obbligatorio")
	}
	if status != "approved" && status != "rejected" {
		return errors.New("stato non valido (ammessi: approved, rejected)")
	}
	return s.repo.UpdateShortLeaveStatus(leaveID, status, notes)
}

func (s *service) RecordLeaveRecovery(leaveID string, req RecoverLeaveRequest) error {
	if leaveID == "" {
		return errors.New("leaveID obbligatorio")
	}
	if req.RecoveredHours <= 0 {
		return errors.New("le ore recuperate devono essere maggiori di zero")
	}
	return s.repo.RecordShortLeaveRecovery(leaveID, req.RecoveredHours, req.RecoveredDate)
}

func (s *service) ListSchoolObligations(schoolID, academicYear string) ([]TeacherHourObligation, error) {
	if schoolID == "" {
		return nil, errors.New("schoolID obbligatorio")
	}
	return s.repo.ListObligationsBySchool(schoolID, academicYear)
}
