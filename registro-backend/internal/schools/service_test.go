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

func TestService_UpdateSchool(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	id := "school-1"
	name := "Updated Name"
	req := &UpdateSchoolRequest{Name: &name}

	// Success case
	mockRepo.On("Update", mock.Anything, id, req).Return(nil).Once()
	err := svc.Update(context.Background(), id, req)
	assert.NoError(t, err)

	// Error case
	mockRepo.On("Update", mock.Anything, id, req).Return(errors.New("db error")).Once()
	err = svc.Update(context.Background(), id, req)
	assert.Error(t, err)
}

func TestService_CreateSchool(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := &CreateSchoolRequest{
		Name: "Test School",
		Code: "TEST001",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*schools.School")).Return(nil)
	res, err := svc.Create(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Test School", res.Name)
}
