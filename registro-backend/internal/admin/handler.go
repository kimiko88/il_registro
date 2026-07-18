package admin

import (
	"fmt"
	"net/http"
	"registro-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

// Handler provides HTTP handlers for admin endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new admin handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetDashboardStats returns dashboard statistics
// GET /api/v1/admin/dashboard/stats
func (h *Handler) GetDashboardStats(c *gin.Context) {
	role, _ := auth.GetUserRole(c)
	isSuperAdmin := role == "superadmin"

	var schoolID *string
	if !isSuperAdmin {
		userSchoolID, exists := auth.GetSchoolID(c)
		if exists {
			schoolID = &userSchoolID
		}
	}

	stats, err := h.service.GetDashboardStats(c.Request.Context(), isSuperAdmin, schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to get dashboard stats",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ListSchools returns a list of schools
// GET /api/v1/admin/schools
func (h *Handler) ListSchools(c *gin.Context) {
	var req SchoolListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	// Get school filter from middleware
	filterSchoolID := GetFilteredSchoolID(c)
	var schoolFilter *string
	if filterSchoolID != "" {
		schoolFilter = &filterSchoolID
	}

	schools, err := h.service.ListSchools(c.Request.Context(), &req, schoolFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to list schools",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schools)
}

// GetSchool returns a single school by ID
// GET /api/v1/admin/schools/:id
func (h *Handler) GetSchool(c *gin.Context) {
	schoolID := c.Param("id")

	// Check access permission
	if !CanAccessSchool(c, schoolID) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "you don't have access to this school",
		})
		return
	}

	school, err := h.service.GetSchool(c.Request.Context(), schoolID, nil)
	if err != nil {
		if err == ErrSchoolNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "school not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to get school",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, school)
}

// CreateSchool creates a new school
// POST /api/v1/admin/schools
func (h *Handler) CreateSchool(c *gin.Context) {
	var req CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	school, err := h.service.CreateSchool(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to create school",
			Message: err.Error(),
		})
		return
	}

	// Log action
	userID, _ := auth.GetUserID(c)
	_ = h.service.LogAdminAction(c.Request.Context(), userID, "create", "school", &school.ID, nil, "Created school: "+school.Name)

	c.JSON(http.StatusCreated, school)
}

// UpdateSchool updates an existing school
// PUT /api/v1/admin/schools/:id
func (h *Handler) UpdateSchool(c *gin.Context) {
	schoolID := c.Param("id")

	// Check access permission
	if !CanAccessSchool(c, schoolID) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "you don't have access to this school",
		})
		return
	}

	var req UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	school, err := h.service.UpdateSchool(c.Request.Context(), schoolID, &req, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to update school",
			Message: err.Error(),
		})
		return
	}

	// Log action
	userID, _ := auth.GetUserID(c)
	_ = h.service.LogAdminAction(c.Request.Context(), userID, "update", "school", &school.ID, &schoolID, "Updated school: "+school.Name)

	c.JSON(http.StatusOK, school)
}

// DeleteSchool deletes a school
// DELETE /api/v1/admin/schools/:id
func (h *Handler) DeleteSchool(c *gin.Context) {
	schoolID := c.Param("id")

	// Check access permission
	if !CanAccessSchool(c, schoolID) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "you don't have access to this school",
		})
		return
	}

	err := h.service.DeleteSchool(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to delete school",
			Message: err.Error(),
		})
		return
	}

	// Log action
	userID, _ := auth.GetUserID(c)
	_ = h.service.LogAdminAction(c.Request.Context(), userID, "delete", "school", &schoolID, &schoolID, "Deleted school")

	c.JSON(http.StatusOK, MessageResponse{Message: "school deleted successfully"})
}

// ListAdminUsers returns a list of admin users
// GET /api/v1/admin/users/admins
func (h *Handler) ListAdminUsers(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if _, err := fmt.Sscanf(p, "%d", &page); err != nil {
			page = 1
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			pageSize = 20
		}
	}

	// Get school filter from middleware/role context
	filterSchoolID := GetFilteredSchoolID(c)
	schoolFilter := c.Query("school_id")

	if filterSchoolID != "" {
		// Non-superadmin: force their own school_id
		schoolFilter = filterSchoolID
	}

	var schoolPtr *string
	if schoolFilter != "" {
		schoolPtr = &schoolFilter
	}

	admins, err := h.service.ListAdminUsers(c.Request.Context(), page, pageSize, schoolPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to list admin users",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, admins)
}

// CreateAdminUser creates a new admin user
// POST /api/v1/admin/users/admins
func (h *Handler) CreateAdminUser(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	admin, err := h.service.CreateAdminUser(c.Request.Context(), &req)
	if err != nil {
		if err == ErrAdminExists {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "admin user already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to create admin user",
			Message: err.Error(),
		})
		return
	}

	// Log action
	userID, _ := auth.GetUserID(c)
	_ = h.service.LogAdminAction(c.Request.Context(), userID, "create", "admin_user", &admin.ID, admin.SchoolID, "Created admin user: "+admin.Email)

	c.JSON(http.StatusCreated, admin)
}

