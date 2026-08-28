package support

import (
	"context"
	"errors"

	"registro-backend/internal/users"
)

var (
	ErrEmptyTopic = errors.New("argomenti e attività del diario di sostegno obbligatori")
)

type Service interface {
	CreateDiaryEntry(ctx context.Context, schoolID, teacherID string, req CreateDiaryEntryRequest) (*SupportDiaryEntry, error)
	ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error)
	DeleteDiaryEntry(ctx context.Context, id, teacherID string) error

	CreatePeiGoal(ctx context.Context, schoolID string, req CreatePeiGoalRequest) (*SupportPeiGoal, error)
	UpdateGoalProgress(ctx context.Context, id, status string) error
	ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error)
	DeletePeiGoal(ctx context.Context, id string) error
	GetUserRepo() users.Repository
}

type service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, uRepo ...users.Repository) Service {
	s := &service{repo: repo}
	if len(uRepo) > 0 && uRepo[0] != nil {
		s.userRepo = uRepo[0]
	}
	return s
}

func (s *service) GetUserRepo() users.Repository {
	return s.userRepo
}

func (s *service) CreateDiaryEntry(ctx context.Context, schoolID, teacherID string, req CreateDiaryEntryRequest) (*SupportDiaryEntry, error) {
	if req.TopicAndActivities == "" {
		return nil, ErrEmptyTopic
	}
	if req.ActivityType == "" {
		req.ActivityType = "in_classe"
	}

	entry := &SupportDiaryEntry{
		SchoolID:           schoolID,
		TeacherID:          teacherID,
		StudentID:          req.StudentID,
		ClassID:            req.ClassID,
		EntryDate:          req.EntryDate,
		TimeSlot:           req.TimeSlot,
		CoTeacherID:        req.CoTeacherID,
		ActivityType:       req.ActivityType,
		TopicAndActivities: req.TopicAndActivities,
		StudentResponses:   req.StudentResponses,
		EducatorNotes:      req.EducatorNotes,
		IsSharedWithFamily: req.IsSharedWithFamily,
	}

	if err := s.repo.CreateDiaryEntry(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *service) ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error) {
	return s.repo.ListDiaryEntries(ctx, schoolID, studentID, teacherID, classID, isFamily)
}

func (s *service) DeleteDiaryEntry(ctx context.Context, id, teacherID string) error {
	return s.repo.DeleteDiaryEntry(ctx, id, teacherID)
}

func (s *service) CreatePeiGoal(ctx context.Context, schoolID string, req CreatePeiGoalRequest) (*SupportPeiGoal, error) {
	if req.PeiType == "" {
		req.PeiType = "equipollente"
	}
	if req.Axis == "" {
		req.Axis = "autonomia"
	}
	if req.ProgressStatus == "" {
		req.ProgressStatus = "non_avviato"
	}
	if req.ExpectedTerm == "" {
		req.ExpectedTerm = "annuale"
	}

	goal := &SupportPeiGoal{
		SchoolID:       schoolID,
		StudentID:      req.StudentID,
		PeiType:        req.PeiType,
		Axis:           req.Axis,
		Title:          req.Title,
		Description:    req.Description,
		ExpectedTerm:   req.ExpectedTerm,
		ProgressStatus: req.ProgressStatus,
	}

	if err := s.repo.CreatePeiGoal(ctx, goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *service) UpdateGoalProgress(ctx context.Context, id, status string) error {
	return s.repo.UpdatePeiGoalProgress(ctx, id, status)
}

func (s *service) ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error) {
	return s.repo.ListPeiGoals(ctx, schoolID, studentID)
}

func (s *service) DeletePeiGoal(ctx context.Context, id string) error {
	return s.repo.DeletePeiGoal(ctx, id)
}
