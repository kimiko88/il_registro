package pcto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func getSchoolID(c *gin.Context) string {
	res, exists := c.Get("school_id")
	if !exists {
		return ""
	}
	return res.(string)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	pcto := r.Group("/pcto")

	// Teacher
	pcto.POST("/projects", h.CreateProject)
	pcto.GET("/projects", h.GetProjects)
	pcto.PUT("/projects/:id", h.UpdateProject)
	pcto.DELETE("/projects/:id", h.DeleteProject)
	pcto.POST("/projects/:id/students", h.AssignStudent)

	// Companies
	pcto.POST("/companies", h.CreateCompany)
	pcto.GET("/companies", h.GetCompanies)

	// Student
	pcto.GET("/my-projects", h.GetMyProjects)
	pcto.POST("/hours", h.LogHours)
	pcto.GET("/my-projects/:id", h.GetProjectDetails)
}

func (h *Handler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	role := c.GetString("role")
	schoolID := getSchoolID(c)
	if err := h.service.CreateProject(c.Request.Context(), schoolID, role, userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *Handler) GetProjects(c *gin.Context) {
	schoolID := getSchoolID(c)
	res, err := h.service.GetProjects(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) AssignStudent(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		StudentID string `json:"student_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := c.GetString("role")
	if err := h.service.AssignStudent(c.Request.Context(), role, id, req.StudentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

func (h *Handler) CreateCompany(c *gin.Context) {
	var req Company
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := c.GetString("role")
	schoolID := getSchoolID(c)
	if err := h.service.CreateCompany(c.Request.Context(), schoolID, role, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "company created"})
}

func (h *Handler) GetCompanies(c *gin.Context) {
	schoolID := getSchoolID(c)
	res, err := h.service.GetCompanies(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetMyProjects(c *gin.Context) {
	userID := c.GetString("userID")
	res, err := h.service.GetMyProjects(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) LogHours(c *gin.Context) {
	var req LogHourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	if err := h.service.LogHours(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "logged"})
}

func (h *Handler) GetProjectDetails(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.GetString("userID")
	part, logs, err := h.service.GetMyProjectDetails(c.Request.Context(), userID, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"participation": part, "logs": logs})
}
func (h *Handler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := c.GetString("role")
	if err := h.service.UpdateProject(c.Request.Context(), role, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	role := c.GetString("role")
	if err := h.service.DeleteProject(c.Request.Context(), role, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
