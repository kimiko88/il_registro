package colloqui

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func parseFlexibleDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" || s == "undefined" || s == "null" {
		return time.Time{}
	}
	if len(s) >= 10 {
		sDate := s[:10]
		if t, err := time.Parse("2006-01-02", sDate); err == nil {
			return t
		}
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return time.Time{}
}

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/colloqui")
	{
		g.POST("/slots", h.CreateSlot)
		g.GET("/slots", h.ListSlots)
		g.GET("/slots/my", h.ListMySlots)
		g.GET("/my-slots", h.ListMySlots)
		g.GET("/available-slots", h.ListSlots)
		g.GET("/availability/:teacherID", h.GetAvailabilityByTeacher)
		g.PATCH("/slots/:id", h.PatchSlot)
		g.DELETE("/slots/:id", h.CancelSlot)

		g.POST("/assemblies", h.CreateAssembly)

		g.POST("/bookings", h.CreateBooking)
		g.POST("/book", h.CreateBooking)
		g.GET("/bookings", h.ListMyBookings)
		g.GET("/my-bookings", h.ListMyBookings)
		g.GET("/bookings/:id", h.GetBookingByID)
		g.GET("/slots/:id/bookings", h.ListSlotBookings)
		g.PUT("/bookings/:id/status", h.UpdateBookingStatus)
		g.PATCH("/bookings/:id/cancel", h.CancelBookingAlias)

		// General Meetings & Live Virtual Queue
		g.POST("/general-meetings", h.CreateGeneralMeeting)
		g.GET("/general-meetings", h.ListGeneralMeetings)
		g.GET("/general-meetings/:id", h.GetGeneralMeeting)
		g.POST("/general-meetings/book-ticket", h.BookQueueTicket)
		g.GET("/general-meetings/:id/tickets", h.ListQueueTickets)
		g.PATCH("/general-meetings/tickets/:ticket_id/status", h.UpdateTicketStatus)
	}
}

func (h *Handler) CreateSlot(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo i docenti o amministratori possono creare slot di colloquio"})
		return
	}
	schoolID := c.GetString("school_id")

	var req CreateSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slot, err := h.service.CreateSlot(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, slot)
}

func (h *Handler) ListSlots(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	schoolID := c.GetString("school_id")
	teacherID := c.Query("teacher_id")
	available := c.Query("available") == "true"

	var fromTime, toTime time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		fromTime = parseFlexibleDate(fromStr)
	}
	if toStr := c.Query("to"); toStr != "" {
		toTime = parseFlexibleDate(toStr)
	}

	if !fromTime.IsZero() && !toTime.IsZero() {
		if toTime.Before(fromTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "la data di fine non può essere precedente alla data di inizio"})
			return
		}
		if toTime.Sub(fromTime) > 366*24*time.Hour {
			c.JSON(http.StatusBadRequest, gin.H{"error": "range di date troppo ampio (massimo 1 anno)"})
			return
		}
	}

	slots, err := h.service.ListSlots(c.Request.Context(), schoolID, teacherID, fromTime, toTime, available)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, slots)
}

func (h *Handler) CancelSlot(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	slotID := c.Param("id")

	if err := h.service.CancelSlot(c.Request.Context(), userID, role, schoolID, slotID); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

func (h *Handler) CreateBooking(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")

	if userID == "" || (role != "parent" && role != "student" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only parents or students can book slots"})
		return
	}

	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.service.BookSlot(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		if err == ErrSlotFull || err == ErrAlreadyBooked || err == ErrSlotCancelled {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *Handler) ListMyBookings(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	schoolID := c.GetString("school_id") // Bug 100: pass schoolID for tenant isolation

	bookings, err := h.service.ListMyBookings(c.Request.Context(), userID, schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *Handler) ListSlotBookings(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	slotID := c.Param("id")

	bookings, err := h.service.ListSlotBookings(c.Request.Context(), userID, role, slotID)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *Handler) UpdateBookingStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	bookingID := c.Param("id")

	var req UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateBookingStatus(c.Request.Context(), userID, role, bookingID, req); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking status updated"})
}

func (h *Handler) PatchSlot(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	slotID := c.Param("id")

	var req struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.PatchSlot(c.Request.Context(), userID, role, schoolID, slotID, req.StartTime, req.EndTime); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "slot timing patched"})
}

func (h *Handler) GetBookingByID(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	booking, err := h.service.GetBookingByID(c.Request.Context(), userID, role, id)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrBookingNotFound) || errors.Is(err, ErrSlotNotFound) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, booking)
}

