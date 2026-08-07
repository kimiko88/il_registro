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
	r.GET("/search/global", h.Search)
}

var validFilterTypes = map[string]bool{
	"":               true,
	"all":            true,
	"students":       true,
	"teachers":       true,
	"classes":        true,
	"communications": true,
	"documents":      true,
	"notes":          true,
}

func (h *Handler) Search(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")

	q := c.Query("q")
	filterType := c.Query("type")

	if q == "" || len(q) < 2 {
		c.JSON(http.StatusOK, gin.H{
			"query":   q,
			"total":   0,
			"results": []SearchResultItem{},
		})
		return
	}

	if userID == "" {
		c.JSON(http.StatusOK, gin.H{
			"query":   q,
			"total":   0,
			"results": []SearchResultItem{},
		})
		return
	}

	if len(q) > 200 {
		q = q[:200]
	}

	if !validFilterTypes[filterType] {
		filterType = "all"
	}

	res, err := h.service.Search(c.Request.Context(), schoolID, q, filterType)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"query":   q,
			"total":   0,
			"results": []SearchResultItem{},
		})
		return
	}
	c.JSON(http.StatusOK, res)
}
