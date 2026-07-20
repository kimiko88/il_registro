package attendance

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	att := r.Group("/attendance")

	// Teacher
	att.POST("/mark", h.MarkAttendance)
	att.POST("/mark-bulk", h.MarkBulk)
	att.GET("/class/:id", h.GetClassAttendance)
	att.GET("/pending-justifications", h.GetPendingJustifications)
	att.POST("/justification/:id/process", h.ProcessJustification)

	// Student/Parent
	att.GET("/my-attendance", h.GetMyAttendance)
	att.GET("/my-attendance/summary", h.GetMySummary)
	att.POST("/justify", h.RequestJustification)

	att.GET("/child-attendance/:studentID", h.GetChildAttendance)
	att.GET("/child-attendance/:studentID/summary", h.GetChildSummary)

	// Admin / Secretary
	att.POST("/justification/:id/approve", h.ApproveJustification)
	att.DELETE("/justification/:id", h.RejectJustification)
	att.GET("/analytics", h.GetAnalytics)
}

// parseWindowParams legge i query param from/to; se assenti usa l'intero anno scolastico corrente.
func parseWindowParams(c *gin.Context) (from, to time.Time, err error) {
	const layout = "2006-01-02"
	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr != "" {
		from, err = time.Parse(layout, fromStr)
		if err != nil {
			return
		}
	} else {
		// Default: inizio anno scolastico (1 settembre dell'anno corrente o precedente).
		now := time.Now()
		year := now.Year()
		if now.Month() < time.September {
			year--
		}
		from = time.Date(year, time.September, 1, 0, 0, 0, 0, time.UTC)
	}
	if toStr != "" {
		to, err = time.Parse(layout, toStr)
		if err != nil {
			return
		}
	} else {
		to = time.Now()
	}
	return
}

func (h *Handler) GetMySummary(c *gin.Context) {
	studentID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	res, err := h.service.GetStudentSummary(c.Request.Context(), studentID, schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetChildAttendance(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	studentID := c.Param("studentID")
	from, to, err := parseWindowParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parametro data non valido: " + err.Error()})
		return
	}
	res, err := h.service.GetChildAttendance(c.Request.Context(), parentID, studentID, from, to)
	if err != nil {
		if err.Error() == "unauthorized: not a guardian of this student" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetChildSummary(c *gin.Context) {
	parentID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	studentID := c.Param("studentID")
	res, err := h.service.GetChildSummary(c.Request.Context(), parentID, studentID, schoolID)
	if err != nil {
		if err.Error() == "unauthorized: not a guardian of this student" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ApproveJustification requires teacher or admin role.
func (h *Handler) ApproveJustification(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" || (actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id := c.Param("id")
	if err := h.service.ProcessJustification(c.Request.Context(), actorID, id, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *Handler) RejectJustification(c *gin.Context) {
	id := c.Param("id")
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" || (actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err := h.service.ProcessJustification(c.Request.Context(), actorID, id, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id mancante nel token"})
		return
	}
	res, err := h.service.GetSchoolAnalytics(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) MarkAttendance(c *gin.Context) {
	var req CreateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.MarkAttendance(c.Request.Context(), userID, schoolID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked"})
}

func (h *Handler) MarkBulk(c *gin.Context) {
	var req BulkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.MarkBulk(c.Request.Context(), userID, schoolID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bulk marked"})
}

// GetClassAttendance requires teacher, admin, or superadmin role.
func (h *Handler) GetClassAttendance(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	classID := c.Param("id")
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date required"})
		return
	}
	res, err := h.service.GetClassAttendance(c.Request.Context(), classID, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetMyAttendance(c *gin.Context) {
	studentID := c.GetString("user_id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	from, to, err := parseWindowParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parametro data non valido: " + err.Error()})
		return
	}
	res, err := h.service.GetStudentAttendance(c.Request.Context(), studentID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) RequestJustification(c *gin.Context) {
	var req JustificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.RequestJustification(c.Request.Context(), parentID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "requested"})
}

func (h *Handler) GetPendingJustifications(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" || (actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	classID := c.Query("class_id")
	res, err := h.service.GetPendingJustifications(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) ProcessJustification(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Approve bool `json:"approve"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	teacherID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if teacherID == "" || (actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err := h.service.ProcessJustification(c.Request.Context(), teacherID, id, req.Approve); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "processed"})
}

func (h *Handler) GetAttendance(c *gin.Context) { h.GetMyAttendance(c) }
