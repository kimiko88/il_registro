package communications

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"registro-backend/pkg/upload"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service  *Service
	uploader upload.StorageUploader
}

func NewHandler(s *Service, u ...upload.StorageUploader) *Handler {
	h := &Handler{service: s}
	if len(u) > 0 {
		h.uploader = u[0]
	}
	return h
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/communications")
	{
		g.GET("", h.List)
		g.GET("/bacheca", h.ListBacheca)
		g.GET("/circolari", h.ListCircolari)
		g.GET("/unread-count", h.GetUnreadCount)
		g.GET("/:id", h.GetByID)
		g.PUT("/:id", h.Update)
		g.POST("", h.Send)
		g.POST("/upload", h.UploadAttachment)
		g.GET("/:id/attachment", h.DownloadAttachment)
		g.DELETE("/:id", h.Delete)
		g.POST("/:id/sign", h.Sign)
		g.POST("/:id/read", h.MarkAsRead)
		g.POST("/:id/ack", h.Ack)
		g.GET("/:id/signatures", h.GetSignatures)
		g.GET("/:id/signature-report", h.GetSignatureReport)
		g.GET("/:id/unread-users", h.GetUnreadUsers)
	}
}

func (h *Handler) List(c *gin.Context) {
	uid := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	msgs, err := h.service.ListMessages(c.Request.Context(), uid, schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.decorateMessages(msgs)
	c.JSON(http.StatusOK, msgs)
}

func (h *Handler) ListBacheca(c *gin.Context) {
	uid := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	msgs, err := h.service.ListBacheca(c.Request.Context(), schoolID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.decorateMessages(msgs)
	c.JSON(http.StatusOK, msgs)
}

func (h *Handler) ListCircolari(c *gin.Context) {
	uid := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	year := c.Query("year")
	msgs, err := h.service.ListCircolari(c.Request.Context(), schoolID, uid, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if msgs == nil {
		msgs = []*Message{}
	}
	h.decorateMessages(msgs)
	c.JSON(http.StatusOK, msgs)
}

func (h *Handler) Send(c *gin.Context) {
	uid := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64*1024)
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "il corpo del messaggio supera il limite massimo consentito (64 KB) o formato JSON non valido"})
		return
	}
	msg, err := h.service.SendMessage(c.Request.Context(), role, schoolID, uid, req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") || strings.HasPrefix(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.decorateMessageAttachment(msg)
	c.JSON(http.StatusCreated, msg)
}

// Delete verifies that the caller is the author or authorized admin of the message before deleting.
func (h *Handler) Delete(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	id := c.Param("id")
	if err := h.service.DeleteMessage(c.Request.Context(), uid, role, schoolID, id); err != nil {
		if err.Error() == "forbidden" || strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) Sign(c *gin.Context) {
	id := c.Param("id")
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ipAddress := c.ClientIP()
	if err := h.service.SignMessageWithIP(c.Request.Context(), id, uid, ipAddress); err != nil {
		if errors.Is(err, ErrNotFound) || err.Error() == "communication not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "communication not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "signed"})
}

// GetSignatures returns the list of users who signed a message.
func (h *Handler) GetSignatures(c *gin.Context) {
	uid := c.GetString("user_id")
	role := c.GetString("role")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id := c.Param("id")
	names, err := h.service.GetMessageSignatures(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, names)
}

// GetSignatureReport returns the full signature report for a message.
func (h *Handler) GetSignatureReport(c *gin.Context) {
	uid := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	role := c.GetString("role")
	report, err := h.service.GetSignatureReport(c.Request.Context(), uid, role, schoolID, id)
	if err != nil {
		if err.Error() == "forbidden: signature reports are restricted to staff" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

// GetByID returns a single message by ID.
func (h *Handler) GetByID(c *gin.Context) {
	uid := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	msg, err := h.service.GetMessageByID(c.Request.Context(), uid, schoolID, role, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) || err.Error() == "communication not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "communication not found"})
			return
		}
		if strings.HasPrefix(err.Error(), "unauthorized") || err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.decorateMessageAttachment(msg)
	c.JSON(http.StatusOK, msg)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	var req struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateMessage(c.Request.Context(), uid, role, schoolID, id, req.Subject, req.Body); err != nil {
		if errors.Is(err, ErrNotFound) || err.Error() == "communication not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "communication not found"})
			return
		}
		if strings.HasPrefix(err.Error(), "unauthorized") || err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	uid := c.GetString("user_id")
	ipAddress := c.ClientIP()
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.MarkAsRead(c.Request.Context(), id, uid, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "read"})
}

// GetUnreadUsers returns users who have not read a message.
// BUG FIX: aggiunto controllo auth e role (solo staff).
func (h *Handler) GetUnreadUsers(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id := c.Param("id")
	unread, err := h.service.GetUnreadUsers(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, unread)
}

