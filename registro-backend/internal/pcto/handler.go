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

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	pcto := r.Group("/pcto")

	// Teacher
	pcto.POST("/projects", h.CreateProject)
	pcto.GET("/projects", h.GetProjects)
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
	if err := h.service.CreateProject(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *Handler) GetProjects(c *gin.Context) {
	res, err := h.service.GetProjects(c.Request.Context())
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
	if err := h.service.AssignStudent(c.Request.Context(), id, req.StudentID); err != nil {
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
	if err := h.service.CreateCompany(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "company created"})
}

func (h *Handler) GetCompanies(c *gin.Context) {
	res, err := h.service.GetCompanies(c.Request.Context())
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
