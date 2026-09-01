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
	role := c.GetString("role")
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo i docenti e la dirigenza possono registrare voci nel diario di sostegno"})
		return
	}

	var req CreateDiaryEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	teacherID := c.GetString("teacher_id")
	if teacherID == "" {
		teacherID = c.GetString("user_id")
	}

	entry, err := h.svc.CreateDiaryEntry(c.Request.Context(), schoolID, teacherID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *Handler) ListDiaryEntries(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	studentID := c.Query("student_id")
	teacherID := c.Query("teacher_id")
	classID := c.Query("class_id")

	role := c.GetString("role")
	if role == "" {
		role = c.GetString("user_role")
	}
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	isFamily := (role == "parent" || role == "student")
	switch role {
	case "student":
		if studentID != "" && studentID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: non puoi consultare i diari di altri studenti"})
			return
		}
		studentID = userID
	case "parent":
		if studentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student_id obbligatorio per i genitori"})
			return
		}
		if uRepo := h.svc.GetUserRepo(); uRepo != nil {
			isG, err := uRepo.IsGuardian(c.Request.Context(), userID, studentID)
			if err != nil || !isG {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: non sei tutore legale di questo studente"})
				return
			}
		}
	}

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
	role := c.GetString("role")
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	teacherID := c.GetString("teacher_id")
	if teacherID == "" {
		teacherID = c.GetString("user_id")
	}

	if err := h.svc.DeleteDiaryEntry(c.Request.Context(), id, teacherID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) CreatePeiGoal(c *gin.Context) {
	role := c.GetString("role")
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo i docenti e la dirigenza possono creare obiettivi PEI"})
		return
	}

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
	role := c.GetString("role")
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo i docenti e la dirigenza possono aggiornare il progresso PEI"})
		return
	}

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
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	studentID := c.Query("student_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	switch role {
	case "student":
		if studentID != "" && studentID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: non puoi visualizzare gli obiettivi PEI di altri studenti"})
			return
		}
		studentID = userID
	case "parent":
		if studentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student_id obbligatorio per i genitori"})
			return
		}
		if uRepo := h.svc.GetUserRepo(); uRepo != nil {
			isG, err := uRepo.IsGuardian(c.Request.Context(), userID, studentID)
			if err != nil || !isG {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: non sei tutore legale di questo studente"})
				return
			}
		}
	}

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
	role := c.GetString("role")
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	if err := h.svc.DeletePeiGoal(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
