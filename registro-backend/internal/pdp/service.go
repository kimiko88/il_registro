package pdp

import (
	"context"
	"errors"
	"fmt"

	"registro-backend/internal/users"
)

var (
	ErrPlanNotFound    = errors.New("pdp plan not found")
	ErrNotSharedYet    = errors.New("pdp plan has not been shared with the family yet")
	ErrAlreadyApproved = errors.New("pdp plan already approved by the family")
	ErrUnauthorized    = errors.New("unauthorized: insufficient role for this pdp operation")
	ErrNotGuardian     = errors.New("unauthorized: non sei il tutore legale dello studente")
)

// Service defines the business logic for PDP/PEI plans.
type Service struct {
	repo     Repository
	userRepo users.Repository
}

// NewService creates a new PDP service.
func NewService(repo Repository, uRepo ...users.Repository) *Service {
	svc := &Service{repo: repo}
	if len(uRepo) > 0 && uRepo[0] != nil {
		svc.userRepo = uRepo[0]
	}
	return svc
}

// CreatePlan creates a new PDP/PEI plan.
// Only coordinators, referenti or dirigenza may create plans.
func (s *Service) CreatePlan(ctx context.Context, actorID, actorRole, schoolID string, req *CreatePdpRequest) (*PdpPlan, error) {
	if !canManagePDP(actorRole) {
		return nil, ErrUnauthorized
	}
	if req.PlanType == "" {
		req.PlanType = PlanTypePDP
	}

	plan := &PdpPlan{
		StudentID:     req.StudentID,
		ClassID:       req.ClassID,
		SchoolID:      schoolID,
		AcademicYear:  req.AcademicYear,
		PlanType:      req.PlanType,
		Diagnosis:     req.Diagnosis,
		Content:       req.Content,
		CoordinatorID: req.CoordinatorID,
		ReferenteID:   req.ReferenteID,
		CreatedBy:     actorID,
	}

	// SchoolID must be populated — caller sets it from the auth context
	return s.repo.Create(ctx, plan)
}

// GetByStudent returns all PDP/PEI plans for a student for a given academic year.
// Parents only see plans that have been shared with the family, and guardianship is verified.
func (s *Service) GetByStudent(ctx context.Context, actorID, actorRole, actorSchoolID, studentID, academicYear string) ([]*PdpPlan, error) {
	switch actorRole {
	case "parent":
		if s.userRepo != nil {
			isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, studentID)
			if err != nil || !isGuardian {
				return nil, ErrNotGuardian
			}
		}
	case "student":
		if actorID != studentID {
			return nil, ErrUnauthorized
		}
	}

	plans, err := s.repo.GetByStudent(ctx, studentID, academicYear)
	if err != nil {
		return nil, err
	}

	// Multi-tenant check & privacy filter
	var visible []*PdpPlan
	for _, p := range plans {
		if actorRole != "superadmin" && actorSchoolID != "" && p.SchoolID != "" && p.SchoolID != actorSchoolID {
			continue
		}
		if actorRole == "parent" || actorRole == "student" {
			if !p.SharedWithFamily {
				continue
			}
			p.Diagnosis = "" // never expose diagnosis to non-staff
		}
		visible = append(visible, p)
	}
	return visible, nil
}

// GetByClass returns all PDP/PEI plans for an entire class.
// Only teachers, coordinators and dirigenza can access this.
func (s *Service) GetByClass(ctx context.Context, actorRole, classID, academicYear string) ([]*PdpPlan, error) {
	if !canViewClassPDP(actorRole) {
		return nil, ErrUnauthorized
	}
	return s.repo.GetByClass(ctx, classID, academicYear)
}

