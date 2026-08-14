package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockRepository) GetPasswordHistory(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockRepository) AddPasswordHistory(ctx context.Context, userID, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
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
func (m *MockRepository) BulkDelete(ctx context.Context, ids []string) (int, error) {
	args := m.Called(ctx, ids)
	return args.Int(0), args.Error(1)
}
func (m *MockRepository) ClearTempMFASecret(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockRepository) HardDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "RevokeAllUserTokens" {
			args := m.Called(ctx, userID)
			return args.Error(0)
		}
	}
	return nil
}
func (m *MockRepository) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetChildren(ctx context.Context, parentID string) ([]StudentChild, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]StudentChild), args.Error(1)
}

func (m *MockRepository) GetStudentsByClass(ctx context.Context, classID string) ([]User, error) {
	args := m.Called(ctx, classID)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockRepository) IsActive(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) AddGuardian(ctx context.Context, studentProfileID, parentProfileID, relationship string) error {
	args := m.Called(ctx, studentProfileID, parentProfileID, relationship)
	return args.Error(0)
}

func (m *MockRepository) GetStudentProfile(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) GetParentProfile(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) RemoveGuardian(ctx context.Context, studentProfileID, parentProfileID string) error {
	args := m.Called(ctx, studentProfileID, parentProfileID)
	return args.Error(0)
}

func (m *MockRepository) GetGuardians(ctx context.Context, studentProfileID string) ([]GuardianInfo, error) {
	args := m.Called(ctx, studentProfileID)
	return args.Get(0).([]GuardianInfo), args.Error(1)
}

