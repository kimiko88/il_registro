package attendance

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

func respond500(c *gin.Context, msg string, err error) {
	logger.Log.Errorf("%s: %v", msg, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

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
	att.PUT("/:id", h.UpdateAttendance)
	att.GET("/class/:id", h.GetClassAttendance)
	att.DELETE("/class/:id/hour/:hour", h.DeleteClassAttendanceHour)
	att.GET("/pending-justifications", h.GetPendingJustifications)
	att.POST("/justification/:id/process", h.ProcessJustification)
	att.POST("/justification/:id/reject", h.RejectJustification)
	att.GET("/export", h.ExportAttendance)

	// Student/Parent
	att.GET("/my-attendance", h.GetMyAttendance)
	att.GET("/my-attendance/summary", h.GetMySummary)
	att.POST("/justify", h.RequestJustification)

	att.GET("/child-attendance/:studentID", h.GetChildAttendance)
	att.GET("/child-attendance/:studentID/summary", h.GetChildSummary)
	att.GET("/child-attendance/:studentID/trends", h.GetChildAttendanceTrends)

	att.GET("/child/:studentID/unjustified", h.GetChildUnjustified)
	att.POST("/child/:studentID/justify/:attendanceID", h.JustifyChildAbsence)
	att.GET("/child/:studentID/stats", h.GetChildAttendanceStats)

	// Admin / Secretary
	att.POST("/justification/:id/approve", h.ApproveJustification)
	att.GET("/analytics", h.GetAnalytics)

	// Monthly Breakdown
	att.GET("/students/:studentID/monthly-breakdown", h.GetMonthlyBreakdown)
	att.GET("/child-attendance/:studentID/monthly-breakdown", h.GetChildMonthlyBreakdown)

	// Teacher: student summary
	att.GET("/students/:studentID/summary", h.GetStudentSummaryForTeacher)
}

// parseWindowParams legge i query param from/to; se assenti usa l'intero anno scolastico corrente.
func parseWindowParams(c *gin.Context) (from, to time.Time, err error) {
	const layout = "2006-01-02"
	loc, lErr := time.LoadLocation("Europe/Rome")
	if lErr != nil {
		loc = time.Local
	}
	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr != "" {
		from, err = time.ParseInLocation(layout, fromStr, loc)
		if err != nil {
			return
		}
	} else {
		now := time.Now().In(loc)
		year := now.Year()
		if now.Month() < time.September {
			year--
		}
		from = time.Date(year, time.September, 1, 0, 0, 0, 0, loc)
	}
	if toStr != "" {
		to, err = time.ParseInLocation(layout, toStr, loc)
		if err != nil {
			return
		}
	} else {
		to = time.Now().In(loc)
	}

	if to.Before(from) {
		err = errors.New("data di inizio successiva alla data di fine")
		return
	}
	if to.After(from.AddDate(1, 3, 0)) {
		err = errors.New("range di date troppo ampio (massimo 1 anno e 3 mesi consentito)")
		return
	}
	return
}

func (h *Handler) GetMySummary(c *gin.Context) {
	studentID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only students can access my-summary (parents must use /attendance/students/:studentID/summary)"})
		return
	}
	res, err := h.service.GetStudentSummary(c.Request.Context(), studentID, role, schoolID, studentID)
	if err != nil {
		respond500(c, "GetMySummary error", err)
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
		respond500(c, "GetChildAttendance error", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetChildSummary(c *gin.Context) {
	parentID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "parent" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: access restricted to parents or administrators"})
		return
	}
	studentID := c.Param("studentID")
	res, err := h.service.GetChildSummary(c.Request.Context(), parentID, studentID, schoolID)
	if err != nil {
		if err.Error() == "unauthorized: not a guardian of this student" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetChildSummary error", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetStudentSummaryForTeacher allows a teacher or staff member to fetch the attendance summary of
// a student in their school.
func (h *Handler) GetStudentSummaryForTeacher(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	studentID := c.Param("studentID")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID mancante"})
		return
	}
	res, err := h.service.GetStudentSummary(c.Request.Context(), actorID, actorRole, schoolID, studentID)
	if err != nil {
		respond500(c, "GetStudentSummaryForTeacher error", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetChildAttendanceTrends(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	studentID := c.Param("studentID")
	res, err := h.service.GetChildAttendanceTrends(c.Request.Context(), parentID, studentID)
	if err != nil {
		if err.Error() == "unauthorized: not a guardian of this student" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetChildAttendanceTrends error", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func handleJustificationError(c *gin.Context, err error) {
	errStr := strings.ToLower(err.Error())
	if errors.Is(err, ErrAlreadyProcessed) || strings.Contains(errStr, "già stata elaborata") || strings.Contains(errStr, "already processed") {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, ErrForbidden) || strings.Contains(errStr, "forbidden") || strings.Contains(errStr, "non assegnato") {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, ErrUnauthorized) || strings.Contains(errStr, "unauthorized") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, ErrNotFound) || strings.Contains(errStr, "not found") || strings.Contains(errStr, "non trovata") {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	respond500(c, "Justification processing error", err)
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
	if err := h.service.ProcessJustification(c.Request.Context(), actorID, actorRole, id, true); err != nil {
		handleJustificationError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

// RejectJustification requires teacher or admin role.
func (h *Handler) RejectJustification(c *gin.Context) {
	id := c.Param("id")
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" || (actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err := h.service.ProcessJustification(c.Request.Context(), actorID, actorRole, id, false); err != nil {
		handleJustificationError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" && role != "teacher" {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only teachers or admins can mark attendance"})
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
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only teachers or admins can mark attendance"})
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
	schoolID := c.GetString("school_id")
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
	res, err := h.service.GetClassAttendance(c.Request.Context(), actorID, actorRole, schoolID, classID, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetMyAttendance(c *gin.Context) {
	studentID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only students can view my-attendance (parents must use /child-attendance/:studentID)"})
		return
	}
	from, to, err := parseWindowParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parametro data non valido: " + err.Error()})
		return
	}
	res, err := h.service.GetStudentAttendance(c.Request.Context(), studentID, role, schoolID, studentID, from, to)
	if err != nil {
		respond500(c, "GetMyAttendance error", err)
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
	role := c.GetString("role")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "parent" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo i genitori possono richiedere giustificazioni"})
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
	schoolID := c.GetString("school_id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	classID := c.Query("class_id")
	res, err := h.service.GetPendingJustifications(c.Request.Context(), actorID, actorRole, classID, schoolID)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") || strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetPendingJustifications error", err)
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
	if err := h.service.ProcessJustification(c.Request.Context(), teacherID, actorRole, id, req.Approve); err != nil {
		handleJustificationError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "processed"})
}

func (h *Handler) UpdateAttendance(c *gin.Context) {
	id := c.Param("id")
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")

	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only teachers or admins can update attendance"})
		return
	}

	var req UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateAttendance(c.Request.Context(), teacherID, schoolID, id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "attendance updated"})
}

func (h *Handler) ExportAttendance(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	classID := c.Query("class_id")
	dateRaw := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	dateParsed, err := time.Parse("2006-01-02", dateRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date parameter format: expected YYYY-MM-DD"})
		return
	}
	date := dateParsed.Format("2006-01-02")

	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_id parameter required"})
		return
	}

	res, err := h.service.GetClassAttendance(c.Request.Context(), actorID, actorRole, schoolID, classID, date)
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "forbidden") || strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "non sei assegnato") || strings.Contains(errStr, "non appartiene") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(errStr, "non trovata") || strings.Contains(errStr, "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "ExportAttendance GetClassAttendance error", err)
		return
	}

	// Sanitize classID to prevent header injection: strip quotes, CR, LF, semicolons.
	safeClassID := strings.Map(func(r rune) rune {
		if r == '"' || r == '\r' || r == '\n' || r == ';' || r == '\\' {
			return '_'
		}
		return r
	}, classID)
	filename := fmt.Sprintf("presenze_%s_%s.csv", safeClassID, date)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"StudentID", "Date", "Hour", "Status", "IsJustified", "EntryTime", "ExitTime", "Notes"})

	for _, rec := range res.Records {
		studentID := rec.StudentID
		if len(studentID) > 0 && (studentID[0] == '=' || studentID[0] == '+' || studentID[0] == '-' || studentID[0] == '@') {
			studentID = "'" + studentID
		}
		statusStr := string(rec.Status)
		if len(statusStr) > 0 && (statusStr[0] == '=' || statusStr[0] == '+' || statusStr[0] == '-' || statusStr[0] == '@') {
			statusStr = "'" + statusStr
		}
		notes := rec.Notes
		if len(notes) > 0 && (notes[0] == '=' || notes[0] == '+' || notes[0] == '-' || notes[0] == '@') {
			notes = "'" + notes
		}
		hourStr := ""
		if rec.Hour > 0 {
			hourStr = fmt.Sprintf("%d", rec.Hour)
		}
		_ = w.Write([]string{
			studentID,
			rec.Date,
			hourStr,
			statusStr,
			fmt.Sprintf("%t", rec.IsJustified),
			rec.EntryTime,
			rec.ExitTime,
			notes,
		})
	}
	w.Flush()

	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// GetMonthlyBreakdown returns per-month attendance statistics for a student.
func (h *Handler) GetMonthlyBreakdown(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	studentID := c.Param("studentID")
	schoolYear := c.DefaultQuery("school_year", "")

	if role != "teacher" && role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" && role != "student" && role != "parent" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if role == "student" && userID != studentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: student can only access own monthly breakdown"})
		return
	}

	if role == "parent" {
		res, err := h.service.GetChildMonthlyBreakdown(c.Request.Context(), userID, studentID, schoolYear)
		if err != nil {
			if strings.Contains(err.Error(), "guardian") || strings.Contains(err.Error(), "access denied") {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			respond500(c, "GetMonthlyBreakdown parent error", err)
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	schoolID := c.GetString("school_id")
	res, err := h.service.GetMonthlyBreakdown(c.Request.Context(), userID, role, schoolID, studentID, schoolYear)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetMonthlyBreakdown error", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetChildMonthlyBreakdown returns per-month attendance statistics for a parent's child.
// Maintained as an explicit legacy alias route delegating directly to GetMonthlyBreakdown for backward compatibility with mobile/frontend clients.
func (h *Handler) GetChildMonthlyBreakdown(c *gin.Context) {
	h.GetMonthlyBreakdown(c)
}

func (h *Handler) GetChildUnjustified(c *gin.Context) {
	parentID := c.GetString("user_id")
	role := c.GetString("role")
	studentID := c.Param("studentID")
	if parentID == "" || studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}
	if role != "parent" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo i genitori o gli amministratori possono accedere a questa funzione"})
		return
	}

	result, err := h.service.GetChildUnjustified(c.Request.Context(), parentID, studentID)
	if err != nil {
		errLower := strings.ToLower(err.Error())
		if strings.Contains(errLower, "guardian") || strings.Contains(errLower, "tutela") || strings.Contains(errLower, "forbidden") || strings.Contains(errLower, "unauthorized") || strings.Contains(errLower, "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetChildUnjustified error", err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) JustifyChildAbsence(c *gin.Context) {
	parentID := c.GetString("user_id")
	role := c.GetString("role")
	studentID := c.Param("studentID")
	attendanceID := c.Param("attendanceID")
	if parentID == "" || studentID == "" || attendanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}
	if role != "parent" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo i genitori o gli amministratori possono giustificare le assenze del figlio"})
		return
	}

	var req JustifyAbsenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo della richiesta non valido: " + err.Error()})
		return
	}

	err := h.service.JustifyChildAbsence(c.Request.Context(), parentID, studentID, attendanceID, req)
	if err != nil {
		errLower := strings.ToLower(err.Error())
		if strings.Contains(errLower, "guardian") || strings.Contains(errLower, "tutela") || strings.Contains(errLower, "forbidden") || strings.Contains(errLower, "unauthorized") || strings.Contains(errLower, "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "JustifyChildAbsence error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "absence justified successfully"})
}

