package interpelli

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterPublicRoutes(r *gin.RouterGroup) {
	pub := r.Group("/public/interpelli")
	{
		pub.GET("", h.ListPublic)
		pub.GET("/:id", h.GetPublic)
		pub.POST("/:id/candidatura", h.SubmitCandidatura)
	}
}

func (h *Handler) RegisterProtectedRoutes(r *gin.RouterGroup) {
	prot := r.Group("/interpelli")
	{
		prot.POST("", h.CreateNotice)
		prot.GET("/:id/graduatoria", h.GetGraduatoria)
		prot.POST("/candidature/:id/convoca", h.Convoca)
		prot.POST("/candidature/:id/risposta", h.Rispondi)
	}
}

func (h *Handler) ListPublic(c *gin.Context) {
	schoolID := c.Query("school_id")
	concorso := c.Query("classe_concorso")
	notices, err := h.service.ListPublicNotices(c.Request.Context(), schoolID, concorso)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notices)
}

func (h *Handler) GetPublic(c *gin.Context) {
	id := c.Param("id")
	notice, err := h.service.GetNotice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notice)
}

func (h *Handler) CreateNotice(c *gin.Context) {
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notice, err := h.service.CreateNotice(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, notice)
}

func (h *Handler) SubmitCandidatura(c *gin.Context) {
	noticeID := c.Param("id")
	var req SubmitCandidaturaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cand, err := h.service.SubmitCandidatura(c.Request.Context(), noticeID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cand)
}

func (h *Handler) GetGraduatoria(c *gin.Context) {
	noticeID := c.Param("id")
	graduatoria, err := h.service.GetGraduatoria(c.Request.Context(), noticeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"notice_id":   noticeID,
		"total":       len(graduatoria),
		"graduatoria": graduatoria,
	})
}

func (h *Handler) Convoca(c *gin.Context) {
	candID := c.Param("id")
	var req ConvocaRequest
	_ = c.ShouldBindJSON(&req)

	cand, err := h.service.ConvocaCandidato(c.Request.Context(), candID, req.HoursToRespond)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cand)
}

func (h *Handler) Rispondi(c *gin.Context) {
	candID := c.Param("id")
	var req RispondiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cand, err := h.service.RispondiConvocazione(c.Request.Context(), candID, req.Risposta, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cand)
}
