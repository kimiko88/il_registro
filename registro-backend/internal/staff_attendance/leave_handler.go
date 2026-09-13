package staff_attendance

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LeaveHandler gestisce le API per richieste ferie/permessi del personale ATA
type LeaveHandler struct {
	repo Repository
}

func NewLeaveHandler(repo Repository) *LeaveHandler {
	return &LeaveHandler{repo: repo}
}

func (h *LeaveHandler) RegisterLeaveRoutes(r *gin.RouterGroup) {
	g := r.Group("/staff-attendance")
	{
		// Cartellino mensile
		g.GET("/timecard", h.GetTimecard)
		g.GET("/timecard/export", h.ExportTimecard)

		// Richieste ferie/permessi
		g.POST("/leaves", h.CreateLeave)
		g.GET("/leaves", h.ListLeaves)
		g.GET("/leaves/:id", h.GetLeave)
		g.PATCH("/leaves/:id/approve", h.ApproveLeave)
		g.PATCH("/leaves/:id/reject", h.RejectLeave)
		g.DELETE("/leaves/:id", h.DeleteLeave)
	}
}

func leaveAuthCheck(c *gin.Context) (userID, schoolID, role string, ok bool) {
	userID = c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return "", "", "", false
	}
	schoolID = c.GetString("school_id")
	role = c.GetString("role")
	if !CanReadAttendance(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso non autorizzato"})
		return "", "", "", false
	}
	return userID, schoolID, role, true
}

// GetTimecard restituisce il cartellino mensile
// Query params: user_id (opzionale, solo per DSGA/admin), month (YYYY-MM, default mese corrente)
func (h *LeaveHandler) GetTimecard(c *gin.Context) {
	userID, schoolID, role, ok := leaveAuthCheck(c)
	if !ok {
		return
	}

	month := c.Query("month")
	if month == "" {
		month = ""
	}

	// Se richiesto riepilogo globale (solo DSGA, Dirigente, Admin)
	if c.Query("all") == "true" {
		if role != "dsga" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
			c.JSON(http.StatusForbidden, gin.H{"error": "accesso non autorizzato al riepilogo complessivo"})
			return
		}
		timecards, err := h.repo.GetAllMonthlyTimecards(c.Request.Context(), schoolID, month)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, timecards)
		return
	}

	targetUserID := c.Query("user_id")

	// Solo DSGA/admin/principal possono vedere il cartellino altrui
	if targetUserID != "" && targetUserID != userID {
		if role != "dsga" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" {
			c.JSON(http.StatusForbidden, gin.H{"error": "puoi visualizzare solo il tuo cartellino"})
			return
		}
	} else {
		targetUserID = userID
	}

	timecard, err := h.repo.GetMonthlyTimecard(c.Request.Context(), schoolID, targetUserID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, timecard)
}

// ExportTimecard esporta il riepilogo mensile in CSV
func (h *LeaveHandler) ExportTimecard(c *gin.Context) {
	_, schoolID, role, ok := leaveAuthCheck(c)
	if !ok {
		return
	}
	if role != "dsga" && role != "admin" && role != "superadmin" && role != "principal" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo il DSGA può esportare il riepilogo"})
		return
	}

	month := c.Query("month")
	summaries, err := h.repo.GetAllMonthlyTimecards(c.Request.Context(), schoolID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"cartellino_"+month+".csv\"")

	// BOM UTF-8 per Excel italiano
	c.Writer.Write([]byte("\xef\xbb\xbf"))
	c.Writer.WriteString("Cognome,Nome,Ruolo,Mese,Ore Contratto,Ore Lavorate,Straordinari,Giorni Assenza,Giorni Ferie,Giorni Malattia,Ore Permesso\n")
	for _, t := range summaries {
		c.Writer.WriteString(
			t.LastName + "," + t.FirstName + "," + t.Role + "," + t.Month + "," +
				formatFloat(t.ContractHours) + "," + formatFloat(t.WorkedHours) + "," +
				formatFloat(t.OvertimeHours) + "," + itoa(t.AbsenceDays) + "," +
				itoa(t.LeaveDays) + "," + itoa(t.SickDays) + "," + formatFloat(t.PermitHours) + "\n",
		)
	}
}

func (h *LeaveHandler) CreateLeave(c *gin.Context) {
	userID, schoolID, _, ok := leaveAuthCheck(c)
	if !ok {
		return
	}
	var req CreateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	leave, err := h.repo.CreateLeaveRequest(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, leave)
}

func (h *LeaveHandler) ListLeaves(c *gin.Context) {
	userID, schoolID, role, ok := leaveAuthCheck(c)
	if !ok {
		return
	}

	// DSGA/admin vedono tutte, gli altri solo le proprie
	targetUserID := ""
	if role != "dsga" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" {
		targetUserID = userID
	} else {
		targetUserID = c.Query("user_id") // opzionale: filtra per utente specifico
	}

	leaves, err := h.repo.ListLeaveRequests(c.Request.Context(), schoolID, targetUserID, c.Query("status"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaves)
}

func (h *LeaveHandler) GetLeave(c *gin.Context) {
	userID, schoolID, role, ok := leaveAuthCheck(c)
	if !ok {
		return
	}
	id := c.Param("id")
	leave, err := h.repo.GetLeaveRequest(c.Request.Context(), schoolID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "richiesta non trovata"})
		return
	}
	// Verifica accesso: solo la propria richiesta o ruoli privilegiati
	if leave.UserID != userID && role != "dsga" && role != "admin" && role != "superadmin" && role != "principal" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso negato"})
		return
	}
	c.JSON(http.StatusOK, leave)
}

func (h *LeaveHandler) ApproveLeave(c *gin.Context) {
	userID, schoolID, role, ok := leaveAuthCheck(c)
	if !ok {
		return
	}
	if role != "dsga" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo il DSGA o il Dirigente possono approvare le richieste"})
		return
	}
	id := c.Param("id")
	var req ApproveLeaveRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.repo.ApproveLeaveRequest(c.Request.Context(), schoolID, id, userID, req.Notes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "richiesta approvata"})
}

func (h *LeaveHandler) RejectLeave(c *gin.Context) {
	userID, schoolID, role, ok := leaveAuthCheck(c)
	if !ok {
		return
	}
	if role != "dsga" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo il DSGA o il Dirigente possono rifiutare le richieste"})
		return
	}
	id := c.Param("id")
	var req RejectLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.RejectLeaveRequest(c.Request.Context(), schoolID, id, userID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "richiesta rifiutata"})
}

func (h *LeaveHandler) DeleteLeave(c *gin.Context) {
	userID, schoolID, _, ok := leaveAuthCheck(c)
	if !ok {
		return
	}
	id := c.Param("id")
	if err := h.repo.DeleteLeaveRequest(c.Request.Context(), schoolID, id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "richiesta eliminata"})
}

// Helper per CSV export
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
