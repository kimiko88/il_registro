package attendance

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

func (h *Handler) GetAttendance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get attendance"})
}

func (h *Handler) MarkAttendance(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "mark attendance"})
}
