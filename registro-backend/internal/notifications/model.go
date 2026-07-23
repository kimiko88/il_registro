package notifications

import (
	"time"
)

type PushToken struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	DeviceToken string    `json:"device_token" db:"device_token"`
	Platform    string    `json:"platform" db:"platform"` // 'android', 'ios', 'web'
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterTokenRequest struct {
	DeviceToken string `json:"device_token" binding:"required"`
	Platform    string `json:"platform"` // 'android', 'ios', 'web'
}

type RegisterDeviceRequest struct {
	FCMToken string `json:"fcm_token" binding:"required"`
}

type SendNotificationRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Body    string `json:"body" binding:"required"`
	Payload map[string]string `json:"payload,omitempty"`
}

type PWAConfig struct {
	Name        string `json:"name"`
	ShortName   string `json:"short_name"`
	StartURL    string `json:"start_url"`
	Display     string `json:"display"`
	ThemeColor  string `json:"theme_color"`
	BGColor     string `json:"background_color"`
}

type DBNotification struct {
	ID        string                 `json:"id" db:"id"`
	UserID    string                 `json:"user_id" db:"user_id"`
	Title     string                 `json:"title" db:"title"`
	Body      string                 `json:"body" db:"body"`
	Type      string                 `json:"type" db:"type"` // 'grade', 'absence', 'homework', 'circular', 'info'
	Payload   map[string]interface{} `json:"payload,omitempty" db:"payload"`
	ReadAt    *time.Time             `json:"read_at,omitempty" db:"read_at"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}
