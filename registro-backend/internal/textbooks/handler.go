package textbooks

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	textbooks := rg.Group("/textbooks")
	{
		textbooks.POST("", h.Create)
		textbooks.PUT("/:id", h.Update)
		textbooks.GET("", h.List)
		textbooks.DELETE("/:id", h.Delete)
		textbooks.GET("/class/:classId", h.ListByClass)
		textbooks.POST("/class/:classId", h.AssignToClass)
		textbooks.DELETE("/class/assignment/:id", h.RemoveFromClass)

		// AIE & Spending Limits endpoints
		textbooks.POST("/aie/import", h.ImportAIE)
		textbooks.GET("/aie/catalog", h.SearchAIECatalog)
		textbooks.GET("/classes/:classId/spending-report", h.GetClassSpendingReport)
		textbooks.GET("/classes/:classId/adoptions", h.ListClassAdoptions)
		textbooks.POST("/classes/:classId/adoptions", h.AdoptBook)
		textbooks.DELETE("/adoptions/:id", h.DeleteAdoption)
		textbooks.POST("/spending-limits", h.UpsertSpendingLimit)
		textbooks.GET("/classes/:classId/aie-export", h.ExportClassAIE)
	}
}

func isTextbookAuthorizedRole(role string) bool {
	switch role {
	case "teacher", "coordinator", "coordinatore_classe", "admin", "superadmin", "secretary", "principal", "vice_principal", "assistente_alunni", "assistente_amministrativo":
		return true
	default:
		return false
	}
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isTextbookAuthorizedRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateTextbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID := c.GetString("school_id")
	if err := h.service.CreateTextbook(c.Request.Context(), schoolID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isTextbookAuthorizedRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	var req CreateTextbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID := c.GetString("school_id")
	if err := h.service.UpdateTextbook(c.Request.Context(), role, schoolID, id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if role == "superadmin" {
		if requestedSchoolID := c.Query("school_id"); requestedSchoolID != "" {
			schoolID = requestedSchoolID
		}
	}

	list, err := h.service.ListTextbooks(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isTextbookAuthorizedRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	schoolID := c.GetString("school_id")
	if err := h.service.DeleteTextbook(c.Request.Context(), role, schoolID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) AssignToClass(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isTextbookAuthorizedRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	classID := c.Param("classId")
	var req AssignTextbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AssignToClass(c.Request.Context(), classID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) RemoveFromClass(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isTextbookAuthorizedRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	if err := h.service.RemoveFromClass(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListByClass(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	classID := c.Param("classId")
	list, err := h.service.ListByClass(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// AIE & Spending Limits Handler Methods

func (h *Handler) ImportAIE(c *gin.Context) {
	role := c.GetString("role")
	if !isTextbookAuthorizedRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required (form-data: 'file')"})
		return
	}
	defer func() { _ = file.Close() }()

	imported, err := h.service.ImportAIECatalog(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("import error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Catalogo AIE importato con successo",
		"imported": imported,
	})
}

func (h *Handler) SearchAIECatalog(c *gin.Context) {
	queryStr := c.Query("q")
	subject := c.Query("subject")
	schoolOrder := c.Query("school_order")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	books, err := h.service.SearchAIECatalog(c.Request.Context(), queryStr, subject, schoolOrder, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *Handler) GetClassSpendingReport(c *gin.Context) {
	classID := c.Param("classId")
	schoolID := c.GetString("school_id")

	report, err := h.service.GetClassSpendingReport(c.Request.Context(), schoolID, classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *Handler) ListClassAdoptions(c *gin.Context) {
	classID := c.Param("classId")
	adoptions, err := h.service.repo.ListClassAdoptions(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, adoptions)
}

func (h *Handler) AdoptBook(c *gin.Context) {
	role := c.GetString("role")
	if !isTextbookAuthorizedRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	classID := c.Param("classId")
	var req ClassAdoptionItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ClassID = classID

	if err := h.service.AdoptBook(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) DeleteAdoption(c *gin.Context) {
	role := c.GetString("role")
	if !isTextbookAuthorizedRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	if err := h.service.DeleteAdoption(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UpsertSpendingLimit(c *gin.Context) {
	role := c.GetString("role")
	if !isTextbookAuthorizedRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req SpendingLimit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	if err := h.service.UpsertSpendingLimit(c.Request.Context(), schoolID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tetto di spesa salvato con successo"})
}

func (h *Handler) ExportClassAIE(c *gin.Context) {
	classID := c.Param("classId")
	txt, err := h.service.ExportClassAIE(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=adozioni_aie_%s.txt", classID))
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(txt))
}
