package classes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, c *Class) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockRepository) List(ctx context.Context, schoolID string) ([]Class, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).([]Class), args.Error(1)
}
func (m *MockRepository) Get(ctx context.Context, id string) (*Class, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Class), args.Error(1)
}
func (m *MockRepository) Update(ctx context.Context, c *Class) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepository) ListByTeacher(ctx context.Context, teacherUserID string) ([]Class, error) {
	args := m.Called(ctx, teacherUserID)
	return args.Get(0).([]Class), args.Error(1)
}

func TestService_CreateClass(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	req := CreateClassRequest{Name: "1A", Section: "A", AcademicYear: "2024"}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*classes.Class")).Return(nil)

	c, err := service.CreateClass(ctx, "school-1", req)
	assert.NoError(t, err)
	assert.Equal(t, "1A", c.Name)
	assert.Equal(t, "school-1", c.SchoolID)

	mockRepo.AssertExpectations(t)
}

func TestService_ListClasses(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	expected := []Class{{Name: "1A", ID: "c1"}}

	mockRepo.On("List", ctx, "school-1").Return(expected, nil)

	list, err := service.ListClasses(ctx, "school-1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "1A", list[0].Name)
}
