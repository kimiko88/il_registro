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
		timetableGroup.POST("/generate/:jobID/adjust", h.requireTimetableManagerRole, h.AdjustSchedule)
		timetableGroup.POST("/generate/:jobID/publish", h.requireTimetableManagerRole, h.PublishSchedule)

		// Preferences (Desiderata)
		timetableGroup.GET("/preferences", h.GetPreferences)
		timetableGroup.POST("/preferences", h.SavePreferences)
		timetableGroup.GET("/preferences/window", h.GetDesiderataWindow)
		timetableGroup.POST("/preferences/window", h.requireTimetableManagerRole, h.SetDesiderataWindow)

		// Room Requirements
		timetableGroup.GET("/room-requirements", h.ListRoomRequirements)
		timetableGroup.POST("/room-requirements", h.requireTimetableManagerRole, h.SaveRoomRequirement)
		timetableGroup.DELETE("/room-requirements/:id", h.requireTimetableManagerRole, h.DeleteRoomRequirement)

		// Constraints
		timetableGroup.GET("/constraints", h.ListConstraints)
		timetableGroup.POST("/constraints", h.requireTimetableManagerRole, h.SaveConstraint)
		timetableGroup.DELETE("/constraints/:id", h.requireTimetableManagerRole, h.DeleteConstraint)

		// Curriculum Plans & Class Daily Limits
		timetableGroup.GET("/academic-years", h.requireTimetableManagerRole, h.ListAcademicYears)
		timetableGroup.GET("/classes-plans", h.requireTimetableManagerRole, h.ListClassesCurriculumPlans)
		timetableGroup.GET("/classes/:classID/plan", h.requireTimetableManagerRole, h.GetClassCurriculumPlan)
		timetableGroup.PUT("/classes/:classID/plan", h.requireTimetableManagerRole, h.SaveClassCurriculumPlan)
		timetableGroup.POST("/classes/:classID/inherit", h.requireTimetableManagerRole, h.InheritClassCurriculumPlan)
		timetableGroup.POST("/inherit-all-plans", h.requireTimetableManagerRole, h.InheritAllClassesCurriculumPlans)

		// Teacher Quick Preferences (Tabular representation)
		timetableGroup.GET("/teachers-quick-preferences", h.requireTimetableManagerRole, h.GetTeachersQuickPreferences)
		timetableGroup.POST("/teachers-quick-preferences", h.requireTimetableManagerRole, h.SaveTeachersQuickPreferences)
		timetableGroup.PUT("/teachers-quick-preferences/:teacherID", h.requireTimetableManagerRole, h.SaveSingleTeacherQuickPreference)
	}
}

// requireTimetableManagerRole restricts to principal, vice_principal, collaboratore_ds, admin, superadmin, secretary
func (h *Handler) requireTimetableManagerRole(c *gin.Context) {
	role := c.GetString("role")
	if h.isManagerRole(role) {
		c.Next()
		return
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "accesso negato: operazione riservata alla dirigenza, vicario o collaboratori"})
	c.Abort()
}

