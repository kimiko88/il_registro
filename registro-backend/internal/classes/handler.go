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

// Helper to get school ID from user context (set by auth middleware) -> usually stored as "school_id" in JWT/Context
func getSchoolID(c *gin.Context) string {
	res, exists := c.Get("school_id")
	if !exists {
		return ""
	}
	return res.(string)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Check if SchoolID is provided in request (SuperAdmin overriding)
	schoolID := req.SchoolID

	// 2. If not in request, get from context (Standard Admin/Teacher flow)
	if schoolID == "" {
		schoolID = getSchoolID(c)
	}

	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required (either in body or context)"})
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
	// 1. Try to get school_id from query (for admins/superadmins viewing specific school)
	schoolID := c.Query("school_id")

	// 2. If not in query, get from context (current user's school)
	if schoolID == "" {
		schoolID = getSchoolID(c)
	}

	// 3. If still empty, strictly require it?
	// For superadmin listing ALL classes globally, we might allow empty.
	// But service likely expects a schoolID. Let's check service later.
	// For now, if empty, return error as before to be safe.
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id parameter or context required"})
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

func (h *Handler) Update(c *gin.Context) {
	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class, err := h.service.UpdateClass(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.DeleteClass(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetTeacherClasses(c *gin.Context) {
	// Need user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		// Try 'userID' (some middleware uses this) or 'sub' or check auth middleware
		// Looking at auth middleware (step 287), it uses "userID" and "role".
		// But in step 298 (auth/middleware.go), it uses "user_id".
		// Main.go uses auth.NewMiddleware from internal/auth.
		// So checking internal/auth/middleware.go (step 298), it sets "user_id".
		// Wait, classes/handler.go (step 269) used:
		// res, exists := c.Get("school_id")

		_, ok := c.Get("userID") // Try legacy key if strictly AuthMiddleware from middleware/auth.go?
		// But main.go uses internal/auth. So "user_id" is correct key from internal/auth/middleware.go
		if !exists && !ok {
			// Try fetching "user_id"
			userID, exists = c.Get("user_id")
		}
	}
	// internal/auth/middleware sets "user_id".
	// Let's assume "user_id".
	userID, exists = c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	classes, err := h.service.GetTeacherClasses(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *Handler) AssignSubject(c *gin.Context) {
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

		// Assignments
		group.POST("/:id/subjects", h.AssignSubject)
		group.GET("/:id/subjects", h.GetClassSubjects)
		group.DELETE("/:id/subjects/:assignmentId", h.RemoveSubject)
	}

	// Teacher specific routes
	rg.GET("/teacher/classes", h.GetTeacherClasses)
}
