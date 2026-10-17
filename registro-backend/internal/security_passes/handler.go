package security_passes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/security-passes")
	{
		group.POST("/hall-passes", h.IssueHallPass)
		group.POST("/hall-passes/:id/return", h.ReturnHallPass)
		group.GET("/hall-passes/active", h.GetActivePasses)
		group.GET("/hall-passes/by-class", h.GetClassPasses)

		group.POST("/delegates", h.AddPickupDelegate)
		group.GET("/delegates/by-student", h.GetStudentDelegates)
		group.POST("/delegates/:id/revoke", h.RevokePickupDelegate)
	}
}

func (h *Handler) IssueHallPass(c *gin.Context) {
	schoolID := c.GetString("school_id")
	teacherID := c.GetString("user_id")

	var p StudentHallPass
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.SchoolID = schoolID
	p.TeacherID = teacherID

	if err := h.svc.IssueHallPass(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) ReturnHallPass(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ReturnHallPass(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "studente rientrato in aula"})
}

func (h *Handler) GetActivePasses(c *gin.Context) {
	schoolID := c.GetString("school_id")
	passes, err := h.svc.GetActivePasses(schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passes": passes})
}

func (h *Handler) GetClassPasses(c *gin.Context) {
	classID := c.Query("class_id")
	limitStr := c.DefaultQuery("limit", "30")
	limit, _ := strconv.Atoi(limitStr)

	passes, err := h.svc.GetClassPasses(classID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passes": passes})
}

func (h *Handler) AddPickupDelegate(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var d StudentPickupDelegate
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d.SchoolID = schoolID

	if err := h.svc.AddPickupDelegate(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *Handler) GetStudentDelegates(c *gin.Context) {
	studentID := c.Query("student_id")
	delegates, err := h.svc.GetStudentDelegates(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"delegates": delegates})
}

func (h *Handler) RevokePickupDelegate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RevokePickupDelegate(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "delega revocata con successo"})
}
