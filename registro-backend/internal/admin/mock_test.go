package admin

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CountSchools(ctx context.Context, schoolID *string) (int64, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CountUsers(ctx context.Context, schoolID *string) (int64, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CountUsersByRole(ctx context.Context, role string, schoolID *string) (int64, error) {
	args := m.Called(ctx, role, schoolID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CountActiveUsers24h(ctx context.Context, schoolID *string) (int64, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) GetRecentEvents(ctx context.Context, limit int, schoolID *string) ([]RecentEvent, error) {
	args := m.Called(ctx, limit, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]RecentEvent), args.Error(1)
}

func (m *MockRepository) GetSystemHealth(ctx context.Context) (*SystemHealthStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SystemHealthStatus), args.Error(1)
}

func (m *MockRepository) ListSchools(ctx context.Context, req *SchoolListRequest, offset int, schoolID *string) ([]SchoolResponse, int64, error) {
	args := m.Called(ctx, req, offset, schoolID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]SchoolResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetSchool(ctx context.Context, schoolID string) (*SchoolResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SchoolResponse), args.Error(1)
}

func (m *MockRepository) CreateSchool(ctx context.Context, req *CreateSchoolRequest) (*SchoolResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SchoolResponse), args.Error(1)
}

func (m *MockRepository) UpdateSchool(ctx context.Context, schoolID string, req *UpdateSchoolRequest) (*SchoolResponse, error) {
	args := m.Called(ctx, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SchoolResponse), args.Error(1)
}

func (m *MockRepository) DeleteSchool(ctx context.Context, schoolID string) error {
	args := m.Called(ctx, schoolID)
	return args.Error(0)
}

func (m *MockRepository) SchoolCodeExists(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) ListAdminUsers(ctx context.Context, offset, limit int, schoolFilter *string) ([]AdminUserResponse, int64, error) {
	args := m.Called(ctx, offset, limit, schoolFilter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]AdminUserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetAdminUserByID(ctx context.Context, adminID string) (*AdminUserResponse, error) {
	args := m.Called(ctx, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AdminUserResponse), args.Error(1)
}

func (m *MockRepository) CreateAdminUser(ctx context.Context, req *CreateAdminRequest) (*AdminUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AdminUserResponse), args.Error(1)
}

func (m *MockRepository) UpdateAdminUser(ctx context.Context, adminID string, req *UpdateAdminRequest) (*AdminUserResponse, error) {
	args := m.Called(ctx, adminID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AdminUserResponse), args.Error(1)
}

func (m *MockRepository) DeleteAdminUser(ctx context.Context, adminID string) error {
	args := m.Called(ctx, adminID)
	return args.Error(0)
}

func (m *MockRepository) UserEmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetAdminActivity(ctx context.Context, adminID string, limit int) ([]ActivityLogEntry, error) {
	args := m.Called(ctx, adminID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ActivityLogEntry), args.Error(1)
}

func (m *MockRepository) ListAuditLogs(ctx context.Context, req *AuditLogListRequest, offset int) ([]ActivityLogEntry, int64, error) {
	args := m.Called(ctx, req, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]ActivityLogEntry), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) LogAdminAction(ctx context.Context, adminID, actionType, target string, targetID *string, schoolID *string, details string) error {
	args := m.Called(ctx, adminID, actionType, target, targetID, schoolID, details)
	return args.Error(0)
}

func (m *MockRepository) CountCommunications(ctx context.Context, schoolID *string) (int64, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CountDocuments(ctx context.Context, schoolID *string) (int64, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CountPendingDocuments(ctx context.Context, schoolID *string) (int64, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) GetSetting(ctx context.Context, schoolID, key string) (string, error) {
	args := m.Called(ctx, schoolID, key)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) UpdateSetting(ctx context.Context, schoolID, key, value string) error {
	args := m.Called(ctx, schoolID, key, value)
	return args.Error(0)
}
