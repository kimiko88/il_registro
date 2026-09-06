package admin

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"registro-backend/internal/auth"
	"registro-backend/pkg/logger"
)

// Service provides admin business logic
type Service struct {
	repo Repository
}

// NewService creates a new admin service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetDashboardStats retrieves dashboard statistics
func (s *Service) GetDashboardStats(ctx context.Context, isSuperAdmin bool, schoolID *string) (*DashboardStatsResponse, error) {
	stats := &DashboardStatsResponse{}

	// Get school count
	if isSuperAdmin {
		count, err := s.repo.CountSchools(ctx, nil)
		if err != nil {
			return nil, err
		}
		stats.TotalSchools = count
	} else {
		stats.TotalSchools = 1 // Admin sees only their school
	}

	// Get user counts
	userCount, err := s.repo.CountUsers(ctx, schoolID)
	if err != nil {
		return nil, err
	}
	stats.TotalUsers = userCount

	studentCount, err := s.repo.CountUsersByRole(ctx, "student", schoolID)
	if err != nil {
		return nil, err
	}
	stats.TotalStudents = studentCount

	teacherCount, err := s.repo.CountUsersByRole(ctx, "teacher", schoolID)
	if err != nil {
		return nil, err
	}
	stats.TotalTeachers = teacherCount

	// Get document counts
	docCount, err := s.repo.CountDocuments(ctx, schoolID)
	if err != nil {
		return nil, err
	}
	stats.TotalDocuments = docCount

	pendingDocCount, err := s.repo.CountPendingDocuments(ctx, schoolID)
	if err != nil {
		return nil, err
	}
	stats.PendingDocumentsCount = pendingDocCount

	// Get communications count
	commCount, err := s.repo.CountCommunications(ctx, schoolID)
	if err != nil {
		return nil, err
	}
	stats.AnnouncementsCount = commCount

	// Get active users in last 24h
	activeCount, err := s.repo.CountActiveUsers24h(ctx, schoolID)
	if err != nil {
		return nil, err
	}
	stats.ActiveUsers24h = activeCount

	// Get recent events
	events, err := s.repo.GetRecentEvents(ctx, 10, schoolID)
	if err != nil {
		return nil, err
	}
	stats.RecentEvents = events

	// Get health status (superadmin only)
	if isSuperAdmin {
		health, err := s.repo.GetSystemHealth(ctx)
		if err != nil {
			logger.Log.Errorf("GetDashboardStats health check error: %v", err)
			// Don't fail the whole request if health check fails
			health = &SystemHealthStatus{
				OverallStatus: "unknown",
			}
		}
		stats.HealthStatus = health
	}

	return stats, nil
}

