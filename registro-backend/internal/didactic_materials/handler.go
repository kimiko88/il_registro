package didactic_materials

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	materials := rg.Group("/didactic-materials")
	{
		materials.GET("/class/:class_id", h.GetByClass)
		materials.POST("", h.Create)
		materials.DELETE("/:id", h.Delete)
	}
}

// GetByClass requires the caller to be authenticated.
func (h *Handler) GetByClass(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	classID := c.Param("class_id")
	res, err := h.service.GetMaterialsByClass(c.Request.Context(), userID, role, classID)
	if err != nil {
		errStr := err.Error()
		if strings.HasPrefix(errStr, "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": errStr})
			return
		}
		if strings.Contains(errStr, "not found") || errStr == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": errStr})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errStr})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Create(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only teachers or admins can create didactic materials"})
		return
	}

	schoolID := c.GetString("school_id")
	var req CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.CreateMaterial(c.Request.Context(), teacherID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only teachers or admins can delete didactic materials"})
		return
	}
	if err := h.service.DeleteMaterial(c.Request.Context(), id, teacherID); err != nil {
		errStr := err.Error()
		if strings.HasPrefix(errStr, "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": errStr})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errStr})
		return
	}
	c.Status(http.StatusNoContent)
}
