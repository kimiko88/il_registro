package search

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
	r.GET("/search", h.Search)
}

func (h *Handler) Search(c *gin.Context) {
	q := c.Query("q")
	filterType := c.Query("type")
	schoolID := c.GetString("school_id")

	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q query parameter is required"})
		return
	}

	res, err := h.service.Search(c.Request.Context(), schoolID, q, filterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
