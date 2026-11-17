package family_desk

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"registro-backend/internal/users"
	"registro-backend/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service  *Service
	userRepo users.Repository
	uploader upload.StorageUploader
}

func NewHandler(service *Service, userRepos ...users.Repository) *Handler {
	h := &Handler{service: service}
	if len(userRepos) > 0 {
		h.userRepo = userRepos[0]
	}
	return h
}

func (h *Handler) SetUploader(uploader upload.StorageUploader) {
	h.uploader = uploader
}

func isFamilyDeskStaffRole(role string) bool {
	switch role {
	case "admin", "superadmin", "secretary", "dsga", "principal", "vice_principal", "dpo":
		return true
	default:
		return false
	}
}

func (h *Handler) currentSchoolUser(c *gin.Context) (*users.User, bool) {
	if h.userRepo == nil {
		return nil, false
	}
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" || schoolID == "" {
		return nil, false
	}
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil || user == nil || user.SchoolID == nil || *user.SchoolID != schoolID {
		return nil, false
	}
	return user, true
}

func (h *Handler) authorizeParentStudent(c *gin.Context, studentID string) bool {
	user, ok := h.currentSchoolUser(c)
	if !ok || user.Role != "parent" || studentID == "" {
		return false
	}
	isGuardian, err := h.userRepo.IsGuardian(c.Request.Context(), user.ID, studentID)
	if err != nil || !isGuardian {
		return false
	}
	student, err := h.userRepo.GetByID(c.Request.Context(), studentID)
	return err == nil && student != nil && student.Role == "student" &&
		student.SchoolID != nil && *student.SchoolID == *user.SchoolID
}

func familyAttachmentKey(storageKey, schoolID, parentID, studentID string) (string, bool) {
	parsed, err := url.Parse(strings.TrimSpace(storageKey))
	if err != nil || parsed.Scheme != "supabase" || parsed.Host == "" ||
		parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" ||
		strings.ContainsAny(storageKey, "\\?#\r\n\x00") {
		return "", false
	}
	prefix := "/family-desk/" + schoolID + "/" + parentID + "/" + studentID + "/"
	if !strings.HasPrefix(parsed.Path, prefix) {
		return "", false
	}
	filename := strings.TrimPrefix(parsed.Path, prefix)
	if filename == "" || filename == "." || filename == ".." || strings.Contains(filename, "/") {
		return "", false
	}
	return filename, true
}

