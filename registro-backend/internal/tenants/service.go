package tenants

import (
	"context"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("tenants.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) CreateTenant(ctx context.Context, actorRole string, req CreateTenantRequest) (*Tenant, error) {
	if actorRole != "superadmin" {
		return nil, errors.New("forbidden: only superadmin can provision tenants")
	}

	t := &Tenant{
		Name: req.Name,
		Code: req.Code,
		Quota: TenantQuota{
			MaxStudents: req.MaxStudents,
			MaxTeachers: req.MaxTeachers,
			MaxStorageMB: req.MaxStorageMB,
		},
	}

	if err := s.repo.CreateTenant(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTenants(ctx context.Context, actorRole string) ([]*Tenant, error) {
	if actorRole != "superadmin" {
		return nil, errors.New("forbidden: only superadmin can list tenants")
	}
	return s.repo.ListTenants(ctx)
}

func (s *Service) GetTenant(ctx context.Context, id string) (*Tenant, error) {
	return s.repo.GetTenantByID(ctx, id)
}
