package tenants

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateTenant(ctx context.Context, t *Tenant) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *MockRepository) GetTenantByID(ctx context.Context, id string) (*Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Tenant), args.Error(1)
}
func (m *MockRepository) ListTenants(ctx context.Context) ([]*Tenant, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Tenant), args.Error(1)
}

func TestTenantProvisioning(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateTenantRequest{
		Name:        "Liceo Galileo",
		Code:        "LG100",
		MaxStudents: 500,
	}

	mockRepo.On("CreateTenant", mock.Anything, mock.MatchedBy(func(tr *Tenant) bool {
		return tr.Name == "Liceo Galileo"
	})).Return(nil).Once()

	tr, err := svc.CreateTenant(context.Background(), "superadmin", req)
	assert.NoError(t, err)
	assert.NotNil(t, tr)

	_, err = svc.CreateTenant(context.Background(), "teacher", req)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
