package agenda

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	ag := r.Group("/agenda")
	{
		ag.POST("", h.Create)
		ag.GET("", h.GetCalendar)
		ag.GET("/calendar", h.GetCalendar)
		ag.GET("/class/:classID", h.GetClassEvents)
		ag.GET("/:id", h.GetByID)
		ag.PUT("/:id", h.Update)
		ag.PATCH("/:id", h.Update)
		ag.DELETE("/:id", h.Delete)
		ag.POST("/:id/complete", h.MarkComplete)
		ag.DELETE("/:id/complete", h.UnmarkComplete)
	}
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || (role != "teacher" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateAgendaItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.CreateAgendaItem(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) GetCalendar(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	classID := c.Query("class_id")
	if classID == "" {
		classID = c.Query("classID")
	}
	subjectID := c.Query("subject_id")
	agendaType := c.Query("type")

	var fromTime, toTime time.Time
	var err error

	if dateStr := c.Query("date"); dateStr != "" && c.Query("from") == "" {
		if dTime, pErr := time.Parse("2006-01-02", dateStr); pErr == nil {
			fromTime = dTime
			toTime = dTime
		}
	}
	if fromStr := c.Query("from"); fromStr != "" {
		fromTime, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date format (expected YYYY-MM-DD)"})
			return
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		toTime, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date format (expected YYYY-MM-DD)"})
			return
		}
	}

	filter := CalendarFilter{
		ClassID:   classID,
		SubjectID: subjectID,
		Type:      agendaType,
		From:      fromTime,
		To:        toTime,
	}

	items, err := h.service.GetCalendar(c.Request.Context(), schoolID, userID, filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	item, err := h.service.GetAgendaItem(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agenda item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	id := c.Param("id")

	var req UpdateAgendaItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.UpdateAgendaItem(c.Request.Context(), userID, role, id, req)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	id := c.Param("id")

	if err := h.service.DeleteAgendaItem(c.Request.Context(), userID, role, id); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) MarkComplete(c *gin.Context) {
	studentID := c.GetString("user_id")
	id := c.Param("id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.SetTaskCompletion(c.Request.Context(), studentID, id, true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked complete"})
}

func (h *Handler) UnmarkComplete(c *gin.Context) {
	studentID := c.GetString("user_id")
	id := c.Param("id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.SetTaskCompletion(c.Request.Context(), studentID, id, false); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unmarked complete"})
}

func (h *Handler) GetClassEvents(c *gin.Context) {
	classID := c.Param("classID")
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filter := CalendarFilter{
		ClassID: classID,
	}

	items, err := h.service.GetCalendar(c.Request.Context(), schoolID, userID, filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if items == nil {
		items = []*AgendaItem{}
	}
	c.JSON(http.StatusOK, items)
}
