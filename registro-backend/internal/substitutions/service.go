package substitutions

import (
	"context"
	"fmt"
	"time"

	"registro-backend/internal/notifications"
)

type Service struct {
	repo     Repository
	notifSvc *notifications.Service
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("substitutions.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) SetNotificationService(ns *notifications.Service) {
	s.notifSvc = ns
}

func (s *Service) CreateSubstitution(ctx context.Context, schoolID string, req CreateSubstitutionRequest) (*Substitution, error) {
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

func (s *Service) ListByTeacher(ctx context.Context, teacherID string) ([]*Substitution, error) {
	return s.repo.ListByTeacher(ctx, teacherID)
}

func (s *Service) ListMyToday(ctx context.Context, teacherID string) ([]*Substitution, error) {
	todayStr := time.Now().Format("2006-01-02")
	subs, err := s.repo.ListByTeacher(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	var filtered []*Substitution
	for _, sub := range subs {
		if sub.Date.Format("2006-01-02") == todayStr {
			filtered = append(filtered, sub)
		}
	}
	return filtered, nil
}

func (s *Service) AssignSubstitute(ctx context.Context, id string, req AssignSubstituteRequest) error {
	return s.repo.AssignSubstitute(ctx, id, req.SubstituteTeacherID, req.Notes)
}

func (s *Service) ConfirmSubstitution(ctx context.Context, id string, teacherID string) error {
	return s.repo.ConfirmSubstitution(ctx, id, teacherID)
}