// GetByID returns a single plan by ID.
func (s *Service) GetByID(ctx context.Context, actorID, actorRole, actorSchoolID, id string) (*PdpPlan, error) {
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, ErrPlanNotFound
	}
	if actorRole != "superadmin" && actorSchoolID != "" && plan.SchoolID != "" && plan.SchoolID != actorSchoolID {
		return nil, ErrUnauthorized
	}

	switch actorRole {
	case "parent":
		if s.userRepo != nil {
			isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, plan.StudentID)
			if err != nil || !isGuardian {
				return nil, ErrNotGuardian
			}
		}
		if !plan.SharedWithFamily {
			return nil, ErrNotSharedYet
		}
		plan.Diagnosis = ""
	case "student":
		if actorID != plan.StudentID {
			return nil, ErrUnauthorized
		}
		if !plan.SharedWithFamily {
			return nil, ErrNotSharedYet
		}
		plan.Diagnosis = ""
	}

	return plan, nil
}

// UpdatePlan updates an existing plan. Only staff roles may update.
func (s *Service) UpdatePlan(ctx context.Context, actorRole, actorSchoolID, id string, req *UpdatePdpRequest) (*PdpPlan, error) {
	if !canManagePDP(actorRole) {
		return nil, ErrUnauthorized
	}
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil || plan == nil {
		return nil, ErrPlanNotFound
	}
	if actorRole != "superadmin" && actorSchoolID != "" && plan.SchoolID != "" && plan.SchoolID != actorSchoolID {
		return nil, ErrUnauthorized
	}
	return s.repo.Update(ctx, id, req)
}

// ShareWithFamily sets the shared_with_family flag.
func (s *Service) ShareWithFamily(ctx context.Context, actorRole, actorSchoolID, id string, share bool) (*PdpPlan, error) {
	if !canManagePDP(actorRole) {
		return nil, ErrUnauthorized
	}
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil || plan == nil {
		return nil, ErrPlanNotFound
	}
	if actorRole != "superadmin" && actorSchoolID != "" && plan.SchoolID != "" && plan.SchoolID != actorSchoolID {
		return nil, ErrUnauthorized
	}
	if err := s.repo.SetSharedWithFamily(ctx, id, share); err != nil {
		return nil, fmt.Errorf("pdp.ShareWithFamily: %w", err)
	}
	return s.repo.GetByID(ctx, id)
}

// ApproveByFamily records family approval of the plan.
func (s *Service) ApproveByFamily(ctx context.Context, actorRole, actorID, actorSchoolID, id string) error {
	if actorRole != "parent" && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil || plan == nil {
		return ErrPlanNotFound
	}
	if actorRole != "superadmin" && actorSchoolID != "" && plan.SchoolID != "" && plan.SchoolID != actorSchoolID {
		return ErrUnauthorized
	}
	if actorRole == "parent" {
		if s.userRepo != nil {
			isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, plan.StudentID)
			if err != nil || !isGuardian {
				return ErrNotGuardian
			}
		}
	}
	if !plan.SharedWithFamily {
		return ErrNotSharedYet
	}
	if plan.FamilyApprovedAt != nil {
		return ErrAlreadyApproved
	}
	return s.repo.ApproveByFamily(ctx, id, actorID)
}

// DeletePlan deletes a plan. Only coordinators and dirigenza can delete.
func (s *Service) DeletePlan(ctx context.Context, actorRole, actorSchoolID, id string) error {
	if !canManagePDP(actorRole) {
		return ErrUnauthorized
	}
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil || plan == nil {
		return ErrPlanNotFound
	}
	if actorRole != "superadmin" && actorSchoolID != "" && plan.SchoolID != "" && plan.SchoolID != actorSchoolID {
		return ErrUnauthorized
	}
	return s.repo.Delete(ctx, id)
}

// ─── Role Helpers ────────────────────────────────────────────────────────────

func canManagePDP(role string) bool {
	switch role {
	case "teacher", "secretary", "principal", "vice_principal", "admin", "superadmin":
		return true
	}
	return false
}

func canViewClassPDP(role string) bool {
	switch role {
	case "teacher", "secretary", "principal", "vice_principal", "admin", "superadmin":
		return true
	}
	return false
}
