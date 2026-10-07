package inventory

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

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/inventory")
	{
		g.POST("/assets", h.CreateAsset)
		g.GET("/assets", h.ListAssets)
		g.GET("/assets/:id/label", h.GetLabel)
		g.POST("/loans", h.CreateLoan)
		g.POST("/loans/:id/return", h.ReturnLoan)
	}
}

func (h *Handler) CreateAsset(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset, err := h.service.CreateAsset(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, asset)
}

func (h *Handler) ListAssets(c *gin.Context) {
	schoolID := c.GetString("school_id")
	category := c.Query("category")
	location := c.Query("location")

	assets, err := h.service.ListAssets(c.Request.Context(), schoolID, category, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(assets), "assets": assets})
}

func (h *Handler) GetLabel(c *gin.Context) {
	id := c.Param("id")
	label, err := h.service.GetLabel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, label)
}

func (h *Handler) CreateLoan(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract, err := h.service.CreateLoanContract(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, contract)
}

func (h *Handler) ReturnLoan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ConditionNotes string `json:"condition_notes"`
	}
	_ = c.ShouldBindJSON(&req)

	contract, err := h.service.ReturnLoanDevice(c.Request.Context(), id, req.ConditionNotes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contract)
}
