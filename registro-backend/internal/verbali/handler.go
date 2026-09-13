package verbali

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"registro-backend/internal/pdfworker"
	"registro-backend/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service         *Service
	pdfWorkerClient *pdfworker.Client
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) SetPdfWorkerClient(client *pdfworker.Client) {
	h.pdfWorkerClient = client
}

func isAllowedVerbaliRole(role string) bool {
	switch strings.ToLower(role) {
	case "teacher", "coordinator", "admin", "superadmin", "secretary", "principal", "vice_principal", "docente",
		"dsga", "collaboratore_ds", "assistente_amministrativo":
		return true
	default:
		return false
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	v := r.Group("/verbali")
	{
		// Meetings
		v.POST("/meetings", h.CreateMeeting)
		v.GET("/meetings", h.ListMeetings)

		// Verbali Templates (Dirigente & School-wide)
		v.GET("/templates", h.ListTemplates)
		v.POST("/templates", h.CreateTemplate)
		v.GET("/templates/:id", h.GetTemplate)
		v.PUT("/templates/:id", h.UpdateTemplate)
		v.DELETE("/templates/:id", h.DeleteTemplate)

		// Verbali Lifecycle
		v.GET("", h.ListAllVerbali)
		v.POST("", h.CreateVerbale)
		v.GET("/:id", h.GetVerbale)
		v.PUT("/:id", h.UpdateVerbale)
		v.DELETE("/:id", h.DeleteVerbale)
		v.GET("/meeting/:meetingId", h.ListVerbali)
		v.POST("/:id/sign", h.SignVerbale)
		v.GET("/:id/signatures", h.GetSignatures)
		v.GET("/:id/pdf", h.ExportPDF)
		v.POST("/:id/async-pdf", h.EnqueueAsyncVerbalePdf)
		v.GET("/pdf-jobs/:job_id", h.GetPdfJobStatus)
	}
}

func (h *Handler) getSchoolID(c *gin.Context, userID string) string {
	schoolID := c.GetString("school_id")
	if schoolID == "" && h.service != nil {
		schoolID = h.service.ResolveSchoolID(c.Request.Context(), userID)
	}
	return schoolID
}

func (h *Handler) CreateMeeting(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := h.getSchoolID(c, userID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := h.service.CreateMeeting(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *Handler) ListMeetings(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	schoolID := h.getSchoolID(c, userID)
	classID := c.Query("class_id")

	meetings, err := h.service.ListMeetings(c.Request.Context(), schoolID, classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if meetings == nil {
		meetings = []*CouncilMeeting{}
	}
	c.JSON(http.StatusOK, meetings)
}

func (h *Handler) CreateVerbale(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateVerbaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, err := h.service.CreateVerbale(c.Request.Context(), userID, role, req)
	if err != nil {
		if errors.Is(err, ErrOnlyCoordinatorOrSecretaryCanEdit) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "code": "COORDINATOR_OR_SECRETARY_ONLY"})
			return
		}
		if errors.Is(err, ErrUnauthorized) || err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "meeting not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *Handler) UpdateVerbale(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req UpdateVerbaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, err := h.service.UpdateVerbale(c.Request.Context(), userID, role, id, req)
	if err != nil {
		if errors.Is(err, ErrVerbaleLocked) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "code": "VERBALE_LOCKED"})
			return
		}
		if errors.Is(err, ErrOnlyCoordinatorOrSecretaryCanEdit) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "code": "COORDINATOR_OR_SECRETARY_ONLY"})
			return
		}
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "verbale not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) DeleteVerbale(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	err := h.service.DeleteVerbale(c.Request.Context(), userID, role, id)
	if err != nil {
		if errors.Is(err, ErrVerbaleLocked) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "code": "VERBALE_LOCKED"})
			return
		}
		if errors.Is(err, ErrOnlyCoordinatorOrSecretaryCanEdit) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "code": "COORDINATOR_OR_SECRETARY_ONLY"})
			return
		}
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "verbale not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "verbale deleted"})
}

