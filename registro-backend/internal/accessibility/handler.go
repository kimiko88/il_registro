package accessibility

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// SubmitPublic consente a chiunque (genitori, studenti, cittadini) di segnalare una barriera digitale
func (h *Handler) SubmitPublic(c *gin.Context) {
	var req CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dati della richiesta non validi: " + err.Error()})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	var uAgentPtr, ipPtr *string
	if userAgent != "" {
		uAgentPtr = &userAgent
	}
	if ipAddress != "" {
		ipPtr = &ipAddress
	}

	// Se l'utente è autenticato tramite cookie o header opzionale
	var userIDPtr, schoolIDPtr *string
	if uid, exists := c.Get("user_id"); exists {
		if s, ok := uid.(string); ok && s != "" {
			userIDPtr = &s
		}
	}
	if sid, exists := c.Get("school_id"); exists {
		if s, ok := sid.(string); ok && s != "" {
			schoolIDPtr = &s
		}
	}

	resp, err := h.svc.SubmitFeedback(c.Request.Context(), &req, userIDPtr, schoolIDPtr, uAgentPtr, ipPtr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// List consente ad amministratori e segreteria di consultare le segnalazioni pervenute
func (h *Handler) List(c *gin.Context) {
	status := c.Query("status")
	schoolID := c.Query("school_id")
	if schoolID == "" {
		if sid, exists := c.Get("school_id"); exists {
			if s, ok := sid.(string); ok {
				schoolID = s
			}
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	feedbacks, total, err := h.svc.ListFeedbacks(c.Request.Context(), schoolID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossibile recuperare le segnalazioni: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"feedbacks": feedbacks,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// UpdateStatus aggiorna lo stato di presa in carico o risoluzione
func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID segnalazione obbligatorio"})
		return
	}

	var req struct {
		Status        string `json:"status" binding:"required"`
		ResponseNotes string `json:"response_notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dati richiesta non validi"})
		return
	}

	if err := h.svc.UpdateFeedbackStatus(c.Request.Context(), id, req.Status, req.ResponseNotes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stato segnalazione aggiornato con successo"})
}
