package textbooks

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

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	textbooks := rg.Group("/textbooks")
	{
		textbooks.POST("", h.Create)
		textbooks.GET("", h.List)
		textbooks.DELETE("/:id", h.Delete)
		textbooks.GET("/class/:classId", h.ListByClass)
		textbooks.POST("/class/:classId", h.AssignToClass)
		textbooks.DELETE("/class/assignment/:id", h.RemoveFromClass)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateTextbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID := c.GetString("school_id")
	if err := h.service.CreateTextbook(c.Request.Context(), schoolID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) List(c *gin.Context) {
	schoolID := c.Query("school_id")
	if schoolID == "" {
		schoolID = c.GetString("school_id")
	}
	res, err := h.service.ListTextbooks(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.DeleteTextbook(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListByClass(c *gin.Context) {
	res, err := h.service.ListByClass(c.Request.Context(), c.Param("classId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) AssignToClass(c *gin.Context) {
	var req AssignTextbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AssignToClass(c.Request.Context(), c.Param("classId"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) RemoveFromClass(c *gin.Context) {
	if err := h.service.RemoveFromClass(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
