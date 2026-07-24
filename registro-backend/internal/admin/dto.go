package admin

import "time"

// SchoolListRequest represents the request to list schools
type SchoolListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search"`
	Status   string `form:"status"`
	SortBy   string `form:"sort_by"`
	SortDir  string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

// SchoolResponse represents a school in responses
type SchoolResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	Province     string    `json:"province"`
	ZipCode      string    `json:"zip_code"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Website      string    `json:"website"`
	StudentCount int       `json:"student_count"`
	TeacherCount int       `json:"teacher_count"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SchoolListResponse represents the response for school list
type SchoolListResponse struct {
	Items      []SchoolResponse `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// CreateSchoolRequest represents the request to create a school
type CreateSchoolRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=200"`
	Code     string `json:"code" binding:"required,min=2,max=50"`
	Address  string `json:"address" binding:"required"`
	City     string `json:"city" binding:"required"`
	Province string `json:"province" binding:"required,len=2"`
	ZipCode  string `json:"zip_code" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"omitempty,email"`
	Website  string `json:"website"`
}

// UpdateSchoolRequest represents the request to update a school
type UpdateSchoolRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2,max=200"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Province string `json:"province" binding:"omitempty,len=2"`
	ZipCode  string `json:"zip_code"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"omitempty,email"`
	Website  string `json:"website"`
	IsActive *bool  `json:"is_active"`
}

// AdminUserResponse represents an admin user in responses
type AdminUserResponse struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Role        string     `json:"role"`
	SchoolID    *string    `json:"school_id"`
	SchoolName  *string    `json:"school_name"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// AdminUserListResponse represents the response for admin user list
type AdminUserListResponse struct {
	Items      []AdminUserResponse `json:"items"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// CreateAdminRequest represents the request to create an admin user
type CreateAdminRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	FirstName string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string  `json:"last_name" binding:"required,min=2,max=100"`
	Password  string  `json:"password" binding:"required,min=8"`
	Role      string  `json:"role" binding:"required,oneof=admin superadmin"`
	SchoolID  *string `json:"school_id"` // Required if role is admin
}

// UpdateAdminRequest represents the request to update an admin user
type UpdateAdminRequest struct {
	FirstName string  `json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName  string  `json:"last_name" binding:"omitempty,min=2,max=100"`
	Password  *string `json:"password" binding:"omitempty,min=8"`
	SchoolID  *string `json:"school_id"`
	IsActive  *bool   `json:"is_active"`
}

// DashboardStatsResponse represents dashboard statistics
type DashboardStatsResponse struct {
	TotalSchools          int64               `json:"total_schools"`
	TotalUsers            int64               `json:"total_users"`
	TotalStudents         int64               `json:"total_students"`
	TotalTeachers         int64               `json:"total_teachers"`
	TotalDocuments        int64               `json:"total_documents"`
	PendingDocumentsCount int64               `json:"pending_documents_count"`
	AnnouncementsCount    int64               `json:"announcements_count"`
	ActiveUsers24h        int64               `json:"active_users_24h"`
	RecentEvents          []RecentEvent       `json:"recent_events"`
	HealthStatus          *SystemHealthStatus `json:"health_status,omitempty"` // Only for superadmin
}

// RecentEvent represents a recent system event
type RecentEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserName    string    `json:"user_name"`
	SchoolName  *string   `json:"school_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// SystemHealthStatus represents system health information
type SystemHealthStatus struct {
	Database      HealthCheck `json:"database"`
	Storage       HealthCheck `json:"storage"`
	API           HealthCheck `json:"api"`
	OverallStatus string      `json:"overall_status"` // healthy, degraded, down
}

// HealthCheck represents a single health check result
type HealthCheck struct {
	Status  string  `json:"status"` // healthy, warning, error
	Message string  `json:"message"`
	Value   *string `json:"value,omitempty"`
}

// ActivityLogEntry represents an admin activity log entry
type ActivityLogEntry struct {
	ID         string    `json:"id"`
	AdminID    string    `json:"admin_id"`
	AdminName  string    `json:"admin_name"`
	ActionType string    `json:"action_type"`
	Target     string    `json:"target"`
	TargetID   *string   `json:"target_id"`
	SchoolName *string   `json:"school_name"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"created_at"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// AuditLogListRequest represents request to list audit logs
type AuditLogListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	AdminID  string `form:"admin_id"`
	Action   string `form:"action"`
	FromDate string `form:"from_date"`
	ToDate   string `form:"to_date"`
}

// AuditLogListResponse represents response for audit logs
type AuditLogListResponse struct {
	Items      []ActivityLogEntry `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
