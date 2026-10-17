package family_desk

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	fd := rg.Group("/family-desk")
	{
		fd.POST("/requests", h.SubmitRequest)
		fd.GET("/requests", h.ListRequests)
		fd.GET("/requests/:id", h.GetRequest)
		fd.PUT("/requests/:id/review", h.ReviewRequest)
		fd.GET("/delegates", h.ListDelegates)
	}
}

func (h *Handler) SubmitRequest(c *gin.Context) {
	var req struct {
		StudentID      string                 `json:"student_id" binding:"required"`
		RequestType    string                 `json:"request_type" binding:"required"`
		FormData       map[string]interface{} `json:"form_data" binding:"required"`
		AttachmentURLs []string               `json:"attachment_urls"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati non validi: " + err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}
	parentID := c.GetString("user_id")
	if parentID == "" {
		parentID = "00000000-0000-0000-0000-000000000002"
	}

	item := &FamilyRequest{
		SchoolID:       schoolID,
		StudentID:      req.StudentID,
		ParentID:       parentID,
		RequestType:    req.RequestType,
		FormData:       req.FormData,
		AttachmentURLs: req.AttachmentURLs,
	}

	if err := h.service.SubmitRequest(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Istanza presentata con successo",
		"request": item,
	})
}

func (h *Handler) ListRequests(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}

	role := c.GetString("role")
	userID := c.GetString("user_id")

	var parentFilter string
	if role == "parent" {
		parentFilter = userID
	} else {
		parentFilter = c.Query("parent_id")
	}

	status := c.Query("status")

	list, err := h.service.ListRequests(c.Request.Context(), schoolID, parentFilter, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

func (h *Handler) GetRequest(c *gin.Context) {
	id := c.Param("id")
	req, err := h.service.GetRequest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "istanza non trovata"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) ReviewRequest(c *gin.Context) {
	role := c.GetString("role")
	switch role {
	case "admin", "superadmin", "secretary", "principal", "vice_principal":
		// Allowed
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "solo la segreteria o il dirigente possono revisionare le istanze"})
		return
	}

	id := c.Param("id")
	var body struct {
		Status          string `json:"status" binding:"required"`
		RejectionReason string `json:"rejection_reason"`
		ProtocolNumber  string `json:"protocol_number"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati richiesta non validi: " + err.Error()})
		return
	}

	reviewerID := c.GetString("user_id")

	if err := h.service.ReviewRequest(c.Request.Context(), id, body.Status, body.RejectionReason, body.ProtocolNumber, reviewerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stato istanza aggiornato con successo",
	})
}

func (h *Handler) ListDelegates(c *gin.Context) {
	studentID := c.Query("student_id")
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = c.Query("school_id")
	}

	delegates, err := h.service.ListPermanentDelegates(c.Request.Context(), studentID, schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": delegates,
	})
}
