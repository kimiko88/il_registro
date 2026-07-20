package schoolsettings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	settings := rg.Group("/school-settings")
	{
		settings.GET("", h.GetSettings)
		settings.PUT("", h.UpdateSettings)
	}
}

func (h *Handler) GetSettings(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id not found in token context"})
		return
	}

	res, err := h.service.GetSettings(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateSettings(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id not found in token context"})
		return
	}

	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only admin or principal can update school settings"})
		return
	}

	var req UpdateSchoolSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "message": err.Error()})
		return
	}

	res, err := h.service.UpdateSettings(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
