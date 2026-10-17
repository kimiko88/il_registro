package enrollment

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo   Repository
	solver *FormationSolver
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:   repo,
		solver: NewFormationSolver(),
	}
}

func (s *Service) ImportSIDIApplications(ctx context.Context, schoolID, academicYear string, r io.Reader) (int, error) {
	apps, err := ParseSIDIExport(r)
	if err != nil {
		return 0, err
	}
	return s.repo.SaveApplications(ctx, apps, schoolID, academicYear)
}

func (s *Service) ListApplications(ctx context.Context, schoolID, academicYear, status string) ([]EnrollmentApplication, error) {
	if academicYear == "" {
		academicYear = "2026/2027"
	}
	return s.repo.ListApplications(ctx, schoolID, academicYear, status)
}

func (s *Service) GenerateFormationDraft(ctx context.Context, schoolID, academicYear, title string, params FormationParams) (*ClassFormationDraft, error) {
	if title == "" {
		title = "Bozza Formazione Classi Prime"
	}
	if academicYear == "" {
		academicYear = "2026/2027"
	}

	apps, err := s.repo.ListApplications(ctx, schoolID, academicYear, "pending")
	if err != nil {
		return nil, err
	}
	if len(apps) == 0 {
		return nil, errors.New("nessuna domanda d'iscrizione disponibile per la ripartizione")
	}

	solvedResult, err := s.solver.Solve(apps, params)
	if err != nil {
		return nil, err
	}

	draft := &ClassFormationDraft{
		ID:           uuid.New().String(),
		SchoolID:     schoolID,
		AcademicYear: academicYear,
		Title:        title,
		Parameters:   params,
		Assignments:  *solvedResult,
		IsFinalized:  false,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.SaveDraft(ctx, draft); err != nil {
		return nil, err
	}

	return draft, nil
}

func (s *Service) GetDraft(ctx context.Context, id string) (*ClassFormationDraft, error) {
	return s.repo.GetDraft(ctx, id)
}

func (s *Service) ListDrafts(ctx context.Context, schoolID string) ([]ClassFormationDraft, error) {
	return s.repo.ListDrafts(ctx, schoolID)
}

func (s *Service) UpdateDraftAssignments(ctx context.Context, id string, assignments ClassFormationDraftResult) error {
	return s.repo.UpdateDraftAssignments(ctx, id, assignments)
}

func (s *Service) FinalizeDraft(ctx context.Context, id string) error {
	return s.repo.FinalizeDraft(ctx, id)
}
