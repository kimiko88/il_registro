package tenants

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	t := r.Group("/tenants")
	{
		t.POST("", h.CreateTenant)
		t.GET("", h.ListTenants)
		t.GET("/:id", h.GetTenant)
	}
}

func (h *Handler) CreateTenant(c *gin.Context) {
	role := c.GetString("role")
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.service.CreateTenant(c.Request.Context(), role, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tenant)
}

func (h *Handler) ListTenants(c *gin.Context) {
	role := c.GetString("role")
	tenants, err := h.service.ListTenants(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tenants)
}

func (h *Handler) GetTenant(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "superadmin" && role != "admin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	if role != "superadmin" && schoolID != "" && schoolID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: cannot view other tenant details"})
		return
	}
	tenant, err := h.service.GetTenant(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tenant)
}
