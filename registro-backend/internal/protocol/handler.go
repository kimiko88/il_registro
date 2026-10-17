package protocol

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	pr := rg.Group("/protocol")
	{
		pr.POST("", h.ProtocolDocument)
		pr.GET("", h.ListEntries)
		pr.GET("/:id", h.GetEntry)
		pr.GET("/:id/segnatura.xml", h.DownloadSegnaturaXML)
		pr.GET("/by-entity/:entityType/:entityId", h.GetEntityProtocol)
	}
}

func (h *Handler) ProtocolDocument(c *gin.Context) {
	role := c.GetString("role")
	switch role {
	case "admin", "superadmin", "secretary", "principal", "vice_principal", "assistente_protocollo":
		// Allowed
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "solo il personale di segreteria o protocollo può protocollare gli atti"})
		return
	}

	var req struct {
		Subject                string `json:"subject" binding:"required"`
		Sender                 string `json:"sender" binding:"required"`
		Recipient              string `json:"recipient" binding:"required"`
		FlowDirection          string `json:"flow_direction"` // in, out, internal
		ClassificationTitle    int    `json:"classification_title"`
		ClassificationClass    string `json:"classification_class"`
		ClassificationFascicle string `json:"classification_fascicle"`
		DocumentHashSHA256     string `json:"document_hash_sha256"`
		DocumentFileURL        string `json:"document_file_url"`
		EntityType             string `json:"entity_type"`
		EntityID               string `json:"entity_id"`
		SchoolName             string `json:"school_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati segnatura non validi: " + err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}
	userID := c.GetString("user_id")

	entry := &ProtocolEntry{
		SchoolID:               schoolID,
		Subject:                req.Subject,
		Sender:                 req.Sender,
		Recipient:              req.Recipient,
		FlowDirection:          req.FlowDirection,
		ClassificationTitle:    req.ClassificationTitle,
		ClassificationClass:    req.ClassificationClass,
		ClassificationFascicle: req.ClassificationFascicle,
		DocumentHashSHA256:     req.DocumentHashSHA256,
		DocumentFileURL:        req.DocumentFileURL,
		ProtocolledBy:          userID,
	}

	res, _, err := h.service.ProtocolDocument(c.Request.Context(), entry, req.EntityType, req.EntityID, req.SchoolName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Documento protocollato a norma AgID",
		"protocol": res,
	})
}

func (h *Handler) ListEntries(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}
	yearStr := c.Query("year")
	year := 0
	if yearStr != "" {
		year, _ = strconv.Atoi(yearStr)
	}
	flow := c.Query("flow_direction")

	entries, err := h.service.ListEntries(c.Request.Context(), schoolID, year, flow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": entries,
	})
}

func (h *Handler) GetEntry(c *gin.Context) {
	id := c.Param("id")
	entry, err := h.service.GetEntry(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "atto di protocollo non trovato"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *Handler) DownloadSegnaturaXML(c *gin.Context) {
	id := c.Param("id")
	entry, err := h.service.GetEntry(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "atto non trovato"})
		return
	}

	xmlBytes, err := GenerateSegnaturaXML(*entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "errore generazione XML"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=Segnatura.xml")
	c.Data(http.StatusOK, "application/xml", xmlBytes)
}

func (h *Handler) GetEntityProtocol(c *gin.Context) {
	entityType := c.Param("entityType")
	entityID := c.Param("entityId")

	entry, err := h.service.GetEntityProtocol(c.Request.Context(), entityType, entityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nessun protocollo associato a questa entità"})
		return
	}

	c.JSON(http.StatusOK, entry)
}