func (h *Handler) GetChildAttendanceStats(c *gin.Context) {
	parentID := c.GetString("user_id")
	role := c.GetString("role")
	studentID := c.Param("studentID")
	if parentID == "" || studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}
	if role != "parent" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo i genitori o gli amministratori possono accedere alle statistiche del figlio"})
		return
	}

	stats, err := h.service.GetChildAttendanceStats(c.Request.Context(), parentID, studentID)
	if err != nil {
		errLower := strings.ToLower(err.Error())
		if strings.Contains(errLower, "guardian") || strings.Contains(errLower, "tutela") || strings.Contains(errLower, "forbidden") || strings.Contains(errLower, "unauthorized") || strings.Contains(errLower, "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetChildAttendanceStats error", err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) DeleteClassAttendanceHour(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	classID := c.Param("id")
	hourStr := c.Param("hour")
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter required"})
		return
	}

	var hourInt int
	if _, convErr := fmt.Sscanf(hourStr, "%d", &hourInt); convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hour parameter"})
		return
	}

	schoolID := c.GetString("school_id")

	if err := h.service.DeleteClassAttendanceHour(c.Request.Context(), actorID, actorRole, schoolID, classID, dateStr, hourInt); err != nil {
		errStr := strings.ToLower(err.Error())
		if errors.Is(err, ErrForbidden) || strings.Contains(errStr, "forbidden") || strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "non sei assegnato") || strings.Contains(errStr, "non assegnato") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "DeleteClassAttendanceHour error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance for hour deleted"})
}
