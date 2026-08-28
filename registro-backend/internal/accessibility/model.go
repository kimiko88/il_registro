package accessibility

import "time"

type AccessibilityFeedback struct {
	ID             string     `json:"id"`
	ProtocolNumber string     `json:"protocol_number"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	BarrierType    string     `json:"barrier_type"`
	Description    string     `json:"description"`
	UserID         *string    `json:"user_id,omitempty"`
	SchoolID       *string    `json:"school_id,omitempty"`
	UserAgent      *string    `json:"user_agent,omitempty"`
	IPAddress      *string    `json:"ip_address,omitempty"`
	Status         string     `json:"status"` // 'open', 'in_progress', 'resolved'
	ResponseNotes  *string    `json:"response_notes,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateFeedbackRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	BarrierType string `json:"barrier_type" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type FeedbackResponse struct {
	ID             string `json:"id"`
	ProtocolNumber string `json:"protocol_number"`
	Message        string `json:"message"`
}

type UserAccessibilityPreferences struct {
	UserID    string    `json:"user_id"`
	Settings  string    `json:"settings"` // Raw JSON or JSON string
	UpdatedAt time.Time `json:"updated_at"`
}

type SavePreferencesRequest struct {
	Settings map[string]interface{} `json:"settings" binding:"required"`
}
