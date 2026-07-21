package reports

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	}
}

func (h *Handler) ExportGradesExcel(c *gin.Context) {
	classID := c.Query("class_id")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

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
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}
