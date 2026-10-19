package homework_submissions

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

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/homework-submissions")
	{
		group.POST("", h.SubmitHomework)
		group.GET("/:id", h.GetSubmission)
		group.GET("/by-homework", h.ListByHomework)
		group.GET("/my-submissions", h.ListMySubmissions)
		group.POST("/:id/grade", h.GradeSubmission)
		group.POST("/:id/sync-classroom", h.SyncClassroom)
	}
}

func (h *Handler) SubmitHomework(c *gin.Context) {
	studentID := c.GetString("user_id")
	var sub HomeworkSubmission
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub.StudentID = studentID

	if err := h.svc.SubmitHomework(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sub)
}

func (h *Handler) GetSubmission(c *gin.Context) {
	id := c.Param("id")
	sub, err := h.svc.GetSubmission(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "consegna non trovata"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

func (h *Handler) ListByHomework(c *gin.Context) {
	hwID := c.Query("homework_id")
	list, err := h.svc.ListHomeworkSubmissions(hwID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"submissions": list})
}

func (h *Handler) ListMySubmissions(c *gin.Context) {
	studentID := c.GetString("user_id")
	list, err := h.svc.ListStudentSubmissions(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"submissions": list})
}

func (h *Handler) GradeSubmission(c *gin.Context) {
	id := c.Param("id")
	teacherID := c.GetString("user_id")

	var req GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.GradeSubmission(id, req, teacherID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "elaborato corretto e restituito con successo"})
}

func (h *Handler) SyncClassroom(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ClassroomSubmissionID string `json:"google_classroom_submission_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.SyncGoogleClassroom(id, req.ClassroomSubmissionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sincronizzazione Google Classroom completata"})
}
