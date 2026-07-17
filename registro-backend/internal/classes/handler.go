package classes

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
	v, ok := res.(string)
	if !ok {
		return ""
	}
	return v
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Always use school_id from the JWT context; ignore any school_id in the request body
	// to prevent IDOR attacks where an admin creates classes in another school.
	schoolID := getSchoolID(c)
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id not found in token"})
		return
	}
	class, err := h.service.CreateClass(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, class)
}

func (h *Handler) List(c *gin.Context) {
	// For non-superadmin roles, always scope to the caller's school.
	schoolID := getSchoolID(c)
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id not found in token"})
		return
	}
	academicYear := c.Query("academic_year")
	classes, err := h.service.ListClasses(c.Request.Context(), schoolID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *Handler) Get(c *gin.Context) {
	class, err := h.service.GetClass(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

// Update verifies the class belongs to the caller's school before updating.
func (h *Handler) Update(c *gin.Context) {
	schoolID := getSchoolID(c)
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existing, err := h.service.GetClass(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}
	if existing.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	class, err := h.service.UpdateClass(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

// Delete verifies the class belongs to the caller's school before deleting.
func (h *Handler) Delete(c *gin.Context) {
	schoolID := getSchoolID(c)
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	existing, err := h.service.GetClass(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}
	if existing.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err := h.service.DeleteClass(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetTeacherClasses(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	classes, err := h.service.GetTeacherClasses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *Handler) AssignSubject(c *gin.Context) {
	schoolID := getSchoolID(c)
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	existing, err := h.service.GetClass(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}
	if existing.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var req AssignSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AssignSubject(c.Request.Context(), c.Param("id"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) GetClassSubjects(c *gin.Context) {
	res, err := h.service.GetClassSubjects(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) RemoveSubject(c *gin.Context) {
	if err := h.service.RemoveSubject(c.Request.Context(), c.Param("assignmentId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/classes")
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)

		group.POST("/:id/subjects", h.AssignSubject)
		group.GET("/:id/subjects", h.GetClassSubjects)
		group.DELETE("/:id/subjects/:assignmentId", h.RemoveSubject)
	}
	rg.GET("/teacher/classes", h.GetTeacherClasses)
}
