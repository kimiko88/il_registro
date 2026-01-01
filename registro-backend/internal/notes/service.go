package notes

import (
	"context"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateNote(ctx context.Context, teacherID, schoolID string, req CreateNoteRequest) (*StudentNote, error) {
	// Validate type? binding already does.
	n := &StudentNote{
		SchoolID:  schoolID,
		TeacherID: teacherID,
		StudentID: req.StudentID,
		ClassID:   req.ClassID,
		SubjectID: req.SubjectID,
		Type:      req.Type,
		Note:      req.Note,
		Date:      req.Date,
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) UpdateNote(ctx context.Context, teacherID, noteID string, req UpdateNoteRequest) (*StudentNote, error) {
	// Get note first to check permission
	n, err := s.repo.Get(ctx, noteID)
	if err != nil {
		return nil, err
	}
	if n.TeacherID != teacherID {
		return nil, errors.New("unauthorized: can only edit own notes")
	}

	if req.Type != "" {
		n.Type = req.Type
	}
	if req.Note != "" {
		n.Note = req.Note
	}
	if req.Date != "" {
		n.Date = req.Date
	}

	if err := s.repo.Update(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) DeleteNote(ctx context.Context, teacherID, noteID string) error {
	n, err := s.repo.Get(ctx, noteID)
	if err != nil {
		return err
	}
	if n.TeacherID != teacherID {
		return errors.New("unauthorized: can only delete own notes")
	}
	return s.repo.Delete(ctx, noteID)
}

func (s *Service) ListNotes(ctx context.Context, filter NoteFilter) ([]StudentNote, error) {
	return s.repo.List(ctx, filter)
}
