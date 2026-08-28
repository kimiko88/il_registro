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
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
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
	actorID := c.GetString("user_id")
	role := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "secretary" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	var req GenerateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "secretary" && role != "superadmin" && role != "student" && role != "parent" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	cert, err := h.service.GetCertificateByID(c.Request.Context(), id)
	if err != nil || cert == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificate not found"})
		return
	}
	if role != "superadmin" && schoolID != "" && cert.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: certificate belongs to another school"})
		return
	}
	if role == "student" && cert.StudentID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: cannot access certificate of another student"})
		return
	}
	if role == "parent" {
		if cert.StudentID != userID {
			if uRepo := h.service.GetUserRepo(); uRepo != nil {
				isG, err := uRepo.IsGuardian(c.Request.Context(), userID, cert.StudentID)
				if err != nil || !isG {
					c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: cannot access certificate of another student"})
					return
				}
			}
		}
	}

	pdfBytes, err := h.service.GeneratePDFBytes(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"certificato_%s.pdf\"", id))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *Handler) Delete(c *gin.Context) {
	actorID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "secretary" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	id := c.Param("id")
	cert, err := h.service.GetCertificateByID(c.Request.Context(), id)
	if err != nil || cert == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificate not found"})
		return
	}
	if role != "superadmin" && schoolID != "" && cert.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: cannot delete certificate of another school"})
		return
	}

	if err := h.service.DeleteCertificate(c.Request.Context(), id, actorID, role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
