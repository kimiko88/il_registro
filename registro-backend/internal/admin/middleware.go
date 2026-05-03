package admin

import (
	"net/http"
	"registro-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

// Middleware provides admin permission checking
type Middleware struct{}

// NewMiddleware creates a new admin middleware
func NewMiddleware() *Middleware {
	return &Middleware{}
}

// RequireSuperAdmin ensures the user is a superadmin
func (m *Middleware) RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := auth.GetUserRole(c)
		if !exists || role != "superadmin" {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Message: "superadmin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAdminOrSuperAdmin ensures the user is an admin or superadmin
func (m *Middleware) RequireAdminOrSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := auth.GetUserRole(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "authentication required",
			})
			c.Abort()
			return
		}

		if role != "admin" && role != "superadmin" {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Message: "admin or superadmin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireStaff ensures the user is an admin, superadmin or secretary
func (m *Middleware) RequireStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := auth.GetUserRole(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "authentication required",
			})
			c.Abort()
			return
		}

		if role != "admin" && role != "superadmin" && role != "secretary" {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Message: "staff access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetSchoolFilter sets school filter in context for admin users
// SuperAdmin users can access all schools, so no filter is set
func (m *Middleware) SetSchoolFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := auth.GetUserRole(c)

		// Only filter for admin role, not for superadmin
		if role == "admin" {
			schoolID, exists := auth.GetSchoolID(c)
			if exists && schoolID != "" {
				c.Set("filter_school_id", schoolID)
			}
		}

		c.Next()
	}
}

// GetFilteredSchoolID retrieves the filtered school ID from context
// Returns empty string if no filter (superadmin case)
func GetFilteredSchoolID(c *gin.Context) string {
	if schoolID, exists := c.Get("filter_school_id"); exists {
		if id, ok := schoolID.(string); ok {
			return id
		}
	}
	return ""
}

// IsSuperAdmin checks if the current user is a superadmin
func IsSuperAdmin(c *gin.Context) bool {
	role, exists := auth.GetUserRole(c)
	return exists && role == "superadmin"
}

// CanAccessSchool checks if the user can access a specific school
func CanAccessSchool(c *gin.Context, schoolID string) bool {
	// Superadmin can access all schools
	if IsSuperAdmin(c) {
		return true
	}

	// Admin can only access their assigned school
	role, exists := auth.GetUserRole(c)
	if !exists || role != "admin" {
		return false
	}

	userSchoolID, exists := auth.GetSchoolID(c)
	if !exists {
		return false
	}

	return userSchoolID == schoolID
}
