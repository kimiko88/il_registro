package users

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func getActorRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	v, ok := role.(string)
	if !ok {
		return ""
	}
	return v
}

func getActorID(c *gin.Context) string {
	id, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	v, ok := id.(string)
	if !ok {
		return ""
	}
	return v
}

func getSchoolID(c *gin.Context) string {
	schoolID, exists := c.Get("school_id")
	if !exists {
		return ""
	}
	v, ok := schoolID.(string)
	if !ok {
		return ""
	}
	return v
}

// 1. POST /api/v1/users
func (h *Handler) Create(c *gin.Context) {
	if getActorID(c) == "" || getActorRole(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

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
		if err == ErrEmailExists || err == ErrFiscalCode || strings.Contains(err.Error(), "password") || strings.Contains(err.Error(), "fiscal code") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// 2. GET /api/v1/users
func (h *Handler) List(c *gin.Context) {
	if getActorID(c) == "" || getActorRole(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if pageSize > 500 {
		pageSize = 500
	}
	if pageSize < 1 {
		pageSize = 50
	}

	isActiveStr := c.Query("is_active")
	var isActive *bool
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	// include_deleted is restricted to admin/superadmin
	includeDeleted := false
	if c.Query("include_deleted") == "true" {
		role := getActorRole(c)
		if role != "admin" && role != "superadmin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: include_deleted requires admin role"})
			return
		}
		includeDeleted = true
	}

	schoolID := c.Query("school_id")
	var schoolIDPtr *string
	if schoolID != "" {
		schoolIDPtr = &schoolID
	}

	sortBy := strings.ToLower(c.Query("sort_by"))
	switch sortBy {
	case "id", "first_name", "last_name", "email", "created_at", "role", "school_id", "updated_at":
		// valid column name
	default:
		sortBy = "created_at"
	}
	sortOrder := strings.ToUpper(c.Query("sort_order"))
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	filter := UserFilter{
		Query:     c.Query("q"),
		Role:      c.Query("role"),
		SchoolID:  schoolIDPtr,
		IsActive:  isActive,
		IsDeleted: includeDeleted,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		ClassID:   c.Query("class_id"),
	}

	users, total, err := h.service.ListUsers(c.Request.Context(), getActorRole(c), getSchoolID(c), filter)
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
	if getActorID(c) == "" || getActorRole(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := h.service.GetUser(c.Request.Context(), getActorRole(c), getSchoolID(c), c.Param("id"))
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
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
	if getActorID(c) == "" || getActorRole(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.UpdateUser(c.Request.Context(), getActorRole(c), getSchoolID(c), c.Param("id"), req)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if err == ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if err == ErrEmailExists || err == ErrFiscalCode || (err.Error() != "" && strings.Contains(err.Error(), "fiscal code")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// 5. DELETE /api/v1/users/{id}
func (h *Handler) Delete(c *gin.Context) {
	if getActorID(c) == "" || getActorRole(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.DeleteUser(c.Request.Context(), getActorRole(c), getSchoolID(c), c.Param("id")); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if err == ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// POST /api/v1/users/bulk-delete
func (h *Handler) BulkDelete(c *gin.Context) {
	if getActorID(c) == "" || getActorRole(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	count, err := h.service.BulkDeleteUsers(c.Request.Context(), getActorRole(c), getSchoolID(c), req.UserIDs)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "users deleted", "count": count})
}

// 6. POST /api/v1/users/{id}/restore
func (h *Handler) Restore(c *gin.Context) {
	if getActorID(c) == "" || getActorRole(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.RestoreUser(c.Request.Context(), getActorRole(c), getSchoolID(c), c.Param("id")); err != nil {
		if errors.Is(err, ErrUnauthorized) || err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user restored"})
}

// 7. POST /api/v1/users/bulk-import
func (h *Handler) BulkImport(c *gin.Context) {
	if getActorID(c) == "" || getActorRole(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer func() { _ = file.Close() }()
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB limit"})
		return
	}
	res, err := h.service.BulkImport(c.Request.Context(), getActorRole(c), getSchoolID(c), file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 8. POST /api/v1/users/{id}/change-password
// Self-service only: a user can only change their own password.
// Admins must use ForceResetPassword instead.
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ERR_REQUIRED_FIELDS",
			"code":    "ERR_REQUIRED_FIELDS",
			"message": "Compilare tutti i campi obbligatori (la nuova password deve contenere almeno 10 caratteri)",
		})
		return
	}
	targetID := c.Param("id")
	actorID := getActorID(c)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "ERR_UNAUTHORIZED",
			"code":    "ERR_UNAUTHORIZED",
			"message": "Non autorizzato",
		})
		return
	}
	// Strictly self-service: no admin bypass allowed here.
	if targetID != actorID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "ERR_FORBIDDEN",
			"code":    "ERR_FORBIDDEN",
			"message": "Puoi modificare solo la tua password",
		})
		return
	}
	if err := h.service.ChangePassword(c.Request.Context(), targetID, req); err != nil {
		errCode := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errCode,
			"code":    errCode,
			"message": getPasswordErrorMessage(errCode),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Password aggiornata con successo",
		"code":    "SUCCESS",
	})
}

