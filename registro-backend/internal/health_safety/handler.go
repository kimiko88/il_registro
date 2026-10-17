package health_safety

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/health-safety")
	{
		group.GET("/medications", h.ListMedications)
		group.POST("/medications", h.RegisterMedication)
		group.GET("/medications/:id", h.GetMedication)
		group.POST("/medications/:id/administer", h.LogAdministration)
		group.GET("/medications/:id/administrations", h.ListAdministrations)

		group.GET("/accidents", h.ListAccidents)
		group.POST("/accidents", h.RecordAccident)
		group.GET("/accidents/:id", h.GetAccident)
		group.POST("/accidents/:id/transmit-inail", h.TransmitInail)
	}
}

func (h *Handler) ListMedications(c *gin.Context) {
	schoolID := c.GetString("school_id")
	classID := c.Query("class_id")

	meds, err := h.svc.ListMedications(schoolID, classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"medications": meds})
}

func (h *Handler) RegisterMedication(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var m StudentMedication
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.SchoolID = schoolID

	if err := h.svc.RegisterMedication(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *Handler) GetMedication(c *gin.Context) {
	id := c.Param("id")
	m, err := h.svc.GetMedication(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "farmaco non trovato"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *Handler) LogAdministration(c *gin.Context) {
	medID := c.Param("id")
	userID := c.GetString("user_id")

	var a MedicationAdministration
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.MedicationID = medID
	a.AdministeredBy = userID

	if err := h.svc.LogAdministration(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *Handler) ListAdministrations(c *gin.Context) {
	medID := c.Param("id")
	list, err := h.svc.ListAdministrations(medID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"administrations": list})
}

func (h *Handler) ListAccidents(c *gin.Context) {
	schoolID := c.GetString("school_id")
	accidents, err := h.svc.ListAccidents(schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accidents": accidents})
}

func (h *Handler) RecordAccident(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var a SchoolAccident
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.SchoolID = schoolID

	if err := h.svc.RecordAccident(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *Handler) GetAccident(c *gin.Context) {
	id := c.Param("id")
	acc, err := h.svc.GetAccident(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "infortunio non trovato"})
		return
	}
	c.JSON(http.StatusOK, acc)
}

func (h *Handler) TransmitInail(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Protocol string `json:"protocol" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.TransmitInail(id, req.Protocol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "denuncia INAIL registrata con successo"})
}
