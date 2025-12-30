package attendance

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

// RegisterRoutes sets up the routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	att := r.Group("/attendance")
	// att.Use(middleware.AuthMiddleware()) // Assumed global or higher level

	// Teacher
	att.POST("/mark", h.MarkAttendance)
	att.POST("/mark-bulk", h.MarkBulk)
	att.GET("/class/:id", h.GetClassAttendance)
	att.GET("/pending-justifications", h.GetPendingJustifications)
	att.POST("/justification/:id/process", h.ProcessJustification)

	// Student/Parent
	att.GET("/my-attendance", h.GetMyAttendance)
	att.GET("/my-attendance/summary", h.GetMySummary)
	att.POST("/justify", h.RequestJustification)

	att.GET("/child-attendance/:studentID", h.GetChildAttendance)
	att.GET("/child-attendance/:studentID/summary", h.GetChildSummary)

	// Admin
	att.POST("/justification/:id/approve", h.ApproveJustification)
	att.DELETE("/justification/:id", h.RejectJustification)
	att.GET("/analytics", h.GetAnalytics)
}

func (h *Handler) GetMySummary(c *gin.Context) {
	studentID := c.GetString("userID")
	res, err := h.service.GetStudentSummary(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetChildAttendance(c *gin.Context) {
	studentID := c.Param("studentID")
	// Verify guardianship here or in service!
	res, err := h.service.GetStudentAttendance(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetChildSummary(c *gin.Context) {
	studentID := c.Param("studentID")
	res, err := h.service.GetStudentSummary(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) ApproveJustification(c *gin.Context) {
	id := c.Param("id")
	teacherID := c.GetString("userID")
	if err := h.service.ProcessJustification(c.Request.Context(), teacherID, id, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *Handler) RejectJustification(c *gin.Context) {
	id := c.Param("id")
	// Using DELETE as Reject per requirements
	if err := h.service.DeleteJustification(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	res, err := h.service.GetSchoolAnalytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) MarkAttendance(c *gin.Context) {
	var req CreateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID") // From middleware
	if err := h.service.MarkAttendance(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked"})
}

func (h *Handler) MarkBulk(c *gin.Context) {
	var req BulkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	if err := h.service.MarkBulk(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bulk marked"})
}

func (h *Handler) GetClassAttendance(c *gin.Context) {
	classID := c.Param("id")
	date := c.Query("date") // YYYY-MM-DD
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date required"})
		return
	}
	res, err := h.service.GetClassAttendance(c.Request.Context(), classID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetMyAttendance(c *gin.Context) {
	// studentID from token
	studentID := c.GetString("userID")
	res, err := h.service.GetStudentAttendance(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) RequestJustification(c *gin.Context) {
	var req JustificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Verify logical parent ownership? Service handles logic
	parentID := c.GetString("userID")
	if err := h.service.RequestJustification(c.Request.Context(), parentID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "requested"})
}

func (h *Handler) GetPendingJustifications(c *gin.Context) {
	classID := c.Query("class_id") // Teacher filters by class
	res, err := h.service.GetPendingJustifications(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) ProcessJustification(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Approve bool `json:"approve"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	teacherID := c.GetString("userID")
	if err := h.service.ProcessJustification(c.Request.Context(), teacherID, id, req.Approve); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "processed"})
}

// Needed for compatibility with main
func (h *Handler) GetAttendance(c *gin.Context) { h.GetMyAttendance(c) }
