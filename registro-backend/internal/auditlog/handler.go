package auditlog

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"registro-backend/pkg/logger"
	"registro-backend/pkg/upload"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
	db      *sql.DB
}

func NewHandler(service Service, db ...*sql.DB) *Handler {
	h := &Handler{service: service}
	if len(db) > 0 {
		h.db = db[0]
	}
	return h
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	group := r.Group("/audit-log")
	{
		group.GET("", h.List)
		group.GET("/export", h.ExportCSV)
		group.GET("/immutability-chain", h.GetImmutabilityChain)
	}
}

func (h *Handler) List(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "system_auditor" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access restricted to admin, superadmin, secretary, and system_auditor"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	params := FilterParams{
		SchoolID:   c.GetString("school_id"),
		ActorID:    c.Query("actor_id"),
		Action:     c.Query("action"),
		EntityType: c.Query("entity_type"),
		From:       c.Query("from"),
		To:         c.Query("to"),
		Page:       page,
		Limit:      limit,
	}

	result, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		logger.Log.Errorf("Audit log list error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "errore interno del server durante il recupero dei log di audit"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) ExportCSV(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "system_auditor" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access restricted to admin, superadmin, secretary, and system_auditor"})
		return
	}

	params := FilterParams{
		SchoolID:   c.GetString("school_id"),
		ActorID:    c.Query("actor_id"),
		Action:     c.Query("action"),
		EntityType: c.Query("entity_type"),
		From:       c.Query("from"),
		To:         c.Query("to"),
	}

	csvBytes, err := h.service.ExportCSV(c.Request.Context(), params)
	if err != nil {
		logger.Log.Errorf("Audit log export error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "errore interno del server durante l'esportazione del registro di audit"})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", upload.FormatContentDisposition(fmt.Sprintf("audit_logs_%s.csv", c.GetString("school_id"))))
	c.Data(http.StatusOK, "text/csv", csvBytes)
}

func (h *Handler) GetImmutabilityChain(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "secretary" && role != "system_auditor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	report, err := VerifyChainIntegrity(c.Request.Context(), h.db)
	if err != nil {
		logger.Log.Errorf("Immutability chain verification error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "errore interno durante la verifica della catena di immutabilità"})
		return
	}
	c.JSON(http.StatusOK, report)
}
