package notes

import (
	"context"
	"errors"
	"strings"
	"time"

	"registro-backend/internal/users"
)

var (
	ErrUnauthorizedEdit    = errors.New("unauthorized: can only edit own notes")
	ErrUnauthorizedDelete  = errors.New("unauthorized: can only delete own notes")
	ErrUnauthorizedApprove = errors.New("unauthorized: only dirigenza or admin can approve notes")
	ErrUnauthorizedParent  = errors.New("unauthorized: parent must specify student_id")
	ErrNotGuardian         = errors.New("unauthorized: not a guardian of this student")
)

type Service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, uRepo users.Repository) *Service {
	if repo == nil {
		panic("notes.NewService: repo must not be nil")
	}
	if uRepo == nil {
		panic("notes.NewService: userRepo must not be nil — required for access and guardianship checks")
	}
	return &Service{repo: repo, userRepo: uRepo}
}

func (s *Service) CreateNote(ctx context.Context, teacherID, schoolID string, req CreateNoteRequest) (*StudentNote, error) {
	if req.ClassID != "" {
		isAssigned, err := s.repo.IsTeacherAssignedToClass(ctx, teacherID, req.ClassID)
		if err != nil {
			return nil, err
		}
		if !isAssigned {
			return nil, errors.New("forbidden: docente non assegnato alla classe dello studente")
		}
	}
	if req.TargetRole == "" {
		req.TargetRole = "all"
	}
	isApproved := req.Type != "disciplinary" // Disciplinary notes require principal/admin approval before being published
	n := &StudentNote{
		SchoolID:   schoolID,
		TeacherID:  teacherID,
		StudentID:  req.StudentID,
		ClassID:    req.ClassID,
		SubjectID:  req.SubjectID,
		Type:       req.Type,
		Note:       req.Note,
		Date:       req.Date,
		IsReserved: req.IsReserved,
		TargetRole: req.TargetRole,
		IsApproved: isApproved,
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) ApproveNote(ctx context.Context, actorID, actorRole, noteID string) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "vice_principal" {
		return ErrUnauthorizedApprove
	}
	n, err := s.repo.Get(ctx, noteID)
	if err != nil {
		return err
	}
	if n.IsApproved {
		return errors.New("note is already approved")
	}
	return s.repo.ApproveNote(ctx, noteID, actorID)
}

func (s *Service) UpdateNote(ctx context.Context, actorID, actorRole, noteID string, req UpdateNoteRequest) (*StudentNote, error) {
	n, err := s.repo.Get(ctx, noteID)
	if err != nil {
		return nil, err
	}
	isOwner := n.TeacherID == actorID
	isAdmin := actorRole == "admin" || actorRole == "superadmin" || actorRole == "principal" || actorRole == "vice_principal"
	if !isOwner && !isAdmin {
		return nil, ErrUnauthorizedEdit
	}
	if !isAdmin && isOwner && n.IsApproved {
		return nil, errors.New("forbidden: note già approvata dall'admin, non più modificabile dal docente")
	}

	if req.Type != "" {
		if req.Type != n.Type {
			n.IsApproved = (req.Type != "disciplinary")
		}
		n.Type = req.Type
	}
	if req.Note != "" {
		n.Note = req.Note
	}
	if req.Date != "" {
		n.Date = req.Date
	}
	if req.IsReserved != nil {
		n.IsReserved = *req.IsReserved
	}
	if req.TargetRole != "" {
		n.TargetRole = req.TargetRole
	}

	if err := s.repo.Update(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) DeleteNote(ctx context.Context, actorID, actorRole, noteID string, reason string) error {
	n, err := s.repo.Get(ctx, noteID)
	if err != nil {
		return err
	}
	isOwner := n.TeacherID == actorID
	isAdmin := actorRole == "admin" || actorRole == "superadmin" || actorRole == "principal" || actorRole == "vice_principal"
	if !isOwner && !isAdmin {
		return ErrUnauthorizedDelete
	}
	if !isAdmin && isOwner && n.IsApproved {
		return errors.New("forbidden: note già approvata dall'admin, non più eliminabile dal docente")
	}
	if isAdmin && n.Type == NoteTypeDisciplinary && strings.TrimSpace(reason) == "" {
		return errors.New("bad request: la cancellazione di una nota disciplinare da parte dell'admin richiede una motivazione")
	}
	return s.repo.DeleteWithReason(ctx, noteID, reason)
}

func (s *Service) ListNotes(ctx context.Context, filter NoteFilter) ([]StudentNote, error) {
	switch filter.ActorRole {
	case "student":
		filter.StudentID = filter.ActorID
	case "parent":
		if filter.StudentID == "" {
			return nil, ErrUnauthorizedParent
		}
		if s.userRepo != nil {
			isGuardian, err := s.userRepo.IsGuardian(ctx, filter.ActorID, filter.StudentID)
			if err != nil {
				return nil, err
			}
			if !isGuardian {
				return nil, ErrNotGuardian
			}
		}
	}
	notes, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if filter.ActorRole == "student" || filter.ActorRole == "parent" {
		var safeNotes []StudentNote
		for _, n := range notes {
			if !n.IsReserved && n.IsApproved {
				safeNotes = append(safeNotes, n)
			}
		}
		notes = safeNotes
	}

	if filter.ActorRole == "parent" {
		// Fix: usa una singola query batch invece di N query individuali all'interno di una GET.
		// Raccoglie gli ID delle note non ancora viste e li aggiorna in un unico UPDATE.
		var unviewedIDs []string
		now := time.Now()
		for i := range notes {
			if !notes[i].IsViewedByParent {
				unviewedIDs = append(unviewedIDs, notes[i].ID)
				notes[i].IsViewedByParent = true
				notes[i].ParentViewedAt = &now
			}
		}
		if len(unviewedIDs) > 0 {
			_ = s.repo.MarkManyAsViewedByParent(ctx, unviewedIDs)
		}
	}

	return notes, nil
}
