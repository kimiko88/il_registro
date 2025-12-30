package documents

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

func (h *Handler) GetDocuments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get documents"})
}

func (h *Handler) UploadDocument(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "upload document"})
}
