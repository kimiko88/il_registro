package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"registro-backend/internal/metrics"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	if err := h.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN", "error": "database unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "READY"})
}

// Ping is a lightweight no-auth, no-DB endpoint used by the frontend to verify
// real backend reachability (navigator.onLine alone is not reliable).
func (h *HealthHandler) Ping(c *gin.Context) {
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *HealthHandler) Metrics(c *gin.Context) {
	output := metrics.DefaultRegistry.GeneratePrometheus(h.db)
	c.Data(http.StatusOK, "text/plain; version=0.0.4; charset=utf-8", []byte(output))
}
