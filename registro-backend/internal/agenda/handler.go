package agenda

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func parseFlexibleDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" || s == "undefined" || s == "null" {
		return time.Time{}
	}
	if len(s) >= 10 {
		sDate := s[:10]
		if t, err := time.Parse("2006-01-02", sDate); err == nil {
			return t
		}
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return time.Time{}
}

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

	if dateStr := c.Query("date"); dateStr != "" && c.Query("from") == "" {
		dTime := parseFlexibleDate(dateStr)
		if !dTime.IsZero() {
			fromTime = dTime
			toTime = dTime
		}
	}
	if fromStr := c.Query("from"); fromStr != "" {
		fromTime = parseFlexibleDate(fromStr)
	}
	if toStr := c.Query("to"); toStr != "" {
		toTime = parseFlexibleDate(toStr)
	}

	if !fromTime.IsZero() && !toTime.IsZero() {
		if toTime.Before(fromTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "la data di fine non può essere precedente alla data di inizio"})
			return
		}
		if toTime.Sub(fromTime) > 366*24*time.Hour {
			c.JSON(http.StatusBadRequest, gin.H{"error": "range di date troppo ampio (massimo 1 anno consentito)"})
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

	role := c.GetString("role")
	items, err := h.service.GetCalendar(c.Request.Context(), schoolID, userID, role, filter)
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
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
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
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
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
	role := c.GetString("role")
	id := c.Param("id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.SetTaskCompletion(c.Request.Context(), studentID, role, id, true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked complete"})
}

func (h *Handler) UnmarkComplete(c *gin.Context) {
	studentID := c.GetString("user_id")
	role := c.GetString("role")
	id := c.Param("id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.SetTaskCompletion(c.Request.Context(), studentID, role, id, false); err != nil {
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

	role := c.GetString("role")
	items, err := h.service.GetCalendar(c.Request.Context(), schoolID, userID, role, filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if items == nil {
		items = []*AgendaItem{}
	}
	c.JSON(http.StatusOK, items)
}
