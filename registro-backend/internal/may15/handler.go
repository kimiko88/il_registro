package may15

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

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	may15Group := rg.Group("/may15")
	{
		may15Group.GET("/class/:classId", h.GetDocument)
		may15Group.PUT("/class/:classId", h.SaveDocument)
		may15Group.POST("/class/:classId/publish", h.PublishDocument)
	}
}

func (h *Handler) GetDocument(c *gin.Context) {
	classID := c.Param("classId")
	academicYear := c.Query("academic_year")
	if academicYear == "" {
		academicYear = "2025/2026"
	}

	doc, err := h.service.GetDocument(c.Request.Context(), classID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doc)
}

func (h *Handler) SaveDocument(c *gin.Context) {
	role := c.GetString("role")
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "coordinator" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso riservato al personale docente e amministrativo"})
		return
	}

	schoolID := c.GetString("school_id")
	classID := c.Param("classId")

	var req SaveMay15Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.service.SaveDocument(c.Request.Context(), schoolID, classID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doc)
}

func (h *Handler) PublishDocument(c *gin.Context) {
	role := c.GetString("role")
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "coordinator" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso riservato al consiglio di classe e presidenza"})
		return
	}

	classID := c.Param("classId")
	academicYear := c.Query("academic_year")
	if academicYear == "" {
		academicYear = "2025/2026"
	}

	doc, err := h.service.PublishDocument(c.Request.Context(), classID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Documento del 15 Maggio approvato e pubblicato con successo",
		"document": doc,
	})
}
