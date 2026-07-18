package lessons

import (
	"net/http"

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
		lessons.GET("/class/:class_id", h.GetLessons)
		lessons.POST("", h.CreateLesson)
	}

	homeworks := rg.Group("/homeworks")
	{
		homeworks.GET("/class/:class_id", h.GetHomeworks)
		homeworks.POST("", h.CreateHomework)
	}
}

func (h *Handler) GetLessons(c *gin.Context) {
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

func (h *Handler) CreateLesson(c *gin.Context) {
	teacherID := c.GetString("user_id") // From auth middleware

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