func (h *Handler) isManagerRole(role string) bool {
	switch role {
	case "principal", "vice_principal", "collaboratore_ds", "admin", "superadmin", "secretary":
		return true
	default:
		return false
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
		if h.isManagerRole(role) {
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
	isManager := h.isManagerRole(role)
	if qID := c.Query("teacher_id"); qID != "" {
		if isManager {
			targetTeacherID = qID
		}
	}

	// If not manager, verify if desiderata window is open
	if !isManager {
		isOpen, err := h.service.IsDesiderataWindowOpen(c.Request.Context(), schoolID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "impossibile verificare lo stato della finestra inserimento desiderata"})
			return
		}
		if !isOpen {
			c.JSON(http.StatusForbidden, gin.H{"error": "La finestra per l'inserimento dei desiderata è attualmente chiusa. Rivolgersi al docente vicario o al responsabile orario."})
			return
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

func (h *Handler) GetDesiderataWindow(c *gin.Context) {
	schoolID := c.GetString("school_id")
	isOpen, err := h.service.IsDesiderataWindowOpen(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, DesiderataWindowResponse{IsOpen: isOpen})
}

func (h *Handler) SetDesiderataWindow(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req SetDesiderataWindowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetDesiderataWindowOpen(c.Request.Context(), schoolID, req.IsOpen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "stato finestra desiderata aggiornato con successo",
		"is_open": req.IsOpen,
	})
}

func (h *Handler) AdjustSchedule(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")
	jobID := c.Param("jobID")

	var req AdjustTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.AdjustJobSlots(c.Request.Context(), schoolID, userID, jobID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
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

// ----------------- Curriculum Plans & Class Daily Limits -----------------

func (h *Handler) ListAcademicYears(c *gin.Context) {
	schoolID := c.GetString("school_id")
	years, err := h.service.ListAcademicYears(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, years)
}

func (h *Handler) ListClassesCurriculumPlans(c *gin.Context) {
	schoolID := c.GetString("school_id")
	academicYear := c.Query("academic_year")
	plans, err := h.service.ListClassesCurriculumPlans(c.Request.Context(), schoolID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *Handler) GetClassCurriculumPlan(c *gin.Context) {
	schoolID := c.GetString("school_id")
	classID := c.Param("classID")
	plan, err := h.service.GetClassCurriculumPlan(c.Request.Context(), schoolID, classID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *Handler) SaveClassCurriculumPlan(c *gin.Context) {
	schoolID := c.GetString("school_id")
	classID := c.Param("classID")

	var req SaveClassCurriculumPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SaveClassCurriculumPlan(c.Request.Context(), schoolID, classID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedPlan, err := h.service.GetClassCurriculumPlan(c.Request.Context(), schoolID, classID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "piano orario salvato con successo"})
		return
	}
	c.JSON(http.StatusOK, updatedPlan)
}

func (h *Handler) InheritClassCurriculumPlan(c *gin.Context) {
	schoolID := c.GetString("school_id")
	classID := c.Param("classID")

	var req InheritPlanRequest
	_ = c.ShouldBindJSON(&req) // Optional body

	if req.SourceAcademicYear == "" {
		req.SourceAcademicYear = c.Query("source_academic_year")
	}

	plan, err := h.service.InheritClassCurriculumPlan(c.Request.Context(), schoolID, classID, req.SourceAcademicYear)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *Handler) InheritAllClassesCurriculumPlans(c *gin.Context) {
	schoolID := c.GetString("school_id")

	var req InheritPlanRequest
	_ = c.ShouldBindJSON(&req)

	if req.SourceAcademicYear == "" {
		req.SourceAcademicYear = c.Query("source_academic_year")
	}

	result, err := h.service.InheritAllClassesCurriculumPlans(c.Request.Context(), schoolID, req.SourceAcademicYear)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ----------------- Teacher Quick Preferences Handlers -----------------

func (h *Handler) GetTeachersQuickPreferences(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var academicYearID *string
	if ayID := c.Query("academic_year_id"); ayID != "" {
		academicYearID = &ayID
	}

	res, err := h.service.GetTeachersQuickPreferences(c.Request.Context(), schoolID, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) SaveTeachersQuickPreferences(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req SaveTeacherQuickPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SaveTeacherQuickPreferences(c.Request.Context(), schoolID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "preferenze rapide docenti salvate con successo"})
}

func (h *Handler) SaveSingleTeacherQuickPreference(c *gin.Context) {
	schoolID := c.GetString("school_id")
	teacherID := c.Param("teacherID")
	if teacherID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacherID obbligatorio"})
		return
	}

	var item TeacherQuickPreferenceItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item.TeacherID = teacherID

	var academicYearID *string
	if ayID := c.Query("academic_year_id"); ayID != "" {
		academicYearID = &ayID
	}

	req := SaveTeacherQuickPreferencesRequest{
		AcademicYearID: academicYearID,
		Preferences:    []TeacherQuickPreferenceItem{item},
	}

	if err := h.service.SaveTeacherQuickPreferences(c.Request.Context(), schoolID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "preferenze rapide docente salvate con successo"})
}
