package colloqui

import (
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
		g.GET("/slots/my", h.ListSlots)
		g.GET("/my-slots", h.ListSlots)
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
	}
}

func (h *Handler) CreateSlot(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	if teacherID == "" && (strings.HasSuffix(c.Request.URL.Path, "/slots/my") || strings.Contains(c.Request.URL.Path, "my-slots")) {
		teacherID = userID
	}
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
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") || err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func (h *Handler) GetAvailabilityByTeacher(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
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
