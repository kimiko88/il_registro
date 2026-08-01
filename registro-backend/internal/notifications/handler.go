package notifications

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	n := r.Group("/notifications")
	{
		n.GET("", h.ListNotifications)
		n.PUT("/read-all", h.MarkAllAsRead)
		n.PUT("/:id/read", h.MarkAsRead)
		n.POST("/push-tokens", h.RegisterToken)
		n.POST("/register-device", h.RegisterDevice)
		n.DELETE("/push-tokens", h.UnregisterToken)
		n.POST("/send-push", h.SendPush)
	}

	// PWA Manifest and SW static/configuration endpoints
	r.GET("/manifest.json", h.GetManifest)
	r.GET("/users/me/notifications", h.ListNotifications)
}

func (h *Handler) ListNotifications(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	unreadOnly := c.Query("unread_only") == "true"
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit < 1 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	list, err := h.service.ListDBNotifications(c.Request.Context(), userID, unreadOnly, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.MarkAsRead(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all marked as read"})
}

func (h *Handler) RegisterToken(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req RegisterTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RegisterToken(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "push token registered"})
}

func (h *Handler) RegisterDevice(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		FCMToken string `json:"fcm_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenReq := RegisterTokenRequest{
		DeviceToken: req.FCMToken,
		Platform:    "mobile",
	}

	if err := h.service.RegisterToken(c.Request.Context(), userID, tokenReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "device registered"})
}

func (h *Handler) UnregisterToken(c *gin.Context) {
	userID := c.GetString("user_id")
	deviceToken := c.Query("device_token")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.UnregisterToken(c.Request.Context(), userID, deviceToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "push token unregistered"})
}

func (h *Handler) SendPush(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Non-admin roles (e.g., teachers) cannot send arbitrary notifications to other users without target verification
	if role == "teacher" && req.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "teachers can only send push notifications to themselves or authorized recipients"})
		return
	}

	sent, err := h.service.SendPushNotification(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "notification dispatched", "dispatched_count": sent})
}

func (h *Handler) GetManifest(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetPWAManifest())
}