func (h *Handler) CreateAssembly(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")

	if userID == "" || (role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo amministratori o dirigenti possono creare assemblee"})
		return
	}

	var req CreateAssemblyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slot, err := h.service.CreateAssembly(c.Request.Context(), role, userID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, slot)
}

func (h *Handler) ListMySlots(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Force teacher_id query param to current user ID
	q := c.Request.URL.Query()
	q.Set("teacher_id", userID)
	c.Request.URL.RawQuery = q.Encode()
	h.ListSlots(c)
}

func (h *Handler) GetAvailabilityByTeacher(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	schoolID := c.GetString("school_id")
	teacherID := c.Param("teacherID")
	if teacherID == "" {
		teacherID = c.Query("teacher_id")
	}

	var fromTime, toTime time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		fromTime = parseFlexibleDate(fromStr)
	} else {
		fromTime = time.Now()
	}
	if toStr := c.Query("to"); toStr != "" {
		toTime = parseFlexibleDate(toStr)
	} else {
		toTime = fromTime.AddDate(0, 3, 0)
	}

	// Enforce max 1-year range to prevent heavy unbounded queries.
	if toTime.Sub(fromTime) > 366*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "range di date troppo ampio (massimo 1 anno)"})
		return
	}

	slots, err := h.service.ListSlots(c.Request.Context(), schoolID, teacherID, fromTime, toTime, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, slots)
}

func (h *Handler) CancelBookingAlias(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	bookingID := c.Param("id")

	req := UpdateBookingStatusRequest{
		Status: StatusCancelled,
	}

	if err := h.service.UpdateBookingStatus(c.Request.Context(), userID, role, bookingID, req); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking status updated"})
}

func (h *Handler) CreateGeneralMeeting(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo la presidenza, segreteria o amministratori possono indire assemblee generali"})
		return
	}

	var req CreateGeneralParentMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID := c.GetString("school_id")
	m, err := h.service.CreateGeneralMeeting(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, m)
}

func (h *Handler) ListGeneralMeetings(c *gin.Context) {
	schoolID := c.GetString("school_id")
	list, err := h.service.ListGeneralMeetings(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []GeneralParentMeeting{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetGeneralMeeting(c *gin.Context) {
	id := c.Param("id")
	m, err := h.service.GetGeneralMeeting(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *Handler) BookQueueTicket(c *gin.Context) {
	var req BookQueueTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parentID := c.GetString("parent_id")
	if parentID == "" {
		parentID = c.GetString("user_id")
	}

	ticket, err := h.service.BookQueueTicket(c.Request.Context(), parentID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (h *Handler) ListQueueTickets(c *gin.Context) {
	meetingID := c.Param("id")
	teacherID := c.Query("teacher_id")
	parentID := c.Query("parent_id")

	role := c.GetString("role")
	if role == "parent" && parentID == "" {
		parentID = c.GetString("parent_id")
		if parentID == "" {
			parentID = c.GetString("user_id")
		}
	} else if role == "teacher" && teacherID == "" {
		teacherID = c.GetString("teacher_id")
		if teacherID == "" {
			teacherID = c.GetString("user_id")
		}
	}

	tickets, err := h.service.ListQueueTickets(c.Request.Context(), meetingID, teacherID, parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tickets == nil {
		tickets = []GeneralMeetingQueueTicket{}
	}
	c.JSON(http.StatusOK, tickets)
}

func (h *Handler) UpdateTicketStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: non autorizzato alla gestione dello stato del ticket"})
		return
	}

	ticketID := c.Param("ticket_id")
	var body struct {
		Status string `json:"status" binding:"required"`
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateTicketStatus(c.Request.Context(), ticketID, body.Status, body.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
