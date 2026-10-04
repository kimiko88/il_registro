package middleschoolexam

import (
	"fmt"
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
	m := rg.Group("/middle-school-exam")
	{
		m.POST("/classes/:classId", h.GetOrCreateExam)
		m.GET("/:examId/candidates", h.ListCandidates)
		m.POST("/:examId/candidates/:studentId/admission", h.SaveAdmission)
		m.POST("/:examId/candidates/:studentId/evaluate", h.EvaluateCandidate)
		m.GET("/candidates/:candidateId/diploma", h.GenerateDiploma)
		m.PUT("/:examId/status", h.UpdateStatus)
	}
}

func isTeacherOrSecretary(role string) bool {
	switch role {
	case "teacher", "coordinator", "coordinatore_classe", "admin", "superadmin", "secretary", "principal", "vice_principal":
		return true
	default:
		return false
	}
}

func (h *Handler) GetOrCreateExam(c *gin.Context) {
	role := c.GetString("role")
	if !isTeacherOrSecretary(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	schoolID := c.GetString("school_id")
	classID := c.Param("classId")
	academicYear := c.DefaultQuery("academic_year", "2026/2027")

	var req struct {
		PresidentName string `json:"president_name"`
	}
	_ = c.ShouldBindJSON(&req)

	exam, err := h.service.GetOrCreateExam(c.Request.Context(), schoolID, classID, academicYear, req.PresidentName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exam)
}

func (h *Handler) ListCandidates(c *gin.Context) {
	role := c.GetString("role")
	if !isTeacherOrSecretary(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	examID := c.Param("examId")
	candidates, err := h.service.ListCandidates(c.Request.Context(), examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, candidates)
}

func (h *Handler) SaveAdmission(c *gin.Context) {
	role := c.GetString("role")
	if !isTeacherOrSecretary(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	examID := c.Param("examId")
	studentID := c.Param("studentId")

	var req struct {
		Grade    int    `json:"grade"`
		Judgment string `json:"judgment"`
		Admitted bool   `json:"admitted"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SaveAdmission(c.Request.Context(), examID, studentID, req.Grade, req.Judgment, req.Admitted); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ammissione salvata"})
}

func (h *Handler) EvaluateCandidate(c *gin.Context) {
	role := c.GetString("role")
	if !isTeacherOrSecretary(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	examID := c.Param("examId")
	studentID := c.Param("studentId")

	var req struct {
		Grades ExamCandidateGrades `json:"grades"`
		Notes  string              `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.EvaluateCandidate(c.Request.Context(), examID, studentID, req.Grades, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GenerateDiploma(c *gin.Context) {
	candidateID := c.Param("candidateId")
	text, err := h.service.GenerateDiploma(c.Request.Context(), candidateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=diploma_%s.txt", candidateID))
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(text))
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	role := c.GetString("role")
	if !isTeacherOrSecretary(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	examID := c.Param("examId")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), examID, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stato esame aggiornato"})
}
