package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Helper to get actor role from context (set by auth middleware)
func getActorRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return "" // Guest? Or should error.
	}
	return role.(string)
}

func getActorID(c *gin.Context) string {
	id, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return id.(string)
}

// 1. POST /api/v1/users
func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), getActorRole(c), req)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// 2. GET /api/v1/users - List paginated
// 14. GET /api/v1/users/search - Advanced search (merged logic)
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	isActiveStr := c.Query("is_active")
	var isActive *bool
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	schoolID := c.Query("school_id")
	var schoolIDPtr *string
	if schoolID != "" {
		schoolIDPtr = &schoolID
	}

	filter := UserFilter{
		Query:     c.Query("q"),
		Role:      c.Query("role"),
		SchoolID:  schoolIDPtr,
		IsActive:  isActive,
		IsDeleted: c.Query("include_deleted") == "true",
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	users, total, err := h.service.ListUsers(c.Request.Context(), getActorRole(c), filter)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ListUsersResponse{
		Users:      toUserResponses(users),
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
	})
}

// 3. GET /api/v1/users/{id}
func (h *Handler) Get(c *gin.Context) {
	user, err := h.service.GetUser(c.Request.Context(), getActorRole(c), c.Param("id"))
	if err != nil {
		if err == ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// 4. PATCH /api/v1/users/{id}
func (h *Handler) Update(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), getActorRole(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// 5. DELETE /api/v1/users/{id} - Soft delete
func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.DeleteUser(c.Request.Context(), getActorRole(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// 6. POST /api/v1/users/{id}/restore
func (h *Handler) Restore(c *gin.Context) {
	if err := h.service.RestoreUser(c.Request.Context(), getActorRole(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user restored"})
}

// 7. POST /api/v1/users/bulk-import
func (h *Handler) BulkImport(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	res, err := h.service.BulkImport(c.Request.Context(), getActorRole(c), file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 8. POST /api/v1/users/{id}/change-password
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure user serves themselves or is admin overriding?
	// Spec says /users/{id}/change-password. Usually self-service via /me/password
	// Here assume checking if ID matches token OR admin
	targetID := c.Param("id")
	actorID := getActorID(c)
	if targetID != actorID && getActorRole(c) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "can only change own password"})
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), targetID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// 9. POST /api/v1/users/{id}/reset-password - Force reset
func (h *Handler) ForceResetPassword(c *gin.Context) {
	var body struct {
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), getActorRole(c), c.Param("id"), body.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password reset"})
}

// 10. PATCH /api/v1/users/{id}/roles
func (h *Handler) AssignRoles(c *gin.Context) {
	var req AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Re-use UpdateUser logic but specifically for role
	user, err := h.service.UpdateUser(c.Request.Context(), getActorRole(c), c.Param("id"), UpdateUserRequest{
		Role: &req.Role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

// 11. GET /api/v1/users/{id}/audit-log
func (h *Handler) GetAuditLog(c *gin.Context) {
	logs, total, err := h.service.repo.GetAuditLogs(c.Request.Context(), c.Param("id"), 100, 0) // Limit hardcoded
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs, "total": total})
}

// 12. POST /api/v1/users/{id}/gdpr-export
func (h *Handler) ExportGDPR(c *gin.Context) {
	data, err := h.service.GDPRDataExport(c.Request.Context(), getActorRole(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// 13. DELETE /api/v1/users/{id}/gdpr-delete
func (h *Handler) DeleteGDPR(c *gin.Context) {
	if err := h.service.GDPRDelete(c.Request.Context(), getActorRole(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user pseudonymized"})
}

// 15. PATCH /api/v1/users/{id}/disable-mfa
func (h *Handler) DisableMFA(c *gin.Context) {
	if err := h.service.DisableMFA(c.Request.Context(), getActorRole(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "mfa disabled"})
}

// Helpers
func toUserResponse(u *User) UserResponse {
	return UserResponse{
		ID: u.ID, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName,
		FiscalCode: u.FiscalCode, Role: u.Role, SchoolID: u.SchoolID,
		IsActive: u.IsActive, EmailVerified: u.EmailVerified, MFAEnabled: u.MFAEnabled,
		PhoneNumber: u.PhoneNumber, JobTitle: u.JobTitle,
		CreatedAt: u.CreatedAt, LastLogin: u.LastLogin, DeletedAt: u.DeletedAt, PseudonymizedAt: u.PseudonymizedAt,
	}
}

func toUserResponses(users []User) []UserResponse {
	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = toUserResponse(&u)
	}
	return res
}
