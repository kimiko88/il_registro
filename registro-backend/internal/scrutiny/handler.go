package scrutiny

import (
	"net/http"
	"strconv"
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
		scrutiny.GET("/matrix/:classId", h.GetMatrix)
		scrutiny.POST("/save", h.Save)
		scrutiny.POST("/class/:classId/start", h.Start)
		scrutiny.POST("/class/:classId/validate", h.Validate)
		scrutiny.POST("/class/:classId/close", h.Close)
		scrutiny.GET("/export/:studentId/pdf", h.ExportPagellaPDF)
	}
}

func (h *Handler) ExportPagellaPDF(c *gin.Context) {
	studentID := c.Param("studentId")
	classID := c.Query("class_id")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	pdfBytes, err := h.service.ExportPagellaPDF(c.Request.Context(), actorID, actorRole, classID, studentID, semester)
	if err != nil {
		if err == ErrScrutinyNotValidated {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=pagella.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *Handler) GetMatrix(c *gin.Context) {
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	classID := c.Param("classId")
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

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
		if err == ErrUnauthorizedScrutiny {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) Start(c *gin.Context) {
	classID := c.Param("classId")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	if err := h.service.StartScrutiny(c.Request.Context(), actorID, actorRole, classID, semester); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "scrutiny started successfully"})
}

func (h *Handler) Validate(c *gin.Context) {
	classID := c.Param("classId")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	if err := h.service.ValidateScrutiny(c.Request.Context(), actorID, actorRole, classID, semester); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "scrutiny validated successfully"})
}

func (h *Handler) Close(c *gin.Context) {
	classID := c.Param("classId")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	if err := h.service.CloseScrutiny(c.Request.Context(), actorID, actorRole, classID, semester); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "scrutiny closed successfully"})
}
