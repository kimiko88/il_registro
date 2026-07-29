package timetables

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(repo Repository) *Handler {
	return &Handler{service: NewService(repo)}
}

func NewHandlerWithService(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetByClass(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	classID := c.Param("id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class ID is required"})
		return
	}

	schedule, err := h.service.GetByClass(c.Request.Context(), userID, role, schoolID, classID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	classID := c.Param("id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class ID is required"})
		return
	}

	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Update(c.Request.Context(), userID, role, schoolID, classID, req.Entries)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") || strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule updated successfully"})
}

func (h *Handler) GetMySchedule(c *gin.Context) {
	uid := c.GetString("user_id")
	role := c.GetString("role")
	schedule, err := h.service.GetMySchedule(c.Request.Context(), uid, role)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/classes/:id/schedule", h.GetByClass)
	r.POST("/classes/:id/schedule", h.Update)
	r.PUT("/classes/:id/schedule", h.Update)
	r.PATCH("/classes/:id/schedule", h.Update)
	r.GET("/timetables/my-schedule", h.GetMySchedule)
}
