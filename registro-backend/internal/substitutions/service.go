package substitutions

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"registro-backend/internal/notifications"
)

type Service interface {
	CreateSubstitution(ctx context.Context, actorRole, schoolID string, req CreateSubstitutionRequest) (*Substitution, error)
	ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error)
	ListByTeacher(ctx context.Context, teacherID string, date string) ([]*Substitution, error)
	ListMyToday(ctx context.Context, teacherID string) ([]*Substitution, error)
	AssignSubstitute(ctx context.Context, id, actorRole string, req AssignSubstituteRequest) error
	ConfirmSubstitution(ctx context.Context, id string, teacherID string) error
	SignRegister(ctx context.Context, id string, teacherID string, notes string) error
	RecommendSubstitutes(ctx context.Context, schoolID, classID, subjectID string, date string, hour int) ([]SubstituteRecommendation, error)
	SetNotificationService(ns notifications.Service)
}

type serviceImpl struct {
	repo     Repository
	notifSvc notifications.Service
}

func NewService(repo Repository, notifSvc ...notifications.Service) Service {
	if repo == nil {
		panic("substitutions.NewService: repo must not be nil")
	}
	s := &serviceImpl{repo: repo}
	if len(notifSvc) > 0 {
		s.notifSvc = notifSvc[0]
	}
	return s
}

func (s *serviceImpl) SetNotificationService(ns notifications.Service) {
	s.notifSvc = ns
}

func (s *serviceImpl) CreateSubstitution(ctx context.Context, actorRole, schoolID string, req CreateSubstitutionRequest) (*Substitution, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "coordinator" && actorRole != "staff" {
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

func (s *serviceImpl) ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error) {
	return s.repo.ListBySchool(ctx, schoolID, date)
}

func (s *serviceImpl) ListByTeacher(ctx context.Context, teacherID string, date string) ([]*Substitution, error) {
	return s.repo.ListByTeacher(ctx, teacherID, date)
}

func (s *serviceImpl) ListMyToday(ctx context.Context, teacherID string) ([]*Substitution, error) {
	todayStr := time.Now().Format("2006-01-02")
	return s.repo.ListByTeacher(ctx, teacherID, todayStr)
}

// AssignSubstitute assigns a substitute teacher to a substitution.
func (s *serviceImpl) AssignSubstitute(ctx context.Context, id, actorRole string, req AssignSubstituteRequest) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "coordinator" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "staff" {
		return fmt.Errorf("unauthorized: solo admin, segreteria e coordinatori possono assegnare sostituzioni")
	}
	return s.repo.AssignSubstitute(ctx, id, req.SubstituteTeacherID, req.Notes)
}

// ConfirmSubstitution allows the assigned substitute teacher to confirm they accept the substitution.
func (s *serviceImpl) ConfirmSubstitution(ctx context.Context, id string, teacherID string) error {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("sostituzione non trovata: %w", err)
	}
	teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, teacherID)
	if sub.SubstituteTeacherID == nil || (*sub.SubstituteTeacherID != teacherID && *sub.SubstituteTeacherID != teacherProfileID) {
		return fmt.Errorf("unauthorized: non sei il docente sostituto assegnato a questa sostituzione")
	}
	return s.repo.ConfirmSubstitution(ctx, id, teacherID)
}

// SignRegister signature method for substitute teacher
func (s *serviceImpl) SignRegister(ctx context.Context, id string, teacherID string, notes string) error {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("sostituzione non trovata: %w", err)
	}
	teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, teacherID)
	if sub.SubstituteTeacherID == nil || (*sub.SubstituteTeacherID != teacherID && *sub.SubstituteTeacherID != teacherProfileID) {
		return fmt.Errorf("unauthorized: solo il docente sostituto assegnato può firmare il registro")
	}
	hashInput := fmt.Sprintf("FEQ-SUB-%s-%s-%d", id, teacherID, time.Now().UnixNano())
	sigHash := fmt.Sprintf("%x", sha256.Sum256([]byte(hashInput)))
	return s.repo.SignRegister(ctx, id, sigHash, notes)
}

// RecommendSubstitutes algorithm — Bug 149 real scoring calculation
func (s *serviceImpl) RecommendSubstitutes(ctx context.Context, schoolID, classID, subjectID string, date string, hour int) ([]SubstituteRecommendation, error) {
	candidates, err := s.repo.GetAvailableTeachers(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	var recs []SubstituteRecommendation
	for _, c := range candidates {
		score := 50
		var reasons []string

		teachesClass, _ := s.repo.IsTeacherAssignedToClass(ctx, c.TeacherID, classID)
		if teachesClass {
			score += 25
			reasons = append(reasons, "Docente della stessa classe")
		}

		teachesSubject, _ := s.repo.IsTeacherAssignedToSubject(ctx, c.TeacherID, subjectID)
		if teachesSubject {
			score += 15
			reasons = append(reasons, "Insegna la stessa materia")
		}

		subCount, _ := s.repo.GetWeeklySubstitutionCount(ctx, c.TeacherID)
		if subCount == 0 {
			score += 10
			reasons = append(reasons, "Nessuna supplenza svolta negli ultimi 7 giorni")
		} else if subCount <= 2 {
			score += 5
			reasons = append(reasons, fmt.Sprintf("Basso carico supplenze (%d questa settimana)", subCount))
		}

		if score > 95 {
			score = 95
		}

		reason := "Disponibile per supplenza"
		if len(reasons) > 0 {
			reason = fmt.Sprintf("%s", reasons[0])
			if len(reasons) > 1 {
				reason += fmt.Sprintf(" • %s", reasons[1])
			}
		}

		recs = append(recs, SubstituteRecommendation{
			TeacherID:      c.TeacherID,
			TeacherName:    c.TeacherName,
			Score:          score,
			Reason:         reason,
			IsFree:         true,
			TeachesClass:   teachesClass,
			TeachesSubject: teachesSubject,
		})
	}
	return recs, nil
}
