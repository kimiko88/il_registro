package teacher_activities

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Handler gestisce le route HTTP per le attività libere del docente.
type Handler struct {
	service Service
}

// NewHandler restituisce un nuovo Handler.
func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes registra le route nel gruppo fornito.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/teacher/free-activities")
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.GET("/:id", h.GetByID)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

// List restituisce le attività libere del docente autenticato.
// Query params: from (YYYY-MM-DD), to (YYYY-MM-DD)
func (h *Handler) List(c *gin.Context) {
	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	from := c.Query("from")
	to := c.Query("to")

	// Default: ultimi 30 giorni – prossimi 30 giorni
	if from == "" && to == "" {
		now := time.Now()
		from = now.AddDate(0, 0, -30).Format("2006-01-02")
		to = now.AddDate(0, 0, 30).Format("2006-01-02")
	}

	res, err := h.service.GetByTeacher(teacherID, from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if res == nil {
		res = []TeacherActivityResponse{}
	}
	c.JSON(http.StatusOK, res)
}

// GetByID restituisce una singola attività per ID.
func (h *Handler) GetByID(c *gin.Context) {
	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	res, err := h.service.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows || strings.Contains(err.Error(), "no rows") {
			c.JSON(http.StatusNotFound, gin.H{"error": "attività non trovata"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Create crea una nuova attività libera per il docente autenticato.
func (h *Handler) Create(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo i docenti possono registrare attività libere"})
		return
	}

	var req CreateTeacherActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Create(teacherID, req)
	if err != nil {
		logger.Log.Errorf("CreateTeacherActivity error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Update aggiorna un'attività libera esistente.
func (h *Handler) Update(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id := c.Param("id")
	var req UpdateTeacherActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.Update(teacherID, id, req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Delete elimina un'attività libera.
func (h *Handler) Delete(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id := c.Param("id")
	if err := h.service.Delete(teacherID, id); err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
