package helpdesk

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	hd := rg.Group("/help-desk")
	{
		hd.POST("/slots", h.CreateSlot)
		hd.GET("/slots", h.ListSlots)
		hd.POST("/slots/:id/book", h.BookSlot)
		hd.GET("/slots/:id/bookings", h.ListBookings)
		hd.PUT("/bookings/:id/attendance", h.MarkAttendance)
		hd.PUT("/slots/:id/complete", h.CompleteSlot)
		hd.GET("/fis-report", h.GetFISReport)
	}
}

func (h *Handler) CreateSlot(c *gin.Context) {
	role := c.GetString("role")
	switch role {
	case "admin", "superadmin", "secretary", "principal", "vice_principal", "teacher", "coordinator":
		// Allowed
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso non autorizzato"})
		return
	}

	var req struct {
		SubjectID   string  `json:"subject_id" binding:"required"`
		SlotDate    string  `json:"slot_date" binding:"required"`
		StartTime   string  `json:"start_time" binding:"required"`
		EndTime     string  `json:"end_time" binding:"required"`
		RoomID      *string `json:"room_id"`
		MaxCapacity int     `json:"max_capacity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dati slot non validi: " + err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}
	teacherID := c.GetString("user_id")

	slot := &HelpDeskSlot{
		SchoolID:    schoolID,
		TeacherID:   teacherID,
		SubjectID:   req.SubjectID,
		SlotDate:    req.SlotDate,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		RoomID:      req.RoomID,
		MaxCapacity: req.MaxCapacity,
	}

	if err := h.service.CreateSlot(c.Request.Context(), slot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Disponibilità sportello help registrata",
		"slot":    slot,
	})
}

func (h *Handler) ListSlots(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}
	subjectID := c.Query("subject_id")
	date := c.Query("date")
	status := c.Query("status")

	slots, err := h.service.ListSlots(c.Request.Context(), schoolID, subjectID, date, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": slots,
	})
}

func (h *Handler) BookSlot(c *gin.Context) {
	slotID := c.Param("id")
	studentID := c.GetString("user_id")
	if studentID == "" {
		studentID = c.Query("student_id")
	}

	var req struct {
		TopicDescription string `json:"topic_description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "specificare l'argomento/dubbio: " + err.Error()})
		return
	}

	booking, err := h.service.BookSlot(c.Request.Context(), slotID, studentID, req.TopicDescription)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prenotazione effettuata con successo",
		"booking": booking,
	})
}

func (h *Handler) ListBookings(c *gin.Context) {
	slotID := c.Param("id")
	bookings, err := h.service.ListBookings(c.Request.Context(), slotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookings,
	})
}

func (h *Handler) MarkAttendance(c *gin.Context) {
	bookingID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stato obbligatorio"})
		return
	}

	if err := h.service.MarkAttendance(c.Request.Context(), bookingID, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Presenza aggiornata",
	})
}

func (h *Handler) CompleteSlot(c *gin.Context) {
	slotID := c.Param("id")
	if err := h.service.CompleteSlot(c.Request.Context(), slotID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sportello concluso e validato per la rendicontazione FIS",
	})
}

func (h *Handler) GetFISReport(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}

	report, err := h.service.GetFISAccountingReport(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": report,
	})
}
