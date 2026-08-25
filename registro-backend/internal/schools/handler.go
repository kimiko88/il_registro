package schools

import (
	"net/http"
	"strings"

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
	if role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only superadmin can create schools"})
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
	schoolID := c.GetString("school_id")
	if role != "superadmin" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only administrative staff can update schools"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school ID is required"})
		return
	}
	if role == "admin" && schoolID != "" && schoolID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: cannot update other schools"})
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

// ListPublic returns a list of public school contacts for unauthenticated users (e.g. login contact secretary dialog)
// GET /public/schools
func (h *Handler) ListPublic(c *gin.Context) {
	params := ListParams{Page: 1, PageSize: 100}
	schools, _, err := h.service.List(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type PublicSchoolInfo struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Code    string `json:"code"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		City    string `json:"city"`
		Address string `json:"address"`
	}

	var list []PublicSchoolInfo
	for _, s := range schools {
		if s == nil {
			continue
		}
		email := s.Email
		if email == "" {
			email = "segreteria@" + strings.ToLower(strings.ReplaceAll(s.Name, " ", "")) + ".it"
		}
		list = append(list, PublicSchoolInfo{
			ID:      s.ID,
			Name:    s.Name,
			Code:    s.Code,
			Email:   email,
			Phone:   s.Phone,
			City:    s.City,
			Address: s.Address,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": list})
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
