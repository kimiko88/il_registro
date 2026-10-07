package privacy

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
	g := r.Group("/privacy")
	{
		g.POST("/consent", h.SaveConsent)
		g.GET("/consent/:student_id", h.GetConsent)
		g.POST("/class-badges", h.GetClassBadges)
		g.POST("/treatments", h.CreateTreatment)
		g.GET("/treatments", h.ListTreatments)
	}
}

func (h *Handler) SaveConsent(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")

	var req SaveConsentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consent, err := h.service.SaveStudentConsent(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, consent)
}

func (h *Handler) GetConsent(c *gin.Context) {
	studentID := c.Param("student_id")
	year := c.DefaultQuery("school_year", "2025/2026")
	consent, err := h.service.GetStudentConsent(c.Request.Context(), studentID, year)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, consent)
}

func (h *Handler) GetClassBadges(c *gin.Context) {
	var req struct {
		StudentIDs []string `json:"student_ids" binding:"required"`
		SchoolYear string   `json:"school_year"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.SchoolYear == "" {
		req.SchoolYear = "2025/2026"
	}

	badges, err := h.service.GetClassBadges(c.Request.Context(), req.StudentIDs, req.SchoolYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"badges": badges})
}

func (h *Handler) CreateTreatment(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req TreatmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := h.service.CreateTreatment(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *Handler) ListTreatments(c *gin.Context) {
	schoolID := c.GetString("school_id")
	treatments, err := h.service.ListTreatments(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(treatments), "treatments": treatments})
}
