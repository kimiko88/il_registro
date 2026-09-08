package visitors

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
	g := r.Group("/visitors")
	{
		// Visitatori esterni
		g.POST("", h.RegisterVisitor)
		g.GET("", h.ListTodayVisitors)
		g.PATCH("/:id/exit", h.RecordVisitorExit)

		// Uscite anticipate studenti
		g.POST("/early-exits", h.RecordEarlyExit)
		g.GET("/early-exits", h.ListTodayEarlyExits)
		g.PATCH("/early-exits/:id/return", h.RecordStudentReturn)

		// Segnalazioni guasti
		g.POST("/maintenance", h.CreateMaintenanceReport)
		g.GET("/maintenance", h.ListMaintenanceReports)
		g.PATCH("/maintenance/:id/status", h.UpdateMaintenanceStatus)
	}
}

func authCheck(c *gin.Context) (userID, schoolID, role string, ok bool) {
	userID = c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return "", "", "", false
	}
	schoolID = c.GetString("school_id")
	role = c.GetString("role")
	if !CanAccessVisitorRegistry(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso riservato al personale scolastico"})
		return "", "", "", false
	}
	return userID, schoolID, role, true
}

func (h *Handler) RegisterVisitor(c *gin.Context) {
	userID, schoolID, _, ok := authCheck(c)
	if !ok {
		return
	}
	var req RegisterVisitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v, err := h.service.RegisterVisitor(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *Handler) ListTodayVisitors(c *gin.Context) {
	_, schoolID, _, ok := authCheck(c)
	if !ok {
		return
	}
	date := c.Query("date")
	visitors, err := h.service.ListTodayVisitors(c.Request.Context(), schoolID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, visitors)
}

func (h *Handler) RecordVisitorExit(c *gin.Context) {
	_, schoolID, _, ok := authCheck(c)
	if !ok {
		return
	}
	id := c.Param("id")
	var body RecordExitRequest
	_ = c.ShouldBindJSON(&body)
	if err := h.service.RecordVisitorExit(c.Request.Context(), id, schoolID, body.Notes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "uscita visitatore registrata"})
}

func (h *Handler) RecordEarlyExit(c *gin.Context) {
	userID, schoolID, _, ok := authCheck(c)
	if !ok {
		return
	}
	var req RecordEarlyExitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exit, err := h.service.RecordEarlyExit(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, exit)
}

func (h *Handler) ListTodayEarlyExits(c *gin.Context) {
	_, schoolID, _, ok := authCheck(c)
	if !ok {
		return
	}
	date := c.Query("date")
	exits, err := h.service.ListTodayEarlyExits(c.Request.Context(), schoolID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exits)
}

func (h *Handler) RecordStudentReturn(c *gin.Context) {
	_, schoolID, _, ok := authCheck(c)
	if !ok {
		return
	}
	id := c.Param("id")
	var body RecordStudentReturnRequest
	_ = c.ShouldBindJSON(&body)
	if err := h.service.RecordStudentReturn(c.Request.Context(), id, schoolID, body.Notes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rientro studente registrato"})
}

func (h *Handler) CreateMaintenanceReport(c *gin.Context) {
	userID, schoolID, _, ok := authCheck(c)
	if !ok {
		return
	}
	var req CreateMaintenanceReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rep, err := h.service.CreateMaintenanceReport(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rep)
}

func (h *Handler) ListMaintenanceReports(c *gin.Context) {
	_, schoolID, _, ok := authCheck(c)
	if !ok {
		return
	}
	statusFilter := c.Query("status")
	reports, err := h.service.ListMaintenanceReports(c.Request.Context(), schoolID, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reports)
}

func (h *Handler) UpdateMaintenanceStatus(c *gin.Context) {
	_, schoolID, _, ok := authCheck(c)
	if !ok {
		return
	}
	id := c.Param("id")
	var req UpdateMaintenanceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateMaintenanceStatus(c.Request.Context(), id, schoolID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "stato segnalazione aggiornato"})
}
