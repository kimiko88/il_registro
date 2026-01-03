package schools

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, school *School) error {
	args := m.Called(ctx, school)
	return args.Error(0)
}
func (m *MockRepository) GetByID(ctx context.Context, id string) (*School, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*School), args.Error(1)
}
func (m *MockRepository) List(ctx context.Context, params *ListParams) ([]*School, int, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return []*School{}, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*School), args.Int(1), args.Error(2)
}
func (m *MockRepository) Update(ctx context.Context, id string, req *UpdateSchoolRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}
func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_CreateSchool(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateSchoolRequest
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful creation with all fields",
			req: &CreateSchoolRequest{
				Name:    "Test School",
				Code:    "TEST001",
				Address: "123 Main St",
				City:    "Rome",
				Phone:   "1234567890",
				Email:   "test@school.it",
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(s *School) bool {
					return s.Name == "Test School" && s.Code == "TEST001"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "successful creation with minimal fields",
			req: &CreateSchoolRequest{
				Name: "Minimal School",
				Code: "MIN001",
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*schools.School")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			req: &CreateSchoolRequest{
				Name: "Test School",
				Code: "TEST001",
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*schools.School")).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			svc := NewService(mockRepo)
			res, err := svc.Create(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.req.Name, res.Name)
				assert.Equal(t, tt.req.Code, res.Code)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful get",
			id:   "school-123",
			mockFn: func(m *MockRepository) {
				school := &School{ID: "school-123", Name: "Test School"}
				m.On("GetByID", mock.Anything, "school-123").Return(school, nil)
			},
			wantErr: false,
		},
		{
			name: "school not found",
			id:   "school-999",
			mockFn: func(m *MockRepository) {
				m.On("GetByID", mock.Anything, "school-999").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			svc := NewService(mockRepo)
			school, err := svc.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, school)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, school)
				assert.Equal(t, tt.id, school.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name      string
		params    *ListParams
		mockFn    func(*MockRepository)
		wantErr   bool
		wantCount int
		wantTotal int
	}{
		{
			name:   "successful list with results",
			params: &ListParams{Page: 1, PageSize: 10},
			mockFn: func(m *MockRepository) {
				schools := []*School{
					{ID: "school-1", Name: "School A"},
					{ID: "school-2", Name: "School B"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("*schools.ListParams")).Return(schools, 2, nil)
			},
			wantErr:   false,
			wantCount: 2,
			wantTotal: 2,
		},
		{
			name:   "empty list",
			params: &ListParams{Page: 1, PageSize: 10},
			mockFn: func(m *MockRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("*schools.ListParams")).Return([]*School{}, 0, nil)
			},
			wantErr:   false,
			wantCount: 0,
			wantTotal: 0,
		},
		{
			name:   "repository error",
			params: &ListParams{Page: 1, PageSize: 10},
			mockFn: func(m *MockRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("*schools.ListParams")).Return(nil, 0, errors.New("database error"))
			},
			wantErr:   true,
			wantCount: 0,
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			svc := NewService(mockRepo)
			schools, total, err := svc.List(context.Background(), tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, schools, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_UpdateSchool(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		req     *UpdateSchoolRequest
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful update",
			id:   "school-123",
			req: func() *UpdateSchoolRequest {
				name := "Updated School"
				return &UpdateSchoolRequest{Name: &name}
			}(),
			mockFn: func(m *MockRepository) {
				m.On("Update", mock.Anything, "school-123", mock.AnythingOfType("*schools.UpdateSchoolRequest")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			id:   "school-123",
			req: func() *UpdateSchoolRequest {
				name := "Updated School"
				return &UpdateSchoolRequest{Name: &name}
			}(),
			mockFn: func(m *MockRepository) {
				m.On("Update", mock.Anything, "school-123", mock.AnythingOfType("*schools.UpdateSchoolRequest")).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			svc := NewService(mockRepo)
			err := svc.Update(context.Background(), tt.id, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   "school-123",
			mockFn: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "school-123").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			id:   "school-123",
			mockFn: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "school-123").Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			svc := NewService(mockRepo)
			err := svc.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewService(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	assert.NotNil(t, svc)
	assert.Equal(t, mockRepo, svc.repo)
}
