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
	case "teacher", "coordinator", "admin", "superadmin", "secretary", "principal", "vice_principal", "docente":
		return true
	default:
		return false
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	v := r.Group("/verbali")
	{
		v.POST("/meetings", h.CreateMeeting)
		v.GET("/meetings", h.ListMeetings)
		v.POST("", h.CreateVerbale)
		v.GET("/:id", h.GetVerbale)
		v.GET("/meeting/:meetingId", h.ListVerbali)
		v.POST("/:id/sign", h.SignVerbale)
		v.GET("/:id/signatures", h.GetSignatures)
		v.GET("/:id/pdf", h.ExportPDF)
		v.POST("/:id/async-pdf", h.EnqueueAsyncVerbalePdf)
		v.GET("/pdf-jobs/:job_id", h.GetPdfJobStatus)
	}
}

func (h *Handler) CreateMeeting(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
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

	schoolID := c.GetString("school_id")
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
		if errors.Is(err, ErrUnauthorized) || err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
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

	v, err := h.service.GetVerbale(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
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

	list, err := h.service.ListVerbali(c.Request.Context(), meetingID, userID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
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
	c.JSON(http.StatusOK, gin.H{"message": "verbale signed"})
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

	verbale, err := h.service.GetVerbale(c.Request.Context(), id, userID)
	if err != nil {
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
