package substitutions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/substitutions")
	{
		g.POST("", h.Create)
		g.GET("", h.ListBySchool)
		g.GET("/my", h.ListByTeacher)
		g.GET("/my-today", h.ListMyToday)
		g.GET("/today-summary", h.TodaySummary)
		g.PUT("/:id/assign", h.AssignSubstitute)
		g.PATCH("/:id/confirm", h.Confirm)
		g.POST("/:id/sign-register", h.SignRegister)
		g.GET("/recommend-substitutes", h.RecommendSubstitutes)
	}
}

func isStaffRole(role string, c *gin.Context) bool {
	if role == "teacher" || role == "admin" || role == "superadmin" || role == "secretary" || role == "principal" || role == "vice_principal" {
		return true
	}
	return c.GetBool("is_staff")
}

func isManagementRole(role string, c *gin.Context) bool {
	if role == "admin" || role == "superadmin" || role == "secretary" || role == "principal" || role == "vice_principal" || role == "collaboratore_ds" || role == "dsga" {
		return true
	}
	return c.GetBool("is_staff")
}

func (h *Handler) Create(c *gin.Context) {
	if c.GetString("user_id") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if !isManagementRole(role, c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateSubstitutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	effectiveRole := role
	if c.GetBool("is_staff") {
		effectiveRole = "staff"
	}

	sub, err := h.service.CreateSubstitution(c.Request.Context(), effectiveRole, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sub)
}

func (h *Handler) ListBySchool(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	isStaff := isStaffRole(role, c)
	if !isStaff {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	schoolID := c.GetString("school_id")
	date := c.Query("date")

	subs, err := h.service.ListBySchool(c.Request.Context(), schoolID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *Handler) ListByTeacher(c *gin.Context) {
	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	date := c.Query("date")

	subs, err := h.service.ListByTeacher(c.Request.Context(), teacherID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *Handler) ListMyToday(c *gin.Context) {
	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subs, err := h.service.ListMyToday(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *Handler) AssignSubstitute(c *gin.Context) {
	if c.GetString("user_id") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	isStaff := isStaffRole(role, c)
	if !isStaff {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	var req AssignSubstituteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	effectiveRole := role
	if isStaff {
		effectiveRole = "staff"
	}

	if err := h.service.AssignSubstitute(c.Request.Context(), id, effectiveRole, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "substitute assigned"})
}

func (h *Handler) Confirm(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	if err := h.service.ConfirmSubstitution(c.Request.Context(), id, teacherID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "substitution confirmed"})
}

func (h *Handler) SignRegister(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id := c.Param("id")
	var body struct {
		Notes string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&body)

	if err := h.service.SignRegister(c.Request.Context(), id, teacherID, body.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registro supplenze firmato con successo"})
}

func (h *Handler) RecommendSubstitutes(c *gin.Context) {
	uid := c.GetString("user_id")
	role := c.GetString("role")
	isStaff := isStaffRole(role, c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isStaff {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	schoolID := c.GetString("school_id")
	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")
	date := c.Query("date")
	hourStr := c.Query("hour")
	hour, _ := strconv.Atoi(hourStr)

	recs, err := h.service.RecommendSubstitutes(c.Request.Context(), schoolID, classID, subjectID, date, hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recs)
}

// TodaySummary restituisce un riepilogo rapido delle sostituzioni per il giorno corrente.
// Usato dal pannello "Emergenza Sostituzioni" del Collaboratore DS.
func (h *Handler) TodaySummary(c *gin.Context) {
	uid := c.GetString("user_id")
	role := c.GetString("role")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !isManagementRole(role, c) && !isStaffRole(role, c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	schoolID := c.GetString("school_id")
	today := c.Query("date") // optional override; se vuoto usa today
	if today == "" {
		today = "" // service ListBySchool con date="" restituisce tutto il giorno corrente
	}

	subs, err := h.service.ListBySchool(c.Request.Context(), schoolID, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calcola conteggi per il pannello emergency
	summary := gin.H{
		"total":     len(subs),
		"pending":   0,
		"assigned":  0,
		"confirmed": 0,
		"cancelled": 0,
	}
	for _, s := range subs {
		switch s.Status {
		case StatusPending:
			summary["pending"] = summary["pending"].(int) + 1
		case StatusAssigned:
			summary["assigned"] = summary["assigned"].(int) + 1
		case StatusConfirmed:
			summary["confirmed"] = summary["confirmed"].(int) + 1
		case StatusCancelled:
			summary["cancelled"] = summary["cancelled"].(int) + 1
		}
	}
	summary["substitutions"] = subs

	c.JSON(http.StatusOK, summary)
}
