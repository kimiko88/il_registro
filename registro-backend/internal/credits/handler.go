package credits

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

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/credits")
	{
		g.GET("/calculate", h.Calculate)
		g.POST("/assign", h.Assign)
		g.GET("/class/:class_id", h.ListByClass)
		g.GET("/student/:student_id/summary", h.GetStudentSummary)
	}
}

func (h *Handler) Calculate(c *gin.Context) {
	gradeLevel, _ := strconv.Atoi(c.DefaultQuery("grade_level", "3"))
	avg, _ := strconv.ParseFloat(c.DefaultQuery("average", "6.0"), 64)
	conduct, _ := strconv.Atoi(c.DefaultQuery("conduct", "8"))
	pcto, _ := strconv.Atoi(c.DefaultQuery("pcto_hours", "0"))
	hasExtra := c.DefaultQuery("has_extracurricular", "false") == "true"

	res := h.svc.CalculateSuggestedCredit(gradeLevel, avg, conduct, pcto, hasExtra)
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Assign(c *gin.Context) {
	var req AssignCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID := c.GetString("user_id")
	schoolID := c.GetString("school_id")

	credit, err := h.svc.AssignCredit(c.Request.Context(), actorID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, credit)
}

func (h *Handler) ListByClass(c *gin.Context) {
	classID := c.Param("class_id")
	academicYear := c.Query("academic_year")

	list, err := h.svc.ListClassCredits(c.Request.Context(), classID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []StudentSchoolCredit{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetStudentSummary(c *gin.Context) {
	studentID := c.Param("student_id")
	summary, err := h.svc.GetStudentSummary(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
