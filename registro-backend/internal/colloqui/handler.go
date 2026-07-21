package colloqui

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/colloqui")
	{
		g.POST("/slots", h.CreateSlot)
		g.GET("/slots", h.ListSlots)
		g.PATCH("/slots/:id", h.PatchSlot)
		g.DELETE("/slots/:id", h.CancelSlot)

		g.POST("/bookings", h.CreateBooking)
		g.GET("/my-bookings", h.ListMyBookings)
		g.GET("/bookings/:id", h.GetBookingByID)
		g.GET("/slots/:id/bookings", h.ListSlotBookings)
		g.PUT("/bookings/:id/status", h.UpdateBookingStatus)
	}
}

func (h *Handler) CreateSlot(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")

	if userID == "" || (role != "teacher" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

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
		var err error
		fromTime, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parametro 'from' non valido, formato atteso: YYYY-MM-DD"})
			return
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		var err error
		toTime, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parametro 'to' non valido, formato atteso: YYYY-MM-DD"})
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
	role := c.GetString("role")
	slotID := c.Param("id")

	if err := h.service.CancelSlot(c.Request.Context(), userID, role, slotID); err != nil {
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

	booking, err := h.service.BookSlot(c.Request.Context(), userID, req)
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

	bookings, err := h.service.ListMyBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *Handler) ListSlotBookings(c *gin.Context) {
	userID := c.GetString("user_id")
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

	if err := h.service.PatchSlot(c.Request.Context(), userID, role, slotID, req.StartTime, req.EndTime); err != nil {
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
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, booking)
}
