package documents

import (
	"net/http"

	"registro-backend/pkg/upload"

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
	docs.POST("/upload", h.UploadFile) // multipart file upload with magic-byte validation
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

// UploadFile handles multipart file uploads with magic-byte MIME validation.
// POST /documents/upload
// Form fields: file (required), document_id (optional, to associate with an existing document)
func (h *Handler) UploadFile(c *gin.Context) {
	// Parse multipart form (limits to MaxUploadSize + small overhead for form fields)
	if err := c.Request.ParseMultipartForm(upload.MaxUploadSize + (1 << 10)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossibile leggere il form: " + err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'file' mancante nel form"})
		return
	}
	defer file.Close()

	// Validate size + real MIME type via magic bytes
	if err := upload.ValidateUpload(file, header); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// At this point the file is validated and the read position is reset to 0.
	// TODO: persist the file to storage (S3, local FS, etc.) and
	// associate it with a document record via document_id query param.
	// For now, we return a success response with metadata.
	c.JSON(http.StatusOK, gin.H{
		"message":   "File caricato e validato con successo",
		"filename":  header.Filename,
		"size_bytes": header.Size,
	})
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
