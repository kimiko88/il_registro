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
	if q := c.Query("school_id"); q != "" {
		return q
	}
	if h := c.GetHeader("X-School-ID"); h != "" {
		return h
	}
	if res, exists := c.Get("school_id"); exists {
		if s, ok := res.(string); ok && s != "" {
			return s
		}
	}
	return "162737ff-081f-436c-8874-11cd57bc60f1"
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := getSchoolID(c)
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

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

	teachers, err := h.service.ListTeachers(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teachers)
}

func (h *Handler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := getSchoolID(c)
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teacher, err := h.service.GetTeacher(c.Request.Context(), schoolID, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teacher)
}

func (h *Handler) GetSubjects(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := getSchoolID(c)
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.service.GetTeacherSubjects(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) AssignSubject(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := getSchoolID(c)
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AssignSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AssignSubject(c.Request.Context(), role, schoolID, c.Param("id"), req.SubjectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) RemoveSubject(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := getSchoolID(c)
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.RemoveSubject(c.Request.Context(), role, schoolID, c.Param("id"), c.Param("subjectId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stats, err := h.service.GetDashboardStats(c.Request.Context(), userID, userID, role)
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