// ListSchools retrieves a paginated list of schools
func (s *Service) ListSchools(ctx context.Context, req *SchoolListRequest, schoolID *string) (*SchoolListResponse, error) {
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Get schools with filtering
	schools, total, err := s.repo.ListSchools(ctx, req, offset, schoolID)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &SchoolListResponse{
		Items:      schools,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetSchool retrieves a single school by ID
func (s *Service) GetSchool(ctx context.Context, schoolID string, userSchoolID *string) (*SchoolResponse, error) {
	// Check access if user is admin
	if userSchoolID != nil && *userSchoolID != schoolID {
		return nil, ErrUnauthorized
	}

	school, err := s.repo.GetSchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	return school, nil
}

// CreateSchool creates a new school (superadmin only)
func (s *Service) CreateSchool(ctx context.Context, req *CreateSchoolRequest) (*SchoolResponse, error) {
	// Check if school code already exists
	exists, err := s.repo.SchoolCodeExists(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("school with code %s already exists", req.Code)
	}

	school, err := s.repo.CreateSchool(ctx, req)
	if err != nil {
		return nil, err
	}

	return school, nil
}

// UpdateSchool updates an existing school
func (s *Service) UpdateSchool(ctx context.Context, schoolID string, req *UpdateSchoolRequest, userSchoolID *string) (*SchoolResponse, error) {
	// Check access if user is admin
	if userSchoolID != nil && *userSchoolID != schoolID {
		return nil, ErrUnauthorized
	}

	school, err := s.repo.UpdateSchool(ctx, schoolID, req)
	if err != nil {
		return nil, err
	}

	return school, nil
}

// DeleteSchool deletes a school (superadmin only)
func (s *Service) DeleteSchool(ctx context.Context, callerRole, schoolID string) error {
	if callerRole != "superadmin" {
		return ErrUnauthorized
	}
	return s.repo.DeleteSchool(ctx, schoolID)
}

// ListAdminUsers retrieves a paginated list of admin users (superadmin only)
func (s *Service) ListAdminUsers(ctx context.Context, page, pageSize int, schoolFilter *string) (*AdminUserListResponse, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	admins, total, err := s.repo.ListAdminUsers(ctx, offset, pageSize, schoolFilter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &AdminUserListResponse{
		Items:      admins,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAdminUserByID retrieves a single admin user by ID.
// Used by handlers for ownership checks before mutating admin records.
func (s *Service) GetAdminUserByID(ctx context.Context, adminID string) (*AdminUserResponse, error) {
	return s.repo.GetAdminUserByID(ctx, adminID)
}

// CreateAdminUser creates a new admin user (superadmin only).
// Validates password strength and hashes password before persisting.
func (s *Service) CreateAdminUser(ctx context.Context, req *CreateAdminRequest) (*AdminUserResponse, error) {
	// Validate that admin role has school_id
	if req.Role == "admin" && (req.SchoolID == nil || *req.SchoolID == "") {
		return nil, fmt.Errorf("admin role requires school_id")
	}

	// Validate that superadmin role doesn't have school_id
	if req.Role == "superadmin" && req.SchoolID != nil {
		return nil, fmt.Errorf("superadmin role cannot have school_id")
	}

	// Validate password strength (must meet the same rules as public registration)
	pv := auth.NewPasswordValidator()
	if err := pv.Validate(req.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := s.repo.UserEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrAdminExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	reqCopy := *req
	reqCopy.Password = string(hash)

	admin, err := s.repo.CreateAdminUser(ctx, &reqCopy)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// UpdateAdminUser updates an existing admin user (superadmin only)
func (s *Service) UpdateAdminUser(ctx context.Context, callerRole, adminID string, req *UpdateAdminRequest) (*AdminUserResponse, error) {
	if callerRole != "superadmin" {
		return nil, errors.New("unauthorized: only superadmin can update admin users")
	}

	reqCopy := *req
	if req.Password != nil && *req.Password != "" {
		pv := auth.NewPasswordValidator()
		if err := pv.Validate(*req.Password); err != nil {
			return nil, err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("password hashing failed: %w", err)
		}
		hashedStr := string(hash)
		reqCopy.Password = &hashedStr
	}

	admin, err := s.repo.UpdateAdminUser(ctx, adminID, &reqCopy)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// DeleteAdminUser deletes an admin user (superadmin only)
func (s *Service) DeleteAdminUser(ctx context.Context, callerRole, adminID, currentUserID string) error {
	if callerRole != "superadmin" {
		return errors.New("unauthorized: only superadmin can delete admin users")
	}

	// Prevent self-deletion
	if adminID == currentUserID {
		return ErrCannotDeleteSelf
	}

	target, err := s.repo.GetAdminUserByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}
	if target.Role != "admin" && target.Role != "superadmin" {
		return errors.New("unauthorized: target user is not an admin user")
	}

	return s.repo.DeleteAdminUser(ctx, adminID)
}

// GetSystemHealth retrieves overall system health from repository
func (s *Service) GetSystemHealth(ctx context.Context) (*SystemHealthStatus, error) {
	return s.repo.GetSystemHealth(ctx)
}

// GetAdminActivity retrieves activity log for an admin user
func (s *Service) GetAdminActivity(ctx context.Context, adminID string, limit int) ([]ActivityLogEntry, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}

	return s.repo.GetAdminActivity(ctx, adminID, limit)
}

// ListAuditLogs retrieves paginated audit logs for superadmin
func (s *Service) ListAuditLogs(ctx context.Context, req *AuditLogListRequest) (*AuditLogListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	logs, total, err := s.repo.ListAuditLogs(ctx, req, offset)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &AuditLogListResponse{
		Items:      logs,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// LogAdminAction logs an admin action for audit trail
func (s *Service) LogAdminAction(ctx context.Context, adminID, actionType, target string, targetID *string, schoolID *string, details string) error {
	return s.repo.LogAdminAction(ctx, adminID, actionType, target, targetID, schoolID, details)
}

// UpdateSetting updates a school setting
func (s *Service) UpdateSetting(ctx context.Context, schoolID, key, value string) error {
	return s.repo.UpdateSetting(ctx, schoolID, key, value)
}

// CheckDataIntegrity runs diagnostic checks across school entities
func (s *Service) CheckDataIntegrity(ctx context.Context, schoolID *string) (*DataIntegrityReport, error) {
	return s.repo.CheckDataIntegrity(ctx, schoolID)
}

// GetSchoolSetting retrieves a school setting
func (s *Service) GetSchoolSetting(ctx context.Context, schoolID, key string) (string, error) {
	return s.repo.GetSetting(ctx, schoolID, key)
}

var allowedSettingKeys = map[string]bool{
	"grading_scale":                   true,
	"semester_count":                  true,
	"language":                        true,
	"attendance_threshold":            true,
	"lock_scrutiny":                   true,
	"require_principal_approval":      true,
	"allow_parents_view_grades":       true,
	"require_mfa":                     true,
	"enable_substitute_notifications": true,
	"enable_pcto":                     true,
	"timetable_visible_to_parents":    true,
	"enable_elearning":                true,
	"enable_uda":                      true,
	"enable_colloqui":                 true,
}

// UpdateSchoolSetting updates a school setting with allowlist validation
func (s *Service) UpdateSchoolSetting(ctx context.Context, schoolID, key, value string) error {
	key = strings.ToLower(key)
	if !allowedSettingKeys[key] {
		return errors.New("invalid or unauthorized setting key")
	}
	return s.repo.UpdateSetting(ctx, schoolID, key, value)
}

// GetUserGrowth retrieves user growth data points
func (s *Service) GetUserGrowth(ctx context.Context, schoolID *string) ([]UserGrowthPoint, error) {
	return s.repo.GetUserGrowth(ctx, schoolID)
}
