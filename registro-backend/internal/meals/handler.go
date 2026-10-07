package meals

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
	g := r.Group("/meals")
	{
		g.POST("/roll-call", h.RecordRollCall)
		g.GET("/catering-report", h.GetCateringReport)
		g.POST("/special-diet", h.RegisterDiet)
		g.GET("/special-diet/:student_id", h.GetDiet)
		g.GET("/special-diets", h.ListDiets)
		g.GET("/wallet/:student_id", h.GetWallet)
		g.POST("/wallet/:student_id/topup", h.TopUpWallet)
	}
}

func (h *Handler) RecordRollCall(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req MealRollCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.service.RecordRollCallAndAggregate(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *Handler) GetCateringReport(c *gin.Context) {
	schoolID := c.GetString("school_id")
	date := c.DefaultQuery("date", "")
	rep, err := h.service.repo.GetCateringReport(c.Request.Context(), schoolID, date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rep)
}

func (h *Handler) RegisterDiet(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req SpecialDietRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	diet, err := h.service.RegisterSpecialDiet(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, diet)
}

func (h *Handler) GetDiet(c *gin.Context) {
	stdID := c.Param("student_id")
	diet, err := h.service.GetSpecialDiet(c.Request.Context(), stdID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, diet)
}

func (h *Handler) ListDiets(c *gin.Context) {
	schoolID := c.GetString("school_id")
	diets, err := h.service.ListSpecialDiets(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(diets), "diets": diets})
}

func (h *Handler) GetWallet(c *gin.Context) {
	stdID := c.Param("student_id")
	wallet, err := h.service.GetWallet(c.Request.Context(), stdID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (h *Handler) TopUpWallet(c *gin.Context) {
	stdID := c.Param("student_id")
	var req struct {
		Amount    float64 `json:"amount" binding:"required,gt=0"`
		PagoPAIUV string  `json:"pagopa_iuv"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.service.TopUpWallet(c.Request.Context(), stdID, req.Amount, req.PagoPAIUV)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallet)
}
