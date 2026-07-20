package groups

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, g *Group) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}
func (m *MockRepository) Update(ctx context.Context, g *Group) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}
func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepository) GetByID(ctx context.Context, id string) (*Group, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Group), args.Error(1)
}
func (m *MockRepository) ListBySchool(ctx context.Context, schoolID string) ([]Group, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Group), args.Error(1)
}
func (m *MockRepository) ListByTeacher(ctx context.Context, teacherID string) ([]Group, error) {
	args := m.Called(ctx, teacherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Group), args.Error(1)
}
func (m *MockRepository) ListByStudent(ctx context.Context, studentID string) ([]Group, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Group), args.Error(1)
}
func (m *MockRepository) AddStudents(ctx context.Context, groupID string, studentIDs []string) error {
	args := m.Called(ctx, groupID, studentIDs)
	return args.Error(0)
}
func (m *MockRepository) RemoveStudent(ctx context.Context, groupID, studentID string) error {
	args := m.Called(ctx, groupID, studentID)
	return args.Error(0)
}
func (m *MockRepository) GetStudentsInGroup(ctx context.Context, groupID string) ([]GroupStudentInfo, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]GroupStudentInfo), args.Error(1)
}

func TestGroupService_CreateGroup(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := CreateGroupRequest{
		SchoolID:     "school-123",
		Name:         "Inglese B2",
		AcademicYear: "2025/2026",
		StudentIDs:   []string{"stu-1", "stu-2"},
	}

	expectedGroup := &Group{
		ID:           "group-1",
		SchoolID:     "school-123",
		Name:         "Inglese B2",
		AcademicYear: "2025/2026",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*groups.Group")).Return(nil)
	mockRepo.On("AddStudents", mock.Anything, mock.Anything, req.StudentIDs).Return(nil)
	mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(expectedGroup, nil)

	res, err := service.CreateGroup(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Inglese B2", res.Name)
}
