package groups

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

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	groups := rg.Group("/groups")
	{
		groups.POST("", h.CreateGroup)
		groups.GET("", h.ListGroups)
		groups.GET("/:id", h.GetGroup)
		groups.PUT("/:id", h.UpdateGroup)
		groups.DELETE("/:id", h.DeleteGroup)
		groups.POST("/:id/students", h.AddStudents)
		groups.DELETE("/:id/students/:student_id", h.RemoveStudent)
	}
}

func (h *Handler) CreateGroup(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || (role != "secretary" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "message": err.Error()})
		return
	}

	group, err := h.service.CreateGroup(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *Handler) ListGroups(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	schoolID := c.Query("school_id")
	teacherID := c.Query("teacher_id")
	studentID := c.Query("student_id")

	var groups []Group
	var err error

	if schoolID != "" {
		groups, err = h.service.ListGroupsBySchool(c.Request.Context(), schoolID)
	} else if teacherID != "" {
		groups, err = h.service.ListGroupsByTeacher(c.Request.Context(), teacherID)
	} else if studentID != "" {
		groups, err = h.service.ListGroupsByStudent(c.Request.Context(), studentID)
	} else {
		sID := c.GetString("school_id")
		if sID != "" {
			groups, err = h.service.ListGroupsBySchool(c.Request.Context(), sID)
		} else {
			groups = []Group{}
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (h *Handler) GetGroup(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	group, err := h.service.GetGroupByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *Handler) UpdateGroup(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || (role != "secretary" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "message": err.Error()})
		return
	}

	group, err := h.service.UpdateGroup(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *Handler) DeleteGroup(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || (role != "secretary" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	if err := h.service.DeleteGroup(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "group deleted successfully"})
}

func (h *Handler) AddStudents(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || (role != "secretary" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	var req AddGroupStudentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "message": err.Error()})
		return
	}

	if err := h.service.AddStudentsToGroup(c.Request.Context(), id, req.StudentIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "students added successfully"})
}

func (h *Handler) RemoveStudent(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || (role != "secretary" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	studentID := c.Param("student_id")

	if err := h.service.RemoveStudentFromGroup(c.Request.Context(), id, studentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "student removed successfully"})
}
