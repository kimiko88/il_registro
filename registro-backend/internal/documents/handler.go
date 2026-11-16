package documents

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"io"
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
	docs.GET("/:id/file", h.DownloadFile)
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

	// The file must be attached to an authorized document; otherwise the upload
	// would create orphaned private objects that no user can retrieve.
	docID := strings.TrimSpace(c.PostForm("document_id"))
	if docID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document_id obbligatorio"})
		return
	}
	if _, err := h.service.GetDocument(c.Request.Context(), role, schoolID, docID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "documento non accessibile"})
		return
	}
	if _, ok := h.uploader.(upload.PrivateStorageDownloader); !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "private document storage is not configured"})
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
	defer func() { _ = file.Close() }()

	// Validate the real MIME and strip privacy-invasive image metadata before
	// sending bytes to storage. Do not trust Content-Type or the filename.
	sanitizedBytes, detectedMIME, err := upload.ValidateAndSanitize(file, header)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique safe storage key using UUID and pure sanitized extension
	cleanFilename := filepath.Base(header.Filename)
	ext := strings.ToLower(filepath.Ext(cleanFilename))
	allowedExts := map[string]map[string]bool{
		".pdf":  {"application/pdf": true},
		".docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true},
		".xlsx": {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true},
		".xls":  {"application/vnd.ms-excel": true},
		".txt":  {"text/plain": true},
		".csv":  {"text/csv": true, "text/plain": true},
		".jpg":  {"image/jpeg": true},
		".jpeg": {"image/jpeg": true},
		".png":  {"image/png": true},
	}
	if !allowedExts[ext][detectedMIME] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "estensione file non consentita"})
		return
	}
	storagePath := uuid.New().String() + ext

	storageKey, err := h.uploader.UploadFile(c.Request.Context(), storagePath, detectedMIME, bytes.NewReader(sanitizedBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore caricamento storage"})
		return
	}
	if err := h.service.AttachFile(c.Request.Context(), role, schoolID, docID, storageKey); err != nil {
		if deleter, ok := h.uploader.(upload.StorageObjectDeleter); ok {
			_ = deleter.DeleteFile(c.Request.Context(), storageKey)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore associazione documento"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "File caricato con successo",
		"filename":    header.Filename,
		"size_bytes":  len(sanitizedBytes),
		"url":         "/api/v1/documents/" + url.PathEscape(docID) + "/file",
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


// DownloadFile authorizes document access before retrieving an object from private storage.
// No public Supabase URL or storage credential is sent to the browser.
func (h *Handler) DownloadFile(c *gin.Context) {
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if role == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	detail, err := h.service.GetDocument(c.Request.Context(), role, schoolID, c.Param("id"))
	if err != nil || detail == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "documento non trovato o non accessibile"})
		return
	}
	var storageKey string
	for _, line := range strings.Split(detail.Content, "\\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Attachment: supabase://") {
			storageKey = strings.TrimSpace(strings.TrimPrefix(line, "Attachment:"))
		}
	}
	if storageKey == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "allegato privato non disponibile"})
		return
	}
	downloader, ok := h.uploader.(upload.PrivateStorageDownloader)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "private document storage is not configured"})
		return
	}
	reader, contentType, contentLength, err := downloader.DownloadFile(c.Request.Context(), storageKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "download documento non riuscito"})
		return
	}
	defer func() { _ = reader.Close() }()
	if contentLength > upload.MaxUploadSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "documento troppo grande"})
		return
	}
	data, err := io.ReadAll(io.LimitReader(reader, upload.MaxUploadSize+1))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "lettura documento non riuscita"})
		return
	}
	if int64(len(data)) > upload.MaxUploadSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "documento troppo grande"})
		return
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	filename := "document-" + c.Param("id") + ".bin"
	if parsed, err := url.Parse(storageKey); err == nil {
		filename = filepath.Base(strings.TrimPrefix(parsed.Path, "/"))
	}
	c.Header("Cache-Control", "private, no-store, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, contentType, data)
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
	if role != "secretary" && role != "principal" && role != "vice_principal" && role != "admin" && role != "superadmin" && role != "dsga" && role != "assistente_amministrativo" && role != "collaboratore_ds" && role != "responsabile_gestione_documentale" && role != "responsabile_conservazione" && role != "assistente_protocollo" && role != "assistente_contabilita" && role != "assistente_alunni" {
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
	if role != "secretary" && role != "principal" && role != "vice_principal" && role != "admin" && role != "superadmin" && role != "dsga" && role != "assistente_amministrativo" && role != "collaboratore_ds" && role != "responsabile_gestione_documentale" && role != "responsabile_conservazione" && role != "assistente_protocollo" && role != "assistente_contabilita" && role != "assistente_alunni" {
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
	if role != "admin" && role != "superadmin" && role != "secretary" && role != "dsga" && role != "assistente_amministrativo" && role != "principal" && role != "vice_principal" && role != "responsabile_gestione_documentale" {
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
	if role != "admin" && role != "superadmin" && role != "secretary" && role != "dsga" && role != "assistente_amministrativo" && role != "principal" && role != "vice_principal" && role != "responsabile_gestione_documentale" {
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
	if role != "admin" && role != "superadmin" && role != "secretary" && role != "dsga" && role != "assistente_amministrativo" && role != "principal" && role != "vice_principal" && role != "responsabile_gestione_documentale" {
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
