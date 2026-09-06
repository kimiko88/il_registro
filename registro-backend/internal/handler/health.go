package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *HealthHandler) Metrics(c *gin.Context) {
	stats := h.db.Stats()
	metrics := fmt.Sprintf(`# HELP db_open_connections The number of established connections both in use and idle.
# TYPE db_open_connections gauge
db_open_connections %d
# HELP db_in_use_connections The number of connections currently in use.
# TYPE db_in_use_connections gauge
db_in_use_connections %d
# HELP db_idle_connections The number of idle connections.
# TYPE db_idle_connections gauge
db_idle_connections %d
`, stats.OpenConnections, stats.InUse, stats.Idle)

	c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte(metrics))
}