func (h *Handler) validateRequestAttachments(req *FamilyRequest) error {
	if len(req.AttachmentURLs) > 10 {
		return fmt.Errorf("non è possibile allegare più di 10 file")
	}
	for _, storageKey := range req.AttachmentURLs {
		if _, ok := familyAttachmentKey(storageKey, req.SchoolID, req.ParentID, req.StudentID); !ok {
			return fmt.Errorf("allegato non valido o non appartenente al genitore e allo studente")
		}
	}
	if raw, ok := req.FormData["id_card_file_url"]; ok {
		cardURL, ok := raw.(string)
		if !ok || cardURL == "" {
			return fmt.Errorf("id_card_file_url non valido")
		}
		if _, ok := familyAttachmentKey(cardURL, req.SchoolID, req.ParentID, req.StudentID); !ok {
			return fmt.Errorf("documento di identità non appartenente al genitore e allo studente")
		}
		found := false
		for _, attachment := range req.AttachmentURLs {
			if attachment == cardURL {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("il documento di identità deve essere incluso negli allegati della pratica")
		}
	}
	return nil
}

func familyAttachmentProxyURL(requestID, filename string) string {
	return "/api/v1/family-desk/requests/" + url.PathEscape(requestID) +
		"/attachments/" + url.PathEscape(filename)
}

func (h *Handler) decorateFamilyRequest(req *FamilyRequest) {
	if req == nil {
		return
	}
	for i, key := range req.AttachmentURLs {
		if filename, ok := familyAttachmentKey(key, req.SchoolID, req.ParentID, req.StudentID); ok {
			req.AttachmentURLs[i] = familyAttachmentProxyURL(req.ID, filename)
		}
	}
	if req.FormData != nil {
		if raw, ok := req.FormData["id_card_file_url"].(string); ok {
			if filename, valid := familyAttachmentKey(raw, req.SchoolID, req.ParentID, req.StudentID); valid {
				req.FormData["id_card_file_url"] = familyAttachmentProxyURL(req.ID, filename)
			}
		}
	}
}

func (h *Handler) authorizeRequest(c *gin.Context, req *FamilyRequest) bool {
	if req == nil {
		return false
	}
	user, ok := h.currentSchoolUser(c)
	if !ok || *user.SchoolID != req.SchoolID {
		return false
	}
	if user.Role == "parent" {
		return req.ParentID == user.ID && h.authorizeParentStudent(c, req.StudentID)
	}
	return isFamilyDeskStaffRole(user.Role)
}

func (h *Handler) UploadAttachment(c *gin.Context) {
	parent, ok := h.currentSchoolUser(c)
	if !ok || parent.Role != "parent" || !h.authorizeParentStudent(c, c.Param("student_id")) {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo il genitore verificato dello studente può caricare allegati"})
		return
	}
	if _, ok := h.uploader.(upload.PrivateStorageDownloader); !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "storage privato non configurato"})
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
	data, detectedMIME, err := upload.ValidateAndSanitize(file, header)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ext := strings.ToLower(filepath.Ext(filepath.Base(header.Filename)))
	allowed := map[string]map[string]bool{
		".pdf":  {"application/pdf": true},
		".jpg":  {"image/jpeg": true},
		".jpeg": {"image/jpeg": true},
		".png":  {"image/png": true},
		".docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true},
		".xlsx": {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true},
		".txt":  {"text/plain": true},
	}
	if !allowed[ext][detectedMIME] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo o estensione di allegato non consentiti"})
		return
	}
	storagePath := fmt.Sprintf("family-desk/%s/%s/%s/%s%s",
		*parent.SchoolID, parent.ID, c.Param("student_id"), uuid.New().String(), ext)
	storageKey, err := h.uploader.UploadFile(c.Request.Context(), storagePath, detectedMIME, bytes.NewReader(data))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "caricamento allegato non riuscito"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"attachment_url": storageKey})
}

