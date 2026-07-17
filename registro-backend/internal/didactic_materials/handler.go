package didactic_materials

import (
	"net/http"

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

func (h *Handler) GetByClass(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	classID := c.Param("class_id")
	res, err := h.service.GetMaterialsByClass(c.Request.Context(), userID, role, classID)
	if err != nil {
		if err.Error() == "unauthorized: student does not belong to this class" || err.Error() == "unauthorized: parent does not have any children in this class" || err.Error() == "unauthorized: invalid role" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Create(c *gin.Context) {
	teacherID := c.GetString("user_id")
	schoolID := c.GetString("school_id")

	var req CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateMaterial(teacherID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	teacherID := c.GetString("user_id")

	if err := h.service.DeleteMaterial(id, teacherID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
