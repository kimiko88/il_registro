package strike

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/strike-notices")
	{
		g.POST("", h.CreateNotice)
		g.GET("", h.ListNotices)
		g.GET("/:id", h.GetNotice)
		g.DELETE("/:id", h.DeleteNotice)
		g.POST("/:id/declare", h.SubmitDeclaration)
		g.GET("/:id/summary", h.GetNoticeSummary)
	}
}

func getSchoolID(c *gin.Context, h *Handler, userID string) string {
	schoolID := c.GetString("school_id")
	if schoolID == "" && h.service != nil {
		schoolID = h.service.repo.ResolveSchoolID(c.Request.Context(), userID)
	}
	return schoolID
}

func (h *Handler) CreateNotice(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := getSchoolID(c, h, userID)

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateStrikeNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notice, err := h.service.CreateNotice(c.Request.Context(), userID, role, schoolID, req)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notice)
}

func (h *Handler) ListNotices(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := getSchoolID(c, h, userID)

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	notices, err := h.service.ListNotices(c.Request.Context(), schoolID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notices)
}

func (h *Handler) GetNotice(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	schoolID := getSchoolID(c, h, userID)

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	notice, err := h.service.GetNotice(c.Request.Context(), id, schoolID, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comunicazione di sciopero non trovata"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notice)
}

func (h *Handler) DeleteNotice(c *gin.Context) {
	id := c.Param("id")
	role := c.GetString("role")
	userID := c.GetString("user_id")
	schoolID := getSchoolID(c, h, userID)

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.DeleteNotice(c.Request.Context(), id, schoolID, role); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comunicazione di sciopero eliminata con successo"})
}

func (h *Handler) SubmitDeclaration(c *gin.Context) {
	noticeID := c.Param("id")
	userID := c.GetString("user_id")
	schoolID := getSchoolID(c, h, userID)
	ipAddress := c.ClientIP()

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SubmitDeclarationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decl, err := h.service.SubmitDeclaration(c.Request.Context(), noticeID, schoolID, userID, ipAddress, req)
	if err != nil {
		if errors.Is(err, ErrDeclarationDeadlinePassed) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "code": "DEADLINE_PASSED"})
			return
		}
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrInvalidIntention) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, decl)
}

func (h *Handler) GetNoticeSummary(c *gin.Context) {
	noticeID := c.Param("id")
	role := c.GetString("role")
	userID := c.GetString("user_id")
	schoolID := getSchoolID(c, h, userID)

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	summary, err := h.service.GetNoticeSummary(c.Request.Context(), noticeID, schoolID, role)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
