package trips

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/trips")
	{
		g.POST("", h.CreateTrip)
		g.GET("", h.ListTrips)
		g.POST("/consent", h.SubmitConsent)
		g.POST("/:id/consent", h.SubmitConsent)
		g.GET("/:id/consents", h.ListConsents)
	}
}

func (h *Handler) CreateTrip(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")

	if userID == "" || (role != "teacher" && role != "admin" && role != "superadmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip, err := h.service.CreateTrip(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, trip)
}

func (h *Handler) ListTrips(c *gin.Context) {
	schoolID := c.GetString("school_id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if schoolID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetStudentID := userID
	if role == "parent" {
		childID := c.Query("student_id")
		if childID != "" {
			targetStudentID = childID
		}
	}
	if role == "teacher" || role == "admin" || role == "superadmin" || role == "secretary" || role == "principal" || role == "vice_principal" {
		targetStudentID = "" // Staff sees all school trips
	}

	trips, err := h.service.ListTrips(c.Request.Context(), schoolID, targetStudentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trips)
}

func (h *Handler) SubmitConsent(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	ipAddress := c.ClientIP()

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "parent" && role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: trip consent submission is restricted to parents and students"})
		return
	}

	var req SubmitConsentRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if req.TripID == "" && c.Param("id") != "" {
		req.TripID = c.Param("id")
	}
	if req.StudentID == "" && role == "student" {
		req.StudentID = userID
	}
	if req.Status == "" {
		req.Status = "Consented"
	}

	if err := h.service.SubmitConsent(c.Request.Context(), userID, role, ipAddress, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trip consent submitted successfully"})
}

func (h *Handler) ListConsents(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" && role != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	tripID := c.Param("id")
	trip, err := h.service.GetTripByID(c.Request.Context(), tripID)
	if err == nil && trip != nil && role != "superadmin" && schoolID != "" && trip.SchoolID != schoolID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: trip belongs to another school"})
		return
	}

	consents, err := h.service.ListConsents(c.Request.Context(), tripID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, consents)
}
