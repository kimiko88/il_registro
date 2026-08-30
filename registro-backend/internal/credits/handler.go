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
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gradeLevel, _ := strconv.Atoi(c.DefaultQuery("grade_level", "3"))
	avg, _ := strconv.ParseFloat(c.DefaultQuery("average", "6.0"), 64)
	conduct, _ := strconv.Atoi(c.DefaultQuery("conduct", "8"))
	pcto, _ := strconv.Atoi(c.DefaultQuery("pcto_hours", "0"))
	hasExtra := c.DefaultQuery("has_extracurricular", "false") == "true"

	res := h.svc.CalculateSuggestedCredit(gradeLevel, avg, conduct, pcto, hasExtra)
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Assign(c *gin.Context) {
	actorID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "coordinator" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: assegnazione crediti scolastici riservata al consiglio di classe e alla dirigenza"})
		return
	}

	var req AssignCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	credit, err := h.svc.AssignCredit(c.Request.Context(), actorID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, credit)
}

func (h *Handler) ListByClass(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" && role != "coordinator" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: consultazione crediti di classe riservata al personale scolastico"})
		return
	}

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
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	studentID := c.Param("student_id")
	if role == "student" {
		if userID != studentID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: non puoi consultare i crediti di un altro studente"})
			return
		}
	} else if role == "parent" {
		if uRepo := h.svc.GetUserRepo(); uRepo != nil {
			isG, err := uRepo.IsGuardian(c.Request.Context(), userID, studentID)
			if err != nil || !isG {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: non sei tutore legale di questo studente"})
				return
			}
		}
	} else if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" && role != "coordinator" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	summary, err := h.svc.GetStudentSummary(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

