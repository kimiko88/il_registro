package personnel_desk

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/personnel-desk")
	{
		g.POST("/requests", h.CreateRequest)
		g.GET("/requests", h.ListRequests)
		g.GET("/requests/:id", h.GetRequest)
		g.PATCH("/requests/:id/submit", h.SubmitRequest)
		g.PATCH("/requests/:id/aa-review", h.AAReview)
		g.PATCH("/requests/:id/dsga-sign", h.DSGASign)
		g.PATCH("/requests/:id/ds-approve", h.DSApprove)
		g.DELETE("/requests/:id", h.DeleteRequest)
	}
}

func getAuthInfo(c *gin.Context) (userID, schoolID, role string, ok bool) {
	userID = c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return "", "", "", false
	}
	schoolID = c.GetString("school_id")
	role = c.GetString("role")
	return userID, schoolID, role, true
}

func (h *Handler) CreateRequest(c *gin.Context) {
	userID, schoolID, _, ok := getAuthInfo(c)
	if !ok {
		return
	}

	var input CreateDeskRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati non validi: " + err.Error()})
		return
	}

	req, err := h.service.CreateRequest(c.Request.Context(), schoolID, userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (h *Handler) ListRequests(c *gin.Context) {
	userID, schoolID, role, ok := getAuthInfo(c)
	if !ok {
		return
	}

	status := c.Query("status")
	list, err := h.service.ListRequests(c.Request.Context(), schoolID, userID, role, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetRequest(c *gin.Context) {
	userID, schoolID, role, ok := getAuthInfo(c)
	if !ok {
		return
	}

	id := c.Param("id")
	req, err := h.service.GetRequest(c.Request.Context(), schoolID, id, userID, role)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) SubmitRequest(c *gin.Context) {
	userID, schoolID, _, ok := getAuthInfo(c)
	if !ok {
		return
	}

	id := c.Param("id")
	req, err := h.service.SubmitRequest(c.Request.Context(), schoolID, id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) AAReview(c *gin.Context) {
	userID, schoolID, role, ok := getAuthInfo(c)
	if !ok {
		return
	}

	if role != "assistente_amministrativo" && role != "secretary" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo l'assistente amministrativo o la segreteria possono istruire la pratica"})
		return
	}

	var input AAReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati non validi: " + err.Error()})
		return
	}

	id := c.Param("id")
	req, err := h.service.AAReview(c.Request.Context(), schoolID, id, userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) DSGASign(c *gin.Context) {
	userID, schoolID, role, ok := getAuthInfo(c)
	if !ok {
		return
	}

	if role != "dsga" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo il DSGA può apporre il visto"})
		return
	}

	var input DSGASignInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati non validi: " + err.Error()})
		return
	}

	id := c.Param("id")
	req, err := h.service.DSGASign(c.Request.Context(), schoolID, id, userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) DSApprove(c *gin.Context) {
	userID, schoolID, role, ok := getAuthInfo(c)
	if !ok {
		return
	}

	if role != "principal" && role != "vice_principal" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo il Dirigente Scolastico può emanare il provvedimento finale"})
		return
	}

	var input DSApproveInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati non validi: " + err.Error()})
		return
	}

	id := c.Param("id")
	req, err := h.service.DSApprove(c.Request.Context(), schoolID, id, userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) DeleteRequest(c *gin.Context) {
	userID, schoolID, _, ok := getAuthInfo(c)
	if !ok {
		return
	}

	id := c.Param("id")
	if err := h.service.DeleteRequest(c.Request.Context(), schoolID, id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "richiesta eliminata"})
}
