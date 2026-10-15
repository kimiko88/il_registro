package sidi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	grp := r.Group("/sidi")

	// Only secretary, admin, superadmin
	grp.POST("/exports/generate", h.GenerateExport)
	grp.GET("/exports", h.GetExports)
	grp.GET("/exports/:id/download", h.DownloadExportXML)

	// WebService MIM Cooperazione Applicativa
	grp.POST("/sync-student-codes", h.SyncStudentCodes)
	grp.POST("/push-scrutiny-results", h.PushScrutinyResults)
	grp.GET("/cooperation-config", h.GetCooperationConfig)
}

func canAccessSidi(role string) bool {
	switch role {
	case "secretary", "admin", "superadmin", "principal", "vice_principal", "dsga", "assistente_amministrativo", "assistente_alunni":
		return true
	default:
		return false
	}
}

func (h *Handler) GenerateExport(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !canAccessSidi(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo la segreteria, il DSGA o l'amministrazione possono generare flussi SIDI"})
		return
	}

	var req GenerateSidiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exportRec, validation, err := h.service.GenerateExport(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"export":     exportRec,
		"validation": validation,
	})
}

func (h *Handler) GetExports(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !canAccessSidi(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	list, err := h.service.GetExports(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) DownloadExportXML(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !canAccessSidi(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	xmlData, err := h.service.GetExportXML(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "export non trovato"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=flusso_sidi.xml")
	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(xmlData))
}

func (h *Handler) SyncStudentCodes(c *gin.Context) {
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if !canAccessSidi(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req struct {
		Students []StudenteSIDI `json:"students"`
	}
	_ = c.ShouldBindJSON(&req)

	resp, err := h.service.SyncStudentCodes(c.Request.Context(), schoolID, req.Students)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) PushScrutinyResults(c *gin.Context) {
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if !canAccessSidi(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req PushScrutinyResultsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.PushScrutinyResults(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCooperationConfig(c *gin.Context) {
	schoolID := c.GetString("school_id")
	cfg, err := h.service.GetCooperationConfig(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}
