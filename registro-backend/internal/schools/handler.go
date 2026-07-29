package schools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for schools
type Handler struct {
	service *Service
}

// NewHandler creates a new schools handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles school creation
// POST /schools
func (h *Handler) Create(c *gin.Context) {
	role := c.GetString("role")
	if role != "superadmin" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only administrative staff can create schools"})
		return
	}

	var req CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	school, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, school)
}

// Get handles retrieving a school by ID
// GET /schools/:id
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school ID is required"})
		return
	}
	school, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if school == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}
	c.JSON(http.StatusOK, school)
}

// List handles listing schools
// GET /schools
func (h *Handler) List(c *gin.Context) {
	var params ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schools, total, err := h.service.List(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": schools,
		"total": total,
		"page":  params.Page,
	})
}

// Update handles updating a school
// PATCH /schools/:id
func (h *Handler) Update(c *gin.Context) {
	role := c.GetString("role")
	if role != "superadmin" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only administrative staff can update schools"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school ID is required"})
		return
	}
	var req UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "school updated"})
}

// Delete handles deleting a school
// DELETE /schools/:id
func (h *Handler) Delete(c *gin.Context) {
	role := c.GetString("role")
	if role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only superadmin can delete schools"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school ID is required"})
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "school deleted"})
}

// RegisterRoutes registers all school routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	schools := router.Group("/schools")
	{
		schools.POST("/", h.Create)
		schools.GET("/", h.List)
		schools.GET("/:id", h.Get)
		schools.PATCH("/:id", h.Update)
		schools.DELETE("/:id", h.Delete)
	}
}
