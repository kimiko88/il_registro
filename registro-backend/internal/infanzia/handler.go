package infanzia

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/infanzia")
	{
		group.POST("/sections", h.CreateSection)
		group.GET("/sections", h.ListSections)
		group.GET("/sections/:id", h.GetSection)

		group.POST("/sections/:id/attendances", h.SaveAttendances)
		group.GET("/sections/:id/attendances", h.GetAttendances)

		group.POST("/observations", h.RecordObservation)
		group.GET("/observations/by-student", h.ListStudentObservations)

		group.POST("/sections/:id/diario", h.AddDiarioStory)
		group.GET("/sections/:id/diario", h.ListDiarioStories)
	}
}

func (h *Handler) CreateSection(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var sec InfanziaSection
	if err := c.ShouldBindJSON(&sec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sec.SchoolID = schoolID

	if err := h.svc.CreateSection(&sec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sec)
}

func (h *Handler) ListSections(c *gin.Context) {
	schoolID := c.GetString("school_id")
	academicYear := c.DefaultQuery("academic_year", "2025/2026")

	list, err := h.svc.ListSections(schoolID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sections": list})
}

func (h *Handler) GetSection(c *gin.Context) {
	id := c.Param("id")
	sec, err := h.svc.GetSection(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sezione infanzia non trovata"})
		return
	}
	c.JSON(http.StatusOK, sec)
}

func (h *Handler) SaveAttendances(c *gin.Context) {
	sectionID := c.Param("id")
	var req struct {
		Date        string                    `json:"date" binding:"required"`
		Attendances []InfanziaDailyAttendance `json:"attendances" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.SaveAttendances(sectionID, req.Date, req.Attendances); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "presenze e pasti salvati per la sezione"})
}

func (h *Handler) GetAttendances(c *gin.Context) {
	sectionID := c.Param("id")
	date := c.Query("date")

	list, err := h.svc.GetAttendances(sectionID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attendances": list})
}

func (h *Handler) RecordObservation(c *gin.Context) {
	teacherID := c.GetString("user_id")
	var obs InfanziaExperienceObservation
	if err := c.ShouldBindJSON(&obs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	obs.RecordedBy = teacherID

	if err := h.svc.RecordObservation(&obs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, obs)
}

func (h *Handler) ListStudentObservations(c *gin.Context) {
	studentID := c.Query("student_id")
	list, err := h.svc.ListStudentObservations(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"observations": list})
}

func (h *Handler) AddDiarioStory(c *gin.Context) {
	sectionID := c.Param("id")
	teacherID := c.GetString("user_id")

	var d InfanziaDiarioBordo
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d.SectionID = sectionID
	d.CreatedBy = teacherID

	if err := h.svc.AddDiarioStory(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *Handler) ListDiarioStories(c *gin.Context) {
	sectionID := c.Param("id")
	limitStr := c.DefaultQuery("limit", "30")
	limit, _ := strconv.Atoi(limitStr)

	list, err := h.svc.ListDiarioStories(sectionID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stories": list})
}
