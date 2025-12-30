package orientamento

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

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	grp := r.Group("/orientamento")

	// Teacher
	grp.POST("/events", h.CreateEvent)
	grp.GET("/events", h.GetEvents)
	grp.POST("/events/:id/attendance", h.MarkAttendance)

	// Student
	grp.POST("/register", h.RegisterStudent)
	grp.GET("/my-events", h.GetMyEvents)
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	if err := h.service.CreateEvent(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *Handler) GetEvents(c *gin.Context) {
	res, err := h.service.GetEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) RegisterStudent(c *gin.Context) {
	var req struct {
		EventID string `json:"event_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	if err := h.service.RegisterStudent(c.Request.Context(), userID, req.EventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registered"})
}

func (h *Handler) GetMyEvents(c *gin.Context) {
	userID := c.GetString("userID")
	res, err := h.service.GetMyEvents(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) MarkAttendance(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		StudentID string `json:"student_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.MarkAttendance(c.Request.Context(), id, req.StudentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "attendance marked"})
}
