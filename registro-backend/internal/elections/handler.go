package elections

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	el := rg.Group("/elections")
	{
		el.POST("", h.CreateElection)
		el.GET("", h.ListElections)
		el.GET("/:id", h.GetElectionDetails)
		el.POST("/:id/vote", h.CastVote)
		el.GET("/:id/scrutiny", h.GetScrutiny)
		el.PUT("/:id/close", h.CloseElection)
	}
}

func (h *Handler) CreateElection(c *gin.Context) {
	role := c.GetString("role")
	switch role {
	case "admin", "superadmin", "secretary", "principal", "vice_principal":
		// Allowed
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso non autorizzato"})
		return
	}

	var req struct {
		Title          string    `json:"title" binding:"required"`
		ElectionTier   string    `json:"election_tier" binding:"required"`
		TargetRole     string    `json:"target_role" binding:"required"`
		ClassID        *string   `json:"class_id"`
		StartTime      time.Time `json:"start_time"`
		EndTime        time.Time `json:"end_time"`
		MaxPreferences int       `json:"max_preferences"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati elezione non validi: " + err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}

	election := &SchoolElection{
		SchoolID:       schoolID,
		Title:          req.Title,
		ElectionTier:   req.ElectionTier,
		TargetRole:     req.TargetRole,
		ClassID:        req.ClassID,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		MaxPreferences: req.MaxPreferences,
	}

	if err := h.service.CreateElection(c.Request.Context(), election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Elezione creata con successo",
		"election": election,
	})
}

func (h *Handler) ListElections(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}

	list, err := h.service.ListElections(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

func (h *Handler) GetElectionDetails(c *gin.Context) {
	electionID := c.Param("id")
	e, lists, err := h.service.GetElectionDetails(c.Request.Context(), electionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "elezione non trovata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"election": e,
		"lists":    lists,
	})
}

func (h *Handler) CastVote(c *gin.Context) {
	electionID := c.Param("id")
	voterID := c.GetString("user_id")
	if voterID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "utente non identificato"})
		return
	}

	var payload CastVotePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "voto non valido: " + err.Error()})
		return
	}

	receipt, err := h.service.CastVote(c.Request.Context(), electionID, voterID, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Voto anonimo registrato nell'urna digitale",
		"receipt": receipt,
	})
}

func (h *Handler) GetScrutiny(c *gin.Context) {
	electionID := c.Param("id")
	seatsStr := c.DefaultQuery("seats", "4")
	seats, _ := strconv.Atoi(seatsStr)

	res, err := h.service.GetScrutiny(c.Request.Context(), electionID, seats)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scrutiny": res,
	})
}

func (h *Handler) CloseElection(c *gin.Context) {
	role := c.GetString("role")
	switch role {
	case "admin", "superadmin", "secretary", "principal", "vice_principal":
		// Allowed
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso non autorizzato"})
		return
	}

	electionID := c.Param("id")
	if err := h.service.CloseElection(c.Request.Context(), electionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Seggio elettorale chiuso con successo",
	})
}
