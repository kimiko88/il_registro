package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockRepository) GetByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}
func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}
func (m *MockRepository) Update(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepository) Restore(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepository) List(ctx context.Context, filter UserFilter) ([]User, int, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]User), args.Int(1), args.Error(2)
}
func (m *MockRepository) ListByIDs(ctx context.Context, ids []string) ([]User, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]User), args.Error(1)
}
func (m *MockRepository) LogAudit(ctx context.Context, log *AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}
func (m *MockRepository) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]AuditLog, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]AuditLog), args.Int(1), args.Error(2)
}
func (m *MockRepository) BulkCreate(ctx context.Context, users []User) (int, []string, error) {
	args := m.Called(ctx, users)
	return args.Int(0), args.Get(1).([]string), args.Error(2)
}
func (m *MockRepository) HardDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepository) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetChildren(ctx context.Context, parentID string) ([]StudentChild, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]StudentChild), args.Error(1)
}

func (m *MockRepository) IsActive(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestService_CreateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name      string
		actorRole string
		input     CreateUserRequest
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Admin can create user",
			actorRole: "admin",
			input: CreateUserRequest{
				Email: "test@example.com", Password: "Password123!", FirstName: "John", LastName: "Doe",
				FiscalCode: "RSSMRA80A01H501U", Role: "teacher",
			},
			mockSetup: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)
				mockRepo.On("LogAudit", mock.Anything, mock.AnythingOfType("*users.AuditLog")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Unauthorized role cannot create",
			actorRole: "teacher",
			input:     CreateUserRequest{Role: "student"},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:      "Weak password fails",
			actorRole: "admin",
			input: CreateUserRequest{
				Email: "t@e.com", Password: "weak", FiscalCode: "RSSMRA80A01H501U", Role: "teacher",
			},
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := service.CreateUser(context.Background(), tt.actorRole, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_ListUsers(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name      string
		actorRole string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Admin can list",
			actorRole: "admin",
			mockSetup: func() {
				mockRepo.On("List", mock.Anything, mock.AnythingOfType("users.UserFilter")).
					Return([]User{}, 0, nil)
			},
			wantErr: false,
		},
		{
			name:      "Student cannot list (unauthorized)",
			actorRole: "student", // Assumed from permissions.go which had no list perm
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, _, err := service.ListUsers(context.Background(), tt.actorRole, UserFilter{})
			if tt.wantErr {
				// We expect ErrUnauthorized
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
