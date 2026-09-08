package subjects

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"registro-backend/internal/cache"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	cache   cache.Cache
}

func NewHandler(service *Service, c ...cache.Cache) *Handler {
	var appCache cache.Cache
	if len(c) > 0 {
		appCache = c[0]
	}
	return &Handler{service: service, cache: appCache}
}

func getSchoolID(c *gin.Context) string {
	role := c.GetString("role")
	if res, exists := c.Get("school_id"); exists {
		if s, ok := res.(string); ok && s != "" {
			if role != "superadmin" {
				return s
			}
		}
	}
	if q := c.Query("school_id"); q != "" {
		return q
	}
	if h := c.GetHeader("X-School-ID"); h != "" {
		return h
	}
	return c.GetString("school_id")
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

	if h.cache != nil && schoolID != "" {
		_ = h.cache.Delete(c.Request.Context(), fmt.Sprintf("subjects:school:%s", schoolID))
	}

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	schoolID := getSchoolID(c)
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required"})
		return
	}

	cacheKey := fmt.Sprintf("subjects:school:%s", schoolID)
	if h.cache != nil {
		if cached, err := h.cache.Get(c.Request.Context(), cacheKey); err == nil && cached != "" {
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
			return
		}
	}

	res, err := h.service.ListSubjects(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if h.cache != nil {
		if jsonBytes, err := json.Marshal(res); err == nil {
			_ = h.cache.Set(c.Request.Context(), cacheKey, string(jsonBytes), 30*time.Minute)
		}
	}
	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := getSchoolID(c)
	if userID == "" || (role != "admin" && role != "superadmin" && role != "secretary") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sID string
	if role != "superadmin" {
		sID = schoolID
	}
	res, err := h.service.UpdateSubject(c.Request.Context(), c.Param("id"), req, sID)
	if err != nil {
		if err.Error() == "forbidden: cannot update subject of another school" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if h.cache != nil && sID != "" {
		_ = h.cache.Delete(c.Request.Context(), fmt.Sprintf("subjects:school:%s", sID))
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
	schoolID := getSchoolID(c)
	if userID == "" || (role != "admin" && role != "superadmin" && role != "secretary") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var sID string
	if role != "superadmin" {
		sID = schoolID
	}
	if err := h.service.DeleteSubject(c.Request.Context(), c.Param("id"), sID); err != nil {
		if err.Error() == "forbidden: cannot delete subject of another school" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if h.cache != nil && sID != "" {
		_ = h.cache.Delete(c.Request.Context(), fmt.Sprintf("subjects:school:%s", sID))
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
