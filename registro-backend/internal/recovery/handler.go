package recovery

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

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/recovery")
	{
		g.POST("/courses", h.CreateCourse)
		g.GET("/courses", h.ListCourses)
		g.GET("/courses/:id", h.GetCourse)
		g.PATCH("/courses/:id/status", h.UpdateStatus)
		g.PATCH("/courses/:id/students/:student_id/attendance", h.UpdateAttendance)
		g.POST("/tests", h.RecordTest)
		g.GET("/tests", h.ListTests)
	}
}

func (h *Handler) CreateCourse(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	course, err := h.svc.CreateCourse(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *Handler) GetCourse(c *gin.Context) {
	id := c.Param("id")
	course, err := h.svc.GetCourse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *Handler) ListCourses(c *gin.Context) {
	schoolID := c.GetString("school_id")
	academicYear := c.Query("academic_year")
	teacherID := c.Query("teacher_id")

	list, err := h.svc.ListCourses(c.Request.Context(), schoolID, academicYear, teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []RecoveryCourse{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), id, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) UpdateAttendance(c *gin.Context) {
	courseID := c.Param("id")
	studentID := c.Param("student_id")
	var body struct {
		AttendanceHours float64 `json:"attendance_hours"`
		Notes           string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateAttendance(c.Request.Context(), courseID, studentID, body.AttendanceHours, body.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) RecordTest(c *gin.Context) {
	var req RecordTestOutcomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	teacherID := c.GetString("teacher_id")

	test, err := h.svc.RecordTestOutcome(c.Request.Context(), schoolID, teacherID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, test)
}

func (h *Handler) ListTests(c *gin.Context) {
	schoolID := c.GetString("school_id")
	classID := c.Query("class_id")
	studentID := c.Query("student_id")

	list, err := h.svc.ListTests(c.Request.Context(), schoolID, classID, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []RecoveryTest{}
	}
	c.JSON(http.StatusOK, list)
}
