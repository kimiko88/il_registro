package documents

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"registro-backend/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type Handler struct {
	service  Service
	uploader upload.StorageUploader
}

func NewHandler(s Service, u ...upload.StorageUploader) *Handler {
	var upl upload.StorageUploader
	if len(u) > 0 && u[0] != nil {
		upl = u[0]
	}
	return &Handler{service: s, uploader: upl}
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
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if userID == "" || role == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.CreateDocument(c.Request.Context(), role, userID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// UploadFile handles multipart file uploads with magic-byte MIME validation and Supabase Storage upload.
// POST /documents/upload
// Form fields: file (required), document_id (optional, to associate with an existing document)
func (h *Handler) UploadFile(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || role == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

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

	// Rewind file reader after validation
	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	// Generate a unique safe storage key using UUID and pure sanitized extension
	cleanFilename := filepath.Base(header.Filename)
	ext := strings.ToLower(filepath.Ext(cleanFilename))
	allowedExts := map[string]bool{
		".pdf": true, ".docx": true, ".doc": true, ".xlsx": true, ".xls": true,
		".txt": true, ".jpg": true, ".jpeg": true, ".png": true, ".csv": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "estensione file non consentita"})
		return
	}
	storagePath := uuid.New().String() + ext

	var publicURL string
	if h.uploader != nil {
		detectedMIME, errMIME := upload.DetectMIME(file)
		if errMIME != nil {
			detectedMIME = header.Header.Get("Content-Type")
		}
		// Rewind file reader again before uploading payload
		if seeker, ok := file.(io.Seeker); ok {
			_, _ = seeker.Seek(0, io.SeekStart)
		}
		publicURL, err = h.uploader.UploadFile(c.Request.Context(), storagePath, detectedMIME, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore caricamento storage"})
			return
		}
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "storage provider not configured"})
		return
	}

	docID := c.PostForm("document_id")
	if docID != "" {
		if err := h.service.AttachFile(c.Request.Context(), role, schoolID, docID, publicURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore associazione documento: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "File caricato con successo",
		"filename":    header.Filename,
		"size_bytes":  header.Size,
		"url":         publicURL,
		"uploaded_by": userID,
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
		if strings.Contains(err.Error(), "not found") || err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateDocument(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || role == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateDocument(c.Request.Context(), role, schoolID, userID, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) ProcessWorkflow(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || role == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req WorkflowActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ProcessWorkflow(c.Request.Context(), role, schoolID, userID, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "processed"})
}

func (h *Handler) SignDocument(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions to sign document"})
		return
	}

	var req SignDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SignDocument(c.Request.Context(), role, schoolID, userID, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "signed"})
}

func (h *Handler) GetInbox(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Teachers submit documents; they do NOT have access to the global school inbox
	// (which shows documents from all staff awaiting review). Teachers access their
	// own documents via GET /documents/my.
	if role != "secretary" && role != "principal" && role != "vice_principal" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: insufficient permissions"})
		return
	}
	res, err := h.service.GetInbox(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if res == nil {
		res = []DocumentListResponse{}
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetReviewQueue(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Review queue shows pending documents from all staff — teachers should not
	// see their colleagues' submitted documents. Only secretarial/managerial staff allowed.
	if role != "secretary" && role != "principal" && role != "vice_principal" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: insufficient permissions"})
		return
	}
	res, err := h.service.GetReviewQueue(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if res == nil {
		res = []DocumentListResponse{}
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if schoolID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: insufficient permissions to manage templates"})
		return
	}

	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if schoolID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: insufficient permissions to manage templates"})
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
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if schoolID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: insufficient permissions to manage templates"})
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
	format := strings.ToLower(c.Query("format"))
	if format == "" {
		format = "pdf"
	}
	allowedFormats := map[string]bool{"pdf": true, "docx": true, "txt": true, "csv": true, "xlsx": true}
	if !allowedFormats[format] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato di esportazione non valido"})
		return
	}
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
	filename := fmt.Sprintf("document_%s.%s", id, format)
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, contentType, data)
}

func (h *Handler) ListDocuments(c *gin.Context) {
	docType := c.Query("type")
	var dt *DocType
	if docType != "" {
		allowedDocTypes := map[string]bool{
			"circolare": true, "modulo": true, "verbale": true,
			"programmazione": true, "pdp": true, "pei": true, "altro": true,
		}
		if !allowedDocTypes[strings.ToLower(docType)] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tipo di documento non valido"})
			return
		}
		val := DocType(docType)
		dt = &val
	}
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || role == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
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
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if role == "" || userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.DeleteDocument(c.Request.Context(), role, schoolID, userID, id); err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
func (h *Handler) ListTemplates(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	res, err := h.service.ListTemplates(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetDocumentVersions(c *gin.Context) {
	id := c.Param("id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if role == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	versions, err := h.service.GetDocumentVersions(c.Request.Context(), role, schoolID, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}
