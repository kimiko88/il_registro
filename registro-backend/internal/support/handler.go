package support

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/support")
	{
		g.POST("/diaries", h.CreateDiaryEntry)
		g.GET("/diaries", h.ListDiaryEntries)
		g.DELETE("/diaries/:id", h.DeleteDiaryEntry)
		g.POST("/pei-goals", h.CreatePeiGoal)
		g.PATCH("/pei-goals/:id/progress", h.UpdatePeiGoalProgress)
		g.GET("/pei-goals", h.ListPeiGoals)
		g.DELETE("/pei-goals/:id", h.DeletePeiGoal)
	}
}

func (h *Handler) CreateDiaryEntry(c *gin.Context) {
	var req CreateDiaryEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	teacherID := c.GetString("teacher_id")

	entry, err := h.svc.CreateDiaryEntry(c.Request.Context(), schoolID, teacherID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *Handler) ListDiaryEntries(c *gin.Context) {
	schoolID := c.GetString("school_id")
	studentID := c.Query("student_id")
	teacherID := c.Query("teacher_id")
	classID := c.Query("class_id")

	role := c.GetString("user_role")
	isFamily := (role == "parent" || role == "student")

	list, err := h.svc.ListDiaryEntries(c.Request.Context(), schoolID, studentID, teacherID, classID, isFamily)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []SupportDiaryEntry{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) DeleteDiaryEntry(c *gin.Context) {
	id := c.Param("id")
	teacherID := c.GetString("teacher_id")

	if err := h.svc.DeleteDiaryEntry(c.Request.Context(), id, teacherID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) CreatePeiGoal(c *gin.Context) {
	var req CreatePeiGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	goal, err := h.svc.CreatePeiGoal(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

func (h *Handler) UpdatePeiGoalProgress(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		ProgressStatus string `json:"progress_status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateGoalProgress(c.Request.Context(), id, body.ProgressStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) ListPeiGoals(c *gin.Context) {
	schoolID := c.GetString("school_id")
	studentID := c.Query("student_id")

	list, err := h.svc.ListPeiGoals(c.Request.Context(), schoolID, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []SupportPeiGoal{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) DeletePeiGoal(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePeiGoal(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
