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

var validFilterTypes = map[string]bool{
	"":              true,
	"all":          true,
	"students":     true,
	"teachers":     true,
	"classes":      true,
	"communications": true,
	"documents":    true,
	"notes":        true,
}

func (h *Handler) Search(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")

	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	q := c.Query("q")
	filterType := c.Query("type")

	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q query parameter is required"})
		return
	}

	if !validFilterTypes[filterType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter type parameter"})
		return
	}

	res, err := h.service.Search(c.Request.Context(), schoolID, q, filterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
