package pcto_tutor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	pt := rg.Group("/pcto-tutor")
	{
		// Internal admin/secretary endpoint to register tutor and generate link
		pt.POST("/tutors", h.RegisterTutor)

		// Public endpoints for company tutor using magic link
		pt.GET("/session", h.GetSession)
		pt.GET("/students", h.ListAssignedStudents)
		pt.POST("/timesheets/verify", h.VerifyTimesheet)
		pt.POST("/evaluations", h.SubmitEvaluation)
	}
}

func (h *Handler) extractTutorToken(c *gin.Context) string {
	token := c.GetHeader("X-Tutor-Token")
	if token == "" {
		token = c.Query("token")
	}
	return token
}

func (h *Handler) RegisterTutor(c *gin.Context) {
	role := c.GetString("role")
	switch role {
	case "admin", "superadmin", "secretary", "principal", "vice_principal", "teacher", "coordinator":
		// Allowed
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso non autorizzato"})
		return
	}

	var req struct {
		CompanyName    string `json:"company_name" binding:"required"`
		TutorFirstName string `json:"tutor_first_name" binding:"required"`
		TutorLastName  string `json:"tutor_last_name" binding:"required"`
		Email          string `json:"email" binding:"required"`
		Phone          string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati tutor non validi: " + err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}

	tutor := &CompanyTutor{
		SchoolID:       schoolID,
		CompanyName:    req.CompanyName,
		TutorFirstName: req.TutorFirstName,
		TutorLastName:  req.TutorLastName,
		Email:          req.Email,
		Phone:          req.Phone,
	}

	if err := h.service.RegisterTutor(c.Request.Context(), tutor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tutor aziendale registrato e magic link generato",
		"tutor":   tutor,
	})
}

func (h *Handler) GetSession(c *gin.Context) {
	token := h.extractTutorToken(c)
	tutor, err := h.service.AuthenticateWithToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sessione tutor non valida o scaduta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tutor": tutor,
	})
}

func (h *Handler) ListAssignedStudents(c *gin.Context) {
	token := h.extractTutorToken(c)
	tutor, err := h.service.AuthenticateWithToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sessione tutor non valida o scaduta"})
		return
	}

	students, err := h.service.GetAssignedStudents(c.Request.Context(), tutor.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": students,
	})
}

func (h *Handler) VerifyTimesheet(c *gin.Context) {
	token := h.extractTutorToken(c)
	tutor, err := h.service.AuthenticateWithToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sessione tutor non valida o scaduta"})
		return
	}

	var req struct {
		ProjectID     string  `json:"project_id" binding:"required"`
		StudentID     string  `json:"student_id" binding:"required"`
		ActivityDate  string  `json:"activity_date" binding:"required"`
		HoursDeclared float64 `json:"hours_declared" binding:"required"`
		HoursApproved float64 `json:"hours_approved" binding:"required"`
		TutorNotes    string  `json:"tutor_notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati verifica ore non validi: " + err.Error()})
		return
	}

	entry := &TimesheetVerification{
		ProjectID:     req.ProjectID,
		StudentID:     req.StudentID,
		ActivityDate:  req.ActivityDate,
		HoursDeclared: req.HoursDeclared,
		HoursApproved: req.HoursApproved,
		TutorID:       tutor.ID,
		TutorNotes:    req.TutorNotes,
	}

	if err := h.service.VerifyTimesheet(c.Request.Context(), entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Presenze certificate dal tutor aziendale",
		"verification": entry,
	})
}

func (h *Handler) SubmitEvaluation(c *gin.Context) {
	token := h.extractTutorToken(c)
	tutor, err := h.service.AuthenticateWithToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sessione tutor non valida o scaduta"})
		return
	}

	var req struct {
		ProjectID        string `json:"project_id" binding:"required"`
		StudentID        string `json:"student_id" binding:"required"`
		ReliabilityLevel int    `json:"reliability_level" binding:"required"`
		TechnicalSkills  int    `json:"technical_skills" binding:"required"`
		TeamworkSkills   int    `json:"teamwork_skills" binding:"required"`
		FinalFeedback    string `json:"final_feedback" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati valutazione non validi: " + err.Error()})
		return
	}

	eval := &CompanyEvaluation{
		TutorID:          tutor.ID,
		StudentID:        req.StudentID,
		ProjectID:        req.ProjectID,
		ReliabilityLevel: req.ReliabilityLevel,
		TechnicalSkills:  req.TechnicalSkills,
		TeamworkSkills:   req.TeamworkSkills,
		FinalFeedback:    req.FinalFeedback,
	}

	if err := h.service.SubmitEvaluation(c.Request.Context(), eval); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Valutazione competenze aziendali registrata con successo",
		"evaluation": eval,
	})
}