func (h *Handler) GetUnreadCount(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	count, err := h.service.GetUnreadCount(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *Handler) UploadAttachment(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "identità utente non valida"})
		return
	}
	if _, err := uuid.Parse(schoolID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "scuola non valida"})
		return
	}
	if _, ok := h.uploader.(upload.PrivateStorageDownloader); !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "private communication storage is not configured"})
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, upload.MaxUploadSize+(1<<10))
	if err := c.Request.ParseMultipartForm(upload.MaxUploadSize + (1 << 10)); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "request body too large") {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "allegato troppo grande"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "form allegato non valido"})
		}
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "allegato mancante"})
		return
	}
	defer func() { _ = file.Close() }()

	sanitizedBytes, detectedMIME, err := upload.ValidateAndSanitize(file, header)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ext := strings.ToLower(filepath.Ext(filepath.Base(header.Filename)))
	allowedExtensions := map[string]map[string]bool{
		".pdf":  {"application/pdf": true},
		".jpg":  {"image/jpeg": true},
		".jpeg": {"image/jpeg": true},
		".png":  {"image/png": true},
		".doc":  {"application/msword": true},
		".docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true},
		".xls":  {"application/vnd.ms-excel": true},
		".xlsx": {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true},
		".txt":  {"text/plain": true},
	}
	if !isAllowedMIMEType(detectedMIME) || !allowedExtensions[ext][detectedMIME] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo o estensione di allegato non consentiti"})
		return
	}

	storagePath := fmt.Sprintf("communications/%s/%s/%s%s", schoolID, userID, uuid.New().String(), ext)
	storageKey, err := h.uploader.UploadFile(c.Request.Context(), storagePath, detectedMIME, bytes.NewReader(sanitizedBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "caricamento allegato non riuscito"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"attachment_url": storageKey})
}

func (h *Handler) decorateMessageAttachment(msg *Message) {
	if msg == nil || msg.AttachmentURL == nil {
		return
	}
	parsed, err := url.Parse(*msg.AttachmentURL)
	if err != nil || parsed.Scheme != "supabase" {
		return
	}
	proxyURL := "/api/v1/communications/" + url.PathEscape(msg.ID) + "/attachment"
	msg.AttachmentURL = &proxyURL
}

func (h *Handler) decorateMessages(messages []*Message) {
	for _, msg := range messages {
		h.decorateMessageAttachment(msg)
	}
}

// DownloadAttachment only serves an attachment after the current user is verified
// as the message sender or a recipient by GetMessageByID.
func (h *Handler) DownloadAttachment(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if userID == "" || schoolID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	msg, err := h.service.GetMessageByID(c.Request.Context(), userID, schoolID, role, c.Param("id"))
	if err != nil || msg == nil || msg.AttachmentURL == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "allegato non trovato o non accessibile"})
		return
	}
	parsed, err := url.Parse(*msg.AttachmentURL)
	if err != nil || parsed.Scheme != "supabase" {
		c.JSON(http.StatusNotFound, gin.H{"error": "allegato privato non disponibile"})
		return
	}
	downloader, ok := h.uploader.(upload.PrivateStorageDownloader)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "private communication storage is not configured"})
		return
	}
	reader, contentType, contentLength, err := downloader.DownloadFile(c.Request.Context(), *msg.AttachmentURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "download allegato non riuscito"})
		return
	}
	defer func() { _ = reader.Close() }()
	if contentLength > upload.MaxUploadSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "allegato troppo grande"})
		return
	}
	data, err := io.ReadAll(io.LimitReader(reader, upload.MaxUploadSize+1))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "lettura allegato non riuscita"})
		return
	}
	if int64(len(data)) > upload.MaxUploadSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "allegato troppo grande"})
		return
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	filename := filepath.Base(strings.TrimPrefix(parsed.Path, "/"))
	c.Header("Cache-Control", "private, no-store, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, contentType, data)
}

// isAllowedMIMEType validates that the MIME type is in the allowed allowlist.
// IMPORTANT: contentType MUST be derived from magic-byte detection (upload.DetectMIME),
// NOT from the client-supplied Content-Type header, to prevent MIME-spoofing attacks.
func isAllowedMIMEType(contentType string) bool {
	contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	allowed := map[string]bool{
		"application/pdf":    true,
		"image/jpeg":         true,
		"image/png":          true,
		"image/gif":          true,
		"image/webp":         true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
		"text/plain": true,
	}
	return allowed[contentType]
}

func (h *Handler) Ack(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	commID := c.Param("id")
	if commID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "communication id required"})
		return
	}

	if err := h.service.AckMessage(c.Request.Context(), commID, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Presa d'atto registrata con successo"})
}
