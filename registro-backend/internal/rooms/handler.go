package rooms

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	roomsGroup := r.Group("/rooms")
	{
		// Buildings (Plessi)
		roomsGroup.GET("/buildings", h.ListBuildings)
		roomsGroup.GET("/buildings/:id", h.GetBuilding)
		roomsGroup.POST("/buildings", h.requireManagementRole, h.CreateBuilding)
		roomsGroup.PUT("/buildings/:id", h.requireManagementRole, h.UpdateBuilding)
		roomsGroup.DELETE("/buildings/:id", h.requireManagementRole, h.DeleteBuilding)

		// Bookings
		roomsGroup.GET("/bookings", h.ListBookings)
		roomsGroup.POST("/bookings", h.CreateBooking)
		roomsGroup.DELETE("/bookings/:id", h.CancelBooking)

		// Rooms
		roomsGroup.GET("", h.ListRooms)
		roomsGroup.GET("/:id", h.GetRoom)
		roomsGroup.GET("/:id/availability", h.GetRoomAvailability)
		roomsGroup.POST("", h.requireManagementRole, h.CreateRoom)
		roomsGroup.PUT("/:id", h.requireManagementRole, h.UpdateRoom)
		roomsGroup.DELETE("/:id", h.requireManagementRole, h.DeleteRoom)
	}
}

func (h *Handler) requireManagementRole(c *gin.Context) {
	role := c.GetString("role")
	switch role {
	case "secretary", "admin", "superadmin", "principal", "vice_principal", "collaboratore_ds":
		c.Next()
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso negato: operazione riservata alla segreteria o alla dirigenza"})
		c.Abort()
	}
}

// ----------------- Buildings Handlers -----------------

func (h *Handler) ListBuildings(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "school context missing"})
		return
	}

	buildings, err := h.service.ListBuildings(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buildings)
}

func (h *Handler) GetBuilding(c *gin.Context) {
	id := c.Param("id")
	building, err := h.service.GetBuilding(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "edificio non trovato"})
		return
	}
	c.JSON(http.StatusOK, building)
}

func (h *Handler) CreateBuilding(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req CreateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	building, err := h.service.CreateBuilding(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, building)
}

func (h *Handler) UpdateBuilding(c *gin.Context) {
	schoolID := c.GetString("school_id")
	id := c.Param("id")
	var req UpdateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	building, err := h.service.UpdateBuilding(c.Request.Context(), schoolID, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, building)
}

func (h *Handler) DeleteBuilding(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteBuilding(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "edificio eliminato con successo"})
}

// ----------------- Rooms Handlers -----------------

func (h *Handler) ListRooms(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "school context missing"})
		return
	}

	var buildingID *string
	if bID := c.Query("building_id"); bID != "" {
		buildingID = &bID
	}

	var roomType *string
	if rt := c.Query("room_type"); rt != "" {
		roomType = &rt
	}

	activeOnly := c.DefaultQuery("active_only", "false") == "true"

	roomsList, err := h.service.ListRooms(c.Request.Context(), schoolID, buildingID, roomType, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roomsList)
}

func (h *Handler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	room, err := h.service.GetRoom(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "aula non trovata"})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (h *Handler) GetRoomAvailability(c *gin.Context) {
	id := c.Param("id")
	fromDate := c.Query("from")
	toDate := c.Query("to")
	if fromDate == "" {
		fromDate = c.Query("date")
	}
	if toDate == "" {
		toDate = fromDate
	}
	if fromDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parametro from o date richiesto (YYYY-MM-DD)"})
		return
	}

	availability, err := h.service.GetRoomAvailability(c.Request.Context(), id, fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, availability)
}

func (h *Handler) CreateRoom(c *gin.Context) {
	schoolID := c.GetString("school_id")
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.service.CreateRoom(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (h *Handler) UpdateRoom(c *gin.Context) {
	schoolID := c.GetString("school_id")
	id := c.Param("id")
	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.service.UpdateRoom(c.Request.Context(), schoolID, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (h *Handler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteRoom(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "aula eliminata con successo"})
}

// ----------------- Bookings Handlers -----------------

func (h *Handler) ListBookings(c *gin.Context) {
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	userID := c.GetString("user_id")

	filter := BookingFilter{
		SchoolID: schoolID,
	}

	if rID := c.Query("room_id"); rID != "" {
		filter.RoomID = &rID
	}
	if bID := c.Query("building_id"); bID != "" {
		filter.BuildingID = &bID
	}
	if cID := c.Query("class_id"); cID != "" {
		filter.ClassID = &cID
	}
	if from := c.Query("from"); from != "" {
		filter.FromDate = &from
	}
	if to := c.Query("to"); to != "" {
		filter.ToDate = &to
	}
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	if parentID := c.Query("parent_booking_id"); parentID != "" {
		filter.ParentBookingID = &parentID
	}

	// Filter by teacher
	if tID := c.Query("teacher_id"); tID != "" {
		filter.TeacherID = &tID
	} else if role == "teacher" && c.Query("all") != "true" {
		// Teachers default to seeing their own bookings unless specified
		filter.TeacherID = &userID
	}

	bookings, err := h.service.ListBookings(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *Handler) CreateBooking(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")

	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	first, all, err := h.service.CreateBooking(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "conflitto") || strings.Contains(err.Error(), "slot already booked") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.IsRecurring {
		c.JSON(http.StatusCreated, gin.H{
			"message":      "prenotazioni ricorrenti create con successo",
			"booking":      first,
			"occurrences":  all,
			"total_booked": len(all),
		})
		return
	}

	c.JSON(http.StatusCreated, first)
}

func (h *Handler) CancelBooking(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	id := c.Param("id")

	cancelSeries := c.DefaultQuery("cancel_series", "false") == "true"

	err := h.service.CancelBooking(c.Request.Context(), userID, role, schoolID, id, cancelSeries)
	if err != nil {
		if err == ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "non puoi cancellare questa prenotazione"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "prenotazione cancellata con successo"})
}