func getPasswordErrorMessage(code string) string {
	switch code {
	case "ERR_CURRENT_PASSWORD_INCORRECT":
		return "La password attuale non è corretta"
	case "ERR_PASSWORD_COMPLEXITY":
		return "La password deve contenere almeno una lettera maiuscola, una minuscola, un numero e un carattere speciale"
	case "ERR_PASSWORD_TOO_SHORT":
		return "La password deve contenere almeno 10 caratteri"
	case "ERR_PASSWORD_TOO_LONG":
		return "La password è troppo lunga (massimo 128 caratteri)"
	case "ERR_PASSWORD_RECENTLY_USED":
		return "La nuova password non può essere uguale a una delle ultime 5 password utilizzate"
	case "ERR_REQUIRED_FIELDS":
		return "Compilare tutti i campi obbligatori"
	default:
		return code
	}
}

// 9. POST /api/v1/users/{id}/reset-password - Force reset (admin, superadmin, secretary)
func (h *Handler) ForceResetPassword(c *gin.Context) {
	actorID := getActorID(c)
	role := getActorRole(c)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only privileged roles can force reset passwords"})
		return
	}
	targetID := c.Param("id")
	if role == "secretary" {
		targetUser, err := h.service.GetUser(c.Request.Context(), role, getSchoolID(c), targetID)
		if err == nil && targetUser != nil && (targetUser.Role == "admin" || targetUser.Role == "superadmin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: secretary cannot reset passwords of administrator accounts"})
			return
		}
	}
	var body struct {
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ResetPassword(c.Request.Context(), role, getSchoolID(c), targetID, body.NewPassword); err != nil {
		if errors.Is(err, ErrUnauthorized) || err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password reset"})
}

// 10. PATCH /api/v1/users/{id}/roles
func (h *Handler) AssignRoles(c *gin.Context) {
	actorID := getActorID(c)
	role := getActorRole(c)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only admins can assign roles"})
		return
	}
	var req AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.UpdateUser(c.Request.Context(), role, getSchoolID(c), c.Param("id"), UpdateUserRequest{
		Role: &req.Role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// 11. GET /api/v1/users/{id}/audit-log
// Restricted to: the user themselves, or admin/superadmin.
func (h *Handler) GetAuditLog(c *gin.Context) {
	actorID := getActorID(c)
	actorRole := getActorRole(c)
	targetID := c.Param("id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	logs, total, err := h.service.GetAuditLogs(c.Request.Context(), actorID, actorRole, targetID, pageSize, offset)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs, "total": total})
}

// 12. POST /api/v1/users/{id}/gdpr-export
func (h *Handler) ExportGDPR(c *gin.Context) {
	actorID := getActorID(c)
	role := getActorRole(c)
	targetID := c.Param("id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorID != targetID && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: can only export your own GDPR data or require admin role"})
		return
	}
	data, err := h.service.GDPRDataExport(c.Request.Context(), actorID, role, targetID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, ErrUserNotFound) || err == ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// 13. DELETE /api/v1/users/{id}/gdpr-delete
func (h *Handler) DeleteGDPR(c *gin.Context) {
	actorID := getActorID(c)
	role := getActorRole(c)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only admins can perform GDPR deletion"})
		return
	}
	if err := h.service.GDPRDelete(c.Request.Context(), role, c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user pseudonymized"})
}

// 15. PATCH /api/v1/users/{id}/disable-mfa
func (h *Handler) DisableMFA(c *gin.Context) {
	actorID := getActorID(c)
	role := getActorRole(c)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only admins can disable MFA"})
		return
	}
	if err := h.service.DisableMFA(c.Request.Context(), role, c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "mfa disabled"})
}

// 16. GET /api/v1/users/me/children
func (h *Handler) GetMyChildren(c *gin.Context) {
	actorID := getActorID(c)
	actorRole := getActorRole(c)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	children, err := h.service.GetChildren(c.Request.Context(), actorRole, actorID, actorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, children)
}

// 17. POST /api/v1/users/me/switch-child/:studentId
func (h *Handler) SwitchChild(c *gin.Context) {
	parentID := getActorID(c)
	studentID := c.Param("studentId")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	student, err := h.service.SwitchChildContext(c.Request.Context(), parentID, studentID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}
		if errors.Is(err, ErrUnauthorized) || strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "guardian") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"active_child": toUserResponse(student),
		"message":      "switched child context successfully",
	})
}

