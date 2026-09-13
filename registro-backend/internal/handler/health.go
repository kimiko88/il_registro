package handler

import (
	"context"
	"database/sql"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"registro-backend/internal/metrics"
	"registro-backend/pkg/version"
)

type ComponentHealth struct {
	Status    string  `json:"status"`
	LatencyMS float64 `json:"latency_ms,omitempty"`
	Error     string  `json:"error,omitempty"`
}

type ReadyResponse struct {
	Status     string                     `json:"status"`
	Timestamp  string                     `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components"`
	System     map[string]interface{}     `json:"system,omitempty"`
}

type HealthHandler struct {
	db          *sql.DB
	redisClient *redis.Client
}

func NewHealthHandler(db *sql.DB, redisClient ...*redis.Client) *HealthHandler {
	h := &HealthHandler{db: db}
	if len(redisClient) > 0 && redisClient[0] != nil {
		h.redisClient = redisClient[0]
	}
	return h
}

// Health provides basic health status
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"version": version.Version,
	})
}

// Live is Kubernetes liveness probe: verifies the process is responsive
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"version":   version.Version,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Ready is Kubernetes readiness probe: performs deep concurrent checks of dependencies
func (h *HealthHandler) Ready(c *gin.Context) {
	reqCtx := context.Background()
	if c.Request != nil {
		reqCtx = c.Request.Context()
	}
	ctx, cancel := context.WithTimeout(reqCtx, 2*time.Second)
	defer cancel()

	components := make(map[string]ComponentHealth)
	allHealthy := true

	// 1. Check PostgreSQL database
	if h.db != nil {
		dbStart := time.Now()
		if err := h.db.PingContext(ctx); err != nil {
			allHealthy = false
			components["database"] = ComponentHealth{
				Status: "DOWN",
				Error:  err.Error(),
			}
		} else {
			components["database"] = ComponentHealth{
				Status:    "UP",
				LatencyMS: float64(time.Since(dbStart).Microseconds()) / 1000.0,
			}
		}
	} else {
		allHealthy = false
		components["database"] = ComponentHealth{
			Status: "DOWN",
			Error:  "database connection not configured",
		}
	}

	// 2. Check Redis if configured
	if h.redisClient != nil {
		redisStart := time.Now()
		if err := h.redisClient.Ping(ctx).Err(); err != nil {
			// Redis failure is non-fatal for readiness if cache is optional, but reported
			components["redis"] = ComponentHealth{
				Status: "DEGRADED",
				Error:  err.Error(),
			}
		} else {
			components["redis"] = ComponentHealth{
				Status:    "UP",
				LatencyMS: float64(time.Since(redisStart).Microseconds()) / 1000.0,
			}
		}
	}

	// 3. System runtime statistics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	systemStats := map[string]interface{}{
		"goroutines": runtime.NumGoroutine(),
		"alloc_mb":   float64(m.Alloc) / (1024 * 1024),
		"sys_mb":     float64(m.Sys) / (1024 * 1024),
		"num_gc":     m.NumGC,
	}

	resp := ReadyResponse{
		Status:     "READY",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Components: components,
		System:     systemStats,
	}

	if !allHealthy {
		resp.Status = "DOWN"
		c.JSON(http.StatusServiceUnavailable, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
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
