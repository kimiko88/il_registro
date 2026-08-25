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
	r.GET("/search/global", h.GlobalSearchEndpoint)
}

var validFilterTypes = map[string]bool{
	"":               true,
	"all":            true,
	"users":          true,
	"students":       true,
	"teachers":       true,
	"classes":        true,
	"communications": true,
	"lessons":        true,
	"documents":      true,
	"notes":          true,
}

// GlobalSearchEndpoint allows cross-tenant search ONLY for superadmins.
func (h *Handler) GlobalSearchEndpoint(c *gin.Context) {
	role := c.GetString("role")
	if role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "global search requires superadmin role"})
		return
	}
	h.Search(c)
}

func (h *Handler) Search(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if role != "superadmin" && schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required for search"})
		return
	}

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

	if len(q) > 200 {
		q = q[:200]
	}

	if !validFilterTypes[filterType] {
		filterType = "all"
	}

	res, err := h.service.Search(c.Request.Context(), role, schoolID, q, filterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	c.JSON(http.StatusOK, res)
}
