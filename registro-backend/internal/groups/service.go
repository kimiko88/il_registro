package groups

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrGroupNotFound = errors.New("group not found")
	ErrInvalidGroup  = errors.New("invalid group parameters")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateGroup(ctx context.Context, req CreateGroupRequest) (*Group, error) {
	if req.Name == "" || req.SchoolID == "" {
		return nil, ErrInvalidGroup
	}
	if req.AcademicYear == "" {
		req.AcademicYear = currentAcademicYear()
	}

	g := &Group{
		SchoolID:     req.SchoolID,
		Name:         req.Name,
		SubjectID:    req.SubjectID,
		TeacherID:    req.TeacherID,
		AcademicYear: req.AcademicYear,
		Description:  req.Description,
	}

	if err := s.repo.Create(ctx, g); err != nil {
		return nil, err
	}

	if len(req.StudentIDs) > 0 {
		_ = s.repo.AddStudents(ctx, g.ID, req.StudentIDs)
	}

	return s.GetGroupByID(ctx, g.ID)
}

func (s *Service) UpdateGroup(ctx context.Context, actorRole, schoolID, id string, req UpdateGroupRequest) (*Group, error) {
	g, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrGroupNotFound
	}
	if actorRole != "superadmin" && schoolID != "" && g.SchoolID != "" && g.SchoolID != schoolID {
		return nil, errors.New("forbidden: il gruppo appartiene ad un'altra scuola")
	}

	if req.Name != "" {
		g.Name = req.Name
	}
	if req.SubjectID != nil {
		g.SubjectID = req.SubjectID
	}
	if req.TeacherID != nil {
		g.TeacherID = req.TeacherID
	}
	if req.Description != "" {
		g.Description = req.Description
	}

	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}

	return s.GetGroupByID(ctx, g.ID)
}

func (s *Service) DeleteGroup(ctx context.Context, actorRole, schoolID, id string) error {
	g, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrGroupNotFound
	}
	if actorRole != "superadmin" && schoolID != "" && g.SchoolID != "" && g.SchoolID != schoolID {
		return errors.New("forbidden: il gruppo appartiene ad un'altra scuola")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetGroupByID(ctx context.Context, id string) (*Group, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListGroupsBySchool(ctx context.Context, schoolID string) ([]Group, error) {
	return s.repo.ListBySchool(ctx, schoolID)
}

func (s *Service) ListGroupsByTeacher(ctx context.Context, teacherID string) ([]Group, error) {
	return s.repo.ListByTeacher(ctx, teacherID)
}

func (s *Service) ListGroupsByStudent(ctx context.Context, studentID string) ([]Group, error) {
	return s.repo.ListByStudent(ctx, studentID)
}

func (s *Service) AddStudentsToGroup(ctx context.Context, actorRole, schoolID, groupID string, studentIDs []string) error {
	g, err := s.repo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}
	if actorRole != "superadmin" && schoolID != "" && g.SchoolID != "" && g.SchoolID != schoolID {
		return errors.New("forbidden: il gruppo appartiene ad un'altra scuola")
	}
	return s.repo.AddStudents(ctx, groupID, studentIDs)
}

func (s *Service) RemoveStudentFromGroup(ctx context.Context, actorRole, schoolID, groupID, studentID string) error {
	g, err := s.repo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}
	if actorRole != "superadmin" && schoolID != "" && g.SchoolID != "" && g.SchoolID != schoolID {
		return errors.New("forbidden: il gruppo appartiene ad un'altra scuola")
	}
	return s.repo.RemoveStudent(ctx, groupID, studentID)
}

// currentAcademicYear returns the school year label for today's date.
// Months September-December belong to the year starting in that calendar year;
// months January-August belong to the year that started in the previous calendar year.
func currentAcademicYear() string {
	now := time.Now()
	year := now.Year()
	if now.Month() < time.September {
		year--
	}
	return fmt.Sprintf("%d/%d", year, year+1)
}
