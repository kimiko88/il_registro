package substitutions

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
	g := r.Group("/substitutions")
	{
		g.POST("", h.Create)
		g.GET("", h.ListBySchool)
		g.GET("/my", h.ListByTeacher)
		g.GET("/my-today", h.ListMyToday)
		g.PUT("/:id/assign", h.AssignSubstitute)
		g.PATCH("/:id/confirm", h.Confirm)
	}
}

func (h *Handler) Create(c *gin.Context) {
	if c.GetString("user_id") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateSubstitutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.CreateSubstitution(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sub)
}

func (h *Handler) ListBySchool(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	schoolID := c.GetString("school_id")
	date := c.Query("date")

	subs, err := h.service.ListBySchool(c.Request.Context(), schoolID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *Handler) ListByTeacher(c *gin.Context) {
	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subs, err := h.service.ListByTeacher(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *Handler) ListMyToday(c *gin.Context) {
	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subs, err := h.service.ListMyToday(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *Handler) AssignSubstitute(c *gin.Context) {
	if c.GetString("user_id") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	var req AssignSubstituteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignSubstitute(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "substitute assigned"})
}

func (h *Handler) Confirm(c *gin.Context) {
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
	if err := h.service.ConfirmSubstitution(c.Request.Context(), id, teacherID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "substitution confirmed"})
}
