package teacher_hours

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/teacher-hours")
	{
		group.GET("/dashboard", h.GetDashboard)
		group.POST("/logs", h.LogActivity)
		group.POST("/leaves", h.RequestShortLeave)
		group.PATCH("/leaves/:id/review", h.ReviewShortLeave)
		group.POST("/leaves/:id/recover", h.RecordLeaveRecovery)
		group.GET("/school-summary", h.ListSchoolSummary)
	}
}

func (h *Handler) GetDashboard(c *gin.Context) {
	schoolID := c.GetString("school_id")
	teacherID := c.GetString("user_id")
	if tID := c.Query("teacher_id"); tID != "" {
		role := c.GetString("user_role")
		if role == "admin" || role == "superadmin" || role == "principal" || role == "dsga" {
			teacherID = tID
		}
	}
	academicYear := c.DefaultQuery("academic_year", "2025/2026")

	obl, logs, leaves, err := h.svc.GetTeacherDashboard(schoolID, teacherID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"obligation": obl,
		"logs":       logs,
		"leaves":     leaves,
	})
}

func (h *Handler) LogActivity(c *gin.Context) {
	teacherID := c.GetString("user_id")
	var req LogHourActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.svc.LogHourActivity(teacherID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, log)
}

func (h *Handler) RequestShortLeave(c *gin.Context) {
	schoolID := c.GetString("school_id")
	teacherID := c.GetString("user_id")

	var req ShortLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leave, err := h.svc.RequestShortLeave(schoolID, teacherID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, leave)
}

func (h *Handler) ReviewShortLeave(c *gin.Context) {
	leaveID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ReviewShortLeave(leaveID, req.Status, req.Notes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "stato aggiornato con successo"})
}

func (h *Handler) RecordLeaveRecovery(c *gin.Context) {
	leaveID := c.Param("id")
	var req RecoverLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.RecordLeaveRecovery(leaveID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "recupero registrato"})
}

func (h *Handler) ListSchoolSummary(c *gin.Context) {
	schoolID := c.GetString("school_id")
	academicYear := c.DefaultQuery("academic_year", "2025/2026")

	obls, err := h.svc.ListSchoolObligations(schoolID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"obligations": obls})
}
