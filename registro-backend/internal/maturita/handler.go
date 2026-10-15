package maturita

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
	g := r.Group("/maturita")
	{
		g.POST("/commission", h.SaveCommission)
		g.GET("/commission", h.GetCommission)
		g.POST("/student-record", h.SaveStudentRecord)
		g.GET("/student-record", h.GetStudentRecord)
		g.GET("/tabellone", h.GetTabellone)
		g.GET("/curriculum-xml", h.ExportCurriculumXML)
		g.POST("/calculate-credits", h.CalculateCredits)
		g.POST("/calculate-scores", h.CalculateScores)
	}
}

func (h *Handler) SaveCommission(c *gin.Context) {
	var comm CommissioneMaturita
	if err := c.ShouldBindJSON(&comm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comm.SchoolID = c.GetString("school_id")
	if err := h.service.SaveCommission(c.Request.Context(), &comm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comm)
}

func (h *Handler) GetCommission(c *gin.Context) {
	classID := c.Query("class_id")
	year := c.DefaultQuery("school_year", "2025/2026")
	comm, err := h.service.GetCommission(c.Request.Context(), classID, year)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comm)
}

func (h *Handler) SaveStudentRecord(c *gin.Context) {
	var rec StudentMaturitaRecord
	if err := c.ShouldBindJSON(&rec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec.SchoolID = c.GetString("school_id")
	if err := h.service.SaveStudentRecord(c.Request.Context(), &rec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rec)
}

func (h *Handler) GetStudentRecord(c *gin.Context) {
	studentID := c.Query("student_id")
	year := c.DefaultQuery("school_year", "2025/2026")
	rec, err := h.service.GetStudentRecord(c.Request.Context(), studentID, year)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rec)
}

func (h *Handler) GetTabellone(c *gin.Context) {
	classID := c.Query("class_id")
	year := c.DefaultQuery("school_year", "2025/2026")
	records, err := h.service.GetTabelloneClasse(c.Request.Context(), classID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"class_id":    classID,
		"school_year": year,
		"total":       len(records),
		"tabellone":   records,
	})
}

func (h *Handler) ExportCurriculumXML(c *gin.Context) {
	studentID := c.Query("student_id")
	year := c.DefaultQuery("school_year", "2025/2026")
	rec, err := h.service.GetStudentRecord(c.Request.Context(), studentID, year)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	xmlStr, err := h.service.ExportCurriculumXML(rec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=curriculum_studente.xml")
	c.String(http.StatusOK, xmlStr)
}

func (h *Handler) CalculateCredits(c *gin.Context) {
	var req struct {
		Grade3rd float64 `json:"grade_3rd"`
		High3rd  bool    `json:"high_3rd"`
		Grade4th float64 `json:"grade_4th"`
		High4th  bool    `json:"high_4th"`
		Grade5th float64 `json:"grade_5th"`
		High5th  bool    `json:"high_5th"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c3, err := CalculateYearCredits(3, req.Grade3rd, req.High3rd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c4, err := CalculateYearCredits(4, req.Grade4th, req.High4th)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c5, err := CalculateYearCredits(5, req.Grade5th, req.High5th)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tot, err := CalculateTrienniumCredits(c3, c4, c5)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TrienniumCredits{
		Grade3rd:     req.Grade3rd,
		Credit3rd:    c3,
		Grade4th:     req.Grade4th,
		Credit4th:    c4,
		Grade5th:     req.Grade5th,
		Credit5th:    c5,
		TotalCredits: tot,
	})
}

func (h *Handler) CalculateScores(c *gin.Context) {
	var req struct {
		Credits       TrienniumCredits `json:"credits"`
		Written1Score float64          `json:"written1_score"`
		Written2Score float64          `json:"written2_score"`
		OralScore     float64          `json:"oral_score"`
		BonusPoints   int              `json:"bonus_points"`
		WantLode      bool             `json:"want_lode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scores, err := CalculateFinalExamScores(req.Credits, req.Written1Score, req.Written2Score, req.OralScore, req.BonusPoints, req.WantLode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scores)
}
