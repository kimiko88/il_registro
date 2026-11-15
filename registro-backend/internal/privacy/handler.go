package privacy

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

func privacySchoolYear(now time.Time) string {
	startYear := now.Year()
	if now.Month() < time.September { startYear-- }
	return fmt.Sprintf("%d/%d", startYear, startYear+1)
}

func isPrivacyStaffRole(role string) bool {
	switch role {
	case "admin", "superadmin", "principal", "dpo", "privacy_officer":
		return true
	default:
		return false
	}
}

func isSchoolStaffRole(role string) bool {
	switch role {
	case "admin", "superadmin", "principal", "vice_principal", "secretary", "dsga", "teacher", "coordinator":
		return true
	default:
		return false
	}
}

func (h *Handler) authorizeStudent(c *gin.Context, studentID string, allowGuardian bool) bool {
	if h.userRepo == nil || c.GetString("user_id") == "" || c.GetString("school_id") == "" {
			return false
		}
	role := c.GetString("role")
	if allowGuardian && role == "parent" {
		parent, err := h.userRepo.GetByID(c.Request.Context(), c.GetString("user_id"))
		if err != nil || parent.SchoolID == nil || *parent.SchoolID != c.GetString("school_id") {
		return false
	}
		ok, err := h.userRepo.IsGuardian(c.Request.Context(), c.GetString("user_id"), studentID)
		return err == nil && ok
	}
	if !isSchoolStaffRole(role) && role != "student" { return false }
	student, err := h.userRepo.GetByID(c.Request.Context(), studentID)
	if err != nil || student.SchoolID == nil || *student.SchoolID != c.GetString("school_id") { return false }
	if role == "student" { return student.ID == c.GetString("user_id") }
	return true
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

	if c.GetString("role") != "parent" || !h.authorizeStudent(c, req.StudentID, true) {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo un genitore/tutore verificato dello studente può registrare il consenso"})
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
	year := c.DefaultQuery("school_year", privacySchoolYear(time.Now()))
	if !h.authorizeStudent(c, studentID, true) {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso al consenso non autorizzato"})
		return
	}
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
	if req.SchoolYear == "" { req.SchoolYear = privacySchoolYear(time.Now()) }
	if len(req.StudentIDs) == 0 || len(req.StudentIDs) > 100 || !isSchoolStaffRole(c.GetString("role")) {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso o parametri non autorizzati"})
		return
	}
	for _, studentID := range req.StudentIDs {
		if !h.authorizeStudent(c, studentID, false) {
			c.JSON(http.StatusForbidden, gin.H{"error": "uno o più studenti non appartengono alla scuola autorizzata"})
			return
		}
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
	if !isPrivacyStaffRole(c.GetString("role")) || schoolID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso al registro dei trattamenti non autorizzato"})
		return
	}
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
	if !isPrivacyStaffRole(c.GetString("role")) || schoolID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso al registro dei trattamenti non autorizzato"})
		return
	}
	treatments, err := h.service.ListTreatments(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(treatments), "treatments": treatments})
}
