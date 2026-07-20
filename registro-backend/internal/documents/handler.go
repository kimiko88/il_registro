package documents

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"registro-backend/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	docs.GET("/:id/versions", h.GetDocumentVersions)
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
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossibile creare la cartella di destinazione: " + err.Error()})
		return
	}

	// Generate a unique safe filename
	destPath := filepath.Join(uploadDir, uuid.New().String()+"-"+filepath.Base(header.Filename))
	out, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossibile creare il file di destinazione: " + err.Error()})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il salvataggio del file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "File caricato e salvato con successo",
		"filename":     header.Filename,
		"size_bytes":    header.Size,
		"storage_path": destPath,
	})
}

func (h *Handler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if role == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.service.GetDocument(c.Request.Context(), role, schoolID, id)
	if err != nil {
		if err.Error() == "unauthorized" || err.Error() == "unauthorized: cannot access documents of another school" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
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
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.CreateTemplate(c.Request.Context(), schoolID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "template created"})
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateTemplate(c.Request.Context(), schoolID, id, req); err != nil {
		if err.Error() == "unauthorized: template belongs to another school" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "template updated"})
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.DeleteTemplate(c.Request.Context(), schoolID, id); err != nil {
		if err.Error() == "unauthorized: template belongs to another school" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "template deleted"})
}

func (h *Handler) ExportDocument(c *gin.Context) {
	id := c.Param("id")
	format := c.Query("format") // pdf, docx
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if role == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	data, contentType, err := h.service.ExportDocument(c.Request.Context(), role, schoolID, id, format)
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

func (h *Handler) GetDocumentVersions(c *gin.Context) {
	id := c.Param("id")
	versions, err := h.service.GetDocumentVersions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}
