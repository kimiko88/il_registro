package scrutiny

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	scrutiny := rg.Group("/scrutiny")
	{
		scrutiny.GET("/matrix/:classId", h.GetMatrix)
		scrutiny.POST("/save", h.Save)
	}
}

func (h *Handler) GetMatrix(c *gin.Context) {
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	classID := c.Param("classId")
	
	matrix, err := h.service.GetMatrix(c.Request.Context(), classID, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matrix)
}

func (h *Handler) Save(c *gin.Context) {
	var req SaveScrutinyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	coordinatorID := c.GetString("user_id")
	if err := h.service.SaveScrutiny(c.Request.Context(), coordinatorID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
