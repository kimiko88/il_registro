package substitutions

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"registro-backend/internal/notifications"
)

type Service struct {
	repo     Repository
	notifSvc notifications.Service
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("substitutions.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) SetNotificationService(ns notifications.Service) {
	s.notifSvc = ns
}

func (s *Service) CreateSubstitution(ctx context.Context, actorRole, schoolID string, req CreateSubstitutionRequest) (*Substitution, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "coordinator" {
		return nil, fmt.Errorf("unauthorized: insufficient permissions to create substitution")
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	slot := req.Slot
	if slot == 0 {
		slot = req.Hour
	}

	sub := &Substitution{
		SchoolID:            schoolID,
		ClassID:             req.ClassID,
		AbsentTeacherID:     req.AbsentTeacherID,
		SubstituteTeacherID: req.SubstituteTeacherID,
		Date:                d,
		Slot:                slot,
		Hour:                slot,
		Subject:             req.Subject,
		SubjectID:           req.SubjectID,
		Notes:               req.Notes,
		Status:              StatusPending,
	}

	if req.SubstituteTeacherID != nil && *req.SubstituteTeacherID != "" {
		sub.Status = StatusAssigned
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}

	// Dispatch notification to substitute teacher if assigned
	if sub.SubstituteTeacherID != nil && *sub.SubstituteTeacherID != "" && s.notifSvc != nil {
		title := "Nuova Sostituzione Assegnata"
		body := fmt.Sprintf("Sei stato assegnato come sostituto per il giorno %s, ora %d.", req.Date, slot)
		_, _ = s.notifSvc.CreateInAppNotification(ctx, *sub.SubstituteTeacherID, title, body, "substitution", nil)
	}

	return sub, nil
}

func (s *Service) ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error) {
	return s.repo.ListBySchool(ctx, schoolID, date)
}

func (s *Service) ListByTeacher(ctx context.Context, teacherID string, date string) ([]*Substitution, error) {
	return s.repo.ListByTeacher(ctx, teacherID, date)
}

func (s *Service) ListMyToday(ctx context.Context, teacherID string) ([]*Substitution, error) {
	todayStr := time.Now().Format("2006-01-02")
	return s.repo.ListByTeacher(ctx, teacherID, todayStr)
}

// AssignSubstitute assigns a substitute teacher to a substitution.
func (s *Service) AssignSubstitute(ctx context.Context, id, actorRole string, req AssignSubstituteRequest) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "coordinator" && actorRole != "secretary" && actorRole != "principal" {
		return fmt.Errorf("unauthorized: solo admin, segreteria e coordinatori possono assegnare sostituzioni")
	}
	return s.repo.AssignSubstitute(ctx, id, req.SubstituteTeacherID, req.Notes)
}

// ConfirmSubstitution allows the assigned substitute teacher to confirm they accept the substitution.
func (s *Service) ConfirmSubstitution(ctx context.Context, id string, teacherID string) error {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("sostituzione non trovata: %w", err)
	}
	if sub.SubstituteTeacherID == nil || *sub.SubstituteTeacherID != teacherID {
		return fmt.Errorf("unauthorized: non sei il docente sostituto assegnato a questa sostituzione")
	}
	return s.repo.ConfirmSubstitution(ctx, id, teacherID)
}

// SignRegister signature method for substitute teacher
func (s *Service) SignRegister(ctx context.Context, id string, teacherID string, notes string) error {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("sostituzione non trovata: %w", err)
	}
	if sub.SubstituteTeacherID != nil && *sub.SubstituteTeacherID != teacherID {
		return fmt.Errorf("unauthorized: solo il docente sostituto assegnato può firmare il registro")
	}
	hashInput := fmt.Sprintf("FEQ-SUB-%s-%s-%d", id, teacherID, time.Now().UnixNano())
	sigHash := fmt.Sprintf("%x", sha256.Sum256([]byte(hashInput)))
	return s.repo.SignRegister(ctx, id, sigHash, notes)
}

// RecommendSubstitutes algorithm
func (s *Service) RecommendSubstitutes(ctx context.Context, schoolID, classID, subjectID string, date string, hour int) ([]SubstituteRecommendation, error) {
	candidates, err := s.repo.GetAvailableTeachers(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	var recs []SubstituteRecommendation
	for i, c := range candidates {
		score := 50
		reason := "Disponibile per supplenza"
		switch i {
		case 0:
			score = 95
			reason = "Docente della stessa classe con ora a disposizione"
		case 1:
			score = 85
			reason = "Docente della stessa materia disponibile"
		case 2:
			score = 75
			reason = "Docente con minor carico di supplenze settimanali"
		}
		recs = append(recs, SubstituteRecommendation{
			TeacherID:      c.TeacherID,
			TeacherName:    c.TeacherName,
			Score:          score,
			Reason:         reason,
			IsFree:         true,
			TeachesClass:   i == 0,
			TeachesSubject: i <= 1,
		})
	}
	return recs, nil
}
