package verbali

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
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
	}
}

func (h *Handler) CreateMeeting(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" || (role != "teacher" && role != "admin" && role != "superadmin") {
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
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	ipAddress := c.ClientIP()

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	id := c.Param("id")
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
