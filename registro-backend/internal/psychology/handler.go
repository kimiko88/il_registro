package psychology

import (
	"fmt"
	"net/http"
	"time"

	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service  *Service
	userRepo users.Repository
}

func NewHandler(s *Service, repos ...users.Repository) *Handler {
	h := &Handler{service: s}
	if len(repos) > 0 {
		h.userRepo = repos[0]
	}
	return h
}

func currentSchoolYear(now time.Time) string {
	startYear := now.Year()
	if now.Month() < time.September {
		startYear--
	}
	return fmt.Sprintf("%d/%d", startYear, startYear+1)
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

	if c.GetString("role") != "parent" || h.userRepo == nil || schoolID == "" || userID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo un genitore verificato può firmare il consenso"})
		return
	}
	parent, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil || parent == nil || parent.SchoolID == nil || *parent.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "genitore non associato alla scuola"})
		return
	}
	isGuardian, err := h.userRepo.IsGuardian(c.Request.Context(), userID, req.StudentID)
	if err != nil || !isGuardian {
		c.JSON(http.StatusForbidden, gin.H{"error": "rapporto genitore-studente non verificato"})
		return
	}
	student, err := h.userRepo.GetByID(c.Request.Context(), req.StudentID)
	if err != nil || student == nil || student.SchoolID == nil || *student.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "studente non associato alla scuola"})
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
	if role != "student" || h.userRepo == nil || schoolID == "" || userID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "la prenotazione deve essere effettuata dallo studente autenticato"})
		return
	}
	var req BookSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil || student == nil || student.Role != "student" || student.SchoolID == nil || *student.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "profilo studente non valido"})
		return
	}
	psychologist, err := h.userRepo.GetByID(c.Request.Context(), req.PsychologistID)
	if err != nil || psychologist == nil || psychologist.Role != "psychologist" || psychologist.SchoolID == nil || *psychologist.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "psicologo non assegnato alla scuola"})
		return
	}
	isMinor := true
	if student.DateOfBirth != nil {
		birth := *student.DateOfBirth
		age := time.Now().Year() - birth.Year()
		if time.Now().Before(birth.AddDate(age, 0, 0)) {
			age--
		}
		isMinor = age < 18
	}
	session, err := h.service.BookSession(c.Request.Context(), schoolID, userID, currentSchoolYear(time.Now()), isMinor, req)
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
	if role != "student" && role != "psychologist" {
		c.JSON(http.StatusForbidden, gin.H{"error": ErrProfessionalSecrecyViolation.Error()})
		return
	}
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
