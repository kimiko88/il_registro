package admin

import (
	"context"
	"fmt"
	"math"
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
func (s *Service) DeleteSchool(ctx context.Context, schoolID string) error {
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

// CreateAdminUser creates a new admin user (superadmin only)
func (s *Service) CreateAdminUser(ctx context.Context, req *CreateAdminRequest) (*AdminUserResponse, error) {
	// Validate that admin role has school_id
	if req.Role == "admin" && (req.SchoolID == nil || *req.SchoolID == "") {
		return nil, fmt.Errorf("admin role requires school_id")
	}

	// Validate that superadmin role doesn't have school_id
	if req.Role == "superadmin" && req.SchoolID != nil {
		return nil, fmt.Errorf("superadmin role cannot have school_id")
	}

	// Check if email already exists
	exists, err := s.repo.UserEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrAdminExists
	}

	admin, err := s.repo.CreateAdminUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// UpdateAdminUser updates an existing admin user (superadmin only)
func (s *Service) UpdateAdminUser(ctx context.Context, adminID string, req *UpdateAdminRequest) (*AdminUserResponse, error) {
	admin, err := s.repo.UpdateAdminUser(ctx, adminID, req)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// DeleteAdminUser deletes an admin user (superadmin only)
func (s *Service) DeleteAdminUser(ctx context.Context, adminID, currentUserID string) error {
	// Prevent self-deletion
	if adminID == currentUserID {
		return ErrCannotDeleteSelf
	}

	return s.repo.DeleteAdminUser(ctx, adminID)
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
