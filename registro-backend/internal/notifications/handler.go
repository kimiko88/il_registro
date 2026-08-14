package notifications

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
	pushLimiters sync.Map
}

func NewHandler(s Service) *Handler {
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

	// PWA Manifest endpoint
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
		log.Printf("[ERROR] notifications.ListNotifications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
		log.Printf("[ERROR] notifications.MarkAsRead: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
		log.Printf("[ERROR] notifications.MarkAllAsRead: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
		Platform string `json:"platform"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	platform := strings.ToLower(req.Platform)
	if platform != "ios" && platform != "android" && platform != "web" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "platform must be 'ios', 'android', or 'web'"})
		return
	}

	tokenReq := RegisterTokenRequest{
		DeviceToken: req.FCMToken,
		Platform:    platform,
	}

	if err := h.service.RegisterToken(c.Request.Context(), userID, tokenReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "device registered"})
}

func (h *Handler) UnregisterToken(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deviceToken := c.GetHeader("X-Device-Token")
	if deviceToken == "" && c.Request.Body != nil && c.Request.ContentLength > 0 {
		var req struct {
			DeviceToken string `json:"device_token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload: " + err.Error()})
			return
		}
		deviceToken = req.DeviceToken
	}

	if deviceToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_token is required"})
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
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: push notifications restricted to administrative staff"})
		return
	}

	// Rate limit: max 10 requests per minute per admin caller
	limiterVal, _ := h.pushLimiters.LoadOrStore(userID, rate.NewLimiter(rate.Every(6*time.Second), 5))
	limiter := limiterVal.(*rate.Limiter)
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "troppi invii di notifiche push, riprova tra qualche istante"})
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

	sent, err := h.service.SendPushNotification(c.Request.Context(), req)
	if err != nil {
		log.Printf("[ERROR] notifications.SendPush: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "notification dispatched", "dispatched_count": sent})
}

func (h *Handler) GetManifest(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetPWAManifest())
}
