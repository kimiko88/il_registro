package reports

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"registro-backend/pkg/upload"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	rep := r.Group("/reports")
	{
		rep.GET("/grades/excel", h.ExportGradesExcel)
		rep.GET("/sidi/students", h.ExportSidiStudentsXML)
		rep.GET("/sidi/scrutini", h.ExportSidiScrutiniXML)
		rep.GET("/sidi/attendance", h.ExportSidiAttendanceCSV)
	}
}

func (h *Handler) ExportGradesExcel(c *gin.Context) {
	classID := c.Query("class_id")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "teacher" && actorRole != "principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_id parameter required"})
		return
	}

	excelBytes, err := h.service.ExportGradesExcel(c.Request.Context(), actorID, actorRole, classID, semester)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("matrice_voti_%s_q%d_%s.xlsx", classID, semester, time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}

func (h *Handler) ExportSidiStudentsXML(c *gin.Context) {
	classID := c.Query("class_id")
	xmlBytes, err := h.service.ExportSidiStudentsXML(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("sidi_anagrafe_alunni_%s.xml", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/xml")
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, "application/xml", xmlBytes)
}

func (h *Handler) ExportSidiScrutiniXML(c *gin.Context) {
	classID := c.Query("class_id")
	sem, _ := strconv.Atoi(c.DefaultQuery("semester", "2"))
	xmlBytes, err := h.service.ExportSidiScrutiniXML(c.Request.Context(), classID, sem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("sidi_scrutini_q%d_%s.xml", sem, time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/xml")
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, "application/xml", xmlBytes)
}

func (h *Handler) ExportSidiAttendanceCSV(c *gin.Context) {
	classID := c.Query("class_id")
	csvBytes, err := h.service.ExportSidiAttendanceCSV(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("sidi_assenze_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, "text/csv", csvBytes)
}

