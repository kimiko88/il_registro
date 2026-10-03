package primaryeval

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/primary")
	{
		group.GET("/objectives", h.ListObjectives)
		group.POST("/objectives", h.CreateObjective)
		group.DELETE("/objectives/:id", h.DeleteObjective)

		group.GET("/evaluations", h.ListEvaluations)
		group.POST("/evaluations/batch", h.SaveEvaluationsBatch)
		group.GET("/matrix", h.GetMatrix)
		group.GET("/student/:studentId", h.GetStudentEvaluations)
	}
}

func (h *Handler) ListObjectives(c *gin.Context) {
	schoolID := c.GetString("school_id")
	subjectID := c.Query("subject_id")
	classIDParam := c.Query("class_id")
	var classID *string
	if classIDParam != "" {
		classID = &classIDParam
	}
	yearGrade, _ := strconv.Atoi(c.Query("year_grade"))
	academicYear := c.Query("academic_year")

	objectives, err := h.service.ListObjectives(c.Request.Context(), schoolID, subjectID, classID, yearGrade, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": "PRIMARY_OBJECTIVES_FETCH_FAILED"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": objectives})
}

func (h *Handler) CreateObjective(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req CreateObjectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "INVALID_OBJECTIVE_PAYLOAD"})
		return
	}

	obj, err := h.service.CreateObjective(c.Request.Context(), schoolID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "CREATE_OBJECTIVE_FAILED"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": obj})
}

func (h *Handler) DeleteObjective(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteObjective(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": "DELETE_OBJECTIVE_FAILED"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Obiettivo didattico eliminato con successo"})
}

func (h *Handler) ListEvaluations(c *gin.Context) {
	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")
	semester, _ := strconv.Atoi(c.Query("semester"))

	evals, err := h.service.ListEvaluations(c.Request.Context(), classID, subjectID, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": "PRIMARY_EVALUATIONS_FETCH_FAILED"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": evals})
}

func (h *Handler) SaveEvaluationsBatch(c *gin.Context) {
	schoolID := c.GetString("school_id")
	teacherID := c.GetString("user_id")

	var req SaveEvaluationsBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "INVALID_PRIMARY_EVAL_PAYLOAD"})
		return
	}

	if err := h.service.SaveEvaluationsBatch(c.Request.Context(), schoolID, teacherID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "SAVE_PRIMARY_EVAL_FAILED"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Valutazioni descrittive registrate con successo"})
}

func (h *Handler) GetMatrix(c *gin.Context) {
	schoolID := c.GetString("school_id")
	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")
	semester, _ := strconv.Atoi(c.Query("semester"))
	if semester <= 0 {
		semester = 1
	}

	matrix, err := h.service.GetPrimaryMatrix(c.Request.Context(), schoolID, classID, subjectID, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": "MATRIX_FETCH_FAILED"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": matrix})
}

func (h *Handler) GetStudentEvaluations(c *gin.Context) {
	studentID := c.Param("studentId")
	semester, _ := strconv.Atoi(c.Query("semester"))

	evals, err := h.service.GetStudentEvaluations(c.Request.Context(), studentID, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": "STUDENT_EVALUATIONS_FETCH_FAILED"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": evals})
}
