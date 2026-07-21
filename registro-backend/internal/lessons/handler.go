package lessons

import (
	"net/http"
	"strings"
	"time"

	"registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	lessons := rg.Group("/lessons")
	{
		lessons.GET("/my-diary", h.GetMyDiary)
		lessons.GET("/class/:class_id", h.GetLessons)
		lessons.GET("/group/:group_id", h.GetLessonsByGroup)
		lessons.GET("/:id", h.GetLessonByID)
		lessons.PUT("/:id", h.UpdateLesson)
		lessons.DELETE("/:id", h.DeleteLesson)
		lessons.POST("", h.CreateLesson)
	}

	homeworks := rg.Group("/homeworks")
	{
		homeworks.GET("/class/:class_id", h.GetHomeworks)
		homeworks.PUT("/:id", h.UpdateHomework)
		homeworks.DELETE("/:id", h.DeleteHomework)
		homeworks.POST("", h.CreateHomework)
	}
}

// parseDate valida un parametro data in formato YYYY-MM-DD.
func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func (h *Handler) GetLessons(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	classID := c.Param("class_id")
	subjectID := c.Query("subject_id")
	date := c.Query("date")

	res, err := h.service.GetLessons(classID, subjectID, date)
	if err != nil {
		logger.Log.Errorf("GetLessons error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetLessonsByGroup(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	groupID := c.Param("group_id")
	date := c.Query("date")

	res, err := h.service.GetLessonsByGroup(groupID, date)
	if err != nil {
		logger.Log.Errorf("GetLessonsByGroup error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateLesson(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only teachers can create lessons"})
		return
	}

	var req CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateLesson(teacherID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetHomeworks(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	classID := c.Param("class_id")

	res, err := h.service.GetHomeworks(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateHomework(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only teachers can create homework"})
		return
	}

	var req CreateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateHomework(teacherID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetLessonByID(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	res, err := h.service.GetLessonByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateLesson(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	var req UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.UpdateLesson(teacherID, role, id, req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteLesson(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	if err := h.service.DeleteLesson(teacherID, role, id); err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateHomework(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	var req UpdateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.UpdateHomework(teacherID, role, id, req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteHomework(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	if err := h.service.DeleteHomework(teacherID, role, id); err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetMyDiary(c *gin.Context) {
	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	from := c.Query("from")
	to := c.Query("to")

	var fromTime, toTime time.Time
	if from != "" {
		var err error
		if fromTime, err = parseDate(from); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parametro 'from' non valido, formato atteso: YYYY-MM-DD"})
			return
		}
	}
	if to != "" {
		var err error
		if toTime, err = parseDate(to); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parametro 'to' non valido, formato atteso: YYYY-MM-DD"})
			return
		}
	}

	if !fromTime.IsZero() && !toTime.IsZero() && toTime.Sub(fromTime) > 365*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "range di date troppo ampio (massimo 1 anno consentito)"})
		return
	}

	res, err := h.service.GetTeacherDiary(teacherID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
