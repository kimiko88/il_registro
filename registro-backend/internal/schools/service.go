package schools

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

// Repository defines the interface for school storage
type Repository interface {
	Create(ctx context.Context, school *School) error
	GetByID(ctx context.Context, id string) (*School, error)
	List(ctx context.Context, params *ListParams) ([]*School, int, error)
	Update(ctx context.Context, id string, req *UpdateSchoolRequest) error
	Delete(ctx context.Context, id string) error
}

// Service handles business logic for schools.
//
// DEPRECATION NOTE (Bug 96): questo Service è un layer di basso livello senza
// validazione RBAC completa né check di duplicazione codice meccanografico.
// Per operazioni amministrative multi-tenant usare internal/admin.Service che
// implementa CreateSchool/UpdateSchool/DeleteSchool con audit log e RBAC completo.
type Service struct {
	repo Repository
}

// NewService creates a new schools service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new school
// Bug 95: validates required fields and email format before persisting.
func (s *Service) Create(ctx context.Context, req *CreateSchoolRequest) (*School, error) {
	// Required field validation
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return nil, errors.New("code (codice meccanografico) is required")
	}
	// Email format validation (optional field, but must be valid if provided)
	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return nil, fmt.Errorf("email format not valid: %w", err)
		}
	}

	school := &School{
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
		City:    req.City,
		Phone:   req.Phone,
		Email:   req.Email,
	}
	if err := s.repo.Create(ctx, school); err != nil {
		return nil, err
	}
	return school, nil
}

// GetByID retrieves a school by ID
func (s *Service) GetByID(ctx context.Context, id string) (*School, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves schools with pagination
func (s *Service) List(ctx context.Context, params *ListParams) ([]*School, int, error) {
	return s.repo.List(ctx, params)
}

// Update updates a school
func (s *Service) Update(ctx context.Context, id string, req *UpdateSchoolRequest) error {
	return s.repo.Update(ctx, id, req)
}

// Delete deletes a school
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