func (h *Handler) SubmitRequest(c *gin.Context) {
	var reqBody struct {
		StudentID      string                 `json:"student_id" binding:"required"`
		RequestType    string                 `json:"request_type" binding:"required"`
		FormData       map[string]interface{} `json:"form_data" binding:"required"`
		AttachmentURLs []string               `json:"attachment_urls"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati non validi: " + err.Error()})
		return
	}
	parent, ok := h.currentSchoolUser(c)
	if !ok || parent.Role != "parent" || !h.authorizeParentStudent(c, reqBody.StudentID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "rapporto genitore-studente o scuola non verificati"})
		return
	}
	item := &FamilyRequest{
		SchoolID: *parent.SchoolID, StudentID: reqBody.StudentID, ParentID: parent.ID,
		RequestType: reqBody.RequestType, FormData: reqBody.FormData, AttachmentURLs: reqBody.AttachmentURLs,
	}
	if err := h.validateRequestAttachments(item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SubmitRequest(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.decorateFamilyRequest(item)
	c.JSON(http.StatusCreated, gin.H{"message": "Istanza presentata con successo", "request": item})
}

func (h *Handler) ListRequests(c *gin.Context) {
	user, ok := h.currentSchoolUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "identità o scuola non valida"})
		return
	}
	var parentFilter string
	if user.Role == "parent" {
		parentFilter = user.ID
	} else if isFamilyDeskStaffRole(user.Role) {
		parentFilter = strings.TrimSpace(c.Query("parent_id"))
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "ruolo non autorizzato"})
		return
	}
	list, err := h.service.ListRequests(c.Request.Context(), *user.SchoolID, parentFilter, c.Query("status"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile elencare le istanze"})
		return
	}
	for i := range list {
		h.decorateFamilyRequest(&list[i])
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *Handler) GetRequest(c *gin.Context) {
	req, err := h.service.GetRequest(c.Request.Context(), c.Param("id"))
	if err != nil || req == nil || !h.authorizeRequest(c, req) {
		c.JSON(http.StatusNotFound, gin.H{"error": "istanza non trovata o non accessibile"})
		return
	}
	h.decorateFamilyRequest(req)
	c.JSON(http.StatusOK, req)
}

func (h *Handler) ReviewRequest(c *gin.Context) {
	user, ok := h.currentSchoolUser(c)
	if !ok || !isFamilyDeskStaffRole(user.Role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo la segreteria o il dirigente possono revisionare le istanze"})
		return
	}
	id := c.Param("id")
	req, err := h.service.GetRequest(c.Request.Context(), id)
	if err != nil || req == nil || req.SchoolID != *user.SchoolID {
		c.JSON(http.StatusNotFound, gin.H{"error": "istanza non trovata o non accessibile"})
		return
	}
	var body struct {
		Status          string `json:"status" binding:"required"`
		RejectionReason string `json:"rejection_reason"`
		ProtocolNumber  string `json:"protocol_number"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati richiesta non validi"})
		return
	}
	if err := h.service.ReviewRequest(c.Request.Context(), id, body.Status, body.RejectionReason, body.ProtocolNumber, user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stato istanza aggiornato con successo"})
}

func (h *Handler) ListDelegates(c *gin.Context) {
	user, ok := h.currentSchoolUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "identità o scuola non valida"})
		return
	}
	studentID := strings.TrimSpace(c.Query("student_id"))
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id obbligatorio"})
		return
	}
	if user.Role == "parent" {
		if !h.authorizeParentStudent(c, studentID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "studente non associato al genitore"})
			return
		}
	} else if isFamilyDeskStaffRole(user.Role) {
		student, err := h.userRepo.GetByID(c.Request.Context(), studentID)
		if err != nil || student == nil || student.SchoolID == nil || *student.SchoolID != *user.SchoolID {
			c.JSON(http.StatusForbidden, gin.H{"error": "studente non appartenente alla scuola"})
			return
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "ruolo non autorizzato"})
		return
	}
	delegates, err := h.service.ListPermanentDelegates(c.Request.Context(), studentID, *user.SchoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile elencare le deleghe"})
		return
	}
	for i := range delegates {
		if delegates[i].IDCardFileURL != "" && delegates[i].ApprovedRequestID != nil {
			// The actual private key stays in the approved request; the request
			// attachment proxy performs the final parent/school authorization.
			if parsed, err := url.Parse(delegates[i].IDCardFileURL); err == nil && parsed.Scheme == "supabase" {
				filename := filepath.Base(parsed.Path)
				delegates[i].IDCardFileURL = familyAttachmentProxyURL(*delegates[i].ApprovedRequestID, filename)
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": delegates})
}

func (h *Handler) DownloadAttachment(c *gin.Context) {
	req, err := h.service.GetRequest(c.Request.Context(), c.Param("id"))
	if err != nil || req == nil || !h.authorizeRequest(c, req) {
		c.JSON(http.StatusNotFound, gin.H{"error": "istanza non trovata o non accessibile"})
		return
	}
	filename := c.Param("filename")
	var storageKey string
	for _, candidate := range req.AttachmentURLs {
		if candidateName, valid := familyAttachmentKey(candidate, req.SchoolID, req.ParentID, req.StudentID); valid && candidateName == filename {
			storageKey = candidate
			break
		}
	}
	if storageKey == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "allegato non trovato o non accessibile"})
		return
	}
	downloader, ok := h.uploader.(upload.PrivateStorageDownloader)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "storage privato non configurato"})
		return
	}
	reader, contentType, contentLength, err := downloader.DownloadFile(c.Request.Context(), storageKey)
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
	if err != nil || int64(len(data)) > upload.MaxUploadSize {
		c.JSON(http.StatusBadGateway, gin.H{"error": "lettura allegato non riuscita"})
		return
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Header("Cache-Control", "private, no-store, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, contentType, data)
}


func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/family-desk")
	{
		g.POST("/requests", h.SubmitRequest)
		g.GET("/requests", h.ListRequests)
		g.GET("/requests/:id", h.GetRequest)
		g.PUT("/requests/:id/review", h.ReviewRequest)
		g.GET("/requests/:id/attachments/:filename", h.DownloadAttachment)
		g.POST("/students/:student_id/attachments", h.UploadAttachment)
		g.GET("/delegates", h.ListDelegates)
	}
}
