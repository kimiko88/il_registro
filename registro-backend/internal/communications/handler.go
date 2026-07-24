package communications

import (
	"errors"
	"fmt"
	"net/http"
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
		g.DELETE("/:id", h.Delete)
		g.POST("/:id/sign", h.Sign)
		g.POST("/:id/read", h.MarkAsRead)
		g.GET("/:id/signatures", h.GetSignatures)
		g.GET("/:id/signature-report", h.GetSignatureReport)
		g.GET("/:id/unread-users", h.GetUnreadUsers)
	}
}

func (h *Handler) List(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	msgs, err := h.service.ListMessages(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	c.JSON(http.StatusOK, msgs)
}

func (h *Handler) Send(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Body) > 64*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "il corpo del messaggio supera il limite massimo consentito (64 KB)"})
		return
	}
	msg, err := h.service.SendMessage(c.Request.Context(), uid, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	role := c.GetString("role")
	report, err := h.service.GetSignatureReport(c.Request.Context(), role, id)
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
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	msg, err := h.service.GetMessageByID(c.Request.Context(), uid, role, id)
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
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "allegato mancante"})
		return
	}
	defer file.Close()

	// Max size check: 10 MB
	const maxFileSize = 10 * 1024 * 1024
	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "la dimensione del file supera il limite massimo consentito (10 MB)"})
		return
	}

	// Extension & MIME check
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".pdf": true, ".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true, ".doc": true, ".docx": true,
		".xls": true, ".xlsx": true, ".txt": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "estensione file non consentita"})
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType != "" && !isValidMIME(contentType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo MIME non consentito"})
		return
	}

	var publicURL string
	if h.uploader != nil {
		storagePath := fmt.Sprintf("communications/%s%s", uuid.New().String(), ext)
		publicURL, err = h.uploader.UploadFile(c.Request.Context(), storagePath, contentType, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "storage provider not configured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attachment_url": publicURL})
}

func isValidMIME(contentType string) bool {
	contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	allowed := map[string]bool{
		"application/pdf":                                                        true,
		"image/jpeg":                                                             true,
		"image/png":                                                              true,
		"image/gif":                                                              true,
		"image/webp":                                                             true,
		"application/msword":                                                     true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel":                                               true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":      true,
		"text/plain":                                                             true,
		"application/octet-stream":                                               true,
	}
	return allowed[contentType]
}
