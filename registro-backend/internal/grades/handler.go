package grades

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

func (h *Handler) GetGrades(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get grades"})
}

func (h *Handler) AddGrade(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "add grade"})
}
