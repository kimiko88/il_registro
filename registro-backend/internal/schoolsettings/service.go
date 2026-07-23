package schoolsettings

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSettings(ctx context.Context, schoolID string) (*SchoolSettings, error) {
	return s.repo.GetBySchoolID(ctx, schoolID)
}

func (s *Service) UpdateSettings(ctx context.Context, schoolID string, req UpdateSchoolSettingsRequest) (*SchoolSettings, error) {
	current, err := s.repo.GetBySchoolID(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	if req.RequirePrincipalApprovalForNotes != nil {
		current.RequirePrincipalApprovalForNotes = *req.RequirePrincipalApprovalForNotes
	}
	if req.AllowParentsViewGrades != nil {
		current.AllowParentsViewGrades = *req.AllowParentsViewGrades
	}
	if req.AllowStudentsViewClassAverages != nil {
		current.AllowStudentsViewClassAverages = *req.AllowStudentsViewClassAverages
	}
	if req.RequireMFAForStaff != nil {
		current.RequireMFAForStaff = *req.RequireMFAForStaff
	}
	if req.LockScrutinyEditingAfterValidation != nil {
		current.LockScrutinyEditingAfterValidation = *req.LockScrutinyEditingAfterValidation
	}
	if req.EnableSubstituteNotifications != nil {
		current.EnableSubstituteNotifications = *req.EnableSubstituteNotifications
	}
	current.UpdatedAt = time.Now()

	if err := s.repo.Upsert(ctx, current); err != nil {
		return nil, err
	}

	return current, nil
}
