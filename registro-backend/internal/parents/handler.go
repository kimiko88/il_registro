package parents

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
	p := r.Group("/parents")
	{
		p.GET("/dashboard", h.GetDashboard)
		p.GET("/dashboard/stats", h.GetDashboardStats)
		p.GET("/child/:studentId/grades-average", h.GetChildGradesAverage)
	}
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stats, err := h.service.GetDashboardStats(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}


func (h *Handler) GetDashboard(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	dash, err := h.service.GetDashboard(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dash)
}

func (h *Handler) GetChildGradesAverage(c *gin.Context) {
	parentID := c.GetString("user_id")
	studentID := c.Param("studentId")

	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	avg, err := h.service.GetChildGradesAverage(c.Request.Context(), parentID, studentID)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"student_id": studentID,
		"average":    avg,
	})
}
