package admin

import "context"

// Repository defines the interface for admin data access
type Repository interface {
	// Dashboard & Stats
	CountSchools(ctx context.Context, schoolID *string) (int64, error)
	CountUsers(ctx context.Context, schoolID *string) (int64, error)
	CountUsersByRole(ctx context.Context, role string, schoolID *string) (int64, error)
	CountActiveUsers24h(ctx context.Context, schoolID *string) (int64, error)
	CountDocuments(ctx context.Context, schoolID *string) (int64, error)
	CountPendingDocuments(ctx context.Context, schoolID *string) (int64, error)
	CountCommunications(ctx context.Context, schoolID *string) (int64, error)
	GetRecentEvents(ctx context.Context, limit int, schoolID *string) ([]RecentEvent, error)
	GetSystemHealth(ctx context.Context) (*SystemHealthStatus, error)
	GetUserGrowth(ctx context.Context, schoolID *string) ([]UserGrowthPoint, error)

	// Schools
	ListSchools(ctx context.Context, req *SchoolListRequest, offset int, schoolID *string) ([]SchoolResponse, int64, error)
	GetSchool(ctx context.Context, schoolID string) (*SchoolResponse, error)
	CreateSchool(ctx context.Context, req *CreateSchoolRequest) (*SchoolResponse, error)
	UpdateSchool(ctx context.Context, schoolID string, req *UpdateSchoolRequest) (*SchoolResponse, error)
	DeleteSchool(ctx context.Context, schoolID string) error
	SchoolCodeExists(ctx context.Context, code string) (bool, error)

	// Admin Users
	ListAdminUsers(ctx context.Context, offset, limit int, schoolFilter *string) ([]AdminUserResponse, int64, error)
	GetAdminUserByID(ctx context.Context, adminID string) (*AdminUserResponse, error)
	CreateAdminUser(ctx context.Context, req *CreateAdminRequest) (*AdminUserResponse, error)
	UpdateAdminUser(ctx context.Context, adminID string, req *UpdateAdminRequest) (*AdminUserResponse, error)
	DeleteAdminUser(ctx context.Context, adminID string) error
	UserEmailExists(ctx context.Context, email string) (bool, error)

	// Activity Log
	GetAdminActivity(ctx context.Context, adminID string, limit int) ([]ActivityLogEntry, error)
	ListAuditLogs(ctx context.Context, req *AuditLogListRequest, offset int) ([]ActivityLogEntry, int64, error)
	LogAdminAction(ctx context.Context, adminID, actionType, target string, targetID *string, schoolID *string, details string) error

	// Settings
	GetSetting(ctx context.Context, schoolID, key string) (string, error)
	UpdateSetting(ctx context.Context, schoolID, key, value string) error

	// Data Integrity
	CheckDataIntegrity(ctx context.Context, schoolID *string) (*DataIntegrityReport, error)
}
