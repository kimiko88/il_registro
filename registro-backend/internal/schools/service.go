package schools

import "context"

// Repository defines the interface for school storage
type Repository interface {
	Create(ctx context.Context, school *School) error
	GetByID(ctx context.Context, id string) (*School, error)
	List(ctx context.Context, params *ListParams) ([]*School, int, error)
	Update(ctx context.Context, id string, req *UpdateSchoolRequest) error
	Delete(ctx context.Context, id string) error
}

// Service handles business logic for schools
type Service struct {
	repo Repository
}

// NewService creates a new schools service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new school
func (s *Service) Create(ctx context.Context, req *CreateSchoolRequest) (*School, error) {
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
