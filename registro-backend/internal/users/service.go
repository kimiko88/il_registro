package users

type Service interface {
	GetAll() ([]UserResponse, error)
	GetByID(id string) (*UserResponse, error)
	Create(req CreateUserRequest) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAll() ([]UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []UserResponse
	for _, u := range users {
		responses = append(responses, UserResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Role:      u.Role,
		})
	}
	return responses, nil
}

func (s *service) GetByID(id string) (*UserResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil // Or specific error
	}

	return &UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      u.Role,
	}, nil
}

func (s *service) Create(req CreateUserRequest) error {
	user := &User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
	}
	return s.repo.Create(user)
}
