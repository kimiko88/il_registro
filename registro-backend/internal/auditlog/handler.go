package auditlog

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	group := r.Group("/audit-log")
	{
		group.GET("", h.List)
		group.GET("/export", h.ExportCSV)
	}
}

func (h *Handler) List(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access restricted to admin and superadmin"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	params := FilterParams{
		SchoolID:   c.GetString("school_id"),
		ActorID:    c.Query("actor_id"),
		Action:     c.Query("action"),
		EntityType: c.Query("entity_type"),
		From:       c.Query("from"),
		To:         c.Query("to"),
		Page:       page,
		Limit:      limit,
	}

	result, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) ExportCSV(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access restricted to admin and superadmin"})
		return
	}

	params := FilterParams{
		SchoolID:   c.GetString("school_id"),
		ActorID:    c.Query("actor_id"),
		Action:     c.Query("action"),
		EntityType: c.Query("entity_type"),
		From:       c.Query("from"),
		To:         c.Query("to"),
	}

	csvBytes, err := h.service.ExportCSV(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"audit_logs_%s.csv\"", c.GetString("school_id")))
	c.Data(http.StatusOK, "text/csv", csvBytes)
}
