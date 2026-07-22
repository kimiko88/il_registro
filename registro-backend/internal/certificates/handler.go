package certificates

import (
	"fmt"
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
	group := r.Group("/certificates")
	{
		group.GET("", h.List)
		group.POST("/generate", h.Generate)
		group.GET("/:id/pdf", h.DownloadPDF)
		group.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) List(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "secretary" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	schoolID := c.GetString("school_id")
	studentID := c.Query("student_id")
	certType := CertificateType(c.Query("type"))
	year := c.Query("academic_year")

	certs, err := h.service.ListCertificates(c.Request.Context(), schoolID, studentID, certType, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if certs == nil {
		certs = []Certificate{}
	}

	c.JSON(http.StatusOK, certs)
}

func (h *Handler) Generate(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "secretary" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	var req GenerateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID := c.GetString("user_id")
	schoolID := c.GetString("school_id")

	cert, _, err := h.service.GenerateCertificate(c.Request.Context(), actorID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      cert.ID,
		"pdf_url": cert.PDFUrl,
		"cert":    cert,
	})
}

func (h *Handler) DownloadPDF(c *gin.Context) {
	id := c.Param("id")
	pdfBytes, err := h.service.GeneratePDFBytes(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificate or PDF not found"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"certificato_%s.pdf\"", id))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *Handler) Delete(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "secretary" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	id := c.Param("id")
	if err := h.service.DeleteCertificate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certificate cancelled"})
}
