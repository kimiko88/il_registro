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

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	docs := r.Group("/documents")

	// Teacher/Creator
	docs.POST("", h.CreateDocument)
	docs.GET("/:id", h.GetDocument)
	docs.PATCH("/:id", h.UpdateDocument)
	docs.POST("/:id/workflow", h.ProcessWorkflow) // Submit

	// Templates
	docs.POST("/template", h.CreateTemplate)

	// Secretary
	docs.GET("/inbox", h.GetInbox)

	// Director
	docs.GET("/review-queue", h.GetReviewQueue)
	docs.POST("/:id/sign", h.SignDocument)

	// Export
	docs.GET("/:id/export", h.ExportDocument)
}

func (h *Handler) CreateDocument(c *gin.Context) {
	var req CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	res, err := h.service.CreateDocument(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	res, err := h.service.GetDocument(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateDocument(c *gin.Context) {
	id := c.Param("id")
	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	if err := h.service.UpdateDocument(c.Request.Context(), userID, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) ProcessWorkflow(c *gin.Context) {
	id := c.Param("id")
	var req WorkflowActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	if err := h.service.ProcessWorkflow(c.Request.Context(), userID, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "processed"})
}

func (h *Handler) SignDocument(c *gin.Context) {
	id := c.Param("id")
	var req SignDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	if err := h.service.SignDocument(c.Request.Context(), userID, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "signed"})
}

func (h *Handler) GetInbox(c *gin.Context) {
	res, err := h.service.GetInbox(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetReviewQueue(c *gin.Context) {
	res, err := h.service.GetReviewQueue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateTemplate(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "template created"})
}

func (h *Handler) ExportDocument(c *gin.Context) {
	id := c.Param("id")
	format := c.Query("format") // pdf, docx
	data, contentType, err := h.service.ExportDocument(c.Request.Context(), id, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, contentType, data)
}