func (h *Handler) GetVerbale(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	v, err := h.service.GetVerbale(c.Request.Context(), id, userID, role)
	if err != nil {
		if errors.Is(err, ErrDraftHiddenFromPrincipal) {
			c.JSON(http.StatusNotFound, gin.H{"error": "verbale in lavorazione: non ancora visibile alla Dirigenza", "code": "DRAFT_HIDDEN"})
			return
		}
		if errors.Is(err, ErrUnauthorized) || err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrNotFound) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "verbale not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) ListVerbali(c *gin.Context) {
	meetingID := c.Param("meetingId")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if meetingID == "" || meetingID == "undefined" || meetingID == "null" {
		c.JSON(http.StatusOK, []*MeetingVerbale{})
		return
	}
	if _, err := uuid.Parse(meetingID); err != nil {
		c.JSON(http.StatusOK, []*MeetingVerbale{})
		return
	}

	list, err := h.service.ListVerbali(c.Request.Context(), meetingID, userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []*MeetingVerbale{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ListAllVerbali(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := h.getSchoolID(c, userID)
	classID := c.Query("class_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	list, err := h.service.ListAllVerbali(c.Request.Context(), schoolID, classID, userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []*MeetingVerbale{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) SignVerbale(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	ipAddress := c.ClientIP()

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.service.SignVerbale(c.Request.Context(), id, userID, ipAddress); err != nil {
		if errors.Is(err, ErrUnauthorized) || err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "verbale firmato con successo con tracciamento IP e archiviato"})
}

func (h *Handler) GetSignatures(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")

	sigs, err := h.service.GetSignatures(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sigs == nil {
		sigs = []VerbaleSignature{}
	}
	c.JSON(http.StatusOK, sigs)
}

func (h *Handler) ExportPDF(c *gin.Context) {
	if h.pdfWorkerClient != nil && c.Query("sync") != "true" {
		h.EnqueueAsyncVerbalePdf(c)
		return
	}

	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	verbale, err := h.service.GetVerbale(c.Request.Context(), id, userID, role)
	if err != nil {
		if errors.Is(err, ErrDraftHiddenFromPrincipal) {
			c.JSON(http.StatusNotFound, gin.H{"error": "verbale in bozza non esportabile per la Dirigenza", "code": "DRAFT_HIDDEN"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "verbale not found"})
		return
	}

	// Fetch signatures and meeting info for the PDF
	sigs, _ := h.service.GetSignatures(c.Request.Context(), id)
	meeting, _ := h.service.GetMeeting(c.Request.Context(), verbale.MeetingID)

	pdfBytes, err := GenerateVerbale(verbale, meeting, sigs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PDF generation failed: " + err.Error()})
		return
	}

	dateStr := verbale.CreatedAt.Format("20060102")
	shortID := verbale.MeetingID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	filename := "verbale_" + shortID + "_" + dateStr + ".pdf"
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *Handler) EnqueueAsyncVerbalePdf(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")

	if h.pdfWorkerClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "async pdf worker queue not configured"})
		return
	}

	jobID := uuid.New().String()
	payload := pdfworker.VerbalePdfPayload{
		JobID:       jobID,
		VerbaleID:   id,
		RequestedBy: userID,
	}

	jobStatus, err := h.pdfWorkerClient.EnqueueVerbalePdf(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to enqueue async verbale pdf job: %v", err)})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "Verbale PDF generation job enqueued successfully",
		"job_id":     jobStatus.JobID,
		"status":     jobStatus.Status,
		"status_url": fmt.Sprintf("/api/v1/verbali/pdf-jobs/%s", jobStatus.JobID),
	})
}

func (h *Handler) GetPdfJobStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	jobID := c.Param("job_id")
	if h.pdfWorkerClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "async pdf worker queue not configured"})
		return
	}

	status, err := h.pdfWorkerClient.GetJobStatus(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found or expired"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// ----------------- Template Handlers (Dirigente Scolastica & Admin) -----------------

func (h *Handler) ListTemplates(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := h.getSchoolID(c, userID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	meetingType := c.Query("meeting_type")
	templates, err := h.service.ListTemplates(c.Request.Context(), schoolID, meetingType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if templates == nil {
		templates = []*MeetingVerbaleTemplate{}
	}
	c.JSON(http.StatusOK, templates)
}

func (h *Handler) GetTemplate(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	t, err := h.service.GetTemplateByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := h.getSchoolID(c, userID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isAllowedVerbaliRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := h.service.CreateTemplate(c.Request.Context(), userID, role, schoolID, req)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "solo la Dirigenza, la DSGA o gli amministratori possono creare modelli di verbale"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	role := c.GetString("role")
	userID := c.GetString("user_id")
	schoolID := h.getSchoolID(c, userID)

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := h.service.UpdateTemplate(c.Request.Context(), role, schoolID, id, req)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "solo la Dirigenza, la DSGA o gli amministratori possono modificare modelli di verbale"})
			return
		}
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	role := c.GetString("role")
	userID := c.GetString("user_id")
	schoolID := h.getSchoolID(c, userID)

	err := h.service.DeleteTemplate(c.Request.Context(), role, schoolID, id)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "solo la Dirigenza, la DSGA o gli amministratori possono eliminare modelli di verbale"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "template deleted"})
}
