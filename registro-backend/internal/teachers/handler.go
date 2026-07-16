package teachers

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

func getSchoolID(c *gin.Context) string {
	res, exists := c.Get("school_id")
	if !exists {
		return ""
	}
	return res.(string)
}

func (h *Handler) List(c *gin.Context) {
	// Check content filter
	subjectID := c.Query("subject_id")
	if subjectID != "" {
		teachers, err := h.service.GetTeachersBySubject(c.Request.Context(), subjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, teachers)
		return
	}

	schoolID := getSchoolID(c)
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required"})
		return
	}
	teachers, err := h.service.ListTeachers(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teachers)
}

func (h *Handler) Get(c *gin.Context) {
	teacher, err := h.service.GetTeacher(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teacher)
}

func (h *Handler) GetSubjects(c *gin.Context) {
	res, err := h.service.GetTeacherSubjects(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) AssignSubject(c *gin.Context) {
	var req AssignSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AssignSubject(c.Request.Context(), c.Param("id"), req.SubjectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) RemoveSubject(c *gin.Context) {
	if err := h.service.RemoveSubject(c.Request.Context(), c.Param("id"), c.Param("subjectId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stats, err := h.service.GetDashboardStats(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/teachers")
	{
		group.GET("", h.List)
		group.GET("/dashboard/stats", h.GetDashboardStats)
		group.GET("/:id", h.Get)
		group.GET("/:id/subjects", h.GetSubjects)
		group.POST("/:id/subjects", h.AssignSubject)
		group.DELETE("/:id/subjects/:subjectId", h.RemoveSubject)
	}
}
