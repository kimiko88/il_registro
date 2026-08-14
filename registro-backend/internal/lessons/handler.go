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
		lessons.GET("/class/:class_id/activity-hours", h.GetActivityHours)
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
		if strings.HasPrefix(err.Error(), "forbidden") || strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		if strings.HasPrefix(err.Error(), "forbidden") || strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		if strings.Contains(err.Error(), "not found") || err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only teachers or staff can update lessons"})
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
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo docenti, dirigenti o segreteria possono eliminare le lezioni"})
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
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo docenti, dirigenti o segreteria possono modificare i compiti"})
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
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo docenti, dirigenti o segreteria possono eliminare i compiti"})
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

	if from == "" && to == "" {
		now := time.Now()
		from = now.AddDate(0, 0, -15).Format("2006-01-02")
		to = now.AddDate(0, 0, 15).Format("2006-01-02")
	}

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

// GetActivityHours restituisce il conteggio delle ore suddivise per tipo di attività
// per una data classe. Usato dal LessonPlanner per mostrare i contatori PCTO/Orientamento.
func (h *Handler) GetActivityHours(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	classID := c.Param("class_id")

	// Recupera tutte le lezioni della classe (senza filtro data per avere il totale)
	lessons, err := h.service.GetLessons(classID, "", "")
	if err != nil {
		logger.Log.Errorf("GetActivityHours error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Conta le ore per tipo di attività
	counters := map[string]int{
		"pcto":         0,
		"orientamento": 0,
		"ptof":         0,
		"assembly":     0,
		"trip":         0,
		"project":      0,
		"lab":          0,
		"standard":     0,
		"other":        0,
	}
	for _, l := range lessons {
		key := l.ActivityType
		if key == "" {
			key = "standard"
		}
		dur := l.Duration
		if dur <= 0 {
			dur = 1
		}
		if _, ok := counters[key]; ok {
			counters[key] += dur
		} else {
			counters["other"] += dur
		}
	}
	c.JSON(http.StatusOK, counters)
}