func toUserResponse(u *User) UserResponse {
	var dob *string
	if u.DateOfBirth != nil {
		s := u.DateOfBirth.Format("2006-01-02")
		dob = &s
	}
	return UserResponse{
		ID: u.ID, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName,
		FiscalCode: u.FiscalCode, Role: u.Role, SchoolID: u.SchoolID,
		ClassID: u.ClassID, ClassName: u.ClassName,
		IsActive: u.IsActive, EmailVerified: u.EmailVerified, MFAEnabled: u.MFAEnabled,
		PhoneNumber: u.PhoneNumber, JobTitle: u.JobTitle, DateOfBirth: dob,
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

func (h *Handler) GetGuardians(c *gin.Context) {
	if getActorID(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	studentID := c.Param("id")
	guardians, err := h.service.GetGuardians(c.Request.Context(), getActorRole(c), studentID)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, guardians)
}

func (h *Handler) AddGuardian(c *gin.Context) {
	if getActorID(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := getActorRole(c)
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only admins and secretary can add guardians"})
		return
	}
	studentID := c.Param("id")
	var req struct {
		ParentUserID     string `json:"parent_user_id" binding:"required"`
		RelationshipType string `json:"relationship_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AddGuardian(c.Request.Context(), role, studentID, req.ParentUserID, req.RelationshipType); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "guardian added"})
}

func (h *Handler) RemoveGuardian(c *gin.Context) {
	if getActorID(c) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := getActorRole(c)
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only admins and secretary can remove guardians"})
		return
	}
	studentID := c.Param("id")
	parentID := c.Param("guardianId")
	if err := h.service.RemoveGuardian(c.Request.Context(), role, studentID, parentID); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "guardian removed"})
}

func (h *Handler) GetFascicolo(c *gin.Context) {
	studentID := c.Param("id")
	actorID := getActorID(c)
	actorRole := getActorRole(c)

	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fascicolo, err := h.service.GetStudentFascicolo(c.Request.Context(), actorID, actorRole, studentID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fascicolo)
}
