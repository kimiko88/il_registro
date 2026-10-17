package enrollment

import (
	"fmt"
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
	e := rg.Group("/enrollment")
	{
		e.POST("/import-sidi", h.ImportSIDI)
		e.GET("/applications", h.ListApplications)
		e.POST("/formation-drafts/generate", h.GenerateFormationDraft)
		e.GET("/formation-drafts", h.ListDrafts)
		e.GET("/formation-drafts/:id", h.GetDraft)
		e.PUT("/formation-drafts/:id", h.UpdateDraftAssignments)
		e.POST("/formation-drafts/:id/finalize", h.FinalizeDraft)
	}
}

func isSecretaryOrAdmin(role string) bool {
	switch role {
	case "admin", "superadmin", "secretary", "principal", "vice_principal", "assistente_alunni", "assistente_amministrativo":
		return true
	default:
		return false
	}
}

func (h *Handler) ImportSIDI(c *gin.Context) {
	role := c.GetString("role")
	if !isSecretaryOrAdmin(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	schoolID := c.GetString("school_id")
	academicYear := c.DefaultQuery("academic_year", "2026/2027")

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required (form-data 'file')"})
		return
	}
	defer func() { _ = file.Close() }()

	count, err := h.service.ImportSIDIApplications(c.Request.Context(), schoolID, academicYear, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("errore importazione SIDI: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Domande SIDI importate con successo",
		"imported": count,
	})
}

func (h *Handler) ListApplications(c *gin.Context) {
	role := c.GetString("role")
	if !isSecretaryOrAdmin(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	schoolID := c.GetString("school_id")
	academicYear := c.DefaultQuery("academic_year", "2026/2027")
	status := c.Query("status")

	apps, err := h.service.ListApplications(c.Request.Context(), schoolID, academicYear, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apps)
}

func (h *Handler) GenerateFormationDraft(c *gin.Context) {
	role := c.GetString("role")
	if !isSecretaryOrAdmin(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	schoolID := c.GetString("school_id")
	var req struct {
		Title        string          `json:"title"`
		AcademicYear string          `json:"academic_year"`
		Parameters   FormationParams `json:"parameters"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draft, err := h.service.GenerateFormationDraft(c.Request.Context(), schoolID, req.AcademicYear, req.Title, req.Parameters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, draft)
}

func (h *Handler) ListDrafts(c *gin.Context) {
	role := c.GetString("role")
	if !isSecretaryOrAdmin(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	schoolID := c.GetString("school_id")
	drafts, err := h.service.ListDrafts(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, drafts)
}

func (h *Handler) GetDraft(c *gin.Context) {
	role := c.GetString("role")
	if !isSecretaryOrAdmin(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	draft, err := h.service.GetDraft(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}
	c.JSON(http.StatusOK, draft)
}

func (h *Handler) UpdateDraftAssignments(c *gin.Context) {
	role := c.GetString("role")
	if !isSecretaryOrAdmin(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	var req ClassFormationDraftResult
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateDraftAssignments(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bozza classi aggiornata con successo"})
}

func (h *Handler) FinalizeDraft(c *gin.Context) {
	role := c.GetString("role")
	if !isSecretaryOrAdmin(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	if err := h.service.FinalizeDraft(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bozza finalizzata. Classi e studenti generati con successo!"})
}
