package psychology

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
	g := r.Group("/psychology")
	{
		g.POST("/consent", h.SignConsent)
		g.POST("/book", h.BookSession)
		g.GET("/sessions", h.ListMySessions)
		g.GET("/sessions/:id", h.GetSession)
		g.PUT("/sessions/:id/notes", h.UpdateClinicalNotes)
	}
}

func (h *Handler) SignConsent(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")

	var req SaveConsentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consent, err := h.service.SignParentConsent(c.Request.Context(), schoolID, req.StudentID, req.SchoolYear, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, consent)
}

func (h *Handler) BookSession(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	var req BookSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In Italian secondary schools students booking CIC are mostly minors unless specified
	isMinor := role == "student"

	session, err := h.service.BookSession(c.Request.Context(), schoolID, userID, "2025/2026", isMinor, req)
	if err != nil {
		if err == ErrBothParentsConsentRequired {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, session)
}

func (h *Handler) ListMySessions(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	sessions, err := h.service.ListMySessions(c.Request.Context(), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(sessions), "sessions": sessions})
}

func (h *Handler) GetSession(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	sess, err := h.service.GetSessionSecure(c.Request.Context(), id, userID, role)
	if err != nil {
		if err == ErrProfessionalSecrecyViolation {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sess)
}

func (h *Handler) UpdateClinicalNotes(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if role != "psychologist" {
		c.JSON(http.StatusForbidden, gin.H{"error": ErrProfessionalSecrecyViolation.Error()})
		return
	}

	var req struct {
		Notes string `json:"notes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sess, err := h.service.UpdateClinicalNotes(c.Request.Context(), id, userID, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sess)
}