// UpdateAdminUser updates an existing admin user
// PUT /api/v1/admin/users/admins/:id
// Only a superadmin may call this endpoint (enforced by RequireSuperAdmin middleware).
// We additionally verify that the target admin belongs to a school the caller
// can access, preventing cross-school privilege escalation.
func (h *Handler) UpdateAdminUser(c *gin.Context) {
	adminID := c.Param("id")

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	// Fetch target admin to verify ownership before mutating.
	target, err := h.service.repo.GetAdminUserByID(c.Request.Context(), adminID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "admin user not found"})
		return
	}
	if target.SchoolID != nil && !CanAccessSchool(c, *target.SchoolID) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "you don't have access to this admin's school",
		})
		return
	}

	admin, err := h.service.UpdateAdminUser(c.Request.Context(), adminID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to update admin user",
			Message: err.Error(),
		})
		return
	}

	// Log action
	userID, _ := auth.GetUserID(c)
	_ = h.service.LogAdminAction(c.Request.Context(), userID, "update", "admin_user", &admin.ID, admin.SchoolID, "Updated admin user: "+admin.Email)

	c.JSON(http.StatusOK, admin)
}

// DeleteAdminUser deletes an admin user
// DELETE /api/v1/admin/users/admins/:id
func (h *Handler) DeleteAdminUser(c *gin.Context) {
	adminID := c.Param("id")
	currentUserID, _ := auth.GetUserID(c)

	err := h.service.DeleteAdminUser(c.Request.Context(), adminID, currentUserID)
	if err != nil {
		if err == ErrCannotDeleteSelf {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "cannot delete your own account",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to delete admin user",
			Message: err.Error(),
		})
		return
	}

	// Log action
	_ = h.service.LogAdminAction(c.Request.Context(), currentUserID, "delete", "admin_user", &adminID, nil, "Deleted admin user")

	c.JSON(http.StatusOK, MessageResponse{Message: "admin user deleted successfully"})
}

// GetAdminActivity returns activity log for an admin user
// GET /api/v1/admin/users/admins/:id/activity
// Verifies that the requesting superadmin can access the target admin's school.
func (h *Handler) GetAdminActivity(c *gin.Context) {
	adminID := c.Param("id")
	limit := 50

	if l := c.Query("limit"); l != "" {
		if _, err := fmt.Sscanf(l, "%d", &limit); err != nil {
			limit = 50
		}
	}

	// Ownership check: verify caller can access the target admin's school.
	target, err := h.service.repo.GetAdminUserByID(c.Request.Context(), adminID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "admin user not found"})
		return
	}
	if target.SchoolID != nil && !CanAccessSchool(c, *target.SchoolID) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "you don't have access to this admin's school",
		})
		return
	}

	activity, err := h.service.GetAdminActivity(c.Request.Context(), adminID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to get admin activity",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": activity,
		"total": len(activity),
	})
}

// RegisterRoutes registers all admin routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup, middleware *Middleware) {
	adminGroup := router.Group("/admin")
	{
		// Dashboard stats (accessible to admin, superadmin, and secretary)
		adminGroup.GET("/dashboard/stats", middleware.RequireStaff(), middleware.SetSchoolFilter(), h.GetDashboardStats)
		adminGroup.GET("/settings/:key", middleware.RequireStaff(), middleware.SetSchoolFilter(), h.GetSchoolSetting)
		adminGroup.PUT("/settings/:key", middleware.RequireStaff(), middleware.SetSchoolFilter(), h.UpdateSchoolSetting)

		// Restricted admin routes (admin and superadmin only)
		restricted := adminGroup.Group("/")
		restricted.Use(middleware.RequireAdminOrSuperAdmin())
		{
			// Schools (with role-based access)
			schools := restricted.Group("/schools")
			schools.Use(middleware.SetSchoolFilter())
			{
				schools.GET("", h.ListSchools)
				schools.GET("/:id", h.GetSchool)
				schools.PUT("/:id", h.UpdateSchool)

				// Superadmin only
				schools.POST("", middleware.RequireSuperAdmin(), h.CreateSchool)
				schools.DELETE("/:id", middleware.RequireSuperAdmin(), h.DeleteSchool)
			}

			// Admin users (superadmin only)
			admins := restricted.Group("/users/admins")
			admins.Use(middleware.RequireSuperAdmin())
			{
				admins.GET("", h.ListAdminUsers)
				admins.POST("", h.CreateAdminUser)
				admins.PUT("/:id", h.UpdateAdminUser)
				admins.DELETE("/:id", h.DeleteAdminUser)
				admins.GET("/:id/activity", h.GetAdminActivity)
			}

			// Audit Logs (SuperAdmin only)
			restricted.GET("/audit-logs", middleware.RequireSuperAdmin(), h.ListAuditLogs)
		}
	}
}

// ListAuditLogs returns global audit logs
// GET /api/v1/admin/audit-logs
func (h *Handler) ListAuditLogs(c *gin.Context) {
	var req AuditLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	logs, err := h.service.ListAuditLogs(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to list audit logs",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetSchoolSetting returns a school setting
// GET /api/v1/admin/settings/:key
func (h *Handler) GetSchoolSetting(c *gin.Context) {
	key := c.Param("key")
	schoolID := GetFilteredSchoolID(c)

	if schoolID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "school_id required"})
		return
	}

	value, err := h.service.GetSchoolSetting(c.Request.Context(), schoolID, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

// UpdateSchoolSetting updates a school setting
// PUT /api/v1/admin/settings/:key
func (h *Handler) UpdateSchoolSetting(c *gin.Context) {
	key := c.Param("key")
	schoolID := GetFilteredSchoolID(c)

	if schoolID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "school_id required"})
		return
	}

	var body struct {
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateSchoolSetting(c.Request.Context(), schoolID, key, body.Value); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "setting updated successfully"})
}
