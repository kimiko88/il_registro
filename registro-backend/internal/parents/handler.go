package parents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     *Service
	dualService *DualSignatureService
	accessGuard *AccessGuard
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) SetDualSignatureService(ds *DualSignatureService) {
	h.dualService = ds
}

func (h *Handler) SetAccessGuard(ag *AccessGuard) {
	h.accessGuard = ag
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	p := r.Group("/parents")
	{
		p.GET("/dashboard", h.GetDashboard)
		p.GET("/dashboard/stats", h.GetDashboardStats)
		p.GET("/child/:studentId/grades-average", h.GetChildGradesAverage)

		// Dual signature & Custody routes
		p.GET("/dual-authorizations", h.ListDualAuthorizations)
		p.POST("/dual-authorizations", h.CreateDualAuthorization)
		p.POST("/dual-authorizations/:id/sign", h.SignDualAuthorization)
		p.POST("/dual-authorizations/:id/reject", h.RejectDualAuthorization)
		p.GET("/child/:studentId/custody", h.GetChildCustodyInfo)
	}
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stats, err := h.service.GetDashboardStats(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	dash, err := h.service.GetDashboard(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dash)
}

func (h *Handler) GetChildGradesAverage(c *gin.Context) {
	parentID := c.GetString("user_id")
	studentID := c.Param("studentId")

	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	avg, err := h.service.GetChildGradesAverage(c.Request.Context(), parentID, studentID)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"student_id": studentID,
		"average":    avg,
	})
}

// Dual Signature Handler Methods

func (h *Handler) ListDualAuthorizations(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if h.dualService == nil {
		c.JSON(http.StatusOK, []DualParentalAuthorization{})
		return
	}

	auths, err := h.dualService.ListAuthorizations(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, auths)
}

func (h *Handler) CreateDualAuthorization(c *gin.Context) {
	parentID := c.GetString("user_id")
	role := c.GetString("role")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_ = role

	if h.dualService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "dual service unavailable"})
		return
	}

	var req CreateDualAuthParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := h.dualService.CreateAuthorization(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}

func (h *Handler) SignDualAuthorization(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if h.dualService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "dual service unavailable"})
		return
	}

	authID := c.Param("id")
	var req SignDualAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	signed, err := h.dualService.SignDocument(c.Request.Context(), authID, parentID, req.PIN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, signed)
}

func (h *Handler) RejectDualAuthorization(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if h.dualService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "dual service unavailable"})
		return
	}

	authID := c.Param("id")
	var req RejectDualAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rejected, err := h.dualService.RejectDocument(c.Request.Context(), authID, parentID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rejected)
}

func (h *Handler) GetChildCustodyInfo(c *gin.Context) {
	parentID := c.GetString("user_id")
	studentID := c.Param("studentId")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if h.accessGuard != nil {
		allowed, err := h.accessGuard.CanAccessStudent(c.Request.Context(), parentID, studentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accesso vietato per provvedimento giudiziario"})
			return
		}
	}

	if h.dualService == nil {
		c.JSON(http.StatusOK, gin.H{"custody_type": "shared"})
		return
	}

	info, err := h.dualService.repo.GetCustodyInfo(c.Request.Context(), parentID, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}
