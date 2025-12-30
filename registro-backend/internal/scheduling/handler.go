package scheduling

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetSchedule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get schedule"})
}

func (h *Handler) AddSchedule(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "add schedule"})
}
