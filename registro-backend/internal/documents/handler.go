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
	docs.GET("", h.ListDocuments)
	docs.GET("/:id", h.GetDocument)
	docs.PATCH("/:id", h.UpdateDocument)
	docs.DELETE("/:id", h.DeleteDocument)
	docs.POST("/:id/workflow", h.ProcessWorkflow) // Submit

	// Templates
	docs.POST("/template", h.CreateTemplate)
	docs.GET("/template", h.ListTemplates)
	docs.PATCH("/template/:id", h.UpdateTemplate)
	docs.DELETE("/template/:id", h.DeleteTemplate)

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
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	res, err := h.service.CreateDocument(c.Request.Context(), role, userID, schoolID, req)
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
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if err := h.service.UpdateDocument(c.Request.Context(), role, userID, id, req); err != nil {
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
	userID := c.GetString("user_id")
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
	userID := c.GetString("user_id")
	if err := h.service.SignDocument(c.Request.Context(), userID, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "signed"})
}

func (h *Handler) GetInbox(c *gin.Context) {
	schoolID := c.GetString("school_id")
	res, err := h.service.GetInbox(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetReviewQueue(c *gin.Context) {
	schoolID := c.GetString("school_id")
	res, err := h.service.GetReviewQueue(c.Request.Context(), schoolID)
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
	schoolID := c.GetString("school_id")
	if err := h.service.CreateTemplate(c.Request.Context(), schoolID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "template created"})
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateTemplate(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "template updated"})
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTemplate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "template deleted"})
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
func (h *Handler) ListDocuments(c *gin.Context) {
	docType := c.Query("type")
	var dt *DocType
	if docType != "" {
		val := DocType(docType)
		dt = &val
	}
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	res, err := h.service.ListDocuments(c.Request.Context(), role, schoolID, dt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	role := c.GetString("role")
	if err := h.service.DeleteDocument(c.Request.Context(), role, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
func (h *Handler) ListTemplates(c *gin.Context) {
	schoolID := c.GetString("school_id")
	res, err := h.service.ListTemplates(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
