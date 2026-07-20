package extracurricular

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
	g := r.Group("/extracurricular")
	{
		g.POST("/courses", h.CreateCourse)
		g.GET("/courses", h.ListCourses)
		g.POST("/courses/:id/enroll", h.EnrollStudent)
		g.GET("/courses/:id/enrollments", h.ListEnrollments)
		g.POST("/attendance", h.MarkAttendance)
	}
}

func (h *Handler) CreateCourse(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")

	if userID == "" || (role != "teacher" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.service.CreateCourse(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)
}

func (h *Handler) ListCourses(c *gin.Context) {
	schoolID := c.GetString("school_id")
	studentID := c.GetString("user_id")

	courses, err := h.service.ListCourses(c.Request.Context(), schoolID, studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func (h *Handler) EnrollStudent(c *gin.Context) {
	courseID := c.Param("id")
	studentID := c.GetString("user_id")

	if err := h.service.EnrollStudent(c.Request.Context(), courseID, studentID); err != nil {
		if err == ErrCourseFull {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "enrolled successfully"})
}

func (h *Handler) ListEnrollments(c *gin.Context) {
	courseID := c.Param("id")
	enrollments, err := h.service.ListEnrollments(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enrollments)
}

func (h *Handler) MarkAttendance(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	var req MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.MarkAttendance(c.Request.Context(), userID, role, req); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "attendance marked successfully"})
}
