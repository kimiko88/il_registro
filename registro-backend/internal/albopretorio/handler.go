package albopretorio

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterPublicRoutes(r *gin.RouterGroup) {
	pub := r.Group("/public")
	{
		pub.GET("/albo-pretorio", h.ListPublicAlbo)
		pub.GET("/amministrazione-trasparente", h.ListTrasparenza)
		pub.GET("/amministrazione-trasparente/anac.xml", h.ExportANACXML)
	}
}

func (h *Handler) RegisterProtectedRoutes(r *gin.RouterGroup) {
	prot := r.Group("/albo-pretorio")
	{
		prot.POST("", h.PublishAct)
		prot.POST("/:id/defissione", h.DefiggiAtto)
		prot.GET("/:id/certificato", h.GetCertificato)
	}
}

func (h *Handler) ListPublicAlbo(c *gin.Context) {
	schoolID := c.Query("school_id")
	category := c.Query("categoria")
	query := c.Query("q")
	archived := c.Query("storico") == "true"
	year, _ := strconv.Atoi(c.Query("anno"))

	items, err := h.service.ListPublicAlbo(c.Request.Context(), schoolID, category, query, year, archived)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": len(items),
		"items": items,
	})
}

func (h *Handler) ListTrasparenza(c *gin.Context) {
	schoolID := c.Query("school_id")
	items, err := h.service.ListPublicAlbo(c.Request.Context(), schoolID, "", "", 0, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var trasparenza []*AlboItem
	for _, it := range items {
		if it.IsTransparencySection {
			trasparenza = append(trasparenza, it)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"total": len(trasparenza),
		"items": trasparenza,
	})
}

func (h *Handler) ExportANACXML(c *gin.Context) {
	schoolID := c.Query("school_id")
	year, _ := strconv.Atoi(c.DefaultQuery("anno", "2026"))

	xmlStr, err := h.service.GenerateANACXML(c.Request.Context(), schoolID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=anac_dataset_l190.xml")
	c.String(http.StatusOK, xmlStr)
}

func (h *Handler) PublishAct(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req PublishActRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.PublishAct(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) DefiggiAtto(c *gin.Context) {
	id := c.Param("id")
	dsName := c.DefaultPostForm("ds_name", "Dirigente Scolastico")
	force := c.Query("force") == "true"

	item, err := h.service.DefiggiAtto(c.Request.Context(), id, dsName, force)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) GetCertificato(c *gin.Context) {
	id := c.Param("id")
	dsName := c.DefaultQuery("ds_name", "Dirigente Scolastico")

	cert, err := h.service.GeneraCertificato(c.Request.Context(), id, dsName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cert)
}
