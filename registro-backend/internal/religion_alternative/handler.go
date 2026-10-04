package religion_alternative

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
	g := rg.Group("/religion-alternative")
	{
		g.POST("/options", h.SaveOption)
		g.GET("/options", h.ListOptions)
		g.GET("/options/:studentId", h.GetOption)
		g.POST("/evaluations", h.SaveEvaluation)
		g.GET("/evaluations", h.ListEvaluations)
	}
}

func isAuthorizedStaff(role string) bool {
	switch role {
	case "admin", "superadmin", "secretary", "principal", "vice_principal", "teacher", "coordinator":
		return true
	default:
		return false
	}
}

func (h *Handler) SaveOption(c *gin.Context) {
	role := c.GetString("role")
	if !isAuthorizedStaff(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso non autorizzato"})
		return
	}

	var req struct {
		StudentID    string `json:"student_id" binding:"required"`
		AcademicYear string `json:"academic_year"`
		OptionType   string `json:"option_type" binding:"required"`
		Notes        string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati non validi: " + err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}

	opt := &StudentReligionOption{
		StudentID:    req.StudentID,
		SchoolID:     schoolID,
		AcademicYear: req.AcademicYear,
		OptionType:   req.OptionType,
		Notes:        req.Notes,
	}

	if err := h.service.RecordOption(c.Request.Context(), opt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Scelta opzione religione salvata con successo",
		"option":  opt,
	})
}

func (h *Handler) ListOptions(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}
	academicYear := c.DefaultQuery("academic_year", "2026/2027")

	list, err := h.service.ListOptions(c.Request.Context(), schoolID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

func (h *Handler) GetOption(c *gin.Context) {
	studentID := c.Param("studentId")
	academicYear := c.DefaultQuery("academic_year", "2026/2027")

	opt, err := h.service.GetOption(c.Request.Context(), studentID, academicYear)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "scelta non trovata"})
		return
	}

	c.JSON(http.StatusOK, opt)
}

func (h *Handler) SaveEvaluation(c *gin.Context) {
	role := c.GetString("role")
	if !isAuthorizedStaff(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso non autorizzato"})
		return
	}

	var req struct {
		StudentID        string  `json:"student_id" binding:"required"`
		GroupID          *string `json:"group_id"`
		ClassID          *string `json:"class_id"`
		Period           string  `json:"period" binding:"required"`
		SubjectKind      string  `json:"subject_kind" binding:"required"`
		JudgmentLevel    string  `json:"judgment_level" binding:"required"`
		DescriptiveNotes string  `json:"descriptive_notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati non validi: " + err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}
	userID := c.GetString("user_id")
	var createdBy *string
	if userID != "" {
		createdBy = &userID
	}

	eval := &AlternativeEvaluation{
		StudentID:        req.StudentID,
		SchoolID:         schoolID,
		GroupID:          req.GroupID,
		ClassID:          req.ClassID,
		Period:           req.Period,
		SubjectKind:      req.SubjectKind,
		JudgmentLevel:    req.JudgmentLevel,
		DescriptiveNotes: req.DescriptiveNotes,
		CreatedBy:        createdBy,
	}

	if err := h.service.RecordEvaluation(c.Request.Context(), eval); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Valutazione registrata con successo",
		"evaluation": eval,
	})
}

func (h *Handler) ListEvaluations(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}
	period := c.Query("period")
	subjectKind := c.Query("subject_kind")

	list, err := h.service.ListEvaluations(c.Request.Context(), schoolID, period, subjectKind)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}
