package scrutiny

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	pkgLogger "registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	scrutiny := rg.Group("/scrutiny")
	{
		scrutiny.GET("/overview", h.GetOverview)
		scrutiny.GET("/matrix/:classId", h.GetMatrix)
		scrutiny.GET("/class/:classId/report", h.GetClassReport)
		scrutiny.POST("/class/:classId/finalize", h.FinalizeClass)
		scrutiny.GET("/export", h.ExportAll)
		scrutiny.POST("/save", h.Save)
		scrutiny.POST("/class/:classId/start", h.Start)
		scrutiny.POST("/class/:classId/validate", h.Validate)
		scrutiny.POST("/class/:classId/close", h.Close)
		scrutiny.GET("/export/:studentId/pdf", h.ExportPagellaPDF)
	}
}

func parseSemester(semStr string) int {
	if s, err := strconv.Atoi(semStr); err == nil && s > 0 {
		return s
	}
	semLower := strings.ToLower(semStr)
	if strings.Contains(semLower, "2") || strings.Contains(semLower, "second") {
		return 2
	}
	return 1
}

func (h *Handler) ExportPagellaPDF(c *gin.Context) {
	studentID := c.Param("studentId")
	classID := c.Query("class_id")
	semester := parseSemester(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	pdfBytes, err := h.service.ExportPagellaPDF(c.Request.Context(), actorID, actorRole, classID, studentID, semester)
	if err != nil {
		if err == ErrScrutinyNotValidated {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("pagella_%s_semestre%d.pdf", studentID, semester)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *Handler) GetMatrix(c *gin.Context) {
	semester := parseSemester(c.DefaultQuery("semester", "1"))
	classID := c.Param("classId")
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	matrix, err := h.service.GetMatrix(c.Request.Context(), actorID, actorRole, classID, semester)
	if err != nil {
		if err == ErrScrutinyNotValidated {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		pkgLogger.Log.Error("failed to get scrutiny matrix", "error", err, "classId", classID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matrix)
}

func (h *Handler) Save(c *gin.Context) {
	var req SaveScrutinyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coordinatorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if err := h.service.SaveScrutiny(c.Request.Context(), coordinatorID, actorRole, req); err != nil {
		if err == ErrUnauthorizedScrutiny || strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "closed") || strings.Contains(err.Error(), "validated") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) Start(c *gin.Context) {
	classID := c.Param("classId")
	semester := parseSemester(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	if err := h.service.StartScrutiny(c.Request.Context(), actorID, actorRole, classID, semester); err != nil {
		if err == ErrUnauthorizedScrutiny || strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "already") || strings.Contains(err.Error(), "closed") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "scrutiny started successfully"})
}

func (h *Handler) Validate(c *gin.Context) {
	classID := c.Param("classId")
	semester := parseSemester(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	if err := h.service.ValidateScrutiny(c.Request.Context(), actorID, actorRole, classID, semester); err != nil {
		if err == ErrUnauthorizedScrutiny || strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "already") || strings.Contains(err.Error(), "closed") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "scrutiny validated successfully"})
}

func (h *Handler) Close(c *gin.Context) {
	classID := c.Param("classId")
	semester := parseSemester(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	if err := h.service.CloseScrutiny(c.Request.Context(), actorID, actorRole, classID, semester); err != nil {
		if err == ErrUnauthorizedScrutiny || strings.HasPrefix(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "already") || strings.Contains(err.Error(), "closed") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "scrutiny closed successfully"})
}

func (h *Handler) GetOverview(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
		return
	}

	overview, err := h.service.GetOverview(c.Request.Context(), actorID, actorRole, schoolID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unauthorized") || strings.HasPrefix(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, overview)
}

func (h *Handler) GetClassReport(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "secretary" && actorRole != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
		return
	}

	classID := c.Param("classId")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	report, err := h.service.GetClassReport(c.Request.Context(), actorID, actorRole, classID, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *Handler) FinalizeClass(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
		return
	}

	classID := c.Param("classId")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "2"))
	if err := h.service.FinalizeClass(c.Request.Context(), actorID, actorRole, classID, semester); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "scrutiny finalized successfully"})
}

func (h *Handler) ExportAll(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "2"))
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
		return
	}

	data, err := h.service.ExportAll(c.Request.Context(), actorID, actorRole, schoolID, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=\"scrutini_overview.csv\"")
	c.Data(http.StatusOK, "text/csv", data)
}
