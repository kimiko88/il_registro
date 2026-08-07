package signatures

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	sigs := rg.Group("/signatures")
	{
		sigs.POST("", h.SignDocument)
		sigs.GET("/document/:id", h.GetSignatures)
		sigs.GET("/cad-preservation/download", h.DownloadCadPackage)
	}
}

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

func (h *Handler) DownloadCadPackage(c *gin.Context) {
	schoolID := c.GetString("school_id")
	year := c.Query("academic_year")

	zipBytes, err := GenerateCadPreservationPackage(c.Request.Context(), schoolID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("pacchetto_conservazione_CAD_%s.zip", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/zip", zipBytes)
}
