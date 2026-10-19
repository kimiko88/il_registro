package competencies_invalsi

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
	group := rg.Group("/competencies-dm14")
	{
		group.POST("/certifications", h.SaveCertification)
		group.GET("/certifications/:id", h.GetCertification)
		group.GET("/certifications/by-student", h.GetCertificationByStudent)
		group.GET("/certifications/by-class", h.ListCertificationsByClass)
		group.POST("/certifications/:id/sign", h.SignCertification)

		group.POST("/invalsi/batch", h.ImportInvalsiBatch)
		group.GET("/invalsi/by-class", h.ListInvalsiByClass)
		group.GET("/invalsi/analytics", h.GetSchoolAnalytics)
	}
}

func (h *Handler) SaveCertification(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var cert CompetencyCertificationDM14
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cert.SchoolID = schoolID

	if err := h.svc.SaveCertification(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cert)
}

func (h *Handler) GetCertification(c *gin.Context) {
	id := c.Param("id")
	cert, err := h.svc.GetCertification(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificazione non trovata"})
		return
	}
	c.JSON(http.StatusOK, cert)
}

func (h *Handler) GetCertificationByStudent(c *gin.Context) {
	studentID := c.Query("student_id")
	schoolYear := c.DefaultQuery("school_year", "2025/2026")
	schoolTier := c.DefaultQuery("school_tier", "primo_ciclo")

	cert, err := h.svc.GetCertificationByStudent(studentID, schoolYear, schoolTier)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificazione non trovata"})
		return
	}
	c.JSON(http.StatusOK, cert)
}

func (h *Handler) ListCertificationsByClass(c *gin.Context) {
	classID := c.Query("class_id")
	schoolYear := c.DefaultQuery("school_year", "2025/2026")

	list, err := h.svc.ListCertificationsByClass(classID, schoolYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"certifications": list})
}

func (h *Handler) SignCertification(c *gin.Context) {
	id := c.Param("id")
	principalID := c.GetString("user_id")

	if err := h.svc.SignCertification(id, principalID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "certificazione firmata digitalmente con successo"})
}

func (h *Handler) ImportInvalsiBatch(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req struct {
		Results []InvalsiResult `json:"results" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range req.Results {
		req.Results[i].SchoolID = schoolID
	}

	if err := h.svc.ImportInvalsiBatch(req.Results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "dati INVALSI importati con successo", "count": len(req.Results)})
}

func (h *Handler) ListInvalsiByClass(c *gin.Context) {
	classID := c.Query("class_id")
	schoolYear := c.DefaultQuery("school_year", "2025/2026")

	list, err := h.svc.ListInvalsiByClass(classID, schoolYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": list})
}

func (h *Handler) GetSchoolAnalytics(c *gin.Context) {
	schoolID := c.GetString("school_id")
	schoolYear := c.DefaultQuery("school_year", "2025/2026")

	list, err := h.svc.GetSchoolAnalytics(schoolID, schoolYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"analytics": list})
}
