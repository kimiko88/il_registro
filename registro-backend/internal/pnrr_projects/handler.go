package pnrr_projects

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
	group := rg.Group("/pnrr-projects")
	{
		group.POST("", h.CreateProject)
		group.GET("", h.ListProjects)
		group.GET("/:id", h.GetProject)
		group.POST("/:id/enroll", h.EnrollStudent)
		group.GET("/:id/students", h.ListStudents)
		group.POST("/:id/lessons", h.RecordLesson)
		group.GET("/:id/lessons", h.ListLessons)
	}
}

func (h *Handler) CreateProject(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var p PNRRProject
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.SchoolID = schoolID

	if err := h.svc.CreateProject(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) ListProjects(c *gin.Context) {
	schoolID := c.GetString("school_id")
	academicYear := c.DefaultQuery("academic_year", "2025/2026")

	list, err := h.svc.ListProjects(schoolID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": list})
}

func (h *Handler) GetProject(c *gin.Context) {
	id := c.Param("id")
	p, err := h.svc.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "progetto PNRR non trovato"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) EnrollStudent(c *gin.Context) {
	projectID := c.Param("id")
	var req struct {
		StudentID string `json:"student_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.EnrollStudent(projectID, req.StudentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "studente iscritto al corso PNRR"})
}

func (h *Handler) ListStudents(c *gin.Context) {
	projectID := c.Param("id")
	list, err := h.svc.ListStudents(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"students": list})
}

func (h *Handler) RecordLesson(c *gin.Context) {
	projectID := c.Param("id")
	var req struct {
		LessonDate         string   `json:"lesson_date" binding:"required"`
		Hours              int      `json:"hours" binding:"required"`
		Topic              string   `json:"topic" binding:"required"`
		AttendedStudentIDs []string `json:"attended_student_ids" binding:"required"`
		SignaturesJSON     string   `json:"signatures_json"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson := &PNRRLessonAttendance{
		ProjectID:      projectID,
		LessonDate:     req.LessonDate,
		Hours:          req.Hours,
		Topic:          req.Topic,
		SignaturesJSON: req.SignaturesJSON,
	}

	if err := h.svc.RecordLesson(lesson, req.AttendedStudentIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "lezione e presenze PNRR registrate"})
}

func (h *Handler) ListLessons(c *gin.Context) {
	projectID := c.Param("id")
	list, err := h.svc.ListLessons(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"lessons": list})
}
