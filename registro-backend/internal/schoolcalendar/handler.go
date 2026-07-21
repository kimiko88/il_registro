package schoolcalendar

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	cal := r.Group("/school-calendar")

	// Anno scolastico (solo segreteria/admin/superadmin)
	cal.PUT("/year", h.SetYear)
	cal.GET("/year", h.GetYear)

	// Giorni non didattici
	cal.POST("/non-teaching-days", h.AddNonTeachingDay)
	cal.DELETE("/non-teaching-days/:id", h.DeleteNonTeachingDay)
	cal.GET("/non-teaching-days", h.ListNonTeachingDays)

	// Periodi valutativi / Quadrimestri
	cal.POST("/periods", h.CreateAcademicPeriod)
	cal.GET("/periods", h.ListAcademicPeriods)
}

func (h *Handler) CreateAcademicPeriod(c *gin.Context) {
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id mancante nel token"})
		return
	}

	var req CreateAcademicPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	period, err := h.service.CreateAcademicPeriod(c.Request.Context(), actorRole, schoolID, req)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, period)
}

func (h *Handler) ListAcademicPeriods(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id mancante nel token"})
		return
	}

	periods, err := h.service.ListAcademicPeriods(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, periods)
}

func (h *Handler) SetYear(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id mancante nel token"})
		return
	}
	var req SetYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.SetSchoolYear(c.Request.Context(), actorID, actorRole, schoolID, req)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetYear(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id mancante nel token"})
		return
	}
	res, err := h.service.GetSchoolYear(c.Request.Context(), schoolID)
	if err != nil {
		if err.Error() == "anno scolastico non ancora configurato" || strings.Contains(err.Error(), "non ancora configurato") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) AddNonTeachingDay(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req AddNonTeachingDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.AddNonTeachingDay(c.Request.Context(), actorID, actorRole, schoolID, req)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) DeleteNonTeachingDay(c *gin.Context) {
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")
	id := c.Param("id")
	if err := h.service.DeleteNonTeachingDay(c.Request.Context(), actorRole, schoolID, id); err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListNonTeachingDays(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id mancante nel token"})
		return
	}
	res, err := h.service.ListNonTeachingDays(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
