package timetablegen

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	timetableGroup := r.Group("/timetable")
	{
		// Generation Job
		timetableGroup.POST("/generate", h.requireTimetableManagerRole, h.StartGeneration)
		timetableGroup.GET("/generate/:jobID", h.GetJobStatus)
		timetableGroup.POST("/generate/:jobID/publish", h.requireTimetableManagerRole, h.PublishSchedule)

		// Preferences (Desiderata)
		timetableGroup.GET("/preferences", h.GetPreferences)
		timetableGroup.POST("/preferences", h.SavePreferences)

		// Room Requirements
		timetableGroup.GET("/room-requirements", h.ListRoomRequirements)
		timetableGroup.POST("/room-requirements", h.requireTimetableManagerRole, h.SaveRoomRequirement)
		timetableGroup.DELETE("/room-requirements/:id", h.requireTimetableManagerRole, h.DeleteRoomRequirement)

		// Constraints
		timetableGroup.GET("/constraints", h.ListConstraints)
		timetableGroup.POST("/constraints", h.requireTimetableManagerRole, h.SaveConstraint)
		timetableGroup.DELETE("/constraints/:id", h.requireTimetableManagerRole, h.DeleteConstraint)
	}
}

// requireTimetableManagerRole restricts to principal, vice_principal, collaboratore_ds, admin, superadmin
func (h *Handler) requireTimetableManagerRole(c *gin.Context) {
	role := c.GetString("role")
	switch role {
	case "principal", "vice_principal", "collaboratore_ds", "admin", "superadmin":
		c.Next()
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso negato: operazione riservata alla dirigenza, vicario o collaboratori"})
		c.Abort()
	}
}

// ----------------- Generation Handlers -----------------

func (h *Handler) StartGeneration(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")

	var req GenerateTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobID, err := h.service.StartGeneration(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"job_id":  jobID,
		"message": "calcolo orario scolastico avviato con successo",
		"status":  JobStatusPending,
	})
}

func (h *Handler) GetJobStatus(c *gin.Context) {
	schoolID := c.GetString("school_id")
	jobID := c.Param("jobID")

	job, err := h.service.GetJobStatus(c.Request.Context(), schoolID, jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *Handler) PublishSchedule(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")
	jobID := c.Param("jobID")

	if err := h.service.PublishSchedule(c.Request.Context(), schoolID, userID, jobID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "orario scolastico generato e pubblicato con successo in class_schedules",
	})
}

// ----------------- Preferences Handlers -----------------

func (h *Handler) GetPreferences(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	targetTeacherID := userID
	if qID := c.Query("teacher_id"); qID != "" {
		if role == "principal" || role == "vice_principal" || role == "collaboratore_ds" || role == "admin" || role == "superadmin" || role == "secretary" {
			targetTeacherID = qID
		}
	}

	var academicYearID *string
	if ayID := c.Query("academic_year_id"); ayID != "" {
		academicYearID = &ayID
	}

	prefs, err := h.service.GetTeacherPreferences(c.Request.Context(), schoolID, targetTeacherID, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prefs)
}

func (h *Handler) SavePreferences(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	targetTeacherID := userID
	if qID := c.Query("teacher_id"); qID != "" {
		if role == "principal" || role == "vice_principal" || role == "collaboratore_ds" || role == "admin" || role == "superadmin" {
			targetTeacherID = qID
		}
	}

	var req SavePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SaveTeacherPreferences(c.Request.Context(), schoolID, targetTeacherID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preferenze desiderata orario salvate con successo"})
}

// ----------------- Room Requirements Handlers -----------------

func (h *Handler) ListRoomRequirements(c *gin.Context) {
	schoolID := c.GetString("school_id")
	reqs, err := h.service.ListRoomRequirements(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reqs)
}

func (h *Handler) SaveRoomRequirement(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req SaveRoomRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	saved, err := h.service.SaveRoomRequirement(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, saved)
}

func (h *Handler) DeleteRoomRequirement(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteRoomRequirement(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "requisito aula rimosso con successo"})
}

// ----------------- Constraints Handlers -----------------

func (h *Handler) ListConstraints(c *gin.Context) {
	schoolID := c.GetString("school_id")
	list, err := h.service.ListConstraints(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) SaveConstraint(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req SaveConstraintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	saved, err := h.service.SaveConstraint(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, saved)
}

func (h *Handler) DeleteConstraint(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteConstraint(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "vincolo orario rimosso con successo"})
}
