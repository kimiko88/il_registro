package signatures

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc     Service
	feqSvc  QualifiedService
	sidiSvc *SidiExportService
}

func NewHandler(svc Service, feqSvc QualifiedService, sidiSvc *SidiExportService) *Handler {
	return &Handler{svc: svc, feqSvc: feqSvc, sidiSvc: sidiSvc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	sigs := rg.Group("/signatures")
	{
		// Firma FEA base (legacy)
		sigs.POST("", h.SignDocument)
		sigs.GET("/document/:id", h.GetSignatures)

		// Firma FEQ/FES con valore legale (CAD art. 21)
		sigs.POST("/qualified", h.SignQualified)
		sigs.GET("/qualified/:id/verify", h.VerifyQualified)
		sigs.GET("/qualified/document/:doc_id", h.GetQualifiedByDocument)

		// Conservazione sostitutiva CAD (DPCM 3/12/2013)
		sigs.GET("/cad-preservation/download", h.DownloadCadPackage)
	}

	// Export SIDI/MIUR
	sidi := rg.Group("/sidi")
	{
		sidi.GET("/export/:school_id", h.ExportSidi)
	}
}

// ── FEA base (legacy) ─────────────────────────────────────────────────────

func (h *Handler) SignDocument(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.IPAddress = c.ClientIP()
	sig, err := h.svc.SignDocument(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sig)
}


func (h *Handler) GetSignatures(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	docID := c.Param("id")
	if docID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document id required"})
		return
	}
	sigs, err := h.svc.GetSignatures(docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sigs)
}

// ── FEQ/FES con valore legale ─────────────────────────────────────────────

// SignQualified appone una firma FEQ o FES su un documento.
// Richiede MFA TOTP abilitato per FEQ (CAD art. 26).
// POST /api/v1/signatures/qualified
func (h *Handler) SignQualified(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req QualifiedSignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sig, err := h.feqSvc.SignQualified(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	sig.IPAddress = c.ClientIP()
	c.JSON(http.StatusCreated, sig)
}

// VerifyQualified verifica crittograficamente una firma FEQ/FES.
// GET /api/v1/signatures/qualified/:id/verify
func (h *Handler) VerifyQualified(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	sigID := c.Param("id")
	result, err := h.feqSvc.VerifyQualified(c.Request.Context(), sigID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetQualifiedByDocument restituisce le firme qualificate su un documento.
// GET /api/v1/signatures/qualified/document/:doc_id
func (h *Handler) GetQualifiedByDocument(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	docID := c.Param("doc_id")
	sigs, err := h.feqSvc.GetQualifiedByDocument(c.Request.Context(), docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sigs)
}

// ── Export SIDI/MIUR ─────────────────────────────────────────────────────

// ExportSidi genera e scarica il pacchetto ZIP per la trasmissione SIDI/MIUR.
// GET /api/v1/sidi/export/:school_id?tipologia=SCRUTINI&academic_year=2025/2026
func (h *Handler) ExportSidi(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	callerSchoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "principal" && role != "vice_principal" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo la presidenza, segreteria o amministratori possono esportare i flussi SIDI"})
		return
	}
	schoolID := c.Param("school_id")
	if role != "superadmin" && callerSchoolID != "" && callerSchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: non autorizzato per questa scuola"})
		return
	}

	tipologia := c.DefaultQuery("tipologia", "SCRUTINI")
	academicYear := c.DefaultQuery("academic_year", "2025/2026")
	schoolName := c.DefaultQuery("school_name", "Istituto Scolastico")

	zipBytes, record, err := h.sidiSvc.GenerateSidiPackage(
		c.Request.Context(), schoolID, schoolName, academicYear, tipologia,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("SIDI_%s_%s_%s.zip", schoolID, tipologia, time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("X-SIDI-Hash-Integrita", record.HashIntegrità)
	c.Header("X-SIDI-Stato", record.StatoTrasmissione)
	c.Data(http.StatusOK, "application/zip", zipBytes)
}

// ── Conservazione sostitutiva CAD ────────────────────────────────────────

// DownloadCadPackage genera e scarica il pacchetto di conservazione sostitutiva.
// GET /api/v1/signatures/cad-preservation/download?academic_year=2025/2026
func (h *Handler) DownloadCadPackage(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	if role != "principal" && role != "vice_principal" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo la presidenza, segreteria o amministratori possono scaricare il pacchetto di conservazione CAD"})
		return
	}

	schoolID := c.GetString("school_id")
	year := c.Query("academic_year")

	zipBytes, err := GenerateCadPreservationPackage(c.Request.Context(), schoolID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("conservazione_CAD_%s_%s.zip", schoolID, time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/zip", zipBytes)
}