func (m *MockRepository) GetFascicoloSummary(ctx context.Context, studentID string, isActive bool) (map[string]interface{}, error) {
	args := m.Called(ctx, studentID, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
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
				mockRepo.On("AddPasswordHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
			_, _, err := service.ListUsers(context.Background(), tt.actorRole, "", UserFilter{})
			if tt.wantErr {
				// We expect ErrUnauthorized
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name      string
		actorRole string
		userID    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Admin can get user",
			actorRole: "admin",
			userID:    "user-123",
			mockSetup: func() {
				user := &User{ID: "user-123", Email: "test@example.com"}
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)
			},
			wantErr: false,
		},
		{
			name:      "Teacher can get user (has UserRead permission)",
			actorRole: "teacher",
			userID:    "user-456",
			mockSetup: func() {
				user := &User{ID: "user-456", Email: "student@example.com"}
				mockRepo.On("GetByID", mock.Anything, "user-456").Return(user, nil)
			},
			wantErr: false,
		},
		{
			name:      "Student cannot get user  (no UserRead permission)",
			actorRole: "student",
			userID:    "user-789",
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := service.GetUser(context.Background(), tt.actorRole, "", tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	newName := "Updated Name"

	tests := []struct {
		name      string
		actorRole string
		userID    string
		req       UpdateUserRequest
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Admin can update user",
			actorRole: "admin",
			userID:    "user-123",
			req:       UpdateUserRequest{FirstName: &newName},
			mockSetup: func() {
				user := &User{ID: "user-123", FirstName: "Old Name"}
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Teacher cannot update (no UserUpdate permission)",
			actorRole: "teacher",
			userID:    "user-123",
			req:       UpdateUserRequest{FirstName: &newName},
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := service.UpdateUser(context.Background(), tt.actorRole, "", tt.userID, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name      string
		actorRole string
		userID    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Admin can delete user",
			actorRole: "admin",
			userID:    "user-123",
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(&User{ID: "user-123"}, nil)
				mockRepo.On("Delete", mock.Anything, "user-123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Teacher cannot delete (no UserDelete permission)",
			actorRole: "teacher",
			userID:    "user-123",
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.DeleteUser(context.Background(), tt.actorRole, "", tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_ChangePassword(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Create a hashed password for "OldPassword123!"
	oldHashedPassword, _ := bcrypt.GenerateFromPassword([]byte("OldPassword123!"), bcrypt.MinCost)

	tests := []struct {
		name      string
		userID    string
		req       ChangePasswordRequest
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "Successful password change",
			userID: "user-123",
			req: ChangePasswordRequest{
				CurrentPassword: "OldPassword123!",
				NewPassword:     "NewPassword123!",
			},
			mockSetup: func() {
				user := &User{ID: "user-123", PasswordHash: string(oldHashedPassword)}
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)
				mockRepo.On("GetPasswordHistory", mock.Anything, "user-123").Return([]string{}, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)
				mockRepo.On("AddPasswordHistory", mock.Anything, "user-123", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Incorrect current password",
			userID: "user-123",
			req: ChangePasswordRequest{
				CurrentPassword: "WrongPassword!",
				NewPassword:     "NewPassword123!",
			},
			mockSetup: func() {
				user := &User{ID: "user-123", PasswordHash: string(oldHashedPassword)}
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)
			},
			wantErr: true,
		},
		{
			name:   "Weak new password",
			userID: "user-123",
			req: ChangePasswordRequest{
				CurrentPassword: "OldPassword123!",
				NewPassword:     "weak",
			},
			mockSetup: func() {
				user := &User{ID: "user-123", PasswordHash: string(oldHashedPassword)}
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)
			},
			wantErr: true,
		},
		{
			name:   "Password missing uppercase character",
			userID: "user-123",
			req: ChangePasswordRequest{
				CurrentPassword: "OldPassword123!",
				NewPassword:     "newpassword123!",
			},
			mockSetup: func() {
				user := &User{ID: "user-123", PasswordHash: string(oldHashedPassword)}
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)
			},
			wantErr: true,
		},
		{
			name:   "Password missing special character",
			userID: "user-123",
			req: ChangePasswordRequest{
				CurrentPassword: "OldPassword123!",
				NewPassword:     "NewPassword1234",
			},
			mockSetup: func() {
				user := &User{ID: "user-123", PasswordHash: string(oldHashedPassword)}
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.ChangePassword(context.Background(), tt.userID, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_ResetPassword_SecretaryPermissions(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	schoolID := "school-1"

	tests := []struct {
		name        string
		actorRole   string
		targetRole  string
		targetID    string
		newPassword string
		mockSetup   func(targetID, targetRole string)
		wantErr     bool
	}{
		{
			name:        "Secretary resets teacher password -> Allowed",
			actorRole:   "secretary",
			targetRole:  "teacher",
			targetID:    "teacher-1",
			newPassword: "NewTeacherPass123!",
			mockSetup: func(tid, role string) {
				user := &User{ID: tid, Role: role, SchoolID: &schoolID}
				mockRepo.On("GetByID", mock.Anything, tid).Return(user, nil)
				mockRepo.On("GetPasswordHistory", mock.Anything, tid).Return([]string{}, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)
				mockRepo.On("AddPasswordHistory", mock.Anything, tid, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Secretary resets student password -> Allowed",
			actorRole:   "secretary",
			targetRole:  "student",
			targetID:    "student-1",
			newPassword: "NewStudentPass123!",
			mockSetup: func(tid, role string) {
				user := &User{ID: tid, Role: role, SchoolID: &schoolID}
				mockRepo.On("GetByID", mock.Anything, tid).Return(user, nil)
				mockRepo.On("GetPasswordHistory", mock.Anything, tid).Return([]string{}, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)
				mockRepo.On("AddPasswordHistory", mock.Anything, tid, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Secretary resets parent password -> Allowed",
			actorRole:   "secretary",
			targetRole:  "parent",
			targetID:    "parent-1",
			newPassword: "NewParentPass123!",
			mockSetup: func(tid, role string) {
				user := &User{ID: tid, Role: role, SchoolID: &schoolID}
				mockRepo.On("GetByID", mock.Anything, tid).Return(user, nil)
				mockRepo.On("GetPasswordHistory", mock.Anything, tid).Return([]string{}, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)
				mockRepo.On("AddPasswordHistory", mock.Anything, tid, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Secretary attempts to reset admin password -> Forbidden",
			actorRole:   "secretary",
			targetRole:  "admin",
			targetID:    "admin-1",
			newPassword: "NewAdminPass123!",
			mockSetup: func(tid, role string) {
				user := &User{ID: tid, Role: role, SchoolID: &schoolID}
				mockRepo.On("GetByID", mock.Anything, tid).Return(user, nil)
			},
			wantErr: true,
		},
		{
			name:        "Secretary attempts to reset superadmin password -> Forbidden",
			actorRole:   "secretary",
			targetRole:  "superadmin",
			targetID:    "superadmin-1",
			newPassword: "NewSuperAdminPass123!",
			mockSetup: func(tid, role string) {
				user := &User{ID: tid, Role: role}
				mockRepo.On("GetByID", mock.Anything, tid).Return(user, nil)
			},
			wantErr: true,
		},
		{
			name:        "Secretary attempts to reset another secretary password -> Forbidden",
			actorRole:   "secretary",
			targetRole:  "secretary",
			targetID:    "secretary-2",
			newPassword: "NewSecretaryPass123!",
			mockSetup: func(tid, role string) {
				user := &User{ID: tid, Role: role, SchoolID: &schoolID}
				mockRepo.On("GetByID", mock.Anything, tid).Return(user, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(tt.targetID, tt.targetRole)
			err := service.ResetPassword(context.Background(), tt.actorRole, schoolID, tt.targetID, tt.newPassword)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_RestoreUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name      string
		actorRole string
		userID    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Admin can restore user",
			actorRole: "admin",
			userID:    "user-123",
			mockSetup: func() {
				user := &User{ID: "user-123"}
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)
				mockRepo.On("Restore", mock.Anything, "user-123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Teacher cannot restore (no UserDelete permission)",
			actorRole: "teacher",
			userID:    "user-123",
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.RestoreUser(context.Background(), tt.actorRole, "", tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_DisableMFA(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name      string
		actorRole string
		userID    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Admin can disable MFA",
			actorRole: "admin",
			userID:    "user-123",
			mockSetup: func() {
				user := &User{ID: "user-123", MFAEnabled: true, MFASecret: "secret"}
				mockRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)
				mockRepo.On("ClearTempMFASecret", mock.Anything, "user-123").Return(nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *User) bool {
					return !u.MFAEnabled && u.MFASecret == ""
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Student cannot disable MFA (no permissions)",
			actorRole: "student",
			userID:    "user-123",
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.DisableMFA(context.Background(), tt.actorRole, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_IsGuardian(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name       string
		parentID   string
		studentID  string
		isGuardian bool
		wantErr    bool
	}{
		{
			name:       "Is guardian",
			parentID:   "parent-1",
			studentID:  "student-1",
			isGuardian: true,
			wantErr:    false,
		},
		{
			name:       "Is not guardian",
			parentID:   "parent-2",
			studentID:  "student-1",
			isGuardian: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("IsGuardian", mock.Anything, tt.parentID, tt.studentID).Return(tt.isGuardian, nil).Once()

			got, err := service.IsGuardian(context.Background(), tt.parentID, tt.studentID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.isGuardian, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetChildren(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name     string
		parentID string
		mockFn   func()
		wantErr  bool
	}{
		{
			name:     "Get children success",
			parentID: "parent-1",
			mockFn: func() {
				children := []StudentChild{{ID: "child-1", FirstName: "Child", LastName: "One"}}
				mockRepo.On("GetChildren", mock.Anything, "parent-1").Return(children, nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			_, err := service.GetChildren(context.Background(), "parent", tt.parentID, tt.parentID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_ExportUsers(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name      string
		actorRole string
		format    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Admin export csv",
			actorRole: "admin",
			format:    "csv",
			mockSetup: func() {
				users := []User{{ID: "u1", Email: "test@example.com", Role: "student"}}
				mockRepo.On("List", mock.Anything, mock.AnythingOfType("users.UserFilter")).Return(users, 1, nil).Once()
			},
			wantErr: false,
		},
		{
			name:      "Unauthorized export",
			actorRole: "student",
			format:    "csv",
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			data, err := service.ExportUsers(context.Background(), tt.actorRole, "", UserFilter{}, tt.format)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, data)
			}
		})
	}
}

func TestService_GDPRConvert(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("GDPR Export", func(t *testing.T) {
		user := &User{ID: "user-1", Email: "u@e.com", FirstName: "U", LastName: "S", Role: "student"}
		logs := []AuditLog{}

		mockRepo.On("GetByID", mock.Anything, "user-1").Return(user, nil).Once()
		mockRepo.On("GetAuditLogs", mock.Anything, "user-1", 1000, 0).Return(logs, 0, nil).Once()

		data, err := service.GDPRDataExport(context.Background(), "user-1", "student", "user-1")
		assert.NoError(t, err)
		assert.NotNil(t, data)
		profile := data["profile"].(map[string]interface{})
		assert.Equal(t, "u@e.com", profile["email"])
	})

	t.Run("GDPR Delete", func(t *testing.T) {
		user := &User{ID: "user-del", Email: "del@e.com", FirstName: "F", LastName: "L"}

		mockRepo.On("GetByID", mock.Anything, "user-del").Return(user, nil).Once()
		// Update should be called with modified user
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *User) bool {
			return u.Email != "del@e.com" // Basic check that it changed
		})).Return(nil).Once()
		mockRepo.On("RevokeAllUserTokens", mock.Anything, "user-del").Return(nil).Once()

		err := service.GDPRDelete(context.Background(), "admin", "user-del")
		assert.NoError(t, err)
	})
}

func TestService_DeleteUser_SchoolIDIsolation(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("Admin deleting user in another school -> Unauthorized", func(t *testing.T) {
		otherSchool := "school-B"
		user := &User{ID: "target-1", SchoolID: &otherSchool}
		mockRepo.On("GetByID", mock.Anything, "target-1").Return(user, nil).Once()

		err := service.DeleteUser(context.Background(), "admin", "school-A", "target-1")
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("SuperAdmin deleting user in another school -> Allowed", func(t *testing.T) {
		otherSchool := "school-B"
		user := &User{ID: "target-2", SchoolID: &otherSchool}
		mockRepo.On("GetByID", mock.Anything, "target-2").Return(user, nil).Once()
		mockRepo.On("Delete", mock.Anything, "target-2").Return(nil).Once()

		err := service.DeleteUser(context.Background(), "superadmin", "school-A", "target-2")
		assert.NoError(t, err)
	})
}

func TestService_RestoreUser_SchoolIDIsolation(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("Secretary restoring user in another school -> Unauthorized", func(t *testing.T) {
		otherSchool := "school-B"
		user := &User{ID: "target-1", SchoolID: &otherSchool}
		mockRepo.On("GetByID", mock.Anything, "target-1").Return(user, nil).Once()

		err := service.RestoreUser(context.Background(), "secretary", "school-A", "target-1")
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("Secretary restoring user in same school -> Allowed", func(t *testing.T) {
		sameSchool := "school-A"
		user := &User{ID: "target-2", SchoolID: &sameSchool}
		mockRepo.On("GetByID", mock.Anything, "target-2").Return(user, nil).Once()
		mockRepo.On("Restore", mock.Anything, "target-2").Return(nil).Once()

		err := service.RestoreUser(context.Background(), "secretary", "school-A", "target-2")
		assert.NoError(t, err)
	})
}
