package payments

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/payments")
	{
		g.GET("", h.List)
		g.GET("/:id", h.Get)
		g.GET("/:id/bollettino", h.GetBollettino)
		g.POST("/:id/checkout", h.Checkout)
		g.POST("/:id/pay", h.Pay)
		g.POST("", h.Create)
		g.POST("/reconcile-opi", h.ReconcileOPI)
		g.GET("/solleciti", h.GetSolleciti)
	}
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetStudentID := c.Query("student_id")
	summary, err := h.service.List(c.Request.Context(), userID, role, schoolID, targetStudentID)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *Handler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	paymentID := c.Param("id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payment, err := h.service.GetByID(c.Request.Context(), userID, role, schoolID, paymentID)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if err == ErrPaymentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "pagamento non trovato"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *Handler) Pay(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	paymentID := c.Param("id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req PayRequest
	if c.Request.ContentLength > 0 {
		_ = c.ShouldBindJSON(&req)
	}
	if req.PaymentMethod == "" {
		req.PaymentMethod = "PagoPA"
	}

	payment, err := h.service.Pay(c.Request.Context(), userID, role, paymentID, req.PaymentMethod)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if strings.Contains(err.Error(), "già stato pagato") || strings.Contains(err.Error(), "già effettuato") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pagamento registrato con successo",
		"payment": payment,
	})
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.Create(c.Request.Context(), userID, role, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *Handler) GetBollettino(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	paymentID := c.Param("id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	notice, err := h.service.GetBollettino(c.Request.Context(), userID, role, schoolID, paymentID)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if err == ErrPaymentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "pagamento non trovato"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notice)
}

func (h *Handler) Checkout(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	paymentID := c.Param("id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		ReturnURL string `json:"return_url"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.ReturnURL == "" {
		req.ReturnURL = "https://scuola.edu.it/pagamenti/conferma"
	}

	session, err := h.service.Checkout(c.Request.Context(), userID, role, paymentID, req.ReturnURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *Handler) ReconcileOPI(c *gin.Context) {
	role := c.GetString("role")
	schoolID := c.GetString("school_id")

	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	format := c.DefaultQuery("format", "OPI_XML")
	data, err := c.GetRawData()
	if err != nil || len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "flusso OPI mancante o non valido"})
		return
	}

	report, err := h.service.ReconcileOPI(c.Request.Context(), schoolID, format, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *Handler) GetSolleciti(c *gin.Context) {
	role := c.GetString("role")
	schoolID := c.GetString("school_id")

	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	solleciti, err := h.service.GetSolleciti(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     len(solleciti),
		"solleciti": solleciti,
	})
}
