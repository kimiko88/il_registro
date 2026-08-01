package subjects

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func getSchoolID(c *gin.Context) string {
	res, exists := c.Get("school_id")
	if !exists {
		return ""
	}
	if s, ok := res.(string); ok {
		return s
	}
	return ""
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || (role != "admin" && role != "superadmin" && role != "secretary") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID := getSchoolID(c)
	if req.SchoolID != "" {
		if role == "superadmin" || req.SchoolID == schoolID {
			schoolID = req.SchoolID
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "only superadmin can specify custom school_id"})
			return
		}
	}
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required"})
		return
	}

	res, err := h.service.CreateSubject(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	schoolID := c.Query("school_id")
	if schoolID == "" {
		schoolID = getSchoolID(c)
	}
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required"})
		return
	}

	res, err := h.service.ListSubjects(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || (role != "admin" && role != "superadmin" && role != "secretary") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.UpdateSubject(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.service.GetSubject(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || (role != "admin" && role != "superadmin" && role != "secretary") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.service.DeleteSubject(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/subjects")
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
	}
}
